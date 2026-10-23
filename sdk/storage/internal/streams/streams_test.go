// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package streams

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSeekableStreamAt0AndGetCount_Nil(t *testing.T) {
	count, err := ValidateSeekableStreamAt0AndGetCount(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestValidateSeekableStreamAt0AndGetCount_Empty(t *testing.T) {
	body := bytes.NewReader(nil)
	count, err := ValidateSeekableStreamAt0AndGetCount(body)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestValidateSeekableStreamAt0AndGetCount_WithContent(t *testing.T) {
	data := []byte("hello world")
	body := bytes.NewReader(data)
	count, err := ValidateSeekableStreamAt0AndGetCount(body)
	require.NoError(t, err)
	require.Equal(t, int64(len(data)), count)

	// verify stream is reset to position 0
	pos, err := body.Seek(0, 1) // SeekCurrent
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
}

func TestValidateSeekableStreamAt0AndGetCount_NotAtZero(t *testing.T) {
	body := bytes.NewReader([]byte("hello"))
	_, _ = body.Seek(3, 0) // move to position 3

	_, err := ValidateSeekableStreamAt0AndGetCount(body)
	require.Error(t, err)
	require.Contains(t, err.Error(), "position 0")
}

func TestValidateSeekableStreamAt0AndGetCount_NonSeekable(t *testing.T) {
	body := &nonSeekableReader{}
	_, err := ValidateSeekableStreamAt0AndGetCount(body)
	require.Error(t, err)
	require.Contains(t, err.Error(), "seekable")
}

func TestValidateSeekableStreamAt0_Nil(t *testing.T) {
	err := validateSeekableStreamAt0(nil)
	require.NoError(t, err)
}

func TestValidateSeekableStreamAt0_AtZero(t *testing.T) {
	body := strings.NewReader("test")
	err := validateSeekableStreamAt0(body)
	require.NoError(t, err)
}

// nonSeekableReader implements io.ReadSeeker but always fails on Seek.
type nonSeekableReader struct{}

func (r *nonSeekableReader) Read(p []byte) (int, error) {
	return 0, nil
}

func (r *nonSeekableReader) Seek(offset int64, whence int) (int64, error) {
	return 0, &seekError{}
}

type seekError struct{}

func (e *seekError) Error() string { return "seek not supported" }
