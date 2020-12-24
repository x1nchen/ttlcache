package lru

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/x1nchen/ttlcache"
)

type node struct {
	key      string
	value    interface{}
	expireAt time.Time
}

type options struct {
	maxCache int
	gcPeriod time.Duration
}

var defaultOptions = options{
	maxCache: 0,
	gcPeriod: 5 * time.Second,
}

type Option func(options *options)

func WithMaxCache(maxCache int) Option {
	return func(options *options) {
		options.maxCache = maxCache
	}
}

func WithGcPeriod(gcPeriod time.Duration) Option {
	return func(options *options) {
		options.gcPeriod = gcPeriod
	}
}

type Cache struct {
	mu sync.Mutex
	options
	l *list.List
	m map[string]*list.Element
}

func New(opts ...Option) *Cache {
	//
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	//
	mc := &Cache{
		options: options,
		m:       make(map[string]*list.Element),
		l:       list.New(),
	}
	mc.runGC()
	return mc
}

func (c *Cache) runGC() {
	time.AfterFunc(c.gcPeriod, func() {
		c.runGC()
		c.GC()
	})
}

const gcSample = 20
const gcPercent = 0.25

func (c *Cache) GC() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for {
		count := 1.0
		gcCount := 0.0
		gcStart := time.Now()
		for _, e := range c.m {
			count++
			n := e.Value.(*node)
			if n.expireAt.Before(gcStart) {
				c.l.Remove(e)
				delete(c.m, n.key)
				gcCount++
			}
			if count > gcSample {
				break
			}
		}
		if gcCount/count < gcPercent {
			break
		}
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if ok {
		c.l.Remove(e)
	}
	e = c.l.PushBack(&node{key, value, time.Now().Add(ttl)})
	c.m[key] = e

	if c.maxCache > 0 && c.l.Len() > c.maxCache {
		e := c.l.Front()
		n := c.l.Remove(e).(*node)
		delete(c.m, n.key)
	}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if !ok {
		return nil, cache.ErrKeyNotFound
	}

	n := e.Value.(*node)
	if n.expireAt.Before(time.Now()) {
		c.l.Remove(e)
		delete(c.m, n.key)
		return nil, cache.ErrKeyNotFound
	}
	c.l.MoveToBack(e)
	return n.value, nil
}

func (c *Cache) Del(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if ok {
		c.l.Remove(e)
		delete(c.m, key)
	}
}

func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.l.Len()
}
