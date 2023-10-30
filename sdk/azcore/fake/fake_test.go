//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package fake_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/require"
)

type widget struct {
	Name string
}

type widgets struct {
	NextPage *string
	Widgets  []widget
}

func TestNewTokenCredential(t *testing.T) {
	cred := fake.TokenCredential{}

	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{})
	require.NoError(t, err)
	require.NotZero(t, tk)

	myErr := errors.New("failed")
	cred.SetError(myErr)
	tk, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{})
	require.ErrorIs(t, err, myErr)
	require.Zero(t, tk)
}

func TestResponder(t *testing.T) {
	respr := fake.Responder[widget]{}
	header := http.Header{}
	header.Set("one", "1")
	header.Set("two", "2")
	respr.SetResponse(http.StatusOK, widget{Name: "foo"}, &fake.SetResponseOptions{Header: header})

	req := &http.Request{}
	resp, err := server.MarshalResponseAsJSON(server.GetResponseContent(respr), server.GetResponse(respr), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req, resp.Request)
	require.Equal(t, "1", resp.Header.Get("one"))
	require.Equal(t, "2", resp.Header.Get("two"))
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, "200 OK", resp.Status)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	w := widget{}
	require.NoError(t, json.Unmarshal(body, &w))
	require.Equal(t, "foo", w.Name)
}

func TestErrorResponder(t *testing.T) {
	req := &http.Request{}

	errResp := fake.ErrorResponder{}
	require.NoError(t, server.GetError(errResp, req))

	myErr := errors.New("failed")
	errResp.SetError(myErr)
	require.ErrorIs(t, server.GetError(errResp, req), myErr)

	errResp.SetResponseError(http.StatusBadRequest, "ErrorInvalidWidget")
	var respErr *azcore.ResponseError
	require.ErrorAs(t, server.GetError(errResp, req), &respErr)
	require.Equal(t, "ErrorInvalidWidget", respErr.ErrorCode)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
	require.NotNil(t, respErr.RawResponse)
	require.Equal(t, req, respErr.RawResponse.Request)
}

func unmarshal[T any](resp *http.Response) (T, error) {
	var t T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t, err
	}
	resp.Body.Close()

	err = json.Unmarshal(body, &t)
	return t, err
}

func TestPagerResponder(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	req.URL.Scheme = "http"
	req.URL.Host = "fakehost.org"
	req.URL.Path = "/lister"

	pagerResp := fake.PagerResponder[widgets]{}

	require.False(t, server.PagerResponderMore(&pagerResp))
	resp, err := server.PagerResponderNext(&pagerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pagerResp.AddError(errors.New("one"))
	pagerResp.AddPage(http.StatusOK, widgets{
		Widgets: []widget{
			{Name: "foo"},
			{Name: "bar"},
		},
	}, nil)
	pagerResp.AddError(errors.New("two"))
	pagerResp.AddPage(http.StatusOK, widgets{
		Widgets: []widget{
			{Name: "baz"},
		},
	}, nil)
	pagerResp.AddResponseError(http.StatusBadRequest, "ErrorPagerBlewUp")

	server.PagerResponderInjectNextLinks(&pagerResp, req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations := 0
	for server.PagerResponderMore(&pagerResp) {
		resp, err := server.PagerResponderNext(&pagerResp, req)
		switch iterations {
		case 0:
			require.Error(t, err)
			require.Equal(t, "one", err.Error())
			require.Nil(t, resp)
		case 1:
			require.NoError(t, err)
			require.NotNil(t, resp)
			page, err := unmarshal[widgets](resp)
			require.NoError(t, err)
			require.NotNil(t, page.NextPage)
			require.Equal(t, []widget{{Name: "foo"}, {Name: "bar"}}, page.Widgets)
		case 2:
			require.Error(t, err)
			require.Equal(t, "two", err.Error())
			require.Nil(t, resp)
		case 3:
			require.NoError(t, err)
			require.NotNil(t, resp)
			page, err := unmarshal[widgets](resp)
			require.NoError(t, err)
			require.NotNil(t, page.NextPage)
			require.Equal(t, []widget{{Name: "baz"}}, page.Widgets)
		case 4:
			require.Error(t, err)
			var respErr *azcore.ResponseError
			require.ErrorAs(t, err, &respErr)
			require.Equal(t, "ErrorPagerBlewUp", respErr.ErrorCode)
			require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
			require.Nil(t, resp)
		default:
			t.Fatalf("unexpected case %d", iterations)
		}
		iterations++
	}
	require.Equal(t, 5, iterations)
}

func TestPollerResponder(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	req.URL.Scheme = "http"
	req.URL.Host = "fakehost.org"
	req.URL.Path = "/lro"

	pollerResp := fake.PollerResponder[widget]{}

	require.False(t, server.PollerResponderMore(&pollerResp))
	resp, err := server.PollerResponderNext(&pollerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.SetTerminalResponse(http.StatusOK, widget{Name: "dodo"}, nil)

	iterations := 0
	for server.PollerResponderMore(&pollerResp) {
		resp, err := server.PollerResponderNext(&pollerResp, req)
		switch iterations {
		case 0:
			require.NoError(t, err)
			require.NotNil(t, resp)
		case 1:
			require.Error(t, err)
			require.Nil(t, resp)
		case 2:
			require.NoError(t, err)
			require.NotNil(t, resp)
		case 3:
			require.NoError(t, err)
			require.NotNil(t, resp)
			w, err := unmarshal[widget](resp)
			require.NoError(t, err)
			require.Equal(t, "dodo", w.Name)
		default:
			t.Fatalf("unexpected case %d", iterations)
		}
		iterations++
	}
	require.Equal(t, 4, iterations)
}

func TestPollerResponderTerminalFailure(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	req.URL.Scheme = "http"
	req.URL.Host = "fakehost.org"
	req.URL.Path = "/lro"

	pollerResp := fake.PollerResponder[widget]{}

	require.False(t, server.PollerResponderMore(&pollerResp))
	resp, err := server.PollerResponderNext(&pollerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.SetTerminalError(http.StatusConflict, "ErrorConflictingOperation")

	iterations := 0
	for server.PollerResponderMore(&pollerResp) {
		resp, err := server.PollerResponderNext(&pollerResp, req)
		switch iterations {
		case 0:
			require.Error(t, err)
			require.Nil(t, resp)
		case 1:
			require.NoError(t, err)
			require.NotNil(t, resp)
		case 2:
			require.Error(t, err)
			require.Nil(t, resp)
			var respErr *azcore.ResponseError
			require.ErrorAs(t, err, &respErr)
			require.Equal(t, "ErrorConflictingOperation", respErr.ErrorCode)
			require.Equal(t, http.StatusConflict, respErr.StatusCode)
			require.Equal(t, req, respErr.RawResponse.Request)
		default:
			t.Fatalf("unexpected case %d", iterations)
		}
		iterations++
	}
	require.Equal(t, 3, iterations)
}
