//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
)

// DeleteRepositoryOptions contains the optional parameters for the RegistryClient.DeleteRepository method and Repository.Delete method.
type DeleteRepositoryOptions = generated.ContainerRegistryClientDeleteRepositoryOptions

// ListRepositoriesOptions contains the optional parameters for the RegistryClient.NewListRepositoriesPager method.
type ListRepositoriesOptions = generated.ContainerRegistryClientListRepositoriesOptions

// Repositories - List of repositories
type Repositories = generated.Repositories

// GetRepositoryPropertiesOptions contains the optional parameters for the Repository.GetRepositoryProperties method.
type GetRepositoryPropertiesOptions = generated.ContainerRegistryClientGetPropertiesOptions

// RepositoryProperties - Properties of this repository.
type RepositoryProperties = generated.ContainerRepositoryProperties

// RepositoryWriteableProperties - Changeable attributes for Repository
type RepositoryWriteableProperties = generated.RepositoryWriteableProperties

// UpdateRepositoryPropertiesOptions contains the optional parameters for the Repository.UpdateRepositoryProperties method.
type UpdateRepositoryPropertiesOptions = generated.ContainerRegistryClientUpdatePropertiesOptions

// ListManifestsOptions contains the optional parameters for the Repository.NewListManifestsPager method.
type ListManifestsOptions = struct {
	// Query parameter for the last item in previous query. Result set will include values lexically after last.
	Last *string
	// query parameter for max number of items
	N *int32
	// order by query parameter
	OrderBy *ManifestOrderBy
}

// Manifests - Manifest attributes
type Manifests = generated.AcrManifests

// ManifestAttributesBase - Manifest details
type ManifestAttributesBase generated.ManifestAttributesBase

// ArtifactManifestPlatform - The artifact's platform, consisting of operating system and architecture.
type ArtifactManifestPlatform = generated.ArtifactManifestPlatform

// DeleteArtifactOptions contains the optional parameters for the Artifact.Delete method.
type DeleteArtifactOptions = generated.ContainerRegistryClientDeleteManifestOptions

// GetManifestPropertiesOptions contains the optional parameters for the Artifact.GetManifestProperties method.
type GetManifestPropertiesOptions = generated.ContainerRegistryClientGetManifestPropertiesOptions

// ManifestProperties - Manifest attributes details
type ManifestProperties = generated.ArtifactManifestProperties

// UpdateManifestPropertiesOptions contains the optional parameters for the Artifact.UpdateManifestProperties method.
type UpdateManifestPropertiesOptions = generated.ContainerRegistryClientUpdateManifestPropertiesOptions

// ManifestWriteableProperties - Changeable attributes
type ManifestWriteableProperties = generated.ManifestWriteableProperties

// GetDigestOptions contains the optional parameters for the Artifact.GetDigest method.
type GetDigestOptions struct {
	// placeholder for future optional parameters
}

// GetTagPropertiesOptions contains the optional parameters for the Artifact.GetTagProperties method.
type GetTagPropertiesOptions = generated.ContainerRegistryClientGetTagPropertiesOptions

// TagProperties - Tag attributes
type TagProperties = generated.ArtifactTagProperties

// TagAttributesBase - Tag attribute details
type TagAttributesBase = generated.TagAttributesBase

// TagWriteableProperties - Changeable attributes
type TagWriteableProperties = generated.TagWriteableProperties

// UpdateTagPropertiesOptions contains the optional parameters for the Artifact.UpdateTagProperties method.
type UpdateTagPropertiesOptions = generated.ContainerRegistryClientUpdateTagAttributesOptions

// ListTagsOptions contains the optional parameters for the Artifact.NewListTagsPager method.
type ListTagsOptions struct {
	// filter by digest
	Digest *string
	// Query parameter for the last item in previous query. Result set will include values lexically after last.
	Last *string
	// query parameter for max number of items
	N *int32
	// order by query parameter
	OrderBy *TagOrderBy
}

// TagList - List of tag details
type TagList = generated.TagList

// DeleteTagOptions contains the optional parameters for the Artifact.DeleteTag method.
type DeleteTagOptions = generated.ContainerRegistryClientDeleteTagOptions

// DeleteBlobOptions contains the optional parameters for the Repository.DeleteBlob method.
type DeleteBlobOptions = generated.ContainerRegistryBlobClientDeleteBlobOptions

// UploadOCIManifestOptions contains the optional parameters for the Repository.UploadOCIManifest method.
type UploadOCIManifestOptions struct {
	// placeholder for future optional parameters
}

// CreateOCIManifestOptions contains the optional parameters for the Repository.CreateOCIManifest method.
type CreateOCIManifestOptions struct {
	// tag of the manifest
	Tag *string
}

// OCIManifest - Returns the requested OCI Manifest file
type OCIManifest = generated.OCIManifest

// Annotations - Additional information provided through arbitrary metadata.
type Annotations = generated.Annotations

// Descriptor - Docker V2 image layer descriptor including config and layers
type Descriptor = generated.Descriptor

// DownloadOCIManifestOptions contains the optional parameters for the Repository.DownloadOCIManifest method.
type DownloadOCIManifestOptions struct {
	// placeholder for future optional parameters
}

// UploadBlobOptions contains the optional parameters for the Repository.UploadBlob method.
type UploadBlobOptions struct {
	// placeholder for future optional parameters
}

// DownloadBlobOptions contains the optional parameters for the Repository.DownloadBlob method.
type DownloadBlobOptions struct {
	// placeholder for future optional parameters
}

// DeleteManifestOptions contains the optional parameters for the Repository.DeleteManifest method.
type DeleteManifestOptions struct {
	// placeholder for future optional parameters
}
