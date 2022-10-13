/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:08
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:08
 * @FilePath: ql-gateway/internal/constant/constant.go
 */

package constant

const (
	DbDefaultName     = "default"
	RedisName         = "default"
	Layout            = "2006-01-02 15:04:05"
	RFC3339Nano       = "2006-01-02T15:04:05.000Z0700"
	UserLoginTotal    = "_password_err_total" // 登陆错误次数
	ContextOrderUUid  = "order_uuid"
	ContextMpUUid     = "mp_uuid"
	UploadBucketQiNiu = "blogdeng"
)
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// grpc prefix
const (
	OrderPrefix = "order"
	MpPrefix    = "mp"
)
