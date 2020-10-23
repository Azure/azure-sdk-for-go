// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultMaxRetries = 3
)

// RetryOptions configures the retry policy's behavior.
type RetryOptions struct {
	// MaxRetries specifies the maximum number of attempts a failed operation will be retried
	// before producing an error.  A value of zero means one try and no retries.
	MaxRetries int32

	// TryTimeout indicates the maximum time allowed for any single try of an HTTP request.
	TryTimeout time.Duration

	// RetryDelay specifies the amount of delay to use before retrying an operation.
	// The delay increases exponentially with each retry up to a maximum specified by MaxRetryDelay.
	// If you specify 0, then you must also specify 0 for MaxRetryDelay.
	// If you specify RetryDelay, then you must also specify MaxRetryDelay, and MaxRetryDelay should be
	// equal to or greater than RetryDelay.
	RetryDelay time.Duration

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
	// If you specify 0, then you must also specify 0 for RetryDelay.
	MaxRetryDelay time.Duration

	// StatusCodes specifies the HTTP status codes that indicate the operation should be retried.
	// If unspecified it will default to the status codes in StatusCodesForRetry.
	StatusCodes []int
}

var (
	// StatusCodesForRetry is the default set of HTTP status code for which the policy will retry.
	// Changing its value will affect future created clients that use the default values.
	StatusCodesForRetry = []int{
		http.StatusRequestTimeout,      // 408
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}
)

// DefaultRetryOptions returns an instance of RetryOptions initialized with default values.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		StatusCodes:   StatusCodesForRetry,
		MaxRetries:    defaultMaxRetries,
		TryTimeout:    1 * time.Minute,
		RetryDelay:    4 * time.Second,
		MaxRetryDelay: 120 * time.Second,
	}
}

// used as a context key for adding/retrieving RetryOptions
type ctxWithRetryOptionsKey struct{}

// WithRetryOptions adds the specified RetryOptions to the parent context.
// Use this to specify custom RetryOptions at the API-call level.
func WithRetryOptions(parent context.Context, options RetryOptions) context.Context {
	return context.WithValue(parent, ctxWithRetryOptionsKey{}, options)
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
// Pass nil to accept the default values; this is the same as passing the result
// from a call to DefaultRetryOptions().
func NewRetryPolicy(o *RetryOptions) Policy {
	if o == nil {
		def := DefaultRetryOptions()
		o = &def
	}
	return &retryPolicy{options: *o}
}

type retryPolicy struct {
	options RetryOptions
}

func (p *retryPolicy) Do(req *Request) (resp *Response, err error) {
	options := p.options
	// check if the retry options have been overridden for this call
	if override := req.Context().Value(ctxWithRetryOptionsKey{}); override != nil {
		options = override.(RetryOptions)
	}
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
	for {
		resp = nil // reset
		Log().Writef(LogRetryPolicy, "\n=====> Try=%d %s %s", try, req.Method, req.URL.String())

		// For each try, seek to the beginning of the Body stream. We do this even for the 1st try because
		// the stream may not be at offset 0 when we first get it and we want the same behavior for the
		// 1st try as for additional tries.
		err = req.RewindBody()
		if err != nil {
			return
		}

		// Set the per-try time for this particular retry operation and then Do the operation.
		tryCtx, tryCancel := context.WithTimeout(req.Context(), options.TryTimeout)
		clone := req.clone(tryCtx)
		resp, err = clone.Next() // Make the request
		tryCancel()
		if err == nil {
			Log().Writef(LogRetryPolicy, "response %d", resp.StatusCode)
		} else {
			Log().Writef(LogRetryPolicy, "error %v", err)
		}

		if err == nil && !resp.HasStatusCode(options.StatusCodes...) {
			// if there is no error and the response code isn't in the list of retry codes then we're done.
			return
		} else if ctxErr := req.Context().Err(); ctxErr != nil {
			// don't retry if the parent context has been cancelled or its deadline exceeded
			err = ctxErr
			Log().Writef(LogRetryPolicy, "abort due to %v", err)
			return
		}

		// check if the error is not retriable
		var nre NonRetriableError
		if errors.As(err, &nre) {
			// the error says it's not retriable so don't retry
			Log().Writef(LogRetryPolicy, "non-retriable error %T", nre)
			return
		}

		if try == options.MaxRetries+1 {
			// max number of tries has been reached, don't sleep again
			Log().Writef(LogRetryPolicy, "MaxRetries %d exceeded", options.MaxRetries)
			return
		}

		// drain before retrying so nothing is leaked
		resp.Drain()

		// use the delay from retry-after if available
		delay := resp.retryAfter()
		if delay <= 0 {
			delay = options.calcDelay(try)
		}
		Log().Writef(LogRetryPolicy, "End Try #%d, Delay=%v", try, delay)
		select {
		case <-time.After(delay):
			try++
		case <-req.Context().Done():
			err = req.Context().Err()
			Log().Writef(LogRetryPolicy, "abort due to %v", err)
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
