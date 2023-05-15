//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"strings"
	"time"
)

// SMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
type SMBProperties struct {
	// NTFSFileAttributes for Files and Directories. Default value is 'None' for file and
	// 'Directory' for directory. ‘None’ can also be specified as default.
	Attributes *NTFSFileAttributes
	// The Coordinated Universal Time (UTC) creation time for the file/directory. Default value is 'now'.
	CreationTime *time.Time
	// The Coordinated Universal Time (UTC) last write time for the file/directory. Default value is 'now'.
	LastWriteTime *time.Time
}

// Format returns file attributes, creation time and last write time.
func (sp *SMBProperties) Format(isDir bool, defaultFileAttributes string, defaultCurrentTimeString string) (fileAttributes string, creationTime string, lastWriteTime string) {
	if sp == nil {
		return defaultFileAttributes, defaultCurrentTimeString, defaultCurrentTimeString
	}

	fileAttributes = defaultFileAttributes
	if sp.Attributes != nil {
		fileAttributes = sp.Attributes.String()
		if fileAttributes == "" {
			fileAttributes = defaultFileAttributes
		} else if isDir && strings.ToLower(fileAttributes) != "none" {
			// Directories need to have this attribute included, if setting any attributes.
			fileAttributes += "|Directory"
		}
	}

	creationTime = defaultCurrentTimeString
	if sp.CreationTime != nil {
		creationTime = sp.CreationTime.UTC().Format(generated.ISO8601)
	}

	lastWriteTime = defaultCurrentTimeString
	if sp.LastWriteTime != nil {
		lastWriteTime = sp.LastWriteTime.UTC().Format(generated.ISO8601)
	}

	return
}

// NTFSFileAttributes for Files and Directories.
// The subset of attributes is listed at: https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties#file-system-attributes.
type NTFSFileAttributes struct {
	ReadOnly, Hidden, System, Directory, Archive, None, Temporary, Offline, NotContentIndexed, NoScrubData bool
}

// String returns a string representation of NTFSFileAttributes.
func (f *NTFSFileAttributes) String() string {
	fileAttributes := ""
	if f.ReadOnly {
		fileAttributes += "ReadOnly|"
	}
	if f.Hidden {
		fileAttributes += "Hidden|"
	}
	if f.System {
		fileAttributes += "System|"
	}
	if f.Directory {
		fileAttributes += "Directory|"
	}
	if f.Archive {
		fileAttributes += "Archive|"
	}
	if f.None {
		fileAttributes += "None|"
	}
	if f.Temporary {
		fileAttributes += "Temporary|"
	}
	if f.Offline {
		fileAttributes += "Offline|"
	}
	if f.NotContentIndexed {
		fileAttributes += "NotContentIndexed|"
	}
	if f.NoScrubData {
		fileAttributes += "NoScrubData|"
	}

	fileAttributes = strings.TrimSuffix(fileAttributes, "|")
	return fileAttributes
}

// ParseNTFSFileAttributes parses the file attributes from *string to *NTFSFileAttributes.
// It returns an error for any unknown file attribute.
func ParseNTFSFileAttributes(attributes *string) (*NTFSFileAttributes, error) {
	if attributes == nil {
		return nil, nil
	}

	ntfsFileAttributes := NTFSFileAttributes{}
	parts := strings.Split(*attributes, "|")

	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(p))
		switch p {
		case "readonly":
			ntfsFileAttributes.ReadOnly = true
		case "hidden":
			ntfsFileAttributes.Hidden = true
		case "system":
			ntfsFileAttributes.System = true
		case "directory":
			ntfsFileAttributes.Directory = true
		case "archive":
			ntfsFileAttributes.Archive = true
		case "none":
			ntfsFileAttributes.None = true
		case "temporary":
			ntfsFileAttributes.Temporary = true
		case "offline":
			ntfsFileAttributes.Offline = true
		case "notcontentindexed":
			ntfsFileAttributes.NotContentIndexed = true
		case "noscrubdata":
			ntfsFileAttributes.NoScrubData = true
		default:
			return nil, fmt.Errorf("unknown file attribute %v", p)
		}
	}

	return &ntfsFileAttributes, nil
}
