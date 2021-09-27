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

type cacheValue struct {
	value         interface{}
	obsoleteValue interface{}
	complete      bool
	fn            cacheValueTask
	ch            <-chan *cacheTaskResult
	err           error
}

type cacheValueTask func() *cacheTaskResult
type cacheTaskResult struct {
	value interface{}
	err   error
}

func newAsyncCache() *asyncCache {
	return &asyncCache{}
}

func (ac *asyncCache) setValue(key interface{}, value interface{}) {
	ac.values.Store(key, cacheValue{value: value})
}

func (ac *asyncCache) set(key interface{}, singleValueInit cacheValueTask, ctx context.Context) error {
	ch := ac.execCacheValueTask(singleValueInit)
	cachedValue := cacheValue{complete: false, fn: singleValueInit, ch: ch}
	ac.values.Store(key, cachedValue)
	_, err := ac.awaitCacheValue(key, ctx)

	if err != nil {
		return err
	}

	return nil
}

func (ac *asyncCache) getValue(key interface{}) (interface{}, bool) {
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

func (ac *asyncCache) getAsync(key interface{}, obsoleteValue interface{}, singleValueInit cacheValueTask) error {
	var cachedValue cacheValue
	value, valueExists := ac.values.Load(key)

	if !valueExists {
		return nil
	}

	cachedValue, converted := value.(cacheValue)

	if !converted {
		return invalidCacheValue{}
	}

	if cachedValue.complete {
		ch := ac.execCacheValueTask(singleValueInit)
		cachedValue.obsoleteValue = obsoleteValue
		cachedValue.complete = false
		cachedValue.fn = singleValueInit
		cachedValue.ch = ch
		ac.values.Store(key, cachedValue)
	} else {
		cachedValue.fn = singleValueInit
		cachedValue.obsoleteValue = obsoleteValue
		ac.values.Store(key, cachedValue)
	}

	return nil
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

func (ac *asyncCache) execCacheValueTask(t cacheValueTask) <-chan *cacheTaskResult {
	ch := make(chan *cacheTaskResult)

	go func() {
		defer close(ch)
		ch <- t()
	}()
	return ch
}

func (ac *asyncCache) awaitCacheValue(key interface{}, ctx context.Context) (interface{}, error) {
	value, exists := ac.values.Load(key)

	if exists {
		cachedValue, converted := value.(cacheValue)

		if !converted {
			return nil, invalidCacheValue{}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-cachedValue.ch:
			if result == nil {
				return cachedValue.value, cachedValue.err
			}

			if !reflect.DeepEqual(cachedValue.obsoleteValue, result.value) {
				cachedValue.value = result.value
				cachedValue.err = result.err
				cachedValue.complete = true
				ac.values.Store(key, cachedValue)
			} else {
				newch := ac.execCacheValueTask(cachedValue.fn)
				cachedValue.ch = newch
				ac.values.Store(key, cachedValue)

				return ac.awaitCacheValue(key, ctx)
			}
		}

		return cachedValue.value, cachedValue.err
	}

	return nil, nil
}
