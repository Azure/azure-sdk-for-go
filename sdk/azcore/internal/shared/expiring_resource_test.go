//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewExpiringResource(t *testing.T) {
	er := NewExpiringResource(func(state interface{}) (newResource interface{}, newExpiration time.Time, err error) {
		s := state.(string)
		switch s {
		case "initial":
			return "updated", time.Now().Add(-time.Minute), nil
		case "updated":
			return "refreshed", time.Now().Add(1 * time.Hour), nil
		default:
			t.Fatalf("unexpected state %s", s)
			return "", time.Time{}, errors.New("unexpected")
		}
	})
	res, err := er.GetResource("initial")
	require.NoError(t, err)
	require.Equal(t, "updated", res)
	res, err = er.GetResource(res)
	require.NoError(t, err)
	require.Equal(t, "refreshed", res)
	res, err = er.GetResource(res)
	require.NoError(t, err)
	require.Equal(t, "refreshed", res)
}

func TestExpiringResourceError(t *testing.T) {
	expectedState := "expected state"
	expectedError := "expected error"
	calls := 0
	er := NewExpiringResource(func(state interface{}) (newResource interface{}, newExpiration time.Time, err error) {
		calls += 1
		if calls == 1 {
			return expectedState, time.Now().Add(time.Minute), nil
		} else {
			return "un" + expectedState, time.Time{}, errors.New(expectedError)
		}
	})
	res, err := er.GetResource(expectedState)
	require.NoError(t, err)
	require.Equal(t, expectedState, res)

	// When an eager update fails, GetResource should return the prior value and no error.
	er.lastAttempt = time.Now().Add(-time.Hour)
	for i := 0; i < 3; i++ {
		res, err = er.GetResource(res)
		require.NoError(t, err)
		require.Equal(t, expectedState, res)
		// GetResource should wait before trying a second eager update i.e. it shouldn't make a third call in this loop
		require.Equal(t, 2, calls)
	}

	// After the resource has expired, GetResource should return any error from updating
	er.expiration = time.Now().Add(-time.Hour)
	_, err = er.GetResource(res)
	require.Error(t, err, expectedError)
	require.Equal(t, 3, calls)
}
