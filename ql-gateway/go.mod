module github.com/restoflife/micro/gateway

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-kit/kit v0.12.0
	github.com/mattn/go-isatty v0.0.14
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	go.uber.org/zap v1.20.0
	google.golang.org/grpc v1.43.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/gin-contrib/pprof v1.3.0
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/goccy/go-json v0.8.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.1.2
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/mitchellh/mapstructure v1.4.3
	github.com/mojocn/base64Captcha v1.3.5
	github.com/olivere/elastic/v7 v7.0.30
	github.com/restoflife/log v1.9.7 // indirect
	github.com/restoflife/micro/protos v1.1.0
	github.com/satori/go.uuid v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/tidwall/gjson v1.12.1
	go.etcd.io/etcd/client/v3 v3.5.1 // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gorm.io/driver/mysql v1.2.3
	gorm.io/gorm v1.22.4
	gorm.io/plugin/dbresolver v1.1.0
	moul.io/zapgorm2 v1.1.1 // indirect
	xorm.io/builder v0.3.9
	xorm.io/xorm v1.2.5
)

replace github.com/restoflife/micro/protos => ../ql-protos
