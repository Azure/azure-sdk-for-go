//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package poller

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsTerminalState(t *testing.T) {
	require.False(t, IsTerminalState("upDAting"), "Updating is not a terminal state")
	require.True(t, IsTerminalState("SuccEEded"), "Succeeded is a terminal state")
	require.True(t, IsTerminalState("completEd"), "Completed is a terminal state")
	require.True(t, IsTerminalState("faIled"), "failed is a terminal state")
	require.True(t, IsTerminalState("canCeled"), "canceled is a terminal state")
	require.True(t, IsTerminalState("canceLLed"), "cancelled is a terminal state")
}

func TestStatusCodeValid(t *testing.T) {
	require.True(t, StatusCodeValid(&http.Response{StatusCode: http.StatusOK}))
	require.True(t, StatusCodeValid(&http.Response{StatusCode: http.StatusAccepted}))
	require.True(t, StatusCodeValid(&http.Response{StatusCode: http.StatusCreated}))
	require.True(t, StatusCodeValid(&http.Response{StatusCode: http.StatusNoContent}))
	require.False(t, StatusCodeValid(&http.Response{StatusCode: http.StatusPartialContent}))
	require.False(t, StatusCodeValid(&http.Response{StatusCode: http.StatusBadRequest}))
	require.False(t, StatusCodeValid(&http.Response{StatusCode: http.StatusInternalServerError}))
}

func TestIsValidURL(t *testing.T) {
	require.False(t, IsValidURL("/foo"))
	require.True(t, IsValidURL("https://foo.bar/baz"))
}

func TestFailed(t *testing.T) {
	require.False(t, Failed("sUcceeded"))
	require.False(t, Failed("ppdATing"))
	require.True(t, Failed("fAilEd"))
	require.True(t, Failed("caNcElled"))
}

func TestSucceeded(t *testing.T) {
	require.True(t, Succeeded("Succeeded"))
	require.False(t, Succeeded("Updating"))
	require.False(t, Succeeded("failed"))
}

func TestGetJSON(t *testing.T) {
	j, err := GetJSON(&http.Response{Body: http.NoBody})
	require.ErrorIs(t, err, ErrNoBody)
	require.Nil(t, j)
	j, err = GetJSON(&http.Response{Body: io.NopCloser(strings.NewReader(`{ "foo": "bar" }`))})
	require.NoError(t, err)
	require.Equal(t, "bar", j["foo"])
}

func TestGetStatusSuccess(t *testing.T) {
	const jsonBody = `{ "status": "InProgress" }`
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsonBody)),
	}
	status, err := GetStatus(resp)
	require.NoError(t, err)
	require.Equal(t, "InProgress", status)
}

func TestGetNoBody(t *testing.T) {
	resp := &http.Response{
		Body: http.NoBody,
	}
	status, err := GetStatus(resp)
	require.ErrorIs(t, err, ErrNoBody)
	require.Empty(t, status)
	status, err = GetProvisioningState(resp)
	require.ErrorIs(t, err, ErrNoBody)
	require.Empty(t, status)
}

func TestGetStatusError(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader("{}")),
	}
	status, err := GetStatus(resp)
	require.NoError(t, err)
	require.Empty(t, status)
}

func TestGetProvisioningState(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Canceled" } }`
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsonBody)),
	}
	state, err := GetProvisioningState(resp)
	require.NoError(t, err)
	require.Equal(t, "Canceled", state)
}

func TestGetResourceLocation(t *testing.T) {
	resp := &http.Response{
		Body: http.NoBody,
	}
	resLoc, err := GetResourceLocation(resp)
	require.Error(t, err)
	require.Empty(t, resLoc)
	resp.Body = io.NopCloser(strings.NewReader(`{"status": "succeeded"}`))
	resLoc, err = GetResourceLocation(resp)
	require.NoError(t, err)
	require.Empty(t, resLoc)
	resp.Body = io.NopCloser(strings.NewReader(`{"resourceLocation": 123}`))
	resLoc, err = GetResourceLocation(resp)
	require.Error(t, err)
	require.Empty(t, resLoc)
	resp.Body = io.NopCloser(strings.NewReader(`{"resourceLocation": "here"}`))
	resLoc, err = GetResourceLocation(resp)
	require.NoError(t, err)
	require.Equal(t, "here", resLoc)
}

func TestGetProvisioningStateError(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader("{}")),
	}
	state, err := GetProvisioningState(resp)
	require.NoError(t, err)
	require.Empty(t, state)
}
