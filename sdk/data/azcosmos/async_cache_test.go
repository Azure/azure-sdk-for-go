// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_set(t *testing.T) {
	key := "someKey"
	expectedValue := ContainerProperties{ID: "someId"}

	cache := newAsyncCache()

	cache.setValue(key, expectedValue)
	value, _ := cache.getValue(key)
	containerProps, _ := value.(ContainerProperties)
	assert.Equal(t, expectedValue.ID, containerProps.ID)
}

func Test_setAsync(t *testing.T) {
	key := "someKeyAsync"
	expectedValue := ContainerProperties{ID: "someIdAsync"}

	cache := newAsyncCache()

	f := func() *cacheTaskResult {

		return &cacheTaskResult{value: expectedValue, err: nil}
	}

	_ = cache.set(key, f, context.Background())
	value, _ := cache.getValue(key)
	containerProps, _ := value.(ContainerProperties)
	assert.Equal(t, expectedValue.ID, containerProps.ID)
}

func Test_getAsync_not_obsolete(t *testing.T) {
	key := "testAsyncKey"
	expectedValue0 := ContainerProperties{ID: "0"}
	expectedValue1 := ContainerProperties{ID: "1"}
	f1Called := false
	f2Called := false

	ctx := context.Background()

	cache := newAsyncCache()

	f0 := func() *cacheTaskResult {
		return &cacheTaskResult{value: expectedValue0, err: nil}
	}

	_ = cache.set(key, f0, ctx)

	f1 := func() *cacheTaskResult {
		f1Called = true
		time.Sleep(3 * time.Second)
		return &cacheTaskResult{value: expectedValue1, err: nil}
	}

	_ = cache.getAsync(key, expectedValue0, f1)

	f2 := func() *cacheTaskResult {
		f2Called = true
		return &cacheTaskResult{value: expectedValue1, err: nil}
	}

	_ = cache.getAsync(key, expectedValue0, f2)

	value2, _ := cache.awaitCacheValue(key, ctx)
	value, _ := cache.awaitCacheValue(key, ctx)

	assert.True(t, f1Called)
	assert.False(t, f2Called)

	containerProps, _ := value.(ContainerProperties)
	assert.Equal(t, expectedValue1.ID, containerProps.ID)

	containerProps2, _ := value2.(ContainerProperties)
	assert.Equal(t, expectedValue1.ID, containerProps2.ID)
}

func Test_getAsync_obsolete(t *testing.T) {
	key := "testAsyncObsoleteKey"
	expectedValue0 := ContainerProperties{ID: "0"}
	expectedValue1 := ContainerProperties{ID: "1"}
	expectedValue2 := ContainerProperties{ID: "2"}
	f1Called := false
	f2Called := false

	ctx := context.Background()

	cache := newAsyncCache()

	f0 := func() *cacheTaskResult {
		return &cacheTaskResult{value: expectedValue0, err: nil}
	}

	_ = cache.set(key, f0, ctx)

	f1 := func() *cacheTaskResult {
		f1Called = true
		time.Sleep(3 * time.Second)
		return &cacheTaskResult{value: expectedValue1, err: nil}
	}

	_ = cache.getAsync(key, expectedValue0, f1)

	f2 := func() *cacheTaskResult {
		f2Called = true
		return &cacheTaskResult{value: expectedValue2, err: nil}
	}

	_ = cache.getAsync(key, expectedValue1, f2)

	value, _ := cache.awaitCacheValue(key, ctx)
	containerProps, _ := value.(ContainerProperties)

	value2, _ := cache.awaitCacheValue(key, ctx)
	containerProps2, _ := value2.(ContainerProperties)

	assert.True(t, f1Called)
	assert.True(t, f2Called)
	assert.Equal(t, expectedValue2.ID, containerProps.ID)
	assert.Equal(t, expectedValue2.ID, containerProps2.ID)
}

