/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-27 14:38
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-27 14:38
 * @FilePath: ql-gateway/internal/service/order/rpc.go
 */

package order

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/restoflife/micro/gateway/internal/component/grpccli"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/errutil"
	orderPb "github.com/restoflife/micro/protos/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
)

func getOrderDetails(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*orderPb.GetOrderDetailsReq)
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

		svc := orderPb.NewOrderSvcClient(conn)
		ctx, cancel, uuid := grpccli.GrpcClientCtx(constant.ContextOrderUUid, 5)
		defer cancel()
		r, err := svc.GetOrderDetails(ctx, req)
		if err != nil {
			log.Error(zap.String("uuid", uuid), zap.Error(err))
			return nil, errutil.ErrRpcRequest
		}
		return r, nil
	}, nil, nil
}
