//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
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
