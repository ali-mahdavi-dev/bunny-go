package unit_of_work

import (
	"sync"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter/repository"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/adapter"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
)

type PGUnitOfWork interface {
	adapter.UnitOfWork
	CollectNewEvents(eventCh chan<- any)
	User() repository.UserRepository
	Token() repository.TokenRepository
}

type pgUnitOfWork struct {
	adapter.UnitOfWork
	log          logging.Logger
	repositories []adapter.SeenedRepository

	// repositories
	user  repository.UserRepository
	token repository.TokenRepository
}

func New(db *gorm.DB, logInstans logging.Logger) PGUnitOfWork {
	uow := &pgUnitOfWork{
		UnitOfWork: adapter.NewBaseUnitOfWork(db),
		log:        logInstans,
	}

	return uow
}

func (uow *pgUnitOfWork) CollectNewEvents(eventCh chan<- any) {
	var wg sync.WaitGroup

	for _, repo := range uow.repositories {
		wg.Go(func() {
			for _, entity := range repo.Seen() {
				for _, event := range entity.Event() {
					uow.log.Info(logging.Internal, logging.Event, "send event", map[logging.ExtraKey]interface{}{
						logging.EventExtraKey: event,
					})
					eventCh <- event
				}
			}
		})
	}

	wg.Wait()
	if err := uow.Commit(); err != nil {
		
	}
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

func (uow *pgUnitOfWork) Token() repository.TokenRepository {
	if uow.token == nil {
		uow.token = repository.NewTokenRepository(uow.GetSession())
	}

	return uow.token
}
