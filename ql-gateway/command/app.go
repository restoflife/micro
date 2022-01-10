/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:48
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:48
 * @FilePath: ql-gateway/command/app.go
 */

package command

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/app"
	"github.com/restoflife/micro/gateway/internal/component/db"
	"github.com/restoflife/micro/gateway/internal/component/elasticsearch"
	"github.com/restoflife/micro/gateway/internal/component/grpccli"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/component/orm"
	"github.com/restoflife/micro/gateway/internal/component/redis"
	"github.com/restoflife/micro/gateway/internal/model"
	"github.com/restoflife/micro/gateway/router"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type mainApp struct {
	*app.Base
}

var srv *http.Server

func NewApp(name string, cmd *cobra.Command) *mainApp {
	return &mainApp{
		Base: &app.Base{
			AppName: name,
			Command: cmd,
		},
	}
}

func (m *mainApp) InitConfig() {
	confFile := m.Command.Flags().Lookup("c")
	if confFile == nil {
		panic("There is no configuration file" + m.Name())
	}
	if _, err := toml.DecodeFile(confFile.Value.String(), &conf.C); err != nil {
		panic(err)
	}
}

func (m *mainApp) BootUpPrepare() {

	log.Infox("initialize connection to redis...")
	if err := redis.MustBootUp(conf.C.Redis); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("grpc client initialized...")
	if err := grpccli.MustBootUp(); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("elasticsearch client initialized...")
	if err := elasticsearch.NewElasticSearchClient(conf.C.Elastic); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("initialize xorm connection to database....")
	//TODO ::db.SetSyncXormFunc(model.SyncXorm) 生产环境不建议开启
	if err := db.MustBootUp(conf.C.DB, db.SetSyncXormFunc(model.SyncXorm)); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("initialize gorm connection to database....")
	//TODO ::orm.SetSyncGormFunc(model.SyncGorm) 生产环境不建议开启
	if err := orm.MustBootUp(conf.C.DB, orm.SetSyncGormFunc(model.SyncGorm)); err != nil {
		log.Panic(zap.Error(err))
	}

}
func (m *mainApp) BootUpServer() {
	go httpServer()
}

func (m *mainApp) Run() {
	f := func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Error(zap.Any("Server Shutdown:", zap.Error(err)))
				return err
			}
			return fmt.Errorf("received signal %s", sig)
		}
	}
	log.Infox("Terminated", zap.Error(f()))
}

func httpServer() {
	if !conf.C.ServerCfg.Mode {
		gin.SetMode(gin.ReleaseMode)
	}

	logger, err := log.NewLogger(conf.C.AccessLogCfg)
	if err != nil {
		return
	}
	handler := gin.New()
	handler.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowWebSockets:  true, //Allow webSocket
	}))

	handler.Use(Logger(logger), Recovery(log.Logger()))
	//pprof
	pprof.Register(handler)
	//Load API route
	router.ApiRouter(handler)

	log.Infox("listening",
		zap.String("transport", "HTTP"),
		zap.String("address", conf.C.ServerCfg.Addr),
	)

	srv = &http.Server{
		Addr:           conf.C.ServerCfg.Addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic(zap.Error(err))
	}
}
