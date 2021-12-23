/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 14:17
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 14:17
 * @FilePath: ql-mp/internal/encoding/http.go
 */

package encoding

import (
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/errutil"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type encResponse struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
	Time   int64       `json:"time"`
}

// Ok gin
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, encResponse{
		Code:   http.StatusOK,
		Data:   data,
		Status: "success",
		Time:   time.Now().Unix(),
	})
}

func Error(c *gin.Context, err error) {
	_ = c.Error(err)
	var e = errutil.ErrInternalServer
	if sc, ok := err.(errutil.Error); ok {
		e = sc
	}

	Code := http.StatusOK
	if err == errutil.ErrUnauthorized {
		Code = http.StatusUnauthorized
	}
	log.Err(zap.String("[GIN-PATH]", c.Request.URL.Path), zap.Error(err))

	c.AbortWithStatusJSON(Code, e)
}

func ErrorMsg(c *gin.Context, err error) {
	log.Err(zap.String("[GIN-PATH]", c.Request.URL.Path), zap.Error(err))
	c.AbortWithStatusJSON(http.StatusOK, errutil.New(errutil.ErrorWithMessage, err.Error()))
}

func ErrorWithGRPC(c *gin.Context, err error) {
	_ = c.Error(err)
	err = gRPCErrorConvert(err)
	var e = errutil.ErrInternalServer
	if sc, ok := err.(errutil.Error); ok {
		e = sc
	}
	//log.Err(zap.String("[GRPC-PATH]", c.Request.URL.Path), zap.Error(err))
	statusCode := http.StatusOK
	if err == errutil.ErrUnauthorized {
		statusCode = http.StatusUnauthorized
	}

	c.AbortWithStatusJSON(statusCode, e)
}

func gRPCErrorConvert(err error) error {
	if err == nil {
		return nil
	}

	if sc, ok := status.FromError(err); ok {
		switch sc.Code() {
		case codes.Unauthenticated:
			err = errutil.ErrUnauthorized
		default:
			err = errutil.New(errutil.Unknown, sc.Message())
		}
	}

	return err
}

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
