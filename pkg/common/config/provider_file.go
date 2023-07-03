package config

import (
	"errors"
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"treafik-api/pkg/logger"
)

var ErrEmptyConfigFile = errors.New("error empty config file")

type FileProvider struct {
	DefaultConfigPath     string
	SkipIfPathEmpty       bool // skip this provider if config file path is empty
	SkipIfDefaultNotExist bool
}

func (p *FileProvider) Name() string {
	return "file"
}

var _ Provider = &FileProvider{}

func (p *FileProvider) Config(helper *providerHelper, cfg interface{}) ([]byte, error) {
	configFile := helper.configFile
	initConfig(configFile, cfg)
	return nil, nil
}

func initConfig(cfgFile string, cfg interface{}) {
	// 加載配置文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile) // 指定配置文件名
		// 如果配置文件名中没有文件扩展名，则需要指定配置文件的格式，告诉viper以何种格式解析文件
		viper.SetConfigType("yaml")
	} else {
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
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	logger.Infow("config file", "path", viper.ConfigFileUsed())
	// 解析配置信息
	err := viper.Unmarshal(cfg)
	logger.Infow("config file content:", "config", cfg)
	if err != nil {
		panic(err)
	}
}
