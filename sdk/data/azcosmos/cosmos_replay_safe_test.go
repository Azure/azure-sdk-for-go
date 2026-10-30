// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newReplaySafeTestContainer(t *testing.T, transport policy.Transporter, retry policy.RetryOptions, multiWrite, crossRegion bool) (*ContainerClient, *locationCache) {
	t.Helper()
	client, _ := newReadAttemptTimeoutTestClient(t, "https://write-a.test", transport, retry)
	preferred := []string{"North Europe", "East US", "West US"}
	lc := newLocationCache(preferred, *client.endpointUrl, crossRegion)
	writes := []accountRegion{{Name: "East US", Endpoint: "https://write-a.test"}}
	if multiWrite {
		writes = append(writes, accountRegion{Name: "West US", Endpoint: "https://write-b.test"})
	}
	reads := append([]accountRegion{{Name: "North Europe", Endpoint: "https://read-only.test"}}, writes...)
	require.NoError(t, lc.update(writes, reads, preferred, &multiWrite))
	client.gem.locationCache = lc
	client.gem.preferredLocations = preferred
	container, err := client.NewContainer("db", "container")
	require.NoError(t, err)
	return container, lc
}

func replaySafeTestFailure(req *http.Request, failure string) (*http.Response, error) {
	switch failure {
	case "not sent":
		return nil, &net.DNSError{Err: "unreachable"}
	case "unclassified":
		return nil, errors.New("unclassified transport failure")
	case "EOF":
		return nil, io.EOF
	case "attempt timeout":
		<-req.Context().Done()
		return nil, req.Context().Err()
	case "408":
		return readAttemptTimeoutResponse(req, http.StatusRequestTimeout), nil
	case "418":
		return readAttemptTimeoutResponse(req, http.StatusTeapot), nil
	case "500":
		return readAttemptTimeoutResponse(req, http.StatusInternalServerError), nil
	case "502":
		return readAttemptTimeoutResponse(req, http.StatusBadGateway), nil
	case "503":
		return readAttemptTimeoutResponse(req, http.StatusServiceUnavailable), nil
	case "504":
		return readAttemptTimeoutResponse(req, http.StatusGatewayTimeout), nil
	default:
		panic("unknown test failure: " + failure)
	}
}

func TestReplaySafeItemWrites(t *testing.T) {
	for _, operation := range []string{"create", "upsert", "replace", "delete", "patch"} {
		for _, failure := range []string{"EOF", "attempt timeout", "408", "500", "502", "504"} {
			for _, optIn := range []string{"nil", "false", "true"} {
				t.Run(operation+"/"+failure+"/"+optIn, func(t *testing.T) {
					t.Parallel()
					var hosts, bodies []string
					transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
						hosts = append(hosts, req.URL.Host)
						assert.Equal(t, req.URL.Host, req.Host)
						var body []byte
						if req.Body != nil {
							var err error
							body, err = io.ReadAll(req.Body)
							require.NoError(t, err)
						}
						bodies = append(bodies, string(body))
						if len(hosts) == 1 {
							return replaySafeTestFailure(req, failure)
						}
						return readAttemptTimeoutResponse(req, http.StatusOK), nil
					})
					retry := policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Millisecond}
					if failure == "attempt timeout" {
						retry.TryTimeout = 10 * time.Millisecond
					}
					container, lc := newReplaySafeTestContainer(t, transport, retry, false, true)
					var options *ItemOptions
					if optIn != "nil" {
						options = &ItemOptions{ReplaySafe: optIn == "true"}
					}
					ctx, pk := context.Background(), NewPartitionKeyString("pk")
					item := []byte(`{"id":"item","value":1}`)
					var err error
					switch operation {
					case "create":
						_, err = container.CreateItem(ctx, pk, item, options)
					case "upsert":
						_, err = container.UpsertItem(ctx, pk, item, options)
					case "replace":
						_, err = container.ReplaceItem(ctx, pk, "item", item, options)
					case "delete":
						_, err = container.DeleteItem(ctx, pk, "item", options)
					case "patch":
						patch := PatchOperations{}
						patch.AppendSet("/value", 1)
						_, err = container.PatchItem(ctx, pk, "item", patch, options)
					}
					if optIn == "true" {
						require.NoError(t, err)
						assert.Equal(t, []string{"write-a.test", "write-a.test"}, hosts)
						require.Len(t, bodies, 2)
						assert.Equal(t, bodies[0], bodies[1], "replays must rewind the original body")
					} else {
						require.Error(t, err)
						assert.Equal(t, []string{"write-a.test"}, hosts)
					}
					assert.Empty(t, lc.locationUnavailabilityInfoMap, "ambiguous failures must not demote endpoints")
				})
			}
		}
	}
}

