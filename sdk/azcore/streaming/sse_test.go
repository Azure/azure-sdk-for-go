// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package streaming_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

// frameView is a decoded SSE frame used only by these tests. It captures the
// wire envelope plus the parsed JSON payload so the runtime can be exercised
// without depending on any generated event union.
type frameView struct {
	eventType string
	data      string
	retry     int
	fields    map[string]any
}

// decodeView is a generic decoder: JSON object payloads are unmarshaled into
// fields, the "[DONE]" sentinel is treated as terminal, and everything else is
// surfaced as raw text.
func decodeView(f streaming.SSEFrame) (frameView, bool, error) {
	v := frameView{eventType: f.Type, data: string(f.Data), retry: f.Retry}
	if v.data == "[DONE]" {
		return frameView{}, false, io.EOF
	}
	if len(f.Data) > 0 && f.Data[0] == '{' {
		if err := json.Unmarshal(f.Data, &v.fields); err != nil {
			return frameView{}, false, err
		}
	}
	return v, false, nil
}

func bodyConnect(body string) streaming.SSEConnector {
	return func(context.Context, string) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(body)), nil
	}
}

func collect[T any](t *testing.T, body string, decode func(streaming.SSEFrame) (T, bool, error)) ([]T, *streaming.SSEReader[T]) {
	t.Helper()
	connect := bodyConnect(body)
	s, err := streaming.NewSSEReader(context.Background(), connect, streaming.SSEHandler[T]{
		Decode: decode,
	}, nil)
	if err != nil {
		t.Fatalf("NewSSEReader: %v", err)
	}
	t.Cleanup(func() { require.NoError(t, s.Close()) })
	var got []T
	for {
		v, err := s.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got = append(got, v)
	}
	return got, s
}

func TestUnnamedEvents(t *testing.T) {
	body := "data: {\"desc\": \"one\"}\n\ndata: {\"desc\": \"two\"}\n\ndata: {\"desc\": \"three\"}\n\n"
	got, _ := collect(t, body, decodeView)
	want := []string{"one", "two", "three"}
	if len(got) != len(want) {
		t.Fatalf("got %d events, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i].fields["desc"] != w {
			t.Fatalf("event %d = %+v, want desc %q", i, got[i], w)
		}
	}
}

func TestNamedEventsWithTerminal(t *testing.T) {
	body := "event: responseCreated\ndata: {\"id\": \"resp_1\"}\n\n" +
		"event: responseDelta\ndata: {\"delta\": \"Hello\"}\n\n" +
		"event: responseDelta\ndata: {\"delta\": \" world\"}\n\n" +
		"data: [DONE]\n\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 3 {
		t.Fatalf("got %d events, want 3 (terminal [DONE] ends the stream)", len(got))
	}
	if got[0].eventType != "responseCreated" || got[0].fields["id"] != "resp_1" {
		t.Fatalf("event 0 = %+v", got[0])
	}
	if got[1].fields["delta"] != "Hello" {
		t.Fatalf("event 1 = %+v", got[1])
	}
	if got[2].fields["delta"] != " world" {
		t.Fatalf("event 2 = %+v", got[2])
	}
}

func TestProtocolEnvelopeMetadata(t *testing.T) {
	body := "id: event-1\nevent: message\ndata: {\"message\": \"hello\"}\n\n"
	got, s := collect(t, body, decodeView)
	if len(got) != 1 || got[0].fields["message"] != "hello" {
		t.Fatalf("got %+v", got)
	}
	if s.LastEventID() != "event-1" {
		t.Fatalf("LastEventID = %q, want event-1", s.LastEventID())
	}
}

func TestProtocolRetryValid(t *testing.T) {
	body := "retry: 1000\nevent: message\ndata: {\"message\": \"hello\"}\n\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 1 || got[0].retry != 1000 {
		t.Fatalf("retry = %d, want 1000", got[0].retry)
	}
}

func TestProtocolRetryInvalidIgnored(t *testing.T) {
	body := "retry: not-a-number\nevent: message\ndata: {\"message\": \"hello\"}\n\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 1 || got[0].retry != -1 {
		t.Fatalf("retry = %d, want -1 for invalid retry", got[0].retry)
	}
}

