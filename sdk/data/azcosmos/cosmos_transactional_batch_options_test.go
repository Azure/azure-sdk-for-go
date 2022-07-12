// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestTransactionalBatchOptionsToHeaders(t *testing.T) {
	options := &TransactionalBatchOptions{}
	options.ConsistencyLevel = ConsistencyLevelSession.ToPtr()
	options.SessionToken = "sessionToken"
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
	if headers[cosmosHeaderIsBatchAtomic] != "True" {
		t.Fatal("IsBatchAtomic should be true")
	}
	if headers[cosmosHeaderIsBatchRequest] != "True" {
		t.Fatal("IsBatchRequest should be true")
	}
	if headers[cosmosHeaderIsBatchOrdered] != "True" {
		t.Fatal("IsBatchOrdered should be true")
	}
}
