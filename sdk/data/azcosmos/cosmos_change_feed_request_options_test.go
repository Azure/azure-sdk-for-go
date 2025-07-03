// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strings"
	"testing"
	"time"
)

func TestChangeFeedOptionsToHeaders(t *testing.T) {
	// Test with no partition key ranges
	options := &ChangeFeedOptions{}
	headers := options.toHeaders(nil)
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

	// Test MaxItemCount
	options.MaxItemCount = 10
	headers = options.toHeaders(nil)
	h = *headers
	if h[cosmosHeaderMaxItemCount] != "10" {
		t.Errorf("Expected MaxItemCount to be 10, got %v", h[cosmosHeaderMaxItemCount])
	}
	if h[cosmosHeaderChangeFeed] != cosmosHeaderValuesChangeFeed {
		t.Errorf("Expected AIM to be Incremental Feed, got %v", h[cosmosHeaderChangeFeed])
	}

	// Test Continuation
	continuation := "test-etag"
	options.Continuation = &continuation
	headers = options.toHeaders(nil)
	h = *headers
	if h[headerIfNoneMatch] != "test-etag" {
		t.Errorf("Expected IfNoneMatch to be \"test-etag\", got %v", h[headerIfNoneMatch])
	}

	// Test ChangeFeedStartFrom (replaced IfModifiedSince)
	now := time.Now().UTC()
	options.ChangeFeedStartFrom = &now
	headers = options.toHeaders(nil)
	h = *headers
	expectedIfModifiedSince := strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
	if h[cosmosHeaderIfModifiedSince] != expectedIfModifiedSince {
		t.Errorf("Expected IfModifiedSince to be %v, got %v", expectedIfModifiedSince, h[cosmosHeaderIfModifiedSince])
	}

	// Test PartitionKey
	pk := NewPartitionKeyString("pkvalue")
	options.PartitionKey = &pk
	headers = options.toHeaders(nil)
	h = *headers
	pkJSON, _ := pk.toJsonString()
	if h[cosmosHeaderPartitionKey] != string(pkJSON) {
		t.Errorf("Expected PartitionKey to be %v, got %v", string(pkJSON), h[cosmosHeaderPartitionKey])
	}

	// Test FeedRange with partition key ranges
	feedRange := &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}
	options.FeedRange = feedRange

	// Test with matching partition key range
	partitionKeyRanges := []PartitionKeyRangeProperties{
		{
			ID:           "0",
			MinInclusive: "00",
			MaxExclusive: "FF",
		},
	}

	headers = options.toHeaders(partitionKeyRanges)
	h = *headers
	if h[headerXmsDocumentDbPartitionKeyRangeId] != "0" {
		t.Errorf("Expected partition key range ID to be 0, got %v", h[headerXmsDocumentDbPartitionKeyRangeId])
	}

	// Test FeedRange with no matching partition key range
	partitionKeyRangesNoMatch := []PartitionKeyRangeProperties{
		{
			ID:           "1",
			MinInclusive: "AA",
			MaxExclusive: "BB",
		},
	}

	headers = options.toHeaders(partitionKeyRangesNoMatch)
	h = *headers
	if _, exists := h[headerXmsDocumentDbPartitionKeyRangeId]; exists {
		t.Errorf("Expected no partition key range ID header when no match found")
	}

	// Test empty continuation
	emptyContinuation := ""
	options.Continuation = &emptyContinuation
	headers = options.toHeaders(nil)
	h = *headers
	if _, exists := h[headerIfNoneMatch]; exists {
		t.Errorf("Expected no IfNoneMatch header for empty continuation")
	}

	// Test nil continuation
	options.Continuation = nil
	headers = options.toHeaders(nil)
	h = *headers
	if _, exists := h[headerIfNoneMatch]; exists {
		t.Errorf("Expected no IfNoneMatch header for nil continuation")
	}
}

func TestChangeFeedOptionsToHeadersWithAllFields(t *testing.T) {
	// Test with all fields populated
	now := time.Now().UTC()
	pk := NewPartitionKeyString("testPK")
	continuation := "test-continuation"
	feedRange := &FeedRange{
		MinInclusive: "10",
		MaxExclusive: "20",
	}

	options := &ChangeFeedOptions{
		MaxItemCount:        25,
		ChangeFeedStartFrom: &now,
		PartitionKey:        &pk,
		FeedRange:           feedRange,
		Continuation:        &continuation,
	}

	partitionKeyRanges := []PartitionKeyRangeProperties{
		{
			ID:           "range1",
			MinInclusive: "10",
			MaxExclusive: "20",
		},
	}

	headers := options.toHeaders(partitionKeyRanges)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	h := *headers

	// Verify all headers are set correctly
	if h[cosmosHeaderMaxItemCount] != "25" {
		t.Errorf("Expected MaxItemCount to be 25, got %v", h[cosmosHeaderMaxItemCount])
	}

	expectedIfModifiedSince := strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
	if h[cosmosHeaderIfModifiedSince] != expectedIfModifiedSince {
		t.Errorf("Expected IfModifiedSince to be %v, got %v", expectedIfModifiedSince, h[cosmosHeaderIfModifiedSince])
	}

	pkJSON, _ := pk.toJsonString()
	if h[cosmosHeaderPartitionKey] != string(pkJSON) {
		t.Errorf("Expected PartitionKey to be %v, got %v", string(pkJSON), h[cosmosHeaderPartitionKey])
	}

	if h[headerXmsDocumentDbPartitionKeyRangeId] != "range1" {
		t.Errorf("Expected partition key range ID to be range1, got %v", h[headerXmsDocumentDbPartitionKeyRangeId])
	}

	if h[headerIfNoneMatch] != continuation {
		t.Errorf("Expected IfNoneMatch to be %v, got %v", continuation, h[headerIfNoneMatch])
	}

	if h[cosmosHeaderChangeFeed] != cosmosHeaderValuesChangeFeed {
		t.Errorf("Expected AIM to be %v, got %v", cosmosHeaderValuesChangeFeed, h[cosmosHeaderChangeFeed])
	}
}
