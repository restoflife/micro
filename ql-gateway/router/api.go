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
	"net/http"
)

var (
	rootPath = "/api/"
)

func ApiRouter(root *gin.Engine) {
	root.NoRoute(func(c *gin.Context) { c.String(http.StatusNotFound, "") })
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
