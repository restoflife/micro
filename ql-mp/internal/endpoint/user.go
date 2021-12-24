/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:04
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:04
 * @FilePath: ql-mp/internal/endpoint/user.go
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/restoflife/micro/mp/internal/errutil"
	"github.com/restoflife/micro/mp/internal/service/user"
	user_pb "github.com/restoflife/micro/protos/mp"
)

func MakeUserListHandler(svc user.PassportAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*user_pb.GetUserListReq)
		if !ok {
			return nil, errutil.ErrEndpointType
		}
		return svc.UserList(ctx, req)
	}
}
