package repositories

import (
	"bunny-go/internal/user_management/domain/entities"
	"bunny-go/pkg/framwork/adapter"

	"gorm.io/gorm"
)

type TradeRepository interface {
	adapter.BaseRepository[*entities.Trade]
}

type tradeGormRepository struct {
	adapter.BaseRepository[*entities.Trade]
	db *gorm.DB
}

func NewTradeGormRepository(db *gorm.DB) TradeRepository {
	return &tradeGormRepository{
		BaseRepository: adapter.NewGormRepository[*entities.Trade](db),
		db:             db,
	}
}
