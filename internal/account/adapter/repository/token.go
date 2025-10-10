package repository

import (
	"context"
	"errors"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/entity"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/adapter"
	"gorm.io/gorm"
)

var ErrTokenNotFound = errors.New("Token not found")

type TokenRepository interface {
	adapter.BaseRepository[*entity.Token]
	FindByUserID(ctx context.Context, userID uint64) (*entity.Token, error)
}

type tokenGormRepository struct {
	adapter.BaseRepository[*entity.Token]
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenGormRepository{
		BaseRepository: adapter.NewGormRepository[*entity.Token](db),
		db:             db,
	}
}

func (u *tokenGormRepository) Model(ctx context.Context) *gorm.DB {
	return u.db.WithContext(ctx).Model(&entity.Token{})
}

func (u *tokenGormRepository) FindByUserID(ctx context.Context, userID uint64) (*entity.Token, error) {
	token, err := u.FindByField(ctx, "user_id", userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}

		return nil, err
	}

	return token, nil
}
