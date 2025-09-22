package repository

import (
	"context"
	"errors"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/entity"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/adapter"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	adapter.BaseRepository[*entity.User]
	FindByUserName(ctx context.Context, username string) (*entity.User, error)
	FindByUsernameExcludingID(ctx context.Context, username string, Id uint) (*entity.User, error)
}

type userGormRepository struct {
	adapter.BaseRepository[*entity.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userGormRepository{
		BaseRepository: adapter.NewGormRepository[*entity.User](db),
		db:             db,
	}
}

func (u *userGormRepository) Model(ctx context.Context) *gorm.DB {
	return u.db.WithContext(ctx).Model(&entity.User{})
}

func (u *userGormRepository) FindByUsernameExcludingID(ctx context.Context, username string, id uint) (*entity.User, error) {
	var user = new(entity.User)
	err := u.Model(ctx).Where("user_name = ? and id != ? and deleted_at is null", username, id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (u *userGormRepository) FindByUserName(ctx context.Context, username string) (*entity.User, error) {
	usesr, err := u.FindByField(ctx, "user_name", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return usesr, nil
}
