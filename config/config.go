package config

import (
	"treafik-api/pkg/common/config"
	"treafik-api/pkg/common/databases"
	"treafik-api/pkg/server"
)

type Config struct {
	Server server.ServiceConfig    `json:"service" yaml:"server"  toml:"server"`
	MySQL  databases.MysqlDBConfig `json:"mysql" yaml:"mysql"  toml:"mysql"`
	Redis  databases.RedisConfig   `json:"redis" yaml:"redis" toml:"redis"`
	Log    config.LogConfig        `json:"log" yaml:"log" toml:"log"`
}
