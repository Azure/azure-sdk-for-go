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
	etag                          = "ETag"
)

// ResponseFormat determines what is returned from a service request
type ResponseFormat string

const (
	ResponseFormatReturnContent   ResponseFormat = "return-content"
	ResponseFormatReturnNoContent ResponseFormat = "return-no-content"
)

// PossibleResponseFormatValues returns the possible values for the ResponseFormat const type.
func PossibleResponseFormatValues() []ResponseFormat {
	return []ResponseFormat{
		ResponseFormatReturnContent,
		ResponseFormatReturnNoContent,
	}
}

// ToPtr returns a *ResponseFormat pointing to the current value.
func (c ResponseFormat) ToPtr() *ResponseFormat {
	return &c
}

// EntityUpdateMode specifies what type of update to do on InsertEntity or UpdateEntity. ReplaceEntity
// will replace an existing entity, MergeEntity will merge properties of the entities.
type EntityUpdateMode string

const (
	ReplaceEntity EntityUpdateMode = "replace"
	MergeEntity   EntityUpdateMode = "merge"
)
