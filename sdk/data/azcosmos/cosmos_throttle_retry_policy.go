// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strconv"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const (
	defaultMaxThrottleRetryAttempts = 9
	defaultMaxThrottleRetryWaitTime = 30 * time.Second
	defaultThrottleRetryDelay       = 5 * time.Second
)

// throttleRetryPolicy retries requests that fail with HTTP 429 (Too Many Requests).
// It honors the Cosmos-specific x-ms-retry-after-ms header to determine the
// delay between attempts and caps the number of attempts and total cumulative
// retry delay. This matches the throttling retry behavior of the other Cosmos
// SDKs (.NET, Java, Python).
type throttleRetryPolicy struct {
	maxRetryAttempts int
	maxRetryWaitTime time.Duration
	// defaultDelay is used when a 429 response is missing the
	// x-ms-retry-after-ms header. Defaults to defaultThrottleRetryDelay.
	defaultDelay time.Duration
}

// newThrottleRetryPolicy constructs a throttleRetryPolicy. Non-positive values
// for either option fall back to defaults.
func newThrottleRetryPolicy(o *ThrottlingRetryOptions) *throttleRetryPolicy {
	p := &throttleRetryPolicy{
		maxRetryAttempts: defaultMaxThrottleRetryAttempts,
		maxRetryWaitTime: defaultMaxThrottleRetryWaitTime,
		defaultDelay:     defaultThrottleRetryDelay,
	}
	if o != nil {
		if o.MaxRetryAttempts > 0 {
			p.maxRetryAttempts = o.MaxRetryAttempts
		} else if o.MaxRetryAttempts < 0 {
			// negative values disable throttling retries entirely
			p.maxRetryAttempts = 0
		}
		if o.MaxRetryWaitTime > 0 {
			p.maxRetryWaitTime = o.MaxRetryWaitTime
		}
	}
	return p
}

func (p *throttleRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	attemptCount := 0
	cumulativeDelay := time.Duration(0)
	for {
		response, err := req.Next()
		// Transport / non-HTTP errors are not throttling; let other policies decide.
		if err != nil || response == nil || response.StatusCode != http.StatusTooManyRequests {
			return response, err
		}

		if attemptCount >= p.maxRetryAttempts {
			log.Writef(azlog.EventRetryPolicy, "Cosmos throttle retry exhausted attempts (%d); returning 429 to caller", p.maxRetryAttempts)
			return response, nil
		}

		delay := readRetryAfterMs(response)
		if delay <= 0 {
			delay = p.defaultDelay
		}

		if cumulativeDelay+delay > p.maxRetryWaitTime {
			log.Writef(azlog.EventRetryPolicy, "Cosmos throttle retry exceeded cumulative wait time (%s); returning 429 to caller", p.maxRetryWaitTime)
			return response, nil
		}

		cumulativeDelay += delay
		attemptCount++

		// drain and close the response body so the connection can be reused
		azruntime.Drain(response)

		if err := req.RewindBody(); err != nil {
			return response, err
		}

		log.Writef(azlog.EventRetryPolicy, "Cosmos throttle retry attempt %d after %s (cumulative %s)", attemptCount, delay, cumulativeDelay)

		select {
		case <-time.After(delay):
		case <-req.Raw().Context().Done():
			return response, req.Raw().Context().Err()
		}
	}
}

// readRetryAfterMs parses the Cosmos x-ms-retry-after-ms header (milliseconds).
// Returns 0 if the header is missing or cannot be parsed.
func readRetryAfterMs(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	v := resp.Header.Get(cosmosHeaderRetryAfterMs)
	if v == "" {
		return 0
	}
	ms, err := strconv.ParseFloat(v, 64)
	if err != nil || ms < 0 {
		return 0
	}
	return time.Duration(ms * float64(time.Millisecond))
}
