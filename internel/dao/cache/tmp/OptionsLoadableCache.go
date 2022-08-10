package tmp

// advanced loadable (eko/gocache) with extra options

import (
	"context"
	"github.com/eko/gocache/v3/cache"
	"sync"

	"github.com/eko/gocache/v3/store"
)

const (
	// LoadableType represents the loadable cache type as a string value
	LoadableType = "loadable"
)

type loadableKeyValue[T any] struct {
	key   any
	value T
}

type LoadFunction[T any] func(ctx context.Context, key any) (T, error)

// OptionsLoadableCache represents a cache that uses a function to load data
type OptionsLoadableCache[T any] struct {
	loadFunc   LoadFunction[T]
	cache      cache.CacheInterface[T]
	setChannel chan *loadableKeyValue[T]
	setterWg   *sync.WaitGroup
}

// NewOptionsLoadableCache instanciates a new cache that uses a function to load data
func NewOptionsLoadableCache[T any](loadFunc LoadFunction[T], cache cache.CacheInterface[T], options ...store.Option) *OptionsLoadableCache[T] {
	loadable := &OptionsLoadableCache[T]{
		loadFunc:   loadFunc,
		cache:      cache,
		setChannel: make(chan *loadableKeyValue[T], 10000),
		setterWg:   &sync.WaitGroup{},
	}

	loadable.setterWg.Add(1)
	go loadable.setter(options...)

	return loadable
}

func (c *OptionsLoadableCache[T]) setter(options ...store.Option) {
	defer c.setterWg.Done()

	for item := range c.setChannel {
		c.Set(context.Background(), item.key, item.value, options...)
	}
}

// Get returns the object stored in cache if it exists
func (c *OptionsLoadableCache[T]) Get(ctx context.Context, key any) (T, error) {
	var err error

	object, err := c.cache.Get(ctx, key)
	if err == nil {
		return object, err
	}

	// Unable to find in cache, try to load it from load function
	object, err = c.loadFunc(ctx, key)
	if err != nil {
		return object, err
	}

	// Then, put it back in cache
	c.setChannel <- &loadableKeyValue[T]{key, object}

	return object, err
}

// Set sets a value in available caches
func (c *OptionsLoadableCache[T]) Set(ctx context.Context, key any, object T, options ...store.Option) error {
	return c.cache.Set(ctx, key, object, options...)
}

// Delete removes a value from cache
func (c *OptionsLoadableCache[T]) Delete(ctx context.Context, key any) error {
	return c.cache.Delete(ctx, key)
}

// Invalidate invalidates cache item from given options
func (c *OptionsLoadableCache[T]) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return c.cache.Invalidate(ctx, options...)
}

// Clear resets all cache data
func (c *OptionsLoadableCache[T]) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

// GetType returns the cache type
func (c *OptionsLoadableCache[T]) GetType() string {
	return LoadableType
}

func (c *OptionsLoadableCache[T]) Close() error {
	close(c.setChannel)
	c.setterWg.Wait()

	return nil
}
