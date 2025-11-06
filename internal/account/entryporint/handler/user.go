package handler

import (
	"image/png"

	"shikposh-backend/internal/account/adapter"
	"shikposh-backend/internal/account/domain/commands"
	httpapi "shikposh-backend/pkg/framework/api/http"
	"shikposh-backend/pkg/framework/errors"
	"shikposh-backend/pkg/framework/errors/phrases"
	"shikposh-backend/pkg/framework/service_layer/messagebus"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cast"
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

func (u *UserController) RegisterRoutes(r fiber.Router) {
	publicRoute := r.Group("/api/v1/public")
	{
		publicRoute.Post("/avatar/:id", u.GenerateAvatarHandler)
		publicRoute.Post("/register", u.Register)
		publicRoute.Post("/login", u.Login)
		publicRoute.Post("/logout", u.Logout)
	}
}

func (u *UserController) GenerateAvatarHandler(c fiber.Ctx) error {
	identifier := c.Params("id")

	img, err := u.ag.Generate(identifier)
	if err != nil {
		return httpapi.ResError(c, err)
	}

	// Set headers for PNG response
	c.Set("Content-Type", "image/png")

	// Encode image directly to response
	if err := png.Encode(c.Response().BodyWriter(), img); err != nil {
		return httpapi.ResError(c, err)
	}

	return nil
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Handles user registration by parsing the request body and invoking the registration command.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		commands.RegisterUser	true	"RegisterUser request"
//	@Success		200		{object}	httpapi.ResponseResult	"Registration successful"
//	@Failure		400		{object}	httpapi.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		409		{object}	httpapi.ResponseResult	"User already exists"
//	@Failure		422		{object}	httpapi.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500		{object}	httpapi.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/register [post]
func (u *UserController) Register(c fiber.Ctx) error {
	ctx := c.Context()
	cmd := new(commands.RegisterUser)

	if err := httpapi.ParseJSON(c, cmd); err != nil {
		return httpapi.ResError(c, err)
	}

	result, err := u.bus.Handle(ctx, cmd)
	if err != nil {
		return httpapi.ResError(c, err)
	}

	if result != nil {
		return httpapi.ResSuccess(c, result)
	}

	return httpapi.ResOK(c)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticates a user and returns an access token.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		commands.LoginUser		true	"LoginUser"
//	@Success		200		{object}	map[string]string		"Access token"
//	@Failure		400		{object}	httpapi.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401		{object}	httpapi.ResponseResult	"Authentication failed"
//	@Failure		422		{object}	httpapi.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500		{object}	httpapi.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/login [post]
func (u *UserController) Login(c fiber.Ctx) error {
	ctx := c.Context()
	cmd := new(commands.LoginUser)

	if err := httpapi.ParseJSON(c, cmd); err != nil {
		return httpapi.ResError(c, err)
	}

	result, err := u.bus.Handle(ctx, cmd)
	if err != nil {
		return httpapi.ResError(c, err)
	}

	return httpapi.ResSuccess(c, result)
}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logs out the authenticated user.
//	@Description	Example success response: {"success": true}
//	@Description	Example error response: {"success": false, "error": {"code": "USER_NOT_FOUND", "message": "User not found", "status": "Not Found"}}
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httpapi.ResponseResult	"Logout completed successfully"
//	@Failure		400	{object}	httpapi.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401	{object}	httpapi.ResponseResult	"User not authenticated"
//	@Failure		404	{object}	httpapi.ResponseResult	"User not found"
//	@Failure		422	{object}	httpapi.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500	{object}	httpapi.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/logout [post]
func (u *UserController) Logout(c fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Get("user_id")
	if userID == "" {
		return httpapi.ResError(c, errors.NotFound(phrases.UserNotFound))
	}

	cmd := new(commands.Logout)
	cmd.UserID = cast.ToUint64(userID)

	_, err := u.bus.Handle(ctx, cmd)
	if err != nil {
		return httpapi.ResError(c, err)
	}

	return httpapi.ResOK(c)
}
