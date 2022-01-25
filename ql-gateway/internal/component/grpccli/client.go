/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-27 09:24
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-27 09:24
 * @FilePath: ql-gateway/internal/component/grpccli/client.go
 */

package grpccli

import (
	"context"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/errutil"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"time"
)

func ExecHandler(src sd.Instancer, factory sd.Factory, req interface{}) (interface{}, error) {
	logger := kitLog.NewNopLogger()
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instanced的变化动态更新Factory创建的endPoint
	endPointer := sd.NewEndpointer(src, factory, logger)
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endPointer)
	// 我们可以通过负载均衡器直接获取请求的endPoint，发起请求
	reqEndPoint, _ := balancer.Endpoint()

	//也可以通过retry定义尝试次数进行请求
	//todo:临时只尝试一次
	reqEndPoint = lb.Retry(1, 5*time.Second, balancer)

	//现在我们可以通过 endPoint 发起请求了
	var (
		err error
		r   interface{}
	)
	if r, err = reqEndPoint(context.Background(), req); err != nil {
		return nil, err
	}
	return r, err
}

func GrpcClientCtx(kv string, t int64) (context.Context, context.CancelFunc, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), kv).String()
	md := metadata.Pairs(kv, UUID)
	return metadata.NewOutgoingContext(ctx, md), cancel, UUID
}

// GrpcClient TODO: grpc.WithInsecure() 已弃用,
func GrpcClient(instanceAddr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(unaryInterceptor),
	}
	conn, err := grpc.Dial(instanceAddr, opts...)

	if err != nil {
		return nil, errutil.ErrRpcInternalServer
	}
	return conn, err
}

func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err := status.Errorf(codes.Internal, "Panic err: %v", e)
			log.Error(zap.Error(err))
		}
	}()
	//ctx = metadata.AppendToOutgoingContext(ctx, "k", "v")
	return invoker(ctx, method, req, reply, cc, opts...)
}
