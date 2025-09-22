package account

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint/controller"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/service_layer/handler"
	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
)

func Bootstrap(router *gin.Engine, db *gorm.DB) error {
	uow := unit_of_work.New(db)
	bus := messagebus.NewMessageBus(uow)
	ag, err := adapter.NewAvatarGenerator()
	if err != nil {
		return err
	}

	// init controller
	userController := controller.NewUserController(bus, ag)

	userHandler := handler.NewUserHandler(uow)

	// init router
	entryporint.NewUserManagementRouter(router, entryporint.UserManagementRouter{
		User: userController,
	})

	// init handler
	bus.AddHandler(
		// avatar
		commandeventhandler.NewCommandHandler(userHandler.Register),
	)
	bus.AddHandlerEvent(
		// avatar
		commandeventhandler.NewEventHandler(userHandler.RegisterEvent),
	)

	return nil
}
