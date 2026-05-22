// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"runtime/debug"
	"sync/atomic"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

type globalEndpointManagerPolicy struct {
	gem *globalEndpointManager
	// asyncRefreshPending gates the spawn of the async refresh goroutine.
	// Without this gate every request that observes ShouldRefresh()==true
	// (potentially thousands during a burst that arrives right as the
	// throttle expires) would spawn its own goroutine, each of which would
	// then queue as a waiter inside gem.Update. The singleflight in Update
	// collapses them to one HTTP call, but the goroutine + select overhead
	// is wasted. CAS this to true before spawning; the goroutine clears it
	// on exit.
	asyncRefreshPending atomic.Bool
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Synchronous bootstrap: while the GEM has never been successfully
	// populated, every request synchronously calls Update. Concurrent
	// callers are coalesced inside gem.Update via the single-in-flight
	// pattern, so at most one HTTP call is in flight. If the call fails
	// the throttle in gem.Update (lastAttemptTime) ensures subsequent
	// bootstrap retries respect refreshTimeInterval, preventing a
	// failure storm.
	var err error
	if !p.gem.populated() {
		// Use the same context, but without the cancellation signal.
		// We DO want to preserve things like context values, but the GEM
		// update needs to complete fully, even if the user cancels the
		// triggering request.
		err = p.gem.Update(context.WithoutCancel(req.Raw().Context()), false)
	}
	if p.gem.ShouldRefresh() && p.asyncRefreshPending.CompareAndSwap(false, true) {
		// Capture the context before launching the goroutine so we do not
		// depend on req's lifetime after the policy returns.
		refreshCtx := context.WithoutCancel(req.Raw().Context())
		go func() {
			defer p.asyncRefreshPending.Store(false)
			// gem.Update's panic-safe defer re-panics after cleanup. We
			// recover here so a panic in the GEM pipeline does not bring
			// down the host process via this detached goroutine. The
			// recovered value is logged (rather than silently dropped) so
			// production crashes remain triageable.
			defer func() {
				if r := recover(); r != nil {
					log.Writef(azlog.EventResponse,
						"panic in azcosmos GEM async refresh: %v\n%s",
						r, debug.Stack())
				}
			}()
			// Log refresh failures so a chronically failing GEM is
			// observable. Without this, post-bootstrap topology drift is
			// silent: the data plane keeps routing to the cached topology
			// and callers see no signal until requests start failing.
			if err := p.gem.Update(refreshCtx, false); err != nil {
				log.Writef(azlog.EventResponse,
					"azcosmos GEM async refresh failed: %v", err)
			}
		}()
	}
	if err != nil {
		return nil, err
	}
	return req.Next()
}
