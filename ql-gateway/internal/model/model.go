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
	CreatedAt time.Time `json:"created_at" xorm:"'created_at' created comment('创建时间') TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" xorm:"'updated_at' updated comment('最后更新时间') TIMESTAMP"`
	DeletedAt time.Time `json:"deleted_at" xorm:"'deleted_at' deleted comment('删除时间') TIMESTAMP"`
}

func Sync(name string, group *xorm.EngineGroup) error {
	switch name {
	case constant.DbDefaultName:
		return group.StoreEngine("InnoDB").Sync2()
	}

	return nil
}
