/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:34
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:34
 * @FilePath: ql-gateway/internal/protocol/common.go
 */

package protocol

import (
	"github.com/uber/jaeger-client-go"
)

type CommonListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
type Config struct {
	ServiceName string                `json:"service_name" yaml:"service_name" toml:"service_name"`
	Endpoint    string                `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	TraceOpts   []jaeger.TracerOption `json:"trace_opts" yaml:"trace_opts" toml:"trace_opts"`
}
type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type EncryptJson struct {
	Param string `json:"param"`
}
type UploadImageResp struct {
	ImgUrl string `json:"img_url"`
}
