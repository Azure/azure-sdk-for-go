// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

// TestChangeFeedOptions_BuildRequestHeaders_Defaults verifies the minimum
// header set produced for an empty options value: only the change-feed AIM
// header is set; no MaxItemCount, no IfModifiedSince, no PartitionKey.
func TestChangeFeedOptions_BuildRequestHeaders_Defaults(t *testing.T) {
	options := &ChangeFeedOptions{}
	headers, err := options.buildRequestHeaders(changeFeedRange{}, "")
	require.NoError(t, err)
	require.Equal(t, cosmosHeaderValuesChangeFeed, headers[cosmosHeaderChangeFeed])
	_, hasMaxItem := headers[cosmosHeaderMaxItemCount]
	require.False(t, hasMaxItem, "MaxItemCount header must be omitted when 0")
	_, hasIfModified := headers[cosmosHeaderIfModifiedSince]
	require.False(t, hasIfModified, "IfModifiedSince must be omitted when StartFrom is nil")
	_, hasPK := headers[cosmosHeaderPartitionKey]
	require.False(t, hasPK, "PartitionKey header must be omitted when PartitionKey is nil")
	_, hasIfNoneMatch := headers[headerIfNoneMatch]
	require.False(t, hasIfNoneMatch, "IfNoneMatch must be omitted when head ContinuationToken is nil")
	_, hasPKRangeID := headers[headerXmsDocumentDbPartitionKeyRangeId]
	require.False(t, hasPKRangeID, "PK-range-id header must be omitted when resolvedPKRangeID is empty")
}

// TestChangeFeedOptions_BuildRequestHeaders_AllFields verifies that every
// supported field of ChangeFeedOptions plus the resolved head/PK-range ID
// is materialized into the expected header.
func TestChangeFeedOptions_BuildRequestHeaders_AllFields(t *testing.T) {
	now := time.Now().UTC()
	pk := NewPartitionKeyString("testPK")
	etag := azcore.ETag("\"etag-12345\"")

	options := &ChangeFeedOptions{
		MaxItemCount: 25,
		StartFrom:    &now,
		PartitionKey: &pk,
	}
	head := changeFeedRange{
		MinInclusive:      "10",
		MaxExclusive:      "20",
		ContinuationToken: &etag,
	}

	headers, err := options.buildRequestHeaders(head, "range1")
	require.NoError(t, err)

	require.Equal(t, cosmosHeaderValuesChangeFeed, headers[cosmosHeaderChangeFeed])
	require.Equal(t, "25", headers[cosmosHeaderMaxItemCount])
	require.Equal(t, now.Format(time.RFC1123), headers[cosmosHeaderIfModifiedSince])

	pkJSON, _ := pk.toJsonString()
	require.Equal(t, string(pkJSON), headers[cosmosHeaderPartitionKey])

	require.Equal(t, "range1", headers[headerXmsDocumentDbPartitionKeyRangeId])
	require.Equal(t, string(etag), headers[headerIfNoneMatch])
}

// TestChangeFeedOptions_BuildRequestHeaders_EmptyContinuationOmitsIfNoneMatch
// verifies that an explicitly-set-but-empty continuation token does not
// produce an IfNoneMatch header (would otherwise cause the server to treat
// every read as conditional against the empty string).
func TestChangeFeedOptions_BuildRequestHeaders_EmptyContinuationOmitsIfNoneMatch(t *testing.T) {
	emptyETag := azcore.ETag("")
	head := changeFeedRange{
		MinInclusive:      "00",
		MaxExclusive:      "FF",
		ContinuationToken: &emptyETag,
	}

	headers, err := (&ChangeFeedOptions{}).buildRequestHeaders(head, "0")
	require.NoError(t, err)
	_, exists := headers[headerIfNoneMatch]
	require.False(t, exists, "empty ContinuationToken must NOT produce an IfNoneMatch header")
}

// TestChangeFeedOptions_BuildRequestHeaders_PartitionKeySerializationError
// exercises the new error-returning contract: a PartitionKey whose
// serialization fails is surfaced to the caller rather than silently
// dropped (the original-bug pattern from toHeaders).
func TestChangeFeedOptions_BuildRequestHeaders_PartitionKeySerializationError(t *testing.T) {
	// json.Marshal cannot serialize a channel; toJsonString surfaces that error.
	pk := NewPartitionKey()
	pk.values = []interface{}{make(chan int)}
	options := &ChangeFeedOptions{PartitionKey: &pk}

	_, err := options.buildRequestHeaders(changeFeedRange{}, "")
	require.Error(t, err, "buildRequestHeaders must surface PartitionKey serialization errors")
}
