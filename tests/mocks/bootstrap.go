package mocks

import (
	"os/user"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/messagebus"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain"

	"gorm.io/gorm"
)

func MockUserManagementBootstrapTestApp() *FakeMessageBus {
	bus := NewFakeMessageBus(NewFakeUnitOfWork())
	bus.Register(domain.CreateUserCommand{}, user.NewCreateUserCommandHandler(bus.Uow))
	return bus
}

func SqliteUserManagementBootstrapTestApp(db *gorm.DB) *messagebus.MessageBus {
	return user_management.Bootstrap(db)
}
