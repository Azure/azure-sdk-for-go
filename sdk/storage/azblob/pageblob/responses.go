//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type CreateResponse = generated.PageBlobClientCreateResponse

type UploadPagesResponse = generated.PageBlobClientUploadPagesResponse

type UploadPagesFromURLResponse = generated.PageBlobClientUploadPagesFromURLResponse

type ClearPagesResponse = generated.PageBlobClientClearPagesResponse

type GetPageRangesResponse = generated.PageBlobClientGetPageRangesResponse

type GetPageRangesDiffResponse = generated.PageBlobClientGetPageRangesDiffResponse

type ResizeResponse = generated.PageBlobClientResizeResponse

type UpdateSequenceNumberResponse = generated.PageBlobClientUpdateSequenceNumberResponse

type CopyIncrementalResponse = generated.PageBlobClientCopyIncrementalResponse
