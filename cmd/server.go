package cmd

import (
	"github.com/spf13/cobra"

	"treafik-api/core"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := core.NewDatabases(AppConfig)
		if err != nil {
			panic(err)
		}
		// 实例化服务
		appService := core.NewAppService(AppConfig)
		// 启动服务
		appService.Start()
	},
}
