//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
)

type EncryptionAlgorithmType = blob.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeAES256
)

// responses models:

type ImmutabilityPolicyMode = blob.ImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeLocked
)

type CopyStatusType = blob.CopyStatusType

const (
	CopyStatusTypePending CopyStatusType = blob.CopyStatusTypePending
	CopyStatusTypeSuccess CopyStatusType = blob.CopyStatusTypeSuccess
	CopyStatusTypeAborted CopyStatusType = blob.CopyStatusTypeAborted
	CopyStatusTypeFailed  CopyStatusType = blob.CopyStatusTypeFailed
)

type LeaseDurationType = lease.DurationType

const (
	LeaseDurationTypeInfinite LeaseDurationType = lease.DurationTypeInfinite
	LeaseDurationTypeFixed    LeaseDurationType = lease.DurationTypeFixed
)

type LeaseStateType = lease.StateType

const (
	LeaseStateTypeAvailable LeaseStateType = lease.StateTypeAvailable
	LeaseStateTypeLeased    LeaseStateType = lease.StateTypeLeased
	LeaseStateTypeExpired   LeaseStateType = lease.StateTypeExpired
	LeaseStateTypeBreaking  LeaseStateType = lease.StateTypeBreaking
	LeaseStateTypeBroken    LeaseStateType = lease.StateTypeBroken
)

type LeaseStatusType = lease.StatusType

const (
	LeaseStatusTypeLocked   LeaseStatusType = lease.StatusTypeLocked
	LeaseStatusTypeUnlocked LeaseStatusType = lease.StatusTypeUnlocked
)
