//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"context"
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"strings"
)

type Artifact struct {
	Endpoint                string
	Name                    string
	Reference               string
	FullyQualifiedReference string
	digest                  string
	containerRegistryClient *generated.ContainerRegistryClient
}

// newArtifact creates an Artifact instance.
func newArtifact(endpoint, name, reference string, client *generated.ContainerRegistryClient) Artifact {
	return Artifact{
		Endpoint:                endpoint,
		Name:                    name,
		Reference:               reference,
		containerRegistryClient: client,
	}
}

// GetDigest - Get digest of this artifact.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetDigestOptions contains the optional parameters for the Artifact.GetDigest method.
func (client *Artifact) GetDigest(ctx context.Context, options *GetDigestOptions) (string, error) {
	if client.digest != "" {
		return client.digest, nil
	}
	if client.isDigest() {
		client.digest = client.Reference
	} else {
		resp, err := client.GetTagProperties(ctx, client.Reference, nil)
		if err != nil {
			return "", err
		}
		client.digest = *resp.Tag.Digest
	}

	return client.digest, nil
}

func (client *Artifact) isDigest() bool {
	return strings.Contains(client.Reference, ":")
}

// Delete - Delete this artifact by deleting its manifest.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - DeleteArtifactOptions contains the optional parameters for the Artifact.Delete method.
func (client *Artifact) Delete(ctx context.Context, options *DeleteArtifactOptions) error {
	digest, err := client.GetDigest(ctx, nil)
	if err != nil {
		return err
	}
	_, err = client.containerRegistryClient.DeleteManifest(ctx, client.Name, digest, nil)
	return err
}

// GetManifestProperties - Get manifest attributes of this artifact.
// If the operation fails it returns an *azcore.ResponseError type.
//   - options - GetManifestPropertiesOptions contains the optional parameters for the Artifact.GetManifestProperties method.
func (client *Artifact) GetManifestProperties(ctx context.Context, options *GetManifestPropertiesOptions) (GetManifestPropertiesResponse, error) {
	digest, err := client.GetDigest(ctx, nil)
	if err != nil {
		return GetManifestPropertiesResponse{}, err
	}
	resp, err := client.containerRegistryClient.GetManifestProperties(ctx, client.Name, digest, nil)
	if err != nil {
		return GetManifestPropertiesResponse{}, err
	}
	return resp, nil
}

// UpdateManifestProperties - Update manifest properties of this artifact.
// If the operation fails it returns an *azcore.ResponseError type.
//   - value - Manifest attribute value
//   - options - UpdateManifestPropertiesOptions contains the optional parameters for the Artifact.UpdateManifestProperties method.
func (client *Artifact) UpdateManifestProperties(ctx context.Context, value ManifestWriteableProperties, options *UpdateManifestPropertiesOptions) (UpdateManifestPropertiesResponse, error) {
	digest, err := client.GetDigest(ctx, nil)
	if err != nil {
		return UpdateManifestPropertiesResponse{}, err
	}
	resp, err := client.containerRegistryClient.UpdateManifestProperties(ctx, client.Name, digest, &generated.ContainerRegistryClientUpdateManifestPropertiesOptions{Value: &value})
	if err != nil {
		return UpdateManifestPropertiesResponse{}, err
	}
	return resp, nil
}

// NewListTagsPager - List tags.
//   - options - ListTagsOptions contains the optional parameters for the Artifact.NewListTagsPager method.
func (client *Artifact) NewListTagsPager(options *ListTagsOptions) *runtime.Pager[ListTagsResponse] {
	var requestOption *generated.ContainerRegistryClientListTagsOptions
	if options != nil {
		requestOption = &generated.ContainerRegistryClientListTagsOptions{
			Last:    options.Last,
			N:       options.N,
			Orderby: (*string)(options.OrderBy),
		}
	}
	pagerInternal := client.containerRegistryClient.NewListTagsPager(client.Name, requestOption)
	return runtime.NewPager(runtime.PagingHandler[ListTagsResponse]{
		More: func(ListTagsResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListTagsResponse) (ListTagsResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListTagsResponse{}, err
			}
			return page, nil
		},
	})
}

// GetTagProperties - Get tag attributes by tag name.
// If the operation fails it returns an *azcore.ResponseError type.
//   - tag - Tag name
//   - options - GetTagPropertiesOptions contains the optional parameters for the Artifact.GetTagProperties
//     method.
func (client *Artifact) GetTagProperties(ctx context.Context, tag string, options *GetTagPropertiesOptions) (GetTagPropertiesResponse, error) {
	resp, err := client.containerRegistryClient.GetTagProperties(ctx, client.Name, tag, nil)
	if err != nil {
		return GetTagPropertiesResponse{}, err
	}
	return resp, nil
}

// UpdateTagProperties - Update tag attributes by tag name.
// If the operation fails it returns an *azcore.ResponseError type.
//   - tag - Tag name
//   - value - Tag attribute value
//   - options - UpdateTagPropertiesOptions contains the optional parameters for the Artifact.UpdateTagProperties method.
func (client *Artifact) UpdateTagProperties(ctx context.Context, tag string, value TagWriteableProperties, options *UpdateTagPropertiesOptions) (UpdateTagPropertiesResponse, error) {
	resp, err := client.containerRegistryClient.UpdateTagAttributes(ctx, client.Name, tag, &generated.ContainerRegistryClientUpdateTagAttributesOptions{Value: &value})
	if err != nil {
		return UpdateTagPropertiesResponse{}, err
	}
	return resp, nil
}

// DeleteTag - Delete tag by tag name.
// If the operation fails it returns an *azcore.ResponseError type.
//   - tag - Tag name
//   - options - DeleteTagOptions contains the optional parameters for the Artifact.DeleteTag method.
func (client *Artifact) DeleteTag(ctx context.Context, tag string, options *DeleteTagOptions) error {
	_, err := client.containerRegistryClient.DeleteTag(ctx, client.Name, tag, nil)
	return err
}