func TestProtocolInvalidIDIgnored(t *testing.T) {
	body := "id: invalid\x00id\nevent: message\ndata: {\"message\": \"hello\"}\n\n"
	_, s := collect(t, body, decodeView)
	if s.LastEventID() != "" {
		t.Fatalf("LastEventID = %q, want empty for id containing NUL", s.LastEventID())
	}
}

func TestIncompleteEventIsDiscarded(t *testing.T) {
	body := "event: withEnvelope\ndata: hello"
	got, _ := collect(t, body, decodeView)
	require.Empty(t, got)
}

func TestDataWithoutEnvelope(t *testing.T) {
	body := "event: withoutEnvelope\ndata: {\"metadata\": {\"source\": \"test\"}, \"contents\": \"world\"}\n\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 1 || got[0].fields["contents"] != "world" {
		t.Fatalf("got %+v", got)
	}
	meta, ok := got[0].fields["metadata"].(map[string]any)
	if !ok || meta["source"] != "test" {
		t.Fatalf("metadata = %+v", got[0].fields["metadata"])
	}
}

func TestCRLFLineEndings(t *testing.T) {
	body := "event: responseDelta\r\ndata: {\"delta\": \"crlf\"}\r\n\r\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 1 || got[0].fields["delta"] != "crlf" {
		t.Fatalf("got %+v", got)
	}
}

func TestMultilineData(t *testing.T) {
	// Two data lines concatenate with a newline; the JSON stays valid.
	body := "event: responseDelta\ndata: {\"delta\":\ndata: \"multi\"}\n\n"
	got, _ := collect(t, body, decodeView)
	if len(got) != 1 || got[0].fields["delta"] != "multi" {
		t.Fatalf("got %+v", got)
	}
}

func TestNewSSEReaderNilConnector(t *testing.T) {
	handler := streaming.SSEHandler[frameView]{Decode: decodeView}
	_, err := streaming.NewSSEReader(context.Background(), nil, handler, nil)
	require.ErrorContains(t, err, "Connect must not be nil")
}

func TestReconnectResumesWithLastEventID(t *testing.T) {
	segments := []string{
		"id: 1\ndata: {\"desc\": \"a\"}\n\nid: 2\ndata: {\"desc\": \"b\"}\n\n",
		"id: 3\ndata: {\"desc\": \"c\"}\n\n",
		"", // Connect returning io.EOF signals the continuous stream is exhausted
	}
	var lastIDs []string
	call := 0
	handler := streaming.SSEHandler[frameView]{
		Decode:    decodeView,
		Reconnect: true,
	}
	connect := func(_ context.Context, lastEventID string) (io.ReadCloser, error) {
		lastIDs = append(lastIDs, lastEventID)
		body := segments[call]
		call++
		if call == len(segments) {
			return nil, io.EOF
		}
		return io.NopCloser(strings.NewReader(body)), nil
	}
	s, err := streaming.NewSSEReader(context.Background(), connect, handler, &streaming.SSEReaderOptions{ReconnectDelay: time.Millisecond})
	if err != nil {
		t.Fatalf("NewSSEReader: %v", err)
	}
	var got []string
	for {
		v, err := s.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got = append(got, v.fields["desc"].(string))
	}
	if strings.Join(got, ",") != "a,b,c" {
		t.Fatalf("got %v, want a,b,c", got)
	}
	// initial connect ("") then reconnect after id 2, then after id 3.
	if strings.Join(lastIDs, ",") != ",2,3" {
		t.Fatalf("lastIDs = %v, want [\"\" 2 3]", lastIDs)
	}
}

func TestFiniteStreamWithoutEventIDDoesNotReplay(t *testing.T) {
	const body = "data: {\"desc\": \"a\"}\n\ndata: {\"desc\": \"b\"}\n\n"
	connects := 0
	connect := func(context.Context, string) (io.ReadCloser, error) {
		connects++
		return io.NopCloser(strings.NewReader(body)), nil
	}
	s, err := streaming.NewSSEReader(context.Background(), connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
	if err != nil {
		t.Fatalf("NewSSEReader: %v", err)
	}
	var got []string
	for {
		v, err := s.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got = append(got, v.fields["desc"].(string))
	}
	if strings.Join(got, ",") != "a,b" {
		t.Fatalf("got %v, want a,b", got)
	}
	if connects != 1 {
		t.Fatalf("connects = %d, want 1 (no replay)", connects)
	}
}
