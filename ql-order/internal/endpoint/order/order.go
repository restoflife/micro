/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 11:29
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 11:29
 * @FilePath: ql-order/internal/endpoint/order/order.go
 */

package order

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/restoflife/micro/order/internal/errutil"
	"github.com/restoflife/micro/order/internal/service/order"
	orderPb "github.com/restoflife/micro/protos/order"
)

func MakeOrderDetailsEndpoint(svc order.IAuthAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*orderPb.GetOrderDetailsReq)
		if !ok {
			return nil, errutil.ErrEndpointType
		}
		resp, err := svc.OrderDetails(ctx, req)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
