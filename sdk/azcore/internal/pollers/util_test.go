//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestIsTerminalState(t *testing.T) {
	require.False(t, IsTerminalState("Updating"), "Updating is not a terminal state")
	require.True(t, IsTerminalState("Succeeded"), "Succeeded is a terminal state")
	require.True(t, IsTerminalState("failed"), "failed is a terminal state")
	require.True(t, IsTerminalState("canceled"), "canceled is a terminal state")
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

func TestIsValidURL(t *testing.T) {
	require.False(t, IsValidURL("/foo"))
	require.True(t, IsValidURL("https://foo.bar/baz"))
}

func TestFailed(t *testing.T) {
	require.False(t, Failed("Succeeded"))
	require.False(t, Failed("Updating"))
	require.True(t, Failed("failed"))
}

func TestGetJSON(t *testing.T) {
	j, err := GetJSON(&http.Response{Body: http.NoBody})
	require.ErrorIs(t, err, ErrNoBody)
	require.Nil(t, j)
	j, err = GetJSON(&http.Response{Body: ioutil.NopCloser(strings.NewReader(`{ "foo": "bar" }`))})
	require.NoError(t, err)
	require.Equal(t, "bar", j["foo"])
}

func TestGetStatusSuccess(t *testing.T) {
	const jsonBody = `{ "status": "InProgress" }`
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
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
		Body: ioutil.NopCloser(strings.NewReader("{}")),
	}
	status, err := GetStatus(resp)
	require.NoError(t, err)
	require.Empty(t, status)
}

func TestGetProvisioningState(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Canceled" } }`
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
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
		Body: ioutil.NopCloser(strings.NewReader("{}")),
	}
	state, err := GetProvisioningState(resp)
	require.NoError(t, err)
	require.Empty(t, state)
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
		return &http.Response{}, nil
	}))
	err = PollHelper(context.Background(), fakeEndpoint, pl, func(*http.Response) (string, error) {
		return "", errors.New("failed")
	})
	require.Error(t, err)

	require.Error(t, err)
	pl = exported.NewPipeline(shared.TransportFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{}, nil
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
