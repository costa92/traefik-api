package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"treafik-api/config"
	"treafik-api/pkg/common/databases"
)

var BaseRedisApi *redis.Client

func initRedisDbs(cfg *config.Config) *redis.Client {
	BaseRedisApi = databases.NewRedisClient(cfg.Redis)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err := BaseRedisApi.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return BaseRedisApi
}
