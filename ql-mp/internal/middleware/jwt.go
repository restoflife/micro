/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 14:24
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 14:24
 * @FilePath: ql-mp/internal/middleware/jwt.go
 */

package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/component/redis"
	"github.com/restoflife/micro/mp/internal/constant"
	"github.com/restoflife/micro/mp/internal/encoding"
	"github.com/restoflife/micro/mp/internal/errutil"
	"github.com/restoflife/micro/mp/utils"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

var (
	//key  = []byte(conf.C.Jwt.Key)
	pool = sync.Pool{
		New: func() interface{} {
			return &userMetaEx{}
		},
	}
)

const (
	jwtOK    = 0
	jwtGuest = -1
)

type userMeta struct {
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
}

type userMetaEx struct {
	userMeta
	Token string `json:"token"`
}

//SignJWT 签发JWT
func SignJWT(uid, nickname string) (string, error) {
	var token *jwt.Token
	newUid, _ := strconv.ParseInt(uid, 10, 64)
	uidHash, err := utils.HashIDEncode(newUid)
	if err != nil {
		return "", nil
	}
	mEx := pool.Get().(*userMetaEx)
	defer pool.Put(mEx)
	_, err = redis.CheckCache(constant.RedisPrefixToken+uid, func() (interface{}, error) {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"uid": uidHash,           // hash混淆的uid
			"iat": time.Now().Unix(), // 签发时间
		})

		tok, err := token.SignedString([]byte(conf.C.Jwt.Key))
		if err != nil {
			return "", err
		}
		mEx.Uid = uid
		mEx.Nickname = nickname
		mEx.Token = tok
		return mEx, nil
	}, time.Duration(conf.C.Jwt.TimeOut), true)
	if err != nil {
		return "", err
	}
	return mEx.Token, nil
}

//CheckJWT 检查JWT是否合法
func CheckJWT(c *gin.Context) bool {

	result := checkJWTHelper(c)
	log.Debug(zap.Int8("检查JWT是否合法", result))
	if result == jwtOK {
		return true
	}

	return false

}

//验证token
func checkJWTHelper(c *gin.Context) int8 {
	tokenVal := c.GetHeader("Authorization")
	log.Debugx("CheckJWT", zap.String("token", tokenVal))
	if tokenVal == "" {
		return jwtGuest
	}

	token, err := jwt.Parse(tokenVal, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(conf.C.Jwt.Key), nil
	})

	if err == nil && token.Valid {
		uidHash, ok := token.Claims.(jwt.MapClaims)["uid"].(string)
		if !ok {
			return -2
		}

		uid, err := utils.HashIDDecode(uidHash)
		if err != nil {
			log.Error(zap.Error(err))
			return -3
		}

		b, err := redis.GetCache(constant.RedisPrefixToken + strconv.FormatInt(uid, 10))
		if err != nil {
			log.Error(zap.Error(err))
			return -1
		}
		mEx := pool.Get().(*userMetaEx)
		defer pool.Put(mEx)
		mEx.Token = ""
		if err = mapstructure.Decode(b, &mEx); err != nil {
			log.Error(zap.Error(err))
			return -5
		}
		if mEx.Token != tokenVal {
			return -6
		}

		c.Set("meta", &userMeta{
			Uid:      mEx.Uid,
			Nickname: mEx.Nickname,
		})
		return jwtOK

	} else if err != nil {
		log.Debugx("解析token", zap.Error(err))
	}
	return -7
}

func JWTValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckJWT(c) {
			encoding.Error(c, errutil.ErrUnauthorized)
			return
		}
	}
}
