/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-01-07 10:49
 * @LastEditors: Administrator
 * @LastEditTime: 2022-01-07 10:49
 * @FilePath: ql-gateway/internal/component/orm/option.go
 */

package orm

import (
	"gorm.io/gorm"
)

type SyncGormFunc func(string, *gorm.DB) error

type Options struct {
	syncGorm SyncGormFunc
}

type Option func(*Options)

func SetSyncGormFunc(f SyncGormFunc) Option {
	return func(o *Options) {
		o.syncGorm = f
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		syncGorm: nil,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
