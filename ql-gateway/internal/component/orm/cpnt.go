/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-01-07 10:49
 * @LastEditors: Administrator
 * @LastEditTime: 2022-01-07 10:49
 * @FilePath: ql-gateway/internal/component/orm/cpnt.go
 */

package orm

import (
	"fmt"
	ormLog "github.com/restoflife/log"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

var ormMgr = map[string]*gorm.DB{}

// MustBootUp Start database by gorm
func MustBootUp(configs map[string]*conf.ConfigLite, opts ...Option) error {
	sqlLog, err := log.NewLogger(conf.C.SQLLogCfg)
	if err != nil {
		return err
	}
	lg := ormLog.NewGormLogger(sqlLog)
	lg.SetAsDefault()
	options := newOptions(opts...)
	for name, config := range configs {
		master := mysql.Open(config.Dsn)
		slaves := make([]gorm.Dialector, len(config.Slave))
		for i, s := range config.Slave {
			slave := mysql.Open(s.Dsn)
			slaves[i] = slave
		}

		db, ex := gorm.Open(master, &gorm.Config{
			Logger: lg,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.Prefix,
				SingularTable: config.Singular,
			},
			SkipDefaultTransaction: true,
			QueryFields:            true,
		})
		if ex != nil {
			return ex
		}
		plugin := dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{master},
			Replicas: slaves,
			Policy:   dbresolver.RandomPolicy{},
		})
		if config.MaxIdle > 0 {
			plugin.SetMaxIdleConns(config.MaxIdle)
		}
		if config.MaxOpen > 0 {
			plugin.SetMaxOpenConns(config.MaxOpen)
		}
		if config.MaxLife > 0 {
			plugin.SetConnMaxLifetime(time.Millisecond * time.Duration(config.MaxLife))
		}
		if config.ShowSql {
			db = db.Debug()
		}
		if err = db.Use(plugin); err != nil {
			return err
		}
		if d, x := db.DB(); x != nil {
			log.Error(zap.Error(x))
			return x
		} else {
			if err = d.Ping(); err != nil {
				log.Error(zap.Error(err))
				return err
			}
		}
		if _, ok := ormMgr[name]; ok {
			return fmt.Errorf("database components loaded twice：[%s]", name)
		}
		if options.syncGorm != nil {
			if err = options.syncGorm(name, db); err != nil {
				return err
			}
		}
		ormMgr[name] = db
	}
	go func() {
		ticker := time.NewTicker(time.Hour * 5)
		for {
			select {
			case <-ticker.C:
				for _, v := range ormMgr {
					if d, x := v.DB(); x == nil {
						if x = d.Ping(); x != nil {
							log.Error(zap.Error(x))
							return
						}
						log.Infox(fmt.Sprintf("%s  %s", "[GORM]", "PING DATABASE mysql"))
					} else {
						log.Error(zap.Error(x))
						return
					}
				}
			}
		}
	}()
	return nil
}

// NewSession 单数据库
func NewSession(name string) (*gorm.DB, error) {
	if g, e := get(name); e == nil {
		return g.Session(&gorm.Session{}), nil
	} else {
		return nil, e
	}
}

// Read 从数据库 只读
func Read(name string) (*gorm.DB, error) {
	if g, e := get(name); e == nil {
		return g.Clauses(dbresolver.Read).Session(&gorm.Session{}), nil
	} else {
		return nil, e
	}
}

// Write 主数据库 只写
func Write(name string) (*gorm.DB, error) {
	if g, e := get(name); e == nil {
		return g.Clauses(dbresolver.Write).Session(&gorm.Session{}), nil
	} else {
		return nil, e
	}
}

// Close 释放连接至连接池
func Close(tx *gorm.DB) {
	if g, e := tx.DB(); e == nil {
		if err := g.Close(); err != nil {
			log.Error(zap.Error(err))
			return
		}
	}
	return
}

func get(name string) (*gorm.DB, error) {
	g, ok := ormMgr[name]
	if !ok {
		log.Error(zap.Error(fmt.Errorf("database does not exist:[%s]", name)))
		return nil, fmt.Errorf("database does not exist:[%s]", name)
	}
	return g, nil
}
