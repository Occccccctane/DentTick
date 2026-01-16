package Ioc

import (
	"DentTick/Handler"
	ijwt "DentTick/Handler/Jwt"
	"DentTick/Handler/MiddleWare"

	"github.com/gin-gonic/gin"
)

func InitWebServer(middlewares []gin.HandlerFunc, userHandler *Handler.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandler.RegisterRoute(server)
	return server
}

func InitMiddlerWares(hdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//这里添加要过全局的中间件
		//跨域
		(&MiddleWare.CrossDomain{}).CrossDomainHandler(),
		//	JWT 验证
		MiddleWare.NewLoginJWTBuilder(hdl).CheckLogin(),
		//	TODO: 限流
	}
}
