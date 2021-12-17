/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 10:56
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 10:56
 * @FilePath: ql-order/cmd/service.go
 */

package main

import "github.com/restoflife/micro/order/command"

func main() {
	_ = command.RootCmd.Execute()
}
