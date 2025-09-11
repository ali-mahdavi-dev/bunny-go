package controller

import (
	"github.com/gin-gonic/gin"

	"bunny-go/internal/user_management/domain"
	"bunny-go/pkg/framwork/service_layer/messagebus"
	"bunny-go/pkg/ginx"
)

type UserController struct {
	bus messagebus.MessageBus
}

func NewUserController(bus messagebus.MessageBus) *UserController {
	return &UserController{
		bus: bus,
	}
}

func (u *UserController) CreateUserController(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := new(domain.CreateUserCommand)
	if err := ginx.ParseJSON(c, cmd); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := u.bus.Handle(ctx, cmd)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResOK(c)
}
