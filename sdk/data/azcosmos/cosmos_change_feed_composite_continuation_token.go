// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type compositeContinuationToken struct {
	// ResourceID is the ID of the resource for which the continuation token is valid.
	ResourceID string `json:"resourceId"`
	// Continuation is the list of Epk Ranges part of the continuation token
	Continuation []ChangeFeedRange `json:"continuation"`
}

// newCompositeContinuationToken creates a new CompositeContinuationToken with the specified resource ID and continuation ranges.
// This function is used to create a continuation token for the Cosmos DB change feed.
// It is designed for internal use only and should not be used directly by clients.
func newCompositeContinuationToken(resourceID string, continuation []ChangeFeedRange) compositeContinuationToken {
	return compositeContinuationToken{
		ResourceID:   resourceID,
		Continuation: continuation,
	}
}
