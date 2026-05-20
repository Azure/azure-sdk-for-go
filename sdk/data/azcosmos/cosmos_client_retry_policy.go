// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"net"
	"net/http"
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
	serverErrorRetryCount  int
}

const maxRetryCount = 120
const defaultBackoff = 1
const maxServerErrorRetryCount = 2

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
			if p.isNetworkConnectionError(err) {
				shouldRetry, errRetry := p.attemptRetryOnNetworkError(req, &retryContext)
				if errRetry != nil {
					return nil, errRetry
				}
				if !shouldRetry {
					return nil, err
				}
				err = req.RewindBody()
				if err != nil {
					return nil, err
				}
				retryContext.retryCount += 1
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
			case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
				shouldRetry, inRegion := p.attemptRetryOnServerError(o.isWriteOperation, &retryContext)
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
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
			continue
		}

		return response, err
	}

}

func (p *clientRetryPolicy) shouldRetryStatus(status int, subStatus string) (shouldRetry bool) {
	if (status == http.StatusForbidden && (subStatus == subStatusWriteForbidden || subStatus == subStatusDatabaseAccountNotFound)) ||
		(status == http.StatusNotFound && subStatus == subStatusReadSessionNotAvailable) ||
		(status == http.StatusServiceUnavailable) ||
		(status == http.StatusInternalServerError) ||
		(status == http.StatusBadGateway) ||
		(status == http.StatusGatewayTimeout) {
		return true
	}
	return false
}

func (p *clientRetryPolicy) attemptRetryOnNetworkError(req *policy.Request, retryContext *retryContext) (bool, error) {
	if (retryContext.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}

	err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
	if err != nil {
		return false, err
	}
	err = p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
	if err != nil {
		return false, err
	}
	err = p.gem.Update(req.Raw().Context(), false)
	if err != nil {
		return false, err
	}

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

// attemptRetryOnServerError applies the 5xx retry policy for transient server errors
// (500 Internal Server Error, 502 Bad Gateway, 504 Gateway Timeout). Consistent with
// the other Cosmos SDKs (.NET, Python, Java), only read operations are retried. The
// retry budget is one in-region retry followed by one cross-region retry, after which
// the error is surfaced to the caller. The cross-region retry is only attempted when
// cross-region retries are enabled and a preferred location is available to fail over
// to. The returned inRegion flag tells the caller whether to keep targeting the current
// endpoint (true) or to advance to the next preferred region (false).
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
	retryContext.serverErrorRetryCount += 1
	retryContext.preferredLocationIndex += 1
	return true, false
}

// isNetworkConnectionError checks if the error is related to failure to connect / resolve DNS
func (p *clientRetryPolicy) isNetworkConnectionError(err error) bool {
	var dnserror *net.DNSError
	return errors.As(err, &dnserror)
}
