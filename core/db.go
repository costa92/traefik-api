package core

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"treafik-api/config"
)

type Databases struct {
	MysqlDb      *gorm.DB
	BaseRedisApi *redis.Client
}

var MySQLStorage *gorm.DB

func NewDatabases(cfg *config.Config) (*Databases, error) {
	var err error
	MySQLStorage, err = initMysqlDB(cfg)
	if err != nil {
		return nil, err
	}
	// baseRedisApi := initRedisDbs(cfg)
	return &Databases{
		MysqlDb: MySQLStorage,
		// BaseRedisApi: baseRedisApi,
	}, nil
}
