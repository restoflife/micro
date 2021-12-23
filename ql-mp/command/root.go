/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:16
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:16
 * @FilePath: ql-mp/command/root.go
 */

package command

import (
	"github.com/restoflife/micro/mp/internal/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "mp",
	Short: "operation background",
	Long:  "no only book operation background",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(NewApp("mp", cmd))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.Flags().String("c", "configs/config.toml", "the path to the config file")
}
