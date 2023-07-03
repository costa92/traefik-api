package controller

import (
	"treafik-api/config"
	v1 "treafik-api/controller/v1"
	"treafik-api/db"
)

type ApiHttp struct {
	Index *v1.Index
	Dbs   *db.Databases
}

func NewApiHttp(cfg *config.Config, dbs *db.Databases) *ApiHttp {
	return &ApiHttp{
		Index: v1.NewIndex(dbs),
	}
}
