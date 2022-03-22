/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-20 13:43
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-20 13:43
 * @FilePath: ql-gateway/internal/middleware/auth.go
 */

package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/db"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/model"
	"github.com/restoflife/micro/gateway/internal/model/auth"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func AuthInit() (*GinJWTMiddleware, error) {
	return New(&GinJWTMiddleware{
		Realm:      "kit",
		Key:        []byte(conf.C.Login.Key),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		// 验证用户
		Authenticator: Authenticator,
		// 退出登录
		Unauthorized: Unauthorized,
		PayloadFunc:  PayloadFunc,
		// 验证token
		Authorizator:    Authorization,
		IdentityHandler: IdentityHandler,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
}

func Authenticator(c *gin.Context) (interface{}, error) {
	req := &protocol.LoginReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(zap.Error(err))
		return nil, errutil.ErrIllegalParameter
	}
	session, err := db.NewSession(constant.DbDefaultName)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, errutil.ErrInternalServer
	}
	defer db.Close(session)

	// if !utils.VerifyCode(req.UID, req.Otp) {
	//	return nil, errutil.ErrVerificationCode
	// }

	// key := fmt.Sprintf("%s%s", req.Account, constant.UserLoginTotal)
	// total, _ := redis.GetCacheByFloat(key)
	// if int(total) >= conf.C.Login.Total {
	//	relieveTime, _ := redis.GetRedisExpTime(key)
	//	log.Error(zap.Error(errors.New(fmt.Sprintf("%s错误次数已达上限[%s]后解除锁定", req.Account, relieveTime))))
	//	return nil, errors.New(fmt.Sprintf("错误次数已达上限[%s]后解除锁定", relieveTime))
	// }
	req.Ip = c.ClientIP()
	u, e := auth.NewAuthModel(session).LoginModel(req)
	if e == nil {
		var mp = make(map[string]interface{})
		mp["user"] = u
		return mp, nil
	} else {
		// _ = redis.SetCache(key, int(total)+1, time.Duration(conf.C.Login.Time)*time.Second)
		return nil, e
	}
}
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}

func PayloadFunc(data interface{}) MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var u model.Account
		// analysis interface to struct
		if err := mapstructure.Decode(v["user"], &u); err != nil {
			log.Error(zap.Any("analysis interface to struct fail", zap.Error(err)))
			return MapClaims{}
		}
		return MapClaims{
			IdentityKey: u.Id,
			UsernameKey: u.Username,
			UIDKey:      u.Uid,
		}
	}
	return MapClaims{}
}

func Authorization(data interface{}, c *gin.Context) bool {
	type User struct {
		IdentityKey int    `json:"IdentityKey"`
		Uid         int    `json:"uid"`
		UserName    string `json:"UserName"`
	}
	var u User
	if err := mapstructure.Decode(data, &u); err != nil {
		log.Error(zap.Error(errors.New("analysis interface to struct fail")))
		return false
	}
	c.Set("uuid", u.Uid)
	c.Set("userName", u.UserName)
	c.Set("id", u.IdentityKey)
	return true

}

func IdentityHandler(c *gin.Context) interface{} {
	claims := ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["username"],
		"Uid":         claims["uid"],
	}
}
