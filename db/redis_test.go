package db

import (
	"context"
	"testing"
	"time"

	"treafik-api/config"
	"treafik-api/pkg/common/databases"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := &config.Config{
		Redis: databases.RedisConfig{
			Addr: "127.0.0.1:6379",
		},
	}
	cli := initRedisDbs(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cli.Set(ctx, "key", 1, 10)

	res := cli.Get(ctx, "key")
	t.Log(res.Result())
}
