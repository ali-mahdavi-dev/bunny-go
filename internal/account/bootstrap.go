package account

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/endpoint"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/endpoint/controller"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/service_layer/handler"
	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
)

func Bootstrap(router *gin.Engine, db *gorm.DB) error {
	logging.Info("Initializing account module").Log()

	logging.Debug("Creating unit of work").Log()
	uow := unit_of_work.New(db)

	logging.Debug("Creating message bus").Log()
	bus := messagebus.NewMessageBus(uow)

	logging.Debug("Initializing avatar generator").Log()
	ag, err := adapter.NewAvatarGenerator()
	if err != nil {
		logging.Error("Failed to initialize avatar generator").
			WithError(err).
			Log()
		return err
	}

	logging.Debug("Creating user controller").Log()
	userController := controller.NewUserController(bus, ag)

	logging.Debug("Creating user handler").Log()
	userHandler := handler.NewUserHandler(uow)

	logging.Debug("Registering routes").Log()
	endpoint.NewUserManagementRouter(router, endpoint.UserManagementRouter{
		User: userController,
	})

	logging.Debug("Registering command and event handlers").Log()
	bus.AddHandler(
		commandeventhandler.NewCommandHandler(userHandler.Register),
	)
	bus.AddHandlerEvent(
		commandeventhandler.NewEventHandler(userHandler.RegisterEvent),
	)

	logging.Info("Account module bootstrapped successfully").Log()
	return nil
}
