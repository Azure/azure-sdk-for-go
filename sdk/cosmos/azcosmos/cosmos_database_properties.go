// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// CosmosDatabaseProperties represents the properties of a database.
type CosmosDatabaseProperties struct {
	// Id contains the unique id of the database.
	Id string `json:"id"`
	// ETag contains the entity etag of the database
	ETag string `json:"_etag,omitempty"`
	// SelfLink contains the self-link of the database
	SelfLink string `json:"_self,omitempty"`
	// ResourceId contains the resource id of the database
	ResourceId string `json:"_rid,omitempty"`
	// LastModified contains the last modified time of the database
	LastModified *UnixTime `json:"_ts,omitempty"`
	// Database represented by these properties
	Database *CosmosDatabase `json:"-"`
}
