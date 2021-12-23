/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 14:14
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 14:14
 * @FilePath: ql-gateway/internal/service/order/server.go
 */

package order

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/protocol"
)

func MakeOrderDetailsHandler(svc API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &protocol.GetOrderDetailsReq{}
		if err := c.ShouldBind(req); err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
		resp, err := svc.OrderDetails(c, req)
		if err != nil {
			encoding.ErrorWithGRPC(c, err)
			return
		}
		encoding.Ok(c, resp)
		c.Next()
	}
}

//func MakeOrderDetailsHandler(c *gin.Context) {
//	req := &protocol.GetOrderDetailsReq{}
//	if err := c.ShouldBind(req); err != nil {
//		encoding.Error(c, errutil.ErrIllegalParameter)
//		return
//	}
//	resp, err := getOrderDetails(req.Id)
//	// todo: check if
//	//resp, err := redis.CheckCache(fmt.Sprintf("id:%d", req.Id), func() (interface{}, error) {
//	//	return getOrderDetails(req.Id)
//	//}, time.Duration(60), true)
//	if err != nil {
//		encoding.ErrorWithGRPC(c, err)
//		return
//	}
//	encoding.Ok(c, resp)
//	c.Next()
//}
