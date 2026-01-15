package Handler

import (
	"DentTick/Constant"
	"DentTick/Domain"
	ijwt "DentTick/Handler/Jwt"
	"DentTick/Package/logger"
	"DentTick/Service"
	"errors"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ijwt.Handler
	passwordRexExp *regexp.Regexp
	phoneRexExp    *regexp.Regexp
	svc            Service.UserService
	l              logger.Logger
}

func NewUserHandler(svc Service.UserService, hdl ijwt.Handler, l logger.Logger) *UserHandler {
	//读取正则
	type Config struct {
		passwordRegex string `yaml:"passwordRegex"`
		phoneRegex    string `yaml:"phoneRegex"`
	}
	var cfg Config
	err := viper.UnmarshalKey("regexp", &cfg)
	if err != nil {
		panic(err)
	}
	return &UserHandler{
		svc:            svc,
		passwordRexExp: regexp.MustCompile(cfg.passwordRegex, regexp.None),
		phoneRexExp:    regexp.MustCompile(cfg.phoneRegex, regexp.None),
		l:              l,
		Handler:        hdl,
	}
}
func (h *UserHandler) RegisterRoute(server *gin.Engine) {
	user := server.Group("/users")
	user.POST("/signup", h.Signup)
	user.POST("logout", h.Logout)

	user.GET("/refresh_token", h.RefreshToken)

}

func (h *UserHandler) Signup(ctx *gin.Context) {
	type signUpReq struct {
		Phone           string `json:"phone"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signUpReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	// 校验手机格式
	isPhoneTrue, err := h.phoneRexExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("手机格式错误", logger.Error(err))
		return
	}
	if !isPhoneTrue {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机格式错误",
		})
		return
	}

	//校验密码
	isPasswordTrue, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isPasswordTrue {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "密码格式错误，应包括大小写字母和数字，并大于8位",
		})
		return
	}

	//校验两次密码
	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "两次密码不一致",
		})
		return
	}

	//ServiceMix
	err = h.svc.Signup(ctx, Domain.User{
		Phone:    req.Phone,
		Password: req.Password,
	})
	//错误处理
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{Code: 2})
	case errors.Is(err, Service.ErrUserUnique):
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 4,
			Msg:  "邮箱已注册",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "服务器出错",
		})
	}

}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return Constant.RcJwtKey, nil
	})
	//解析失败，401,未授权
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// jwt没承诺非法就返回错误，加入校验保底
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//校验 SSID
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "OK",
	})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "已登出",
	})
}
