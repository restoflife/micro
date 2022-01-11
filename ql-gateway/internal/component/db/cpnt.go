/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:42
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:42
 * @FilePath: ql-gateway/internal/component/db/cpnt.go
 */

package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	ormLog "github.com/restoflife/log"
	"github.com/restoflife/micro/gateway/conf"
	l "github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
	"time"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var dbMgr = map[string]*xorm.EngineGroup{}

// MustBootUp Start database by xorm
func MustBootUp(configs map[string]*conf.ConfigLite, opts ...Option) error {
	options := newOptions(opts...)
	sqlLog, err := l.NewLogger(conf.C.SQLLogCfg)
	if err != nil {
		return err
	}
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
		db.Logger().SetLevel(xlog.LOG_ERR)
		db.SetLogger(ormLog.NewXormLogger(sqlLog))
		db.ShowSQL(config.ShowSql)
		if config.ShowSql {
			db.Logger().SetLevel(xlog.LOG_INFO)
		}
		if config.MaxIdle > 0 {
			db.SetMaxIdleConns(config.MaxIdle)
		}
		if config.MaxOpen > 0 {
			db.SetMaxOpenConns(config.MaxOpen)
		}
		if config.MaxLife > 0 {
			db.SetConnMaxLifetime(time.Millisecond * time.Duration(config.MaxLife))
		}
		if _, ok := dbMgr[name]; ok {
			return fmt.Errorf("database components loaded twice：[%s]", name)
		}
		if options.syncXorm != nil {
			if err = options.syncXorm(name, db); err != nil {
				return err
			}
		}
		dbMgr[name] = db
	}
	go func() {
		ticker := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-ticker.C:
				for _, v := range dbMgr {
					if err := v.Ping(); err != nil {
						l.Error(zap.Error(err))
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
	g, ok := dbMgr[name]
	if !ok {
		l.Error(zap.Error(fmt.Errorf("database does not exist:[%s]", name)))
		return nil, fmt.Errorf("database does not exist:[%s]", name)
	}
	return g, nil
}

func Close(session *xorm.Session) {
	if err := session.Close(); err != nil {
		l.Error(zap.Error(err))
		return
	}
	return
}
