package handler

import (
	"context"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/unit_of_work"
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
	h.uow.User().Save(ctx, entity.NewUser(cmd.AvatarIdentifier,
		cmd.UserName,
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		cmd.Password))

	return nil
}
