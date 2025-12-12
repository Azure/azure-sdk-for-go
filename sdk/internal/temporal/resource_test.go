// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package temporal

import (
	"errors"
	"fmt"
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

func TestShouldRefresh(t *testing.T) {
	for _, shouldRefresh := range []bool{true, false} {
		t.Run(fmt.Sprint(shouldRefresh), func(t *testing.T) {
			before := backoff
			defer func() { backoff = before }()
			backoff = func(time.Time, time.Time) bool { return false }
			called := false
			initial, updated := "initial", "updated"
			states := []string{"a", "b"}
			r := NewResourceWithOptions(
				func(state string) (string, time.Time, error) {
					rsrc := initial
					if called {
						require.Equal(t, states[1], state)
						rsrc = updated
					} else {
						require.Equal(t, states[0], state)
					}
					return rsrc, time.Now().Add(time.Hour), nil
				},
				ResourceOptions[string, string]{
					ShouldRefresh: func(resource, state string) bool {
						require.False(t, called)
						require.Equal(t, states[1], state)
						require.Equal(t, initial, resource)
						called = true
						return shouldRefresh
					},
				},
			)
			resources := []string{initial, updated}
			for _, state := range states {
				actual, err := r.Get(state)
				require.NoError(t, err)
				expected := resources[0]
				if shouldRefresh {
					resources = resources[1:]
				}
				require.Equal(t, expected, actual)
			}
			require.True(t, called)
		})
	}
	t.Run("backoff", func(t *testing.T) {
		gets := 0
		r := NewResourceWithOptions(
			func(string) (string, time.Time, error) {
				gets++
				return "", time.Now().Add(time.Hour), nil
			},
			ResourceOptions[string, string]{
				ShouldRefresh: func(string, string) bool { return true },
			},
		)
		for i := 0; i < 42; i++ {
			_, err := r.Get("")
			require.NoError(t, err)
		}
		require.Equal(t, 1, gets, "Get should apply backoff when ShouldRefresh returns true")
	})
	t.Run("expired resource", func(t *testing.T) {
		r := NewResourceWithOptions(
			func(string) (string, time.Time, error) {
				return "expired", time.Now().Add(-time.Hour), nil
			},
			ResourceOptions[string, string]{
				ShouldRefresh: func(string, string) bool {
					t.Fatal("Resource shouldn't call ShouldRefresh when it's expired")
					return true
				},
			},
		)
		for i := 0; i < 3; i++ {
			_, err := r.Get("")
			require.NoError(t, err)
		}
	})
}
