/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-08-09 09:26
 * @LastEditors: Administrator
 * @LastEditTime: 2022-08-09 09:26
 * @FilePath: ql-gateway/utils/upload.go
 */

package utils

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"io"
)

func zone() *storage.Zone {
	switch conf.C.QiNiu.Bucket {
	case constant.UploadBucketQiNiu:
		return &storage.ZoneHuanan
	default:
		return &storage.Zone{}
	}
}

// UploadImageByQiNiu 七牛云上传
func UploadImageByQiNiu(data io.Reader, key string, size int64) (*protocol.UploadImageResp, error) {
	// 上传
	mac := qbox.NewMac(conf.C.QiNiu.AccessKey, conf.C.QiNiu.SecretKey)
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", conf.C.QiNiu.Bucket, key),
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          zone(), // 空间对应的机房
		UseHTTPS:      true,   // 是否使用https域名
		UseCdnDomains: true,   // 上传是否使用CDN上传加速
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, size, &storage.PutExtra{})
	if err != nil {
		return nil, err
	}

	return &protocol.UploadImageResp{
		ImgUrl: conf.C.QiNiu.Domain + "/" + ret.Key,
	}, nil
}

// UploadImageByAli 阿里云oss 上传
func UploadImageByAli(data io.Reader, filename string, _ int64) (*protocol.UploadImageResp, error) {
	client, err := oss.New(
		conf.C.AliOss.Domain,
		conf.C.AliOss.AccessKey,
		conf.C.AliOss.SecretKey,
		oss.UseCname(true),
	)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(conf.C.AliOss.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.PutObject(filename, data)
	if err != nil {
		return nil, err
	}
	return &protocol.UploadImageResp{
		ImgUrl: conf.C.AliOss.Domain + "/" + filename,
	}, nil
}