func TestReplaySafeRetryParity(t *testing.T) {
	for _, tt := range []struct {
		name        string
		failure     string
		crossRegion bool
		multiWrite  bool
		noPreferred bool
		// Counts are ordered read, replay-safe write, default write.
		sends     [3]int
		callbacks [3]int
		delegated bool
	}{
		{name: "network disabled", failure: "EOF", sends: [3]int{1, 1, 1}},
		{name: "network exhaustion", failure: "EOF", crossRegion: true, multiWrite: true, sends: [3]int{5, 5, 1}},
		{name: "network single write region", failure: "EOF", crossRegion: true, sends: [3]int{5, 4, 1}},
		{name: "not sent disabled", failure: "not sent", sends: [3]int{1, 1, 1}},
		{name: "not sent exhaustion", failure: "not sent", crossRegion: true, multiWrite: true, sends: [3]int{5, 5, 5}},
		{name: "not sent single write region", failure: "not sent", crossRegion: true, sends: [3]int{5, 4, 4}},
		{name: "408 disabled", failure: "408", sends: [3]int{1, 1, 1}},
		{name: "408 exhaustion", failure: "408", crossRegion: true, multiWrite: true, sends: [3]int{2, 2, 1}},
		{name: "408 single write region", failure: "408", crossRegion: true, sends: [3]int{2, 2, 1}},
		{name: "500 in region only", failure: "500", sends: [3]int{2, 2, 1}},
		{name: "500 exhaustion", failure: "500", crossRegion: true, multiWrite: true, sends: [3]int{3, 3, 1}},
		{name: "500 single write region", failure: "500", crossRegion: true, sends: [3]int{3, 2, 1}},
		{name: "503 disabled", failure: "503", sends: [3]int{1, 1, 1}},
		{name: "503 no preferred regions", failure: "503", crossRegion: true, multiWrite: true, noPreferred: true, sends: [3]int{1, 1, 1}},
		{name: "503 exhaustion", failure: "503", crossRegion: true, multiWrite: true, sends: [3]int{4, 4, 4}},
		{name: "503 single write region", failure: "503", crossRegion: true, sends: [3]int{4, 1, 1}},
		{name: "unclassified delegated", failure: "unclassified", sends: [3]int{3, 3, 3}, callbacks: [3]int{3, 3, 3}, delegated: true},
		{name: "attempt timeout delegated", failure: "attempt timeout", sends: [3]int{3, 3, 1}, callbacks: [3]int{3, 3, 0}, delegated: true},
	} {
		for i, operation := range []string{"read", "replay-safe write", "default write"} {
			t.Run(tt.name+"/"+operation, func(t *testing.T) {
				t.Parallel()
				var hosts []string
				callbacks := 0
				item := []byte(`{"id":"item"}`)
				transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
					hosts = append(hosts, req.URL.Host)
					assert.Equal(t, req.URL.Host, req.Host)
					if operation != "read" {
						assert.Contains(t, []string{"write-a.test", "write-b.test"}, req.URL.Host)
						if !tt.multiWrite {
							assert.Equal(t, "write-a.test", req.URL.Host)
						}
						body, err := io.ReadAll(req.Body)
						require.NoError(t, err)
						assert.Equal(t, item, body, "every native or outer attempt must rewind the write")
					}
					return replaySafeTestFailure(req, tt.failure)
				})
				retry := policy.RetryOptions{
					MaxRetries: 2, RetryDelay: time.Millisecond,
					ShouldRetry: func(*http.Response, error) bool {
						callbacks++
						return true
					},
				}
				if tt.failure == "attempt timeout" {
					retry.TryTimeout = 10 * time.Millisecond
				}
				container, lc := newReplaySafeTestContainer(t, transport, retry, tt.multiWrite, tt.crossRegion)
				if tt.noPreferred {
					container.database.client.gem.preferredLocations = nil
				}
				var err error
				if operation == "read" {
					_, err = container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item", nil)
				} else {
					_, err = container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), item, &ItemOptions{ReplaySafe: operation == "replay-safe write"})
				}
				require.Error(t, err)
				assert.Len(t, hosts, tt.sends[i], "replay safety cannot increase the read retry budget")
				if operation != "read" && tt.multiWrite && tt.crossRegion && !tt.noPreferred && len(hosts) > 1 {
					assert.Contains(t, hosts, "write-b.test", "eligible native failover must reach the other writable region")
				}
				assert.Equal(t, tt.callbacks[i], callbacks, "shared native refusal must not reach outer ShouldRetry")
				wantVeto := !tt.delegated || (tt.failure == "attempt timeout" && operation == "default write")
				assert.Equal(t, wantVeto, errors.As(err, new(errorinfo.NonRetriable)))
				if tt.failure != "not sent" {
					assert.Empty(t, lc.locationUnavailabilityInfoMap, "ambiguous failures must not mark endpoints unavailable")
				}
			})
		}
	}
}

