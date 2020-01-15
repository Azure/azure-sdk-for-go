// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultMaxTries = 4
)

// RetryOptions configures the retry policy's behavior.
type RetryOptions struct {
	// MaxTries specifies the maximum number of attempts an operation will be tried before producing an error (0=default).
	// A value of zero means that you accept our default policy. A value of 1 means 1 try and no retries.
	MaxTries int32

	// TryTimeout indicates the maximum time allowed for any single try of an HTTP request.
	// A value of zero means that you accept our default timeout. NOTE: When transferring large amounts
	// of data, the default TryTimeout will probably not be sufficient. You should override this value
	// based on the bandwidth available to the host machine and proximity to the service. A good
	// starting point may be something like (60 seconds per MB of anticipated-payload-size).
	TryTimeout time.Duration

	// RetryDelay specifies the amount of delay to use before retrying an operation (0=default).
	// The delay increases exponentially with each retry up to a maximum specified by MaxRetryDelay.
	// If you specify 0, then you must also specify 0 for MaxRetryDelay.
	// If you specify RetryDelay, then you must also specify MaxRetryDelay, and MaxRetryDelay should be
	// equal to or greater than RetryDelay.
	RetryDelay time.Duration

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation (0=default).
	// If you specify 0, then you must also specify 0 for RetryDelay.
	MaxRetryDelay time.Duration

	// StatusCodes specifies the HTTP status codes that indicate the operation should be retried.
	// If unspecified it will default to the status codes in StatusCodesForRetry.
	StatusCodes []int
}

var (
	// StatusCodesForRetry is the default set of HTTP status code for which the policy will retry.
	StatusCodesForRetry = [6]int{
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}
)

func (o RetryOptions) defaults() RetryOptions {
	// We assume the following:
	// 1. o.MaxTries >= 0
	// 2. o.TryTimeout, o.RetryDelay, and o.MaxRetryDelay >=0
	// 3. o.RetryDelay <= o.MaxRetryDelay
	// 4. Both o.RetryDelay and o.MaxRetryDelay must be 0 or neither can be 0

	if len(o.StatusCodes) == 0 {
		o.StatusCodes = StatusCodesForRetry[:]
	}

	IfDefault := func(current *time.Duration, desired time.Duration) {
		if *current == time.Duration(0) {
			*current = desired
		}
	}

	// Set defaults if unspecified
	if o.MaxTries == 0 {
		o.MaxTries = defaultMaxTries
	}

	IfDefault(&o.TryTimeout, 1*time.Minute)
	IfDefault(&o.RetryDelay, 4*time.Second)
	IfDefault(&o.MaxRetryDelay, 120*time.Second)

	return o
}

func (o RetryOptions) calcDelay(try int32) time.Duration { // try is >=1; never 0
	pow := func(number int64, exponent int32) int64 { // pow is nested helper function
		var result int64 = 1
		for n := int32(0); n < exponent; n++ {
			result *= number
		}
		return result
	}

	delay := time.Duration(pow(2, try)-1) * o.RetryDelay

	// Introduce some jitter:  [0.0, 1.0) / 2 = [0.0, 0.5) + 0.8 = [0.8, 1.3)
	delay = time.Duration(delay.Seconds() * (rand.Float64()/2 + 0.8) * float64(time.Second)) // NOTE: We want math/rand; not crypto/rand
	if delay > o.MaxRetryDelay {
		delay = o.MaxRetryDelay
	}
	return delay
}

// NewRetryPolicy creates a policy object configured using the specified options.
func NewRetryPolicy(o RetryOptions) Policy {
	return &retryPolicy{options: o.defaults()} // Force defaults to be calculated
}

type retryPolicy struct {
	options RetryOptions
}

func (p *retryPolicy) Do(ctx context.Context, req *Request) (resp *Response, err error) {
	// Exponential retry algorithm: ((2 ^ attempt) - 1) * delay * random(0.8, 1.2)
	// When to retry: connection failure or temporary/timeout.
	if req.Body != nil {
		// wrap the body so we control when it's actually closed
		rwbody := &retryableRequestBody{body: req.Body.(ReadSeekCloser)}
		req.Body = rwbody
		req.Request.GetBody = func() (io.ReadCloser, error) {
			_, err := rwbody.Seek(0, io.SeekStart) // Seek back to the beginning of the stream
			return rwbody, err
		}
		defer rwbody.realClose()
	}
	try := int32(1)
	shouldLog := Log().Should(LogRetryPolicy)
	for {
		resp = nil // reset
		if shouldLog {
			Log().Write(LogRetryPolicy, fmt.Sprintf("\n=====> Try=%d\n", try))
		}

		// For each try, seek to the beginning of the Body stream. We do this even for the 1st try because
		// the stream may not be at offset 0 when we first get it and we want the same behavior for the
		// 1st try as for additional tries.
		err = req.RewindBody()
		if err != nil {
			return
		}

		// Set the time for this particular retry operation and then Do the operation.
		tryCtx, tryCancel := context.WithTimeout(ctx, p.options.TryTimeout)
		resp, err = req.Next(tryCtx) // Make the request
		tryCancel()
		if shouldLog {
			Log().Write(LogRetryPolicy, fmt.Sprintf("Err=%v, response=%v\n", err, resp))
		}

		if err == nil && !resp.HasStatusCode(p.options.StatusCodes...) {
			// if there is no error and the response code isn't in the list of retry codes then we're done.
			return
		} else if ctx.Err() != nil {
			// don't retry if the parent context has been cancelled or its deadline exceeded
			return
		} else if retrier, ok := err.(Retrier); ok && retrier.IsNotRetriable() {
			// the error says it's not retriable so don't retry
			return
		}

		// drain before retrying so nothing is leaked
		resp.Drain()

		if try == p.options.MaxTries {
			// max number of tries has been reached, don't sleep again
			return
		}

		// use the delay from retry-after if available
		delay, ok := resp.RetryAfter()
		if !ok {
			delay = p.options.calcDelay(try)
		}
		if shouldLog {
			Log().Write(LogRetryPolicy, fmt.Sprintf("Try=%d, Delay=%v\n", try, delay))
		}
		select {
		case <-time.After(delay):
			try++
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

// ********** The following type/methods implement the retryableRequestBody (a ReadSeekCloser)

// This struct is used when sending a body to the network
type retryableRequestBody struct {
	body io.ReadSeeker // Seeking is required to support retries
}

// Read reads a block of data from an inner stream and reports progress
func (b *retryableRequestBody) Read(p []byte) (n int, err error) {
	return b.body.Read(p)
}

func (b *retryableRequestBody) Seek(offset int64, whence int) (offsetFromStart int64, err error) {
	return b.body.Seek(offset, whence)
}

func (b *retryableRequestBody) Close() error {
	// We don't want the underlying transport to close the request body on transient failures so this is a nop.
	// The retry policy closes the request body upon success.
	return nil
}

func (b *retryableRequestBody) realClose() error {
	if c, ok := b.body.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
