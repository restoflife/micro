/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 16:34
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 16:34
 * @FilePath: ql-gateway/internal/protocol/common.go
 */

package protocol

type CommonListResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
