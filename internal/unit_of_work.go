package internal

import (
	"context"

	"gorm.io/gorm"

	"bunny-go/internal/user_management/adapter/repositories"
	"bunny-go/pkg/framwork/service_layer/types"
	"bunny-go/pkg/framwork/service_layer/unit_of_work"
)

var (
	_ UnitOfWorkImp = &unitOfWorkImp{}
)

type UnitOfWorkImp interface {
	unit_of_work.UnitOfWork
	Register()
	User() repositories.UserRepository
	Trade() repositories.TradeRepository
}

type unitOfWorkImp struct {
	unit_of_work.UnitOfWork
	user  repositories.UserRepository
	trade repositories.TradeRepository
}

func NewGormUnitOfWorkImp(db *gorm.DB) UnitOfWorkImp {
	return &unitOfWorkImp{UnitOfWork: unit_of_work.NewGormUnitOfWork(db)}
}
func (uow *unitOfWorkImp) Register() {
	uow.user = repositories.NewUserGormRepository(uow.UnitOfWork.GetSession())
	uow.trade = repositories.NewTradeGormRepository(uow.UnitOfWork.GetSession())
}
func (uow *unitOfWorkImp) Do(ctx context.Context, fc types.UowUseCase) error {
	defer func() {
		if recover() != nil {
			_ = uow.Rollback()
		}
	}()

	if err := uow.Begin(); err != nil {
		return err
	}
	// initial repository
	uow.Register()

	if err := uow.UnitOfWork.Do(ctx, fc); err != nil {
		return err
	}
	if err := uow.Commit(); err != nil {
		return err
	}

	return nil
}

func (uow *unitOfWorkImp) User() repositories.UserRepository {
	return uow.user
}

func (uow *unitOfWorkImp) Trade() repositories.TradeRepository {
	return uow.trade
}
