//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// ListSharesIncludeType defines values for ListSharesIncludeType
type ListSharesIncludeType = generated.ListSharesIncludeType

const (
	ListSharesIncludeTypeSnapshots ListSharesIncludeType = generated.ListSharesIncludeTypeSnapshots
	ListSharesIncludeTypeMetadata  ListSharesIncludeType = generated.ListSharesIncludeTypeMetadata
	ListSharesIncludeTypeDeleted   ListSharesIncludeType = generated.ListSharesIncludeTypeDeleted
)

// PossibleListSharesIncludeTypeValues returns the possible values for the ListSharesIncludeType const type.
func PossibleListSharesIncludeTypeValues() []ListSharesIncludeType {
	return generated.PossibleListSharesIncludeTypeValues()
}

// ShareRootSquash defines values for the root squashing behavior on the share when NFS is enabled. If it's not specified, the default is NoRootSquash.
type ShareRootSquash = generated.ShareRootSquash

const (
	RootSquashNoRootSquash ShareRootSquash = generated.ShareRootSquashNoRootSquash
	RootSquashRootSquash   ShareRootSquash = generated.ShareRootSquashRootSquash
	RootSquashAllSquash    ShareRootSquash = generated.ShareRootSquashAllSquash
)

// PossibleShareRootSquashValues returns the possible values for the RootSquash const type.
func PossibleShareRootSquashValues() []ShareRootSquash {
	return generated.PossibleShareRootSquashValues()
}

// ShareTokenIntent is required if authorization header specifies an OAuth token.
type ShareTokenIntent = generated.ShareTokenIntent

const (
	ShareTokenIntentBackup ShareTokenIntent = generated.ShareTokenIntentBackup
)

// PossibleShareTokenIntentValues returns the possible values for the ShareTokenIntent const type.
func PossibleShareTokenIntentValues() []ShareTokenIntent {
	return generated.PossibleShareTokenIntentValues()
}
