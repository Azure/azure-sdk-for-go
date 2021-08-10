// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// CosmomAccountCredential is a wrapper interface for SharedKeyCredential
type CosmosAccountCredential interface {
	computeHMACSHA256(message string) (base64String string)
}
