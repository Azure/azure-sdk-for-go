// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

const ServiceVersion = "2026-06-06"

type ExpiryOptions string

const (
	ExpiryOptionsAbsolute           ExpiryOptions = "Absolute"
	ExpiryOptionsNeverExpire        ExpiryOptions = "NeverExpire"
	ExpiryOptionsRelativeToCreation ExpiryOptions = "RelativeToCreation"
	ExpiryOptionsRelativeToNow      ExpiryOptions = "RelativeToNow"
)

type LeaseDurationType string

const (
	LeaseDurationTypeFixed    LeaseDurationType = "fixed"
	LeaseDurationTypeInfinite LeaseDurationType = "infinite"
)

// PossibleLeaseDurationTypeValues returns the possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return []LeaseDurationType{
		LeaseDurationTypeFixed,
		LeaseDurationTypeInfinite,
	}
}

type LeaseStateType string

const (
	LeaseStateTypeAvailable LeaseStateType = "available"
	LeaseStateTypeBreaking  LeaseStateType = "breaking"
	LeaseStateTypeBroken    LeaseStateType = "broken"
	LeaseStateTypeExpired   LeaseStateType = "expired"
	LeaseStateTypeLeased    LeaseStateType = "leased"
)

type LeaseStatusType string

const (
	LeaseStatusTypeLocked   LeaseStatusType = "locked"
	LeaseStatusTypeUnlocked LeaseStatusType = "unlocked"
)

// PossibleLeaseStatusTypeValues returns the possible values for the LeaseStatusType const type.
func PossibleLeaseStatusTypeValues() []LeaseStatusType {
	return []LeaseStatusType{
		LeaseStatusTypeLocked,
		LeaseStatusTypeUnlocked,
	}
}

type PublicAccessType string

const (
	PublicAccessTypeFile       PublicAccessType = "blob"
	PublicAccessTypeFileSystem PublicAccessType = "container"
)

type ListContainersIncludeType string

const (
	ListContainersIncludeTypeDeleted  ListContainersIncludeType = "deleted"
	ListContainersIncludeTypeMetadata ListContainersIncludeType = "metadata"
	ListContainersIncludeTypeSystem   ListContainersIncludeType = "system"
)
