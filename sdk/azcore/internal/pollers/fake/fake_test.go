//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package fake

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/poller"
	"github.com/stretchr/testify/require"
)

const (
	fakePollingURL  = "https://foo.bar.baz/status"
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(ctx context.Context, method string, resp io.Reader) *http.Response {
	req, err := http.NewRequestWithContext(ctx, method, fakeResourceURL, nil)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Body:    io.NopCloser(resp),
		Header:  http.Header{},
		Request: req,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	require.False(t, Applicable(resp), "missing Fake-Poller-Status should not be applicable")
	resp.Header.Set(shared.HeaderFakePollerStatus, fakePollingURL)
	require.True(t, Applicable(resp), "having Fake-Poller-Status should be applicable")
}

func TestCanResume(t *testing.T) {
	token := map[string]any{}
	require.False(t, CanResume(token))
	token["fakeURL"] = fakePollingURL
	require.True(t, CanResume(token))
}

func TestNew(t *testing.T) {
	fp, err := New[struct{}](exported.Pipeline{}, nil)
	require.NoError(t, err)
	require.Empty(t, fp.FakeStatus)

	fp, err = New[struct{}](exported.Pipeline{}, &http.Response{Header: http.Header{}})
	require.Error(t, err)
	require.Nil(t, fp)

	resp := initialResponse(context.Background(), http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, "faking")
	fp, err = New[struct{}](exported.Pipeline{}, resp)
	require.Error(t, err)
	require.Nil(t, fp)

	resp = initialResponse(context.WithValue(context.Background(), shared.CtxAPINameKey{}, 123), http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, "faking")
	fp, err = New[struct{}](exported.Pipeline{}, resp)
	require.Error(t, err)
	require.Nil(t, fp)

	resp = initialResponse(context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI"), http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, "faking")
	fp, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.NotNil(t, fp)
	require.False(t, fp.Done())
}

func TestSynchronousCompletion(t *testing.T) {
	resp := initialResponse(context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI"), http.MethodPut, http.NoBody)
	resp.StatusCode = http.StatusNoContent
	resp.Header.Set(shared.HeaderFakePollerStatus, poller.StatusSucceeded)
	fp, err := New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, poller.StatusSucceeded, fp.FakeStatus)
	require.True(t, fp.Done())
	require.NoError(t, fp.Result(context.Background(), nil))
}

type widget struct {
	Shape string `json:"shape"`
}

func TestPollSucceeded(t *testing.T) {
	pollCtx := context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI")
	resp := initialResponse(pollCtx, http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, poller.StatusInProgress)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{shared.HeaderFakePollerStatus: []string{"Succeeded"}},
			Body:       io.NopCloser(strings.NewReader(`{ "shape": "triangle" }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	sanitizedPollerPath := SanitizePollerPath(poller.FakeURL)
	require.NotEqualValues(t, sanitizedPollerPath, poller.FakeStatus)
	require.EqualValues(t, fakeResourceURL, sanitizedPollerPath)
	require.True(t, strings.HasPrefix(poller.FakeURL, sanitizedPollerPath))
	resp, err = poller.Poll(pollCtx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	require.NoError(t, poller.Result(context.Background(), &result))
	require.EqualValues(t, "triangle", result.Shape)
}

func TestPollError(t *testing.T) {
	pollCtx := context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI")
	resp := initialResponse(pollCtx, http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, poller.StatusInProgress)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{shared.HeaderFakePollerStatus: []string{poller.StatusFailed}},
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "NotFound", "message": "the item doesn't exist" } }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(pollCtx)
	require.Error(t, err)
	require.Nil(t, resp)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.False(t, poller.Done())
	var result widget
	require.Error(t, poller.Result(context.Background(), &result))
	require.ErrorAs(t, err, &respErr)
}

func TestPollFailed(t *testing.T) {
	pollCtx := context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI")
	resp := initialResponse(pollCtx, http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, poller.StatusInProgress)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{shared.HeaderFakePollerStatus: []string{poller.StatusFailed}},
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "FakeFailure", "message": "couldn't do the thing" } }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(pollCtx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	var respErr *exported.ResponseError
	err = poller.Result(context.Background(), &result)
	require.Error(t, err)
	require.ErrorAs(t, err, &respErr)
}

func TestPollErrorNoHeader(t *testing.T) {
	pollCtx := context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI")
	resp := initialResponse(pollCtx, http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderFakePollerStatus, poller.StatusInProgress)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "NotFound", "message": "the item doesn't exist" } }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(pollCtx)
	require.Error(t, err)
	require.Nil(t, resp)
}
