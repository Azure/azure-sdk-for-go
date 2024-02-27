// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type clientRetryPolicy struct {
	gem                    *globalEndpointManager
	retryCount             int
	sessionRetryCount      int
	preferredLocationIndex int
}

const maxRetryCount = 120
const defaultBackoff = 1

func (p *clientRetryPolicy) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if !req.OperationValue(&o) {
		return nil, fmt.Errorf("Failed to obtain request options, please check request being sent: %s", req.Body())
	}
	resolvedEndpoint := p.gem.ResolveServiceEndpoint(p.retryCount, o.isWriteOperation)
	req.Raw().Host = resolvedEndpoint.Host
	req.Raw().URL.Host = resolvedEndpoint.Host
	response, err := req.Next() // err can happen in weird scenarios (connectivity, etc) - need to test
	subStatus := response.Header.Get(cosmosHeaderSubstatus)
	fmt.Println(response.StatusCode)
	if p.shouldRetryStatus(response.StatusCode, subStatus) {
		p.retryCount = 0
		p.sessionRetryCount = 0
		for {
			subStatus = response.Header.Get(cosmosHeaderSubstatus)
			if p.shouldRetryStatus(response.StatusCode, subStatus) {
				fmt.Println("Policy TIME")
				if response.StatusCode == statusForbidden {
					if !p.attemptRetryOnEndpointFailure(req, o.isWriteOperation) {
						break
					}
				} else if response.StatusCode == statusNotFound {
					if !p.attemptRetryOnSessionUnavailable(req, o.isWriteOperation) {
						break
					}
				} else if response.StatusCode == statusServiceUnavailable {
					if !p.attemptRetryOnServiceUnavailable(req, o.isWriteOperation) {
						break
					}
				}
				fmt.Println("bout to retry this")
			} else {
				fmt.Println("supposed to break this")
				break
			}
			err = req.RewindBody()
			if err != nil {
				return response, err
			}
			resolvedEndpoint := p.gem.ResolveServiceEndpoint(p.retryCount, o.isWriteOperation)
			req.Raw().Host = resolvedEndpoint.Host
			req.Raw().URL.Host = resolvedEndpoint.Host
			fmt.Println("retryngggg")
			response, err = req.Next()
			fmt.Println("should have retried")
		}
	}
	fmt.Println("returning response from policy")
	return response, err
}

func (p *clientRetryPolicy) shouldRetryStatus(status int, subStatus string) (shouldRetry bool) {
	if (status == statusForbidden && (subStatus == subStatusWriteForbidden || subStatus == subStatusDatabaseAccountNotFound)) ||
		(status == statusNotFound && subStatus == subStatusReadSessionNotAvailable) ||
		(status == statusServiceUnavailable) {
		return true
	}
	return false
}

func (p *clientRetryPolicy) attemptRetryOnEndpointFailure(req *policy.Request, isWriteOperation bool) bool {
	if (p.retryCount > maxRetryCount) || !p.gem.locationCache.enableCrossRegionRetries {
		return false
	}
	if isWriteOperation {
		p.gem.MarkEndpointUnavailableForWrite(*req.Raw().URL)
	} else {
		p.gem.MarkEndpointUnavailableForRead(*req.Raw().URL)
	}
	p.gem.Update(req.Raw().Context())

	p.retryCount += 1
	time.Sleep(defaultBackoff * time.Second)
	return true
}

func (p *clientRetryPolicy) attemptRetryOnSessionUnavailable(req *policy.Request, isWriteOperation bool) bool {
	if p.gem.CanUseMultipleWriteLocations() {
		endpoints := []string{}
		if isWriteOperation {
			endpoints = p.gem.locationCache.locationInfo.availWriteLocations
		} else {
			endpoints = p.gem.locationCache.locationInfo.availReadLocations
		}
		if p.sessionRetryCount >= len(endpoints) {
			return false
		}
	} else {
		if p.sessionRetryCount > 0 {
			return false
		}
	}
	p.sessionRetryCount += 1
	p.retryCount += 1
	return true
}

func (p *clientRetryPolicy) attemptRetryOnServiceUnavailable(req *policy.Request, isWriteOperation bool) bool {
	//On HTTP 503 response, if it's a read request and preferredRegions > 1,
	//retry on the next preferredRegion. If it's a write request and account is multi master
	//and preferredRegions > 1, retry on the next preferredRegion.
	if !p.gem.locationCache.enableCrossRegionRetries || p.preferredLocationIndex >= len(p.gem.preferredLocations) {
		return false
	}
	if isWriteOperation {
		if p.gem.CanUseMultipleWriteLocations() {
			locationalEndpoint := p.gem.GetPreferredLocationEndpoint(p.preferredLocationIndex, *req.Raw().URL)
			req.Raw().URL = &locationalEndpoint
		} else {
			return false
		}
	} else {
		locationalEndpoint := p.gem.GetPreferredLocationEndpoint(p.preferredLocationIndex, *req.Raw().URL)
		req.Raw().URL = &locationalEndpoint
	}
	p.preferredLocationIndex += 1
	return true
}
