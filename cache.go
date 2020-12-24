package cache

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration)
	Del(ctx context.Context, key string)
}
