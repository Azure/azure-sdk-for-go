//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package op

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
	"github.com/stretchr/testify/require"
)

const (
	fakePollingURL     = "https://foo.bar.baz/status"
	fakeLocationURL    = "https://foo.bar.baz/location"
	fakeResourceURL    = "https://foo.bar.baz/resource"
	fakeResourceLocURL = "https://foo.bar.baz/resourceLocation"
)

func initialResponse(method string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, fakeResourceURL, nil)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Body:    io.NopCloser(body),
		Header:  http.Header{},
		Request: req,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	require.False(t, Applicable(resp), "missing Operation-Location should not be applicable")
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	require.True(t, Applicable(resp), "having Operation-Location should be applicable")
}

func TestCanResume(t *testing.T) {
	token := map[string]interface{}{}
	require.False(t, CanResume(token))
	token["oplocURL"] = fakePollingURL
	require.True(t, CanResume(token))
}

func TestNew(t *testing.T) {
	poller, err := New[struct{}](exported.Pipeline{}, nil, "")
	require.NoError(t, err)
	require.Empty(t, poller.CurState)

	poller, err = New[struct{}](exported.Pipeline{}, &http.Response{Header: http.Header{}}, "")
	require.Error(t, err)
	require.Nil(t, poller)

	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, "this is an invalid polling URL")
	poller, err = New[struct{}](exported.Pipeline{}, resp, "")
	require.Error(t, err)
	require.Nil(t, poller)

	resp = initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, "this is an invalid polling URL")
	poller, err = New[struct{}](exported.Pipeline{}, resp, "")
	require.Error(t, err)
	require.Nil(t, poller)

	resp = initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err = New[struct{}](exported.Pipeline{}, resp, "")
	require.NoError(t, err)
	require.Equal(t, "Updating", poller.CurState)
	require.False(t, poller.Done())
}

type widget struct {
	Shape string `json:"shape"`
}

func TestFinalStateViaLocation(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded" }`)),
			}, nil
		} else if surl == fakeLocationURL {
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

func TestFinalStateViaOperationLocationWithPost(t *testing.T) {
	resp := initialResponse(http.MethodPost, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded", "shape": "rhombus" }`)),
		}, nil
	})), resp, pollers.FinalStateViaOpLocation)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "rhombus", result.Shape)
}

func TestFinalStateViaResourceLocation(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded", "resourceLocation": "https://foo.bar.baz/resourceLocation" }`)),
			}, nil
		} else if surl == fakeResourceLocURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "shape": "square" }`)),
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
	require.Equal(t, "square", result.Shape)
}

func TestResultForPatch(t *testing.T) {
	resp := initialResponse(http.MethodPatch, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded" }`)),
			}, nil
		} else if surl == fakeResourceURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "shape": "square" }`)),
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "square", result.Shape)
}

func TestPostWithLocation(t *testing.T) {
	resp := initialResponse(http.MethodPost, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakePollingURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "status": "succeeded" }`)),
			}, nil
		} else if surl == fakeLocationURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "shape": "triangle" }`)),
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp, "")
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

func TestOperationFailed(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "status": "Failed", "error": { "code": "InvalidSomething" } }`)),
		}, nil
	})), resp, pollers.FinalStateViaLocation)
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
	require.Equal(t, "InvalidSomething", respErr.ErrorCode)
	require.Empty(t, result)
}

func TestPollFailed(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("failed")
	})), resp, pollers.FinalStateViaLocation)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
	require.False(t, poller.Done())
}

func TestPollError(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "NotFound", "message": "the item doesn't exist" } }`)),
		}, nil
	})), resp, pollers.FinalStateViaLocation)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.False(t, poller.Done())
}

func TestMissingStatus(t *testing.T) {
	resp := initialResponse(http.MethodPatch, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "shape": "square" }`)),
		}, nil
	})), resp, "")
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
	require.False(t, poller.Done())
}
