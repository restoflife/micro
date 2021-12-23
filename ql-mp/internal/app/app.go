/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:16
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:16
 * @FilePath: ql-mp/internal/app/app.go
 */

package app

import "github.com/spf13/viper"

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
