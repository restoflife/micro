/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 15:58
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 15:58
 * @FilePath: ql-mp/internal/transport/user.go
 */

package transport

import (
	"context"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/restoflife/micro/mp/internal/encoding"
	userEndpoint "github.com/restoflife/micro/mp/internal/endpoint"
	"github.com/restoflife/micro/mp/internal/service/user"
	user_pb "github.com/restoflife/micro/protos/mp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userSvc struct {
	getUserList grpcTransport.Handler
}

func NewUserServer(opts ...grpcTransport.ServerOption) user_pb.UserSvcServer {
	return &userSvc{
		getUserList: grpcTransport.NewServer(
			userEndpoint.MakeUserListHandler(user.NewUserSvc()),
			decodeUserListRequest,
			encodeUserListResponse,
			opts...,
		),
	}
}

func (u userSvc) GetUserList(ctx context.Context, req *user_pb.GetUserListReq) (*user_pb.GetUserListResp, error) {
	_, resp, err := u.getUserList.ServeGRPC(ctx, req)
	errObj := encoding.Wrap(err)
	if errObj != nil {
		return nil, status.Error(codes.Code(errObj.Code), errObj.Msg)
	}

	return resp.(*user_pb.GetUserListResp), err
}
