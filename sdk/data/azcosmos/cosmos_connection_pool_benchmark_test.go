// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// BenchmarkCosmosConnectionPool compares a fixed local workload, not Cosmos
// service performance. GetConn-to-GotConn includes both queueing and connection
// establishment. All latencies include reading the response body.
func BenchmarkCosmosConnectionPool(b *testing.B) {
	for _, protocol := range []struct {
		name    string
		h2      bool
		streams int
	}{
		{"H1", false, 1},
		{"H2_streams1", true, 1},
		{"H2_streams100", true, 100},
	} {
		for _, limit := range []int{0, 32} {
			for _, concurrency := range []int{64, 256} {
				b.Run(fmt.Sprintf("%s/limit%d/concurrency%d", protocol.name, limit, concurrency), func(b *testing.B) {
					payload := strings.Repeat("x", 1024)
					server := newPoolTestServer(b, protocol.h2, protocol.streams, func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path != "/warm" {
							time.Sleep(2 * time.Millisecond)
						}
						_, _ = io.WriteString(w, payload)
					})
					client, counters := newPoolTestClient(b, server, limit)
					warm, err := client.Get(server.URL + "/warm")
					require.NoError(b, err)
					drainPoolResponse(b, warm)
					latency := make([]time.Duration, b.N)
					queue := make([]time.Duration, b.N)
					var next atomic.Int64
					var failures atomic.Int64
					firstError := make(chan error, 1)
					var wg sync.WaitGroup
					b.ReportAllocs()
					b.ResetTimer()
					start := time.Now()
					for range concurrency {
						wg.Go(func() {
							for {
								index := int(next.Add(1) - 1)
								if index >= b.N {
									return
								}
								var getConn time.Time
								trace := &httptrace.ClientTrace{
									GetConn: func(string) { getConn = time.Now() },
									GotConn: func(httptrace.GotConnInfo) { queue[index] = time.Since(getConn) },
								}
								ctx := httptrace.WithClientTrace(b.Context(), trace)
								req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
								if err != nil {
									failures.Add(1)
									select {
									case firstError <- err:
									default:
									}
									continue
								}
								requestStart := time.Now()
								response, err := client.Do(req)
								if err == nil {
									_, err = io.Copy(io.Discard, response.Body)
									if closeErr := response.Body.Close(); err == nil {
										err = closeErr
									}
								}
								if err != nil {
									failures.Add(1)
									select {
									case firstError <- err:
									default:
									}
									continue
								}
								latency[index] = time.Since(requestStart)
							}
						})
					}
					wg.Wait()
					elapsed := time.Since(start)
					b.StopTimer()
					if failures.Load() != 0 {
						b.Logf("%d requests failed; first error: %v", failures.Load(), <-firstError)
					}
					if limit > 0 {
						require.Zero(b, failures.Load())
					}
					successLatency := latency[:0]
					successQueue := queue[:0]
					for i, duration := range latency {
						if duration != 0 {
							successLatency = append(successLatency, duration)
							successQueue = append(successQueue, queue[i])
						}
					}
					latency, queue = successLatency, successQueue
					require.NotEmpty(b, latency, "no successful requests to measure")
					slices.Sort(latency)
					slices.Sort(queue)
					b.ReportMetric(float64(failures.Load())/float64(b.N), "errors/op")
					b.ReportMetric(float64(len(latency))/elapsed.Seconds(), "success/s")
					b.ReportMetric(float64(latency[(len(latency)-1)*50/100].Microseconds()), "p50-us")
					b.ReportMetric(float64(latency[(len(latency)-1)*95/100].Microseconds()), "p95-us")
					b.ReportMetric(float64(latency[(len(latency)-1)*99/100].Microseconds()), "p99-us")
					b.ReportMetric(float64(queue[(len(queue)-1)*99/100].Microseconds()), "queue-p99-us")
					b.ReportMetric(float64(counters.peak.Load()), "peak-conns")
					b.ReportMetric(float64(counters.live.Load()), "retained-conns")
					b.ReportMetric(float64(counters.dials.Load()), "dials")
				})
			}
		}
	}
}
