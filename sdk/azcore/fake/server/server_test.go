//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package server

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/require"
)

type widget struct {
	Name string
}

type badWidget struct {
	Count int
}

func (badWidget) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed")
}

func (badWidget) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return errors.New("failed")
}

func (*badWidget) UnmarshalJSON([]byte) error {
	return errors.New("failed")
}

func (*badWidget) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	return errors.New("failed")
}

func TestNewResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	resp, err := NewResponse(ResponseContent{HTTPStatus: http.StatusNoContent}, req, nil)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
}

func TestNewResponseWithOptions(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	resp, err := NewResponse(ResponseContent{HTTPStatus: http.StatusOK}, req, &ResponseOptions{
		Body:        io.NopCloser(strings.NewReader("the body")),
		ContentType: shared.ContentTypeTextPlain,
	})
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, shared.ContentTypeTextPlain, resp.Header.Get(shared.HeaderContentType))
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, "the body", string(body))
}

func TestMarshalUnmarshalAsJSON(t *testing.T) {
	thing := widget{Name: "foo"}
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	resp, err := MarshalResponseAsJSON(ResponseContent{HTTPStatus: http.StatusOK}, thing, req)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, `{"Name":"foo"}`, string(body))

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", io.NopCloser(bytes.NewReader(body)))
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

func TestMarshalUnmarshalAsText(t *testing.T) {
	const thing = "some text"
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	resp, err := MarshalResponseAsText(ResponseContent{HTTPStatus: http.StatusOK}, to.Ptr(thing), req)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, thing, string(body))

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", io.NopCloser(bytes.NewReader(body)))
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

func TestMarshalUnmarshalAsXML(t *testing.T) {
	thing := widget{Name: "foo"}
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	resp, err := MarshalResponseAsXML(ResponseContent{HTTPStatus: http.StatusOK}, thing, req)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, `<widget><Name>foo</Name></widget>`, string(body))

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", io.NopCloser(bytes.NewReader(body)))
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

func TestUnmarshalRequestReadFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", &readFailer{})
	require.NoError(t, err)
	require.NotNil(t, req)

	var nre errorinfo.NonRetriable

	b, err := UnmarshalRequestAsByteArray(req, exported.Base64StdFormat)
	require.Error(t, err)
	require.Zero(t, b)
	require.ErrorAs(t, err, &nre)

	w, err := UnmarshalRequestAsJSON[widget](req)
	require.Error(t, err)
	require.Zero(t, w)
	require.ErrorAs(t, err, &nre)

	s, err := UnmarshalRequestAsText(req)
	require.Error(t, err)
	require.Zero(t, s)
	require.ErrorAs(t, err, &nre)

	w, err = UnmarshalRequestAsXML[widget](req)
	require.Error(t, err)
	require.Zero(t, w)
	require.ErrorAs(t, err, &nre)
}

func TestMarshalUnmarshalFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", strings.NewReader("won't get here"))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err := MarshalResponseAsJSON(ResponseContent{}, badWidget{}, nil)
	require.Error(t, err)
	require.Nil(t, resp)

	resp, err = MarshalResponseAsXML(ResponseContent{}, badWidget{}, nil)
	require.Error(t, err)
	require.Nil(t, resp)

	w, err := UnmarshalRequestAsJSON[badWidget](req)
	require.Error(t, err)
	require.Zero(t, w)

	w, err = UnmarshalRequestAsXML[badWidget](req)
	require.Error(t, err)
	require.Zero(t, w)
}

func TestMarshalUnmarshalAsByteArray(t *testing.T) {
	const encodeVal = "encode me"
	req, err := http.NewRequest(http.MethodPut, "https://foo.bar/baz", nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	body, err := UnmarshalRequestAsByteArray(req, exported.Base64StdFormat)
	require.NoError(t, err)
	require.Nil(t, body)
	resp, err := MarshalResponseAsByteArray(ResponseContent{HTTPStatus: http.StatusOK}, []byte(encodeVal), exported.Base64StdFormat, req)
	require.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.EqualValues(t, "ZW5jb2RlIG1l", string(body))

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", io.NopCloser(bytes.NewReader(body)))
	require.NoError(t, err)
	require.NotNil(t, req)
	body, err = UnmarshalRequestAsByteArray(req, exported.Base64StdFormat)
	require.NoError(t, err)
	require.EqualValues(t, encodeVal, string(body))

	req, err = http.NewRequest(http.MethodPut, "https://foo.bar/baz", io.NopCloser(strings.NewReader("not base64 encoded")))
	require.NoError(t, err)
	require.NotNil(t, req)
	body, err = UnmarshalRequestAsByteArray(req, exported.Base64StdFormat)
	require.Error(t, err)
	require.Nil(t, body)
}

func TestResponderHelpers(t *testing.T) {
	thing := widget{Name: "foo"}
	header := http.Header{
		"header1": []string{"value1"},
	}
	respr := fake.Responder[widget]{}
	respr.SetResponse(http.StatusOK, thing, &fake.SetResponseOptions{Header: header})
	require.EqualValues(t, thing, GetResponse(respr))
	require.EqualValues(t, http.StatusOK, GetResponseContent(respr).HTTPStatus)
	require.EqualValues(t, header, GetResponseContent(respr).Header)
}

func TestErrorResponderHelpers(t *testing.T) {
	errResp := fake.ErrorResponder{}
	errResp.SetError(io.EOF)
	require.ErrorIs(t, GetError(errResp, nil), io.EOF)
}

func TestPagerResponderHelpers(t *testing.T) {
	pagerResp := fake.PagerResponder[widget]{}
	require.False(t, PagerResponderMore(&pagerResp))
	resp, err := PagerResponderNext(&pagerResp, nil)
	require.Error(t, err)
	require.Nil(t, resp)
	PagerResponderInjectNextLinks(&pagerResp, nil, func(page *widget, createLink func() string) {})
}

func TestPollerResponderHelpers(t *testing.T) {
	pollerResp := fake.PollerResponder[widget]{}
	require.False(t, PollerResponderMore(&pollerResp))
	resp, err := PollerResponderNext(&pollerResp, nil)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSanitizePagerPollerPath(t *testing.T) {
	const untouched = "/this/path/wont/change"
	require.EqualValues(t, untouched, SanitizePagerPollerPath(untouched))
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
