//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.ShareClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.ShareClientDeleteResponse

// RestoreResponse contains the response from method Client.Restore.
type RestoreResponse = generated.ShareClientRestoreResponse
