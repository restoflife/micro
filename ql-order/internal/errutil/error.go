/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 15:46
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 15:46
 * @FilePath: ql-order/internal/errutil/error.go
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
		Code: code,
		Msg:  msg,
		Time: time.Now().Unix(),
	}
}

//业务错误
const (
	Unknown      = iota + 2000
	unauthorized //未授权
	endpointType
)

var (
	ErrEndpointType = New(endpointType, "endpoint request type error")
	ErrUnauthorized = New(unauthorized, "unauthorized")
)
