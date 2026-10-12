// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timeoutTestTransport func(*http.Request) (*http.Response, error)

func (f timeoutTestTransport) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTimeoutTestContainer(t *testing.T, transport policy.Transporter, retry policy.RetryOptions) (*ContainerClient, *globalEndpointManager) {
	t.Helper()
	endpoint, err := url.Parse("https://cosmos.test")
	require.NoError(t, err)
	gem := &globalEndpointManager{
		clientEndpoint:      endpoint.String(),
		locationCache:       CreateMockLC(*endpoint, true),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Now(),
	}
	// Use cached account metadata so the production pipeline needs no service.
	gem.everPopulated.Store(true)
	client, err := newClient(policyFunc(func(req *policy.Request) (*http.Response, error) {
		return req.Next()
	}), gem, &ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport, Retry: retry}})
	require.NoError(t, err)
	c := &Client{endpoint: endpoint.String(), endpointUrl: endpoint, internal: client, gem: gem}
	container, err := c.NewContainer("database", "container")
	require.NoError(t, err)
	return container, gem
}

func timeoutTestResponse(req *http.Request, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"id":"item"}`)),
		Request:    req,
	}
}

func timeoutTestOperation(ctx context.Context, container *ContainerClient, write bool) error {
	if write {
		_, err := container.CreateItem(ctx, NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), nil)
		return err
	}
	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "item", nil)
	return err
}

func TestClientRetryPolicyTryTimeout(t *testing.T) {
	for _, write := range []bool{false, true} {
		for _, failure := range []string{"transport", "408Backoff"} {
			if write && failure == "408Backoff" {
				// HTTP 408 writes must not enter a backoff.
				continue
			}
			for _, persistent := range []bool{false, true} {
				name := failure + "/read"
				if write {
					name = failure + "/write"
				}
				if persistent {
					name += "/exhausted"
				} else {
					name += "/recovered"
				}
				t.Run(name, func(t *testing.T) {
					const maxRetries = 2
					var attempts []context.Context
					transport := timeoutTestTransport(func(req *http.Request) (*http.Response, error) {
						attempts = append(attempts, req.Context())
						if write {
							body, err := io.ReadAll(req.Body)
							require.NoError(t, err)
							assert.JSONEq(t, `{"id":"item"}`, string(body), "write body must be rewound for each outer attempt")
						}
						if !persistent && len(attempts) > 1 {
							return timeoutTestResponse(req, http.StatusOK), nil
						}
						switch failure {
						case "transport":
							<-req.Context().Done()
							return nil, &url.Error{Op: req.Method, URL: req.URL.String(), Err: req.Context().Err()}
						default:
							return timeoutTestResponse(req, http.StatusRequestTimeout), nil
						}
					})
					container, gem := newTimeoutTestContainer(t, transport, policy.RetryOptions{
						MaxRetries: maxRetries,
						TryTimeout: 100 * time.Millisecond,
						RetryDelay: time.Millisecond,
					})
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					err := timeoutTestOperation(ctx, container, write)
					require.NoError(t, ctx.Err(), "only the per-try context should expire")
					if persistent {
						require.ErrorIs(t, err, context.DeadlineExceeded)
						assert.Len(t, attempts, maxRetries+1, "azcore must enforce the outer retry budget")
						var nonRetriable errorinfo.NonRetriable
						assert.NotErrorAs(t, err, &nonRetriable)
					} else {
						require.NoError(t, err)
						require.Len(t, attempts, 2, "a timed-out try must reach the next outer attempt")
						assert.True(t, attempts[0].Done() != attempts[1].Done(), "each outer attempt must get a fresh timeout context")
					}
					assert.Empty(t, gem.locationCache.locationUnavailabilityInfoMap,
						"an expired try must not mark an endpoint unavailable")
				})
			}
		}
	}
}

func TestClientRetryPolicyTryTimeoutPreservesWriteSafety(t *testing.T) {
	for _, failure := range []string{"408", "500", "ambiguousTransport", "transportTimeout"} {
		t.Run(failure, func(t *testing.T) {
			requests := 0
			transport := timeoutTestTransport(func(req *http.Request) (*http.Response, error) {
				requests++
				require.NoError(t, req.Context().Err(), "the try is still alive when the failure occurs")
				switch failure {
				case "408":
					return timeoutTestResponse(req, http.StatusRequestTimeout), nil
				case "500":
					return timeoutTestResponse(req, http.StatusInternalServerError), nil
				case "transportTimeout":
					return nil, &net.OpError{Op: "read", Net: "tcp", Err: context.DeadlineExceeded}
				default:
					return nil, io.EOF
				}
			})
			container, gem := newTimeoutTestContainer(t, transport, policy.RetryOptions{
				MaxRetries: 3,
				TryTimeout: 100 * time.Millisecond,
				RetryDelay: time.Millisecond,
			})
			err := timeoutTestOperation(context.Background(), container, true)
			require.Error(t, err)
			var nonRetriable errorinfo.NonRetriable
			assert.ErrorAs(t, err, &nonRetriable)
			assert.Equal(t, 1, requests, "ambiguous writes must not enter either retry loop")
			assert.Empty(t, gem.locationCache.locationUnavailabilityInfoMap)
		})
	}
}

func TestClientRetryPolicyCallerContext(t *testing.T) {
	for _, write := range []bool{false, true} {
		for _, failure := range []string{"transport", "408Backoff"} {
			if write && failure == "408Backoff" {
				continue
			}

			for _, deadline := range []bool{false, true} {
				name := failure + "/read"
				if write {
					name = failure + "/write"
				}
				if deadline {
					name += "/deadline"
				} else {
					name += "/cancel"
				}
				t.Run(name, func(t *testing.T) {
					ctx, cancel := context.WithCancel(context.Background())
					wantErr := context.Canceled
					if deadline {
						cancel()
						ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
						wantErr = context.DeadlineExceeded
					}
					defer cancel()
					if !deadline {
						timer := time.AfterFunc(100*time.Millisecond, cancel)
						defer timer.Stop()
					}
					requests := 0
					transport := timeoutTestTransport(func(req *http.Request) (*http.Response, error) {
						requests++
						switch failure {
						case "transport":
							<-req.Context().Done()
							return nil, req.Context().Err()
						default:
							return timeoutTestResponse(req, http.StatusRequestTimeout), nil
						}
					})
					container, _ := newTimeoutTestContainer(t, transport, policy.RetryOptions{
						MaxRetries: 10,
						TryTimeout: 2 * time.Second,
						RetryDelay: time.Millisecond,
					})
					start := time.Now()
					err := timeoutTestOperation(ctx, container, write)
					require.ErrorIs(t, err, wantErr)
					assert.Equal(t, 1, requests, "caller cancellation must stop both retry loops")
					assert.Less(t, time.Since(start), time.Second, "caller cancellation must interrupt backoff")
				})
			}
		}
	}
}
