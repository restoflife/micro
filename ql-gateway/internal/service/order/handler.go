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
	"github.com/restoflife/micro/gateway/internal/component/grpccli/order"
	orderPb "github.com/restoflife/micro/protos/order"
)

func getOrderDetails(id int64) (*orderPb.GetOrderDetailsResp, error) {
	resp, err := order.ExecHandler(order.GetOrderDetails, &orderPb.GetOrderDetailsReq{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.(*orderPb.GetOrderDetailsResp), nil
}
