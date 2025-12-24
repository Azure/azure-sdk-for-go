// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// AccessTier defines values for the access tier of the share.
type AccessTier = generated.ShareAccessTier

const (
	AccessTierCool                 AccessTier = generated.ShareAccessTierCool
	AccessTierHot                  AccessTier = generated.ShareAccessTierHot
	AccessTierPremium              AccessTier = generated.ShareAccessTierPremium
	AccessTierTransactionOptimized AccessTier = generated.ShareAccessTierTransactionOptimized
)

// PermissionFormat contains the format of the file permissions, Can be sddl (Default) or Binary.
type PermissionFormat = generated.FilePermissionFormat

const (
	FilePermissionFormatBinary PermissionFormat = generated.FilePermissionFormatBinary
	FilePermissionFormatSddl   PermissionFormat = generated.FilePermissionFormatSddl
)

// PossibleFilePermissionFormatValues returns the possible values for the FilePermissionFormat const type.
func PossibleFilePermissionFormatValues() []PermissionFormat {
	return generated.PossibleFilePermissionFormatValues()
}

// PossibleAccessTierValues returns the possible values for the AccessTier const type.
func PossibleAccessTierValues() []AccessTier {
	return generated.PossibleShareAccessTierValues()
}

// RootSquash defines values for the root squashing behavior on the share when NFS is enabled. If it's not specified, the default is NoRootSquash.
type RootSquash = generated.ShareRootSquash

const (
	RootSquashNoRootSquash RootSquash = generated.ShareRootSquashNoRootSquash
	RootSquashRootSquash   RootSquash = generated.ShareRootSquashRootSquash
	RootSquashAllSquash    RootSquash = generated.ShareRootSquashAllSquash
)

// PossibleRootSquashValues returns the possible values for the RootSquash const type.
func PossibleRootSquashValues() []RootSquash {
	return generated.PossibleShareRootSquashValues()
}

// DeleteSnapshotsOptionType defines values for DeleteSnapshotsOptionType
type DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionType

const (
	DeleteSnapshotsOptionTypeInclude       DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeInclude
	DeleteSnapshotsOptionTypeIncludeLeased DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeIncludeLeased
)

// PossibleDeleteSnapshotsOptionTypeValues returns the possible values for the DeleteSnapshotsOptionType const type.
func PossibleDeleteSnapshotsOptionTypeValues() []DeleteSnapshotsOptionType {
	return generated.PossibleDeleteSnapshotsOptionTypeValues()
}

// TokenIntent is required if authorization header specifies an OAuth token.
type TokenIntent = generated.ShareTokenIntent

const (
	TokenIntentBackup TokenIntent = generated.ShareTokenIntentBackup
)

// PossibleTokenIntentValues returns the possible values for the TokenIntent const type.
func PossibleTokenIntentValues() []TokenIntent {
	return generated.PossibleShareTokenIntentValues()
}
