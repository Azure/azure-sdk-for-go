//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSectionWriter(t *testing.T) {
	b := [10]byte{}
	buffer := NewBytesWriter(b[:])

	section := NewSectionWriter(buffer, 0, 5)
	require.Equal(t, section.Count, int64(5))
	require.Equal(t, section.Offset, int64(0))
	require.Equal(t, section.Position, int64(0))

	count, err := section.Write([]byte{1, 2, 3})
	require.NoError(t, err)
	require.Equal(t, count, 3)
	require.Equal(t, section.Position, int64(3))
	require.Equal(t, b, [10]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{4, 5, 6})
	require.Contains(t, err.Error(), "not enough space for all bytes")
	require.Equal(t, count, 2)
	require.Equal(t, section.Position, int64(5))
	require.Equal(t, b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{6, 7, 8})
	require.Contains(t, err.Error(), "end of section reached")
	require.Equal(t, count, 0)
	require.Equal(t, section.Position, int64(5))
	require.Equal(t, b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	// Intentionally create a section writer which will attempt to write
	// outside the bounds of the buffer.
	section = NewSectionWriter(buffer, 5, 6)
	require.Equal(t, section.Count, int64(6))
	require.Equal(t, section.Offset, int64(5))
	require.Equal(t, section.Position, int64(0))

	count, err = section.Write([]byte{6, 7, 8})
	require.NoError(t, err)
	require.Equal(t, count, 3)
	require.Equal(t, section.Position, int64(3))
	require.Equal(t, b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0})

	// Attempt to write past the end of the section. Since the underlying
	// buffer rejects the write it gives the same error as in the normal case.
	count, err = section.Write([]byte{9, 10, 11})
	require.Contains(t, err.Error(), "not enough space for all bytes")
	require.Equal(t, count, 2)
	require.Equal(t, section.Position, int64(5))
	require.Equal(t, b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Attempt to write past the end of the buffer. In this case the buffer
	// rejects the write completely since it falls completely out of bounds.
	count, err = section.Write([]byte{11, 12, 13})
	require.Contains(t, err.Error(), "offset value is out of range")
	require.Equal(t, count, 0)
	require.Equal(t, section.Position, int64(5))
	require.Equal(t, b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func TestSectionWriterCopySrcDestEmpty(t *testing.T) {
	input := make([]byte, 0)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := NewBytesWriter(output)
	section := NewSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	require.NoError(t, err)
	require.Equal(t, count, int64(0))
}

func TestSectionWriterCopyDestEmpty(t *testing.T) {
	input := make([]byte, 10)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := NewBytesWriter(output)
	section := NewSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	require.Contains(t, err.Error(), "end of section reached")
	require.Equal(t, count, int64(0))
}
