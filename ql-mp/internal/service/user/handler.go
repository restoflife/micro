/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:18
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:18
 * @FilePath: ql-mp/internal/service/user/handler.go
 */

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/internal/component/db"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/component/redis"
	"github.com/restoflife/micro/mp/internal/constant"
	"github.com/restoflife/micro/mp/internal/middleware"
	"github.com/restoflife/micro/mp/internal/model"
	"github.com/restoflife/micro/mp/internal/protocol"
	"github.com/restoflife/micro/mp/utils"
	user_pb "github.com/restoflife/micro/protos/mp"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type PassportAPI interface {
	login(ctx context.Context, req *protocol.MpLoginReq) (*protocol.MpLoginResp, error)
	UserList(ctx context.Context, req *user_pb.GetUserListReq) (*user_pb.GetUserListResp, error)
}

type IPassportAPI struct{}

func NewUserSvc() PassportAPI {
	return &IPassportAPI{}
}

func (I *IPassportAPI) login(_ context.Context, req *protocol.MpLoginReq) (*protocol.MpLoginResp, error) {
	appid := conf.C.Wechat.APPID
	secret := conf.C.Wechat.SECRET
	url := fmt.Sprintf(constant.UrlCode2Session, appid, secret, req.Code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Error(zap.Error(err))
		}
	}(resp.Body)

	ret := protocol.WxSessionResp{}
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	if ret.ErrCode != 0 {
		return nil, fmt.Errorf(ret.ErrMsg)
	}

	wxUserInfo, err := wxMPEncryptedDataDecode(appid, req.EncryptedData, req.Iv, ret.SessionKey)
	if err != nil {
		return nil, err
	}

	wxUserInfo.Openid = ret.Openid
	wxUserInfo.SessionKey = ret.SessionKey

	session, err := db.NewSession(constant.DbDefaultName)
	if err != nil {
		return nil, err
	}
	defer db.Close(session)
	//用户信息
	uInfo, err := model.NewUserModel(session).GetUserByThird(wxUserInfo.Openid)
	if err != nil {
		return nil, err
	}
	result := protocol.MpLoginResp{
		Openid: ret.Openid,
	}
	//修改
	if uInfo != nil {
		bean := map[string]interface{}{}
		if uInfo.Nickname != wxUserInfo.Nickname {
			bean["nickname"] = wxUserInfo.Nickname
			uInfo.Nickname = wxUserInfo.Nickname
		}
		if uInfo.Avatar != wxUserInfo.AvatarUrl {
			bean["avatar"] = wxUserInfo.AvatarUrl
		}
		if len(bean) > 0 {
			err = model.NewUserModel(session).UpdateUserByUid(uInfo.Uid, bean)
			if err != nil {
				log.Error(
					zap.Int64("uid", uInfo.Uid),
					zap.Any("userInfo", uInfo),
					zap.Error(err))
				return nil, err
			}
		}
		result.Uid = strconv.FormatInt(uInfo.Uid, 10)
		result.Nickname = uInfo.Nickname
		result.Avatar = uInfo.Avatar
	} else {
		//新增
		newUid, err := utils.GetUUID()
		if err != nil {
			return nil, err
		}
		u, err := model.NewUserModel(session).InsertUser(newUid, wxUserInfo)
		if err != nil {
			return nil, err
		}
		result.Uid = strconv.FormatInt(u.Uid, 10)
		result.Nickname = u.Nickname
		result.Avatar = u.Avatar
		//记录session_key
		_, err = redis.CheckCache(constant.RedisKeyWxSessionKey+result.Uid, func() (interface{}, error) {
			return ret.SessionKey, nil
		}, time.Duration(-1), false)
		if err != nil {
			log.Error(zap.String("uuid", result.Uid),
				zap.String("session", ret.SessionKey),
				zap.Error(err))
		}
	}
	// 签发token
	token, err := middleware.SignJWT(result.Uid, result.Nickname)
	if err != nil {
		return nil, err
	}
	result.Token = token
	return &result, nil
}

func (I *IPassportAPI) UserList(ctx context.Context, req *user_pb.GetUserListReq) (*user_pb.GetUserListResp, error) {
	session, err := db.NewSession(constant.DbDefaultName)
	if err != nil {
		return nil, err
	}
	defer db.Close(session)
	return model.NewUserModel(session).GetUserList(req)

}

func wxMPEncryptedDataDecode(appId, encryptedData, iv, sessionKey string) (*protocol.WxMPSensitivityData, error) {
	//解密密文
	pc := utils.WxBizDataCrypt{
		AppId:      appId,
		SessionKey: sessionKey,
	}
	a, err := pc.Decrypt(encryptedData, iv, true)
	if err != nil {
		return nil, err
	}
	data := protocol.WxMPSensitivityData{}
	if err = json.Unmarshal([]byte(a.(string)), &data); err != nil {
		return nil, err
	}

	return &data, err
}
