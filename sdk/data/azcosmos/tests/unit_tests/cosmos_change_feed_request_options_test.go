// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestChangeFeedOptionsToHeaders(t *testing.T) {
	options := &ChangeFeedOptions{}
	headers := options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	h := *headers
	if h[cosmosHeaderChangeFeed] != cosmosHeaderValuesChangeFeed {
		t.Errorf("Expected default AIM to be %v, got %v", cosmosHeaderValuesChangeFeed, h[cosmosHeaderChangeFeed])
	}

	options.MaxItemCount = 10
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	if h[cosmosHeaderMaxItemCount] != "10" {
		t.Errorf("Expected MaxItemCount to be 10, got %v", h[cosmosHeaderMaxItemCount])
	}

	continuation := "test-etag"
	options.Continuation = &continuation
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	if h[headerIfNoneMatch] != "test-etag" {
		t.Errorf("Expected IfNoneMatch to be \"test-etag\", got %v", h[headerIfNoneMatch])
	}

	now := time.Now().UTC()
	options.StartFrom = &now
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	expectedIfModifiedSince := now.Format(time.RFC1123)
	if h[cosmosHeaderIfModifiedSince] != expectedIfModifiedSince {
		t.Errorf("Expected IfModifiedSince to be %v, got %v", expectedIfModifiedSince, h[cosmosHeaderIfModifiedSince])
	}

	pk := NewPartitionKeyString("pkvalue")
	options.PartitionKey = &pk
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	pkJSON, _ := pk.toJsonString()
	if h[cosmosHeaderPartitionKey] != string(pkJSON) {
		t.Errorf("Expected PartitionKey to be %v, got %v", string(pkJSON), h[cosmosHeaderPartitionKey])
	}

	feedRange := &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}
	options.FeedRange = feedRange

	partitionKeyRanges := []partitionKeyRange{
		{
			ID:           "0",
			MinInclusive: "00",
			MaxExclusive: "FF",
		},
	}

	headers = options.toHeaders(partitionKeyRanges)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	if h[headerXmsDocumentDbPartitionKeyRangeId] != "0" {
		t.Errorf("Expected partition key range ID to be 0, got %v", h[headerXmsDocumentDbPartitionKeyRangeId])
	}

	partitionKeyRangesNoMatch := []partitionKeyRange{
		{
			ID:           "1",
			MinInclusive: "AA",
			MaxExclusive: "BB",
		},
	}

	headers = options.toHeaders(partitionKeyRangesNoMatch)
	if headers != nil {
		t.Errorf("Expected nil headers when no matching partition key range found")
	}

	options.FeedRange = nil

	emptyContinuation := ""
	options.Continuation = &emptyContinuation
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	if _, exists := h[headerIfNoneMatch]; exists {
		t.Errorf("Expected no IfNoneMatch header for empty continuation")
	}

	options.Continuation = nil
	headers = options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}
	h = *headers
	if _, exists := h[headerIfNoneMatch]; exists {
		t.Errorf("Expected no IfNoneMatch header for nil continuation")
	}
}

func TestChangeFeedOptionsToHeadersWithAllFields(t *testing.T) {
	now := time.Now().UTC()
	pk := NewPartitionKeyString("testPK")
	continuation := "test-continuation"
	feedRange := &FeedRange{
		MinInclusive: "10",
		MaxExclusive: "20",
	}

	options := &ChangeFeedOptions{
		MaxItemCount: 25,
		StartFrom:    &now,
		PartitionKey: &pk,
		FeedRange:    feedRange,
		Continuation: &continuation,
	}

	partitionKeyRanges := []partitionKeyRange{
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
	if h[cosmosHeaderMaxItemCount] != "25" {
		t.Errorf("Expected MaxItemCount to be 25, got %v", h[cosmosHeaderMaxItemCount])
	}

	expectedIfModifiedSince := now.Format(time.RFC1123)
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

func TestChangeFeedOptionsCompositeContinuationToken(t *testing.T) {
	etag := azcore.ETag("test-etag")
	cfRange := newChangeFeedRange("00", "FF", &ChangeFeedRangeOptions{
		ContinuationToken: &etag,
	})
	compositeToken := newCompositeContinuationToken("test-resource-id", []changeFeedRange{cfRange})

	tokenBytes, err := json.Marshal(compositeToken)
	if err != nil {
		t.Fatalf("Failed to marshal composite token: %v", err)
	}
	tokenString := string(tokenBytes)

	options := &ChangeFeedOptions{
		Continuation: &tokenString,
	}

	headers := options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	h := *headers

	if h[headerIfNoneMatch] != string(etag) {
		t.Errorf("Expected IfNoneMatch to be %v, got %v", string(etag), h[headerIfNoneMatch])
	}

	if options.FeedRange == nil {
		t.Fatal("Expected FeedRange to be set from composite token")
	}
	if options.FeedRange.MinInclusive != "00" {
		t.Errorf("Expected FeedRange.MinInclusive to be 00, got %v", options.FeedRange.MinInclusive)
	}
	if options.FeedRange.MaxExclusive != "FF" {
		t.Errorf("Expected FeedRange.MaxExclusive to be FF, got %v", options.FeedRange.MaxExclusive)
	}
}

func TestChangeFeedOptionsCompositeContinuationTokenWithExistingFeedRange(t *testing.T) {
	etag := azcore.ETag("test-etag")
	cfRange := newChangeFeedRange("00", "FF", &ChangeFeedRangeOptions{
		ContinuationToken: &etag,
	})
	compositeToken := newCompositeContinuationToken("test-resource-id", []changeFeedRange{cfRange})

	tokenBytes, err := json.Marshal(compositeToken)
	if err != nil {
		t.Fatalf("Failed to marshal composite token: %v", err)
	}
	tokenString := string(tokenBytes)

	explicitFeedRange := &FeedRange{
		MinInclusive: "AA",
		MaxExclusive: "BB",
	}

	options := &ChangeFeedOptions{
		Continuation: &tokenString,
		FeedRange:    explicitFeedRange,
	}

	headers := options.toHeaders(nil)
	if headers == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	h := *headers

	if h[headerIfNoneMatch] != string(etag) {
		t.Errorf("Expected IfNoneMatch to be %v, got %v", string(etag), h[headerIfNoneMatch])
	}

	if options.FeedRange.MinInclusive != "AA" {
		t.Errorf("Expected FeedRange.MinInclusive to remain AA, got %v", options.FeedRange.MinInclusive)
	}
	if options.FeedRange.MaxExclusive != "BB" {
		t.Errorf("Expected FeedRange.MaxExclusive to remain BB, got %v", options.FeedRange.MaxExclusive)
	}
}
