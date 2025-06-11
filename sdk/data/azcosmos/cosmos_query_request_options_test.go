// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strconv"
	"testing"
	"time"
)

func TestQueryRequestOptionsToHeaders(t *testing.T) {
	options := &QueryOptions{}
	options.ConsistencyLevel = ConsistencyLevelSession.ToPtr()
	sessionToken := "sessionToken"
	options.SessionToken = &sessionToken
	options.PageSizeHint = 20
	options.EnableScanInQuery = true
	options.ResponseContinuationTokenLimitInKB = 100
	options.PopulateIndexMetrics = true
	continuation := "continuationToken"
	options.ContinuationToken = &continuation
	maxIntegratedCacheStalenessDuration := time.Duration(5 * time.Minute)
	options.DedicatedGatewayRequestOptions = &DedicatedGatewayRequestOptions{
		MaxIntegratedCacheStaleness: &maxIntegratedCacheStalenessDuration,
	}
	options.DedicatedGatewayRequestOptions.BypassIntegratedCache = true
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if headers[cosmosHeaderConsistencyLevel] != "Session" {
		t.Errorf("ConsistencyLevel should be Session but got %v", headers[cosmosHeaderConsistencyLevel])
	}
	if headers[cosmosHeaderSessionToken] != "sessionToken" {
		t.Errorf("SessionToken should be sessionToken but got %v", headers[cosmosHeaderSessionToken])
	}
	if headers[cosmosHeaderMaxItemCount] != "20" {
		t.Errorf("PageSizeHint should be 20 but got %v", headers[cosmosHeaderMaxItemCount])
	}
	if headers[cosmosHeaderEnableScanInQuery] != "true" {
		t.Errorf("EnableScanInQuery should be true but got %v", headers[cosmosHeaderEnableScanInQuery])
	}
	if headers[cosmosHeaderResponseContinuationTokenLimitInKb] != "100" {
		t.Errorf("ResponseContinuationTokenLimitInKb should be 100 but got %v", headers[cosmosHeaderResponseContinuationTokenLimitInKb])
	}
	if headers[cosmosHeaderPopulateIndexMetrics] != "true" {
		t.Errorf("PopulateIndexMetrics should be true but got %v", headers[cosmosHeaderPopulateIndexMetrics])
	}
	if headers[cosmosHeaderContinuationToken] != "continuationToken" {
		t.Errorf("ContinuationToken should be continuationToken but got %v", headers[cosmosHeaderContinuationToken])
	}
	if headers[cosmosHeaderPopulateQueryMetrics] != "true" {
		t.Errorf("PopulateQueryMetrics should be true but got %v", headers[cosmosHeaderPopulateQueryMetrics])
	}
	if headers[headerDedicatedGatewayMaxAge] != strconv.FormatInt(300000, 10) {
		t.Errorf("headerDedicatedGatewayMaxAge should be 300000 but got %v", headers[headerDedicatedGatewayMaxAge])
	}
	if headers[headerDedicatedGatewayBypassCache] != "true" {
		t.Errorf("headerDedicatedGatewayBypassCache should be true but got %v", headers[headerDedicatedGatewayBypassCache])
	}
}

func TestQueryRequestOptionsToHeaders_bypassIntegratedCacheNotSet(t *testing.T) {
	options := &QueryOptions{}
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if _, exists := headers[headerDedicatedGatewayBypassCache]; exists {
		t.Errorf("headerDedicatedGatewayBypassCache should not exist when BypassIntegratedCache is not set")
	}
}
