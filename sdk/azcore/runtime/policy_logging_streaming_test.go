// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestResponseBodyLoggingHonorsSkipBodyDownload(t *testing.T) {
	for _, tc := range []struct {
		name        string
		skip        bool
		status      int
		contentType string
		logBody     bool
	}{
		{"streaming SSE", true, 200, "text/event-stream", false},
		{"streaming JSONL", true, 200, "application/jsonl", false},
		{"streaming text", true, 200, "text/plain", false},
		{"redirect response", true, 302, "text/plain", false},
		{"error response", true, 400, "application/json", true},
		{"server error response", true, 500, "application/json", true},
		{"normal response", false, 200, "application/json", true},
		{"non-text error response", true, 400, "application/octet-stream", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var responseLog string
			var requestLog string
			log.SetListener(func(event log.Event, message string) {
				switch event {
				case log.EventResponse:
					responseLog = message
				case log.EventRequest:
					requestLog = message
				}
			})
			t.Cleanup(func() { log.SetListener(nil) })
			server, closeServer := mock.NewServer()
			defer closeServer()
			server.SetResponse(mock.WithStatusCode(tc.status), mock.WithHeader("Content-Type", tc.contentType), mock.WithBody([]byte("response-payload")))
			pipeline := exported.NewPipeline(nil, server, NewLogPolicy(&policy.LogOptions{IncludeBody: true}))
			request, err := NewRequest(context.Background(), http.MethodPost, server.URL())
			require.NoError(t, err)
			require.NoError(t, request.SetBody(exported.NopCloser(strings.NewReader("request-payload")), "text/plain"))
			if tc.skip {
				SkipBodyDownload(request)
			}
			response, err := pipeline.Do(request)
			require.NoError(t, err)
			defer func() {
				_ = response.Body.Close()
			}()
			require.Contains(t, requestLog, "request-payload", "response streaming must not suppress request body logging")
			require.Contains(t, responseLog, "RESPONSE Status:")
			require.Contains(t, responseLog, tc.contentType)
			if tc.logBody {
				require.Contains(t, responseLog, "response-payload")
			} else {
				require.NotContains(t, responseLog, "response-payload")
			}
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			require.Equal(t, "response-payload", string(body))
		})
	}
}
