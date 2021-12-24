/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:18
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:18
 * @FilePath: ql-mp/internal/model/model.go
 */

package model

import (
	"github.com/restoflife/micro/mp/internal/constant"
	"time"
	"xorm.io/xorm"
)

type TimeStruct struct {
	CreatedAt time.Time `json:"created_at" xorm:"'created_at' created comment('创建时间') DATETIME"`
	UpdatedAt time.Time `json:"updated_at" xorm:"'updated_at' updated comment('最后更新时间') DATETIME"`
	DeletedAt time.Time `json:"deleted_at" xorm:"'deleted_at' deleted comment('删除时间') DATETIME"`
}
type MpUser struct {
	Id         int    `json:"id" xorm:"'id' not null pk autoincr INT"`
	Uid        int64  `json:"uid" xorm:"'uid' not null comment('平台唯一id') index BIGINT"`
	Phone      string `json:"phone" xorm:"'phone' not null comment('电话') VARCHAR(50)"`
	Openid     string `json:"openid" xorm:"'openid' not null comment('用户的标识，对当前公众号唯一') index VARCHAR(30)"`
	Nickname   string `json:"nickname" xorm:"'nickname' not null comment('用户的昵称') VARCHAR(64)"`
	Avatar     string `json:"avatar" xorm:"'avatar' not null comment('头像') VARCHAR(255)"`
	Sex        int    `json:"sex" xorm:"'sex' not null default 0 comment('用户的性别，值为1时是男性，值为2时是女性，值为0时是未知') TINYINT(1)"`
	City       string `json:"city" xorm:"'city' not null comment('用户所在城市') VARCHAR(64)"`
	Language   string `json:"language" xorm:"'language' not null comment('用户的语言，简体中文为zh_CN') VARCHAR(64)"`
	Province   string `json:"province" xorm:"'province' not null comment('用户所在省份') VARCHAR(64)"`
	Country    string `json:"country" xorm:"'country' not null comment('用户所在国家') VARCHAR(64)"`
	SessionKey string `json:"session_key" xorm:"'session_key' comment('小程序用户会话密匙') VARCHAR(32)"`
	TimeStruct `xorm:"extends"`
}

func (*MpUser) TableName() string {
	return "mp_user"
}
func Sync(name string, group *xorm.EngineGroup) error {
	switch name {
	case constant.DbDefaultName:
		return group.StoreEngine("InnoDB").Sync2(
			new(MpUser),
		)
	}

	return nil
}
