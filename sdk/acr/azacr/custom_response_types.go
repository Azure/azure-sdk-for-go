//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

// UploadBlobResponse contains the response from method ContainerRegistryBlobClient.UploadBlob.
type UploadBlobResponse struct {
	// The blob's digest, calculated by the registry.
	Digest string
}
