//go:build wireinject

package main

import (
	"DentTick/Ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWireServer() *gin.Engine {
	wire.Build(
		userSet,
		//Ioc
		Ioc.InitWebServer,
		Ioc.InitMiddlerWares,
	)
	return gin.Default()
}

var userSet = wire.NewSet(
	//Service
	Service.NewUserService,
	//Handler
	Handler.NewUserHandler,
)
