//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package temporal

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewExpiringResource(t *testing.T) {
	er := NewResource(func(state string) (newResource string, newExpiration time.Time, err error) {
		switch state {
		case "initial":
			return "updated", time.Now().Add(-time.Minute), nil
		case "updated":
			return "refreshed", time.Now().Add(1 * time.Hour), nil
		default:
			t.Fatalf("unexpected state %s", state)
			return "", time.Time{}, errors.New("unexpected")
		}
	})
	res, err := er.Get("initial")
	require.NoError(t, err)
	require.Equal(t, "updated", res)
	res, err = er.Get(res)
	require.NoError(t, err)
	require.Equal(t, "refreshed", res)
	res, err = er.Get(res)
	require.NoError(t, err)
	require.Equal(t, "refreshed", res)
}

func TestExpiringResourceError(t *testing.T) {
	expectedState := "expected state"
	expectedError := "expected error"
	calls := 0
	er := NewResource(func(state string) (newResource string, newExpiration time.Time, err error) {
		calls += 1
		if calls == 1 {
			return expectedState, time.Now().Add(time.Minute), nil
		} else {
			return "un" + expectedState, time.Time{}, errors.New(expectedError)
		}
	})
	res, err := er.Get(expectedState)
	require.NoError(t, err)
	require.Equal(t, expectedState, res)

	// When an eager update fails, GetResource should return the prior value and no error.
	er.lastAttempt = time.Now().Add(-time.Hour)
	for i := 0; i < 3; i++ {
		res, err = er.Get(res)
		require.NoError(t, err)
		require.Equal(t, expectedState, res)
		// GetResource should wait before trying a second eager update i.e. it shouldn't make a third call in this loop
		require.Equal(t, 2, calls)
	}

	// After the resource has expired, GetResource should return any error from updating
	er.expiration = time.Now().Add(-time.Hour)
	_, err = er.Get(res)
	require.Error(t, err, expectedError)
	require.Equal(t, 3, calls)
}

func TestNewExpiringResourceConcurrent(t *testing.T) {
	er := NewResource(func(i int) (newResource int, newExpiration time.Time, err error) {
		return i + 1, time.Now().Add(time.Millisecond), nil
	})

	wg := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			_, _ = er.Get(0)
			wg.Done()
		}()
	}

	wg.Wait()
}
