// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"reflect"
	"sync"
)

type invalidCacheValue struct{}

func (i invalidCacheValue) Error() string { return "Invalid cache value" }

type asyncCache struct {
	values sync.Map
}

type cacheValueTask func() interface{}

type cacheValue struct {
	value interface{}
	ch    <-chan interface{}
}

func newAsyncCache() *asyncCache {
	return &asyncCache{}
}

func (ac *asyncCache) set(key interface{}, value interface{}) {
	ac.values.Store(key, cacheValue{value: value})
}

func (ac *asyncCache) setAsync(key interface{}, singleValueInit cacheValueTask, ctx context.Context) error {
	value, ok := ac.values.Load(key)

	if !ok {
		ch := ac.execCacheValueTask(singleValueInit)
		cachedValue, err := ac.awaitCacheValue(ch, ctx)

		if err != nil {
			return err
		}

		ac.values.Store(key, cacheValue{value: cachedValue, ch: ch})

		return nil
	}

	cachedValue, converted := value.(cacheValue)

	if !converted {
		return invalidCacheValue{}
	}

	ac.awaitCacheValue(cachedValue.ch, ctx)

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

func (ac *asyncCache) getAsync(key interface{}, obsoleteValue interface{}, singleValueInit cacheValueTask, ctx context.Context) (interface{}, error) {
	var cachedValue cacheValue
	value, ok := ac.values.Load(key)

	if !ok {
		return nil, nil
	}

	cachedValue, converted := value.(cacheValue)

	if !converted {
		return nil, invalidCacheValue{}
	}

	awaitedValue, err := ac.awaitCacheValue(cachedValue.ch, ctx)

	if err != nil {
		return nil, err
	}

	if awaitedValue != nil && !reflect.DeepEqual(awaitedValue, obsoleteValue) {
		return awaitedValue, nil
	}

	ch := ac.execCacheValueTask(singleValueInit)
	ac.values.Store(key, cacheValue{ch: ch})
	awaitedValue, err = ac.awaitCacheValue(ch, ctx)

	if err != nil {
		return nil, err
	}

	return awaitedValue, nil
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

func (ac *asyncCache) awaitCacheValue(ch <-chan interface{}, ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-ch:
		return result, nil
	}
}