func TestReplaySafeRetrySettings(t *testing.T) {
	for _, tt := range []struct {
		name        string
		failure     string
		crossRegion bool
		maxRetries  int32
		statusCodes []int
		reject      bool
		wantSends   int
	}{
		{name: "network cross-region disabled", failure: "EOF", maxRetries: 2, wantSends: 1},
		{name: "not sent cross-region disabled", failure: "not sent", maxRetries: 2, wantSends: 1},
		{name: "not sent retries disabled", failure: "not sent", maxRetries: -1, wantSends: 1},
		{name: "not sent ShouldRetry veto", failure: "not sent", maxRetries: 2, reject: true, wantSends: 1},
		{name: "503 cross-region disabled", failure: "503", maxRetries: 2, wantSends: 1},
		{name: "503 retries disabled", failure: "503", maxRetries: -1, wantSends: 1},
		{name: "503 excluded status", failure: "503", maxRetries: 2, statusCodes: []int{502}, wantSends: 1},
		{name: "503 included status", failure: "503", maxRetries: 2, statusCodes: []int{503}, wantSends: 1},
		{name: "503 ShouldRetry veto", failure: "503", maxRetries: 2, reject: true, wantSends: 1},
		{name: "408 cross-region disabled", failure: "408", maxRetries: 2, wantSends: 1},
		{name: "500 uses native budget", failure: "500", maxRetries: 2, wantSends: 2},
		{name: "network native budget with outer retries disabled", failure: "EOF", crossRegion: true, maxRetries: -1, wantSends: 5},
		{name: "HTTP native budget with outer retries disabled", failure: "500", crossRegion: true, maxRetries: -1, wantSends: 3},
		{name: "empty status list does not alter native 408", failure: "408", crossRegion: true, maxRetries: 2, statusCodes: []int{}, wantSends: 2},
		{name: "excluded status does not alter native 500", failure: "500", crossRegion: true, maxRetries: 2, statusCodes: []int{502}, wantSends: 3},
		{name: "included status", failure: "502", crossRegion: true, maxRetries: 2, statusCodes: []int{502}, wantSends: 3},
		{name: "ShouldRetry cannot extend native exhaustion", failure: "EOF", crossRegion: true, maxRetries: 2, reject: true, wantSends: 5},
		{name: "unclassified uses outer budget", failure: "unclassified", maxRetries: 2, wantSends: 3},
		{name: "unclassified retries disabled", failure: "unclassified", maxRetries: -1, wantSends: 1},
		{name: "unclassified ShouldRetry veto", failure: "unclassified", maxRetries: 2, reject: true, wantSends: 1},
		{name: "terminal status default", failure: "418", maxRetries: 2, wantSends: 1},
		{name: "outer empty status list", failure: "418", maxRetries: 2, statusCodes: []int{}, wantSends: 1},
		{name: "outer excluded status", failure: "418", maxRetries: 2, statusCodes: []int{502}, wantSends: 1},
		{name: "outer included status", failure: "418", maxRetries: 2, statusCodes: []int{418}, wantSends: 3},
		{name: "attempt retries disabled", failure: "attempt timeout", crossRegion: true, maxRetries: -1, wantSends: 1},
		{name: "attempt retries bounded", failure: "attempt timeout", crossRegion: true, maxRetries: 1, wantSends: 2},
	} {
		for _, write := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/write=%t", tt.name, write), func(t *testing.T) {
				t.Parallel()
				sends := 0
				transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
					sends++
					return replaySafeTestFailure(req, tt.failure)
				})
				retry := policy.RetryOptions{
					MaxRetries: tt.maxRetries, RetryDelay: time.Millisecond, StatusCodes: tt.statusCodes,
				}
				if tt.failure == "attempt timeout" {
					retry.TryTimeout = 20 * time.Millisecond
				}
				if tt.reject {
					retry.ShouldRetry = func(*http.Response, error) bool { return false }
				}
				container, _ := newReplaySafeTestContainer(t, transport, retry, true, tt.crossRegion)
				var err error
				if write {
					_, err = container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
				} else {
					_, err = container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item", nil)
				}
				require.Error(t, err)
				assert.Equal(t, tt.wantSends, sends)
			})
		}
	}
}

