/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 14:00
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 14:00
 * @FilePath: ql-gateway/internal/component/grpccli/init.go
 */

package grpccli

import (
	"github.com/restoflife/micro/gateway/internal/component/grpccli/order"
	"github.com/restoflife/micro/gateway/internal/component/grpccli/user"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
)

func MustSetup() error {
	//订单中心客户端
	if err := order.InitOrder(); err != nil {
		log.Error(zap.Error(err))
		return err
	}
	//小程序
	if err := user.InitUser(); err != nil {
		log.Error(zap.Error(err))
		return err
	}
	return nil
}
