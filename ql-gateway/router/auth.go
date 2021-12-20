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
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/middleware"
	"github.com/restoflife/micro/gateway/internal/service/auth"
	"go.uber.org/zap"
)

var (
	userPath       = "/passport/"
	adminPath      = "/user/"
	authMiddleware *middleware.GinJWTMiddleware
	err            error
)

func authGroup(root *gin.RouterGroup) {
	authApi := root.Group(userPath).Use()
	authMiddleware, err = middleware.AuthInit()
	if err != nil {
		log.Panic(zap.Error(err))
	}
	//登陆
	authApi.POST("/login", authMiddleware.LoginHandler)
	//验证码
	authApi.GET("/captcha", auth.MakeCaptchaHandler)
	//注册
	authApi.POST("/register", auth.MakeRegisterHandler)
}

func adminGroup(root *gin.RouterGroup) {
	adminApi := root.Group(adminPath).Use(authMiddleware.MiddlewareFunc())
	//管理员列表
	adminApi.GET("/list", auth.MakeUserListHandler)
}
