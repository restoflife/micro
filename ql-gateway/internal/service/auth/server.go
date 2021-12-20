/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:40
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:40
 * @FilePath: ql-gateway/internal/service/auth/server.go
 */

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
)

//func MakeLoginHandler(c *gin.Context) {
//	req := &protocol.LoginReq{}
//	if err := c.ShouldBindJSON(req); err != nil {
//		encoding.Error(c, errutil.ErrIllegalParameter)
//		return
//	}
//	req.Ip = c.ClientIP()
//	resp, err := makeLoginService(req)
//	if err != nil {
//		encoding.Error(c, err)
//		return
//	}
//	encoding.Ok(c, resp)
//}

func MakeCaptchaHandler(c *gin.Context) {
	id, b64s, err := utils.DriverDigitFunc()
	if err != nil {
		log.Error(zap.Error(err))
		return
	}
	resp := &protocol.CaptchaResp{
		Id:  id,
		Url: b64s,
	}
	encoding.Ok(c, resp)

}

func MakeRegisterHandler(c *gin.Context) {
	req := &protocol.RegisterReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		encoding.Error(c, errutil.ErrIllegalParameter)
		return
	}
	req.Ip = c.ClientIP()
	uid, _ := utils.GetUUID()
	req.UID = uid
	if err := makeRegisterService(req); err != nil {
		encoding.Error(c, err)
		return
	}
	encoding.Ok(c, "")

}
func MakeUserListHandler(c *gin.Context) {}