func Test_getAsync_obsolete_with_error(t *testing.T) {
	key := "testAsyncObsoleteKey"
	expectedValue0 := ContainerProperties{ID: "0"}
	expectedValue1 := ContainerProperties{ID: "1"}
	expectedValue2 := ContainerProperties{ID: "2"}
	f1Called := false
	f2Called := false

	ctx := context.Background()

	cache := newAsyncCache()

	f0 := func() *cacheTaskResult {
		return &cacheTaskResult{value: expectedValue0, err: nil}
	}

	_ = cache.set(key, f0, ctx)

	f1 := func() *cacheTaskResult {
		f1Called = true
		time.Sleep(3 * time.Second)
		return &cacheTaskResult{value: nil, err: errors.New("some error")}
	}

	_ = cache.getAsync(key, expectedValue0, f1)

	f2 := func() *cacheTaskResult {
		f2Called = true
		return &cacheTaskResult{value: expectedValue2, err: nil}
	}

	_ = cache.getAsync(key, expectedValue1, f2)

	_, err := cache.awaitCacheValue(key, ctx)

	_, err2 := cache.awaitCacheValue(key, ctx)

	assert.True(t, f1Called)
	assert.False(t, f2Called)
	assert.Error(t, err)
	assert.Error(t, err2)
}

func Test_getAsync_obsolete_with_context_error(t *testing.T) {
	key := "testAsyncObsoleteKey"
	expectedValue0 := ContainerProperties{ID: "0"}
	expectedValue1 := ContainerProperties{ID: "1"}
	expectedValue2 := ContainerProperties{ID: "2"}
	f1Called := false
	f2Called := false

	ctx := context.Background()

	cache := newAsyncCache()

	f0 := func() *cacheTaskResult {
		return &cacheTaskResult{value: expectedValue0, err: nil}
	}

	_ = cache.set(key, f0, ctx)

	f1 := func() *cacheTaskResult {
		f1Called = true
		time.Sleep(3 * time.Second)
		return &cacheTaskResult{value: nil, err: errors.New("some error")}
	}

	_ = cache.getAsync(key, expectedValue0, f1)

	f2 := func() *cacheTaskResult {
		f2Called = true
		return &cacheTaskResult{value: expectedValue2, err: nil}
	}

	ctx.Done()

	_ = cache.getAsync(key, expectedValue1, f2)

	_, err := cache.awaitCacheValue(key, ctx)

	_, err2 := cache.awaitCacheValue(key, ctx)

	assert.True(t, f1Called)
	assert.False(t, f2Called)
	assert.Error(t, err)
	assert.Error(t, err2)
}

func Test_remove(t *testing.T) {
	key := "someKeyToRemove"
	expectedValue := ContainerProperties{ID: "someIdToRemove"}

	cache := newAsyncCache()

	cache.setValue(key, expectedValue)
	value, _ := cache.getValue(key)
	containerProps, _ := value.(ContainerProperties)
	assert.Equal(t, expectedValue.ID, containerProps.ID)

	cache.remove(key)

	_, ok := cache.getValue(key)

	assert.False(t, ok)
}

func Test_clear(t *testing.T) {
	key := "someKeyToClear"
	expectedValue := ContainerProperties{ID: "someIdToDelete"}
	key2 := "someKeyToClear2"
	expectedValue2 := ContainerProperties{ID: "someIdToDelete2"}

	cache := newAsyncCache()

	cache.setValue(key, expectedValue)
	value, _ := cache.getValue(key)
	containerProps, _ := value.(ContainerProperties)
	assert.Equal(t, expectedValue.ID, containerProps.ID)

	cache.setValue(key2, expectedValue2)
	value2, _ := cache.getValue(key2)
	containerProps2, _ := value2.(ContainerProperties)
	assert.Equal(t, expectedValue2.ID, containerProps2.ID)

	cache.clear()

	_, ok := cache.getValue(key)

	assert.False(t, ok)

	_, ok2 := cache.getValue(key2)

	assert.False(t, ok2)
}
