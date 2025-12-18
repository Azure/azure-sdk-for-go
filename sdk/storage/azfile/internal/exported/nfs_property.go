//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"time"
)

// NFSProperties contains the optional parameters regarding the NFS properties for a file.
type NFSProperties struct {
	// The Coordinated Universal Time (UTC) creation time for the file/directory. Default value is 'now'.
	CreationTime *time.Time
	// The Coordinated Universal Time (UTC) last write time for the file/directory. Default value is 'now'.
	LastWriteTime *time.Time
	// The file mode of the file or directory
	FileMode *string
	// The owner of the file or directory.
	Owner *string
	// The owning group of the file or directory.
	Group *string
}

// FormatNFSProperties returns creation time, last write time.
func FormatNFSProperties(np *NFSProperties, isDir bool) (creationTime *string, lastWriteTime *string) {
	if np == nil {
		return nil, nil
	}

	creationTime = nil
	if np.CreationTime != nil {
		creationTime = to.Ptr(np.CreationTime.UTC().Format(generated.ISO8601))
	}

	lastWriteTime = nil
	if np.LastWriteTime != nil {
		lastWriteTime = to.Ptr(np.LastWriteTime.UTC().Format(generated.ISO8601))
	}

	return
}
