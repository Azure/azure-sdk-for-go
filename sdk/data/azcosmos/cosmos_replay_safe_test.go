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
	case "EOF":
		return nil, io.EOF
	case "attempt timeout":
		<-req.Context().Done()
		return nil, req.Context().Err()
	case "408":
		return readAttemptTimeoutResponse(req, http.StatusRequestTimeout), nil
	case "500":
		return readAttemptTimeoutResponse(req, http.StatusInternalServerError), nil
	case "502":
		return readAttemptTimeoutResponse(req, http.StatusBadGateway), nil
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

func TestReplaySafeWriteExhaustion(t *testing.T) {
	for _, multiWrite := range []bool{false, true} {
		for _, failure := range []string{"EOF", "attempt timeout", "408", "500", "502", "504"} {
			t.Run(fmt.Sprintf("multiWrite=%t/%s", multiWrite, failure), func(t *testing.T) {
				t.Parallel()
				var hosts []string
				callbacks := 0
				transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
					hosts = append(hosts, req.URL.Host)
					return replaySafeTestFailure(req, failure)
				})
				retry := policy.RetryOptions{
					MaxRetries: 2, RetryDelay: time.Millisecond,
					ShouldRetry: func(resp *http.Response, err error) bool {
						callbacks++
						if failure == "EOF" {
							assert.ErrorIs(t, err, io.EOF)
						} else if failure != "attempt timeout" {
							require.NotNil(t, resp)
							assert.NoError(t, err, "HTTP responses must reach azcore without error wrappers")
						}
						return true
					},
				}
				want := []string{"write-a.test", "write-a.test", "write-a.test"}
				if failure == "attempt timeout" {
					retry.TryTimeout = 10 * time.Millisecond
				}
				container, lc := newReplaySafeTestContainer(t, transport, retry, multiWrite, true)
				_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
				require.Error(t, err)
				assert.Equal(t, want, hosts, "write retries must never target the preferred read-only region")
				if failure == "attempt timeout" {
					assert.ErrorIs(t, err, context.DeadlineExceeded)
				}
				assert.Equal(t, 3, callbacks)
				var nonRetriable errorinfo.NonRetriable
				assert.False(t, errors.As(err, &nonRetriable), "eligible failures must reach outer retries unwrapped")
				assert.Empty(t, lc.locationUnavailabilityInfoMap)
			})
		}
	}
}

func TestReplaySafeWriteRetrySettings(t *testing.T) {
	for _, tt := range []struct {
		name        string
		failure     string
		crossRegion bool
		maxRetries  int32
		statusCodes []int
		reject      bool
		wantSends   int
	}{
		{name: "network cross-region disabled", failure: "EOF", maxRetries: 2, wantSends: 3},
		{name: "408 cross-region disabled", failure: "408", maxRetries: 2, wantSends: 3},
		{name: "500 uses outer budget", failure: "500", maxRetries: 2, wantSends: 3},
		{name: "network retries disabled", failure: "EOF", crossRegion: true, maxRetries: -1, wantSends: 1},
		{name: "HTTP retries disabled", failure: "500", crossRegion: true, maxRetries: -1, wantSends: 1},
		{name: "empty status list", failure: "408", crossRegion: true, maxRetries: 2, statusCodes: []int{}, wantSends: 1},
		{name: "excluded status", failure: "500", crossRegion: true, maxRetries: 2, statusCodes: []int{502}, wantSends: 1},
		{name: "included status", failure: "502", crossRegion: true, maxRetries: 2, statusCodes: []int{502}, wantSends: 3},
		{name: "ShouldRetry veto", failure: "EOF", crossRegion: true, maxRetries: 2, reject: true, wantSends: 1},
		{name: "attempt retries disabled", failure: "attempt timeout", crossRegion: true, maxRetries: -1, wantSends: 1},
		{name: "attempt retries bounded", failure: "attempt timeout", crossRegion: true, maxRetries: 1, wantSends: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sends := 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				return replaySafeTestFailure(req, tt.failure)
			})
			retry := policy.RetryOptions{
				MaxRetries: tt.maxRetries, RetryDelay: time.Millisecond, TryTimeout: 20 * time.Millisecond, StatusCodes: tt.statusCodes,
			}
			if tt.reject {
				retry.ShouldRetry = func(*http.Response, error) bool { return false }
			}
			container, _ := newReplaySafeTestContainer(t, transport, retry, true, tt.crossRegion)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			require.Error(t, err)
			assert.Equal(t, tt.wantSends, sends)
		})
	}
}

func TestReplaySafeWriteOuterDelay(t *testing.T) {
	for _, failure := range []string{"EOF", "408"} {
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
			}, false, true)
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
	for _, capped := range []bool{false, true} {
		t.Run(fmt.Sprintf("RetryAfterCapped=%t", capped), func(t *testing.T) {
			t.Parallel()
			var first *http.Response
			sends, callbacks := 0, 0
			transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
				sends++
				if sends == 1 {
					first = readAttemptTimeoutResponse(req, http.StatusBadGateway)
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
					return resp.StatusCode == http.StatusBadGateway
				},
			}
			if capped {
				retry.MaxRetryDelay = 10 * time.Millisecond
			}
			container, _ := newReplaySafeTestContainer(t, transport, retry, true, true)
			start := time.Now()
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			if capped {
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
		return nil, io.EOF
	})
	container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{MaxRetries: -1}, false, false)
	ctx := policy.WithRetryOptions(context.Background(), policy.RetryOptions{MaxRetries: 2, RetryDelay: time.Millisecond})
	_, err := container.UpsertItem(ctx, NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
	assert.ErrorIs(t, err, io.EOF)
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
	for _, failure := range []string{"EOF", "408", "attempt timeout"} {
		for _, deadline := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/deadline=%t", failure, deadline), func(t *testing.T) {
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
				}, true, true)
				_, err := container.UpsertItem(ctx, NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
				assert.ErrorIs(t, err, wantErr)
				assert.Equal(t, 1, sends)
				assert.Empty(t, lc.locationUnavailabilityInfoMap)
			})
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
			for _, maxRetries := range []int32{-1, 1} {
				t.Run(fmt.Sprintf("ReplaySafe=%t/crossRegion=%t/MaxRetries=%d", optIn, crossRegion, maxRetries), func(t *testing.T) {
					t.Parallel()
					var hosts []string
					callbacks := 0
					transport := retryPolicyTransportFunc(func(req *http.Request) (*http.Response, error) {
						hosts = append(hosts, req.URL.Host)
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
					}, false, crossRegion)
					_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: optIn})
					require.Error(t, err)
					want := []string{"write-a.test"}
					if optIn && crossRegion {
						if maxRetries > 0 {
							want = append(want, "write-a.test")
						}
						assert.Equal(t, len(want), callbacks, "bypass the write-specific single-region veto, not routing")
					} else {
						assert.Zero(t, callbacks, "retain independent cross-region-disabled veto and unsafe-write default")
					}
					assert.Equal(t, want, hosts)
					assert.Empty(t, lc.locationUnavailabilityInfoMap)
				})
			}
		}
	}
}
