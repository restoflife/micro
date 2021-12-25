/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:29
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:29
 * @FilePath: ql-mp/router/router.go
 */

package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/mp/internal/middleware"
	"github.com/restoflife/micro/mp/internal/service/user"
	"net/http"
)

var (
	needTokenPath = "/api/v1/"
	noTokenPath   = "/api/"
)

func API(root *gin.Engine) {
	root.NoRoute(func(c *gin.Context) { c.String(http.StatusNotFound, "") })
	//不要token
	noTokenApi := root.Group(noTokenPath)
	{
		noTokenApi.POST("/login", user.MakeLoginHandler(user.NewUserSvc()))
	}
	tokenApi := root.Group(needTokenPath).Use(middleware.JWTValidationMiddleware())
	{
		tokenApi.GET("/x", func(context *gin.Context) {
			fmt.Println(context)
		})
	}
}
