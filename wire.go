//go:build wireinject

package main

import (
	"DentTick/Handler"
	"DentTick/Handler/Jwt"
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
		thirdParty,
		ServerSet,
	)
	return gin.Default()
}

var userSet = wire.NewSet(
	//Dao
	Dao.NewUserGormDao,
	//Repository
	Repository.NewUserRepository,
	//Service
	Service.NewUserService,
	//Handler
	Handler.NewUserHandler,
	Jwt.NewRedisJWTHandler,
)

var ServerSet = wire.NewSet(
	Ioc.InitLogger,
	Ioc.InitWebServer,
	Ioc.InitMiddlerWares,
)

var thirdParty = wire.NewSet(
	Ioc.InitDB,
	Ioc.InitRedis,
)
