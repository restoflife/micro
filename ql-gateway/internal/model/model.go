/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:40
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:40
 * @FilePath: ql-gateway/internal/model/model.go
 */

package model

import (
	"github.com/restoflife/micro/gateway/internal/constant"
	"time"
	"xorm.io/xorm"
)

type TimeStruct struct {
	CreatedAt time.Time `json:"created_at" xorm:"'created_at' created comment('创建时间') DATETIME"`
	UpdatedAt time.Time `json:"updated_at" xorm:"'updated_at' updated comment('最后更新时间') DATETIME"`
	DeletedAt time.Time `json:"deleted_at" xorm:"'deleted_at' deleted comment('删除时间') DATETIME"`
}

type Account struct {
	Id         int       `json:"id" xorm:"'id' not null pk autoincr INT"`
	Uid        uint64    `json:"uid" xorm:"'uid' not null comment('uid 唯一') unique BIGINT"`
	Account    string    `json:"account" xorm:"'account' not null comment('账号 唯一') unique VARCHAR(125)"`
	Password   string    `json:"password" xorm:"'password' not null comment('密码') VARCHAR(255)"`
	Username   string    `json:"username" xorm:"'username' not null comment('昵称') VARCHAR(125)"`
	Ip         string    `json:"ip" xorm:"'ip' not null comment('登陆ip') VARCHAR(50)"`
	Time       time.Time `json:"time" xorm:"'time' not null comment('登陆时间') DATETIME"`
	Status     string    `json:"status" xorm:"'status' not null default 'active' comment('active启用,inactive禁用') VARCHAR(50)"`
	Avatar     string    `json:"avatar" xorm:"'avatar' not null comment('头像') VARCHAR(255)"`
	TimeStruct `xorm:"extends"`
}

func (m *Account) TableName() string {
	return "account"
}
func Sync(name string, group *xorm.EngineGroup) error {
	switch name {
	case constant.DbDefaultName:
		return group.StoreEngine("InnoDB").Sync2(
			new(Account),
		)
	}

	return nil
}
