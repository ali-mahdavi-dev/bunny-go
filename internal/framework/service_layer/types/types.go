package types

import (
	"context"
	"sync"
)

type RedisUseCase func(ctx context.Context) (interface{}, error)
type UowUseCase func(ctx context.Context) error

type EventUC struct {
	Channel chan any
	Wg    sync.WaitGroup
}
