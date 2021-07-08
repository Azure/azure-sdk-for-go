// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "errors"

const (
	timestamp     = "Timestamp"
	partitionKey  = "PartitionKey"
	rowKey        = "RowKey"
	etagOdata     = "odata.etag"
	etag          = "ETag"
	odataMetadata = "odata.metadata"
	OdataType     = "@odata.type"
	edmBinary     = "Edm.Binary"
	edmBoolean    = "Emd.Boolean"
	edmDateTime   = "Edm.DateTime"
	edmDouble     = "Edm.Double"
	edmGuid       = "Edm.Guid"
	edmInt32      = "Edm.Int32"
	edmInt64      = "Edm.Int64"
	edmString     = "Edm.String"
	ISO8601       = "2006-01-02T15:04:05.9999999Z"
)

var partitionKeyRowKeyError = errors.New("Entity must have a PartitionKey and RowKey")
