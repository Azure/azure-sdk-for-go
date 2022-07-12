//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type CreateResponse = generated.AppendBlobClientCreateResponse

type AppendBlockResponse = generated.AppendBlobClientAppendBlockResponse

type AppendBlockFromURLResponse = generated.AppendBlobClientAppendBlockFromURLResponse

type SealResponse = generated.AppendBlobClientSealResponse
