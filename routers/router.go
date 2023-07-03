package routers

import (
	"github.com/gin-gonic/gin"
)

// InitApiRouter 路由初始化
func InitApiRouter(r *gin.Engine) {
	initApi(r)
}

func RegisterRouter(r *gin.Engine) {
	initApi(r)
}
