// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Version 1 is the initial version of the composite continuation token.
const cosmosCompositeContinuationTokenVersion = 1

type compositeContinuationToken struct {
	// Version is the version of the continuation token format.
	Version int `json:"version,omitempty"`
	// ResourceID is the ID of the resource for which the continuation token is valid.
	ResourceID string `json:"resourceId"`
	// Continuation is the list of Epk Ranges part of the continuation token
	Continuation []changeFeedRange `json:"continuation"`
}

// newCompositeContinuationToken creates a new CompositeContinuationToken with the specified resource ID and continuation ranges.
// This function is used to create a continuation token for the Cosmos DB change feed.
// It is designed for internal use only and should not be used directly by clients.
func newCompositeContinuationToken(resourceID string, continuation []changeFeedRange) compositeContinuationToken {
	return compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		ResourceID:   resourceID,
		Continuation: continuation,
	}
}
