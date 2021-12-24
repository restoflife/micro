/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 11:05
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 11:05
 * @FilePath: ql-mp/internal/model/user.go
 */

package model

import (
	"github.com/restoflife/micro/mp/internal/protocol"
	user_pb "github.com/restoflife/micro/protos/mp"
	"strconv"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type UserModel struct {
	session *xorm.Session
}

func NewUserModel(session *xorm.Session) *UserModel {
	return &UserModel{session: session}
}
func (m *UserModel) GetUserList(r *user_pb.GetUserListReq) (*user_pb.GetUserListResp, error) {
	u := make([]MpUser, 0)
	if r.Uid > 0 {
		m.session.Where(builder.Eq{"uid": r.Uid})
	}
	total, err := m.session.Limit(int(r.Page), int(r.PageSize)).FindAndCount(&u)
	if err != nil {
		return nil, err
	}
	item := make([]*user_pb.GetUserListItem, total)

	for k, v := range u {
		item[k] = &user_pb.GetUserListItem{
			Uid:      strconv.FormatInt(v.Uid, 10),
			Avatar:   v.Avatar,
			Nickname: v.Nickname,
		}
	}
	return &user_pb.GetUserListResp{
		Total: total,
		List:  item,
	}, nil
}

func (m *UserModel) GetUserByThird(Openid string) (*MpUser, error) {
	u := MpUser{}
	exist, err := m.session.Where(builder.Eq{"openid": Openid}).Get(&u)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	return &u, nil
}

func (m *UserModel) UpdateUserByUid(uid int64, bean map[string]interface{}) error {
	_, err := m.session.Table(&MpUser{}).Where(builder.Eq{"uid": uid}).Update(bean)
	return err
}

func (m *UserModel) InsertUser(uid uint64, r *protocol.WxMPSensitivityData) (*MpUser, error) {
	u := &MpUser{
		Uid:        int64(uid),
		Phone:      "",
		Openid:     r.Openid,
		Nickname:   r.Nickname,
		Avatar:     r.AvatarUrl,
		Sex:        int(r.Gender),
		City:       r.City,
		Language:   r.Language,
		Province:   r.Province,
		Country:    r.Country,
		SessionKey: r.SessionKey,
	}
	_, err := m.session.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
