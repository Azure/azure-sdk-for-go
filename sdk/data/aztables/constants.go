// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

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

// GeoReplicationStatus - The status of the secondary location.
type GeoReplicationStatus string

const (
	GeoReplicationStatusBootstrap   GeoReplicationStatus = "bootstrap"
	GeoReplicationStatusLive        GeoReplicationStatus = "live"
	GeoReplicationStatusUnavailable GeoReplicationStatus = "unavailable"
)

// PossibleGeoReplicationStatusValues returns the possible values for the GeoReplicationStatus const type.
func PossibleGeoReplicationStatusValues() []GeoReplicationStatus {
	return []GeoReplicationStatus{
		GeoReplicationStatusBootstrap,
		GeoReplicationStatusLive,
		GeoReplicationStatusUnavailable,
	}
}

func toGeneratedStatusType(g *generated.GeoReplicationStatusType) *GeoReplicationStatus {
	if g == nil {
		return nil
	}
	if *g == generated.GeoReplicationStatusTypeBootstrap {
		return to.Ptr(GeoReplicationStatusBootstrap)
	}
	if *g == generated.GeoReplicationStatusTypeLive {
		return to.Ptr(GeoReplicationStatusLive)
	}
	if *g == generated.GeoReplicationStatusTypeUnavailable {
		return to.Ptr(GeoReplicationStatusUnavailable)
	}
	return nil
}

// MetadataFormat specifies the level of OData metadata returned with an entity.
// https://learn.microsoft.com/rest/api/storageservices/payload-format-for-table-service-operations#json-format-applicationjson-versions-2013-08-15-and-later
type MetadataFormat = generated.ODataMetadataFormat

const (
	MetadataFormatFull    MetadataFormat = generated.ODataMetadataFormatApplicationJSONODataFullmetadata
	MetadataFormatMinimal MetadataFormat = generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata
	MetadataFormatNone    MetadataFormat = generated.ODataMetadataFormatApplicationJSONODataNometadata
)

// SASProtocol indicates the SAS protocol
type SASProtocol string

const (
	// SASProtocolHTTPS can be specified for a SAS protocol
	SASProtocolHTTPS SASProtocol = "https"

	// SASProtocolHTTPSandHTTP can be specified for a SAS protocol
	SASProtocolHTTPSandHTTP SASProtocol = "https,http"
)

// PossibleSASProtocolValues returns the possible values for the SASProtocol const type.
func PossibleSASProtocolValues() []SASProtocol {
	return []SASProtocol{
		SASProtocolHTTPS,
		SASProtocolHTTPSandHTTP,
	}
}

// TransactionType is the type for a specific transaction operation.
type TransactionType string

const (
	TransactionTypeAdd           TransactionType = "add"
	TransactionTypeUpdateMerge   TransactionType = "updatemerge"
	TransactionTypeUpdateReplace TransactionType = "updatereplace"
	TransactionTypeDelete        TransactionType = "delete"
	TransactionTypeInsertMerge   TransactionType = "insertmerge"
	TransactionTypeInsertReplace TransactionType = "insertreplace"
)

// PossibleTransactionTypeValues returns the possible values for the TransactionType const type.
func PossibleTransactionTypeValues() []TransactionType {
	return []TransactionType{
		TransactionTypeAdd,
		TransactionTypeUpdateMerge,
		TransactionTypeUpdateReplace,
		TransactionTypeDelete,
		TransactionTypeInsertMerge,
		TransactionTypeInsertReplace,
	}
}

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
