// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestDefaultRetryConfiguration(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	retryPolicy := newResourceThrottleRetryPolicy(cosmosClientOptions)
	if retryPolicy.MaxRetryCount != defaultResourceThrottleRetryPolicyRetryCount {
		t.Errorf("Expected the MaxRetryCount to match the default but got, but got %d", retryPolicy.MaxRetryCount)
	}

	if retryPolicy.MaxWaitTime != defaultResourceThrottleRetryPolicyMaxWaitTime {
		t.Errorf("Expected the MaxWaitTime to match the default but got, but got %s", retryPolicy.MaxWaitTime)
	}
}

func TestCustomRetryConfiguration(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	maxRetryCount := 5
	maxRetryDuration, _ := time.ParseDuration("1s")
	cosmosClientOptions.RateLimitedRetry = &CosmosClientOptionsRateLimitedRetry{
		MaxRetryAttempts: maxRetryCount,
		MaxRetryWaitTime: maxRetryDuration,
	}
	retryPolicy := newResourceThrottleRetryPolicy(cosmosClientOptions)
	if retryPolicy.MaxRetryCount != maxRetryCount {
		t.Errorf("Expected the MaxRetryCount to match the %d but got, but got %d", maxRetryCount, retryPolicy.MaxRetryCount)
	}

	if retryPolicy.MaxWaitTime != maxRetryDuration {
		t.Errorf("Expected the MaxWaitTime to match the %s but got, but got %s", maxRetryDuration, retryPolicy.MaxWaitTime)
	}
}

func TestRetryOn429WithCustomCount(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	maxRetryCount := 5
	maxRetryDuration, _ := time.ParseDuration("5s")
	cosmosClientOptions.RateLimitedRetry = &CosmosClientOptionsRateLimitedRetry{
		MaxRetryAttempts: maxRetryCount,
		MaxRetryWaitTime: maxRetryDuration,
	}
	retryPolicy := newResourceThrottleRetryPolicy(cosmosClientOptions)

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusTooManyRequests))

	pl := azruntime.NewPipeline(srv, retryPolicy)
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, _ := pl.Do(req)
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != maxRetryCount {
		t.Fatalf("wrong request count, got %d expected %d", r, maxRetryCount)
	}
}

func TestRetryOn429WithCustomTime(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	maxRetryCount := 5
	maxRetryDuration, _ := time.ParseDuration("1s")
	cosmosClientOptions.RateLimitedRetry = &CosmosClientOptionsRateLimitedRetry{
		MaxRetryAttempts: maxRetryCount,
		MaxRetryWaitTime: maxRetryDuration,
	}
	retryPolicy := newResourceThrottleRetryPolicy(cosmosClientOptions)

	srv, close := mock.NewTLSServer()
	defer close()
	// Should wait only 1 second and when the retry comes, it should stop
	srv.SetResponse(mock.WithStatusCode(http.StatusTooManyRequests), mock.WithHeader(cosmosHeaderRetryAfter, "1000"))

	pl := azruntime.NewPipeline(srv, retryPolicy)
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, _ := pl.Do(req)
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 2 {
		t.Fatalf("wrong request count, got %d expected %d", r, 2)
	}
}
