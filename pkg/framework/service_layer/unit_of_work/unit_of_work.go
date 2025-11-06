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

// PGUnitOfWork extends the base UnitOfWork with PostgreSQL-specific functionality.
type PGUnitOfWork interface {
	adapter.UnitOfWork
	CollectNewEvents(ctx context.Context, eventCh chan<- any)
	User(ctx context.Context) repository.UserRepository
	Token(ctx context.Context) repository.TokenRepository
	Profile(ctx context.Context) repository.ProfileRepository
}

type pgUnitOfWork struct {
	adapter.UnitOfWork
	db           *gorm.DB
	repositories map[string]adapter.SeenedRepository
	mu           sync.RWMutex
}

// New creates a new PostgreSQL UnitOfWork instance.
func New(db *gorm.DB) PGUnitOfWork {
	return &pgUnitOfWork{
		UnitOfWork:   adapter.NewBaseUnitOfWork(db),
		db:           db,
		repositories: make(map[string]adapter.SeenedRepository),
	}
}

// Do executes a function within a database transaction.
func (uow *pgUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	return uow.UnitOfWork.Do(ctx, fc)
}

// CollectNewEvents collects domain events from repositories and sends them to the channel.
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

	allEntities := make([]adapter.Entity, 0, 10)
	for _, repo := range repos {
		entities := repo.Seen() 
		allEntities = append(allEntities, entities...)
	}

	var wg sync.WaitGroup
	var eventCount int64
	for _, entity := range allEntities {
		wg.Add(1)
		go func(e adapter.Entity) {
			defer wg.Done()
			events := e.Event()
			for _, event := range events {
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
		}(entity)
	}

	wg.Wait()

	logging.Debug("Finished collecting domain events").
		WithInt64("total_events", atomic.LoadInt64(&eventCount)).
		Log()

	uow.clearRepositories()
}

// clearRepositories clears all cached repositories after event collection.
func (uow *pgUnitOfWork) clearRepositories() {
	uow.mu.Lock()
	defer uow.mu.Unlock()
	uow.repositories = make(map[string]adapter.SeenedRepository)
}

// getOrCreateRepository returns a cached repository or creates a new one using double-check locking.
func (uow *pgUnitOfWork) getOrCreateRepository(
	ctx context.Context,
	key string,
	factory func(*gorm.DB) adapter.SeenedRepository,
) adapter.SeenedRepository {
	// Fast path: check with read lock
	uow.mu.RLock()
	if repo, ok := uow.repositories[key]; ok {
		uow.mu.RUnlock()
		return repo
	}
	uow.mu.RUnlock()

	// Slow path: acquire write lock and create repository
	uow.mu.Lock()
	defer uow.mu.Unlock()
	// Create new repository instance
	session := uow.UnitOfWork.GetSession(ctx)
	repo := factory(session)
	uow.repositories[key] = repo

	return repo
}

// User returns the UserRepository instance for the current transaction.
func (uow *pgUnitOfWork) User(ctx context.Context) repository.UserRepository {
	return uow.getOrCreateRepository(ctx, "user", func(session *gorm.DB) adapter.SeenedRepository {
		return repository.NewUserRepository(session)
	}).(repository.UserRepository)
}

// Token returns the TokenRepository instance for the current transaction.
func (uow *pgUnitOfWork) Token(ctx context.Context) repository.TokenRepository {
	return uow.getOrCreateRepository(ctx, "token", func(session *gorm.DB) adapter.SeenedRepository {
		return repository.NewTokenRepository(session)
	}).(repository.TokenRepository)
}

// Profile returns the ProfileRepository instance for the current transaction.
func (uow *pgUnitOfWork) Profile(ctx context.Context) repository.ProfileRepository {
	return uow.getOrCreateRepository(ctx, "profile", func(session *gorm.DB) adapter.SeenedRepository {
		return repository.NewProfileRepository(session)
	}).(repository.ProfileRepository)
}
