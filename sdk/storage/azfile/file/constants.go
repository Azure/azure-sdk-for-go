//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

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
