// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore azlog Writef Retriable unrecovered

package azcosmos

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"syscall"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

type clientRetryPolicy struct {
	gem *globalEndpointManager
	// asyncRefreshState tracks the in-flight goroutine spawned by
	// asyncForceRefreshGEM (Idle/Pending/Failed). See its doc comment.
	asyncRefreshState atomic.Int32
	// lastForcedRefreshUnixNano is the completion time of the most
	// recent asyncForceRefreshGEM. Read by staleForcedRefresh to
	// rate-limit repeat refreshes against the same endpoint.
	lastForcedRefreshUnixNano atomic.Int64
}

const (
	asyncRefreshIdle    int32 = 0
	asyncRefreshPending int32 = 1
	asyncRefreshFailed  int32 = 2

	// forcedRefreshMinInterval rate-limits repeat forced refreshes
	// against an already-unavailable endpoint. Must be >=
	// defaultBackoff*time.Second so a tight 403 loop cannot bypass it.
	forcedRefreshMinInterval = 2 * time.Second
)

// asyncForceRefreshGEM kicks off a forced GEM topology refresh in a
// detached goroutine. The refresh must never block a data-plane retry:
// during a regional outage the global FQDN often resolves to the same
// regional FE pool we just marked unavailable, so a synchronous Update
// can stall and prevent failover.
//
// asyncRefreshState (CAS-gated) caps in-flight refreshes at one per
// policy. We run on context.Background() so a near-expired caller
// deadline cannot abort the refresh. Panics from gem.Update are
// recovered + logged but NOT re-panicked (an unrecovered panic in a
// detached goroutine terminates the process).
//
// Returns true if a refresh was actually spawned.
func (p *clientRetryPolicy) asyncForceRefreshGEM() bool {
	for {
		state := p.asyncRefreshState.Load()
		if state == asyncRefreshPending {
			return false
		}
		if p.asyncRefreshState.CompareAndSwap(state, asyncRefreshPending) {
			break
		}
	}
	go func() {
		err := error(nil)
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic in azcosmos retry-policy async GEM refresh: %v", r)
				log.Writef(azlog.EventResponse, "%v\n%s", err, debug.Stack())
			}
			// Record completion time BEFORE flipping state so callers
			// that observe Idle also see the freshly-updated timestamp.
			p.lastForcedRefreshUnixNano.Store(time.Now().UnixNano())
			if err != nil {
				p.asyncRefreshState.Store(asyncRefreshFailed)
			} else {
				p.asyncRefreshState.Store(asyncRefreshIdle)
			}
		}()
		err = p.gem.Update(context.Background(), true)
		if err != nil {
			log.Writef(azlog.EventResponse,
				"azcosmos retry-policy async GEM refresh failed: %v", err)
		}
	}()
	return true
}

// staleForcedRefresh reports whether the rate-limit window has
// elapsed since the last completed asyncForceRefreshGEM (or no refresh
// has run yet). Used to permit follow-up refreshes for repeat 403s
// against an already-unavailable endpoint -- critical for single-master
// writes, which cannot reroute locally.
func (p *clientRetryPolicy) staleForcedRefresh() bool {
	last := p.lastForcedRefreshUnixNano.Load()
	if last == 0 {
		return true
	}
	return time.Since(time.Unix(0, last)) >= forcedRefreshMinInterval
}

// Retry context for the request
type retryContext struct {
	useWriteEndpoint       bool
	retryCount             int
	sessionRetryCount      int
	preferredLocationIndex int
	// serverErrorRetryCount tracks the number of retries attempted for a
	// transient 5xx server error (500/502/504). Only reads are retried;
	// the budget is one in-region retry followed by one cross-region retry.
	serverErrorRetryCount int
	// sameRegionRetryCount tracks the number of consecutive retries we have
	// attempted against the currently-resolved endpoint for a connection
	// error chain. It resets to 0 whenever we fail over to another region
	// or whenever an HTTP-status retry changes the endpoint.
	sameRegionRetryCount int
	// crossRegionFailoverDone is set once this request has performed its
	// single cross-region failover attempt for a connection error. After
	// this is set, further connection errors are returned to the caller
	// without additional retries.
	crossRegionFailoverDone bool
	// requestTimeoutRetryDone is set once this request has performed its
	// single cross-region retry for an HTTP 408. Only reads are retried
	// on 408; writes are returned to the caller immediately.
	requestTimeoutRetryDone bool
	// resolveFromHead is a one-shot signal to the outer Do loop to use
	// locationIndex 0 instead of retryCount. Set by retry paths that
	// demote-to-tail (MarkEndpointUnavailable* moves the bad endpoint
	// to the tail of the route list).
	resolveFromHead bool
}

