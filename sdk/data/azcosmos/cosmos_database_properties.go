// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// DatabaseProperties represents the properties of a database.
type DatabaseProperties struct {
	// ID contains the unique id of the database.
	ID string `json:"id"`
	// ETag contains the entity etag of the database
	ETag azcore.ETag `json:"_etag,omitempty"`
	// SelfLink contains the self-link of the database
	SelfLink string `json:"_self,omitempty"`
	// ResourceID contains the resource id of the database
	ResourceID string `json:"_rid,omitempty"`
	// LastModified contains the last modified time of the database
	LastModified int64 `json:"_ts,omitempty"`
}
