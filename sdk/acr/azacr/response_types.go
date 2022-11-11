//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
	"io"
)

// ListRepositoriesResponse contains the response from method RegistryClient.ListRepositories.
type ListRepositoriesResponse struct {
	Repositories
}

func listRepositoriesResponseFromGenerated(generated generated.ContainerRegistryClientListRepositoriesResponse) ListRepositoriesResponse {
	return ListRepositoriesResponse{
		Repositories: generated.Repositories,
	}
}

// GetRepositoryPropertiesResponse contains the response from method Repository.GetRepositoryProperties.
type GetRepositoryPropertiesResponse = generated.ContainerRegistryClientGetPropertiesResponse

// UpdateRepositoryPropertiesResponse contains the response from method Repository.UpdateRepositoryProperties.
type UpdateRepositoryPropertiesResponse = generated.ContainerRegistryClientUpdatePropertiesResponse

// ListManifestsResponse contains the response from method Repository.NewListManifestsPager.
type ListManifestsResponse = generated.ContainerRegistryClientListManifestsResponse

// GetManifestPropertiesResponse contains the response from method Artifact.GetManifestProperties.
type GetManifestPropertiesResponse = generated.ContainerRegistryClientGetManifestPropertiesResponse

// UpdateManifestPropertiesResponse contains the response from method Artifact.UpdateManifestProperties.
type UpdateManifestPropertiesResponse = generated.ContainerRegistryClientUpdateManifestPropertiesResponse

// GetTagPropertiesResponse contains the response from method Artifact.GetTagProperties.
type GetTagPropertiesResponse = generated.ContainerRegistryClientGetTagPropertiesResponse

// UpdateTagPropertiesResponse contains the response from method Artifact.UpdateTagProperties.
type UpdateTagPropertiesResponse = generated.ContainerRegistryClientUpdateTagAttributesResponse

// ListTagsResponse contains the response from method Artifact.NewListTagsPager.
type ListTagsResponse = generated.ContainerRegistryClientListTagsResponse

// CreateManifestResponse contains the response from method Repository.UploadManifest.
type CreateManifestResponse = generated.ContainerRegistryClientCreateManifestWithBinaryResponse

// DownloadOCIManifestResponse contains the response from method Repository.DownloadOCIManifest.
type DownloadOCIManifestResponse struct {
	// The digest of the downloaded manifest as calculated by the registry.
	Digest string

	// The OCI manifest that was downloaded.
	OCIManifest OCIManifest

	// The manifest stream that was downloaded.
	OCIManifestStream io.ReadCloser
}

// UploadBlobResponse contains the response from method Repository.UploadBlob.
type UploadBlobResponse struct {
	// The blob's digest, calculated by the registry.
	Digest string
}

// DownloadBlobResponse contains the response from method Repository.DownloadBlob.
type DownloadBlobResponse struct {
	// The blob content.
	Content io.ReadCloser

	// The blob's digest, calculated by the registry.
	Digest string
}
