package unit_of_work

import (
	"context"
	"sync"
	"sync/atomic"

	"gorm.io/gorm"

	"shikposh-backend/internal/account/adapter/repository"
	"shikposh-backend/pkg/framework/adapter"
	"shikposh-backend/pkg/framework/infrastructure/logging"
	"shikposh-backend/pkg/framework/service_layer/types"
)

// PGUnitOfWork extends the base UnitOfWork with PostgreSQL-specific functionality
type PGUnitOfWork interface {
	adapter.UnitOfWork
	CollectNewEvents(ctx context.Context, eventCh chan<- any)
	User(ctx context.Context) repository.UserRepository
	Token(ctx context.Context) repository.TokenRepository
}

type pgUnitOfWork struct {
	adapter.UnitOfWork
	db           *gorm.DB
	repositories map[string]adapter.SeenedRepository
	mu           sync.RWMutex
}

// New creates a new PostgreSQL UnitOfWork instance
func New(db *gorm.DB) PGUnitOfWork {
	return &pgUnitOfWork{
		UnitOfWork:   adapter.NewBaseUnitOfWork(db),
		db:           db,
		repositories: make(map[string]adapter.SeenedRepository),
	}
}

// Do executes a function within a database transaction and collects events
func (uow *pgUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	return uow.UnitOfWork.Do(ctx, fc)
}

// CollectNewEvents collects domain events from all repositories that were accessed
// during the transaction. Events are sent to the provided channel.
func (uow *pgUnitOfWork) CollectNewEvents(ctx context.Context, eventCh chan<- any) {
	uow.mu.RLock()
	repos := make([]adapter.SeenedRepository, 0, len(uow.repositories))
	for _, repo := range uow.repositories {
		repos = append(repos, repo)
	}
	uow.mu.RUnlock()

	if len(repos) == 0 {
		logging.Debug("No repositories accessed, no events to collect").Log()
		return
	}

	logging.Debug("Collecting domain events").
		WithInt("repository_count", len(repos)).
		Log()

	var wg sync.WaitGroup
	var eventCount int64
	for _, repo := range repos {
		wg.Add(1)
		go func(r adapter.SeenedRepository) {
			defer wg.Done()
			for _, entity := range r.Seen() {
				for _, event := range entity.Event() {
					select {
					case eventCh <- event:
						atomic.AddInt64(&eventCount, 1)
						logging.Debug("Domain event collected and sent").Log()
					case <-ctx.Done():
						logging.Warn("Context cancelled while collecting events").
							WithError(ctx.Err()).
							Log()
						return
					}
				}
			}
		}(repo)
	}

	wg.Wait()

	logging.Debug("Finished collecting domain events").
		WithInt64("total_events", atomic.LoadInt64(&eventCount)).
		Log()

	uow.clearRepositories()
}

// clearRepositories clears all cached repositories after event collection
func (uow *pgUnitOfWork) clearRepositories() {
	uow.mu.Lock()
	defer uow.mu.Unlock()
	uow.repositories = make(map[string]adapter.SeenedRepository)
}

func (uow *pgUnitOfWork) User(ctx context.Context) repository.UserRepository {
	session := uow.UnitOfWork.GetSession(ctx)
	userRepo := repository.NewUserRepository(session)
	uow.mu.Lock()
	defer uow.mu.Unlock()

	key := "user"
	uow.repositories[key] = userRepo

	return userRepo
}

func (uow *pgUnitOfWork) Token(ctx context.Context) repository.TokenRepository {
	session := uow.UnitOfWork.GetSession(ctx)
	tokenRepo := repository.NewTokenRepository(session)
	uow.mu.Lock()
	defer uow.mu.Unlock()

	key := "token"
	uow.repositories[key] = tokenRepo

	return tokenRepo
}
