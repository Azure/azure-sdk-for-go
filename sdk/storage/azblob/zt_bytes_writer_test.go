// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint
func TestBytesWriterWriteAt(t *testing.T) {
	b := make([]byte, 10)
	buffer := newBytesWriter(b)

	count, err := buffer.WriteAt([]byte{1, 2}, 10)
	require.Contains(t, err.Error(), "offset value is out of range")
	require.Equal(t, count, 0)

	count, err = buffer.WriteAt([]byte{1, 2}, -1)
	require.Contains(t, err.Error(), "offset value is out of range")
	require.Equal(t, count, 0)

	count, err = buffer.WriteAt([]byte{1, 2}, 9)
	require.Contains(t, err.Error(), "not enough space for all bytes")
	require.Equal(t, count, 1)
	require.Equal(t, bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}), 0)

	count, err = buffer.WriteAt([]byte{1, 2}, 8)
	require.NoError(t, err)
	require.Equal(t, count, 2)
	require.Equal(t, bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 2}), 0)
}
