package adapter

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/types"
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

func (uow *BaseUnitOfWork) Begin() error {
	uow.DB = uow.DB.Begin()
	if uow.DB.Error != nil {
		return uow.DB.Error
	}

	return nil
}

func (uow *BaseUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	err := uow.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	err = fc(ctx)
	defer func() {
		errRecover := recover()
		if err != nil || errRecover != nil {
			e := uow.Rollback()
			if e != nil {
				log.Println("error rolling back transaction:", e)
			}
		}
	}()

	return err
}

func (uow *BaseUnitOfWork) Commit() error {
	return uow.DB.Commit().Error
}

func (uow *BaseUnitOfWork) Rollback() error {
	return uow.DB.Rollback().Error
}
