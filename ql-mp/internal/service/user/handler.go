/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:18
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:18
 * @FilePath: ql-mp/internal/service/user/handler.go
 */

package user

import (
	"context"
	"fmt"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/protocol"
	"github.com/restoflife/micro/mp/utils"
)

type PassportAPI interface {
	login(ctx context.Context, req *protocol.MpLoginReq) (*protocol.MpLoginResp, error)
}

type IPassportAPI struct{}

func (I *IPassportAPI) login(ctx context.Context, req *protocol.MpLoginReq) (*protocol.MpLoginResp, error) {
	uri := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", "wx156a410597adc5e5", "86435cebea493a2e141dc258a62f2cfa", req.Code)
	_, data := utils.Get(uri)
	fmt.Println(string(data))
	//TODO implement me
	log.Infox(req.Code)
	return nil, nil
}

func NewOrderSvc() PassportAPI {
	return &IPassportAPI{}
}
