/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:14
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:14
 * @FilePath: ql-mp/cmd/api.go
 */

package main

import "github.com/restoflife/micro/mp/command"

func main() {
	_ = command.RootCmd.Execute()
}
