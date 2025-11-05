package handler

import (
	"image/png"
	"net/http"

	"shikposh-backend/internal/account/adapter"
	"shikposh-backend/internal/account/domain/commands"
	"shikposh-backend/internal/account/service_layer/command_handler"
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
	uh  *command_handler.UserHandler
}

func NewUserController(bus messagebus.MessageBus, ag *adapter.AvatarGenerator, uh *command_handler.UserHandler) *UserController {
	return &UserController{
		bus: bus,
		ag:  ag,
		uh:  uh,
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set headers for PNG response
	c.Set("Content-Type", "image/png")

	// Encode image directly to response
	if err := png.Encode(c.Response().BodyWriter(), img); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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
//	@Param			request	body		commands.RegisterUser		true	"RegisterUser"
//	@Success		200		{object}	httputils.ResponseResult	"Registration successful"
//	@Failure		400		{object}	httputils.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		422		{object}	httputils.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500		{object}	httputils.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/register [post]
func (u *UserController) Register(c fiber.Ctx) error {
	ctx := c.Context()
	cmd := new(commands.RegisterUser)

	if err := httpapi.ParseJSON(c, cmd); err != nil {
		return httpapi.ResError(c, err)
	}

	if err := u.bus.Handle(ctx, cmd); err != nil {
		return httpapi.ResError(c, err)
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
//	@Param			request	body		commands.LoginUser			true	"LoginUser"
//	@Success		200		{object}	map[string]string			"Access token"
//	@Failure		400		{object}	httputils.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401		{object}	httputils.ResponseResult	"Authentication failed"
//	@Failure		422		{object}	httputils.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500		{object}	httputils.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/login [post]
func (u *UserController) Login(c fiber.Ctx) error {
	ctx := c.Context()
	cmd := new(commands.LoginUser)

	if err := httpapi.ParseJSON(c, cmd); err != nil {
		return httpapi.ResError(c, err)
	}

	token, err := u.uh.LoginUseCase(ctx, cmd)
	if err != nil {
		return httpapi.ResError(c, err)
	}

	return httpapi.ResJSON(c, http.StatusOK, fiber.Map{
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
//	@Success		200	{object}	httputils.ResponseResult	"Logout completed successfully"
//	@Failure		400	{object}	httputils.ResponseResult	"Invalid request body or unknown provider"
//	@Failure		401	{object}	httputils.ResponseResult	"User not authenticated"
//	@Failure		422	{object}	httputils.ResponseResult	"Unprocessable input (validation failed)"
//	@Failure		500	{object}	httputils.ResponseResult	"Internal server error"
//	@Router			/api/v1/public/logout [post]
func (u *UserController) Logout(c fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Get("user_id")
	if userID == "" {
		return httpapi.ResError(c, errors.NotFound(phrases.UserNotFound))
	}

	cmd := new(commands.Logout)
	cmd.UserID = cast.ToUint64(userID)

	if err := u.bus.Handle(ctx, cmd); err != nil {
		return httpapi.ResError(c, err)
	}

	return httpapi.ResOK(c)
}
