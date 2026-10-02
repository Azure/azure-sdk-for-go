// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResultAfterCancellationReturnsTerminalCompletion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result := completionResult{body: []byte("created")}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.NoError(t, err)
		require.Equal(t, result.body, got.body)
	})

	t.Run("failure", func(t *testing.T) {
		serviceErr := &Error{Code: CodeConflict}
		result := completionResult{err: serviceErr}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.NoError(t, err)
		require.Same(t, serviceErr, got.err)
	})

	t.Run("cancelled", func(t *testing.T) {
		diagnostics := &Diagnostics{AttemptCount: 2, StatusCode: 429, SubStatus: 3200}
		result := completionResult{
			cancelled: true,
			err: &Error{
				Code:          CodeOperationCancelled,
				RequestCharge: 1.5,
				ActivityID:    "activity-id",
				Diagnostics:   diagnostics,
				AttemptCount:  2,
				StatusCode:    429,
				SubStatus:     3200,
			},
		}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.Empty(t, got)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NotErrorIs(t, err, context.Canceled)
		var cosmosErr *Error
		require.True(t, errors.As(err, &cosmosErr))
		require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
		require.Equal(t, 1.5, cosmosErr.RequestCharge)
		require.Equal(t, "activity-id", cosmosErr.ActivityID)
		require.Same(t, diagnostics, cosmosErr.Diagnostics)
		require.Equal(t, uint32(2), cosmosErr.AttemptCount)
		require.Equal(t, 429, cosmosErr.StatusCode)
		require.Equal(t, 3200, cosmosErr.SubStatus)
	})
}
