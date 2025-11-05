package adapter

import (
	"context"

	"gorm.io/gorm"

	"shikposh-backend/pkg/framework/service_layer/types"
)

// txKey is a key type for storing transaction in context
type txKey struct{}

// UnitOfWork defines the interface for managing database transactions
type UnitOfWork interface {
	Do(ctx context.Context, fc types.UowUseCase) error
	GetSession(ctx context.Context) *gorm.DB
	Commit() error
	Rollback() error
}

// BaseUnitOfWork implements the UnitOfWork pattern with GORM
type BaseUnitOfWork struct {
	db *gorm.DB
}

// NewBaseUnitOfWork creates a new instance of BaseUnitOfWork
func NewBaseUnitOfWork(db *gorm.DB) UnitOfWork {
	return &BaseUnitOfWork{
		db: db,
	}
}

// GetSession returns the database session from context if available (transaction),
// otherwise returns the base database connection
func (uow *BaseUnitOfWork) GetSession(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return uow.db
}

// Do executes a function within a database transaction
// If the function returns an error, the transaction will be rolled back
func (uow *BaseUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Store transaction in context so GetSession can retrieve it
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fc(txCtx)
	})
}

// Commit commits the current transaction (if in manual mode)
// Note: With Transaction() method, commit is automatic on success
func (uow *BaseUnitOfWork) Commit() error {
	return uow.db.Commit().Error
}

// Rollback rolls back the current transaction (if in manual mode)
// Note: With Transaction() method, rollback is automatic on error
func (uow *BaseUnitOfWork) Rollback() error {
	return uow.db.Rollback().Error
}
