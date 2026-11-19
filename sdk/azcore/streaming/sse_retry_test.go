// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"io"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

func TestReconnectDelayPersistsAcrossSegments(t *testing.T) {
	for _, tc := range []struct {
		name    string
		first   string
		second  string
		options *streaming.SSEReaderOptions
		want    []time.Duration
	}{
		{name: "default", want: []time.Duration{3 * time.Second, 3 * time.Second}},
		{name: "configured", options: &streaming.SSEReaderOptions{ReconnectDelay: 2 * time.Second}, want: []time.Duration{2 * time.Second, 2 * time.Second}},
		{name: "server zero", first: "retry: 0\n", want: []time.Duration{0, 0}},
		{name: "server delay", first: "retry: 1500\n", want: []time.Duration{1500 * time.Millisecond, 1500 * time.Millisecond}},
		{name: "invalid retry", first: "retry: 1500\n", second: "retry: invalid\n", want: []time.Duration{1500 * time.Millisecond, 1500 * time.Millisecond}},
		{name: "negative retry", first: "retry: 1500\n", second: "retry: -1\n", want: []time.Duration{1500 * time.Millisecond, 1500 * time.Millisecond}},
		{name: "server override", first: "retry: 1500\n", second: "retry: 250\n", want: []time.Duration{1500 * time.Millisecond, 250 * time.Millisecond}},
		{name: "server override to zero", first: "retry: 1500\n", second: "retry: 0\n", want: []time.Duration{1500 * time.Millisecond, 0}},
		{name: "override configured delay", first: "retry: 0\n", options: &streaming.SSEReaderOptions{ReconnectDelay: 2 * time.Second}, want: []time.Duration{0, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var delays []time.Duration
				previous := time.Now()
				connect := initialBodyConnector(eventBody(tc.first+"data: one\n\n"),
					func(context.Context, string) (io.ReadCloser, error) {
						delays = append(delays, time.Since(previous))
						if len(delays) == 1 {
							return eventBody(tc.second + "data: two\n\n"), nil
						}
						return nil, io.EOF
					})
				reader := newSSEReaderWithContext(t, context.Background(), connect, streaming.SSEHandler[frameView]{Reconnect: true}, tc.options)
				for _, data := range []string{"one", "two"} {
					event, err := reader.Next()
					require.NoError(t, err)
					require.Equal(t, data, event.data)
					previous = time.Now()
				}
				_, err := reader.Next()
				require.ErrorIs(t, err, io.EOF)
				require.Equal(t, tc.want, delays)
			})
		})
	}
}
