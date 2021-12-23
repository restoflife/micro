/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:15
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:15
 * @FilePath: ql-mp/conf/conf.go
 */

package conf

var C *Config

type (
	Config struct {
		ServerCfg    *ServerConfig          `toml:"server"`
		RunLogCfg    *LogConfig             `toml:"runLog"`
		AccessLogCfg *LogConfig             `toml:"accessLog"`
		DB           map[string]*ConfigLite `toml:"db"`
		Redis        *RedisConfig           `toml:"redis"`
		Login        *login                 `toml:"login"`
	}
)

type (
	login struct {
		Total   int    `toml:"total"`
		Time    int    `toml:"time"`
		Key     string `toml:"key"`
		Timeout int64  `toml:"timeout"`
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
		Addr        string `toml:"addr"`
		Mode        bool   `toml:"mode"`
		RPCAddr     string `toml:"rpc_addr"`
		Etcd        string `toml:"etcd"`
		EtcdCert    string `toml:"etcd_cert"`
		EtcdKey     string `toml:"etcd_key"`
		EtcdCaCert  string `toml:"etcd_ca_cert"`
		Prefix      string `toml:"prefix"`
		OrderPrefix string `toml:"order_prefix"`
		LogPrefix   string `toml:"log_prefix"`
	}
)
