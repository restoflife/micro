/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 11:33
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 11:33
 * @FilePath: ql-order/transport/order/auth.go
 */

package order

import (
	"context"
	"fmt"
	orderPb "github.com/restoflife/micro/protos/order"
)

func decodeOrderDetailsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*orderPb.GetOrderDetailsReq)
	if !ok {
		return nil, fmt.Errorf("grpc server decode request出错！")
	}
	return req, nil
}
func encodeOrderDetailsResponse(c context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*orderPb.GetOrderDetailsResp)
	if !ok {
		return nil, fmt.Errorf("grpc server encode response error (%T)", response)
	}
	return resp, nil
}
