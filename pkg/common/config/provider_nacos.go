package config

import (
	"os"
	"strconv"

	remote "github.com/yoyofxteam/nacos-viper-remote"
)

type NacosProvider struct{}

func (p *NacosProvider) Name() string {
	return "nacos"
}

var _ Provider = &NacosProvider{}

func (p *NacosProvider) Config(helper *providerHelper, cfg interface{}) ([]byte, error) {
	NacosConfig(helper.configFile, cfg)
	return nil, nil
}

func NacosConfig(configFile string, cfg interface{}) {
	host := os.Getenv(EnvNacosHost)
	portStr := os.Getenv(EnvNacosPort)
	port, _ := strconv.ParseInt(portStr, 10, 64)
	//// 配置 Vip for Nacos 的远程
	remote.SetOptions(&remote.Option{
		Url:  host,
		Port: uint64(port),
	})
}
