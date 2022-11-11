//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
)

type Repository struct {
	Endpoint                    string
	Name                        string
	containerRegistryClient     *generated.ContainerRegistryClient
	containerRegistryBlobClient *generated.ContainerRegistryBlobClient
}

// newRepository creates a Repository instance.
func newRepository(endpoint, name string, client *generated.ContainerRegistryClient, blobClient *generated.ContainerRegistryBlobClient) Repository {
	return Repository{
		Endpoint:                    endpoint,
		Name:                        name,
		containerRegistryClient:     client,
		containerRegistryBlobClient: blobClient,
	}
}

// Delete - Delete this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - DeleteRepositoryOptions contains the optional parameters for the Repository.DeleteRepository method.
func (client *Repository) Delete(ctx context.Context, options *DeleteRepositoryOptions) error {
	_, err := client.containerRegistryClient.DeleteRepository(ctx, client.Name, nil)
	return err
}

// GetArtifact - Return an instance of Artifact for calling operation related to the artifact.
//   - reference - A tag or a digest, pointing to a specific image
func (client *Repository) GetArtifact(reference string) Artifact {
	return newArtifact(client.Endpoint, client.Name, reference, client.containerRegistryClient)
}

// GetRepositoryProperties - Get attributes of this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetRepositoryPropertiesOptions contains the optional parameters for the Repository.GetRepositoryProperties method.
func (client *Repository) GetRepositoryProperties(ctx context.Context, options *GetRepositoryPropertiesOptions) (GetRepositoryPropertiesResponse, error) {
	resp, err := client.containerRegistryClient.GetProperties(ctx, client.Name, nil)
	if err != nil {
		return GetRepositoryPropertiesResponse{}, err
	}
	return resp, nil
}

// UpdateRepositoryProperties - Update the attribute of this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - value - Repository attribute value
//   - options - UpdatePropertiesOptions contains the optional parameters for the Repository.UpdateRepositoryProperties method.
func (client *Repository) UpdateRepositoryProperties(ctx context.Context, value RepositoryWriteableProperties, options *UpdateRepositoryPropertiesOptions) (UpdateRepositoryPropertiesResponse, error) {
	resp, err := client.containerRegistryClient.UpdateProperties(ctx, client.Name, &generated.ContainerRegistryClientUpdatePropertiesOptions{Value: &value})
	if err != nil {
		return UpdateRepositoryPropertiesResponse{}, err
	}
	return resp, nil
}

// NewListManifestsPager - List manifests of this repository.
//   - options - ListManifestsOptions contains the optional parameters for the Repository.NewListManifestsPager method.
func (client *Repository) NewListManifestsPager(options *ListManifestsOptions) *runtime.Pager[ListManifestsResponse] {
	var requestOption *generated.ContainerRegistryClientListManifestsOptions
	if options != nil {
		requestOption = &generated.ContainerRegistryClientListManifestsOptions{
			Last:    options.Last,
			N:       options.N,
			Orderby: (*string)(options.OrderBy),
		}
	}
	pagerInternal := client.containerRegistryClient.NewListManifestsPager(client.Name, requestOption)
	return runtime.NewPager(runtime.PagingHandler[ListManifestsResponse]{
		More: func(ListManifestsResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListManifestsResponse) (ListManifestsResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListManifestsResponse{}, err
			}
			return page, nil
		},
	})
}

// DeleteManifest - Delete manifest from this repository by digest.
// If the operation fails it returns an *azcore.ResponseError type.
//   - digest - Digest of a manifest
//   - options - DeleteManifestOptions contains the optional parameters for the Repository.DeleteManifest method.
func (client *Repository) DeleteManifest(ctx context.Context, digest string, options *DeleteManifestOptions) error {
	_, err := client.containerRegistryClient.DeleteManifest(ctx, client.Name, digest, nil)
	return err
}

// DeleteBlob - Removes an already uploaded blob of this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - digest - Digest of a BLOB
//   - options - DeleteBlobOptions contains the optional parameters for the Repository.DeleteBlob
//     method.
func (client *Repository) DeleteBlob(ctx context.Context, digest string, options *DeleteBlobOptions) error {
	_, err := client.containerRegistryBlobClient.DeleteBlob(ctx, client.Name, digest, nil)
	return err
}

// UploadOCIManifest - Upload an OCI manifest of this repository with reference where reference can be a tag or digest.
// If the operation fails it returns an *azcore.ResponseError type.
//   - reference - A tag or a digest, pointing to a specific image
//   - payload - OCI manifest file io
//   - options - UploadOCIManifestOptions contains the optional parameters for the Repository.UploadOCIManifest method.
func (client *Repository) UploadOCIManifest(ctx context.Context, reference string, payload io.ReadSeekCloser, options *UploadOCIManifestOptions) (CreateManifestResponse, error) {
	resp, err := client.containerRegistryClient.CreateManifestWithBinary(ctx, client.Name, reference, generated.ContentTypeApplicationVndOciImageManifestV1JSON, payload, nil)
	if err != nil {
		return CreateManifestResponse{}, nil
	}
	return resp, nil
}

