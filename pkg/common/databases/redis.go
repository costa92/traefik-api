package databases

import (
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string `toml:"addr" json:"addr" yaml:"addr" mapstructure:"addr" validate:"hostname_port" long:"addr" description:"redis server addr,format is host:port"`
	PoolSize int    `toml:"pool_size" json:"pool_size"  mapstructure:"pool_size" yaml:"pool_size" long:"pool_size" description:"redis connection pool size"`
	Passwd   string `toml:"passwd" json:"passwd"  mapstructure:"passwd" yaml:"passwd" long:"passwd" description:"redis auth passwd, leave it empty if no auth needed"`
	Tracing  bool   `toml:"tracing" json:"tracing" mapstructure:"tracing" yaml:"tracing"`
}

func NewRedisClient(redisCfg RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     redisCfg.Addr,
		Password: redisCfg.Passwd,
		PoolSize: redisCfg.PoolSize,
	})
	if redisCfg.Tracing {
		client.AddHook(redisotel.NewTracingHook())
	}
	return client
}
