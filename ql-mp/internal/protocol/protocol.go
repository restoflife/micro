/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:39
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:39
 * @FilePath: ql-mp/internal/protocol/protocol.go
 */

package protocol

type (
	MpLoginReq struct {
		Code string `from:"code"`
	}
	MpLoginResp struct {
	}
)
