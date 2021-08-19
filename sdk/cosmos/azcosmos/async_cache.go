// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
)

type invalidCacheValue struct{}

func (i invalidCacheValue) Error() string { return "Invalid cache value" }

type asyncCache struct {
	values sync.Map
}

type cacheValueTask func() interface{}

type cacheValue struct {
	value    interface{}
	complete bool
	ch       <-chan interface{}
}

func newAsyncCache() *asyncCache {
	return &asyncCache{}
}

func (ac *asyncCache) set(key interface{}, value interface{}) {
	ac.values.Store(key, cacheValue{value: value})
}

func (ac *asyncCache) setAsync(key interface{}, singleValueInit cacheValueTask, ctx context.Context) error {
	ch := ac.execCacheValueTask(singleValueInit)
	cachedValue := cacheValue{value: nil, complete: false, ch: ch}
	ac.values.Store(key, cachedValue)
	_, err := ac.awaitCacheValue(key, ch, ctx)

	if err != nil {
		return err
	}

	return nil
}

func (ac *asyncCache) get(key interface{}) (interface{}, bool) {
	var cachedValue cacheValue
	value, ok := ac.values.Load(key)

	if !ok {
		return nil, false
	}

	cachedValue, ok = value.(cacheValue)

	if ok {
		return cachedValue.value, ok
	}

	return nil, false
}

func (ac *asyncCache) getAsync(key interface{}, singleValueInit cacheValueTask) (*cacheValue, error) {
	var cachedValue cacheValue
	value, valueExists := ac.values.Load(key)

	if !valueExists {
		return nil, nil
	}

	cachedValue, converted := value.(cacheValue)

	if !converted {
		return nil, invalidCacheValue{}
	}

	if cachedValue.complete {
		ch := ac.execCacheValueTask(singleValueInit)
		cachedValue.complete = false
		cachedValue.ch = ch
		ac.values.Store(key, cachedValue)
	}

	return &cachedValue, nil
}

func (ac *asyncCache) remove(key interface{}) {
	ac.values.Delete(key)
}

func (ac *asyncCache) clear() {
	ac.values.Range(func(key interface{}, value interface{}) bool {
		ac.values.Delete(key)
		return true
	})

}

func (ac *asyncCache) execCacheValueTask(t cacheValueTask) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)
		ch <- t()
	}()
	return ch
}

func (ac *asyncCache) awaitCacheValue(key interface{}, ch <-chan interface{}, ctx context.Context) (interface{}, error) {
	value, exists := ac.values.Load(key)
	if exists {
		cachedValue, converted := value.(cacheValue)

		if !converted {
			return nil, invalidCacheValue{}
		}

		if !cachedValue.complete {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case result := <-ch:
				cachedValue.complete = true
				cachedValue.value = result
				cachedValue.ch = ch
				ac.values.Store(key, cachedValue)
			}
		}

		return cachedValue.value, nil
	}

	return nil, nil
}
