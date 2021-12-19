/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:30
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:30
 * @FilePath: ql-gateway/router/auth.go
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/service/auth"
)

var (
	userPath    = "/user"
	notAuthPath = "/passport"
)

func authGroup(root *gin.RouterGroup) {
	authApi := root.Group(notAuthPath)

	//登陆
	authApi.POST("/login", auth.Login)
	//注册
	authApi.POST("/register", auth.Register)
	//验证码
	authApi.GET("/check", auth.Register)
}

func adminGroup(root *gin.RouterGroup) {
	userApi := root.Group(userPath)
	//管理员
	userApi.GET("/list", auth.Login)
}
