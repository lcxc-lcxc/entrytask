package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

/**
缓存组件，如果缓存miss，可以自动加载数据库中的内容。
*/

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

/**
该函数从管道里面获取内容，set到redis里面。相当于异步set cache
*/
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

// MGet 该函数需要的loadFunction应该是获取单个key-value的loadFunction
func (c *LoadableCache[T]) MGet(ctx context.Context, keys ...string) ([]T, error) {
	startTime := time.Now()
	values, err := c.cache.MGet(ctx, keys...).Result()
	dur := time.Since(startTime)
	fmt.Println("MGET: " + dur.String())
	if err != nil {
		return nil, err
	}
	var returnValues []T
	for idx, value := range values {
		if value == nil {
			object, err := c.loadFunc(ctx, keys[idx])
			if err != nil {
				return nil, err
			}
			c.setChannel <- &loadableKeyValue[T]{keys[idx], object}
			returnValues = append(returnValues, object)
		} else {
			t := new(T)
			valueBytes := []byte(value.(string))
			err := json.Unmarshal(valueBytes, t)
			if err != nil {
				return nil, err
			}
			returnValues = append(returnValues, *t)

		}
	}
	return returnValues, nil
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

// Close close the channel and let the go routine return
func (c *LoadableCache[T]) Close() error {
	close(c.setChannel)
	c.setterWg.Wait()

	return nil
}
