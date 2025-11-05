//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNopCloser(t *testing.T) {
	nc := NopCloser(strings.NewReader("foo"))
	if err := nc.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestHasStatusCode(t *testing.T) {
	if HasStatusCode(nil, http.StatusAccepted) {
		t.Fatal("unexpected success")
	}
	if HasStatusCode(&http.Response{}) {
		t.Fatal("unexpected success")
	}
	if HasStatusCode(&http.Response{StatusCode: http.StatusBadGateway}, http.StatusBadRequest) {
		t.Fatal("unexpected success")
	}
	if !HasStatusCode(&http.Response{StatusCode: http.StatusOK}, http.StatusAccepted, http.StatusOK, http.StatusNoContent) {
		t.Fatal("unexpected failure")
	}
}

func TestDecodeByteArray(t *testing.T) {
	out := []byte{}
	require.NoError(t, DecodeByteArray("", &out, Base64StdFormat))
	require.Empty(t, out)
	const (
		stdEncoding = "VGVzdERlY29kZUJ5dGVBcnJheQ=="
		urlEncoding = "VGVzdERlY29kZUJ5dGVBcnJheQ"
		decoded     = "TestDecodeByteArray"
	)
	require.NoError(t, DecodeByteArray(stdEncoding, &out, Base64StdFormat))
	require.EqualValues(t, decoded, string(out))
	require.NoError(t, DecodeByteArray(urlEncoding, &out, Base64URLFormat))
	require.EqualValues(t, decoded, string(out))
	require.NoError(t, DecodeByteArray(fmt.Sprintf("\"%s\"", stdEncoding), &out, Base64StdFormat))
	require.EqualValues(t, decoded, string(out))
	require.Error(t, DecodeByteArray(stdEncoding, &out, 123))
	require.Error(t, DecodeByteArray("\"", &out, Base64StdFormat))
}

func TestNewKeyCredential(t *testing.T) {
	const val1 = "foo"
	cred := NewKeyCredential(val1)
	require.NotNil(t, cred)
	require.EqualValues(t, val1, KeyCredentialGet(cred))
	const val2 = "bar"
	cred.Update(val2)
	require.EqualValues(t, val2, KeyCredentialGet(cred))
}

func TestNewSASCredential(t *testing.T) {
	const val1 = "foo"
	cred := NewSASCredential(val1)
	require.NotNil(t, cred)
	require.EqualValues(t, val1, SASCredentialGet(cred))
	const val2 = "bar"
	cred.Update(val2)
	require.EqualValues(t, val2, SASCredentialGet(cred))
}

func TestNewRequestFromRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedData := bytes.NewReader([]byte{1, 2, 3, 4, 5})

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", "https://example.com", expectedData)
	require.NoError(t, err)

	req, err := NewRequestFromRequest(httpRequest)
	require.NoError(t, err)

	// our stream has been drained - the func has to make a copy of the body so it can be seekable.
	// so our stream should be at end.
	currentPos, err := expectedData.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(5), currentPos)

	actualData, err := io.ReadAll(req.Body())
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5}, actualData)

	// now we change stuff in the policy.Request...
	replacementBuff := bytes.NewReader([]byte{6})
	err = req.SetBody(NopCloser(replacementBuff), "application/coolstuff")
	require.NoError(t, err)

	// and it's automatically reflected in the http.Request, which helps us with interop
	// with other HTTP pipelines.
	require.Equal(t, "application/coolstuff", httpRequest.Header.Get("Content-Type"))
	newBytes, err := io.ReadAll(httpRequest.Body)
	require.NoError(t, err)
	require.Equal(t, []byte{6}, newBytes)
}

func TestNewRequestFromRequest_AvoidExtraCopyIfReadSeekCloser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedData := NopCloser(bytes.NewReader([]byte{1, 2, 3, 4, 5}))

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", "https://example.com", expectedData)
	require.NoError(t, err)

	req, err := NewRequestFromRequest(httpRequest)
	require.NoError(t, err)

	// our stream should _NOT_ get drained since it was already an io.ReadSeekCloser
	currentPos, err := expectedData.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(0), currentPos)

	actualData, err := io.ReadAll(req.Body())
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5}, actualData)

	// now we change stuff in the policy.Request...
	replacementBuff := bytes.NewReader([]byte{6})
	err = req.SetBody(NopCloser(replacementBuff), "application/coolstuff")
	require.NoError(t, err)

	// and it's automatically reflected in the http.Request, which helps us with interop
	// with other HTTP pipelines.
	require.Equal(t, "application/coolstuff", httpRequest.Header.Get("Content-Type"))
	newBytes, err := io.ReadAll(httpRequest.Body)
	require.NoError(t, err)
	require.Equal(t, []byte{6}, newBytes)
}
