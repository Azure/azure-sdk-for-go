// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

type EncryptionAlgorithmType = generated.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = generated.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = generated.EncryptionAlgorithmTypeAES256
)

type ImmutabilityPolicyMode = blob.ImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeLocked
)

// CopyStatusType defines values for CopyStatusType
type CopyStatusType = blob.CopyStatusType

const (
	CopyStatusTypePending CopyStatusType = blob.CopyStatusTypePending
	CopyStatusTypeSuccess CopyStatusType = blob.CopyStatusTypeSuccess
	CopyStatusTypeAborted CopyStatusType = blob.CopyStatusTypeAborted
	CopyStatusTypeFailed  CopyStatusType = blob.CopyStatusTypeFailed
)

// StatusType defines values for StatusType
type StatusType = azdatalake.StatusType

const (
	StatusTypeLocked   StatusType = azdatalake.StatusTypeLocked
	StatusTypeUnlocked StatusType = azdatalake.StatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return azdatalake.PossibleStatusTypeValues()
}

// DurationType defines values for DurationType
type DurationType = azdatalake.DurationType

const (
	DurationTypeInfinite DurationType = azdatalake.DurationTypeInfinite
	DurationTypeFixed    DurationType = azdatalake.DurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return azdatalake.PossibleDurationTypeValues()
}

// StateType defines values for StateType
type StateType = azdatalake.StateType

const (
	StateTypeAvailable StateType = azdatalake.StateTypeAvailable
	StateTypeLeased    StateType = azdatalake.StateTypeLeased
	StateTypeExpired   StateType = azdatalake.StateTypeExpired
	StateTypeBreaking  StateType = azdatalake.StateTypeBreaking
	StateTypeBroken    StateType = azdatalake.StateTypeBroken
)

type LeaseAction = generated.LeaseAction

const (
	LeaseActionAcquire        = generated.LeaseActionAcquire
	LeaseActionRelease        = generated.LeaseActionRelease
	LeaseActionAcquireRelease = generated.LeaseActionAcquireRelease
	LeaseActionRenew          = generated.LeaseActionAutoRenew
)

// LayoutAwareRouting defines whether downloads should attempt to be routed to the ideal
// endpoint for each chunk, based on the file's layout.
type LayoutAwareRouting = blob.LayoutAwareRouting

const (
	// LayoutAwareRoutingAuto is the zero value, and therefore the default when no value is
	// specified. Currently, the SDK resolves Auto to enabled, so it behaves identically to
	// LayoutAwareRoutingEnabled. Specify LayoutAwareRoutingEnabled or
	// LayoutAwareRoutingDisabled to pin the behavior.
	LayoutAwareRoutingAuto LayoutAwareRouting = blob.LayoutAwareRoutingAuto

	// LayoutAwareRoutingEnabled always attempts to route requests to the ideal endpoint for each chunk.
	LayoutAwareRoutingEnabled LayoutAwareRouting = blob.LayoutAwareRoutingEnabled

	// LayoutAwareRoutingDisabled never uses layout aware routing; requests are sent to the client's configured endpoint.
	LayoutAwareRoutingDisabled LayoutAwareRouting = blob.LayoutAwareRoutingDisabled
)

// PossibleLayoutAwareRoutingValues returns the possible values for the LayoutAwareRouting const type.
func PossibleLayoutAwareRoutingValues() []LayoutAwareRouting {
	return blob.PossibleLayoutAwareRoutingValues()
}
