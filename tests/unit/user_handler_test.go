package unit

import (
	"context"
	"testing"

	cerrors "github.com/ali-mahdavi-dev/bunny-go/internal/framwork/errors"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain/entities"
	"github.com/ali-mahdavi-dev/bunny-go/tests/mocks"

	"github.com/stretchr/testify/assert"
)

var bus = mocks.MockUserManagementBootstrapTestApp()

func TestAddUser(t *testing.T) {
	command, ctx := CreateUserCommandCreationMethod("NewAli", 0)

	result, err := bus.Handle(ctx, command)
	user, ok := result.(*entities.User)
	assert.Nil(t, err)
	assert.True(t, true)
	assert.True(t, ok)
	assert.Equal(t, user.UserName, command.UserName)
	assert.Equal(t, user.Age, command.Age)
}

func TestForUserExisting(t *testing.T) {
	command, ctx := CreateUserCommandCreationMethod("", 0)

	result, err := bus.Handle(ctx, command)

	assert.Equal(t, err, cerrors.BadRequest("User.AlreadyExists"))
	assert.Nil(t, result)
}

func CreateUserCommandCreationMethod(userName string, age int) (*domain.CreateUserCommand, context.Context) {
	if userName == "" {
		userName = "ali"
	}
	if age == 0 {
		age = 20
	}
	ctx := context.Background()
	command := &domain.CreateUserCommand{
		UserName: userName,
		Age:      age,
	}
	return command, ctx
}
