package entryporint

import (
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint/routes"
	"github.com/gin-gonic/gin"
)

type UserManagementRouter struct {
	User *routes.UserController
}

func NewUserManagementRouter(router *gin.Engine, controller UserManagementRouter) {
	publicRoute := router.Group("/api/v1/public")
	{
		publicRoute.POST("/avatar/:id", controller.User.GenerateAvatarHandler)
		publicRoute.POST("/register", controller.User.Register)
		publicRoute.POST("/login", controller.User.Login)
	}

}
