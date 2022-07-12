//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

const (
	CountToEnd = 0

	SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"

	// DefaultDownloadBlockSize is default block size
	DefaultDownloadBlockSize = int64(4 * 1024 * 1024) // 4MB
)
