//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

type fakeResult[T any] struct {
	Result T
}

func TestNewResumeToken(t *testing.T) {
	n, err := NewResumeToken[struct{}](fakeResult[struct{}]{})
	require.Error(t, err)
	require.Empty(t, n)
	n, err = NewResumeToken[interface{}](fakeResult[interface{}]{})
	require.Error(t, err)
	require.Empty(t, n)
	n, err = NewResumeToken[int](fakeResult[int]{})
	require.NoError(t, err)
	require.Equal(t, `{"type":"int","token":{"Result":0}}`, n)
	n, err = NewResumeToken[*float64](fakeResult[*float64]{})
	require.NoError(t, err)
	require.Equal(t, `{"type":"*float64","token":{"Result":null}}`, n)
}

func TestExtractToken(t *testing.T) {
	tk, err := ExtractToken("not a JSON object")
	require.Error(t, err)
	require.Nil(t, tk)
	tk, err = ExtractToken(`{ "not": "a token" }`)
	require.Error(t, err)
	require.Nil(t, tk)
	tk, err = ExtractToken(`{"type":"int","token":{"Result":0}}`)
	require.NoError(t, err)
	require.Equal(t, `{"Result":0}`, string(tk))
}

func TestIsTokenValid(t *testing.T) {
	err := IsTokenValid[int]("not a JSON object")
	require.Error(t, err)
	err = IsTokenValid[int](`{ "not": "a token" }`)
	require.Error(t, err)
	err = IsTokenValid[int](`{ "type": 123 }`)
	require.Error(t, err)
	err = IsTokenValid[struct{}](`{ "type": "empty" }`)
	require.Error(t, err)
	err = IsTokenValid[int](`{"type":"*float64","token":{"Result":null}}`)
	require.Error(t, err)
	err = IsTokenValid[int](`{"type":"int","token":{"Result":0}}`)
	require.NoError(t, err)
}

func TestNopPoller(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
	}
	np, err := NewNopPoller[struct{}](resp)
	require.NoError(t, err)
	require.NotNil(t, np)
	require.True(t, np.Done())
	pollResp, err := np.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, resp, pollResp)
	var result struct{}
	err = np.Result(context.Background(), &result)
	require.NoError(t, err)

	resp.StatusCode = http.StatusOK
	np, err = NewNopPoller[struct{}](resp)
	require.NoError(t, err)
	require.NotNil(t, np)
	require.True(t, np.Done())
	pollResp, err = np.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, resp, pollResp)
	err = np.Result(context.Background(), &result)
	require.NoError(t, err)

	resp.Body = io.NopCloser(strings.NewReader(`"value"`))
	np2, err := NewNopPoller[string](resp)
	require.NoError(t, err)
	require.NotNil(t, np2)
	require.True(t, np2.Done())
	pollResp, err = np2.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, resp, pollResp)
	var result2 string
	err = np2.Result(context.Background(), &result2)
	require.NoError(t, err)
	require.Equal(t, "value", result2)
}

func TestPollHelper(t *testing.T) {
	const fakeEndpoint = "https://fake.polling/endpoint"
	err := PollHelper(context.Background(), "invalid endpoint", exported.Pipeline{}, func(*http.Response) (string, error) {
		t.Fatal("shouldn't have been called")
		return "", nil
	})
	require.Error(t, err)

	pl := exported.NewPipeline(shared.TransportFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("failed")
	}))
	err = PollHelper(context.Background(), fakeEndpoint, pl, func(*http.Response) (string, error) {
		t.Fatal("shouldn't have been called")
		return "", nil
	})
	require.Error(t, err)

	require.Error(t, err)
	pl = exported.NewPipeline(shared.TransportFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
		}, nil
	}))
	err = PollHelper(context.Background(), fakeEndpoint, pl, func(*http.Response) (string, error) {
		return "", errors.New("failed")
	})
	require.Error(t, err)

	require.Error(t, err)
	pl = exported.NewPipeline(shared.TransportFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
		}, nil
	}))
	err = PollHelper(context.Background(), fakeEndpoint, pl, func(*http.Response) (string, error) {
		return "inProgress", nil
	})
	require.NoError(t, err)
}

type widget struct {
	Result        string
	Precalculated int
}

func TestResultHelper(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
	}
	var result string
	err := ResultHelper(resp, false, &result)
	require.NoError(t, err)
	require.Empty(t, result)

	resp.StatusCode = http.StatusBadRequest
	resp.Body = io.NopCloser(strings.NewReader(`{ "code": "failed", "message": "bad stuff" }`))
	err = ResultHelper(resp, false, &result)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "failed", respErr.ErrorCode)
	require.Empty(t, result)

	resp.StatusCode = http.StatusOK
	resp.Body = http.NoBody
	err = ResultHelper(resp, false, &result)
	require.NoError(t, err)
	require.Empty(t, result)

	resp.Body = io.NopCloser(strings.NewReader(`{ "Result": "happy" }`))
	widgetResult := widget{Precalculated: 123}
	err = ResultHelper(resp, false, &widgetResult)
	require.NoError(t, err)
	require.Equal(t, "happy", widgetResult.Result)
	require.Equal(t, 123, widgetResult.Precalculated)
}
