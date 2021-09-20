// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"

	"github.com/jpillora/backoff"
	"github.com/stretchr/testify/require"
)

func TestRetrier(t *testing.T) {
	retrier := BackoffRetrier{
		MaxRetries: 5,
		Backoff: backoff.Backoff{
			Factor: 0,
		},
	}

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
}
