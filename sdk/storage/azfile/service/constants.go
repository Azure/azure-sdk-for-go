//go:build go1.18
// +build go1.18

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
