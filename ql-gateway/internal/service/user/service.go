/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:27
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:27
 * @FilePath: ql-gateway/internal/service/user/service.go
 */

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/restoflife/micro/gateway/utils"
)

func MakeUserListHandler(svc API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &protocol.GetMpUserListReq{}
		if err := c.ShouldBind(req); err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
		limit, offset := utils.PageIndex(int(req.Page), int(req.PageSize))
		req.Page = int32(limit)
		req.PageSize = int32(offset)
		resp, err := svc.mpUserList(c, req)
		if err != nil {
			encoding.ErrorWithGRPC(c, err)
			return
		}
		encoding.Ok(c, resp)
		c.Next()
	}
}
