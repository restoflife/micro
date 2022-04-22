/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-04-22 16:50
 * @LastEditors: Administrator
 * @LastEditTime: 2022-04-22 16:50
 * @FilePath: ql-gateway/internal/ihttp/client_test.go
 */

package ihttp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type cResponse struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
	Time   int64       `json:"time"`
}

func TestHttpGet(t *testing.T) {
	client := NewClient()
	client.Timeout = 10 * time.Second
	_, bs, err := client.Get("http://127.0.0.1:1800/api/passport/captcha").EndBytes(context.Background())
	if err != nil {
		return
	}
	resp := new(cResponse)
	_ = json.Unmarshal(bs, resp)
	fmt.Println(resp)
	return
}
