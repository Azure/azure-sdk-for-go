// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestTransactionalBatchOptionsToHeaders(t *testing.T) {
	options := &TransactionalBatchOptions{}
	options.ConsistencyLevel = ConsistencyLevelSession.ToPtr()
	options.SessionToken = "sessionToken"
	etagValue := azcore.ETag("someEtag")
	options.IfMatchEtag = &etagValue
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
	if headers[headerIfMatch] != string(*options.IfMatchEtag) {
		t.Errorf("IfMatchEtag should be someEtag but got %v", headers[headerIfMatch])
	}
}
