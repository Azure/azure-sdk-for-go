//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package fake

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/require"
)

type scalar struct {
	Value *string
}

type widget struct {
	Name string
}

type widgets struct {
	NextPage *string
	Widgets  []widget
}

func TestNewTokenCredential(t *testing.T) {
	cred := NewTokenCredential()
	require.NotNil(t, cred)

	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{})
	require.NoError(t, err)
	require.NotZero(t, tk)

	myErr := errors.New("failed")
	cred.SetError(myErr)
	tk, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{})
	require.ErrorIs(t, err, myErr)
	require.Zero(t, tk)
}

func TestResponderJSON(t *testing.T) {
	respr := Responder[widget]{}
	respr.SetResponse(widget{Name: "foo"}, nil)
	respr.SetHeader("one", "1")
	respr.SetHeader("two", "2")

	req := &http.Request{}
	resp, err := MarshalResponseAsJSON(GetResponseContent(respr), GetResponse(respr), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req, resp.Request)
	require.Equal(t, "1", resp.Header.Get("one"))
	require.Equal(t, "2", resp.Header.Get("two"))
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, "OK", resp.Status)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	w := widget{}
	require.NoError(t, json.Unmarshal(body, &w))
	require.Equal(t, "foo", w.Name)
}

func TestResponderText(t *testing.T) {
	respr := Responder[scalar]{}
	respr.SetResponse(scalar{Value: to.Ptr("success")}, nil)
	respr.SetHeader("one", "1")
	respr.SetHeader("two", "2")

	req := &http.Request{}
	resp, err := MarshalResponseAsText(GetResponseContent(respr), GetResponse(respr).Value, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req, resp.Request)
	require.Equal(t, "1", resp.Header.Get("one"))
	require.Equal(t, "2", resp.Header.Get("two"))
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, "OK", resp.Status)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	require.Equal(t, "success", string(body))
}

func TestResponderTextNil(t *testing.T) {
	respr := Responder[scalar]{}
	respr.SetResponse(scalar{}, nil)
	respr.SetHeader("one", "1")
	respr.SetHeader("two", "2")

	req := &http.Request{}
	resp, err := MarshalResponseAsText(GetResponseContent(respr), GetResponse(respr).Value, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req, resp.Request)
	require.Equal(t, "1", resp.Header.Get("one"))
	require.Equal(t, "2", resp.Header.Get("two"))
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, "OK", resp.Status)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	require.Empty(t, string(body))
}

func TestResponderXML(t *testing.T) {
	respr := Responder[widget]{}
	respr.SetResponse(widget{Name: "foo"}, nil)
	respr.SetHeader("one", "1")
	respr.SetHeader("two", "2")

	req := &http.Request{}
	resp, err := MarshalResponseAsXML(GetResponseContent(respr), GetResponse(respr), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req, resp.Request)
	require.Equal(t, "1", resp.Header.Get("one"))
	require.Equal(t, "2", resp.Header.Get("two"))
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, "OK", resp.Status)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	w := widget{}
	require.NoError(t, xml.Unmarshal(body, &w))
	require.Equal(t, "foo", w.Name)
}

type badWidget struct {
	Count int
}

func (badWidget) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed")
}

func (*badWidget) UnmarshalJSON([]byte) error {
	return errors.New("failed")
}

func TestResponderMarshallingError(t *testing.T) {
	respr := Responder[badWidget]{}

	req := &http.Request{}
	resp, err := MarshalResponseAsJSON(GetResponseContent(respr), GetResponse(respr), req)
	require.Error(t, err)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)
}

