//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/poller"
	"github.com/stretchr/testify/require"
)

const (
	fakeLocationURL  = "https://foo.bar.baz/status"
	fakeLocationURL2 = "https://foo.bar.baz/status/other"
)

func initialResponse() *http.Response {
	return &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
		Body:       http.NoBody,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	require.False(t, Applicable(resp), "missing Location should not be applicable")
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	require.True(t, Applicable(resp), "having Location should be applicable")
}

func TestCanResume(t *testing.T) {
	token := map[string]any{}
	require.False(t, CanResume(token))
	token["type"] = kind
	require.True(t, CanResume(token))
	token["type"] = "something_else"
	require.False(t, CanResume(token))
	token["type"] = 123
	require.False(t, CanResume(token))
}

func TestNew(t *testing.T) {
	poller, err := New[struct{}](exported.Pipeline{}, nil)
	require.NoError(t, err)
	require.Empty(t, poller.CurState)

	poller, err = New[struct{}](exported.Pipeline{}, initialResponse())
	require.Error(t, err)
	require.Nil(t, poller)

	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.NotNil(t, poller)

	resp = initialResponse()
	resp.Header.Set(shared.HeaderLocation, "this is a bad polling URL")
	poller, err = New[struct{}](exported.Pipeline{}, resp)
	require.Error(t, err)
	require.Nil(t, poller)
}

func TestUpdateSucceeded(t *testing.T) {
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       http.NoBody,
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	err = poller.Result(context.Background(), nil)
	require.NoError(t, err)
}

func TestUpdateFailed(t *testing.T) {
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakeLocationURL {
			resp := &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       http.NoBody,
				Header:     http.Header{},
			}
			resp.Header.Set(shared.HeaderLocation, fakeLocationURL2)
			return resp, nil
		} else if surl == fakeLocationURL2 {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       http.NoBody,
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	err = poller.Result(context.Background(), nil)
	require.Error(t, err)
}

func TestUpdateFailedWithProvisioningState(t *testing.T) {
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		if surl := req.URL.String(); surl == fakeLocationURL {
			resp := &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       http.NoBody,
				Header:     http.Header{},
			}
			resp.Header.Set(shared.HeaderLocation, fakeLocationURL2)
			return resp, nil
		} else if surl == fakeLocationURL2 {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "failed" } }`)),
			}, nil
		} else {
			return nil, fmt.Errorf("test bug, unhandled URL %s", surl)
		}
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	err = poller.Result(context.Background(), nil)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
}

func TestSynchronousCompletion(t *testing.T) {
	resp := initialResponse()
	resp.Body = io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "Succeeded" } }`))
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	lp, err := New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, fakeLocationURL, lp.PollURL)
	require.Equal(t, poller.StatusSucceeded, lp.CurState)
	require.True(t, lp.Done())
}

func TestWithThrottling(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, srv.URL())
	lp, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return srv.Do(req)
	})), resp)
	require.NoError(t, err)
	respCount := 0
	for !lp.Done() {
		_, err = lp.Poll(context.Background())
		require.NoError(t, err)
		respCount++
	}
	require.EqualValues(t, 4, respCount)
	require.EqualValues(t, poller.StatusSucceeded, lp.CurState)
}

func TestWithTimeout(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, srv.URL())
	lp, err := New[struct{}](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return srv.Do(req)
	})), resp)
	require.NoError(t, err)
	respCount := 0
	for !lp.Done() {
		_, err = lp.Poll(context.Background())
		require.NoError(t, err)
		respCount++
	}
	require.EqualValues(t, 4, respCount)
	require.EqualValues(t, poller.StatusSucceeded, lp.CurState)
}
