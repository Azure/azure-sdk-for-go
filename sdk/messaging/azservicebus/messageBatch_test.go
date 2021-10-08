// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageBatchUnitTests(t *testing.T) {
	t.Run("addMessages", func(t *testing.T) {
		mb := newMessageBatch(1000)

		added, err := mb.Add(&Message{
			Body: []byte("hello world"),
		})

		require.NoError(t, err)
		require.True(t, added)
		require.EqualValues(t, 1, mb.Len())
		require.EqualValues(t, 195, mb.Size())
	})

	t.Run("addTooManyMessages", func(t *testing.T) {
		mb := newMessageBatch(1)

		added, err := mb.Add(&Message{
			Body: []byte("hello world"),
		})

		require.False(t, added)
		require.Nil(t, err)
	})
}
