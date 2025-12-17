package Ioc

import (
	"DentTick/Handler"
	"DentTick/Handler/MiddleWare"

	"github.com/gin-gonic/gin"
)

func InitWebServer(middlewares []gin.HandlerFunc, userHandler *Handler.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandler.RegisterRoute(server)
	return server
}

func InitMiddlerWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		(&MiddleWare.CrossDomain{}).CrossDomainHandler(),
	}
}
