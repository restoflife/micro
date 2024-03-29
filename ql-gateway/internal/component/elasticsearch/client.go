/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-30 10:51
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-30 10:51
 * @FilePath: ql-gateway/internal/component/elasticsearch/.go
 */

package elasticsearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Cli *Elastic7Endpoint

func NewElasticSearchClient(config *conf.ElasticConfig) error {
	options := []elastic.ClientOptionFunc{
		elastic.SetErrorLog(new(errLog)),
		elastic.SetURL(config.Host),
		elastic.SetBasicAuth(config.Username, config.Password),
		// 是否开启集群嗅探
		elastic.SetSniff(config.Sniff),
		// 设置两次运行状况检查之间的间隔, 默认60s
		elastic.SetHealthcheckInterval(time.Duration(config.Heal) * time.Second),
		// 启用或禁用gzip压缩
		elastic.SetGzip(config.Gzip),
		elastic.SetInfoLog(new(infoLog)),
	}
	con, err := elastic.NewClient(options...)
	if err != nil {
		return err
	}
	Cli = &Elastic7Endpoint{
		Client: con,
		ctx:    context.Background(),
	}
	result, code, err := Cli.Ping(config.Host).Do(context.Background())
	if err != nil || code != http.StatusOK {
		return err
	}

	log.Infox("elasticsearch returned wit", zap.Int("code", code), zap.String("version", result.Version.Number))
	return nil
}

// NewClient 获取一个连接
func NewClient() *elastic.Client {
	if Cli.Client == nil {
		log.Error(zap.Error(fmt.Errorf("get elastic error")))
		return nil
	}
	return Cli.Client
}

func (c *Elastic7Endpoint) existsIndex(index string) (exists bool, err error) {
	return c.IndexExists(index).Do(c.ctx)

}

// Insert 写入数据
func (c *Elastic7Endpoint) Insert(index string, value interface{}) (resp *elastic.IndexResponse, err error) {
	var exists bool
	exists, err = c.existsIndex(index)
	if err != nil {
		return
	}
	if !exists {
		err = c.CreateIndexes(index, value)
		if err != nil {
			return
		}
	}
	resp, err = c.Index().
		Index(index).
		BodyJson(value).
		Do(c.ctx)
	if err != nil {
		log.Error(zap.Error(err))
		return
	}
	return
}

// CreateIndexes 创建索引
func (c *Elastic7Endpoint) CreateIndexes(index string, mapping interface{}) (err error) {
	var exists bool
	exists, err = c.existsIndex(index)
	if err != nil {
		return
	}
	if exists {
		return
	}
	_, err = c.CreateIndex(index).
		BodyJson(mapping).
		Do(c.ctx)
	if err != nil {
		return
	}
	return
}

// SearchQuery 查询数据
// elastic.TermQuery() 精确匹配单个字段
// elastic.TermsQuery() 精确匹配单个字段，但使用多值进行匹配，类似于 SQL 中的 in 操作
// elastic.MatchQuery() 单个字段匹配查询（匹配分词结果，不需要全文匹配）
// elastic.RangeQuery() 范围查询
// elastic.BoolQuery() 组合查询
func (c *Elastic7Endpoint) SearchQuery(field *SearchQueryReq) (*elastic.SearchResult, int64, error) {
	limit, offset := utils.PageIndex(field.Limit, field.Offset)
	resp, err := c.Search(field.Index).
		From(offset).
		Size(limit).
		Query(field.Query).
		Pretty(true).
		// Sort(field.Sort, true).
		Do(c.ctx)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, 0, err
	}
	total, err := c.Count(field.Index).Do(c.ctx)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, 0, err
	}
	return resp, total, nil
}

// DeleteQuery 删除单个文档
func (c *Elastic7Endpoint) DeleteQuery(field *DeleteQueryReq) error {
	_, err := c.DeleteByQuery(field.Index).
		Query(field.Query).
		Refresh("true").
		Do(c.ctx)
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}
	return nil
}
