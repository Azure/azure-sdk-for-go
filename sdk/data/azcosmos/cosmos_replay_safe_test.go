// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
					ShouldRetry: func(*http.Response, error) bool {
						callbacks++
						return true
					},
				}
				want := []string{"write-a.test", "write-a.test"}
				if failure == "EOF" {
					want = []string{"write-a.test", "write-a.test", "write-a.test", "write-a.test"}
				}
				if failure == "attempt timeout" {
					retry.TryTimeout = 10 * time.Millisecond
					want = []string{"write-a.test", "write-a.test", "write-a.test"}
				} else if multiWrite {
					if failure == "408" {
						want[1] = "write-b.test"
					} else {
						want = append(want, "write-b.test")
					}
				}
				container, lc := newReplaySafeTestContainer(t, transport, retry, multiWrite, true)
				_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
				require.Error(t, err)
				assert.Equal(t, want, hosts, "write retries must never target the preferred read-only region")
				if failure == "attempt timeout" {
					assert.ErrorIs(t, err, context.DeadlineExceeded)
					assert.Equal(t, 3, callbacks)
				} else {
					var nonRetriable errorinfo.NonRetriable
					assert.ErrorAs(t, err, &nonRetriable)
					assert.Zero(t, callbacks, "azcore must not reset an exhausted Cosmos retry budget")
				}
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
		wantSends   int
	}{
		{name: "network cross-region disabled", failure: "EOF", maxRetries: 2, wantSends: 1},
		{name: "408 cross-region disabled", failure: "408", maxRetries: 2, wantSends: 1},
		{name: "500 retains in-region retry", failure: "500", maxRetries: 2, wantSends: 2},
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
			container, _ := newReplaySafeTestContainer(t, transport, policy.RetryOptions{
				MaxRetries: tt.maxRetries, RetryDelay: time.Millisecond, TryTimeout: 20 * time.Millisecond,
			}, true, tt.crossRegion)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			require.Error(t, err)
			assert.Equal(t, tt.wantSends, sends)
		})
	}
}

func TestReplaySafeWriteAttemptTimeoutDuringBackoff(t *testing.T) {
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
				MaxRetries: 1, RetryDelay: time.Millisecond, TryTimeout: 20 * time.Millisecond,
			}, false, true)
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), &ItemOptions{ReplaySafe: true})
			require.NoError(t, err)
			assert.Equal(t, []string{"write-a.test", "write-a.test"}, hosts)
			assert.Empty(t, lc.locationUnavailabilityInfoMap)
		})
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
					MaxRetries: 3, RetryDelay: time.Millisecond, TryTimeout: time.Second,
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
	for _, status := range []int{http.StatusNotFound, http.StatusConflict, http.StatusPreconditionFailed, http.StatusServiceUnavailable} {
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
