module github.com/restoflife/micro/order

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/go-kit/kit v0.12.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/restoflife/micro/protos v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/tidwall/gjson v1.12.1
	go.uber.org/zap v1.19.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.43.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	xorm.io/xorm v1.2.5
)

replace github.com/restoflife/micro/protos => ../ql-protos
