package user_management

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bunny-go/internal"
	"bunny-go/internal/user_management/entryporint"
	"bunny-go/internal/user_management/entryporint/controller"
	"bunny-go/internal/user_management/service_layer/handler"
	"bunny-go/pkg/framwork/service_layer/messagebus"
)

func Bootstrap(router *gin.Engine, db *gorm.DB) messagebus.MessageBus {
	uow := internal.NewGormUnitOfWorkImp(db)
	bus := messagebus.NewMessageBus()

	// init controller
	userController := controller.NewUserController(bus)


    entryporint.NewUserManagementRouter(router, entryporint.UserManagementRouter{
        User: userController,
    })

    // init handler
	userHandler := handler.NewUserCommandHandler(uow)
	bus.AddHandler(
		// avatar
		messagebus.NewCommandHandler(userHandler.CreateUserHandle),
	)

	return bus
}
