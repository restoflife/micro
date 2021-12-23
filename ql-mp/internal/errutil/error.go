/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 14:17
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 14:17
 * @FilePath: ql-mp/internal/errutil/error.go
 */

package errutil

import "time"

type Error struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Time    int64       `json:"time"`
	Content string      `json:"content,omitempty"`
	Stack   string      `json:"stack,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
}

// Error implements error.
func (e Error) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "unknown error"
}
func New(code int, msg string) Error {
	return Error{
		Code:   code,
		Msg:    msg,
		Time:   time.Now().Unix(),
		Status: "fail",
		Data:   "",
	}
}

// 系统错误
const (
	internalServer = iota + 1000
	rpcInternalServer
	rpcRequest
	ErrorWithMessage
)

var (
	ErrInternalServer    = New(internalServer, "internal server error")
	ErrRpcInternalServer = New(rpcInternalServer, "rpc client connection failed")
	ErrRpcRequest        = New(rpcRequest, "rpc request failed")
)

// 业务错误
const (
	Unknown = iota + 2000
	illegalParameter
	unauthorized
	accountExist
)

var (
	ErrUnauthorized     = New(unauthorized, "unauthorized")
	ErrIllegalParameter = New(illegalParameter, "illegal parameter")
	ErrAccountExist     = New(accountExist, "account exist")
)
