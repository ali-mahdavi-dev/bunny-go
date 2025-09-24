package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter/repository"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/commands"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/entity"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/events"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/api/jwt"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors/phrases"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
)

type UserHandler struct {
	uow unit_of_work.PGUnitOfWork
	cfg *config.Config
}

func NewUserHandler(uow unit_of_work.PGUnitOfWork, cfg *config.Config) *UserHandler {
	return &UserHandler{uow: uow, cfg: cfg}
}

func (h *UserHandler) RegisterHandler(ctx context.Context, cmd *commands.RegisterUser) error {
	return h.uow.Do(ctx, func(ctx context.Context) error {
		_, err := h.uow.User().FindByUserName(ctx, cmd.UserName)
		if err != nil {
			if err != repository.ErrUserNotFound {
				return fmt.Errorf("UserHandler.Register error checking existing username: %w", err)
			}
		} else {
			return cerrors.BadRequest(phrases.UserAlreadyExists)
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

		return nil
	})
}

func (h *UserHandler) LogoutHandler(ctx context.Context, cmd *commands.Logout) error {
	token, err := h.uow.Token().FindByUserID(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrTokenNotFound) {
			return cerrors.NotFound(phrases.UserNotFound)
		}

		return fmt.Errorf("UserHandler.LogoutHandler failed to get token by userID: %w", err)
	}

	if err := h.uow.Token().Remove(ctx, token, false); err != nil {
		return fmt.Errorf("UserHandler.LoginUseCase failed to remove existing token: %w", err)
	}

	return nil
}

func (h *UserHandler) LoginUseCase(ctx context.Context, cmd *commands.LoginUser) (string, error) {
	var accessToken string
	err := h.uow.Do(ctx, func(ctx context.Context) error {
		user, err := h.uow.User().FindByUserName(ctx, cmd.UserName)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return cerrors.NotFound(phrases.UserNotFound)
			}

			return fmt.Errorf("UserHandler.LoginUseCase fail get user by username: %w", err)
		}
		token, err := h.uow.Token().FindByUserID(ctx, user.ID)
		if err != nil && !errors.Is(err, repository.ErrTokenNotFound) {
			return fmt.Errorf("UserHandler.LoginUseCase failed to get token by userID: %w", err)
		}

		if token != nil {
			if err := h.uow.Token().Remove(ctx, token, false); err != nil {
				return fmt.Errorf("UserHandler.LoginUseCase failed to remove existing token: %w", err)
			}
		}

		accessToken, err = jwt.GenerateToken(h.cfg.JWT.AccessTokenExpireDuration, h.cfg.JWT.Secret, user.ID)
		if err != nil {
			return fmt.Errorf("UserHandler.LoginUseCase fail generate token: %w", err)
		}

		err = h.uow.Token().Save(ctx, entity.NewToken(accessToken, user.ID))
		if err != nil {
			return fmt.Errorf("UserHandler.LoginUseCase fail save token to db: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("UserHandler.LoginUseCase fail transaction: %w", err)

	}
	return accessToken, nil
}

func (h *UserHandler) RegisterEvent(ctx context.Context, cmd *events.RegisterUserEvent) error {

	return nil
}
