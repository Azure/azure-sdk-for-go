//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

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
			return "updated", time.Now(), nil
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

func TestNewExpiringResourceError(t *testing.T) {
	er := NewExpiringResource(func(state interface{}) (newResource interface{}, newExpiration time.Time, err error) {
		return "", time.Time{}, errors.New("failed")
	})
	res, err := er.GetResource("stale")
	require.Error(t, err)
	require.Equal(t, "", res)
}
