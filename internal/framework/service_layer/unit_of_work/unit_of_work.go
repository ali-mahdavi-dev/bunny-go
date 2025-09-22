package unit_of_work

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter/repository"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/types"
)

type PGUnitOfWork interface {
	adapter.UnitOfWork
	CollectNewEvents(eventCh chan<- any)
	User() repository.UserRepository
}

type pgUnitOfWork struct {
	adapter.UnitOfWork
	repositories []adapter.SeenedRepository

	// repositories
	user repository.UserRepository
}

func New(db *gorm.DB) PGUnitOfWork {
	uow := &pgUnitOfWork{
		UnitOfWork: adapter.NewBaseUnitOfWork(db),
	}

	return uow
}

func (uow *pgUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	err := uow.UnitOfWork.Do(ctx, fc)
	if err != nil {
		return err
	}

	return nil
}

func (uow *pgUnitOfWork) CollectNewEvents(eventCh chan<- any) {
	var wg sync.WaitGroup

	for _, repo := range uow.repositories {
		wg.Go(func() {
			for _, entity := range repo.Seen() {
				for _, event := range entity.Event() {
					fmt.Println("...event: ", event)
					eventCh <- event
				}
			}
		})
	}

	wg.Wait()
	uow.clearRepo()
}

func (uow *pgUnitOfWork) clearRepo() {
	uow.user = nil
}

func (uow *pgUnitOfWork) User() repository.UserRepository {
	if uow.user == nil {
		uow.user = repository.NewUserRepository(uow.GetSession())
		uow.repositories = append(uow.repositories, uow.user)
	}

	return uow.user
}
