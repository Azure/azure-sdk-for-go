// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
)

// Addressing describes how an operation requires a file or directory client to address its resource.
type Addressing int

const (
	// PathOnly operations require the client to address the resource by its path.
	PathOnly Addressing = iota

	// PathOrFileID operations support a client that addresses the resource either by its path or by its file ID.
	PathOrFileID

	// FileIDOnly operations require the client to address the resource by its file ID.
	FileIDOnly
)

// CheckAddressing returns an error if the operation cannot be called on a client with the given URL, so that the
// operation fails before a request is sent to the service.
func CheckAddressing(resourceURL string, operation string, addressing Addressing) error {
	urlParts, err := sas.ParseURL(resourceURL)
	if err != nil {
		// An unparsable URL is reported by the request itself, as it was before file ID support was added.
		return nil
	}

	if urlParts.FileID == "" {
		if addressing == FileIDOnly {
			return fmt.Errorf("%s: %w", operation, fileerror.ErrFileIDRequired)
		}
		return nil
	}

	if addressing == PathOnly {
		return fmt.Errorf("%s: %w", operation, fileerror.ErrFileIDNotSupported)
	}
	if urlParts.DirectoryOrFilePath != "" {
		// A file ID identifies the resource within the share on its own, e.g. a client navigated to with
		// directory.Client.NewFileClient from a client that addresses a directory by its file ID.
		return fmt.Errorf("%s: the URL must not contain both a path and a file ID: %w", operation, fileerror.ErrFileIDNotSupported)
	}
	return nil
}
