// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const defaultMaxWaitTimeInSeconds time.Duration = 60 * time.Second
const defaultRetryInSeconds int = 5

// cosmosRetryPolicyThrottle retries on HTTP 429.
type cosmosRetryPolicyThrottle struct {
	MaxWaitTime   time.Duration
	MaxRetryCount int
}

func NewCosmosRetryPolicyThrottle(o *CosmosClientOptions) *cosmosRetryPolicyThrottle {
	if o.RateLimitedRetry == nil {
		return &cosmosRetryPolicyThrottle{MaxWaitTime: defaultMaxWaitTimeInSeconds, MaxRetryCount: defaultRetryInSeconds}
	}

	return &cosmosRetryPolicyThrottle{MaxWaitTime: o.RateLimitedRetry.MaxRetryWaitTime, MaxRetryCount: o.RateLimitedRetry.MaxRetryAttempts}
}

func (p *cosmosRetryPolicyThrottle) Do(req *azcore.Request) (*azcore.Response, error) {
	// Policy disabled
	if p.MaxRetryCount == 0 {
		return req.Next()
	}

	var resp *azcore.Response
	var err error
	var cummulativeWaitTime time.Duration
	for attempts := 0; attempts < p.MaxRetryCount; attempts++ {
		// make the original request
		resp, err = req.Next()

		if err != nil || resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}

		retryAfter := resp.Header.Get(cosmosHeaderRetryAfter)
		retryAfterDuration := parseRetryAfter(retryAfter)
		cummulativeWaitTime += retryAfterDuration

		if retryAfterDuration > p.MaxWaitTime || cummulativeWaitTime > p.MaxWaitTime {
			return resp, err
		}
	}

	return resp, err
}

func parseRetryAfter(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 0
	}

	retryAfterInMilliseconds, err := time.ParseDuration(retryAfter + "ms")
	if err != nil {
		return 0
	}

	return retryAfterInMilliseconds
}
