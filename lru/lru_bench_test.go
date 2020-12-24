package lru_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/x1nchen/ttlcache/lru"
)

func BenchmarkLRU_Set(b *testing.B) {
	cache := lru.New(lru.WithMaxCache(1000))  // default gcPeriod: 5 * time.Second
	ctx := context.Background()

	// 初始状态
	for i := 0; i < 500; i++ {
		key := "key_base_" + strconv.Itoa(i)
		val := "val_base_" + strconv.Itoa(i)

		cache.Set(ctx, key, val, 500*time.Millisecond)
	}

	// 重置计时器
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		key := "key_" + strconv.Itoa(n)
		val := "val_" + strconv.Itoa(n)

		cache.Set(ctx, key, val, 500*time.Millisecond)
	}
}

func BenchmarkLRU_Get(b *testing.B) {
	cache := lru.New(lru.WithMaxCache(1000))  // default gcPeriod: 5 * time.Second
	ctx := context.Background()

	// 初始状态
	for i := 0; i < 500; i++ {
		key := "key_" + strconv.Itoa(i)
		val := "val_" + strconv.Itoa(i)

		cache.Set(ctx, key, val, 500*time.Second)
	}

	// 重置计时器
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		key := "key_" + strconv.Itoa(n%500)

		cache.Get(ctx, key)
	}
}

func BenchmarkLRU_Del(b *testing.B) {
	cache := lru.New(lru.WithMaxCache(1000))  // default gcPeriod: 5 * time.Second
	ctx := context.Background()

	// 初始状态
	for i := 0; i < 500; i++ {
		key := "key_base_" + strconv.Itoa(i)
		val := "val_base_" + strconv.Itoa(i)

		cache.Set(ctx, key, val, 500*time.Millisecond)
	}

	// 重置计时器
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		key := "key_" + strconv.Itoa(n)
		val := "val_" + strconv.Itoa(n)

		cache.Set(ctx, key, val, 500*time.Millisecond)
		cache.Del(ctx, key)
	}
}
