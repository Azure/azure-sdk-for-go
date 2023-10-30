//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
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

func TestResponder(t *testing.T) {
	respr := Responder[widget]{}
	header := http.Header{}
	header.Set("one", "1")
	header.Set("two", "2")
	thing := widget{Name: "foo"}
	respr.SetResponse(http.StatusOK, thing, &SetResponseOptions{Header: header})
	require.EqualValues(t, thing, respr.GetResponse())
	require.EqualValues(t, http.StatusOK, respr.GetResponseContent().HTTPStatus)
	require.EqualValues(t, header, respr.GetResponseContent().Header)
}

func TestErrorResponder(t *testing.T) {
	req := &http.Request{}

	errResp := ErrorResponder{}
	require.NoError(t, errResp.GetError(req))

	myErr := errors.New("failed")
	errResp.SetError(myErr)
	require.ErrorIs(t, errResp.GetError(req), myErr)

	errResp.SetResponseError(http.StatusBadRequest, "ErrorInvalidWidget")
	var respErr *azcore.ResponseError
	require.ErrorAs(t, errResp.GetError(req), &respErr)
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

	pagerResp := PagerResponder[widgets]{}

	require.False(t, pagerResp.More())
	resp, err := pagerResp.Next(req)
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

	pagerResp.InjectNextLinks(req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations := 0
	for pagerResp.More() {
		resp, err := pagerResp.Next(req)
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
			sanitizedNextPage := SanitizePagerPath(*page.NextPage)
			require.NotEqualValues(t, sanitizedNextPage, *page.NextPage)
			require.True(t, strings.HasPrefix(*page.NextPage, sanitizedNextPage))
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

	// single page with subsequent error
	pagerResp = PagerResponder[widgets]{}

	pagerResp.AddPage(http.StatusOK, widgets{
		Widgets: []widget{
			{Name: "foo"},
			{Name: "bar"},
		},
	}, nil)
	pagerResp.AddError(errors.New("two"))

	pagerResp.InjectNextLinks(req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations = 0
	for pagerResp.More() {
		resp, err := pagerResp.Next(req)
		switch iterations {
		case 0:
			require.NoError(t, err)
			require.NotNil(t, resp)
			page, err := unmarshal[widgets](resp)
			require.NoError(t, err)
			require.NotNil(t, page.NextPage)
			require.Equal(t, []widget{{Name: "foo"}, {Name: "bar"}}, page.Widgets)
		case 1:
			require.Error(t, err)
			require.Nil(t, resp)
		}
		iterations++
	}
	require.EqualValues(t, 2, iterations)

	// single page with subsequent response error
	pagerResp = PagerResponder[widgets]{}

	pagerResp.AddPage(http.StatusOK, widgets{
		Widgets: []widget{
			{Name: "foo"},
			{Name: "bar"},
		},
	}, nil)
	pagerResp.AddResponseError(http.StatusBadRequest, "BadRequest")

	pagerResp.InjectNextLinks(req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations = 0
	for pagerResp.More() {
		resp, err := pagerResp.Next(req)
		switch iterations {
		case 0:
			require.NoError(t, err)
			require.NotNil(t, resp)
			page, err := unmarshal[widgets](resp)
			require.NoError(t, err)
			require.NotNil(t, page.NextPage)
			require.Equal(t, []widget{{Name: "foo"}, {Name: "bar"}}, page.Widgets)
		case 1:
			require.Error(t, err)
			require.Nil(t, resp)
		}
		iterations++
	}
	require.EqualValues(t, 2, iterations)

	// single page
	pagerResp = PagerResponder[widgets]{}

	pagerResp.AddPage(http.StatusOK, widgets{
		Widgets: []widget{
			{Name: "foo"},
			{Name: "bar"},
		},
	}, nil)

	pagerResp.InjectNextLinks(req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations = 0
	for pagerResp.More() {
		resp, err := pagerResp.Next(req)
		switch iterations {
		case 0:
			require.NoError(t, err)
			require.NotNil(t, resp)
			page, err := unmarshal[widgets](resp)
			require.NoError(t, err)
			require.Nil(t, page.NextPage)
			require.Equal(t, []widget{{Name: "foo"}, {Name: "bar"}}, page.Widgets)
		}
		iterations++
	}
	require.EqualValues(t, 1, iterations)
}

func TestPollerResponder(t *testing.T) {
	req := &http.Request{URL: &url.URL{}}
	req.URL.Scheme = "http"
	req.URL.Host = "fakehost.org"
	req.URL.Path = "/lro"

	pollerResp := PollerResponder[widget]{}

	require.False(t, pollerResp.More())
	resp, err := pollerResp.Next(req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.SetTerminalResponse(http.StatusOK, widget{Name: "dodo"}, nil)

	iterations := 0
	for pollerResp.More() {
		resp, err := pollerResp.Next(req)
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

	pollerResp := PollerResponder[widget]{}

	require.False(t, pollerResp.More())
	resp, err := pollerResp.Next(req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(http.StatusOK, nil)
	pollerResp.SetTerminalError(http.StatusConflict, "ErrorConflictingOperation")

	iterations := 0
	for pollerResp.More() {
		resp, err := pollerResp.Next(req)
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

func TestNewResponse(t *testing.T) {
	resp, err := NewResponse(ResponseContent{}, nil)
	require.Error(t, err)
	require.Nil(t, resp)

	resp, err = NewResponse(ResponseContent{HTTPStatus: http.StatusNoContent}, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, resp.Header)
}

func TestNewErrorResponse(t *testing.T) {
	resp, err := newErrorResponse(0, "", nil)
	require.Error(t, err)
	require.Nil(t, resp)
	const errorCode = "YouCantDoThat"
	resp, err = newErrorResponse(http.StatusForbidden, errorCode, nil)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusForbidden, resp.StatusCode)
	require.EqualValues(t, errorCode, resp.Header.Get(shared.HeaderXMSErrorCode))
}
