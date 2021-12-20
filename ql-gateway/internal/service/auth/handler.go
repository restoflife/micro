/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:40
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:40
 * @FilePath: ql-gateway/internal/service/auth/handler.go
 */

package auth

import (
	"github.com/restoflife/micro/gateway/internal/component/db"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/model/auth"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"go.uber.org/zap"
)

func makeRegisterService(r *protocol.RegisterReq) error {
	session, err := db.NewSession(constant.DbDefaultName)
	if err != nil {
		return err
	}
	defer db.Close(session)
	if err = auth.NewAuthModel(session).RegisterModel(r); err != nil {
		log.Error(zap.Error(err))
		return err
	}
	return nil
}

//func makeLoginService(r *protocol.LoginReq) (*protocol.LoginResp, error) {
//	session, err := db.NewSession(constant.DbDefaultName)
//	if err != nil {
//		return nil, err
//	}
//	defer db.Close(session)
//	resp, err := auth.NewAuthModel(session).LoginModel(r)
//	if err != nil {
//		log.Error(zap.Error(err))
//		return nil, err
//	}
//	//签发token
//	token, err := middleware.GenerateToken(&protocol.TokenClaims{
//		Account: resp.Account,
//		Ip:      r.Ip,
//		Uid:     resp.Uid,
//	})
//	if err != nil {
//		log.Error(zap.Error(err))
//		return nil, err
//	}
//	resp.Token = token.Token
//	resp.ExpireTime = token.ExpireTime
//	return resp, nil
//}