const maxRetryCount = 120
const defaultBackoff = 1

// maxServerErrorRetryCount is the total number of retries attempted for a
// transient 5xx server error: one in-region retry followed by one
// cross-region retry.
const maxServerErrorRetryCount = 2

// sleepWithContext sleeps for d, but returns early with the context's error
// if ctx is cancelled or its deadline expires. Use this in retry paths so
// the policy honors caller-set context deadlines instead of consuming the
// caller's e2e timeout budget asleep.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// maxSameRegionConnectionRetries is the number of times a connection-level
// failure is retried against the current region before considering a
// cross-region failover.
const maxSameRegionConnectionRetries = 3

// connectionErrorKind classifies a transport-level error for the purposes
// of deciding whether it is safe to retry a write across regions.
type connectionErrorKind int

const (
	// connectionErrorNone indicates the error is not a transport-level
	// connection error and should not be handled by this policy's
	// network-error path.
	connectionErrorNone connectionErrorKind = iota
	// connectionErrorNotSent indicates we can prove the request never
	// reached the service (e.g. DNS failure, TCP connect refused, TLS
	// handshake failure). Safe to retry writes on another region.
	connectionErrorNotSent
	// connectionErrorAmbiguous indicates a transport failure where the
	// request may or may not have been received and processed by the
	// service. Safe to retry reads cross-region; not safe for writes.
	connectionErrorAmbiguous
)

