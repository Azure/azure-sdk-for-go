// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

// func (o azcore.RetryOptions) msiDefaults() azcore.RetryOptions {
// 	// We assume the following:
// 	// 1. o.Policy should either be RetryPolicyExponential or RetryPolicyFixed
// 	// 2. o.MaxTries >= 0
// 	// 3. o.TryTimeout, o.RetryDelay, and o.MaxRetryDelay >=0
// 	// 4. o.RetryDelay <= o.MaxRetryDelay
// 	// 5. Both o.RetryDelay and o.MaxRetryDelay must be 0 or neither can be 0

// 	if len(o.StatusCodes) == 0 {
// 		o.StatusCodes = StatusCodesForRetry[:]
// 	}

// 	IfDefault := func(current *time.Duration, desired time.Duration) {
// 		if *current == time.Duration(0) {
// 			*current = desired
// 		}
// 	}

// 	// Set defaults if unspecified
// 	if o.MaxTries == 0 {
// 		o.MaxTries = 8
// 	}
// 	switch o.Policy {
// 	default:
// 		fallthrough
// 	case azcore.RetryPolicyExponential:
// 		IfDefault(&o.TryTimeout, 1*time.Minute)
// 		IfDefault(&o.RetryDelay, 4*time.Second)
// 		IfDefault(&o.MaxRetryDelay, 120*time.Second)

// 	case azcore.RetryPolicyFixed:
// 		IfDefault(&o.TryTimeout, 1*time.Minute)
// 		IfDefault(&o.RetryDelay, 30*time.Second)
// 		IfDefault(&o.MaxRetryDelay, 120*time.Second)
// 	}
// 	return o
// }

// func (o RetryOptions) calcDelay(try int32) time.Duration { // try is >=1; never 0
// 	pow := func(number int64, exponent int32) int64 { // pow is nested helper function
// 		var result int64 = 1
// 		for n := int32(0); n < exponent; n++ {
// 			result *= number
// 		}
// 		return result
// 	}

// 	delay := time.Duration(0)
// 	switch o.Policy {
// 	default:
// 		fallthrough
// 	case RetryPolicyExponential:
// 		delay = time.Duration(pow(2, try)-1) * o.RetryDelay

// 	case RetryPolicyFixed:
// 		delay = o.RetryDelay
// 	}

// 	// Introduce some jitter:  [0.0, 1.0) / 2 = [0.0, 0.5) + 0.8 = [0.8, 1.3)
// 	delay = time.Duration(delay.Seconds() * (rand.Float64()/2 + 0.8) * float64(time.Second)) // NOTE: We want math/rand; not crypto/rand
// 	if delay > o.MaxRetryDelay {
// 		delay = o.MaxRetryDelay
// 	}
// 	return delay
// }

// // NewRetryPolicy creates a policy object configured using the specified options.
// func NewRetryPolicy(o RetryOptions) azcore.Policy {
// 	return &retryPolicy{options: o.defaults()} // Force defaults to be calculated
// }

// func NewMSIRetryPolicy(o RetryOptions) azcore.Policy {
// 	return &retryPolicy{options: o.msiDefaults()}
// }

// type retryPolicy struct {
// 	options RetryOptions
// }

// func (p *retryPolicy) Do(ctx context.Context, req *azcore.Request) (resp *azcore.Response, err error) {
// 	// Exponential retry algorithm: ((2 ^ attempt) - 1) * delay * random(0.8, 1.2)
// 	// When to retry: connection failure or temporary/timeout.
// 	defer req.Close()
// 	for try := int32(1); try <= p.options.MaxTries; try++ {
// 		resp = nil // reset
// 		logf("\n=====> Try=%d\n", try)

// 		// For each try, seek to the beginning of the Body stream. We do this even for the 1st try because
// 		// the stream may not be at offset 0 when we first get it and we want the same behavior for the
// 		// 1st try as for additional tries.
// 		err = req.RewindBody()
// 		if err != nil {
// 			return
// 		}

// 		// Set the time for this particular retry operation and then Do the operation.
// 		tryCtx, tryCancel := context.WithTimeout(ctx, p.options.TryTimeout)
// 		resp, err = req.Do(tryCtx) // Make the request
// 		tryCancel()
// 		logf("Err=%v, response=%v\n", err, resp)

// 		// if there is no error and the response code isn't in the list of retry codes then we're done
// 		// TODO: if this is a failure to get an access token don't retry
// 		if (err == nil && !hasStatusCode(resp, p.options.StatusCodes...)) || errors.Is(err, &CredentialUnavailableError{}) || errors.Is(err, &AuthenticationFailedError{}) {
// 			return
// 		}

// 		// drain before retrying so nothing is leaked
// 		resp.Drain()

// 		// use the delay from retry-after if available
// 		delay, ok := resp.RetryAfter()
// 		if !ok {
// 			delay = p.options.calcDelay(try)
// 		}
// 		logf("Try=%d, Delay=%v\n", try, delay)
// 		select {
// 		case <-time.After(delay):
// 			// no-op
// 		case <-ctx.Done():
// 			err = ctx.Err()
// 			return
// 		}
// 	}
// 	return // Not retryable or too many retries; return the last response/error
// }

// // According to https://github.com/golang/go/wiki/CompilerOptimizations, the compiler will inline this method and hopefully optimize all calls to it away
// var logf = func(format string, a ...interface{}) {}
