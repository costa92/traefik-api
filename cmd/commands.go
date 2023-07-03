package cmd

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"treafik-api/config"
)

var (
	cfgFile   string
	AppConfig *config.Config // 配置对应的结构体
	Verbose   bool
)

// rootCmd is the root command of the application
var rootCmd = &cobra.Command{
	Use:                "ctrapi",
	Short:              "ctrapi is a tool for managing CTR API",
	Long:               `ctrapi is a tool for managing CTR API`,
	DisableSuggestions: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.InitDefaultHelpFlag()
		return cmd.Help()
	},
	// 不需要出现cobra默认的completion子命令
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config.yaml | config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
}

func initConfig() {
	// 加載配置文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile) // 指定配置文件名
		// 如果配置文件名中没有文件扩展名，则需要指定配置文件的格式，告诉viper以何种格式解析文件
		viper.SetConfigType("yaml")
	} else {
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../")
		viper.AddConfigPath("")
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config") // 指定配置文件名
	}
	// 自动加载环境变量
	viper.AutomaticEnv()
	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("Fatal error config file: %s \n", err)
		return
	}
	// 解析配置信息
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}
}
