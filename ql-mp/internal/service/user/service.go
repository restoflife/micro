/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:18
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:18
 * @FilePath: ql-mp/internal/service/user/service.go
 */

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/encoding"
	"github.com/restoflife/micro/mp/internal/errutil"
	"github.com/restoflife/micro/mp/internal/protocol"
	"go.uber.org/zap"
)

func MakeLoginHandler(svc PassportAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &protocol.MpLoginReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
		resp, err := svc.login(c, req)
		if err != nil {
			log.Error(zap.Error(err))
			encoding.ErrorMsg(c, errutil.ErrInternalServer)
			return
		}
		encoding.Ok(c, resp)
		c.Next()
	}
}
