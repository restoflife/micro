/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:42
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:42
 * @FilePath: ql-mp/internal/component/db/cpnt.go
 */

package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/restoflife/micro/mp/conf"
	l "github.com/restoflife/micro/mp/internal/component/log"
	"go.uber.org/zap"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var ormMgr = map[string]*xorm.EngineGroup{}

// MustBootUp Start database
func MustBootUp(configs map[string]*conf.ConfigLite, opts ...Option) error {
	options := newOptions(opts...)
	for name, config := range configs {
		master, err := xorm.NewEngine(config.Driver, config.Dsn)
		if err != nil {
			return err
		}
		slaves := make([]*xorm.Engine, len(config.Slave))
		for i, s := range config.Slave {
			slave, err := xorm.NewEngine(config.Driver, s.Dsn)
			if err != nil {
				return err
			}
			slaves[i] = slave
		}

		db, err := xorm.NewEngineGroup(master, slaves)
		if err != nil {
			return err
		}
		db.ShowSQL(config.ShowSql)
		db.Logger().SetLevel(log.LOG_ERR)
		if config.ShowSql {
			db.Logger().SetLevel(log.LOG_INFO)
		}
		if config.MaxIdle > 0 {
			db.SetMaxIdleConns(config.MaxIdle)
		}
		if config.MaxOpen > 0 {
			db.SetMaxOpenConns(config.MaxOpen)
		}
		if _, ok := ormMgr[name]; ok {
			return fmt.Errorf("database components loaded twice：[%s]", name)
		}
		if options.sync2 != nil {
			if err = options.sync2(name, db); err != nil {
				return err
			}
		}
		ormMgr[name] = db
	}
	go func() {
		ticker := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-ticker.C:
				for _, v := range ormMgr {
					if err := v.Ping(); err != nil {
						l.Error(zap.Error(err))
						return
					}
				}
			}
		}
	}()
	return nil
}

func Read(name string) (*xorm.Engine, error) {
	if g, e := get(name); e == nil {
		return g.Slave(), nil
	} else {
		return nil, e
	}
}

func Write(name string) (*xorm.Engine, error) {
	if g, e := get(name); e == nil {
		return g.Master(), nil
	} else {
		return nil, e
	}
}

func NewSession(name string) (*xorm.Session, error) {
	if g, e := get(name); e == nil {
		return g.NewSession(), nil
	} else {
		return nil, e
	}
}
func get(name string) (*xorm.EngineGroup, error) {
	g, ok := ormMgr[name]
	if !ok {
		l.Error(zap.Error(fmt.Errorf("database does not exist:[%s]", name)))
		return nil, fmt.Errorf("database does not exist:[%s]", name)
	}
	return g, nil
}

func Close(session *xorm.Session) {
	err := session.Close()
	if err != nil {
		l.Error(zap.Error(err))
		return
	}
	return
}
