/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:27
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:27
 * @FilePath: ql-gateway/internal/service/user/handler.go
 */

package user

import (
	"context"
	"github.com/restoflife/micro/gateway/internal/component/grpccli"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/protocol"
	userPb "github.com/restoflife/micro/protos/mp"
)

type API interface {
	mpUserList(ctx context.Context, req *protocol.GetMpUserListReq) (*protocol.CommonListResp, error)
}

type IUserAPI struct{}

func NewUserSvc() API {
	return &IUserAPI{}
}

func (I *IUserAPI) mpUserList(ctx context.Context, r *protocol.GetMpUserListReq) (*protocol.CommonListResp, error) {
	resp, err := grpccli.ExecHandler(grpccli.InstancedMgr[constant.MpPrefix], getUserList, &userPb.GetUserListReq{
		Page:     r.Page,
		PageSize: r.PageSize,
		Uid:      uint64(r.Uid),
		Nickname: r.Nickname,
	})
	if err != nil {
		return nil, err
	}
	result := resp.(*userPb.GetUserListResp)

	return &protocol.CommonListResp{Total: result.Total, List: result.List}, nil
}
