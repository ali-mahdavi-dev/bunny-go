package adapter

import (
	"context"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/types"
)

type UnitOfWork interface {
	Do(ctx context.Context, fc types.UowUseCase) error
	GetSession() *gorm.DB
	Commit() error
	Rollback() error
}

type BaseUnitOfWork struct {
	DB *gorm.DB
}

func NewBaseUnitOfWork(db *gorm.DB) UnitOfWork {
	return &BaseUnitOfWork{
		DB: db,
	}
}

func (uow *BaseUnitOfWork) GetSession() *gorm.DB {
	return uow.DB
}

func (uow *BaseUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	err := uow.DB.Transaction(func(tx *gorm.DB) error {
		return fc(ctx)
	})

	return err
}

func (uow *BaseUnitOfWork) Commit() error {
	return uow.DB.Commit().Error
}

func (uow *BaseUnitOfWork) Rollback() error {
	return uow.DB.Rollback().Error
}
