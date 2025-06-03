// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package amqpwrap

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnsureContextHasTimeout(t *testing.T) {
	t.Run("context without deadline gets default timeout", func(t *testing.T) {
		ctx := context.Background()
		
		// Verify the context doesn't have a deadline initially
		_, hasDeadline := ctx.Deadline()
		require.False(t, hasDeadline)
		
		// Apply the timeout
		ctxWithTimeout, cancel := EnsureContextHasTimeout(ctx, DefaultManagementTimeout)
		defer cancel()
		
		// Verify the context now has a deadline
		deadline, hasDeadline := ctxWithTimeout.Deadline()
		require.True(t, hasDeadline)
		require.True(t, time.Until(deadline) > 0)
		require.True(t, time.Until(deadline) <= DefaultManagementTimeout)
	})

	t.Run("context with existing deadline is preserved", func(t *testing.T) {
		// Create a context with an existing timeout
		existingTimeout := 30 * time.Second
		ctx, originalCancel := context.WithTimeout(context.Background(), existingTimeout)
		defer originalCancel()
		
		originalDeadline, hasDeadline := ctx.Deadline()
		require.True(t, hasDeadline)
		
		// Apply the timeout (should be a no-op)
		ctxWithTimeout, cancel := EnsureContextHasTimeout(ctx, DefaultManagementTimeout)
		defer cancel()
		
		// Verify the deadline is unchanged
		deadline, hasDeadline := ctxWithTimeout.Deadline()
		require.True(t, hasDeadline)
		require.Equal(t, originalDeadline, deadline)
		
		// Verify it's the same context
		require.Equal(t, ctx, ctxWithTimeout)
	})

	t.Run("canceled context with deadline is preserved", func(t *testing.T) {
		// Create a canceled context with deadline
		ctx, originalCancel := context.WithTimeout(context.Background(), time.Hour)
		originalCancel() // cancel immediately
		
		originalDeadline, hasDeadline := ctx.Deadline()
		require.True(t, hasDeadline)
		
		// Apply the timeout (should be a no-op)
		ctxWithTimeout, cancel := EnsureContextHasTimeout(ctx, DefaultManagementTimeout)
		defer cancel()
		
		// Verify the deadline is unchanged
		deadline, hasDeadline := ctxWithTimeout.Deadline()
		require.True(t, hasDeadline)
		require.Equal(t, originalDeadline, deadline)
		
		// Verify it's the same context
		require.Equal(t, ctx, ctxWithTimeout)
	})
}