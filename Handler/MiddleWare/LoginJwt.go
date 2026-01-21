package MiddleWare

import (
	"DentTick/Constant"
	ijwt "DentTick/Handler/Jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTBuilder struct {
	ijwt.Handler
}

func NewLoginJWTBuilder(hdl ijwt.Handler) *LoginJWTBuilder {
	return &LoginJWTBuilder{
		Handler: hdl,
	}
}
func (b *LoginJWTBuilder) CheckLogin() (jwtHdlFunc gin.HandlerFunc) {
	jwtHdlFunc = func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/refresh_token" {
			return
		}

		//约定 token在Authorization的Bearer一起请求
		tokenStr := b.ExtractToken(ctx)

		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return Constant.JwtKey, nil
		})

		//token不对是伪造的 || token没解析出来 || token是非法的或是过期的
		if err != nil || token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//校验完 token后再访问redis可降低一些无效的访问场景
		err = b.CheckSession(ctx, uc.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//严格做法
		if err != nil {
			// token 无效或是redis出问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", uc) //将其放入上下文中
		ctx.Next()
	}
	return
}