func TestReplaySafeWriteOuterDelay(t *testing.T) {
	for _, failure := range []string{"unclassified", "attempt timeout"} {
		t.Run(failure, func(t *testing.T) {
			t.Parallel()
			var hosts []string
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				hosts = append(hosts, req.URL.Host)
				if len(hosts) == 1 {
					return replaySafeTestFailure(req, failure)
				}
				return readAttemptTimeoutResponse(req, http.StatusOK), nil
			})
			container, lc := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
				MaxRetries: 1, RetryDelay: 40 * time.Millisecond, TryTimeout: 10 * time.Millisecond,
			}, false, false)
			start := time.Now()
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			require.NoError(t, err)
			assert.GreaterOrEqual(t, time.Since(start), 30*time.Millisecond, "outer backoff must be used, outside TryTimeout")
			assert.Equal(t, []string{"write-a.test", "write-a.test"}, hosts)
			assert.Empty(t, lc.locationUnavailabilityInfoMap)
		})
	}
}

func TestReplaySafeWriteHTTPResponseHandoff(t *testing.T) {
	for _, tt := range []struct {
		write  bool
		capped bool
	}{
		{},
		{capped: true},
		{write: true},
		{write: true, capped: true},
	} {
		t.Run(fmt.Sprintf("write=%t/RetryAfterCapped=%t", tt.write, tt.capped), func(t *testing.T) {
			t.Parallel()
			var first *http.Response
			sends, callbacks := 0, 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				if sends == 1 {
					first = readAttemptTimeoutResponse(req, http.StatusTeapot)
					first.Header.Set("Retry-After", "1")
					first.Header.Set("x-ms-activity-id", "original-response")
					return first, nil
				}
				return readAttemptTimeoutResponse(req, http.StatusOK), nil
			})
			retry := policy.RetryOptions{
				MaxRetries: 1, RetryDelay: time.Millisecond, MaxRetryDelay: 2 * time.Second,
				ShouldRetry: func(resp *http.Response, err error) bool {
					callbacks++
					require.NoError(t, err)
					if callbacks == 1 {
						assert.Same(t, first, resp)
						assert.Equal(t, "original-response", resp.Header.Get("x-ms-activity-id"))
					}
					return resp.StatusCode == http.StatusTeapot
				},
			}
			if tt.capped {
				retry.MaxRetryDelay = 10 * time.Millisecond
			}
			container, _ := newReplaySafeTestContainer(t, transport, retry, false, false)
			start := time.Now()
			var err error
			if tt.write {
				_, err = container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			} else {
				_, err = container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item", nil)
			}
			if tt.capped {
				require.Error(t, err)
				assert.Equal(t, 1, sends, "Retry-After above the configured cap must not retry")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 2, sends)
				assert.GreaterOrEqual(t, time.Since(start), time.Second, "Retry-After must control the outer delay")
			}
		})
	}
}

