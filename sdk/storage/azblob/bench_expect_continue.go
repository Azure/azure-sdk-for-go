//go:build ignore

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type readSeekCloser struct{ *bytes.Reader }

func (r *readSeekCloser) Close() error { return nil }

func blockID(n int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%06d", n)))
}

type stats struct {
	durations []time.Duration
}

func (s *stats) add(d time.Duration) { s.durations = append(s.durations, d) }
func (s *stats) sortDurations() {
	sort.Slice(s.durations, func(i, j int) bool { return s.durations[i] < s.durations[j] })
}
func (s *stats) p(pct int) time.Duration {
	if len(s.durations) == 0 {
		return 0
	}
	idx := pct * len(s.durations) / 100
	if idx >= len(s.durations) {
		idx = len(s.durations) - 1
	}
	return s.durations[idx]
}
func (s *stats) mean() time.Duration {
	if len(s.durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range s.durations {
		total += d
	}
	return total / time.Duration(len(s.durations))
}

// continueTracker is a per-retry policy that uses httptrace to measure the
// time spent waiting for the server's "100 Continue" response. This is the
// real overhead introduced by the Expect header.
type continueTracker struct {
	mu             sync.Mutex
	continueDelays []time.Duration
}

func (ct *continueTracker) Do(req *policy.Request) (*http.Response, error) {
	raw := req.Raw()
	if raw.Header.Get("Expect") == "" {
		return req.Next()
	}

	var headersSent time.Time
	var got100 time.Time

	trace := &httptrace.ClientTrace{
		WroteHeaders: func() {
			headersSent = time.Now()
		},
		Got100Continue: func() {
			got100 = time.Now()
		},
	}
	ctx := httptrace.WithClientTrace(raw.Context(), trace)
	*req.Raw() = *raw.WithContext(ctx)

	resp, err := req.Next()

	if !headersSent.IsZero() && !got100.IsZero() {
		delay := got100.Sub(headersSent)
		ct.mu.Lock()
		ct.continueDelays = append(ct.continueDelays, delay)
		ct.mu.Unlock()
	}

	return resp, err
}

func (ct *continueTracker) stats() *stats {
	st := &stats{durations: ct.continueDelays}
	st.sortDurations()
	return st
}

func runBench(client *blockblob.Client, body []byte, requests, concurrency int) *stats {
	st := &stats{durations: make([]time.Duration, 0, requests)}
	var mu sync.Mutex

	ch := make(chan int, requests)
	for i := 0; i < requests; i++ {
		ch <- i
	}
	close(ch)

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for w := 0; w < concurrency; w++ {
		go func() {
			defer wg.Done()
			for id := range ch {
				reader := &readSeekCloser{bytes.NewReader(body)}
				start := time.Now()
				_, err := client.StageBlock(context.Background(), blockID(id), reader, nil)
				elapsed := time.Since(start)
				if err != nil {
					fmt.Printf("    error (id=%d): %v\n", id, err)
					continue
				}
				mu.Lock()
				st.add(elapsed)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	st.sortDurations()
	return st
}

func main() {
	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if connStr == "" {
		fmt.Fprintln(os.Stderr, "ERROR: Set AZURE_STORAGE_CONNECTION_STRING")
		os.Exit(1)
	}

	requests := 50
	concurrency := 8
	if v := os.Getenv("BENCH_REQUESTS"); v != "" {
		requests, _ = strconv.Atoi(v)
	}
	if v := os.Getenv("BENCH_CONCURRENCY"); v != "" {
		concurrency, _ = strconv.Atoi(v)
	}

	containerName := "bench100continue"
	cntClient, _ := container.NewClientFromConnectionString(connStr, containerName, nil)
	cntClient.Create(context.Background(), nil)

	sizes := []int{8 * 1024 * 1024, 16 * 1024 * 1024, 32 * 1024 * 1024, 64 * 1024 * 1024} // 8 MiB, 16 MiB, 32 MiB, 64 MiB

	fmt.Printf("Benchmark: %d requests per mode, %d concurrency\n\n", requests, concurrency)
	fmt.Printf("%-10s %-20s %10s %10s %10s %10s %10s\n", "Size", "Mode", "Mean", "P50", "P90", "P99", "Max")
	fmt.Printf("%-10s %-20s %10s %10s %10s %10s %10s\n", "----", "----", "----", "---", "---", "---", "---")

	for _, size := range sizes {
		body := make([]byte, size)
		rand.Read(body)
		sizeMB := fmt.Sprintf("%d MiB", size/(1024*1024))

		modes := []struct {
			name    string
			disable string
		}{
			{"Without 100-cont", "1"},
			{"With 100-continue", ""},
		}

		var baseline time.Duration

		for i, mode := range modes {
			// Set env BEFORE creating the client so the policy picks it up
			os.Setenv("AZURE_STORAGE_DISABLE_EXPECT_CONTINUE", mode.disable)

			// Sleep 10s before starting test.
			time.Sleep(10 * time.Second)

			tracker := &continueTracker{}
			opts := &blockblob.ClientOptions{}
			opts.PerRetryPolicies = []policy.Policy{tracker}

			blobName := fmt.Sprintf("bench-%d-%d-%d", size, i, time.Now().UnixNano())
			client, err := blockblob.NewClientFromConnectionString(connStr, containerName, blobName, opts)
			if err != nil {
				panic(err)
			}

			// Warmup: 3 requests to establish connections
			for w := 0; w < 3; w++ {
				reader := &readSeekCloser{bytes.NewReader(body)}
				client.StageBlock(context.Background(), blockID(9000+w), reader, nil)
			}

			st := runBench(client, body, requests, concurrency)

			overhead := ""
			if i == 0 {
				baseline = st.mean()
			} else if baseline > 0 {
				pct := float64(st.mean()-baseline) / float64(baseline) * 100
				overhead = fmt.Sprintf("  (%+.1f%%)", pct)
			}

			fmt.Printf("%-10s %-20s %10s %10s %10s %10s %10s%s\n",
				sizeMB, mode.name,
				st.mean().Round(time.Millisecond),
				st.p(50).Round(time.Millisecond),
				st.p(90).Round(time.Millisecond),
				st.p(99).Round(time.Millisecond),
				st.p(100).Round(time.Millisecond),
				overhead,
			)

			// Print 100-continue round-trip overhead if we captured any
			cs := tracker.stats()
			if len(cs.durations) > 0 {
				fmt.Printf("%-10s   └─ 100-cont wait   %10s %10s %10s %10s %10s  (%d samples)\n",
					"",
					cs.mean().Round(time.Millisecond),
					cs.p(50).Round(time.Millisecond),
					cs.p(90).Round(time.Millisecond),
					cs.p(99).Round(time.Millisecond),
					cs.p(100).Round(time.Millisecond),
					len(cs.durations),
				)
			}
		}
		fmt.Println()
	}

	os.Unsetenv("AZURE_STORAGE_DISABLE_EXPECT_CONTINUE")
	fmt.Println("Done!")
}
