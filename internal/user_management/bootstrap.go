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
)

func Bootstrap(router *gin.Engine, db *gorm.DB) error {
	eventCh := make(chan commandeventhandler.EventHandler, 100)
	uow := unit_of_work.New(db)
	bus := messagebus.NewMessageBus(uow, eventCh)

	// init controller
	// userHandler, err := handler.NewUserCommandHandler("_data-dev/avatar-data")
	// if err != nil {
	// 	return fmt.Errorf("failed to create user command handler: %w", err)
	// }
	ag, err := adapter.NewAvatarGenerator("_data-dev/avatar-data")
	if err != nil {
		return err
	}
	userController := controller.NewUserController(bus, ag)

	// init router
	entryporint.NewUserManagementRouter(router, entryporint.UserManagementRouter{
		User: userController,
	})

	// init handler
	bus.AddHandler(
	// avatar
	// commandeventhandler.NewCommandHandler(userHandler.CreateUserHandle),
	)
	// fmt.Println("---------------------------- ----------------------------")
	// err:=userHandler.GenerateAndSaveAvatar("ali", "ali") // test
	// fmt.Println("err:", err)
	return nil
}
