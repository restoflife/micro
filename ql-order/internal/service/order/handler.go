/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 11:20
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 11:20
 * @FilePath: ql-order/internal/service/order/handler.go
 */

package order

import (
	"fmt"
	orderPb "github.com/restoflife/micro/protos/order"
)

func getOrderDetails(id int64) (*orderPb.GetOrderDetailsResp, error) {
	// session, err := db.NewSession(constant.DbDefaultName)
	// if err != nil {
	//	return nil, err
	// }
	// defer db.Close(session)
	return &orderPb.GetOrderDetailsResp{Id: id, OrderId: fmt.Sprintf("1%03d", id)}, nil
}
