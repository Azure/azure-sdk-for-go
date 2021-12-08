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

		err := mb.AddMessage(&Message{
			Body: []byte("hello world"),
		})

		require.NoError(t, err)
		require.EqualValues(t, 1, mb.NumMessages())
		require.EqualValues(t, 183, mb.NumBytes())
	})

	t.Run("addTooManyMessages", func(t *testing.T) {
		mb := newMessageBatch(1)

		err := mb.AddMessage(&Message{
			Body: []byte("hello world"),
		})

		require.EqualError(t, err, ErrMessageTooLarge.Error())
	})
}
