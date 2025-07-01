// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestChangeFeedOptionsToHeaders(t *testing.T) {
	options := &ChangeFeedOptions{}
	headers := options.toHeaders()
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	h := *headers
	if h[cosmosHeaderMaxItemCount] != "-1" {
		t.Errorf("Expected default MaxItemCount to be -1, got %v", h[cosmosHeaderMaxItemCount])
	}
	if h[cosmosHeaderChangeFeed] != cosmosHeaderValuesChangeFeed {
		t.Errorf("Expected default AIM to be %v, got %v", cosmosHeaderValuesChangeFeed, h[cosmosHeaderChangeFeed])
	}

	options.MaxItemCount = 10
	options.AIM = "Incremental Feed"
	headers = options.toHeaders()
	h = *headers
	if h[cosmosHeaderMaxItemCount] != "10" {
		t.Errorf("Expected MaxItemCount to be 10, got %v", h[cosmosHeaderMaxItemCount])
	}
	if h[cosmosHeaderChangeFeed] != "Incremental Feed" {
		t.Errorf("Expected AIM to be Incremental Feed, got %v", h[cosmosHeaderChangeFeed])
	}

	etag := azcore.ETag("test-etag")
	options.IfNoneMatch = &etag
	headers = options.toHeaders()
	h = *headers
	if h[headerIfNoneMatch] != `"test-etag"` {
		t.Errorf("Expected IfNoneMatch to be \"test-etag\", got %v", h[headerIfNoneMatch])
	}

	now := time.Now().UTC()
	options.IfModifiedSince = &now
	headers = options.toHeaders()
	h = *headers
	expectedIfModifiedSince := strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
	if h[cosmosHeaderIfModifiedSince] != expectedIfModifiedSince {
		t.Errorf("Expected IfModifiedSince to be %v, got %v", expectedIfModifiedSince, h[cosmosHeaderIfModifiedSince])
	}

	pk := NewPartitionKeyString("pkvalue")
	options.PartitionKey = &pk
	headers = options.toHeaders()
	h = *headers
	pkJSON, _ := pk.toJsonString()
	if h[cosmosHeaderPartitionKey] != string(pkJSON) {
		t.Errorf("Expected PartitionKey to be %v, got %v", string(pkJSON), h[cosmosHeaderPartitionKey])
	}

	pkrid := "range-id"
	options.PartitionKeyRangeID = &pkrid
	headers = options.toHeaders()
	h = *headers
	if h[cosmosHeaderPartitionKeyRangeId] != "range-id" {
		t.Errorf("Expected PartitionKeyRangeID to be range-id, got %v", h[cosmosHeaderPartitionKeyRangeId])
	}
}
