package db

import (
	"gorm.io/gorm"

	"treafik-api/config"
	"treafik-api/pkg/common/databases"
)

func initMysqlDB(cfg *config.Config) (*gorm.DB, error) {
	mysqlDatabase, err := databases.InitGorm(&cfg.MySQL)
	if err != nil {
		return nil, err
	}
	return mysqlDatabase, nil
}
