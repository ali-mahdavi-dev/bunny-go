package handler

import (
	"image/png"
	"net/http"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/command"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/service_layer/command_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors/phrases"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type UserController struct {
	bus messagebus.MessageBus
	ag  *adapter.AvatarGenerator
	uh  *command_handler.UserHandler
}

func NewUserController(bus messagebus.MessageBus, ag *adapter.AvatarGenerator, uh *command_handler.UserHandler) *UserController {
	return &UserController{
		bus: bus,
		ag:  ag,
		uh:  uh,
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

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Handles user registration by parsing the request body and invoking the registration command.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		command.RegisterUser	true	"RegisterUser"
//	@Success		200		{object}	ginx.ResponseResult		"Registration successful"
//	@Failure		400		{object}	ginx.ResponseResult		"Invalid request body or unknown provider"
//	@Failure		422		{object}	ginx.ResponseResult		"Unprocessable input (validation failed)"
//	@Failure		500		{object}	ginx.ResponseResult		"Internal server error"
//	@Router			/api/v1/public/register [post]
func (u *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := new(command.RegisterUser)
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

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticates a user and returns an access token.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		command.LoginUser	true	"LoginUser"
//	@Success		200		{object}	map[string]string	"Access token"
//	@Failure		400		{object}	ginx.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401		{object}	ginx.ResponseResult	"Authentication failed"
//	@Failure		422		{object}	ginx.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500		{object}	ginx.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/login [post]
func (u *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := new(command.LoginUser)
	if err := ginx.ParseJSON(c, cmd); err != nil {
		ginx.ResError(c, err)
		return
	}

	token, err := u.uh.LoginUseCase(ctx, cmd)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResJSON(c, http.StatusOK, map[string]string{
		"access": token,
	})

}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logs out the authenticated user.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ginx.ResponseResult	"Logout completed successfully"
//	@Failure		400	{object}	ginx.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401	{object}	ginx.ResponseResult	"User not authenticated"
//	@Failure		422	{object}	ginx.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500	{object}	ginx.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/logout [post]
func (u *UserController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		ginx.ResError(c, cerrors.NotFound(phrases.UserNotFound))
		return
	}
	cmd := new(command.Logout)
	cmd.UserID = cast.ToUint64(userID)

	err := u.bus.Handle(ctx, cmd)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResOK(c)
}
