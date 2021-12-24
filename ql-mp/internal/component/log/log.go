/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:57
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:57
 * @FilePath: ql-mp/internal/component/log/log.go
 */

package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/internal/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"time"
)

// Logger 全局日志对象
var logger *zap.Logger

type ErrorHandler struct {
	log *zap.Logger
}

func NewZapLogErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		log: logger,
	}
}

func (h *ErrorHandler) Handle(ctx context.Context, err error) {
	h.log.Error(
		fmt.Sprintf("[%s]--%s", ctx.Value(constant.ContextOrderKey), ctx.Value(constant.ContextMpUUid)), zap.Error(err),
	)
}

func Init() {
	l, err := NewLogger(conf.C.RunLogCfg)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	logger = l
}

func Logger() *zap.Logger {
	return logger
}

func NewLogger(logCfg *conf.LogConfig) (*zap.Logger, error) {

	encoder := createFileEncoder()
	consoleEncoder := createConsoleEncoder()
	cores := make([]zapcore.Core, 0)
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.ErrorLevel
	})
	cores = append(
		cores,
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   logCfg.Filename,   //日志文件路径
				MaxSize:    logCfg.MaxSize,    //每个日志文件保存的大小 单位:M
				MaxBackups: logCfg.MaxBackups, //日志文件最多保存多少个备份
				MaxAge:     logCfg.MaxAge,     //文件最多保存多少天
			}),
			createLevelEnablerFunc(logCfg.Level),
		),
		zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stderr),
			debugPriority,
		),
	)
	return zap.New(zapcore.NewTee(cores...)), nil
}

//日志控制台配置
func createConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	//日志时间格式
	encoderConfig.EncodeTime = timeEncoder
	//将级别序列化为全大写字符串并添加颜色
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//日志文件配置
func createFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	//日志时间格式
	encoderConfig.EncodeTime = timeEncoder
	//将级别序列化为全大写字符串
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(constant.Layout))
}

func createLevelEnablerFunc(input string) zap.LevelEnablerFunc {
	var lv = new(zapcore.Level)
	if err := lv.UnmarshalText([]byte(input)); err != nil {
		return nil
	}
	return func(lev zapcore.Level) bool {
		return lev >= *lv
	}
}

// Debug 调试日志
func Debug(f ...zapcore.Field) {
	logger.Debug("[DEBUG]", f...)
}

// Debugx 调试日志
func Debugx(msg string, f ...zapcore.Field) {
	logger.Debug(msg, f...)
}

// Error 错误日志
func Error(f ...zapcore.Field) {
	_, file, line, _ := runtime.Caller(1)
	f = append(f, zap.String("func", fmt.Sprintf("%s:%d", file, line)))
	logger.Error("[ERROR]", f...)
}

// Err 错误日志
func Err(f ...zapcore.Field) {
	_, file, line, _ := runtime.Caller(2)
	f = append(f, zap.String("func", fmt.Sprintf("%s:%d", file, line)))
	logger.Error("[ERROR]", f...)
}

// Info 信息日志
func Info(f ...zapcore.Field) {
	logger.Info("[INFO]", f...)
}

// Infox 信息日志
func Infox(msg string, f ...zapcore.Field) {
	logger.Info(msg, f...)
}

func DumpJson(v interface{}) {
	data, _ := json.Marshal(v)
	Info(zap.String("request data", string(data)))
}

func Fatal(f ...zapcore.Field) {
	logger.Fatal("[FATAL]", f...)
}

func Panic(f ...zapcore.Field) {
	_, file, line, _ := runtime.Caller(1)
	f = append(f, zap.String("func", fmt.Sprintf("%s:%d", file, line)))
	logger.Panic("[PANIC]", f...)
}
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}
