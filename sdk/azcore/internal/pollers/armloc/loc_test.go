//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armloc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

const (
	fakePollingURL1 = "https://foo.bar.baz/status"
)

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	require.False(t, Applicable(resp), "missing Location should not be applicable")
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	require.True(t, Applicable(resp), "having Location should be applicable")
}

func TestCanResume(t *testing.T) {
	token := map[string]interface{}{}
	require.False(t, CanResume(token))
	token["type"] = kind
	require.True(t, CanResume(token))
	token["type"] = "something_else"
	require.False(t, CanResume(token))
	token["type"] = 123
	require.False(t, CanResume(token))
}

func TestNew(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, pollers.StatusInProgress, poller.CurState)
	require.Equal(t, fakePollingURL1, poller.PollURL)
	poller, err = New[struct{}](exported.Pipeline{}, nil)
	require.NoError(t, err)
	require.Equal(t, "", poller.CurState)
	poller, err = New[struct{}](exported.Pipeline{}, &http.Response{Header: http.Header{}})
	require.Error(t, err)
	require.Nil(t, poller)
	resp.Header.Set(shared.HeaderLocation, "this is a bad polling URL")
	poller, err = New[struct{}](exported.Pipeline{}, resp)
	require.Error(t, err)
	require.Nil(t, poller)
}

func TestProvisioningStateSuccessNoContent(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "Succeeded" } }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.True(t, poller.Done())
	result, err := poller.Result(context.Background(), nil)
	require.NoError(t, err)
	require.Empty(t, result)
}

type widget struct {
	Shape string `json:"shape"`
}

func TestProvisioningStateSuccess(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "Succeeded" }, "shape": "round" }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	result, err := poller.Result(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, "round", result.Shape)
}

func TestProvisioningStateSuccessNoProvisioningState(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "shape": "round" }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	result, err := poller.Result(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, "round", result.Shape)
}

func TestPollFailedBadRequest(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{},
			Body:       http.NoBody,
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.True(t, poller.Done())
	result, err := poller.Result(context.Background(), nil)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Empty(t, result)
}

func TestPollFailedError(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("failed")
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
}
