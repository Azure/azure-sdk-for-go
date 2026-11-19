// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

func TestSSEBodyLoggingPreservesIncrementalDelivery(t *testing.T) {
	log.SetListener(func(log.Event, string) {})
	t.Cleanup(func() { log.SetListener(nil) })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprint(w, "data: hello\n\n")
		w.(http.Flusher).Flush()
		<-req.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Logging: policy.LogOptions{IncludeBody: true},
		Retry:   policy.RetryOptions{MaxRetries: -1},
	})
	reader, err := streaming.NewSSEReader(ctx, func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
		return sendSSERequest(ctx, pipeline, server.URL, lastEventID)
	}, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	require.NoError(t, ctx.Err(), "construction must complete before the server ends its response")
	event, err := reader.Next()
	require.NoError(t, err)
	require.Equal(t, "hello", event.data)
	require.NoError(t, ctx.Err(), "the first event must be delivered while the response remains open")
	require.NoError(t, reader.Close())
}
