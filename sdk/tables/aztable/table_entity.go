// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "time"

type Entity struct {
	PartitionKey string
	RowKey       string
	TimeStampt   time.Time
}

type EdmType string

const (
	BINARY   EdmType = "Edm.Binary"
	BOOLEAN  EdmType = "Edm.Boolean"
	DATETIME EdmType = "Edm.DateTime"
	DOUBLE   EdmType = "Edm.Double"
	GUID     EdmType = "Edm.Guid"
	INT32    EdmType = "Edm.Int32"
	INT64    EdmType = "Edm.Int64"
	STRING   EdmType = "Edm.String"
)
