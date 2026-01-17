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
	user.POST("/login", h.Login)
	user.POST("logout", h.Logout)
	user.POST("/edit", h.EditProfile)
	user.GET("/profile", h.Profile)
	user.GET("/refresh_token", h.RefreshToken)

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	// 从上下文中解析用户信息
	uc, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	UC, ok := uc.(ijwt.UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 查询用户资料
	u, err := h.svc.GetProfile(ctx, UC.Uid)
	if err != nil {
		switch {
		case errors.Is(err, Service.ErrUserNotFound):
			ctx.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "user not found",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, Result{
				Code: 5,
				Msg:  "server error",
			})
		}
		return
	}

	type profileResp struct {
		Name     string `json:"name"`
		NickName string `json:"nickname"`
		Info     string `json:"info"`
		Avatar   string `json:"avatar"`
		Phone    string `json:"phone"`
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Data: profileResp{
			Name:     u.Name,
			NickName: u.NickName,
			Info:     u.Info,
			Avatar:   u.Avatar,
			Phone:    u.Phone,
		},
	})
}

func (h *UserHandler) EditProfile(ctx *gin.Context) {

	type editProfileReq struct {
		Name     string `json:"name"`
		NickName string `json:"nickname"`
		Info     string `json:"info"`
		Avatar   string `json:"avatar"`
	}
	var req editProfileReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	// 从上下文中解析用户信息
	uc, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	UC, ok := uc.(ijwt.UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	uid := UC.Uid
	// 只允许编辑资料字段
	err = h.svc.EditProfile(ctx, Domain.User{
		Id:       uid,
		Name:     req.Name,
		NickName: req.NickName,
		Info:     req.Info,
		Avatar:   req.Avatar,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "OK",
	})
}

func (h *UserHandler) Signup(ctx *gin.Context) {
	type signUpReq struct {
		Phone           string `json:"phone"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_passwd"`
	}

	var req signUpReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	h.CheckPhoneForm(ctx, req.Phone, req.Password)
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

func (h *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	h.CheckPhoneForm(ctx, req.Phone, req.Password)
	// service: Login
	u, err := h.svc.Login(ctx, req.Phone, req.Password)

	switch {
	case err == nil:
		h.SetLoginToken(ctx, u.Id)
		ctx.JSON(http.StatusOK, Result{Code: 2, Msg: "登录成功"})
	case errors.Is(err, Service.ErrInvalidUserOrPassword):
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "手机号或密码错误"})
	default:
		ctx.JSON(http.StatusInternalServerError, Result{Code: 5, Msg: "服务器出错"})
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

func (h *UserHandler) CheckPhoneForm(ctx *gin.Context, phone, password string) {
	// 校验手机格式
	isPhoneTrue, err := h.phoneRexExp.MatchString(phone)
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
	isPasswordTrue, err := h.passwordRexExp.MatchString(password)
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
}
