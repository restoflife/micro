/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 11:31
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 11:31
 * @FilePath: ql-order/transport/order/order.go
 */

package order

import (
	"context"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	"github.com/restoflife/micro/order/internal/encoding"
	orderEndpoint "github.com/restoflife/micro/order/internal/endpoint/order"
	"github.com/restoflife/micro/order/internal/service/order"
	orderPb "github.com/restoflife/micro/protos/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderSvc struct {
	getOrderDetails grpcTransport.Handler
}

func (o *orderSvc) GetOrderDetails(ctx context.Context, req *orderPb.GetOrderDetailsReq) (*orderPb.GetOrderDetailsResp, error) {
	_, resp, err := o.getOrderDetails.ServeGRPC(ctx, req)
	errObj := encoding.Wrap(err)
	if errObj != nil {
		return nil, status.Error(codes.Code(errObj.Code), errObj.Msg)
	}

	return resp.(*orderPb.GetOrderDetailsResp), err
}

func NewUserServer(opts ...grpcTransport.ServerOption) orderPb.OrderSvcServer {
	return &orderSvc{
		getOrderDetails: grpcTransport.NewServer(
			orderEndpoint.MakeOrderDetailsEndpoint(order.NewAuthSvc()),
			decodeOrderDetailsRequest,
			encodeOrderDetailsResponse,
			opts...,
		),
	}
}
