/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 10:56
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 10:56
 * @FilePath: ql-order/command/app.go
 */

package command

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/restoflife/micro/order/conf"
	"github.com/restoflife/micro/order/internal/app"
	"github.com/restoflife/micro/order/internal/component/db"
	"github.com/restoflife/micro/order/internal/component/log"
	"github.com/restoflife/micro/order/internal/component/redis"
	"github.com/restoflife/micro/order/internal/constant"
	"github.com/restoflife/micro/order/internal/model"
	"github.com/restoflife/micro/order/transport/order"
	"github.com/restoflife/micro/order/utils"
	order_pb "github.com/restoflife/micro/protos/order"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

type mainApp struct {
	*app.Base
}

var (
	gRpcServer *grpc.Server
	registrar  *etcdv3.Registrar
)

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
	log.Infox("configuration loading completed")

}

func (m *mainApp) Run() {
	f := func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			// 等待grpc连接断开
			gRpcServer.GracefulStop()
			// 注销etcd
			registrar.Deregister()
			return fmt.Errorf("received signal %s", sig)
		}
	}
	log.Infox("terminated grpc server", zap.Error(f()))
}

func (m *mainApp) BootUpServer() {
	// ETCD connection parameters
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

	// Create a connection
	client, err := etcdv3.NewClient(context.Background(), addr, option)
	if err != nil {
		log.Panic(zap.Error(err))
	}

	// Create a registration
	registrar = etcdv3.NewRegistrar(client, etcdv3.Service{
		Key:   conf.C.ServerCfg.Prefix,
		Value: conf.C.ServerCfg.RPCAddr,
	}, kitLog.NewNopLogger())
	// Start registration service
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
			}),
	}

	log.Infox("listening",
		zap.String("transport", "GRPC"),
		zap.String("address", conf.C.ServerCfg.RPCAddr),
		zap.String("prefix", conf.C.ServerCfg.Prefix),
	)
	gRpcServer = grpc.NewServer(opts...)

	RegisterAllHandlers(gRpcServer)

	// reflection
	reflection.Register(gRpcServer)

	go func() {
		if err = gRpcServer.Serve(lis); err != nil {
			log.Panic(zap.Error(err))
		}
	}()

}

func RegisterAllHandlers(s *grpc.Server) {
	opts := []grpcTransport.ServerOption{
		grpcTransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx = context.WithValue(ctx, constant.ContextOrderUUid, md.Get(constant.ContextOrderUUid))
			return ctx
		}),
		grpcTransport.ServerErrorHandler(log.NewZapLogErrorHandler()),
	}

	order_pb.RegisterOrderSvcServer(s, order.NewUserServer(opts...))
}

// UnaryServerInterceptor Interceptor log printing
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
			log.Error(zap.Error(err))
		}
	}()
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "[getClientIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return nil, status.Errorf(codes.Unauthenticated, "[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) < 1 {
		addSlice = append(addSlice, "未知ip来源")
	}

	ctx = context.WithValue(ctx, constant.ContextOrderKey, info.FullMethod)
	log.Infox(info.FullMethod, zap.String("ip", addSlice[0]), zap.Any("request", fmt.Sprintf("%+v", req)))
	return handler(ctx, req)
}
