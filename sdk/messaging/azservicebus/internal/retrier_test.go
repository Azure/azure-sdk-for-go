// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetrier(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		retrier := NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			MaxRetries: 5,
			Factor:     0,
		})

		require := require.New(t)

		// first iteration is always free (ie, that's not the
		// retry part)
		require.True(retrier.Try(context.Background()))

		// now we're doing retries
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))

		// and it's the 6th retry that fails since we've exhausted
		// the retries we're allotted.
		require.False(retrier.Try(context.Background()))
		require.True(retrier.Exhausted())
	})

	t.Run("Cancellation", func(t *testing.T) {
		retrier := NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			MaxRetries: 5,
			Factor:     0,
		})

		// first iteration is always free (ie, that's not the
		// retry part)
		cancelledContext, cancel := context.WithCancel(context.Background())
		cancel()
		require.False(t, retrier.Try(cancelledContext))
	})
}
