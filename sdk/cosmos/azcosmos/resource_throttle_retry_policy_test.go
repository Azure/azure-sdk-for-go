// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
	"time"
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
