package v1

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"treafik-api/core/code"
	"treafik-api/db"
	"treafik-api/pkg/utils"
)

type Auth struct {
	MysqlDb *db.Databases
}

func NewAuth(dbs *db.Databases) *Auth {
	return &Auth{
		MysqlDb: dbs,
	}
}

type LoginRequest struct {
	Username string `json:"username,omitempty" form:"username" binding:"required" `
	Password string `json:"password" form:"password" binding:"required"`
}

func (a *Auth) Login(ctx *gin.Context) {
	var req LoginRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		utils.WriteErrResponse(ctx, errors.WithCode(code.ErrBind, "账号密码错误"))
		return
	}

	utils.WriteSuccessResponse(ctx, "3213")
}
