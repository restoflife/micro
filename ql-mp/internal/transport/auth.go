/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 15:58
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 15:58
 * @FilePath: ql-mp/internal/transport/auth.go
 */

package transport

import (
	"context"
	"fmt"
	user_pb "github.com/restoflife/micro/protos/mp"
)

func decodeUserListRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*user_pb.GetUserListReq)
	if !ok {
		return nil, fmt.Errorf("grpc server decode request出错！")
	}
	return req, nil
}
func encodeUserListResponse(c context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*user_pb.GetUserListResp)
	if !ok {
		return nil, fmt.Errorf("grpc server encode response error (%T)", response)
	}
	return resp, nil
}
