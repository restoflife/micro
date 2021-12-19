/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:48
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:48
 * @FilePath: ql-gateway/command/root.go
 */

package command

import (
	"github.com/restoflife/micro/gateway/internal/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "operation background",
	Long:  "no only book operation background",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(NewApp("gateway", cmd))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.Flags().String("c", "configs/config.toml", "the path to the config file")
}
