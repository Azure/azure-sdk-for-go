// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type globalEndpointManagerPolicy struct {
	gem *globalEndpointManager
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Synchronous bootstrap: while the GEM has never been successfully
	// populated, every request synchronously calls Update. Concurrent
	// callers are coalesced inside gem.Update via the single-in-flight
	// pattern, so at most one HTTP call is in flight. If the call fails
	// the throttle in gem.Update (lastAttemptTime) ensures subsequent
	// bootstrap retries respect refreshTimeInterval -- preventing the
	// failure storm described in
	// https://github.com/Azure/azure-sdk-for-go/issues/25468.
	var err error
	if !p.gem.populated() {
		// Use the same context, but without the cancellation signal.
		// We DO want to preserve things like context values, but the GEM
		// update needs to complete fully, even if the user cancels the
		// triggering request.
		err = p.gem.Update(context.WithoutCancel(req.Raw().Context()), false)
	}
	if p.gem.ShouldRefresh() {
		go func() {
			// Concurrent goroutines spawned here are coalesced inside
			// gem.Update via the single-in-flight pattern.
			_ = p.gem.Update(context.WithoutCancel(req.Raw().Context()), false)
		}()
	}
	if err != nil {
		return nil, err
	}
	return req.Next()
}
