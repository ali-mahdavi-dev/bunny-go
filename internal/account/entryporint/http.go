package entryporint

import (
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint/handler"
	"github.com/gofiber/fiber/v2"
)

type UserManagementRouter struct {
	User *handler.UserController
}

func NewAccountRouter(router fiber.Router, controller UserManagementRouter) {
	controller.User.RegisterRoutes(router)
}
