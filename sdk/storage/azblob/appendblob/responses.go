//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package appendblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type CreateResponse = generated.AppendBlobClientCreateResponse

type AppendBlockResponse = generated.AppendBlobClientAppendBlockResponse
