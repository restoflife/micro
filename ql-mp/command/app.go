/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:16
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:16
 * @FilePath: ql-mp/command/app.go
 */

package command

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/internal/app"
	"github.com/restoflife/micro/mp/internal/component/db"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/component/redis"
	"github.com/restoflife/micro/mp/internal/model"
	"github.com/restoflife/micro/mp/router"
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
	log.Infox("initialize xorm connection to database....")
	if err := db.MustBootUp(conf.C.DB, db.SetSync2Func(model.Sync)); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("initialize connection to redis...")
	if err := redis.MustBootUp(conf.C.Redis); err != nil {
		log.Panic(zap.Error(err))
	}
}
func (m *mainApp) BootUpServer() {
	go httpServer()
	//go gRPC()
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

	//Load API route
	router.API(handler)

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