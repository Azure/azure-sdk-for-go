// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// FilePermissionFormat contains the format of the file permissions, Can be sddl (Default) or Binary.
type FilePermissionFormat = generated.FilePermissionFormat

const (
	FilePermissionFormatBinary FilePermissionFormat = generated.FilePermissionFormatBinary
	FilePermissionFormatSddl   FilePermissionFormat = generated.FilePermissionFormatSddl
)

// PossibleFilePermissionFormatValues returns the possible values for the FilePermissionFormat const type.
func PossibleFilePermissionFormatValues() []FilePermissionFormat {
	return generated.PossibleFilePermissionFormatValues()
}

// ListFilesIncludeType defines values for ListFilesIncludeType
type ListFilesIncludeType = generated.ListFilesIncludeType

const (
	ListFilesIncludeTypeTimestamps    ListFilesIncludeType = generated.ListFilesIncludeTypeTimestamps
	ListFilesIncludeTypeETag          ListFilesIncludeType = generated.ListFilesIncludeTypeEtag
	ListFilesIncludeTypeAttributes    ListFilesIncludeType = generated.ListFilesIncludeTypeAttributes
	ListFilesIncludeTypePermissionKey ListFilesIncludeType = generated.ListFilesIncludeTypePermissionKey
)

// PossibleListFilesIncludeTypeValues returns the possible values for the ListFilesIncludeType const type.
func PossibleListFilesIncludeTypeValues() []ListFilesIncludeType {
	return generated.PossibleListFilesIncludeTypeValues()
}

// ShareTokenIntent is required if authorization header specifies an OAuth token.
type ShareTokenIntent = generated.ShareTokenIntent

const (
	ShareTokenIntentBackup ShareTokenIntent = generated.ShareTokenIntentBackup
)

// PossibleShareTokenIntentValues returns the possible values for the ShareTokenIntent const type.
func PossibleShareTokenIntentValues() []ShareTokenIntent {
	return generated.PossibleShareTokenIntentValues()
}
