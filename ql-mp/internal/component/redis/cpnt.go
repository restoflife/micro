/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021/12/6 13:55
 * @LastEditors: Administrator
 * @LastEditTime: 2021/12/6 13:55
 * @FilePath: internal/component/redis/cpnt.go
 */

package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/restoflife/micro/mp/conf"
	"github.com/restoflife/micro/mp/utils"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ctx   = context.Background()
	Cli   redis.UniversalClient
	lockG = &singleflight.Group{}
)

// MustBootUp Start redis
func MustBootUp(config *conf.RedisConfig) error {
	Cli = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        config.Addr,
		DB:           config.DB,
		Password:     config.Password,
		MasterName:   config.MasterName,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.IdleSize,
	})
	_, err := Cli.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func SetCache(key string, data interface{}, duration time.Duration) error {
	key = utils.MD5String(key)
	dataMap := make(map[string]interface{})
	dataMap["data"] = data
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}
	return Cli.Set(ctx, key, jsonData, duration).Err()
}

func CheckCache(key string, fn func() (interface{}, error), duration time.Duration, needCache bool) (interface{}, error) {
	s, err := GetCache(key)
	if needCache && err == nil {
		return s, nil
	} else {
		var re interface{}
		//At the same time, only one function with the same key is executed to prevent breakdown
		Num, ok, _ := lockG.Do(key, fn)
		if ok == nil {
			_ = SetCache(key, Num, duration*time.Second)
			re = Num
		} else {
			re = Num
		}

		return re, ok
	}
}

func GetCache(key string) (interface{}, error) {
	key = utils.MD5String(key)
	data, err := Cli.Get(ctx, key).Result()
	if err == nil && data != "" {
		dom := gjson.Parse(data)
		return dom.Get("data").Value(), err
	}

	return data, err
}

func GetCacheByFloat(key string) (float64, error) {
	key = utils.MD5String(key)
	data, err := Cli.Get(ctx, key).Result()
	if err == nil && data != "" {
		dom := gjson.Parse(data)
		return dom.Get("data").Num, err
	}

	return 0, err
}

// GetRedisExpTime Expiration time
func GetRedisExpTime(key string) (time.Duration, error) {
	key = utils.MD5String(key)
	data, err := Cli.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return data, nil
}

func DelCache(key ...string) error {
	var keys []string
	if len(key) == 1 {
		keys = append(keys, utils.MD5String(key[0]))
	} else {
		for _, v := range key {
			keys = append(keys, utils.MD5String(v))
		}
	}
	return Cli.Del(ctx, keys...).Err()
}

// DepositZIncrBy Write Leaderboard
func DepositZIncrBy(key string, z *redis.Z) error {
	key = utils.MD5String(key)
	return Cli.ZIncrBy(ctx, key, z.Score, z.Member.(string)).Err()
}

// GetDepositByRedis Query Leaderboard
func GetDepositByRedis(key string, start, stop int64) ([]redis.Z, error) {
	key = utils.MD5String(key)
	return Cli.ZRevRangeWithScores(ctx, key, start, stop).Result()
}
