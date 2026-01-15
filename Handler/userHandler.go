package Handler

import (
	"DentTick/Domain"
	"DentTick/Service"
	"errors"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	passwordRexExp *regexp.Regexp
	phoneRexExp    *regexp.Regexp
	svc            Service.UserService
}

func NewUserHandler(svc Service.UserService) *UserHandler {
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
	}
}
func (h *UserHandler) RegisterRoute(server *gin.Engine) {
	user := server.Group("/users")
	user.POST("", h.Signup)
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

	// 校验邮箱格式
	isEmailTrue, err := h.phoneRexExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isEmailTrue {
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
