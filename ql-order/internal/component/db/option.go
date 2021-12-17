/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021/12/6 13:55
 * @LastEditors: Administrator
 * @LastEditTime: 2021/12/6 13:55
 * @FilePath: internal/component/orm/option.go
 */

package db

import "xorm.io/xorm"

type Sync2Func func(string, *xorm.EngineGroup) error

type Options struct {
	sync2 Sync2Func
}

type Option func(*Options)

func SetSync2Func(f Sync2Func) Option {
	return func(o *Options) {
		o.sync2 = f
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		sync2: nil,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
