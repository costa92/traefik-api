package controller

import (
	"treafik-api/config"
	v1 "treafik-api/controller/v1"
	"treafik-api/db"
	"treafik-api/pkg/server/gin-server/auth"
)

type ApiHttp struct {
	Authorizer *auth.Authorization
	Dbs        *db.Databases
	Index      *v1.Index
	Auth       *v1.Auth
}

func NewApiHttp(cfg *config.Config, dbs *db.Databases) *ApiHttp {
	authorization := auth.NewAuthorization(cfg.Jwt)
	return &ApiHttp{
		Authorizer: authorization,
		Index:      v1.NewIndex(dbs),
		Auth:       v1.NewAuth(dbs),
	}
}
