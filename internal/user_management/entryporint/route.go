package entryporint

import (
	"github.com/gin-gonic/gin"

	"bunny-go/internal/user_management/entryporint/controller"
)

type UserManagementRouter struct {
	User *controller.UserController
}

func NewUserManagementRouter(router *gin.Engine, controller UserManagementRouter) {
	versionRoute := router.Group("/api/v1")
	{
		versionRoute.POST("/avatar/:id", controller.User.GenerateAvatarHandler)
	}
}
