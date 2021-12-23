/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 11:20
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 11:20
 * @FilePath: ql-order/internal/service/order/server.go
 */

package order

import (
	"context"
	orderPb "github.com/restoflife/micro/protos/order"
)

type AuthAPI struct{}

func NewAuthSvc() IAuthAPI {
	return &AuthAPI{}
}

type IAuthAPI interface {
	OrderDetails(ctx context.Context, req *orderPb.GetOrderDetailsReq) (resp *orderPb.GetOrderDetailsResp, err error)
}

func (a *AuthAPI) OrderDetails(ctx context.Context, req *orderPb.GetOrderDetailsReq) (resp *orderPb.GetOrderDetailsResp, err error) {
	return getOrderDetails(req.Id)
}
