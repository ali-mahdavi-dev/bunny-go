package handler

import (
	"context"
	"fmt"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/unit_of_work"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/adapter/repository"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain/commands"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain/entity"
)

type UserHandler struct {
	uow unit_of_work.PGUnitOfWork
}

func NewUserHandler(uow unit_of_work.PGUnitOfWork) *UserHandler {
	return &UserHandler{uow: uow}
}

func (h *UserHandler) Register(ctx context.Context, cmd *commands.RegisterUser) error {

	h.uow.Do(ctx, func(ctx context.Context) error {
		user, err := h.uow.User().FindByUserName(ctx, cmd.UserName)
		if err != nil {
			if err != repository.ErrUserNotFound {
				return fmt.Errorf("UserHandler.Register error checking existing username: %w", err)
			}
		} else {
			return fmt.Errorf("username %s is already taken", user.UserName)
		}

		err = h.uow.User().Save(ctx, entity.NewUser(
			cmd.AvatarIdentifier,
			cmd.UserName,
			cmd.FirstName,
			cmd.LastName,
			cmd.Email,
			cmd.Password,
		))
		if err != nil {
			return fmt.Errorf("UserHandler.Register error saving user: %w", err)
		}

		h.uow.Commit()
		return nil
	})
	return nil
}