func TestReplaySafeWritePerCallRetryOverride(t *testing.T) {
	sends := 0
	transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
		sends++
		return replaySafeTestFailure(req, "unclassified")
	})
	container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{MaxRetries: -1}, false, false)
	ctx := policy.WithRetryOptions(context.Background(), policy.RetryOptions{MaxRetries: 2, RetryDelay: time.Millisecond})
	_, err := container.UpsertItem(ctx, NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
	assert.ErrorContains(t, err, "unclassified transport failure")
	assert.Equal(t, 3, sends)
}

func TestReplaySafeWriteExistingRecovery(t *testing.T) {
	for _, optIn := range []bool{false, true} {
		for _, failure := range []string{"not sent", "session unavailable", "503"} {
			t.Run(fmt.Sprintf("%s/ReplaySafe=%t", failure, optIn), func(t *testing.T) {
				t.Parallel()
				var hosts []string
				transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
					hosts = append(hosts, req.URL.Host)
					if failure == "not sent" && len(hosts) <= maxSameRegionConnectionRetries+1 {
						return nil, &net.DNSError{Err: "unreachable"}
					}
					if len(hosts) == 1 {
						status := http.StatusServiceUnavailable
						if failure == "session unavailable" {
							status = http.StatusNotFound
						}
						resp := readAttemptTimeoutResponse(req, status)
						if failure == "session unavailable" {
							resp.Header.Set(cosmosHeaderSubstatus, subStatusReadSessionNotAvailable)
						}
						return resp, nil
					}
					return readAttemptTimeoutResponse(req, http.StatusOK), nil
				})
				container, lc := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
					MaxRetries: -1,
					ShouldRetry: func(*http.Response, error) bool {
						// Only the final success should leave the internal Cosmos recovery.
						return false
					},
				}, true, true)
				_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: optIn})
				require.NoError(t, err)
				want := []string{"write-a.test", "write-b.test"}
				if failure == "not sent" {
					want = []string{"write-a.test", "write-a.test", "write-a.test", "write-a.test", "write-b.test"}
					assert.NotEmpty(t, lc.locationUnavailabilityInfoMap)
				} else {
					assert.Empty(t, lc.locationUnavailabilityInfoMap)
				}
				assert.Equal(t, want, hosts, "ReplaySafe must not bypass existing recovery or alter its write routing")
			})
		}
	}
}

func TestReplaySafeWriteCallerContext(t *testing.T) {
	for _, failure := range []string{"not sent", "EOF", "408", "unclassified", "attempt timeout"} {
		for _, deadline := range []bool{false, true} {
			for _, write := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/deadline=%t/write=%t", failure, deadline, write), func(t *testing.T) {
					t.Parallel()
					ctx, cancel := context.WithCancel(context.Background())
					wantErr := context.Canceled
					if deadline {
						cancel()
						ctx, cancel = context.WithTimeout(context.Background(), 40*time.Millisecond)
						wantErr = context.DeadlineExceeded
					}
					defer cancel()
					sends := 0
					transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
						sends++
						if !deadline {
							timer := time.AfterFunc(10*time.Millisecond, cancel)
							t.Cleanup(func() { timer.Stop() })
						}
						return replaySafeTestFailure(req, failure)
					})
					container, lc := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
						MaxRetries: 3, RetryDelay: time.Second, TryTimeout: time.Second,
					}, false, true)
					var err error
					if write {
						_, err = container.UpsertItem(ctx, NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
					} else {
						_, err = container.ReadItem(ctx, NewPartitionKeyString("pk"), "item", nil)
					}
					assert.ErrorIs(t, err, wantErr)
					assert.Equal(t, 1, sends)
					assert.Empty(t, lc.locationUnavailabilityInfoMap)
				})
			}
		}
	}
}

