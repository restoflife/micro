module github.com/restoflife/micro/mp

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/mattn/go-isatty v0.0.14
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/tidwall/gjson v1.12.1
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.43.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	xorm.io/builder v0.3.9
	xorm.io/xorm v1.2.5
)

require (
	github.com/go-kit/kit v0.12.0
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
)

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/mitchellh/mapstructure v1.4.3
	github.com/restoflife/micro/protos v1.0.0
	github.com/speps/go-hashids/v2 v2.0.1
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb // indirect
)

replace github.com/restoflife/micro/protos => ../ql-protos
