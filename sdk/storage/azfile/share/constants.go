//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// AccessTier defines values for the access tier of the share.
type AccessTier = generated.ShareAccessTier

const (
	AccessTierCool                 AccessTier = generated.ShareAccessTierCool
	AccessTierHot                  AccessTier = generated.ShareAccessTierHot
	AccessTierTransactionOptimized AccessTier = generated.ShareAccessTierTransactionOptimized
)

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

// LeaseDurationType - When a share is leased, specifies whether the lease is of infinite or fixed duration.
type LeaseDurationType = generated.LeaseDurationType

const (
	LeaseDurationTypeInfinite LeaseDurationType = generated.LeaseDurationTypeInfinite
	LeaseDurationTypeFixed    LeaseDurationType = generated.LeaseDurationTypeFixed
)

// PossibleLeaseDurationTypeValues returns the possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return generated.PossibleLeaseDurationTypeValues()
}

// LeaseStateType - Lease state of the share.
type LeaseStateType = generated.LeaseStateType

const (
	LeaseStateTypeAvailable LeaseStateType = generated.LeaseStateTypeAvailable
	LeaseStateTypeLeased    LeaseStateType = generated.LeaseStateTypeLeased
	LeaseStateTypeExpired   LeaseStateType = generated.LeaseStateTypeExpired
	LeaseStateTypeBreaking  LeaseStateType = generated.LeaseStateTypeBreaking
	LeaseStateTypeBroken    LeaseStateType = generated.LeaseStateTypeBroken
)

// PossibleLeaseStateTypeValues returns the possible values for the LeaseStateType const type.
func PossibleLeaseStateTypeValues() []LeaseStateType {
	return generated.PossibleLeaseStateTypeValues()
}

// LeaseStatusType - The current lease status of the share.
type LeaseStatusType = generated.LeaseStatusType

const (
	LeaseStatusTypeLocked   LeaseStatusType = generated.LeaseStatusTypeLocked
	LeaseStatusTypeUnlocked LeaseStatusType = generated.LeaseStatusTypeUnlocked
)

// PossibleLeaseStatusTypeValues returns the possible values for the LeaseStatusType const type.
func PossibleLeaseStatusTypeValues() []LeaseStatusType {
	return generated.PossibleLeaseStatusTypeValues()
}
