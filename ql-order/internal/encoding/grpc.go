/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 15:46
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 15:46
 * @FilePath: ql-order/internal/encoding/grpc.go
 */

package encoding

import "github.com/restoflife/micro/order/internal/errutil"

func Wrap(err error) *errutil.Error {

	if err == nil {
		return nil
	}

	errObj, _ := err.(*errutil.Error)
	if errObj != nil {
		return errObj
	}

	return &errutil.Error{
		Code: -1,
		Msg:  "系统错误",
	}
}
