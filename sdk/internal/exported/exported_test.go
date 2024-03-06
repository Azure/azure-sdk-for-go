//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasStatusCode(t *testing.T) {
	require.False(t, HasStatusCode(nil, http.StatusAccepted))
	require.False(t, HasStatusCode(&http.Response{}))
	require.False(t, HasStatusCode(&http.Response{StatusCode: http.StatusBadGateway}, http.StatusBadRequest))
	require.True(t, HasStatusCode(&http.Response{StatusCode: http.StatusOK}, http.StatusAccepted, http.StatusOK, http.StatusNoContent))
}

func TestPayload(t *testing.T) {
	const val = "payload"
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(val)),
	}
	b, err := Payload(resp, nil)
	require.NoError(t, err)
	if string(b) != val {
		t.Fatalf("got %s, want %s", string(b), val)
	}
	b, err = Payload(resp, nil)
	require.NoError(t, err)
	if string(b) != val {
		t.Fatalf("got %s, want %s", string(b), val)
	}
}

func TestPayloadDownloaded(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader("payload")),
	}
	require.False(t, PayloadDownloaded(resp))
	_, err := Payload(resp, nil)
	require.NoError(t, err)
	require.True(t, PayloadDownloaded((resp)))
}

func TestPayloadBytesModifier(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader("oldpayload")),
	}
	const newPayload = "newpayload"
	b, err := Payload(resp, &PayloadOptions{
		BytesModifier: func(b []byte) []byte { return []byte(newPayload) },
	})
	require.NoError(t, err)
	require.EqualValues(t, newPayload, string(b))
}

func TestPayloadNilBody(t *testing.T) {
	b, err := Payload(&http.Response{}, nil)
	require.NoError(t, err)
	require.Nil(t, b)
}

func TestNopClosingBytesReader(t *testing.T) {
	const val1 = "the data"
	ncbr := &nopClosingBytesReader{s: []byte(val1)}
	require.NotNil(t, ncbr.Bytes())
	b, err := io.ReadAll(ncbr)
	require.NoError(t, err)
	require.EqualValues(t, val1, b)
	const val2 = "something else"
	ncbr.Set([]byte(val2))
	b, err = io.ReadAll(ncbr)
	require.NoError(t, err)
	require.EqualValues(t, val2, b)
	require.NoError(t, ncbr.Close())
	// seek to beginning and read again
	i, err := ncbr.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Zero(t, i)
	b, err = io.ReadAll(ncbr)
	require.NoError(t, err)
	require.EqualValues(t, val2, b)
	// seek to middle from the end
	i, err = ncbr.Seek(-4, io.SeekEnd)
	require.NoError(t, err)
	require.EqualValues(t, i, len(val2)-4)
	b, err = io.ReadAll(ncbr)
	require.NoError(t, err)
	require.EqualValues(t, "else", b)
	// underflow
	_, err = ncbr.Seek(-int64(len(val2)+1), io.SeekCurrent)
	require.Error(t, err)
}
