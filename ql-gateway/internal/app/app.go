/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:52
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:52
 * @FilePath: ql-gateway/internal/app/app.go
 */

package app

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	var consoleEncoder zapcore.Encoder
	consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	logger := zap.New(core)

	zap.ReplaceGlobals(logger)

}

type Application interface {
	Name() string
	Config() *viper.Viper
	InitConfig()
	BootUpPrepare()
	BootUpServer()
	BootUpAfter()
	InitLogger()
	Run()
}

func Run(app Application) {
	app.InitConfig()

	app.InitLogger()

	app.BootUpPrepare()

	// boot up server
	app.BootUpServer()

	app.BootUpAfter()

	// run servers
	app.Run()
}
