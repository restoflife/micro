/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-27 11:39
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-27 11:39
 * @FilePath: ql-gateway/internal/component/grpccli/etcd.go
 */

package grpccli

import (
	"context"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"time"
)

var InstancedMgr = map[string]*etcdv3.Instancer{}

func MustBootUp() error {
	var (
		//注册中心地址
		etcdServer = conf.C.ServerCfg.Etcd
		ctx        = context.Background()
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
		log.Error(zap.Error(err))
		return err
	}

	//连接注册中心
	client, err := etcdv3.NewClient(ctx, addr, options)
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据

	for name, v := range conf.C.GRpcCli {
		instanced, err := etcdv3.NewInstancer(client, v.Prefix, kitLog.NewNopLogger())
		if err != nil {
			log.Error(zap.Error(err))
			return err
		}
		InstancedMgr[name] = instanced
	}
	return nil
}
