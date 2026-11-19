// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedTerminalPayloadIsDelivered(t *testing.T) {
	type completion struct {
		Total int `json:"total"`
	}
	reader, err := streaming.NewSSEReader(context.Background(), bodyConnect("event: complete\ndata: {\"total\":7}\n\n"),
		streaming.SSEHandler[completion]{
			Decode: func(frame streaming.SSEFrame) (completion, bool, error) {
				var value completion
				err := json.Unmarshal(frame.Data, &value)
				return value, frame.Type == "complete", err
			},
		}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, reader.Close()) })
	event, err := reader.Next()
	require.NoError(t, err)
	assert.Equal(t, 7, event.Total)
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
}

func TestTerminalPayloadStopsIterationAndReconnect(t *testing.T) {
	for _, eventType := range []string{"", "complete"} {
		t.Run(eventType, func(t *testing.T) {
			connects := 0
			decodes := 0
			reader, err := streaming.NewSSEReader(context.Background(), func(context.Context, string) (io.ReadCloser, error) {
				connects++
				return eventBody("event: " + eventType + "\nid: final\ndata: payload\n\ndata: after-terminal\n\n"), nil
			}, streaming.SSEHandler[string]{
				Decode: func(frame streaming.SSEFrame) (string, bool, error) {
					decodes++
					return string(frame.Data), true, nil
				},
				Reconnect: true,
			}, nil)
			require.NoError(t, err)
			t.Cleanup(func() { assert.NoError(t, reader.Close()) })
			value, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, "payload", value)
			for range 2 {
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
			}
			assert.Equal(t, "final", reader.LastEventID())
			assert.Equal(t, 1, connects)
			assert.Equal(t, 1, decodes)
		})
	}
}

func TestSSELineEndingsAcrossReadBoundaries(t *testing.T) {
	for _, ending := range []string{"\r", "\n", "\r\n"} {
		wire := "id: first" + ending + "data: hello" + ending + "data: continued" + ending + ending +
			"data: world" + ending + ending
		for boundary := 1; boundary < len(wire); boundary++ {
			read := io.MultiReader(strings.NewReader(wire[:boundary]), strings.NewReader(wire[boundary:]))
			reader := newSSEReader(t, io.NopCloser(read), streaming.SSEHandler[frameView]{}, nil)
			for _, want := range []string{"hello\ncontinued", "world"} {
				event, err := reader.Next()
				require.NoError(t, err, "ending %q, boundary %d", ending, boundary)
				assert.Equal(t, want, event.data)
				assert.Equal(t, "first", reader.LastEventID())
			}
			_, err := reader.Next()
			require.ErrorIs(t, err, io.EOF)
		}
	}
}

func TestSplitCRLFDoesNotDispatchPrematurely(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		read, write := io.Pipe()
		defer func() {
			_ = write.Close()
		}()
		reader := newSSEReader(t, read, streaming.SSEHandler[frameView]{}, nil)
		defer func() {
			_ = reader.Close()
		}()
		type result struct {
			event frameView
			err   error
		}
		next := make(chan result, 1)
		go func() {
			event, err := reader.Next()
			next <- result{event: event, err: err}
		}()
		for _, chunk := range []string{"data: hello\r", "\n", "data: world\r", "\n"} {
			_, err := io.WriteString(write, chunk)
			require.NoError(t, err)
			synctest.Wait()
			select {
			case got := <-next:
				t.Fatalf("event dispatched before the blank line: %+v", got)
			default:
			}
		}
		_, err := io.WriteString(write, "\r")
		require.NoError(t, err)
		synctest.Wait()
		select {
		case got := <-next:
			require.NoError(t, got.err)
			assert.Equal(t, "hello\nworld", got.event.data)
		default:
			assert.Fail(t, "a complete event requires another read")
			require.NoError(t, reader.Close())
			<-next
		}
	})
}

func TestLoneCRDispatchesWithoutAnotherRead(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		read, write := io.Pipe()
		defer func() {
			_ = write.Close()
		}()
		reader := newSSEReader(t, read, streaming.SSEHandler[frameView]{}, nil)
		defer func() {
			_ = reader.Close()
		}()
		type result struct {
			event frameView
			err   error
		}
		next := make(chan result, 1)
		go func() {
			event, err := reader.Next()
			next <- result{event: event, err: err}
		}()
		_, err := io.WriteString(write, "data: hello\r\r")
		require.NoError(t, err)
		synctest.Wait()
		select {
		case got := <-next:
			require.NoError(t, got.err)
			assert.Equal(t, "hello", got.event.data)
		default:
			assert.Fail(t, "a complete CR-delimited event is blocked until another byte arrives")
			require.NoError(t, reader.Close())
			<-next
		}
	})
}
