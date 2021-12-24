/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:17
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:17
 * @FilePath: ql-mp/internal/constant/constant.go
 */

package constant

const (
	DbDefaultName        = "default"
	RedisName            = "default"
	Layout               = "2006-01-02 15:04:05"
	ContextMpUUid        = "mp_uuid"
	UrlCode2Session      = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	UrlGetAccessToken    = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	RedisPrefixToken     = "ql_wechat_token:"
	RedisKeyWxSessionKey = "ql_wx_session_key:"
	ContextOrderKey      = "context_key"
)
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)
