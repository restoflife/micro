/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:10
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:10
 * @FilePath: ql-gateway/router/api.go
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"net/http"
)

var (
	rootPath = "/api/"
)

func ApiRouter(root *gin.Engine) {
	root.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "")
		encoding.Error(c, errutil.ErrPageNotFound)
		return
	})

	api := root.Group(rootPath)

	// Authentic route
	authGroup(api)
	// administrators route
	adminGroup(api)
	// Order route
	orderGroup(api)
	// mp route
	mpApiGroup(api)

}
