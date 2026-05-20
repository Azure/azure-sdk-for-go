// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

type clientRetryPolicy struct {
	gem *globalEndpointManager
}

// Retry context for the request
type retryContext struct {
	useWriteEndpoint       bool
	retryCount             int
	sessionRetryCount      int
	preferredLocationIndex int
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
}

const maxRetryCount = 120
const defaultBackoff = 1

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
		resolvedEndpoint := p.gem.ResolveServiceEndpoint(retryContext.retryCount, o.resourceType, o.isWriteOperation, retryContext.useWriteEndpoint)
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
			// with our retries.
			if ctxErr := req.Raw().Context().Err(); ctxErr != nil {
				return nil, errorinfo.NonRetriableError(err)
			}
			kind := classifyNetworkError(err)
			if kind != connectionErrorNone {
				shouldRetry, errRetry := p.attemptRetryOnNetworkError(req, kind, o.isWriteOperation, &retryContext)
				if errRetry != nil {
					return nil, errRetry
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
					return nil, err
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			}
			err = req.RewindBody()
			if err != nil {
				return response, err
			}
			retryContext.retryCount += 1
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
		(status == http.StatusRequestTimeout) {
		return true
	}
	return false
}

// attemptRetryOnNetworkError decides how to respond to a transport-level
// failure. The first maxSameRegionConnectionRetries attempts always retry
// against the same region (the currently-resolved endpoint), without
// touching the location cache. Once that budget is exhausted, exactly one
// cross-region failover is attempted, subject to write-safety rules:
//   - reads always fail over;
//   - writes only fail over when the error is classified as
//     connectionErrorNotSent (i.e. we are sure the request never reached
//     the service). Writes on ambiguous errors stop retrying to avoid
//     duplicate side-effects.
//
// After the single cross-region failover, any further connection error
// stops retrying — the policy does not chain failovers across regions.
func (p *clientRetryPolicy) attemptRetryOnNetworkError(req *policy.Request, kind connectionErrorKind, isWriteOperation bool, retryContext *retryContext) (bool, error) {
	if retryContext.retryCount > maxRetryCount {
		return false, nil
	}

	// While still on the original region, allow the same-region budget.
	if !retryContext.crossRegionFailoverDone && retryContext.sameRegionRetryCount < maxSameRegionConnectionRetries {
		retryContext.sameRegionRetryCount += 1
		time.Sleep(defaultBackoff * time.Second)
		return true, nil
	}

	// We've either exhausted the same-region budget or already failed
	// over once. We only ever perform a single cross-region failover
	// from this policy; further connection failures bubble up to the
	// caller.
	if retryContext.crossRegionFailoverDone {
		return false, nil
	}
	if !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}
	if isWriteOperation && kind != connectionErrorNotSent {
		// Ambiguous failure for a write: we cannot safely retry on
		// another region without risking a duplicate.
		return false, nil
	}

	err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
	if err != nil {
		return false, err
	}
	if isWriteOperation {
		if err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL); err != nil {
			return false, err
		}
	}
	if err := p.gem.Update(req.Raw().Context(), false); err != nil {
		return false, err
	}

	retryContext.sameRegionRetryCount = 0
	retryContext.retryCount += 1
	retryContext.crossRegionFailoverDone = true
	time.Sleep(defaultBackoff * time.Second)
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnEndpointFailure(req *policy.Request, isWriteOperation bool, retryContext *retryContext) (bool, error) {
	if (retryContext.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}
	if isWriteOperation {
		err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
		if err != nil {
			return false, err
		}
	} else {
		err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
		if err != nil {
			return false, err
		}
	}

	err := p.gem.Update(req.Raw().Context(), isWriteOperation)
	if err != nil {
		return false, err
	}

	time.Sleep(defaultBackoff * time.Second)
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnSessionUnavailable(isWriteOperation bool, retryContext *retryContext) bool {
	if p.gem.CanUseMultipleWriteLocations() {
		endpoints := p.gem.locationCache.locationInfo.availReadLocations
		if isWriteOperation {
			endpoints = p.gem.locationCache.locationInfo.availWriteLocations
		}
		if retryContext.sessionRetryCount >= len(endpoints) {
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

	if err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL); err != nil {
		return false, err
	}
	if err := p.gem.Update(req.Raw().Context(), false); err != nil {
		return false, err
	}
	retryContext.requestTimeoutRetryDone = true
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

	// Failures during dial / TCP connect mean the request bytes were
	// never put on the wire.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Op == "dial" {
			return connectionErrorNotSent
		}
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

	// OS-level "connection could not be established" signals.
	switch {
	case errors.Is(err, syscall.ECONNREFUSED),
		errors.Is(err, syscall.EHOSTUNREACH),
		errors.Is(err, syscall.ENETUNREACH),
		errors.Is(err, syscall.ENETDOWN):
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

	// Other net.Error timeouts / generic OpErrors are ambiguous.
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return connectionErrorAmbiguous
	}
	if opErr != nil {
		return connectionErrorAmbiguous
	}

	return connectionErrorNone
}
