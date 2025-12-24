// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// DurationType - When a share is leased, specifies whether the lease is of infinite or fixed duration.
type DurationType = generated.LeaseDurationType

const (
	DurationTypeInfinite DurationType = generated.LeaseDurationTypeInfinite
	DurationTypeFixed    DurationType = generated.LeaseDurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return generated.PossibleLeaseDurationTypeValues()
}

// StateType - Lease state of the share.
type StateType = generated.LeaseStateType

const (
	StateTypeAvailable StateType = generated.LeaseStateTypeAvailable
	StateTypeLeased    StateType = generated.LeaseStateTypeLeased
	StateTypeExpired   StateType = generated.LeaseStateTypeExpired
	StateTypeBreaking  StateType = generated.LeaseStateTypeBreaking
	StateTypeBroken    StateType = generated.LeaseStateTypeBroken
)

// PossibleStateTypeValues returns the possible values for the StateType const type.
func PossibleStateTypeValues() []StateType {
	return generated.PossibleLeaseStateTypeValues()
}

// StatusType - The current lease status of the share.
type StatusType = generated.LeaseStatusType

const (
	StatusTypeLocked   StatusType = generated.LeaseStatusTypeLocked
	StatusTypeUnlocked StatusType = generated.LeaseStatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return generated.PossibleLeaseStatusTypeValues()
}
