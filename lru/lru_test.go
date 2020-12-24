package lru_test

import (
	"context"
	"testing"
	"time"

	"github.com/x1nchen/ttlcache"
	"github.com/x1nchen/ttlcache/lru"

	"github.com/stretchr/testify/require"
)

func TestLRU(t *testing.T) {
	cache := lru.New(lru.WithMaxCache(3), lru.WithGcPeriod(1*time.Second))
	ctx := context.Background()
	//
	cache.Set(ctx, "k1", "v1", 500*time.Millisecond)
	cache.Set(ctx, "k2", "v2", 500*time.Millisecond)
	cache.Set(ctx, "k3", "v3", 500*time.Millisecond)
	cache.Set(ctx, "k4", "v4", 500*time.Millisecond)
	require.Equal(t, 3, cache.Len())
	_, err := cache.Get(ctx, "k1")
	require.EqualError(t, err, ttlcache.ErrKeyNotFound.Error())
	//
	cache.Set(ctx, "k2", "v22", 500*time.Millisecond)
	cache.Set(ctx, "k1", "v1", 500*time.Millisecond)
	v2, err := cache.Get(ctx, "k2")
	require.NoError(t, err)
	require.Equal(t, "v22", v2)
	_, err = cache.Get(ctx, "k3")
	require.EqualError(t, err, ttlcache.ErrKeyNotFound.Error())
	//
	cache.Del(ctx, "k1")
	cache.Del(ctx, "k2")
	require.Equal(t, 1, cache.Len())
	//
	time.Sleep(1010 * time.Millisecond)
	require.Equal(t, 0, cache.Len())
}
