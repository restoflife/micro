/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021/12/6 9:47
 * @LastEditors: Administrator
 * @LastEditTime: 2021/12/6 9:47
 * @FilePath: internal/middleware/logger.go
 */

package command

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"time"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

//MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return LoggerWithWriter(logger, gin.DefaultWriter)
}

// LoggerWithWriter instance a Logger middleware with the specified writer buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(logger *zap.Logger, out io.Writer, notlogged ...string) gin.HandlerFunc {
	return LoggerWithConfig(logger, LoggerConfig{
		Output:    out,
		SkipPaths: notlogged,
	})
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(logger *zap.Logger, conf LoggerConfig) gin.HandlerFunc {
	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	netlogo := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(netlogo); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range netlogo {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}
			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			if param.ErrorMessage == "" {
				Ip := utils.ClientIp(c.Request)
				logger.Info("GIN",
					zap.String("Path", path),
					zap.Int("Code", param.StatusCode),
					zap.String("Method", param.Method),
					zap.String("ClientIP", Ip),
					zap.String("User-Agent", c.Request.UserAgent()),
					zap.Any("Latency", fmt.Sprintf("%dms", param.Latency.Milliseconds())),
				)
			} else {
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			}
		}
	}
}

// Recovery 使用zap替换gin内部的recovery模块
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Error("[Recovery From Panic]",
					zap.String("Time", time.Now().Format(constant.Layout)),
					zap.Any("Error", err),
					zap.String("Request", string(httpRequest)),
					zap.String("Stack", string(debug.Stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
