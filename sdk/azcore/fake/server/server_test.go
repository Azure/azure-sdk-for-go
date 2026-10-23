// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package server

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
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

type sseTestValue struct {
	eventType string
	data      string
}

func encodeSSETestValue(v sseTestValue) (streaming.SSEFrame, error) {
	return streaming.SSEFrame{Type: v.eventType, Data: []byte(v.data)}, nil
}

func TestMarshalSSEResponder(t *testing.T) {
	var responder fake.SSEResponder[sseTestValue]
	responder.AddEvent(sseTestValue{eventType: "responseCreated", data: `{"id":"resp_1"}`})
	responder.AddEvent(sseTestValue{eventType: "responseDelta", data: `{"delta":"Hello"}`})
	responder.AddFrame(streaming.SSEFrame{ID: to.Ptr("event-1"), Type: "message", Retry: to.Ptr(1000), Data: []byte(`{"message":"hello"}`)})

	body, err := MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)

	reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
		return body, nil
	}, streaming.SSEHandler[streaming.SSEFrame]{
		Decode: func(f streaming.SSEFrame) (streaming.SSEFrame, bool, error) {
			return f, false, nil
		},
	}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	var got []streaming.SSEFrame
	for {
		f, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		got = append(got, f)
	}
	require.Len(t, got, 3)
	require.Equal(t, "responseCreated", got[0].Type)
	require.JSONEq(t, `{"id":"resp_1"}`, string(got[0].Data))
	require.Equal(t, "responseDelta", got[1].Type)
	require.Equal(t, "message", got[2].Type)
	require.Equal(t, to.Ptr("event-1"), got[2].ID)
	require.Equal(t, to.Ptr(1000), got[2].Retry)
}

func TestMarshalSSEResponderMetadata(t *testing.T) {
	var responder fake.SSEResponder[sseTestValue]
	responder.AddFrame(streaming.SSEFrame{ID: to.Ptr("checkpoint"), Retry: to.Ptr(1000)})
	responder.AddFrame(streaming.SSEFrame{Retry: to.Ptr(0)})
	responder.AddFrame(streaming.SSEFrame{Data: []byte{}})
	responder.AddEvent(sseTestValue{data: "hello\nworld"})

	body, err := MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)
	wire, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())
	require.Equal(t, "id: checkpoint\nretry: 1000\n\nretry: 0\n\ndata: \n\ndata: hello\ndata: world\n\n", string(wire))

	body, err = MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)
	var decoded []streaming.SSEFrame
	reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
		return body, nil
	}, streaming.SSEHandler[streaming.SSEFrame]{
		Decode: func(frame streaming.SSEFrame) (streaming.SSEFrame, bool, error) {
			decoded = append(decoded, frame)
			return frame, false, nil
		},
	}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()

	frame, err := reader.Next()
	require.NoError(t, err)
	require.NotNil(t, frame.Data)
	require.Empty(t, frame.Data)
	require.Nil(t, frame.ID)
	require.Nil(t, frame.Retry)
	frame, err = reader.Next()
	require.NoError(t, err)
	require.Equal(t, "hello\nworld", string(frame.Data))
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
	require.Len(t, decoded, 2, "metadata-only frames must not invoke the decoder")
	require.Equal(t, "checkpoint", reader.LastEventID())
}

func TestMarshalSSEResponderZeroRetry(t *testing.T) {
	var responder fake.SSEResponder[sseTestValue]
	responder.AddFrame(streaming.SSEFrame{Retry: to.Ptr(1000), Data: []byte("first")})
	responder.AddFrame(streaming.SSEFrame{Retry: to.Ptr(0), Data: []byte("second")})
	responder.AddEvent(sseTestValue{data: "third"})
	body, err := MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)
	reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
		return body, nil
	}, streaming.SSEHandler[streaming.SSEFrame]{
		Decode: func(frame streaming.SSEFrame) (streaming.SSEFrame, bool, error) {
			return frame, false, nil
		},
	}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	for _, retry := range []*int{to.Ptr(1000), to.Ptr(0), nil} {
		frame, err := reader.Next()
		require.NoError(t, err)
		require.Equal(t, retry, frame.Retry)
	}
}

