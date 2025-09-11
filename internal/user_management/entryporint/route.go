package entryporint

import (
	"github.com/gin-gonic/gin"

	"bunny-go/internal/user_management/entryporint/controller"
)

type UserManagementRouter struct {
	User *controller.UserController
}

func NewUserManagementRouter(router *gin.Engine, controller UserManagementRouter) {
	userRoute := router.Group("/v1/user")
	{
		userRoute.POST("", controller.User.CreateUserController)
	}
}
