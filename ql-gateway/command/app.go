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
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/app"
	"github.com/restoflife/micro/gateway/internal/component/db"
	"github.com/restoflife/micro/gateway/internal/component/grpccli"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/component/redis"
	"github.com/restoflife/micro/gateway/internal/model"
	"github.com/restoflife/micro/gateway/router"
	"github.com/restoflife/micro/gateway/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
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
	log.Infox("initialize xorm connection to database")
	if err := db.MustBootUp(conf.C.DB, db.SetSync2Func(model.Sync)); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("initialize connection to redis")
	if err := redis.MustBootUp(conf.C.Redis); err != nil {
		log.Panic(zap.Error(err))
	}

	log.Infox("grpc client initialized...")
	if err := grpccli.MustSetup(); err != nil {
		log.Error(zap.Error(err))
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
		ReadTimeout:    600 * time.Second,
		WriteTimeout:   600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic(zap.Error(err))
	}
}

func gRPC() {
	//ETCD connection parameters
	option := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
		Cert:          conf.C.ServerCfg.EtcdCert,
		Key:           conf.C.ServerCfg.EtcdKey,
		CACert:        conf.C.ServerCfg.EtcdCaCert,
	}
	addr, err := utils.GetUrls(conf.C.ServerCfg.Etcd)
	if err != nil {
		log.Panic(zap.Error(err))
	}

	//Create a connection
	client, err := etcdv3.NewClient(context.Background(), addr, option)
	if err != nil {
		log.Panic(zap.Error(err))
	}
	//Create a registration
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
		Key:   conf.C.ServerCfg.Prefix,
		Value: conf.C.ServerCfg.RPCAddr,
	}, kitLog.NewNopLogger())

	//Start registration service
	registrar.Register()

	lis, err := net.Listen("tcp", conf.C.ServerCfg.RPCAddr)
	if err != nil {
		log.Panic(zap.Error(err))
	}
	opts := []grpc.ServerOption{
		grpcMiddleware.WithUnaryServerChain(
			UnaryServerInterceptor,
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time: time.Second * 10,
			})}

	log.Infox("listening",
		zap.String("transport", "GRPC"),
		zap.String("address", conf.C.ServerCfg.RPCAddr),
		zap.String("prefix", conf.C.ServerCfg.Prefix),
	)

	gRpcServer := grpc.NewServer(opts...)

	RegisterAllHandlers(gRpcServer)

	//reflection
	reflection.Register(gRpcServer)

	if err = gRpcServer.Serve(lis); err != nil {
		log.Panic(zap.Error(err))
	}

	f := func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			registrar.Deregister()
			return fmt.Errorf("received signal %s", sig)
		}
	}
	log.Info(zap.Error(f()))
}

func RegisterAllHandlers(s *grpc.Server) {}

// UnaryServerInterceptor Interceptor log printing
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
			log.Error(zap.Error(err))
		}
	}()
	log.Infox(info.FullMethod, zap.Any("request", fmt.Sprintf("%+v", req)))

	return handler(ctx, req)
}
