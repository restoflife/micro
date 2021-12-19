/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 15:32
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 15:32
 * @FilePath: ql-gateway/internal/protocol/order.go
 */

package protocol

type (
	GetOrderDetailsReq struct {
		Id int64 `form:"id" binding:"required"`
	}
)
