// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type cosmosOperationContext struct {
	resourceType    resourceType
	resourceAddress string
	isRidBased      bool
}
