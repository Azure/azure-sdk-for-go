// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

// StatusType defines values for StatusType
type StatusType = generated.LeaseStatus

const (
	StatusTypeLocked   StatusType = generated.LeaseStatusLocked
	StatusTypeUnlocked StatusType = generated.LeaseStatusUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return generated.PossibleLeaseStatusValues()
}

// DurationType defines values for DurationType
type DurationType = generated.LeaseDuration

const (
	DurationTypeInfinite DurationType = generated.LeaseDurationInfinite
	DurationTypeFixed    DurationType = generated.LeaseDurationFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return generated.PossibleLeaseDurationValues()
}

// StateType defines values for StateType
type StateType = generated.LeaseState

const (
	StateTypeAvailable StateType = generated.LeaseStateAvailable
	StateTypeLeased    StateType = generated.LeaseStateLeased
	StateTypeExpired   StateType = generated.LeaseStateExpired
	StateTypeBreaking  StateType = generated.LeaseStateBreaking
	StateTypeBroken    StateType = generated.LeaseStateBroken
)

// PossibleStateTypeValues returns the possible values for the StateType const type.
func PossibleStateTypeValues() []StateType {
	return generated.PossibleLeaseStateValues()
}
