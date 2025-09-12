package user_management

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/unit_of_work"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/entryporint"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/entryporint/controller"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/service_layer/handler"
)

func Bootstrap(router *gin.Engine, db *gorm.DB) error {
	eventCh := make(chan commandeventhandler.EventHandler, 100)
	uow := unit_of_work.New(db)
	bus := messagebus.NewMessageBus(uow, eventCh)
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

	return nil
}
