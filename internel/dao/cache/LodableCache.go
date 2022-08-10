package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type loadableKeyValue[T any] struct {
	key   string
	value T
}

type LoadFunction[T any] func(ctx context.Context, key any) (T, error)

// LoadableCache represents a cache that uses a function to load data
type LoadableCache[T any] struct {
	loadFunc   LoadFunction[T]
	cache      *redis.Client
	setChannel chan *loadableKeyValue[T]
	setterWg   *sync.WaitGroup
}

// NewLoadableCache instanciates a new cache that uses a function to load data
func NewLoadableCache[T any](loadFunc LoadFunction[T], cache *redis.Client, timeout time.Duration) *LoadableCache[T] {
	loadable := &LoadableCache[T]{
		loadFunc:   loadFunc,
		cache:      cache,
		setChannel: make(chan *loadableKeyValue[T], 10000),
		setterWg:   &sync.WaitGroup{},
	}

	loadable.setterWg.Add(1)
	go loadable.setter(timeout)

	return loadable
}

func (c *LoadableCache[T]) setter(timeout time.Duration) {
	defer c.setterWg.Done()

	for item := range c.setChannel {
		c.Set(context.Background(), item.key, item.value, timeout)
	}
}

// Get returns the object stored in cache if it exists
func (c *LoadableCache[T]) Get(ctx context.Context, key string) (T, error) {
	var err error
	t := new(T)

	// get data from cache
	objectStr, err := c.cache.Get(ctx, key).Result()
	if err == nil {
		err2 := json.Unmarshal([]byte(objectStr), t)
		if err2 == nil {
			return *t, nil
		} else {
			return *new(T), err2
		}
	}

	// Unable to find in cache, try to load it from load function
	object, err := c.loadFunc(ctx, key)
	if err != nil {
		return object, err
	}

	// Then, put it back in cache
	c.setChannel <- &loadableKeyValue[T]{key, object}

	return object, err
}

// Set sets a value in available caches
func (c *LoadableCache[T]) Set(ctx context.Context, key string, object T, timeout time.Duration) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return c.cache.Set(ctx, key, objectBytes, timeout).Err()
}

// Delete removes a value from cache
func (c *LoadableCache[T]) Delete(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

func (c *LoadableCache[T]) Close() error {
	close(c.setChannel)
	c.setterWg.Wait()

	return nil
}
