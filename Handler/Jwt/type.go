package Jwt

import "github.com/gin-gonic/gin"

type Handler interface {
	// ExtractToken 提取token
	ExtractToken(ctx *gin.Context) string
	// SetLoginToken 登录设置token
	SetLoginToken(ctx *gin.Context, uid int64)
	// SetJWTToken 用来刷新短token
	SetJWTToken(c *gin.Context, uid int64, ssid string)
	// SetRefreshToken 用来设置长token
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	// ClearToken 登出清除token
	ClearToken(ctx *gin.Context) error
	// CheckSession 验证长token
	CheckSession(ctx *gin.Context, ssid string) error
}
