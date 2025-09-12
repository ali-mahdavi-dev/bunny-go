package types

import (
	"context"
)

type RedisUseCase func(ctx context.Context) (interface{}, error)
type UowUseCase func(ctx context.Context) error
