//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

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

// Deprecated: Internal implementation; use FormatPermissions instead.
// Format returns file permission string and permission key.
func (p *Permissions) Format(defaultFilePermissionStr string) (*string, *string) {
	return nil, nil
}

// FormatPermissions returns file permission string and permission key.
func FormatPermissions(p *Permissions) (*string, *string) {
	if p == nil {
		return nil, nil
	}

	if p.Permission == nil && p.PermissionKey == nil {
		return nil, nil
	} else {
		return p.Permission, p.PermissionKey
	}
}