func TestErrorResponder(t *testing.T) {
	req := &http.Request{}

	errResp := ErrorResponder{}
	require.NoError(t, GetError(errResp, req))

	myErr := errors.New("failed")
	errResp.SetError(myErr)
	require.ErrorIs(t, GetError(errResp, req), myErr)

	errResp.SetResponseError("ErrorInvalidWidget", http.StatusBadRequest)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, GetError(errResp, req), &respErr)
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

	require.False(t, PagerResponderMore(&pagerResp))
	resp, err := PagerResponderNext(&pagerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pagerResp.AddError(errors.New("one"))
	pagerResp.AddPage(widgets{
		Widgets: []widget{
			{Name: "foo"},
			{Name: "bar"},
		},
	}, nil)
	pagerResp.AddError(errors.New("two"))
	pagerResp.AddPage(widgets{
		Widgets: []widget{
			{Name: "baz"},
		},
	}, nil)
	pagerResp.AddResponseError("ErrorPagerBlewUp", http.StatusBadRequest)

	PagerResponderInjectNextLinks(&pagerResp, req, func(p *widgets, create func() string) {
		p.NextPage = to.Ptr(create())
	})

	iterations := 0
	for PagerResponderMore(&pagerResp) {
		resp, err := PagerResponderNext(&pagerResp, req)
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
			require.Nil(t, page.NextPage)
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

	pollerResp := PollerResponder[widget]{}

	require.False(t, PollerResponderMore(&pollerResp))
	resp, err := PollerResponderNext(&pollerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddNonTerminalResponse(nil)
	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(nil)
	pollerResp.SetTerminalResponse(widget{Name: "dodo"}, nil)

	iterations := 0
	for PollerResponderMore(&pollerResp) {
		resp, err := PollerResponderNext(&pollerResp, req)
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

	require.False(t, PollerResponderMore(&pollerResp))
	resp, err := PollerResponderNext(&pollerResp, req)
	var nre errorinfo.NonRetriable
	require.ErrorAs(t, err, &nre)
	require.Nil(t, resp)

	pollerResp.AddPollingError(errors.New("network glitch"))
	pollerResp.AddNonTerminalResponse(nil)
	pollerResp.SetTerminalError("ErrorConflictingOperation", http.StatusConflict)

	iterations := 0
	for PollerResponderMore(&pollerResp) {
		resp, err := PollerResponderNext(&pollerResp, req)
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
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	resp, err := NewResponse(ResponseContent{}, req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}

func TestNewBinaryResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	resp, err := NewBinaryResponse(ResponseContent{}, io.NopCloser(strings.NewReader("the body")), req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, "the body", string(body))
}

func TestUnmarshalRequestAsJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", strings.NewReader(`{"Name": "foo"}`))
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err := UnmarshalRequestAsJSON[widget](req)
	require.NoError(t, err)
	require.Equal(t, "foo", w.Name)

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err = UnmarshalRequestAsJSON[widget](req)
	require.NoError(t, err)
	require.Zero(t, w)
}

func TestUnmarshalRequestAsText(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", strings.NewReader("some text"))
	require.NoError(t, err)
	require.NotNil(t, req)

	txt, err := UnmarshalRequestAsText(req)
	require.NoError(t, err)
	require.Equal(t, "some text", txt)

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	txt, err = UnmarshalRequestAsText(req)
	require.NoError(t, err)
	require.Zero(t, txt)
}

func TestUnmarshalRequestAsXML(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", strings.NewReader(`<widget><Name>foo</Name></widget>`))
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err := UnmarshalRequestAsXML[widget](req)
	require.NoError(t, err)
	require.Equal(t, "foo", w.Name)

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err = UnmarshalRequestAsXML[widget](req)
	require.NoError(t, err)
	require.Zero(t, w)
}

func TestUnmarshalRequestAsJSONReadFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", &readFailer{})
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err := UnmarshalRequestAsJSON[widget](req)
	require.Error(t, err)
	require.Zero(t, w)
}

func TestUnmarshalRequestAsJSONUnmarshalFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", strings.NewReader(`{"Name": "foo"}`))
	require.NoError(t, err)
	require.NotNil(t, req)

	w, err := UnmarshalRequestAsJSON[badWidget](req)
	require.Error(t, err)
	require.Zero(t, w)
}

type readFailer struct {
	wrapped io.ReadCloser
}

func (r *readFailer) Close() error {
	return r.wrapped.Close()
}

func (r *readFailer) Read(p []byte) (int, error) {
	return 0, errors.New("mock read failure")
}
