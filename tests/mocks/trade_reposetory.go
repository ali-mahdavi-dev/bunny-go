package mocks

import (
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/adapter/repositories"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/domain/entities"
)

type FakeTradeRepository struct {
	FakRepository[*entities.Trade]
}

func NewFakeTradeRepository() repositories.TradeRepository {
	userRepo := &FakeTradeRepository{
		FakRepository: *NewFakeRepository[*entities.Trade](),
	}

	return userRepo
}
