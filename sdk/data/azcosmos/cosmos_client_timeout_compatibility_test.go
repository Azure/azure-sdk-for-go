// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type compatibilityTransport func(*http.Request) (*http.Response, error)

func (f compatibilityTransport) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type compatibilityTimeoutBody struct {
	ctx context.Context
}

func (b *compatibilityTimeoutBody) Read([]byte) (int, error) {
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *compatibilityTimeoutBody) Close() error { return nil }

type compatibilityNonRetriableError struct{ error }

func (e *compatibilityNonRetriableError) NonRetriable() {}
func (e *compatibilityNonRetriableError) Unwrap() error { return e.error }

func TestTransportTimeoutCompatibility(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPatch} {
		for _, failure := range []string{"transport", "markedTransport", "responseBody"} {
			t.Run(method+"/"+failure, func(t *testing.T) {
				requests := 0
				var originalErr error
				transport := compatibilityTransport(func(req *http.Request) (*http.Response, error) {
					response := &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Request:    req,
					}
					if req.URL.Path == "/" {
						response.Body = io.NopCloser(strings.NewReader(`{
							"readableLocations":[{"name":"East US","databaseAccountEndpoint":"https://cosmos.test/"}],
							"writableLocations":[{"name":"East US","databaseAccountEndpoint":"https://cosmos.test/"}]
						}`))
						return response, nil
					}
					requests++
					if requests > 1 {
						response.Body = io.NopCloser(strings.NewReader(`{"id":"item"}`))
						return response, nil
					}
					if failure == "responseBody" {
						response.Body = &compatibilityTimeoutBody{ctx: req.Context()}
						return response, nil
					}
					<-req.Context().Done()
					originalErr = req.Context().Err()
					if failure == "markedTransport" {
						originalErr = &compatibilityNonRetriableError{error: originalErr}
					}
					return nil, originalErr
				})
				credential, err := azcosmos.NewKeyCredential("dGVzdA==")
				require.NoError(t, err)
				client, err := azcosmos.NewClientWithKey("https://cosmos.test/", credential, &azcosmos.ClientOptions{
					ClientOptions: azcore.ClientOptions{
						Transport: transport,
						Retry: policy.RetryOptions{
							MaxRetries: 2,
							TryTimeout: 100 * time.Millisecond,
							RetryDelay: time.Millisecond,
						},
					},
				})
				require.NoError(t, err)
				container, err := client.NewContainer("database", "container")
				require.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				switch method {
				case http.MethodGet:
					_, err = container.ReadItem(ctx, azcosmos.NewPartitionKeyString("pk"), "item", nil)
				case http.MethodPost:
					_, err = container.CreateItem(ctx, azcosmos.NewPartitionKeyString("pk"), []byte(`{"id":"item"}`), nil)
				case http.MethodPatch:
					patch := azcosmos.PatchOperations{}
					patch.AppendIncrement("/value", 1)
					_, err = container.PatchItem(ctx, azcosmos.NewPartitionKeyString("pk"), "item", patch, nil)
				}
				require.NoError(t, ctx.Err())
				if failure == "markedTransport" || (failure == "responseBody" && method != http.MethodGet) {
					require.ErrorIs(t, err, context.DeadlineExceeded)
					var marker interface{ NonRetriable() }
					assert.True(t, errors.As(err, &marker), "downstream non-retriable classification must survive")
					assert.Equal(t, 1, requests, "non-retriable errors must not replay the operation")
					if originalErr != nil {
						assert.ErrorIs(t, err, originalErr, "preserve the original downstream error")
					}
				} else {
					require.NoError(t, err)
					assert.Equal(t, 2, requests, "unmarked timeouts retain v1.4.2 retry behavior")
				}
			})
		}
	}
}
