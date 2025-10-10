package entryporint

import (
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint/handler"
	"github.com/gin-gonic/gin"
)

type UserManagementRouter struct {
	User *handler.UserController
}

func NewUserManagementRouter(router *gin.Engine, controller UserManagementRouter) {
	publicRoute := router.Group("/api/v1/public")
	{
		publicRoute.POST("/avatar/:id", controller.User.GenerateAvatarHandler)
		publicRoute.POST("/register", controller.User.Register)
		publicRoute.POST("/login", controller.User.Login)
		publicRoute.POST("/logout", controller.User.Logout)
	}

}
