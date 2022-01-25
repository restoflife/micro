/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-01-10 17:40
 * @LastEditors: Administrator
 * @LastEditTime: 2022-01-10 17:40
 * @FilePath: ql-gateway/internal/component/orm/logger.go
 */

package orm

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func New(zapLogger *zap.Logger) Logger {
	return Logger{
		ZapLogger:                 zapLogger.Named("[GORM]"),
		LogLevel:                  logger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

func (l Logger) SetAsDefault() {
	logger.Default = l
}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return Logger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.ZapLogger.Sugar().Info(fmt.Sprintf(str, args...))
}

func (l Logger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.ZapLogger.Sugar().Warnf(fmt.Sprintf(str, args...))
}

func (l Logger) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.ZapLogger.Sugar().Error(fmt.Sprintf(str, args...))
}

func (l Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	lg := l.ZapLogger
	var level zapcore.Level
	sql, rows := fc()
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		level = zapcore.ErrorLevel
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		level = zapcore.WarnLevel
	case l.LogLevel >= logger.Info:
		level = zapcore.DebugLevel
	}
	lg.Check(level, "").Write(zap.String("sql", sql), zap.Int64("rows", rows), zap.String("latency", elapsed.String()), zap.Error(err))
}
