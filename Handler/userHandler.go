package Handler

import (
	"DentTick/Service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service Service.UserService
}

func NewUserHandler(svc Service.UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}
func (h *UserHandler) RegisterRoute(server *gin.Engine) {
	user := server.Group("/users")
	user.POST("", h.Signup)
}

func (h *UserHandler) Signup(ctx *gin.Context) {

}
