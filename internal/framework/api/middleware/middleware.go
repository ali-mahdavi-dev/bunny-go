package middleware

import (
	"errors"

	"shikposh-backend/config"
	"shikposh-backend/internal/framework/service_layer/unit_of_work"

	"github.com/gofiber/fiber/v3"
)

var errFailGetTokenFromDB = errors.New("fail to get token from DB")
var errTokenDoesNotExist = errors.New("token does not exist")

type Middleware struct {
	Cfg *config.Config
	Uow unit_of_work.PGUnitOfWork
}

func NewMiddleware(cfg *config.Config, uow unit_of_work.PGUnitOfWork) *Middleware {
	return &Middleware{
		Cfg: cfg,
		Uow: uow,
	}
}

func (m *Middleware) Register(app *fiber.App) {
	app.Use(m.DefaultStructuredLogger())
}
