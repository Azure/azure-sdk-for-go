//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"encoding/binary"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

const (
	_1MiB      = 1024 * 1024
	CountToEnd = 0

	// MaxUpdateRangeBytes indicates the maximum number of bytes that can be updated in a call to Client.UploadRange.
	MaxUpdateRangeBytes = 4 * 1024 * 1024 // 4MiB

	// MaxFileSize indicates the maximum size of the file allowed.
	MaxFileSize = 4 * 1024 * 1024 * 1024 * 1024 // 4 TiB

	// DefaultDownloadChunkSize is default chunk size
	DefaultDownloadChunkSize = int64(4 * 1024 * 1024) // 4MiB
)

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

// NfsFileType specifies the type of the file or directory.
type NfsFileType = generated.NfsFileType

// ModeCopyMode specifies the mode of the file or directory.
type ModeCopyMode = generated.ModeCopyMode

// OwnerCopyMode specifies the copy mode source or override.
type OwnerCopyMode = generated.OwnerCopyMode

const (
	// NFSFileTypeRegular Default and only value for the parameter NFS File Type.
	NFSFileTypeRegular   NfsFileType = generated.NfsFileTypeRegular
	NfsFileTypeDirectory NfsFileType = generated.NfsFileTypeDirectory
	NfsFileTypeSymlink   NfsFileType = generated.NfsFileTypeSymLink

	OwnerCopyModeOverride OwnerCopyMode = generated.OwnerCopyModeOverride
	OwnerCopyModeSource   OwnerCopyMode = generated.OwnerCopyModeSource

	ModeCopyModeOverride ModeCopyMode = generated.ModeCopyModeOverride
	ModeCopyModeSource   ModeCopyMode = generated.ModeCopyModeSource
)

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

// PermissionFormat contains the format of the file permissions, Can be sddl (Default) or Binary.
type PermissionFormat = generated.FilePermissionFormat

const (
	FilePermissionFormatBinary PermissionFormat = generated.FilePermissionFormatBinary
	FilePermissionFormatSddl   PermissionFormat = generated.FilePermissionFormatSddl
)

// PossibleFilePermissionFormatValues returns the possible values for the FilePermissionFormat const type.
func PossibleFilePermissionFormatValues() []PermissionFormat {
	return generated.PossibleFilePermissionFormatValues()
}

// PossibleRangeWriteTypeValues returns the possible values for the RangeWriteType const type.
func PossibleRangeWriteTypeValues() []RangeWriteType {
	return generated.PossibleFileRangeWriteTypeValues()
}

// SourceContentValidationType abstracts mechanisms used to validate source content
type SourceContentValidationType interface {
	apply(generated.SourceContentSetter) error
	notPubliclyImplementable()
}

// SourceContentValidationTypeCRC64 is a SourceContentValidationType used to provide a precomputed CRC64.
type SourceContentValidationTypeCRC64 uint64

// Apply implements the SourceContentValidationType interface for type SourceContentValidationTypeCRC64.
func (s SourceContentValidationTypeCRC64) apply(cfg generated.SourceContentSetter) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(s))
	cfg.SetSourceContentCRC64(buf)
	return nil
}

func (SourceContentValidationTypeCRC64) notPubliclyImplementable() {}

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType = exported.TransferValidationType

// TransferValidationTypeMD5 is a TransferValidationType used to provide a precomputed MD5.
type TransferValidationTypeMD5 = exported.TransferValidationTypeMD5

// ShareTokenIntent is required if authorization header specifies an OAuth token.
type ShareTokenIntent = generated.ShareTokenIntent

const (
	ShareTokenIntentBackup ShareTokenIntent = generated.ShareTokenIntentBackup
)

// PossibleShareTokenIntentValues returns the possible values for the ShareTokenIntent const type.
func PossibleShareTokenIntentValues() []ShareTokenIntent {
	return generated.PossibleShareTokenIntentValues()
}

// LastWrittenMode specifies if the file last write time should be preserved or overwritten
type LastWrittenMode = generated.FileLastWrittenMode

const (
	LastWrittenModeNow      LastWrittenMode = generated.FileLastWrittenModeNow
	LastWrittenModePreserve LastWrittenMode = generated.FileLastWrittenModePreserve
)

// PossibleLastWrittenModeValues returns the possible values for the LastWrittenMode const type.
func PossibleLastWrittenModeValues() []LastWrittenMode {
	return generated.PossibleFileLastWrittenModeValues()
}
