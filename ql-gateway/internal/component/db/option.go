/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021/12/6 13:55
 * @LastEditors: Administrator
 * @LastEditTime: 2021/12/6 13:55
 * @FilePath: internal/component/orm/option.go
 */

package db

import (
	"xorm.io/xorm"
)

type SyncXormFunc func(string, *xorm.EngineGroup) error

type Options struct {
	syncXorm SyncXormFunc
}

type Option func(*Options)

func SetSyncXormFunc(f SyncXormFunc) Option {
	return func(o *Options) {
		o.syncXorm = f
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		syncXorm: nil,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
