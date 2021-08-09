// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// CosmomAccountCredential is a wrapper interface for SharedKeyCredential and UserDelegationCredential
type CosmosAccountCredential interface {
	AccountName() string
	ComputeHMACSHA256(message string) (base64String string)
}
