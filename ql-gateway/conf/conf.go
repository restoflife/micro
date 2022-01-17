/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:50
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:50
 * @FilePath: ql-gateway/conf/conf.go
 */

package conf

var C *Config

type (
	Config struct {
		ServerCfg    *ServerConfig           `toml:"server"`
		RunLogCfg    *LogConfig              `toml:"runLog"`
		AccessLogCfg *LogConfig              `toml:"accessLog"`
		SQLLogCfg    *LogConfig              `toml:"sqlLog"`
		DB           map[string]*ConfigLite  `toml:"db"`
		Redis        *RedisConfig            `toml:"redis"`
		Mongo        map[string]*MongoConfig `toml:"mongo"`
		Login        *login                  `toml:"login"`
		GRpcCli      map[string]grpcCli      `toml:"grpc_cli"`
		Elastic      *ElasticConfig          `toml:"elastic"`
	}
)

type (
	ElasticConfig struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Sniff    bool   `json:"sniff"`
		Heal     int    `json:"heal"`
		Gzip     bool   `json:"gzip"`
	}
	grpcCli struct {
		Prefix string `toml:"prefix"`
	}
	login struct {
		Total   int    `toml:"total"`
		Time    int    `toml:"time"`
		Key     string `toml:"key"`
		Timeout int64  `toml:"timeout"`
	}
	MongoConfig struct {
		Addr        string `toml:"addr"`
		MaxIdleTime int    `toml:"max_idle_time"`
		MaxPool     uint64 `toml:"max_pool"`
		MinPool     uint64 `toml:"min_pool"`
		Username    string `toml:"username"`
		Password    string `toml:"password"`
	}
	ConfigLite struct {
		Driver  string `toml:"driver"`
		Dsn     string `toml:"dsn"`
		MaxIdle int    `toml:"max_idle"`
		MaxOpen int    `toml:"max_open"`
		ShowSql bool   `toml:"show_sql"`
		Slave   []struct {
			Dsn string `toml:"dsn"`
		}
		Prefix   string `toml:"prefix"`
		Singular bool   `toml:"singular"`
		MaxLife  int    `toml:"max_life"`
	}
	RedisConfig struct {
		Addr       []string `toml:"addr"`
		Password   string   `toml:"password"`
		DB         int      `toml:"db"`
		MasterName string   `toml:"master_name"`
		PoolSize   int      `toml:"pool_size"`
		IdleSize   int      `toml:"idle_conns"`
	}

	LogConfig struct {
		Level      string `toml:"level"`
		Filename   string `toml:"file"`
		MaxSize    int    `toml:"maxSize"`
		MaxBackups int    `toml:"maxBackups"`
		MaxAge     int    `toml:"maxAge"`
	}

	ServerConfig struct {
		Addr       string `toml:"addr"`
		Mode       bool   `toml:"mode"`
		Etcd       string `toml:"etcd"`
		EtcdCert   string `toml:"etcd_cert"`
		EtcdKey    string `toml:"etcd_key"`
		EtcdCaCert string `toml:"etcd_ca_cert"`
	}
)
