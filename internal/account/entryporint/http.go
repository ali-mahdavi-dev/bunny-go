package entryporint

import (
	"shikposh-backend/internal/account/entryporint/handler"

	"github.com/gofiber/fiber/v3"
)

type UserManagementRouter struct {
	User *handler.UserController
}

func NewAccountRouter(router fiber.Router, controller UserManagementRouter) {
	controller.User.RegisterRoutes(router)
}
