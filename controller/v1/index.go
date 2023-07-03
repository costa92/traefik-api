package v1

import (
	"github.com/gin-gonic/gin"

	"treafik-api/db"
	"treafik-api/pkg/utils"
)

type Index struct {
	MysqlDb *db.Databases
}

func NewIndex(mysqlDb *db.Databases) *Index {
	return &Index{
		MysqlDb: mysqlDb,
	}
}

func (a *Index) Index(ctx *gin.Context) {
	utils.WriteSuccessResponse(ctx, "success")
}
