// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package streaming_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

func TestSSEHTTPContinuousReaderReconnects(t *testing.T) {
	for _, unexpectedEOF := range []bool{false, true} {
		name := "EOF"
		if unexpectedEOF {
			name = "unexpected EOF"
		}
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			headers := make(chan string, 3)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				request := requests.Add(1)
				if request <= 3 {
					headers <- r.Header.Get("Last-Event-ID")
				}
				switch request {
				case 1:
					w.Header().Set("Content-Type", "text/event-stream")
					if unexpectedEOF {
						w.Header().Set("Content-Length", "1000")
					}
					_, _ = fmt.Fprint(w, "retry: 0\nid: event-1\nevent: message\ndata: {\"message\":\"hello\"}\n\n")
				case 2:
					w.Header().Set("Content-Type", "text/event-stream")
					_, _ = fmt.Fprint(w, "id: event-2\nevent: message\ndata: {\"message\":\"world\"}\n\n")
				case 3:
					w.WriteHeader(http.StatusNoContent)
				default:
					http.Error(w, "unexpected reconnect", http.StatusBadRequest)
				}
			}))
			t.Cleanup(server.Close)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)
			var readerCtx, attemptCtx context.Context
			request := sseHTTPConnector(server.URL, &policy.ClientOptions{
				Retry: policy.RetryOptions{TryTimeout: time.Minute},
				Transport: sseTransport(func(req *http.Request) (*http.Response, error) {
					attemptCtx = req.Context()
					return http.DefaultClient.Do(req)
				}),
			})
			connects := 0
			connect := func(connectCtx context.Context, lastEventID string) (io.ReadCloser, error) {
				connects++
				require.NoError(t, connectCtx.Err())
				require.NoError(t, ctx.Err())
				if attemptCtx != nil {
					require.ErrorIs(t, attemptCtx.Err(), context.Canceled)
					require.Same(t, readerCtx, connectCtx)
				} else {
					readerCtx = connectCtx
				}
				return request(connectCtx, lastEventID)
			}
			reader, err := streaming.NewSSEReader(ctx, connect, streaming.SSEHandler[frameView]{
				Decode:    decodeView,
				Reconnect: true,
			}, &streaming.SSEReaderOptions{LastEventID: "saved-event"})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, reader.Close()) })
			require.Equal(t, 1, connects, "the initial connection must be eager")
			for _, message := range []string{"hello", "world"} {
				event, err := reader.Next()
				require.NoError(t, err)
				require.Equal(t, message, event.fields["message"])
			}
			for range 2 {
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
			}
			require.Equal(t, "event-2", reader.LastEventID())
			require.Equal(t, 3, connects)
			require.EqualValues(t, 3, requests.Load())
			for _, wantID := range []string{"saved-event", "event-1", "event-2"} {
				require.Equal(t, wantID, <-headers)
			}
			require.ErrorIs(t, readerCtx.Err(), context.Canceled)
			require.ErrorIs(t, attemptCtx.Err(), context.Canceled)
			require.NoError(t, ctx.Err())
		})
	}
}
