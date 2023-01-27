//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// NtfsFileAttributes for Files and Directories.
// The subset of attributes is listed at: https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties#file-system-attributes.
// Their respective values are listed at: https://learn.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants.
type NtfsFileAttributes uint32

const (
	Readonly          NtfsFileAttributes = 1
	Hidden            NtfsFileAttributes = 2
	System            NtfsFileAttributes = 4
	Directory         NtfsFileAttributes = 16
	Archive           NtfsFileAttributes = 32
	None              NtfsFileAttributes = 128
	Temporary         NtfsFileAttributes = 256
	Offline           NtfsFileAttributes = 4096
	NotContentIndexed NtfsFileAttributes = 8192
	NoScrubData       NtfsFileAttributes = 131072
)

// PossibleNtfsFileAttributesValues returns the possible values for the NtfsFileAttributes const type.
func PossibleNtfsFileAttributesValues() []NtfsFileAttributes {
	return []NtfsFileAttributes{
		Readonly,
		Hidden,
		System,
		Directory,
		Archive,
		None,
		Temporary,
		Offline,
		NotContentIndexed,
		NoScrubData,
	}
}

// CopyStatusType defines the states of the copy operation.
type CopyStatusType = generated.CopyStatusType

const (
	CopyStatusTypePending CopyStatusType = generated.CopyStatusTypePending
	CopyStatusTypeSuccess CopyStatusType = generated.CopyStatusTypeSuccess
	CopyStatusTypeAborted CopyStatusType = generated.CopyStatusTypeAborted
	CopyStatusTypeFailed  CopyStatusType = generated.CopyStatusTypeFailed
)

// PossibleCopyStatusTypeValues returns the possible values for the CopyStatusType const type.
func PossibleCopyStatusTypeValues() []CopyStatusType {
	return generated.PossibleCopyStatusTypeValues()
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

// PermissionCopyModeType determines the copy behavior of the security descriptor of the file.
//   - source: The security descriptor on the destination file is copied from the source file.
//   - override: The security descriptor on the destination file is determined via the x-ms-file-permission or x-ms-file-permission-key header.
type PermissionCopyModeType = generated.PermissionCopyModeType

const (
	PermissionCopyModeTypeSource   PermissionCopyModeType = generated.PermissionCopyModeTypeSource
	PermissionCopyModeTypeOverride PermissionCopyModeType = generated.PermissionCopyModeTypeOverride
)

// PossiblePermissionCopyModeTypeValues returns the possible values for the PermissionCopyModeType const type.
func PossiblePermissionCopyModeTypeValues() []PermissionCopyModeType {
	return generated.PossiblePermissionCopyModeTypeValues()
}

// RangeWriteType represents one of the following options.
//   - update: Writes the bytes specified by the request body into the specified range. The Range and Content-Length headers must match to perform the update.
//   - clear: Clears the specified range and releases the space used in storage for that range. To clear a range, set the Content-Length header to zero,
//     and set the Range header to a value that indicates the range to clear, up to maximum file size.
type RangeWriteType = generated.FileRangeWriteType

const (
	RangeWriteTypeUpdate RangeWriteType = generated.FileRangeWriteTypeUpdate
	RangeWriteTypeClear  RangeWriteType = generated.FileRangeWriteTypeClear
)

// PossibleRangeWriteTypeValues returns the possible values for the RangeWriteType const type.
func PossibleRangeWriteTypeValues() []RangeWriteType {
	return generated.PossibleFileRangeWriteTypeValues()
}