func TestReplaySafeWriteDownstreamVeto(t *testing.T) {
	cause := errors.New("downstream veto")
	sends := 0
	transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
		sends++
		<-req.Context().Done()
		return nil, fmt.Errorf("wrapped: %w", errorinfo.NonRetriableError(cause))
	})
	container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
		MaxRetries: 3, TryTimeout: 10 * time.Millisecond,
		ShouldRetry: func(*http.Response, error) bool {
			t.Error("a downstream veto must not reach ShouldRetry")
			return true
		},
	}, true, true)
	_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
	assert.ErrorIs(t, err, cause)
	assert.Equal(t, 1, sends)
}

func TestReplaySafeWriteRecoveryDownstreamVeto(t *testing.T) {
	for _, failure := range []string{"not sent", "EOF", "408", "500", "503"} {
		t.Run(failure, func(t *testing.T) {
			t.Parallel()
			cause := errors.New("downstream veto")
			sends, callbacks := 0, 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				wantHost := "write-a.test"
				if sends == 2 && (failure == "408" || failure == "503") {
					wantHost = "write-b.test"
				}
				assert.Equal(t, wantHost, req.URL.Host)
				if sends == 1 {
					return replaySafeTestFailure(req, failure)
				}
				return nil, fmt.Errorf("wrapped: %w", errorinfo.NonRetriableError(cause))
			})
			container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
				MaxRetries: 3, RetryDelay: time.Millisecond,
				ShouldRetry: func(*http.Response, error) bool {
					callbacks++
					t.Error("native recovery and downstream veto must not reach ShouldRetry")
					return true
				},
			}, true, true)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			assert.ErrorIs(t, err, cause)
			assert.Equal(t, 2, sends)
			assert.Zero(t, callbacks)
		})
	}
}

func TestReplaySafeWriteTerminalResponse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusConflict, http.StatusPreconditionFailed, http.StatusNotImplemented} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			sends := 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				return readAttemptTimeoutResponse(req, status), nil
			})
			container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{}, false, true)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			var responseErr *azcore.ResponseError
			require.ErrorAs(t, err, &responseErr)
			assert.Equal(t, status, responseErr.StatusCode)
			assert.Equal(t, 1, sends)
		})
	}

}

func TestReplaySafeWriteBackoffTimeout(t *testing.T) {
	for _, optIn := range []bool{false, true} {
		t.Run(fmt.Sprintf("ReplaySafe=%t", optIn), func(t *testing.T) {
			sends := 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				if sends == 1 {
					return nil, &net.DNSError{Err: "unreachable"}
				}
				return readAttemptTimeoutResponse(req, http.StatusOK), nil
			})
			container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
				MaxRetries: 1, RetryDelay: time.Millisecond, TryTimeout: 20 * time.Millisecond,
			}, false, true)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: optIn})
			if optIn {
				require.NoError(t, err)
				assert.Equal(t, 2, sends)
			} else {
				require.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, 1, sends, "existing unflagged behavior is outside this feature")
			}
		})
	}
}

func TestReplaySafeWriteAmbiguousTimeouts(t *testing.T) {
	for _, transportErr := range []error{
		syscall.ETIMEDOUT,
		&net.OpError{Op: "read", Err: syscall.ETIMEDOUT},
		&net.OpError{Op: "write", Err: syscall.ETIMEDOUT},
		context.DeadlineExceeded,
	} {
		for _, optIn := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/ReplaySafe=%t", transportErr, optIn), func(t *testing.T) {
				t.Parallel()
				sends := 0
				transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
					sends++
					<-req.Context().Done()
					return nil, transportErr
				})
				container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
					MaxRetries: 1, RetryDelay: time.Millisecond, TryTimeout: 20 * time.Millisecond,
				}, false, true)
				_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: optIn})
				require.Error(t, err)
				wantSends := 1
				if optIn {
					wantSends = 2
				}
				assert.Equal(t, wantSends, sends, "ReplaySafe is required to bypass the expired-write guard")
			})
		}
	}
}

