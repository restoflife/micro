/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 11:50
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 11:50
 * @FilePath: ql-mp/internal/component/hold/hold.go
 */

package hold

import (
	"encoding/json"
	"fmt"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/restoflife/micro/mp/internal/constant"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	tokenLock sync.RWMutex
	tokenMgr  string
)

func MustStartup() error {
	if err := holdAccessToken(conf.C.Wechat.APPID, conf.C.Wechat.SECRET); err != nil {
		log.Error(zap.String("app_Id", conf.C.Wechat.APPID),
			zap.Error(err))
		return err
	}
	return nil
}
func holdAccessToken(appId, appSecret string) error {
	token, expire, err := getWXAccessToken(appId, appSecret)
	if err != nil {
		return err
	}
	tokenLock.Lock()
	tokenMgr = token
	tokenLock.Unlock()
	log.Infox("init accessToken",
		zap.String("token", token),
		zap.Int("expire", expire))
	go func() {
		ticker := time.NewTicker(time.Second * 7000)
		for {
			select {
			case <-ticker.C:
				token, expire, err = getWXAccessToken(appId, appSecret)
				if err != nil {
					log.Error(zap.Error(err))
					continue
				}

				tokenLock.Lock()
				tokenMgr = token
				tokenLock.Unlock()

				log.Infox("refresh accessToken",
					zap.String("token", token),
					zap.Int("expire", expire))
			}
		}
	}()
	return nil
}

func getWXAccessToken(appId, secret string) (token string, expire int, err error) {
	url := fmt.Sprintf(constant.UrlGetAccessToken, appId, secret)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error(zap.Error(err))
		}
	}(resp.Body)

	result := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Token   string `json:"access_token"`
		Expire  int    `json:"expires_in"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != 0 {
		err = fmt.Errorf(strconv.Itoa(result.ErrCode) + "_" + result.ErrMsg)
		return
	}
	token = result.Token
	expire = result.Expire
	return
}

func GetAccessToken() string {
	tokenLock.RLock()
	token := tokenMgr
	tokenLock.RUnlock()

	return token
}
