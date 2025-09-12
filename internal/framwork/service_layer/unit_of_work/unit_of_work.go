package unit_of_work

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/adapter"
	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/types"
	"github.com/ali-mahdavi-dev/bunny-go/internal/user_management/adapter/repository"
)

type PGUnitOfWork interface {
	adapter.UnitOfWork
	CollectNewEvents(eventCh chan<- commandeventhandler.EventHandler)
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

func (uow *pgUnitOfWork) CollectNewEvents(eventCh chan<- commandeventhandler.EventHandler) {
	var wg sync.WaitGroup

	for _, repo := range uow.repositories {
		wg.Add(1)
		go func(r adapter.SeenedRepository) {
			defer wg.Done()
			for _, entity := range r.Seen() {
				for _, event := range entity.Event() {
					eventCh <- event
				}
			}
		}(repo)
	}

	wg.Wait()
	uow.clear()
}

func (uow *pgUnitOfWork) clear() {
	uow.user = nil
}

func (uow *pgUnitOfWork) User() repository.UserRepository {
	if uow.user == nil {
		uow.user = repository.NewUserRepository(uow.GetSession())
		uow.repositories = append(uow.repositories, uow.user)
	}

	return uow.user
}
