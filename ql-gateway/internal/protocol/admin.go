/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-20 09:39
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-20 09:39
 * @FilePath: ql-gateway/internal/protocol/admin.go
 */

package protocol

import (
	"github.com/golang-jwt/jwt"
	"github.com/mojocn/base64Captcha"
)

type RegisterReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Ip       string `json:"ip" xorm:"-"`
	UID      uint64 `json:"uid"`
}

type ConfigJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

type CaptchaResp struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Otp      string `json:"otp"`
	UID      string `json:"uid"`
	Ip       string `json:"ip" xorm:"-"`
}
type LoginResp struct {
	Uid        string `json:"uid"`
	Account    string `json:"account"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Time       string `json:"time"`
	Token      string `json:"token"`
	ExpireTime string `json:"expire_time"`
}

type TokenResp struct {
	Token      string `json:"token"`
	ExpireTime string `json:"expire_time"`
}
type TokenClaims struct {
	Account string `json:"account"`
	Ip      string `json:"ip"`
	Uid     string `json:"uid"`
	jwt.StandardClaims
}
