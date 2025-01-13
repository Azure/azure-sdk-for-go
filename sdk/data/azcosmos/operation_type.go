// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// operationType defines supported values for operations.
type operationType int

const (
	operationTypeCreate  operationType = 0
	operationTypePatch   operationType = 1
	operationTypeRead    operationType = 2
	operationTypeReplace operationType = 5
	operationTypeDelete  operationType = 4
	operationTypeUpsert  operationType = 20
	operationTypeQuery   operationType = 15
	operationTypeBatch   operationType = 40
)
