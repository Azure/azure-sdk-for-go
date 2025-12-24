// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// CopyFileNFSProperties contains the optional parameters regarding the NFS properties for a file.
type CopyFileNFSProperties struct {
	// Specifies either the option to copy file creation time from a source file(source) to a target file or a time value in ISO
	// 8601 format to set as creation time on a target file.
	// CopyFileCreationTime is an interface and its underlying implementation are:
	//   - SourceCopyFileCreationTime - specifies to copy file creation time from a source file to a target file.
	//   - DestinationCopyFileCreationTime - specifies a time value in ISO 8601 format to set as creation time on a target file.
	CreationTime CopyFileCreationTime
	// Specifies either the option to copy file last write time from a source file(source) to a target file or a time value in
	// ISO 8601 format to set as last write time on a target file.
	// CopyFileLastWriteTime is an interface and its underlying implementation are:
	//   - SourceCopyFileLastWriteTime - specifies to copy file last write time from a source file to a target file.
	//   - DestinationCopyFileLastWriteTime - specifies a time value in ISO 8601 format to set as last write time on a target file.
	LastWriteTime CopyFileLastWriteTime
	// The file mode of the file or directory
	FileMode *string
	// The owner of the file or directory.
	Owner *string
	// The owning group of the file or directory.
	Group *string
	// NFS only. Applicable only when the copy source is a File. Determines the copy behavior of the mode bits of the file.
	// source: The mode on the destination file is copied from the source file.
	// override: The mode on the destination file is determined via the x-ms-mode header.
	FileModeCopyMode *generated.ModeCopyMode
	// NFS only. Determines the copy behavior of the owner user identifier (UID) and group identifier (GID) of the file.
	// source: The owner user identifier (UID) and group identifier (GID) on the destination
	// file is copied from the source file. override: The owner user identifier (UID) and group identifier (GID) on the destination
	// file is determined via the x-ms-owner and x-ms-group headers.
	FileOwnerCopyMode *generated.OwnerCopyMode
}

// FormatCopyFileNFSProperties returns creation time, last write time.
func FormatCopyFileNFSProperties(np *CopyFileNFSProperties) (opts *generated.CopyFileSMBInfo) {
	opts = &generated.CopyFileSMBInfo{}

	if np == nil {
		return nil
	}

	if np.CreationTime != nil {
		opts.FileCreationTime = np.CreationTime.FormatCreationTime()
	}

	if np.LastWriteTime != nil {
		opts.FileLastWriteTime = np.LastWriteTime.FormatLastWriteTime()
	}

	return
}
