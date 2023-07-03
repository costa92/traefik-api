package routers

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"treafik-api/controller"
	"treafik-api/pkg/utils"
)

func RegisterRouter(r *gin.Engine, api *controller.ApiHttp) {
	// 写一个启动 http 服务代码
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/traefik", func(c *gin.Context) {
			utils.WriteResponse(c, errors.WithCode(500, "参数错误"), "start traefik v2")
			// utils.WriteSuccessResponse(c, "start traefik v2")
		})
		apiv1.GET("/index", api.Index.Index)
	}
}