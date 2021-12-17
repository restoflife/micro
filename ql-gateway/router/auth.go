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
	"net/http"
	"time"
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
		time.Sleep(time.Second * 3)
		c.JSON(http.StatusOK, gin.H{"name": userPath})
	})
}
