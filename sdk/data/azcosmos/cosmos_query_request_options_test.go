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
	options.MaxItemCount = 20
	options.EnableScanInQuery = true
	options.ResponseContinuationTokenLimitInKb = 100
	options.PopulateIndexMetrics = true
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
		t.Errorf("MaxItemCount should be 20 but got %v", headers[cosmosHeaderMaxItemCount])
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
}
