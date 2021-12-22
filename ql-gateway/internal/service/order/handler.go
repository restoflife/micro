/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 14:14
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 14:14
 * @FilePath: ql-gateway/internal/service/order/handler.go
 */

package order

import (
	"context"
	"github.com/restoflife/micro/gateway/internal/component/grpccli/order"
	"github.com/restoflife/micro/gateway/internal/protocol"
	orderPb "github.com/restoflife/micro/protos/order"
)

type API interface {
	OrderDetails(ctx context.Context, req *protocol.GetOrderDetailsReq) (*orderPb.GetOrderDetailsResp, error)
}

type IOrderAPI struct{}

func NewOrderSvc() API {
	return &IOrderAPI{}
}

func (o *IOrderAPI) OrderDetails(ctx context.Context, req *protocol.GetOrderDetailsReq) (*orderPb.GetOrderDetailsResp, error) {
	resp, err := order.ExecHandler(order.GetOrderDetails, &orderPb.GetOrderDetailsReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return resp.(*orderPb.GetOrderDetailsResp), nil
}

func getOrderDetails(id int64) (*orderPb.GetOrderDetailsResp, error) {
	resp, err := order.ExecHandler(order.GetOrderDetails, &orderPb.GetOrderDetailsReq{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.(*orderPb.GetOrderDetailsResp), nil
}
