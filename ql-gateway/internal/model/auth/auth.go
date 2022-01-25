/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:41
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:41
 * @FilePath: ql-gateway/internal/model/auth/auth.go
 */

package auth

import (
	"errors"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/model"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Model struct {
	session *xorm.Session
}

func NewAuthModel(session *xorm.Session) *Model {
	return &Model{session: session}
}

func (m *Model) RegisterModel(r *protocol.RegisterReq) error {
	exist, err := m.session.Where(builder.Eq{"account": r.Account}).Exist(&model.Account{})
	if err != nil {
		return err
	}
	if exist {
		log.Debug(zap.String(r.Account, "Already exists"))
		return errutil.ErrAccountExist
	}
	pwd, err := utils.EncryptionPassword(r.Password)
	if err != nil {
		return err
	}
	_, err = m.session.Omit("status").Insert(&model.Account{
		Uid:      r.UID,
		Account:  r.Account,
		Password: pwd,
		Ip:       r.Ip,
		Time:     time.Now(),
		Avatar:   r.Avatar,
		Username: r.Username,
	})
	if err != nil {
		return err
	}
	//m.session.Where(builder.Between{Col: "id", LessVal: 1, MoreVal: 10}).Sum(&model.Account{}, "id")
	return nil
}

func (m *Model) LoginModel(r *protocol.LoginReq) (*model.Account, error) {
	u := &model.Account{}
	_, err := m.session.Where(builder.Eq{"account": r.Account}).Get(u)
	if err != nil {
		return nil, err
	}
	hash, err := utils.CompareHashAndPassword(u.Password, r.Password)
	if err != nil || !hash {
		return nil, errors.New("password mismatch")
	}

	if u.Status != "active" {
		return nil, errors.New("账号状态异常")
	}
	_, err = m.session.Where(builder.Eq{"account": r.Account}).
		Update(&model.Account{
			Ip:   r.Ip,
			Time: time.Now(),
		})

	if err != nil {
		log.Debug(zap.String(r.Account, "Log login failed"))
	}
	return u, nil
}
