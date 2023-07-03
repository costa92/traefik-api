package v1

import (
	"github.com/gin-gonic/gin"

	"treafik-api/core"
)

type Index struct {
	MysqlDb *core.MysqlDb
}

func NewIndex(mysqlDb *core.MysqlDb) *Index {
	return &Index{
		MysqlDb: mysqlDb,
	}
}

func (a *Index) Index(ctx *gin.Context) {
	ctx.JSON(200, "success")
}
