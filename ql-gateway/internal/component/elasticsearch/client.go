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
	"github.com/olivere/elastic/v7"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"go.uber.org/zap"
	"net/http"
)

var Cli *Elastic7Endpoint

type Elastic7Endpoint struct {
	*elastic.Client
	ctx context.Context
}

const host = "http://127.0.0.1:9200"

type Logger struct{}

func (*Logger) Printf(format string, v ...interface{}) {
	log.Error(zap.Any(format, v))
}
func NewElasticSearchClient() error {
	newClient, err := elastic.NewClient(
		elastic.SetErrorLog(new(Logger)),
		elastic.SetURL(host),
		//elastic.SetBasicAuth("", ""),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Err(zap.Error(err))
		return err
	}
	Cli = &Elastic7Endpoint{Client: newClient, ctx: context.Background()}
	result, code, err := Cli.Ping(host).Do(context.Background())
	if err != nil || code != http.StatusOK {
		return err
	}
	log.Infox("elasticsearch returned wit", zap.Int("code", code), zap.String("version", result.Version.Number))
	return nil
}

func (c *Elastic7Endpoint) Insert(index string, value interface{}) (*elastic.IndexResponse, error) {
	response, err := c.Index().
		Index(index).
		BodyJson(value).
		Do(c.ctx)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (c *Elastic7Endpoint) Create(index string, mapping interface{}) error {
	_, err := c.CreateIndex(index).
		BodyJson(mapping).
		Do(c.ctx)
	if err != nil {
		return err
	}
	return nil
}

const Mapping = `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "influenced_id": {
        "type": "integer"
      },
      "time": {
       "type": "date",
       "format": "yyyy-MM-dd"
      },
      "fans": {
        "type": "integer"
      },
      "play": {
        "type": "integer"
      },
      "like": {
        "type": "integer"
      },
      "opus": {
        "type": "integer"
      },
      "live_money": {
        "type": "float"
      },
      "ad_money": {
        "type": "float"
      },
      "shop_money": {
        "type": "float"
      },
      "total_money": {
        "type": "float"
      },
      "live_time": {
        "type": "float"
      },
      "ad_order": {
        "type": "integer"
      },
      "shop_sales": {
        "type": "float"
      },
      "shop_order": {
        "type": "integer"
      },
      "shop_gold": {
        "type": "float"
      },
      "created_at": {
        "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss"
      },
      "updated_at": {
        "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss"
      }
    }
  }
}`
