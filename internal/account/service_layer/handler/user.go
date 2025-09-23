package handler

import (
	"context"
	"fmt"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter/repository"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/commands"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/entity"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/events"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors/phrases"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
	"github.com/hashicorp/go-uuid"
)

type UserHandler struct {
	uow unit_of_work.PGUnitOfWork
}

func NewUserHandler(uow unit_of_work.PGUnitOfWork) *UserHandler {
	return &UserHandler{uow: uow}
}

func (h *UserHandler) Register(ctx context.Context, cmd *commands.RegisterUser) error {
	return h.uow.Do(ctx, func(ctx context.Context) error {
		user, err := h.uow.User().FindByUserName(ctx, cmd.UserName)
		if err != nil {
			if err != repository.ErrUserNotFound {
				return cerrors.BadRequest(phrases.UserNotFound)
				// return fmt.Errorf("UserHandler.Register error checking existing username: %w", err)
			}
		} else {
			return fmt.Errorf("username %s is already taken", user.UserName)
		}

		a, _ := uuid.GenerateUUID()
		err = h.uow.User().Save(ctx, entity.NewUser(
			cmd.AvatarIdentifier,
			cmd.UserName+string(a),
			cmd.FirstName,
			cmd.LastName,
			cmd.Email+string(a),
			cmd.Password,
		))
		if err != nil {
			return fmt.Errorf("UserHandler.Register error saving user: %w", err)
		}

		return nil
	})
}

func (h *UserHandler) RegisterEvent(ctx context.Context, cmd *events.RegisterUserEvent) error {
	a, _ := uuid.GenerateUUID()
	err := h.uow.User().Save(ctx, entity.NewUser(
		cmd.AvatarIdentifier,
		cmd.UserName+string(a),
		cmd.FirstName,
		cmd.LastName,
		cmd.Email+string(a),
		cmd.Password,
	))
	fmt.Println("...RegisterEvent: ", a)
	fmt.Println("...RegisterEvent error: ", err)

	return nil
}
