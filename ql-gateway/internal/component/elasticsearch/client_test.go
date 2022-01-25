/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-31 15:42
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-31 15:42
 * @FilePath: ql-gateway/internal/component/elasticsearch/client_test.go
 */

package elasticsearch

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/restoflife/micro/gateway/internal/component/log"
	userPb "github.com/restoflife/micro/protos/mp"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestElastic7Endpoint_SearchQuery(t *testing.T) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery("nickname", "余生"))
	boolQuery.Must(elastic.NewMatchQuery("uid", 211703197922754562))
	result, total, err := Cli.SearchQuery(&SearchQueryReq{
		Index:  "mp",
		Query:  boolQuery,
		Limit:  10,
		Offset: 1,
	})
	if err != nil {
		log.Error(zap.Error(err))
		return
	}
	resp := make([]*userPb.GetUserListItem, 0)
	for _, v := range result.Each(reflect.TypeOf(&userPb.GetUserListItem{})) {
		resp = append(resp, v.(*userPb.GetUserListItem))
	}
	fmt.Println(resp)
	fmt.Println(total)
}
