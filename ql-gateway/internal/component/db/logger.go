/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-01-11 11:29
 * @LastEditors: Administrator
 * @LastEditTime: 2022-01-11 11:29
 * @FilePath: ql-gateway/internal/component/db/logger.go
 */

package db

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"xorm.io/xorm/log"
)

type Logger struct {
	logger *zap.Logger
	off    bool
	show   bool
	level  log.LogLevel
}

func New(zl *zap.Logger) *Logger {
	return &Logger{
		logger: zl.Named("[XORM]"),
		off:    false,
		show:   true,
	}
}

func (zl *Logger) BeforeSQL(_ log.LogContext) {}

func (zl *Logger) AfterSQL(ctx log.LogContext) {
	sql := fmt.Sprintf("%v %v", ctx.SQL, ctx.Args)
	var level zapcore.Level
	if ctx.Err != nil {
		level = zapcore.ErrorLevel
	} else {
		level = zapcore.DebugLevel
	}
	lg := zl.logger
	lg.Check(level, "").Write(zap.String("sql", sql), zap.String("latency", ctx.ExecuteTime.String()), zap.Error(ctx.Err))
}

func (zl *Logger) Debugf(format string, v ...interface{}) {
	zl.logger.Debug(fmt.Sprintf(format, v...))
}

func (zl *Logger) Infof(format string, v ...interface{}) {
	zl.logger.Info(fmt.Sprintf(format, v...))
}

func (zl *Logger) Warnf(format string, v ...interface{}) {
	zl.logger.Warn(fmt.Sprintf(format, v...))
}

func (zl *Logger) Errorf(format string, v ...interface{}) {
	zl.logger.Error(fmt.Sprintf(format, v...))
}

func (zl *Logger) Level() log.LogLevel {
	if zl.off {
		return log.LOG_OFF
	}

	for _, l := range []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel} {
		if zl.logger.Core().Enabled(l) {
			switch l {
			case zapcore.DebugLevel:
				return log.LOG_DEBUG

			case zapcore.InfoLevel:
				return log.LOG_INFO

			case zapcore.WarnLevel:
				return log.LOG_WARNING

			case zapcore.ErrorLevel:
				return log.LOG_ERR
			}
		}
	}
	return log.LOG_UNKNOWN
}

func (zl *Logger) SetLevel(l log.LogLevel) {
	zl.level = l
}

func (zl *Logger) ShowSQL(b ...bool) {
	zl.show = b[0]
}
func (zl *Logger) IsShowSQL() bool {
	return zl.show
}
