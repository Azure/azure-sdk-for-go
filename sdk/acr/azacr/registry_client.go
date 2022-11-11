//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"context"
	"errors"
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/url"
	"strings"
)

// RegistryClient represents a client to the Azure Container Registry acrService.
// Don't use this type directly, use NewRegistryClient() instead.
type RegistryClient struct {
	Endpoint                    string
	containerRegistryClient     *generated.ContainerRegistryClient
	containerRegistryBlobClient *generated.ContainerRegistryBlobClient
}

type RegistryClientOptions struct {
	azcore.ClientOptions
}

func getDefaultScope(endpoint string) (string, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.New("error parsing endpoint url")
	}

	return parsedURL.Scheme + "://" + parsedURL.Host + "/.default", nil
}

// NewRegistryClient creates a RegistryClient struct using the...
func NewRegistryClient(endpoint string, credential azcore.TokenCredential, options *RegistryClientOptions) (*RegistryClient, error) {
	if options == nil {
		options = &RegistryClientOptions{}
	}

	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	authPl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &options.ClientOptions)
	authClient := generated.NewAuthenticationClient(endpoint, authPl)
	tokenScope, err := getDefaultScope(endpoint)
	if err != nil {
		return nil, err
	}
	authPolicy := NewAuthenticationPolicy(
		credential,
		[]string{tokenScope},
		authClient,
		nil,
	)

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	containerRegistryClient := generated.NewContainerRegistryClient(endpoint, pl)
	containerRegistryBlobClient := generated.NewContainerRegistryBlobClient(endpoint, pl)
	return &RegistryClient{
		endpoint,
		containerRegistryClient,
		containerRegistryBlobClient,
	}, nil
}

// DeleteRepository - Delete the repository identified by name.
// If the operation fails it returns an *azcore.ResponseError type.
//   - name - Name of the image (including the namespace)
//   - options - DeleteRepositoryOptions contains the optional parameters for the RegistryClient.DeleteRepository method.
func (client *RegistryClient) DeleteRepository(ctx context.Context, name string, options *DeleteRepositoryOptions) error {
	_, err := client.containerRegistryClient.DeleteRepository(ctx, name, nil)
	return err
}

// GetArtifact - Return an instance of Artifact for calling operation related to the artifact.
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
func (client *RegistryClient) GetArtifact(name, reference string) Artifact {
	return newArtifact(client.Endpoint, name, reference, client.containerRegistryClient)
}

// GetRepository - Return an instance of Repository for calling operation related to the repository.
func (client *RegistryClient) GetRepository(name string) Repository {
	return newRepository(client.Endpoint, name, client.containerRegistryClient, client.containerRegistryBlobClient)
}

// NewListRepositoriesPager - List repositories of this registry.
//   - options - ListRepositoriesOptions contains the optional parameters for the RegistryClient.NewListRepositoriesPager method.
func (client *RegistryClient) NewListRepositoriesPager(options *ListRepositoriesOptions) *runtime.Pager[ListRepositoriesResponse] {
	pagerInternal := client.containerRegistryClient.NewListRepositoriesPager(options)
	return runtime.NewPager(runtime.PagingHandler[ListRepositoriesResponse]{
		More: func(ListRepositoriesResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListRepositoriesResponse) (ListRepositoriesResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListRepositoriesResponse{}, err
			}
			return listRepositoriesResponseFromGenerated(page), nil
		},
	})
}
