/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-17 10:56
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-17 10:56
 * @FilePath: ql-order/command/root.go
 */

package command

import (
	"github.com/restoflife/micro/order/internal/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "order",
	Short: "operation background",
	Long:  "no only book operation background",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(NewApp("order", cmd))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.Flags().String("c", "configs/config.toml", "the path to the config file")
}
