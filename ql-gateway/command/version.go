/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 16:48
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 16:48
 * @FilePath: ql-gateway/command/version.go
 */

package command

import "github.com/spf13/cobra"

const version = "1.0.0"

// The version command prints this service.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "operation background",
	Long:  "no only book operation background",
	Run: func(cmd *cobra.Command, args []string) {
		println("version ", version)
	},
}
