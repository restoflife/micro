/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:30
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:30
 * @FilePath: ql-gateway/router/server.go
 */

package router

import "github.com/gin-gonic/gin"

func RegisterPOSTHandler(r *gin.RouterGroup, path string, handlers ...gin.HandlerFunc) {
	r.POST(path, handlers...)
}

func RegisterGETHandler(r *gin.RouterGroup, path string, handlers ...gin.HandlerFunc) {
	r.GET(path, handlers...)
}

func RegisterPUTHandler(r *gin.RouterGroup, path string, handlers ...gin.HandlerFunc) {
	r.PUT(path, handlers...)
}
func RegisterDELETEHandler(r *gin.RouterGroup, path string, handlers ...gin.HandlerFunc) {
	r.DELETE(path, handlers...)
}