func TestReplaySafeWriteServiceUnavailableGuard(t *testing.T) {
	for _, optIn := range []bool{false, true} {
		for _, crossRegion := range []bool{false, true} {
			for _, regions := range []string{"present", "absent", "exhausted"} {
				for _, maxRetries := range []int32{-1, 1} {
					t.Run(fmt.Sprintf("ReplaySafe=%t/crossRegion=%t/regions=%s/MaxRetries=%d", optIn, crossRegion, regions, maxRetries), func(t *testing.T) {
						t.Parallel()
						var hosts []string
						callbacks := 0
						item := []byte(`{"id":"item"}`)
						transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
							hosts = append(hosts, req.URL.Host)
							assert.Equal(t, req.URL.Host, req.Host)
							body, err := io.ReadAll(req.Body)
							require.NoError(t, err)
							assert.Equal(t, item, body, "native and outer retries must rewind the body")
							return readAttemptTimeoutResponse(req, http.StatusServiceUnavailable), nil
						})
						container, lc := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
							MaxRetries: maxRetries, RetryDelay: time.Millisecond,
							ShouldRetry: func(resp *http.Response, err error) bool {
								callbacks++
								require.NoError(t, err)
								assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
								return true
							},
						}, regions == "exhausted", crossRegion)
						if regions == "absent" {
							container.database.client.gem.preferredLocations = nil
						} else if regions == "exhausted" {
							container.database.client.gem.preferredLocations = []string{"East US"}
						}
						_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), item, &ItemOptions{ReplaySafe: optIn})
						require.Error(t, err)
						want := []string{"write-a.test"}
						if regions == "exhausted" && crossRegion {
							want = append(want, "write-b.test")
						}
						assert.Zero(t, callbacks, "shared 503 refusal must veto outer retries regardless of replay safety")
						assert.Equal(t, want, hosts)
						assert.Empty(t, lc.locationUnavailabilityInfoMap)
					})
				}
			}
		}
	}
}

func TestReplaySafeWriteNotSentRecoveryRefusal(t *testing.T) {
	for _, optIn := range []bool{false, true} {
		for _, crossRegion := range []bool{false, true} {
			for _, maxRetries := range []int32{-1, 1, 2} {
				t.Run(fmt.Sprintf("ReplaySafe=%t/crossRegion=%t/MaxRetries=%d", optIn, crossRegion, maxRetries), func(t *testing.T) {
					t.Parallel()
					var hosts []string
					callbacks := 0
					item := []byte(`{"id":"item"}`)
					cause := &net.DNSError{Err: "unreachable"}
					transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
						hosts = append(hosts, req.URL.Host)
						assert.Equal(t, "write-a.test", req.URL.Host, "single-write recovery must not use a read-only target")
						assert.Equal(t, req.URL.Host, req.Host)
						body, err := io.ReadAll(req.Body)
						require.NoError(t, err)
						assert.Equal(t, item, body, "native and outer retries must rewind the body")
						return nil, cause
					})
					container, lc := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
						MaxRetries: maxRetries, RetryDelay: time.Millisecond,
						ShouldRetry: func(resp *http.Response, err error) bool {
							callbacks++
							assert.Nil(t, resp)
							assert.Same(t, cause, err, "native refusal must preserve the transport error")
							return true
						},
					}, false, crossRegion)
					_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), item, &ItemOptions{ReplaySafe: optIn})
					require.ErrorIs(t, err, cause)
					sendsPerAttempt := 1
					if crossRegion {
						sendsPerAttempt += maxSameRegionConnectionRetries
						assert.NotEmpty(t, lc.locationUnavailabilityInfoMap, "preserve native endpoint recovery")
					} else {
						assert.Empty(t, lc.locationUnavailabilityInfoMap)
					}
					assert.Len(t, hosts, sendsPerAttempt, "native refusal must not start another outer budget")
					assert.Zero(t, callbacks)
					assert.True(t, errors.As(err, new(errorinfo.NonRetriable)))
				})
			}
		}
	}
}
