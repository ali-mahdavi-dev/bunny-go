package repository

import (
	"context"
	"errors"

	"shikposh-backend/internal/account/domain/entity"
	"shikposh-backend/pkg/framework/adapter"

	"gorm.io/gorm"
)

var ErrProfileNotFound = errors.New("profile not found")

type ProfileRepository interface {
	adapter.BaseRepository[*entity.Profile]
	FindByUserID(ctx context.Context, userID uint64) (*entity.Profile, error)
}

type profileGormRepository struct {
	adapter.BaseRepository[*entity.Profile]
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileGormRepository{
		BaseRepository: adapter.NewGormRepository[*entity.Profile](db),
		db:             db,
	}
}

func (p *profileGormRepository) Model(ctx context.Context) *gorm.DB {
	return p.db.WithContext(ctx).Model(&entity.Profile{})
}

func (p *profileGormRepository) FindByUserID(ctx context.Context, userID uint64) (*entity.Profile, error) {
	profile, err := p.FindByField(ctx, "user_id", userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProfileNotFound
		}

		return nil, err
	}

	return profile, nil
}

