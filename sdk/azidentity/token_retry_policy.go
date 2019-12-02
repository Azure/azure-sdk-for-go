// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewMSIRetryPolicy checks for authentication errors and stops retrying in case of
// a failure in trying to get a token
func NewMSIRetryPolicy(o azcore.RetryOptions) azcore.Policy {
	return &msiRetryPolicy{options: o}
}

type msiRetryPolicy struct {
	options azcore.RetryOptions
}

func (p *msiRetryPolicy) Do(ctx context.Context, req *azcore.Request) (resp *azcore.Response, err error) {
	retries := []int{
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503  TODO look into this
		http.StatusGatewayTimeout,      // 504
	}
	// extra retry status codes specific to IMDS
	retries = append(retries,
		http.StatusNotFound,
		http.StatusGone,
		// all remaining 5xx
		http.StatusNotImplemented,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired)

	// see https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/how-to-use-vm-token#retry-guidance

	const maxDelay time.Duration = 60 * time.Second

	attempt := int32(0)
	delay := time.Duration(0)

	for attempt < p.options.MaxTries {
		resp, err = req.Do(ctx)
		// we want to retry if err is not nil or the status code is in the list of retry codes
		// we will fail if we receive a token credential related error
		var credUnavailable *CredentialUnavailableError
		var credFailure *AuthenticationFailedError
		if (err == nil && !resp.HasStatusCode(retries...)) || errors.As(err, &credUnavailable) || errors.As(err, &credFailure) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		// perform exponential backoff with a cap.
		// must increment attempt before calculating delay.
		attempt++
		// the base value of 2 is the "delta backoff" as specified in the guidance doc
		delay += (time.Duration(math.Pow(2, float64(attempt))) * time.Second)
		if delay > maxDelay {
			delay = maxDelay
		}

		select {
		case <-time.After(delay):
			// intentionally left blank
		case <-req.Context().Done():
			err = req.Context().Err()
			return
		}
	}
	return
}
