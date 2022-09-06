// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

const (
	headerXmsDate                 = "x-ms-date"
	headerAuthorization           = "Authorization"
	etagOData                     = "odata.etag"
	rfc3339                       = "2006-01-02T15:04:05.9999999Z"
	legacyCosmosTableDomain       = ".table.cosmosdb."
	cosmosTableDomain             = ".table.cosmos."
	headerContentType             = "Content-Type"
	headerContentTransferEncoding = "Content-Transfer-Encoding"
	timestamp                     = "Timestamp"
	partitionKey                  = "PartitionKey"
	rowKey                        = "RowKey"
)

// UpdateMode specifies what type of update to do on UpsertEntity or UpdateEntity. UpdateModeReplace
// will replace an existing entity, UpdateModeMerge will merge properties of the entities.
type UpdateMode string

const (
	UpdateModeReplace UpdateMode = "replace"
	UpdateModeMerge   UpdateMode = "merge"
)

// PossibleUpdateModeValues returns the possible values for the EntityUpdateMode const type.
func PossibleUpdateModeValues() []UpdateMode {
	return []UpdateMode{
		UpdateModeMerge,
		UpdateModeReplace,
	}
}
