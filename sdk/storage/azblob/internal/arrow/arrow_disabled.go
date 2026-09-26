// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:build azblob_noarrow

package arrow

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// Supported reports whether Arrow list responses are compiled in.
// With the azblob_noarrow tag, list operations use the XML format.
const Supported = false

// HandleFlatListResponse is never called when Supported is false.
func HandleFlatListResponse(resp *http.Response) (generated.ContainerClientListBlobFlatSegmentResponse, error) {
	_ = resp.Body.Close()
	return generated.ContainerClientListBlobFlatSegmentResponse{}, ErrNotSupported
}

// HandleHierarchyListResponse is never called when Supported is false.
func HandleHierarchyListResponse(resp *http.Response) (generated.ContainerClientListBlobHierarchySegmentResponse, error) {
	_ = resp.Body.Close()
	return generated.ContainerClientListBlobHierarchySegmentResponse{}, ErrNotSupported
}