func (p *clientRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if !req.OperationValue(&o) {
		return nil, fmt.Errorf("failed to obtain request options, please check request being sent: %s", req.Body())
	}

	retryContext := retryContext{}
	for {
		// Update the retry context with the latest retry values
		req.SetOperationValue(retryContext)
		// Consume the one-shot resolveFromHead override.
		locationIndex := retryContext.retryCount
		if retryContext.resolveFromHead {
			locationIndex = 0
			retryContext.resolveFromHead = false
		}
		resolvedEndpoint := p.gem.ResolveServiceEndpoint(locationIndex, o.resourceType, o.isWriteOperation, retryContext.useWriteEndpoint)
		regionName := p.gem.GetEndpointLocation(resolvedEndpoint)
		req.Raw().Host = resolvedEndpoint.Host
		req.Raw().URL.Host = resolvedEndpoint.Host
		attemptStartTime := time.Now().UTC()
		response, err := req.Next() // err can happen in weird scenarios (connectivity, etc)
		if err != nil {
			if state := requestDiagnosticsStateFromContext(req.Raw().Context()); state != nil && state.clientSideStats != nil {
				state.clientSideStats.recordHTTPError(attemptStartTime, req.Raw(), err, o.resourceType, regionName)
			}
			// Honor the caller's context: if their deadline expired or
			// they cancelled the request, do not consume their budget
			// with our retries. Preserve both the cancellation reason
			// (so errors.Is(returned, context.DeadlineExceeded) works)
			// and the underlying transport error for diagnostics.
			if ctxErr := req.Raw().Context().Err(); ctxErr != nil {
				return nil, errorinfo.NonRetriableError(fmt.Errorf("%w: underlying transport error: %v", ctxErr, err))
			}
			kind := classifyNetworkError(err)
			if kind != connectionErrorNone {
				shouldRetry, errRetry := p.attemptRetryOnNetworkError(req, kind, o.isWriteOperation, err, &retryContext)
				if errRetry != nil {
					return nil, errorinfo.NonRetriableError(errRetry)
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(err)
				}
				err = req.RewindBody()
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}
		if state := requestDiagnosticsStateFromContext(req.Raw().Context()); state != nil && state.clientSideStats != nil {
			state.clientSideStats.recordHTTPResponse(attemptStartTime, response, o.resourceType, regionName)
		}
		subStatus := response.Header.Get(cosmosHeaderSubstatus)
		if p.shouldRetryStatus(response.StatusCode, subStatus) {
			retryContext.useWriteEndpoint = false
			// advanceLocation gates whether the post-switch logic advances
			// retryCount (which moves the resolved endpoint to the next
			// region). An in-region 5xx retry leaves it true=>false so the
			// retry targets the same endpoint.
			advanceLocation := true
			switch response.StatusCode {
			case http.StatusForbidden:
				shouldRetry, err := p.attemptRetryOnEndpointFailure(req, o.isWriteOperation, &retryContext)
				if err != nil {
					return nil, err
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			case http.StatusNotFound:
				if !p.attemptRetryOnSessionUnavailable(o.isWriteOperation, &retryContext) {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			case http.StatusServiceUnavailable:
				if !p.attemptRetryOnServiceUnavailable(o.isWriteOperation, &retryContext) {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			case http.StatusRequestTimeout:
				shouldRetry, err := p.attemptRetryOnRequestTimeout(req, o.isWriteOperation, &retryContext)
				if err != nil {
					return nil, errorinfo.NonRetriableError(err)
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
				shouldRetry, inRegion := p.attemptRetryOnServerError(o.isWriteOperation, &retryContext)
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
				// The in-region retry targets the same endpoint, so do not
				// advance retryCount. The cross-region retry advances the
				// location via preferredLocationIndex and retryCount.
				if inRegion {
					advanceLocation = false
				}
			}
			err = req.RewindBody()
			if err != nil {
				return response, err
			}
			if advanceLocation {
				retryContext.retryCount += 1
			}
			// HTTP-status retries can change the endpoint (via retryCount
			// or preferredLocationIndex). Reset the connection-error
			// same-region budget so a fresh chain of connection errors
			// against the new endpoint gets its full set of same-region
			// retries.
			retryContext.sameRegionRetryCount = 0
			continue
		}

		return response, err
	}

}

func (p *clientRetryPolicy) shouldRetryStatus(status int, subStatus string) (shouldRetry bool) {
	if (status == http.StatusForbidden && (subStatus == subStatusWriteForbidden || subStatus == subStatusDatabaseAccountNotFound)) ||
		(status == http.StatusNotFound && subStatus == subStatusReadSessionNotAvailable) ||
		(status == http.StatusServiceUnavailable) ||
		(status == http.StatusRequestTimeout) ||
		(status == http.StatusInternalServerError) ||
		(status == http.StatusBadGateway) ||
		(status == http.StatusGatewayTimeout) {
		return true
	}
	return false
}

// attemptRetryOnNetworkError handles transport-level failures. With
// cross-region retries enabled, it allows up to maxSameRegionConnectionRetries
// against the current region (reads always, writes only when
// connectionErrorNotSent so non-idempotent mutations are never replayed),
// then performs at most one cross-region failover.
//
// MarkEndpointUnavailable* is only invoked for connectionErrorNotSent;
// a single mid-exchange failure is too weak a signal to declare the
// whole region unavailable for concurrent and future traffic. As a
// result, the ambiguous-read failover cannot rely on demote-to-tail —
// it bumps retryContext.retryCount so the Do loop resolves a different
// locationIndex on the next iteration (mirroring the 503/408 paths).
//
// transportErr is preserved alongside the caller's context error when
// the backoff is interrupted by ctx cancellation, so callers can
// errors.Is against both context.DeadlineExceeded and the transport
// error class.
func (p *clientRetryPolicy) attemptRetryOnNetworkError(req *policy.Request, kind connectionErrorKind, isWriteOperation bool, transportErr error, retryContext *retryContext) (bool, error) {
	if retryContext.retryCount > maxRetryCount {
		return false, nil
	}
	// Caller opted out of any retries.
	if !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}

	// Same-region budget: reads always, writes only when we can prove
	// the request never reached the service (avoids replaying
	// non-idempotent mutations).
	if !retryContext.crossRegionFailoverDone &&
		retryContext.sameRegionRetryCount < maxSameRegionConnectionRetries &&
		(!isWriteOperation || kind == connectionErrorNotSent) {
		retryContext.sameRegionRetryCount += 1
		if sleepErr := sleepWithContext(req.Raw().Context(), defaultBackoff*time.Second); sleepErr != nil {
			return false, fmt.Errorf("%w: underlying transport error: %v", sleepErr, transportErr)
		}
		return true, nil
	}

	// At most one cross-region failover per request.
	if retryContext.crossRegionFailoverDone {
		return false, nil
	}

	canCrossRegionWrite := !isWriteOperation || p.gem.CanUseMultipleWriteLocations()
	if isWriteOperation && (kind != connectionErrorNotSent || !canCrossRegionWrite) {
		// Ambiguous write OR single-master write: cannot safely retry.
		// Only mark when NotSent (single-master write); ambiguous is
		// too weak a signal. No forced gem.Update: invalidate() inside
		// MarkEndpointUnavailable* arms the next non-force Update.
		if kind == connectionErrorNotSent {
			if _, err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL); err != nil {
				return false, err
			}
			if canCrossRegionWrite {
				if _, err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL); err != nil {
					return false, err
				}
			}
		}
		return false, nil
	}

	// Cross-region failover: reads (any kind) or NotSent multi-master
	// writes (ambiguous writes are gated above).
	if kind == connectionErrorNotSent {
		// Mark + demote-to-tail; resolveFromHead pins the next resolve
		// to index 0, which now points at the failover region.
		if _, err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL); err != nil {
			return false, err
		}
		if isWriteOperation {
			if _, err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL); err != nil {
				return false, err
			}
		}
		retryContext.resolveFromHead = true
	} else {
		// Ambiguous read: don't mark. Route to the next region by
		// bumping retryCount (Do loop uses it as locationIndex when
		// resolveFromHead is not set). Mirrors the 503/408 paths.
		retryContext.retryCount += 1
	}

	retryContext.sameRegionRetryCount = 0
	retryContext.crossRegionFailoverDone = true
	if sleepErr := sleepWithContext(req.Raw().Context(), defaultBackoff*time.Second); sleepErr != nil {
		return false, fmt.Errorf("%w: underlying transport error: %v", sleepErr, transportErr)
	}
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnEndpointFailure(req *policy.Request, isWriteOperation bool, retryContext *retryContext) (bool, error) {
	if (retryContext.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}
	var wasAlreadyUnavailable bool
	var err error
	if isWriteOperation {
		wasAlreadyUnavailable, err = p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
		if err != nil {
			return false, err
		}
	} else {
		wasAlreadyUnavailable, err = p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
		if err != nil {
			return false, err
		}
	}

	// Kick off a forced async refresh on:
	//   - a NEW unavailability event for this endpoint (first
	//     transition; we always want fresh topology after a brand-new
	//     mark), OR
	//   - a repeat mark when no refresh is currently in flight AND
	//     the last completed forced refresh is older than
	//     forcedRefreshMinInterval. This single condition covers
	//     both recovery from a successful-but-stale prior refresh
	//     (single-master writes can't reroute locally) and recovery
	//     from a failed prior refresh (metadata endpoint was
	//     transiently unhealthy) without storming GetDatabaseAccount
	//     when the metadata endpoint is sustained-unhealthy.
	//
	// MarkEndpointUnavailable* already calls invalidate() on the first
	// transition, so the next non-force Update will refresh anyway --
	// but for single-master writes the local route list cannot reroute
	// around the bad write endpoint, so without these additional
	// forced refreshes the client could be stuck on the failed write
	// region for refreshTimeInterval (default 5 min).
	//
	// Fire-and-forget: we do NOT block the retry on its outcome.
	// MarkEndpointUnavailable* has already invalidated the GEM cache
	// and demoted the bad endpoint locally, so the next
	// ResolveServiceEndpoint will pick the failover region (in
	// multi-region scenarios) whether or not the metadata refresh
	// succeeds. Blocking here would surface a transient metadata
	// failure to the caller and skip the very cross-region retry this
	// function is supposed to perform.
	state := p.asyncRefreshState.Load()
	shouldForceRefresh := !wasAlreadyUnavailable ||
		(state != asyncRefreshPending && p.staleForcedRefresh())
	if shouldForceRefresh {
		p.asyncForceRefreshGEM()
	}

	// Force the next resolve to use locationIndex 0. Without this, the
	// outer Do() loop bumps retryCount += 1 after we return true, which
	// for a two-region account turns readEndpoints[1 % 2] back into the
	// just-marked unhealthy endpoint that MarkEndpointUnavailable*
	// demoted to the tail. resolveFromHead is a one-shot consumed by
	// the outer loop's ResolveServiceEndpoint call.
	retryContext.resolveFromHead = true

	if sleepErr := sleepWithContext(req.Raw().Context(), defaultBackoff*time.Second); sleepErr != nil {
		return false, sleepErr
	}
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnSessionUnavailable(isWriteOperation bool, retryContext *retryContext) bool {
	// Snapshot multi-write capability AND the relevant slice length
	// under a single RLock. The async refresh paths (in this file and
	// in globalEndpointManagerPolicy) can call locationCache.update
	// concurrently, which rewrites enableMultipleWriteLocations and
	// availRead/WriteLocations under mapMutex.Lock(). Sampling these
	// across two separate lock acquisitions can yield a mixed snapshot
	// (multi-write decision from before a refresh + slice length from
	// after it, or vice versa), causing the wrong branch to be taken.
	multiWrite, readN, writeN := p.gem.locationCache.sessionRetrySnapshot()
	if multiWrite {
		n := readN
		if isWriteOperation {
			n = writeN
		}
		if retryContext.sessionRetryCount >= n {
			return false
		}
	} else {
		if retryContext.sessionRetryCount > 0 {
			return false
		}
		retryContext.useWriteEndpoint = true
	}
	retryContext.sessionRetryCount += 1
	return true
}

func (p *clientRetryPolicy) attemptRetryOnServiceUnavailable(isWriteOperation bool, retryContext *retryContext) bool {
	if !p.gem.locationCache.enableCrossRegionRetries || retryContext.preferredLocationIndex >= len(p.gem.preferredLocations) {
		return false
	}
	if isWriteOperation && !p.gem.CanUseMultipleWriteLocations() {
		return false
	}
	retryContext.preferredLocationIndex += 1
	return true
}

// attemptRetryOnServerError applies the 5xx retry policy for transient server
// errors (500 Internal Server Error, 502 Bad Gateway, 504 Gateway Timeout).
// Consistent with the other Cosmos SDKs (.NET, Python, Java), only read
// operations are retried. The retry budget is one in-region retry followed by
// one cross-region retry, after which the error is surfaced to the caller. The
// cross-region retry is only attempted when cross-region retries are enabled, a
// preferred location is available to fail over to, and the location cache has
// resolved more than one read endpoint -- otherwise the "cross-region" retry
// would just hit the same endpoint as the in-region retry. The returned
// inRegion flag tells the caller whether to keep targeting the current endpoint
// (true) or to advance to the next preferred region (false).
func (p *clientRetryPolicy) attemptRetryOnServerError(isWriteOperation bool, retryContext *retryContext) (shouldRetry bool, inRegion bool) {
	if isWriteOperation {
		return false, false
	}
	if retryContext.serverErrorRetryCount >= maxServerErrorRetryCount {
		return false, false
	}
	if retryContext.serverErrorRetryCount == 0 {
		retryContext.serverErrorRetryCount += 1
		return true, true
	}
	if !p.gem.locationCache.enableCrossRegionRetries || retryContext.preferredLocationIndex >= len(p.gem.preferredLocations) {
		return false, false
	}
	if p.gem.locationCache.readEndpointCount() <= 1 {
		return false, false
	}
	retryContext.serverErrorRetryCount += 1
	retryContext.preferredLocationIndex += 1
	return true, false
}

// attemptRetryOnRequestTimeout handles an HTTP 408 from the service. A
// 408 is ambiguous from a write-safety standpoint (the request may or
// may not have been processed before the server timed out), so only
// reads are retried — and at most once, against another region. Writes
// are returned to the caller immediately so a duplicate write cannot
// occur.
func (p *clientRetryPolicy) attemptRetryOnRequestTimeout(req *policy.Request, isWriteOperation bool, retryContext *retryContext) (bool, error) {
	if isWriteOperation {
		return false, nil
	}
	if !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}
	if retryContext.requestTimeoutRetryDone {
		return false, nil
	}

	retryContext.requestTimeoutRetryDone = true
	// Preserve the caller's cancellation cause if their context fires
	// during the backoff so errors.Is(returned, context.DeadlineExceeded)
	// works at the call site. There is no transport error to compose
	// with here (the 408 path is reached from a successful HTTP
	// exchange), so the bare sleep error is sufficient.
	if sleepErr := sleepWithContext(req.Raw().Context(), defaultBackoff*time.Second); sleepErr != nil {
		return false, sleepErr
	}
	return true, nil
}

// classifyNetworkError categorizes a transport-level error so the retry
// policy can decide whether a write is safe to retry on another region.
//
//   - connectionErrorNotSent  : we are sure the request never reached the
//     service (DNS failure, TCP connect refused/unreachable, TLS handshake
//     failure, any failure during the dial phase).
//   - connectionErrorAmbiguous: a transport failure that may have occurred
//     after the request was placed on the wire (EOF, connection reset,
//     broken pipe, transport-level timeouts).
//   - connectionErrorNone     : not a transport-level connection error.
//
// The syscall constants (ECONNREFUSED, ECONNRESET, etc.) used below are
// normalized by the Go runtime across Unix and Windows, so this
// classifier is portable without per-OS build tags.
func classifyNetworkError(err error) connectionErrorKind {
	if err == nil {
		return connectionErrorNone
	}

	// Definitely not sent: DNS resolution failures occur before any
	// connection is established.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return connectionErrorNotSent
	}

	// TLS handshake failures happen before any request bytes are flushed.
	var tlsRecordErr tls.RecordHeaderError
	if errors.As(err, &tlsRecordErr) {
		return connectionErrorNotSent
	}
	var certErr *tls.CertificateVerificationError
	if errors.As(err, &certErr) {
		return connectionErrorNotSent
	}

	// OS-level "connection could not be established" signals. ETIMEDOUT
	// here covers TCP connect timeouts (kernel SYN timeout) — we treat
	// those as not-sent because no application bytes were ever written.
	switch {
	case errors.Is(err, syscall.ECONNREFUSED),
		errors.Is(err, syscall.EHOSTUNREACH),
		errors.Is(err, syscall.ENETUNREACH),
		errors.Is(err, syscall.ENETDOWN),
		errors.Is(err, syscall.ETIMEDOUT):
		return connectionErrorNotSent
	}

	// Ambiguous: the connection was up but failed mid-exchange.
	switch {
	case errors.Is(err, io.EOF),
		errors.Is(err, io.ErrUnexpectedEOF),
		errors.Is(err, syscall.ECONNRESET),
		errors.Is(err, syscall.EPIPE):
		return connectionErrorAmbiguous
	}

	// Transport-level deadlines (e.g. http.Transport.ResponseHeaderTimeout)
	// surface as context.DeadlineExceeded but without the caller's context
	// being done; the caller-context check is performed by the retry loop
	// before this function is called.
	if errors.Is(err, context.DeadlineExceeded) {
		return connectionErrorAmbiguous
	}

	// net.OpError covers dial / read / write / accept failures. A "dial"
	// failure proves the request never left this host. Any other
	// OpError happened after the connection was established and is
	// therefore ambiguous.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Op == "dial" {
			return connectionErrorNotSent
		}
		return connectionErrorAmbiguous
	}

	// Other net.Error timeouts (e.g. http.Client.Timeout) are ambiguous.
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return connectionErrorAmbiguous
	}

	return connectionErrorNone
}
