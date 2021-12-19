/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 13:58
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 13:58
 * @FilePath: ql-gateway/internal/component/grpccli/grpccli.go
 */

package order

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/utils"
	orderPb "github.com/restoflife/micro/protos/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"time"
)

var Instanced *etcdv3.Instancer

func InitOrder() error {
	var err error
	Instanced, err = NewOrderClient()
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}
	return nil
}
func NewOrderClient() (*etcdv3.Instancer, error) {
	var (
		//注册中心地址
		etcdServer = conf.C.ServerCfg.Etcd
		//监听的服务前缀
		prefix = conf.C.ServerCfg.OrderPrefix
		ctx    = context.Background()
	)

	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
		Cert:          conf.C.ServerCfg.EtcdCert,
		Key:           conf.C.ServerCfg.EtcdKey,
		CACert:        conf.C.ServerCfg.EtcdCaCert,
	}

	addr, err := utils.GetUrls(etcdServer)
	if err != nil {
		log.Panic(zap.Error(err))
	}

	//连接注册中心
	client, err := etcdv3.NewClient(ctx, addr, options)
	if err != nil {
		log.Panic(zap.Error(err))
	}

	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	instanced, err := etcdv3.NewInstancer(client, prefix, kitLog.NewNopLogger())
	if err != nil {
		return nil, err
	}
	return instanced, err

}

func ExecHandler(factory sd.Factory, req interface{}) (interface{}, error) {
	logger := kitLog.NewNopLogger()
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instanced的变化动态更新Factory创建的endPoint
	endPointer := sd.NewEndpointer(Instanced, factory, logger)
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endPointer)
	/**
	  我们可以通过负载均衡器直接获取请求的endPoint，发起请求*/
	reqEndPoint, _ := balancer.Endpoint()

	/**
	  也可以通过retry定义尝试次数进行请求
	*/
	reqEndPoint = lb.Retry(3, 3*time.Second, balancer)
	//现在我们可以通过 endPoint 发起请求了 \
	//req := struct{}{}
	var r interface{}
	var err error
	if r, err = reqEndPoint(context.Background(), req); err != nil {
		return nil, err
	}
	return r, err
}

func GetOrderDetails(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*orderPb.GetOrderDetailsReq)
		if !ok {

			return nil, errutil.ErrIllegalParameter
		}

		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, errutil.ErrRpcInternalServer
		}
		defer func() {
			_ = conn.Close()
		}()

		svc := orderPb.NewOrderSvcClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := svc.GetOrderDetails(ctx, req)
		if err != nil {
			return nil, errutil.ErrRpcRequest
		}
		return r, nil
	}, nil, nil
}
