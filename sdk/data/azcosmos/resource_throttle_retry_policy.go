// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const (
	defaultResourceThrottleRetryPolicyMaxWaitTime time.Duration = 60 * time.Second
	defaultResourceThrottleRetryPolicyRetryCount  int           = 9
)

// resourceThrottleRetryPolicy retries on HTTP 429.
type resourceThrottleRetryPolicy struct {
	MaxWaitTime   time.Duration
	MaxRetryCount int
}

func newResourceThrottleRetryPolicy(o *CosmosClientOptions) *resourceThrottleRetryPolicy {
	if o.RateLimitedRetry == nil {
		return &resourceThrottleRetryPolicy{
			MaxWaitTime:   defaultResourceThrottleRetryPolicyMaxWaitTime,
			MaxRetryCount: defaultResourceThrottleRetryPolicyRetryCount}
	}

	return &resourceThrottleRetryPolicy{
		MaxWaitTime:   o.RateLimitedRetry.MaxRetryWaitTime,
		MaxRetryCount: o.RateLimitedRetry.MaxRetryAttempts}
}

func (p *resourceThrottleRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Policy disabled
	if p.MaxRetryCount == 0 {
		return req.Next()
	}

	var resp *http.Response
	var err error
	var cumulativeWaitTime time.Duration
	for attempts := 0; attempts < p.MaxRetryCount; attempts++ {
		err = req.RewindBody()
		if err != nil {
			return resp, err
		}

		resp, err = req.Next()
		if err != nil || resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}

		retryAfter := resp.Header.Get(cosmosHeaderRetryAfter)
		retryAfterDuration := parseRetryAfter(retryAfter)
		cumulativeWaitTime += retryAfterDuration

		if retryAfterDuration > p.MaxWaitTime || cumulativeWaitTime > p.MaxWaitTime {
			return resp, err
		}

		// drain before retrying so nothing is leaked
		azruntime.Drain(resp)

		select {
		case <-time.After(retryAfterDuration):
			// retry
		case <-req.Raw().Context().Done():
			err = req.Raw().Context().Err()
			log.Writef(log.RetryPolicy, "ResourceThrottleRetryPolicy abort due to %v", err)
			return resp, err
		}
	}

	return resp, err
}

func parseRetryAfter(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 0
	}

	retryAfterDuration, err := time.ParseDuration(retryAfter + "ms")
	if err != nil {
		return 0
	}

	return retryAfterDuration
}
