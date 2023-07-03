package config

import (
	"treafik-api/pkg/common/config"
	databases2 "treafik-api/pkg/common/databases"
)

type Config struct {
	Server Server                   `json:"service" yaml:"server"  toml:"server"`
	MySQL  databases2.MysqlDBConfig `json:"mysql" yaml:"mysql"  toml:"mysql"`
	Redis  databases2.RedisConfig   `json:"redis" yaml:"redis" toml:"redis"`
	Log    config.LogConfig         `json:"log" yaml:"log" toml:"log"`
}

type Server struct {
	Port         string   `json:"port" yaml:"port" toml:"port"`
	Address      string   `json:"address" yaml:"address" toml:"address"`
	ReadTimeout  int      `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout int      `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	Mode         string   `json:"mode" yaml:"mode" toml:"mode"`
	Middlewares  []string `json:"middlewares" yaml:"middlewares" toml:"middlewares"`
}
