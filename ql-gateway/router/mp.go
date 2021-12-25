/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:24
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:24
 * @FilePath: ql-gateway/router/mp.go
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/service/user"
)

var miniPath = "/mini/"

func mpApiGroup(root *gin.RouterGroup) {
	mpApi := root.Group(miniPath) //.Use(authMiddleware.MiddlewareFunc())
	//小程序用户列表
	mpApi.GET("/list", user.MakeUserListHandler(user.NewUserSvc()))
}
