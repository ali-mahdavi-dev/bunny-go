package entryporint

import (
	"github.com/gin-gonic/gin"

	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/entryporint/controller"
)

type UserManagementRouter struct {
	User *controller.UserController
}

func NewUserManagementRouter(router *gin.Engine, controller UserManagementRouter) {
	publicRoute := router.Group("/api/v1/public")
	{
		publicRoute.POST("/avatar/:id", controller.User.GenerateAvatarHandler)
		publicRoute.POST("/register", controller.User.Register)
	}
}
