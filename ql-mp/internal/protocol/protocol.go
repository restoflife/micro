/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:39
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:39
 * @FilePath: ql-mp/internal/protocol/protocol.go
 */

package protocol

type (
	CommonListResp struct {
		Total int64       `json:"total"`
		List  interface{} `json:"list"`
	}
	MpLoginReq struct {
		Code          string `json:"code" binding:"required"`
		EncryptedData string `json:"encrypted_data" binding:"required"`
		Iv            string `json:"iv" binding:"required"`
	}
	MpLoginResp struct {
		Uid      string `json:"uid"`
		Openid   string `json:"openid"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Token    string `json:"token"`
	}
	WxSessionResp struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		Unionid    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errMsg"`
	}
	WxMPSensitivityData struct {
		Openid     string `json:"openid"`
		Nickname   string `json:"nickName"`
		Gender     int32  `json:"gender"`
		City       string `json:"city"`
		Province   string `json:"province"`
		Country    string `json:"country"`
		AvatarUrl  string `json:"avatarUrl"`
		UnionId    string `json:"unionId"`
		Language   string `json:"language"`
		SessionKey string `json:"session_key"`
	}
)
