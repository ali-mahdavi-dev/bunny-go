package handler

import (
	"context"

	"gorm.io/gorm"

	"bunny-go/internal"
	"bunny-go/internal/user_management/domain"
	"bunny-go/internal/user_management/domain/entities"
	"bunny-go/pkg/framwork/errors"
	"bunny-go/pkg/framwork/helpers/is"
)

type UserCommandHandler struct {
	uow internal.UnitOfWorkImp
}

func NewUserCommandHandler(uow internal.UnitOfWorkImp) *UserCommandHandler {
	return &UserCommandHandler{uow: uow}
}

func (u UserCommandHandler) CreateUserHandle(ctx context.Context, cmd *domain.CreateUserCommand) error {

	err := u.uow.Do(ctx, func(ctx context.Context, tx *gorm.DB) error {
		
		_, err := u.uow.User().FindByUserName(ctx, cmd.UserName)
		if !is.Error(err, gorm.ErrRecordNotFound) {
			return errors.BadRequest("User.AlreadyExists")
		}

		user, err := entities.NewUser(cmd.UserName, cmd.Age, cmd.Amount)
		if !is.Empty(err) {
			return err
		}

		err = u.uow.User().Save(ctx, user)
		if !is.Empty(err) {
			return errors.BadRequest("Operation.CanNot")
		}
		return nil
	})

	return err
}
