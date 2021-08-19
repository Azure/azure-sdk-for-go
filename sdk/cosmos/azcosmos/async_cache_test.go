// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCacheKey struct {
	name string
}

func Test_set(t *testing.T) {
	key := "someKey"
	expectedValue := CosmosContainerProperties{Id: "someId"}

	cache := newAsyncCache()

	cache.set(key, expectedValue)
	value, _ := cache.get(key)
	containerProps, _ := value.(CosmosContainerProperties)
	assert.Equal(t, expectedValue.Id, containerProps.Id)
}

func Test_setAsync(t *testing.T) {
	key := "someKeyAsync"
	expectedValue := CosmosContainerProperties{Id: "someIdAsync"}

	cache := newAsyncCache()

	f := func() interface{} {
		return expectedValue
	}
	cache.setAsync(key, f, context.Background())
	value, _ := cache.get(key)
	containerProps, _ := value.(CosmosContainerProperties)
	assert.Equal(t, expectedValue.Id, containerProps.Id)
}

func Test_remove(t *testing.T) {
	key := "someKeyToRemove"
	expectedValue := CosmosContainerProperties{Id: "someIdToRemove"}

	cache := newAsyncCache()

	cache.set(key, expectedValue)
	value, _ := cache.get(key)
	containerProps, _ := value.(CosmosContainerProperties)
	assert.Equal(t, expectedValue.Id, containerProps.Id)

	cache.remove(key)

	_, ok := cache.get(key)

	assert.False(t, ok)
}

func Test_clear(t *testing.T) {
	key := "someKeyToClear"
	expectedValue := CosmosContainerProperties{Id: "someIdToDelete"}
	key2 := "someKeyToClear2"
	expectedValue2 := CosmosContainerProperties{Id: "someIdToDelete2"}

	cache := newAsyncCache()

	cache.set(key, expectedValue)
	value, _ := cache.get(key)
	containerProps, _ := value.(CosmosContainerProperties)
	assert.Equal(t, expectedValue.Id, containerProps.Id)

	cache.set(key2, expectedValue2)
	value2, _ := cache.get(key2)
	containerProps2, _ := value2.(CosmosContainerProperties)
	assert.Equal(t, expectedValue2.Id, containerProps2.Id)

	cache.clear()

	_, ok := cache.get(key)

	assert.False(t, ok)

	_, ok2 := cache.get(key2)

	assert.False(t, ok2)
}
