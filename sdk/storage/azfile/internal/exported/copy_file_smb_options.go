//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"time"
)

// CopyFileCreationTime specifies either the option to copy file creation time from a source file(source) to a target file or
// a time value in ISO 8601 format to set as creation time on a target file.
type CopyFileCreationTime interface {
	FormatCreationTime() *string
	notPubliclyImplementable()
}

// SourceCopyFileCreationTime specifies to copy file creation time from a source file(source) to a target file.
type SourceCopyFileCreationTime struct {
}

func (s SourceCopyFileCreationTime) FormatCreationTime() *string {
	return to.Ptr("source")
}

func (s SourceCopyFileCreationTime) notPubliclyImplementable() {}

// DestinationCopyFileCreationTime specifies a time value in ISO 8601 format to set as creation time on a target file.
type DestinationCopyFileCreationTime time.Time

func (d DestinationCopyFileCreationTime) FormatCreationTime() *string {
	return to.Ptr(time.Time(d).UTC().Format(generated.ISO8601))
}

func (d DestinationCopyFileCreationTime) notPubliclyImplementable() {}

// ---------------------------------------------------------------------------------------------------------------------

// CopyFileLastWriteTime specifies either the option to copy file last write time from a source file(source) to a target file or
// a time value in ISO 8601 format to set as last write time on a target file.
type CopyFileLastWriteTime interface {
	FormatLastWriteTime() *string
	notPubliclyImplementable()
}

// SourceCopyFileLastWriteTime specifies to copy file last write time from a source file(source) to a target file.
type SourceCopyFileLastWriteTime struct {
}

func (s SourceCopyFileLastWriteTime) FormatLastWriteTime() *string {
	return to.Ptr("source")
}

func (s SourceCopyFileLastWriteTime) notPubliclyImplementable() {}

// DestinationCopyFileLastWriteTime specifies a time value in ISO 8601 format to set as last write time on a target file.
type DestinationCopyFileLastWriteTime time.Time

func (d DestinationCopyFileLastWriteTime) FormatLastWriteTime() *string {
	return to.Ptr(time.Time(d).UTC().Format(generated.ISO8601))
}

func (d DestinationCopyFileLastWriteTime) notPubliclyImplementable() {}

// ---------------------------------------------------------------------------------------------------------------------

// CopyFileAttributes specifies either the option to copy file attributes from a source file(source) to a target file or
// a list of attributes to set on a target file.
type CopyFileAttributes interface {
	FormatAttributes() *string
	notPubliclyImplementable()
}

// SourceCopyFileAttributes specifies to copy file attributes from a source file(source) to a target file
type SourceCopyFileAttributes struct {
}

func (s SourceCopyFileAttributes) FormatAttributes() *string {
	return to.Ptr("source")
}

func (s SourceCopyFileAttributes) notPubliclyImplementable() {}

// DestinationCopyFileAttributes specifies a list of attributes to set on a target file.
type DestinationCopyFileAttributes NTFSFileAttributes

func (d DestinationCopyFileAttributes) FormatAttributes() *string {
	attributes := NTFSFileAttributes(d)
	return to.Ptr(attributes.String())
}

func (d DestinationCopyFileAttributes) notPubliclyImplementable() {}
