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
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
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
		ZapLogger:                 zapLogger,
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
	l.logger().Sugar().Debugf(str, args...)
}

func (l Logger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l Logger) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.logger().Error("[GORM]", zap.Error(err), zap.String("latency", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.logger().Warn("[GORM]", zap.String("latency", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= logger.Info:
		sql, rows := fc()
		l.logger().Debug("[GORM]", zap.String("latency", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
	ormPackage  = filepath.Join("ql-gateway", "internal", "component", "orm")
)

func (l Logger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, ormPackage):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
