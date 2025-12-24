// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
)

// PublicAccessType defines values for AccessType - private (default) or file or filesystem.
type PublicAccessType = filesystem.PublicAccessType

// Not to be used anymore as public access is disabled.
const (
	File       PublicAccessType = filesystem.File
	FileSystem PublicAccessType = filesystem.FileSystem
)

// StatusType defines values for StatusType
type StatusType = generated_blob.LeaseStatusType

const (
	StatusTypeLocked   StatusType = generated_blob.LeaseStatusTypeLocked
	StatusTypeUnlocked StatusType = generated_blob.LeaseStatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return generated_blob.PossibleLeaseStatusTypeValues()
}

// DurationType defines values for DurationType
type DurationType = generated_blob.LeaseDurationType

const (
	DurationTypeInfinite DurationType = generated_blob.LeaseDurationTypeInfinite
	DurationTypeFixed    DurationType = generated_blob.LeaseDurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return generated_blob.PossibleLeaseDurationTypeValues()
}

// StateType defines values for StateType
type StateType = generated_blob.LeaseStateType

const (
	StateTypeAvailable StateType = generated_blob.LeaseStateTypeAvailable
	StateTypeLeased    StateType = generated_blob.LeaseStateTypeLeased
	StateTypeExpired   StateType = generated_blob.LeaseStateTypeExpired
	StateTypeBreaking  StateType = generated_blob.LeaseStateTypeBreaking
	StateTypeBroken    StateType = generated_blob.LeaseStateTypeBroken
)
