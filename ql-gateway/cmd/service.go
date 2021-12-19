/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:27
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:27
 * @FilePath: ql-gateway/cmd/service.go
 */

package main

import "github.com/restoflife/micro/gateway/command"

func main() {
	_ = command.RootCmd.Execute()
}
