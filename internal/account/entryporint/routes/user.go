package routes

import (
	"image/png"
	"net/http"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/commands"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/service_layer/handler"
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
	uh  *handler.UserHandler
}

func NewUserController(bus messagebus.MessageBus, ag *adapter.AvatarGenerator, uh *handler.UserHandler) *UserController {
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

// Transfer godoc
//
//	@Summary		Transfer funds
//	@Description	Transfers money to the specified destination account. The provider is automatically determined using the given FinancialServiceProviderID. Additionally, today's and this month's transaction amounts are calculated and tracked.
//	@Tags			Banking
//	@Accept			json
//	@Produce		json
//	@Param			request	body		commands.RegisterUser	true	"Transfer request data including destination account, amount, and FinancialServiceProviderID"
//	@Success		200		{object}	ginx.ResponseResult		"Transfer completed successfully"
//	@Failure		400		{object}	ginx.ResponseResult		"Invalid request body or unknown provider"
//	@Failure		422		{object}	ginx.ResponseResult		"Unprocessable input (validation failed)"
//	@Failure		500		{object}	ginx.ResponseResult		"Internal server error"
//	@Router			/internal/api/v1/banking/transfer/auto [post]
func (u *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := new(commands.RegisterUser)
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

func (u *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := new(commands.LoginUser)
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

func (u *UserController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		ginx.ResError(c, cerrors.NotFound(phrases.UserNotFound))
		return
	}
	cmd := new(commands.Logout)
	cmd.UserID = cast.ToUint64(userID)

	err := u.bus.Handle(ctx, cmd)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResOK(c)
}
