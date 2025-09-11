package mocks

import (
	"os/user"

	"bunny-go/internal/user_management"
	"bunny-go/internal/user_management/domain"
	"bunny-go/pkg/framwork/service_layer/messagebus"

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
