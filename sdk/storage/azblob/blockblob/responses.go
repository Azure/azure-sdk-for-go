//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type UploadResponse = generated.BlockBlobClientUploadResponse

type StageBlockResponse = generated.BlockBlobClientStageBlockResponse

type CommitBlockListResponse = generated.BlockBlobClientCommitBlockListResponse

type StageBlockFromURLResponse = generated.BlockBlobClientStageBlockFromURLResponse

type GetBlockListResponse = generated.BlockBlobClientGetBlockListResponse
