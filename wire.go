//go:build wireinject

package main

import (
	"DentTick/Handler"
	"DentTick/Ioc"
	"DentTick/Repository"
	"DentTick/Repository/Dao"
	"DentTick/Service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWireServer() *gin.Engine {
	wire.Build(
		userSet,
		//Ioc
		Ioc.InitDB, Ioc.InitLogger,
		Ioc.InitWebServer,
		Ioc.InitMiddlerWares,
	)
	return gin.Default()
}

var userSet = wire.NewSet(
	//Dao
	Dao.NewUserDao,
	//Repository
	Repository.NewUserRepository,
	//Service
	Service.NewUserService,
	//Handler
	Handler.NewUserHandler,
)
