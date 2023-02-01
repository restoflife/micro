/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-01-17 14:14
 * @LastEditors: Administrator
 * @LastEditTime: 2022-01-17 14:14
 * @FilePath: ql-gateway/internal/component/mongo/cpnt.go
 */

package mongo

import (
	"context"
	"fmt"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

var (
	mongoMgr = map[string]*mongo.Client{}
	ctx      = context.Background()
)

func MustBootUp(configs map[string]*conf.MongoConfig) error {
	for name, config := range configs {
		opts := options.Client()
		opts.ApplyURI(config.Addr)
		// opts.SetAuth(options.Credential{
		//	//AuthMechanism:           "MONGODB-CR",
		//	//AuthMechanismProperties: nil,
		//	AuthSource:  "admin",
		//	Username:    config.Username,
		//	Password:    config.Password,
		//	PasswordSet: true,
		// })
		if config.MaxIdleTime > 0 {
			opts.SetMaxConnIdleTime(time.Duration(config.MaxIdleTime) * time.Second)
		}
		if config.MaxPool > 0 {
			opts.SetMaxPoolSize(config.MaxPool)
		}
		if config.MinPool > 0 {
			opts.SetMinPoolSize(config.MinPool)
		}
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			return err
		}
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			return err
		}
		// 列出所以数据库
		// client.ListDatabaseNames()
		log.Infox(fmt.Sprintf("mongodb database：%s", name))
		mongoMgr[name] = client
	}
	defer func() {
		for _, m := range mongoMgr {
			log.Error(zap.Error(m.Disconnect(ctx)))
		}
	}()
	return nil
}

func NewDatabase(name string) (*mongo.Database, error) {
	if g, e := get(name); e == nil {
		return g, nil
	} else {
		return nil, e
	}
}

func get(name string) (*mongo.Database, error) {
	g, ok := mongoMgr[name]
	if !ok {
		return nil, fmt.Errorf("the mongodb database does not exist：[%s]", name)
	}
	return g.Database(name), nil
}
