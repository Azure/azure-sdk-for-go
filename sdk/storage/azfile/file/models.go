//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import "time"

// SMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
type SMBProperties struct {
	// NtfsFileAttributes for Files and Directories. Default value is ‘None’ for file and ‘Directory’
	// for directory. ‘None’ can also be specified as default.
	Attributes *NtfsFileAttributes
	// The Coordinated Universal Time (UTC) creation time for the file/directory. Default value is 'now'.
	CreationTime *time.Time
	// The Coordinated Universal Time (UTC) last write time for the file/directory. Default value is 'now'.
	LastWriteTime *time.Time
}

// Permissions contains the optional parameters for the permissions on the file.
type Permissions struct {
	// If specified the permission (security descriptor) shall be set for the directory/file. This header can be used if Permission
	// size is <= 8KB, else x-ms-file-permission-key header shall be used. Default
	// value: Inherit. If SDDL is specified as input, it must have owner, group and dacl. Note: Only one of the x-ms-file-permission
	// or x-ms-file-permission-key should be specified.
	Permission *string
	// Key of the permission to be set for the directory/file.
	// Note: Only one of the x-ms-file-permission or x-ms-file-permission-key should be specified.
	PermissionKey *string
}
