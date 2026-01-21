package MiddleWare

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CrossDomain struct {
}

func (r *CrossDomain) CrossDomainHandler() gin.HandlerFunc {

	corsHdlFunc := cors.New(cors.Config{
		AllowCredentials: true, //是否允许带cookie等用户凭据，正常都需要允许

		// 允许的请求头,并希望在前端请求时把token从Authorization带回来
		AllowHeaders: []string{"content-type", "Authorization"},

		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"}, //允许前端访问后端响应头部,让前端能看到这个头部
		AllowOriginFunc: func(origin string) bool {
			// if strings.HasPrefix(origin, "http://localhost")判定包含前缀
			if strings.Contains(origin, "localhost") { //判断包含该字段
				return true
			}
			if strings.Contains(origin, "192.168.80.1") { //判断包含该字段
				return true
			}
			if strings.Contains(origin, "101.34.33.105") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour, //检测时间长度
	})
	return func(ctx *gin.Context) {
		corsHdlFunc(ctx)
	}
}
