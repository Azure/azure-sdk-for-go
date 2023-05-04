//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package body

import (
	"context"
	"errors"
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
		Request: &http.Request{
			Method: http.MethodDelete,
		},
	}
	require.False(t, Applicable(resp), "method DELETE should not be applicable")
	resp.Request.Method = http.MethodPatch
	require.True(t, Applicable(resp), "method PATCH should be applicable")
	resp.Request.Method = http.MethodPut
	require.True(t, Applicable(resp), "method PUT should be applicable")
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
	bp, err := New[struct{}](exported.Pipeline{}, nil)
	require.NoError(t, err)
	require.Empty(t, bp.CurState)

	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	resp.StatusCode = http.StatusCreated
	bp, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, "Started", bp.CurState)

	resp = initialResponse(http.MethodPut, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	resp.StatusCode = http.StatusOK
	bp, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, "Started", bp.CurState)

	resp = initialResponse(http.MethodPut, http.NoBody)
	resp.StatusCode = http.StatusOK
	bp, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, poller.StatusSucceeded, bp.CurState)

	resp = initialResponse(http.MethodPut, http.NoBody)
	resp.StatusCode = http.StatusNoContent
	bp, err = New[struct{}](exported.Pipeline{}, resp)
	require.NoError(t, err)
	require.Equal(t, poller.StatusSucceeded, bp.CurState)
}

type widget struct {
	Shape string `json:"shape"`
}

func TestUpdateNoProvStateFail(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	resp.StatusCode = http.StatusOK
	bp, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, bp.Done())
	resp, err = bp.Poll(context.Background())
	require.ErrorIs(t, err, poller.ErrNoBody)
	require.Nil(t, resp)
	require.False(t, bp.Done())
}

func TestUpdateNoProvStateSuccess(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	resp.StatusCode = http.StatusOK
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{ "shape": "rectangle" }`)),
		}, nil
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, poller.Done())
	var result widget
	err = poller.Result(context.Background(), &result)
	require.NoError(t, err)
	require.Equal(t, "rectangle", result.Shape)
}

func TestUpdateNoProvState204(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	resp.StatusCode = http.StatusOK
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
	require.True(t, poller.Done())
	err = poller.Result(context.Background(), nil)
	require.NoError(t, err)
}

func TestPollFailed(t *testing.T) {
	resp := initialResponse(http.MethodPatch, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{ "properties": { "provisioningState": "failed" } }`)),
		}, nil
	})), resp)
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

func TestPollFailedError(t *testing.T) {
	resp := initialResponse(http.MethodPatch, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("failed")
	})), resp)
	require.NoError(t, err)
	require.False(t, poller.Done())
	resp, err = poller.Poll(context.Background())
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestPollError(t *testing.T) {
	resp := initialResponse(http.MethodPatch, strings.NewReader(`{ "properties": { "provisioningState": "Started" } }`))
	poller, err := New[widget](exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "NotFound", "message": "the item doesn't exist" } }`)),
		}, nil
	})), resp)
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
