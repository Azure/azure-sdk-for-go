// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestQueryRequestOptionsToHeaders(t *testing.T) {
	options := &QueryOptions{}
	options.ConsistencyLevel = ConsistencyLevelSession.ToPtr()
	options.SessionToken = "sessionToken"
	options.PageSizeHint = 20
	options.EnableScanInQuery = true
	options.ResponseContinuationTokenLimitInKB = 100
	options.PopulateIndexMetrics = true
	options.ContinuationToken = "continuationToken"
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
}
