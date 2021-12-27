/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-27 14:39
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-27 14:39
 * @FilePath: ql-gateway/internal/service/user/rpc.go
 */

package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/restoflife/micro/gateway/internal/component/grpccli"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/errutil"
	userPb "github.com/restoflife/micro/protos/mp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
)

func getUserList(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*userPb.GetUserListReq)
		if !ok {
			return nil, errutil.ErrIllegalParameter
		}

		conn, err := grpccli.GrpcClient(instanceAddr)
		if err != nil {
			return nil, errutil.ErrRpcInternalServer
		}
		defer func(conn *grpc.ClientConn) {
			err = conn.Close()
			if err != nil {
				log.Error(zap.Error(err))
			}
		}(conn)

		svc := userPb.NewUserSvcClient(conn)

		ctx, cancel, uuid := grpccli.GrpcClientCtx(constant.ContextMpUUid, 5)
		defer cancel()

		r, err := svc.GetUserList(ctx, req)
		if err != nil {
			log.Error(zap.String("uuid", uuid), zap.Error(err))
			return nil, errutil.ErrRpcRequest
		}
		return r, nil
	}, nil, nil
}
