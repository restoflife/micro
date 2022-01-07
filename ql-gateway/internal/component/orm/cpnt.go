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
	"database/sql"
	"fmt"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var (
	ormMgr    = map[string]*gorm.DB{}
	l         dbLog
	newLogger logger.Interface
)

type LogLevel int

const (
	Silent LogLevel = iota + 1
	Error
	Warn
	Info
)

type dbLog struct {
	logger.Writer
	logger.Config
}

//MustBootUp Start database by gorm
func MustBootUp(configs map[string]*conf.ConfigLite, opts ...Option) error {
	options := newOptions(opts...)
	for name, config := range configs {
		l.LogLevel = logger.Error

		if config.ShowSql {
			l.LogLevel = logger.Info
		}

		newLogger = logger.New(
			l,
			logger.Config{
				SlowThreshold:             time.Millisecond * 300,
				LogLevel:                  l.LogLevel,
				IgnoreRecordNotFoundError: true, //忽略ErrRecordNotFound（记录未找到）错误
			},
		)
		db, err := gorm.Open(
			mysql.New(
				mysql.Config{
					DSN: config.Dsn,
				}), &gorm.Config{Logger: newLogger,
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   config.Prefix,
					SingularTable: config.Singular,
				},
				SkipDefaultTransaction: true,
			})
		if err != nil {
			log.Error(zap.Error(err))
			return err
		}

		db = db.Debug()
		DB, err := db.DB()
		if err != nil {
			log.Error(zap.Error(err))
			return err
		}
		if config.MaxIdle > 0 {
			DB.SetMaxIdleConns(config.MaxIdle)
		}
		if config.MaxOpen > 0 {
			DB.SetMaxOpenConns(config.MaxOpen)
		}
		if config.MaxLife > 0 {
			DB.SetConnMaxLifetime(time.Millisecond * time.Duration(config.MaxLife))
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
		ticker := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-ticker.C:
				for _, v := range ormMgr {
					if d, err := v.DB(); err != nil {
						log.Error(zap.Error(err))
					} else {
						if err = d.Ping(); err != nil {
							log.Error(zap.Error(err))
						}
					}
					log.Infox("[gorm]", zap.String("mysql", "PING DATABASE BY GORM"))
				}
			}
		}
	}()
	return nil
}

// ConvertLevel todo abandoned
func ConvertLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "err":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}

	return logger.Silent
}

func (l dbLog) Printf(msg string, data ...interface{}) {
	lg := log.Logger()
	if l.LogLevel >= logger.Info {
		lg.Info("[GORM]", zap.String("info", fmt.Sprintf(msg, data...)))
	}
	if l.LogLevel < logger.Info {
		lg.Error("[GORM]", zap.String("error", fmt.Sprintf(msg, data...)))
	}
}

func NewSession(name string) (*gorm.DB, error) {
	if g, e := get(name); e == nil {
		return g.Session(&gorm.Session{Logger: newLogger}) /*g.WithContext(context.Background())*/, nil
	} else {
		return nil, e
	}
}

func Close(tx *sql.DB) {
	if err := tx.Close(); err != nil {
		log.Error(zap.Error(err))
		return
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
