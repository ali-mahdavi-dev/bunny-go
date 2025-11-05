package account

import (
	"shikposh-backend/config"
	"shikposh-backend/internal/account/adapter"
	"shikposh-backend/internal/account/entryporint/handler"
	"shikposh-backend/internal/account/service_layer/command_handler"
	"shikposh-backend/internal/framework/service_layer/unit_of_work"
	"shikposh-backend/pkg/framework/infrastructure/logging"
	commandeventhandler "shikposh-backend/pkg/framework/service_layer/command_event_handler"
	"shikposh-backend/pkg/framework/service_layer/messagebus"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Bootstrap(router fiber.Router, db *gorm.DB, cfg *config.Config, logger logging.Logger) error {
	logging.Info("Initializing account module").Log()

	logging.Debug("Creating unit of work").Log()
	uow := unit_of_work.New(db)

	logging.Debug("Creating message bus").Log()
	bus := messagebus.NewMessageBus(uow)

	logging.Debug("Initializing avatar generator").Log()
	ag, err := adapter.NewAvatarGenerator(AssetsFS)
	if err != nil {
		logging.Error("Failed to initialize avatar generator").WithError(err).Log()
		return err
	}

	logging.Debug("Creating user handler").Log()
	userHandler := command_handler.NewUserHandler(uow, cfg)

	logging.Debug("Creating user controller").Log()
	userController := handler.NewUserController(bus, ag, userHandler)

	logging.Debug("Registering routes").Log()
	userController.RegisterRoutes(router)

	logging.Debug("Registering command and event handlers").Log()
	bus.AddHandler(
		commandeventhandler.NewCommandHandler(userHandler.RegisterHandler),
	)
	bus.AddHandler(
		commandeventhandler.NewCommandHandler(userHandler.LogoutHandler),
	)
	bus.AddHandlerEvent(
		commandeventhandler.NewEventHandler(userHandler.RegisterEvent),
	)

	logging.Info("Account module bootstrapped successfully").Log()
	return nil
}
