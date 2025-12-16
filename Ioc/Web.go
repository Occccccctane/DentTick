package Ioc

import "github.com/gin-gonic/gin"

func InitWebServer(middlewares []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	return server
}

func InitMiddlerWares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
