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
	"github.com/restoflife/micro/gateway/internal/encoding"
)

var (
	userPath = "/passport/"
)

func authGroup(root *gin.RouterGroup) {
	authApi := root.Group(userPath)

	RegisterAuthAPIHandler(authApi)
}

func RegisterAuthAPIHandler(r *gin.RouterGroup) {

	RegisterPOSTHandler(r, "/login", func(c *gin.Context) {
		encoding.Ok(c, userPath)
		c.Next()
	})
}