// CreateOCIManifest - Create an OCI manifest of this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - manifest - OCI manifest
//   - options - CreateOCIManifestOptions contains the optional parameters for the Repository.UploadManifest method.
func (client *Repository) CreateOCIManifest(ctx context.Context, manifest OCIManifest, options *CreateOCIManifestOptions) (CreateManifestResponse, error) {
	payload, err := json.Marshal(manifest)
	if err != nil {
		return CreateManifestResponse{}, nil
	}
	var tagOrDigest string
	if options.Tag != nil {
		tagOrDigest = *options.Tag
	} else {
		tagOrDigest = fmt.Sprintf("%x", sha256.Sum256(payload))
	}
	resp, err := client.containerRegistryClient.CreateManifestWithBinary(ctx, client.Name, tagOrDigest, generated.ContentTypeApplicationVndOciImageManifestV1JSON, streaming.NopCloser(bytes.NewReader(payload)), nil)
	if err != nil {
		return CreateManifestResponse{}, nil
	}
	return resp, nil
}

// DownloadOCIManifest - Download an OCI manifest of this repository with reference where reference can be a tag or digest.
// If the operation fails it returns an *azcore.ResponseError type.
//   - reference - A tag or a digest, pointing to a specific image
//   - options - DownloadOCIManifestOptions contains the optional parameters for the Repository.DownloadOCIManifest method.
func (client *Repository) DownloadOCIManifest(ctx context.Context, reference string, options *DownloadOCIManifestOptions) (DownloadOCIManifestResponse, error) {
	resp, err := client.containerRegistryClient.GetManifest(ctx, client.Name, reference, &generated.ContainerRegistryClientGetManifestOptions{Accept: to.Ptr(string(generated.ContentTypeApplicationVndOciImageManifestV1JSON))})
	if err != nil {
		return DownloadOCIManifestResponse{}, nil
	}
	manifest := OCIManifest{
		Annotations:   resp.Annotations,
		Config:        resp.Config,
		Layers:        resp.Layers,
		SchemaVersion: resp.SchemaVersion,
	}
	payload, err := json.Marshal(manifest)
	if err != nil {
		return DownloadOCIManifestResponse{}, nil
	}
	return DownloadOCIManifestResponse{
		*resp.DockerContentDigest,
		manifest,
		streaming.NopCloser(bytes.NewReader(payload)),
	}, nil
}

// UploadBlob - Upload a blob to this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - blob - blob file io
//   - options - UploadBlobOptions contains the optional parameters for the Repository.UploadBlob method.
func (client *Repository) UploadBlob(ctx context.Context, blob io.ReadSeekCloser, options *UploadBlobOptions) (UploadBlobResponse, error) {
	startResp, err := client.containerRegistryBlobClient.StartUpload(ctx, client.Name, nil)
	if err != nil {
		return UploadBlobResponse{}, nil
	}
	payload, err := io.ReadAll(blob)
	if err != nil {
		return UploadBlobResponse{}, nil
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	uploadResp, err := client.containerRegistryBlobClient.UploadChunkWithBinary(ctx, *startResp.Location, blob, nil)
	if err != nil {
		return UploadBlobResponse{}, nil
	}
	complateResp, err := client.containerRegistryBlobClient.CompleteUploadWithBinary(ctx, digest, *uploadResp.Location, nil)
	if err != nil {
		return UploadBlobResponse{}, nil
	}
	if digest != *complateResp.DockerContentDigest {
		return UploadBlobResponse{}, errors.New("digest of blob to upload does not match the digest from the server")
	}
	return UploadBlobResponse{Digest: digest}, nil
}

// DownloadBlob - Download a blob from this repository.
// If the operation fails it returns an *azcore.ResponseError type.
//   - digest - digest of the blob
//   - options - DownloadBlobOptions contains the optional parameters for the Repository.DownloadBlob method.
func (client *Repository) DownloadBlob(ctx context.Context, digest string, options *DownloadBlobOptions) (DownloadBlobResponse, error) {
	resp, err := client.containerRegistryBlobClient.GetBlob(ctx, client.Name, digest, nil)
	if err != nil {
		return DownloadBlobResponse{}, nil
	}
	return DownloadBlobResponse{
		Content: resp.Body,
		Digest:  *resp.DockerContentDigest,
	}, nil
}
