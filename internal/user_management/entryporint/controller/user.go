package controller

import (
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"

	"bunny-go/internal/framwork/service_layer/messagebus"
	"bunny-go/internal/user_management/adapter"
)

type UserController struct {
	bus messagebus.MessageBus
	ag  *adapter.AvatarGenerator
}

func NewUserController(bus messagebus.MessageBus, ag *adapter.AvatarGenerator) *UserController {
	return &UserController{
		bus: bus,
		ag:  ag,
	}
}

// Gin handler to generate avatar and return as PNG
func (s *UserController) GenerateAvatarHandler(c *gin.Context) {
	identifier := c.Param("id")

	img, err := s.ag.Generate(identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PNG response
	c.Header("Content-Type", "image/png")

	// Encode directly to response
	if err := png.Encode(c.Writer, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// func (u *UserController) CreateUserController(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	cmd := new(domain.GenerateAvatarCommand)
// 	if err := ginx.ParseJSON(c, cmd); err != nil {
// 		ginx.ResError(c, err)
// 		return
// 	}

// 	err := u.bus.Handle(ctx, cmd)
// 	if err != nil {
// 		fmt.Println("Error handling command:", err)
// 		ginx.ResError(c, fmt.Errorf("failed to create user"))
// 		return
// 	}

// 	ginx.ResOK(c)
// }
