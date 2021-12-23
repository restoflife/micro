/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 10:57
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 10:57
 * @FilePath: ql-order/internal/app/base.go
 */

package app

import (
	"github.com/restoflife/micro/order/internal/component/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Base struct {
	AppName string
	Command *cobra.Command
	conf    *viper.Viper
}

func New(name string, cmd *cobra.Command) *Base {
	return &Base{
		AppName: name,
		Command: cmd,
	}
}
func (b *Base) Name() string {
	return b.AppName
}
func (b *Base) Config() *viper.Viper {
	return b.conf
}
func (b *Base) InitConfig() {}

func (b *Base) BootUpPrepare() {}

func (b *Base) BootUpServer() {}

func (b *Base) BootUpAfter() {}

func (b *Base) InitLogger() {
	log.Init()
	defer log.Sync()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(zap.Error(err))
	}
	zap.ReplaceGlobals(logger)
}

func (b *Base) Run() {}
