package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"treafik-api/pkg/server/gin-server/auth"
)

type JwtMiddleware struct {
	Auth *auth.Authorization
}

func NewJwtMiddleware(auth *auth.Authorization) *JwtMiddleware {
	return &JwtMiddleware{
		Auth: auth,
	}
}

func (j *JwtMiddleware) JwtAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    4001,
				"message": "Token Export",
			})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    2004,
				"message": "请求头中的auth格式错误",
			})
			// 阻止调用后续的函数
			ctx.Abort()
			return
		}
		mc, err := j.Auth.ParseJwt(parts[1])
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    2005,
				"message": "无效的token",
			})
			ctx.Abort()
			return
		}
		ctx.Set(auth.UserKey, mc.ID)
		ctx.Set("username", mc.Username)
		ctx.Next()
	}
}
