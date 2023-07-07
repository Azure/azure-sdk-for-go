//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"

// PublicAccessType defines values for AccessType - private (default) or file or filesystem.
type PublicAccessType = azblob.PublicAccessType

const (
	File       PublicAccessType = azblob.PublicAccessTypeBlob
	Filesystem PublicAccessType = azblob.PublicAccessTypeContainer
)

// TODO: figure out a way to import this from datalake rather than blob again

// StatusType defines values for StatusType
type StatusType = lease.StatusType

const (
	StatusTypeLocked   StatusType = lease.StatusTypeLocked
	StatusTypeUnlocked StatusType = lease.StatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return lease.PossibleStatusTypeValues()
}

// DurationType defines values for DurationType
type DurationType = lease.DurationType

const (
	DurationTypeInfinite DurationType = lease.DurationTypeInfinite
	DurationTypeFixed    DurationType = lease.DurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return lease.PossibleDurationTypeValues()
}

// StateType defines values for StateType
type StateType = lease.StateType

const (
	StateTypeAvailable StateType = lease.StateTypeAvailable
	StateTypeLeased    StateType = lease.StateTypeLeased
	StateTypeExpired   StateType = lease.StateTypeExpired
	StateTypeBreaking  StateType = lease.StateTypeBreaking
	StateTypeBroken    StateType = lease.StateTypeBroken
)

// PossibleStateTypeValues returns the possible values for the StateType const type.
func PossibleStateTypeValues() []StateType {
	return lease.PossibleStateTypeValues()
}