func TestMarshalSSEResponderDataLineEndings(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
		want string
	}{
		{"empty", "", ""},
		{"LF", "hello\nworld", "hello\nworld"},
		{"CRLF", "hello\r\nworld", "hello\nworld"},
		{"CR", "hello\rworld", "hello\nworld"},
		{"mixed", "one\r\ntwo\rthree\nfour", "one\ntwo\nthree\nfour"},
		{"leading CR", "\rhello", "\nhello"},
		{"trailing CR", "hello\r", "hello\n"},
		{"trailing CRLF", "hello\r\n", "hello\n"},
		{"trailing LF", "hello\n", "hello\n"},
		{"consecutive CR", "hello\r\rworld", "hello\n\nworld"},
		{"consecutive CRLF", "hello\r\n\r\nworld", "hello\n\nworld"},
		{"only CR", "\r", "\n"},
		{"only CRLF", "\r\n", "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, typed := range []bool{false, true} {
				name := "raw"
				if typed {
					name = "typed"
				}
				t.Run(name, func(t *testing.T) {
					var responder fake.SSEResponder[sseTestValue]
					if typed {
						responder.AddEvent(sseTestValue{data: tc.data})
					} else {
						responder.AddFrame(streaming.SSEFrame{Data: []byte(tc.data)})
					}
					responder.AddEvent(sseTestValue{data: "next event"})
					body, err := MarshalSSEResponder(&responder, encodeSSETestValue)
					require.NoError(t, err)
					wire, err := io.ReadAll(body)
					require.NoError(t, err)
					require.NoError(t, body.Close())
					require.Equal(t, "data: "+strings.ReplaceAll(tc.want, "\n", "\ndata: ")+"\n\ndata: next event\n\n", string(wire))

					body, err = MarshalSSEResponder(&responder, encodeSSETestValue)
					require.NoError(t, err)
					reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
						return body, nil
					}, streaming.SSEHandler[streaming.SSEFrame]{
						Decode: func(frame streaming.SSEFrame) (streaming.SSEFrame, bool, error) {
							return frame, false, nil
						},
					}, nil)
					require.NoError(t, err)
					defer func() {
						_ = reader.Close()
					}()
					for _, want := range []string{tc.want, "next event"} {
						frame, err := reader.Next()
						require.NoError(t, err)
						require.NotNil(t, frame.Data)
						require.Equal(t, want, string(frame.Data))
					}
					_, err = reader.Next()
					require.ErrorIs(t, err, io.EOF)
				})
			}
		})
	}
}

func TestMarshalSSEResponderIDReset(t *testing.T) {
	var responder fake.SSEResponder[sseTestValue]
	responder.AddFrame(streaming.SSEFrame{ID: to.Ptr("checkpoint"), Data: []byte("set")})
	responder.AddEvent(sseTestValue{data: "inherit"})
	responder.AddFrame(streaming.SSEFrame{ID: to.Ptr("")})
	responder.AddEvent(sseTestValue{data: "reset"})
	responder.AddEvent(sseTestValue{data: "after reset"})
	body, err := MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)
	wire, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())
	require.Equal(t, "id: checkpoint\ndata: set\n\ndata: inherit\n\nid: \n\ndata: reset\n\ndata: after reset\n\n", string(wire))
	body, err = MarshalSSEResponder(&responder, encodeSSETestValue)
	require.NoError(t, err)
	reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
		return body, nil
	}, streaming.SSEHandler[streaming.SSEFrame]{
		Decode: func(frame streaming.SSEFrame) (streaming.SSEFrame, bool, error) {
			return frame, false, nil
		},
	}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	for _, cursor := range []string{"checkpoint", "checkpoint", "", ""} {
		_, err := reader.Next()
		require.NoError(t, err)
		require.Equal(t, cursor, reader.LastEventID())
	}
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
}
