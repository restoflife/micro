/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-31 16:14
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-31 16:14
 * @FilePath: ql-gateway/internal/component/elasticsearch/struct.go
 */

package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
)

func (*errLog) Printf(format string, v ...interface{}) {
	log.Error(zap.Any(format, v))
}
func (*infoLog) Printf(format string, v ...interface{}) {
	log.Info(zap.Any(format, v))
}

type (
	errLog  struct{}
	infoLog struct{}

	Elastic7Endpoint struct {
		*elastic.Client
		ctx context.Context
	}

	SearchQueryReq struct {
		Index  string        `json:"index"`
		Query  elastic.Query `json:"query"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
		Sort   string        `json:"sort"`
	}
	DeleteQueryReq struct {
		Index string        `json:"index"`
		Query elastic.Query `json:"query"`
	}
)
