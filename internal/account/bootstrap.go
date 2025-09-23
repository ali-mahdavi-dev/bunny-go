package account

import (
	"embed"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/entryporint/controller"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/service_layer/handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
)

//go:embed assets/images/*
var imagesFS embed.FS

var LogInstans logging.Logger

func Bootstrap(router *gin.Engine, db *gorm.DB, cfg *config.Config) error {
	LogInstans = logging.NewLogger(cfg)
	uow := unit_of_work.New(db, LogInstans)
	bus := messagebus.NewMessageBus(uow, LogInstans)

	ag, err := adapter.NewAvatarGenerator(imagesFS)
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
