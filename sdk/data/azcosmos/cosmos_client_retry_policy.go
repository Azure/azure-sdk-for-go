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
	gem                    *globalEndpointManager
	useWriteEndpoint       bool
	retryCount             int
	sessionRetryCount      int
	preferredLocationIndex int
}

const maxRetryCount = 120
const defaultBackoff = 1

func (p *clientRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.resetPolicyCounters()
	o := pipelineRequestOptions{}
	if !req.OperationValue(&o) {
		return nil, fmt.Errorf("failed to obtain request options, please check request being sent: %s", req.Body())
	}
	for{
		resolvedEndpoint := p.gem.ResolveServiceEndpoint(p.retryCount, o.isWriteOperation, p.useWriteEndpoint)
		req.Raw().Host = resolvedEndpoint.Host
		req.Raw().URL.Host = resolvedEndpoint.Host
		response, err := req.Next() // err can happen in weird scenarios (connectivity, etc)
		if err != nil {
			if (p.isNetworkConnectionError(err)) {
				shouldRetry, err := p.attemptRetryOnNetworkError(req)
				if err != nil {
					return nil, err
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
				err = req.RewindBody()
				if err != nil {
					return nil, err
				}
				p.retryCount += 1
				continue;
			}

			return nil, err
		}
		subStatus := response.Header.Get(cosmosHeaderSubstatus)
		if p.shouldRetryStatus(response.StatusCode, subStatus) {
			p.useWriteEndpoint = false
			if response.StatusCode == http.StatusForbidden {
				shouldRetry, err := p.attemptRetryOnEndpointFailure(req, o.isWriteOperation)
				if err != nil {
					return nil, err
				}
				if !shouldRetry {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			} else if response.StatusCode == http.StatusNotFound {
				if !p.attemptRetryOnSessionUnavailable(req, o.isWriteOperation) {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			} else if response.StatusCode == http.StatusServiceUnavailable {
				if !p.attemptRetryOnServiceUnavailable(req, o.isWriteOperation) {
					return nil, errorinfo.NonRetriableError(azruntime.NewResponseErrorWithErrorCode(response, response.Status))
				}
			}
			err = req.RewindBody()
			if err != nil {
				return response, err
			}
			p.retryCount += 1
			continue;
		}

		return response, err
	}
	
}

func (p *clientRetryPolicy) shouldRetryStatus(status int, subStatus string) (shouldRetry bool) {
	if (status == http.StatusForbidden && (subStatus == subStatusWriteForbidden || subStatus == subStatusDatabaseAccountNotFound)) ||
		(status == http.StatusNotFound && subStatus == subStatusReadSessionNotAvailable) ||
		(status == http.StatusServiceUnavailable) {
		return true
	}
	return false
}

func (p *clientRetryPolicy) attemptRetryOnNetworkError(req *policy.Request) (bool, error) {
	if (p.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}

	err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
	if (err != nil) {
		return false, err
	}
	err = p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
	if (err != nil) {
		return false, err
	}
	err = p.gem.Update(req.Raw().Context(), false)
	if (err != nil) {
		return false, err
	}

	time.Sleep(defaultBackoff * time.Second)
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnEndpointFailure(req *policy.Request, isWriteOperation bool) (bool, error) {
	if (p.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false, nil
	}
	if isWriteOperation {
		err := p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
		if (err != nil) {
			return false, err
		}
	} else {
		err := p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
		if (err != nil) {
			return false, err
		}
	}

	err := p.gem.Update(req.Raw().Context(), isWriteOperation)
	if (err != nil) {
		return false, err
	}

	time.Sleep(defaultBackoff * time.Second)
	return true, nil
}

func (p *clientRetryPolicy) attemptRetryOnSessionUnavailable(req *policy.Request, isWriteOperation bool) bool {
	if p.gem.CanUseMultipleWriteLocations() {
		endpoints := p.gem.locationCache.locationInfo.availReadLocations
		if isWriteOperation {
			endpoints = p.gem.locationCache.locationInfo.availWriteLocations
		}
		if p.sessionRetryCount >= len(endpoints) {
			return false
		}
	} else {
		if p.sessionRetryCount > 0 {
			return false
		}
		p.useWriteEndpoint = true
	}
	p.sessionRetryCount += 1
	return true
}

func (p *clientRetryPolicy) attemptRetryOnServiceUnavailable(req *policy.Request, isWriteOperation bool) bool {
	if !p.gem.locationCache.enableCrossRegionRetries || p.preferredLocationIndex >= len(p.gem.preferredLocations) {
		return false
	}
	if isWriteOperation && !p.gem.CanUseMultipleWriteLocations() {
		return false
	}
	p.preferredLocationIndex += 1
	return true
}

func (p *clientRetryPolicy) resetPolicyCounters() {
	p.retryCount = 0
	p.sessionRetryCount = 0
	p.preferredLocationIndex = 0
}

// isNetworkConnectionError checks if the error is related to failure to connect / resolve DNS
func (p *clientRetryPolicy) isNetworkConnectionError(err error) bool {
	var dnserror *net.DNSError 
	return errors.As(err, &dnserror)
}
