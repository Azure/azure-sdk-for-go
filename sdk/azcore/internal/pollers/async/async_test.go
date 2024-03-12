//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package async

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/poller"
	"github.com/stretchr/testify/require"
)

const (
	fakePollingURL  = "https://foo.bar.baz/status"
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(method string, resp io.Reader) *http.Response {
	req, err := http.NewRequest(method, fakeResourceURL, nil)
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
	require.False(t, Applicable(resp), "missing Azure-AsyncOperation should not be applicable")
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	require.True(t, Applicable(resp), "having Azure-AsyncOperation should be applicable")
}

func TestCanResume(t *testing.T) {
	token := map[string]any{}
	require.False(t, CanResume(token))
	token["asyncURL"] = fakePollingURL
	require.True(t, CanResume(token))
}

func TestNew(t *testing.T) {
	ap, err := New[struct{}](exported.Pipeline{}, nil, "")
	require.NoError(t, err)
	require.Empty(t, ap.CurState)

	ap, err = New[struct{}](exported.Pipeline{}, &http.Response{Header: http.Header{}}, "")
	require.Error(t, err)
	require.Nil(t, ap)

	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, "this is an invalid polling URL")
	ap, err = New[struct{}](exported.Pipeline{}, resp, "")
	require.Error(t, err)
	require.Nil(t, ap)

	resp = initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeResourceURL)
	ap, err = New[struct{}](exported.Pipeline{}, resp, "")
	require.NoError(t, err)
	require.Equal(t, fakePollingURL, ap.AsyncURL)
	require.Equal(t, fakeResourceURL, ap.LocURL)
	require.Equal(t, poller.StatusInProgress, ap.CurState)
	require.False(t, ap.Done())
}

func TestNewDeleteNoProvState(t *testing.T) {
	resp := initialResponse(http.MethodDelete, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[struct{}](exported.Pipeline{}, resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
}

func TestNewPutNoProvState(t *testing.T) {
	// missing provisioning state on initial response
	// NOTE: ARM RPC forbids this but we allow it for back-compat
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[struct{}](exported.Pipeline{}, resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
}

type widget struct {
	Shape string `json:"shape"`
}

func TestFinalGetLocation(t *testing.T) {
	const (
		locURL = "https://foo.bar.baz/location"
	)
	resp := initialResponse(http.MethodPost, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, locURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded" }`)),
			}, nil
		} else if surl == locURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "shape": "triangle" }`)),
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp, pollers.FinalStateViaLocation)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "triangle", result.Shape)
}

func TestFinalGetOrigin(t *testing.T) {
	resp := initialResponse(http.MethodPost, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded" }`)),
			}, nil
		} else if surl == fakeResourceURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "shape": "circle" }`)),
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp, pollers.FinalStateViaOriginalURI)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "circle", result.Shape)
}

func TestNoFinalGet(t *testing.T) {
	resp := initialResponse(http.MethodPost, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded", "shape": "circle" }`)),
		}, nil
	})), resp, pollers.FinalStateViaAzureAsyncOp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "circle", result.Shape)
}

func TestPatchNoFinalGet(t *testing.T) {
	resp := initialResponse(http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded", "shape": "circle" }`)),
		}, nil
	})), resp, pollers.FinalStateViaAzureAsyncOp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "circle", result.Shape)
}

func TestPollFailed(t *testing.T) {
	resp := initialResponse(http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{ "status": "failed" }`)),
		}, nil
	})), resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Empty(t, result)
}

func TestPollError(t *testing.T) {
	resp := initialResponse(http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "NotFound", "message": "the item doesn't exist" } }`)),
		}, nil
	})), resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.False(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.ErrorAs(t, err, &respErr)
	require.Empty(t, result)
}

func TestPollFailedError(t *testing.T) {
	resp := initialResponse(http.MethodPatch, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("failed")
	})), resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSynchronousCompletion(t *testing.T) {
	resp := initialResponse(http.MethodPut, io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "Succeeded" } }`)))
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeResourceURL)
	ap, err := New[struct{}](exported.Pipeline{}, resp, "")
	require.NoError(t, err)
	require.Equal(t, fakePollingURL, ap.AsyncURL)
	require.Equal(t, fakeResourceURL, ap.LocURL)
	require.Equal(t, poller.StatusSucceeded, ap.CurState)
	require.True(t, ap.Done())
}
