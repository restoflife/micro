/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 15:28
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 15:28
 * @FilePath: ql-gateway/router/order.go
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/service/order"
)

var orderPath = "/order/"

func orderGroup(root *gin.RouterGroup) {
	orderAPI := root.Group(orderPath)
	//订单详情
	orderAPI.GET("/detail", order.GetOrderDetails)
}
