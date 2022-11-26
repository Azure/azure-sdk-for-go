//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/url"
	"strings"
)

func getDefaultScope(endpoint string) (string, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.New("error parsing endpoint url")
	}

	return parsedURL.Scheme + "://" + parsedURL.Host + "/.default", nil
}

// ContainerRegistryClientOptions contains the optional parameters for the NewContainerRegistryClient method.
type ContainerRegistryClientOptions struct {
	azcore.ClientOptions
}

// NewContainerRegistryClient creates a new instance of ContainerRegistryClient with the specified values.
//   - endpoint - registry login URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewContainerRegistryClient(endpoint string, credential azcore.TokenCredential, options *ContainerRegistryClientOptions) (*ContainerRegistryClient, error) {
	if options == nil {
		options = &ContainerRegistryClientOptions{}
	}

	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	authClient := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{
		options.ClientOptions,
	})
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
	return &ContainerRegistryClient{
		endpoint,
		pl,
	}, nil
}

func isDigest(reference string) bool {
	return strings.Contains(reference, ":")
}

// GetDigest - Get digest of the manifest.
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - options - GetDigestOptions contains the optional parameters for the ContainerRegistryClient.GetDigest method.
func (client *ContainerRegistryClient) GetDigest(ctx context.Context, name string, reference string, options *GetDigestOptions) (string, error) {
	if isDigest(reference) {
		return reference, nil
	} else {
		resp, err := client.GetTagProperties(ctx, name, reference, nil)
		if err != nil {
			return "", err
		}
		return *resp.Tag.Digest, nil
	}
}

// DeleteManifest - Delete the manifest identified by name and reference.
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - options - ContainerRegistryClientDeleteManifestOptions contains the optional parameters for the ContainerRegistryClient.DeleteManifest
//     method.
func (client *ContainerRegistryClient) DeleteManifest(ctx context.Context, name string, reference string, options *ContainerRegistryClientDeleteManifestOptions) (ContainerRegistryClientDeleteManifestResponse, error) {
	digest, err := client.GetDigest(ctx, name, reference, nil)
	if err != nil {
		return ContainerRegistryClientDeleteManifestResponse{}, err
	}
	return client.deleteManifest(ctx, name, digest, options)
}

// GetManifestProperties - Get manifest properties
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - options - ContainerRegistryClientGetManifestPropertiesOptions contains the optional parameters for the ContainerRegistryClient.GetManifestProperties
//     method.
func (client *ContainerRegistryClient) GetManifestProperties(ctx context.Context, name string, reference string, options *ContainerRegistryClientGetManifestPropertiesOptions) (ContainerRegistryClientGetManifestPropertiesResponse, error) {
	digest, err := client.GetDigest(ctx, name, reference, nil)
	if err != nil {
		return ContainerRegistryClientGetManifestPropertiesResponse{}, err
	}
	return client.getManifestProperties(ctx, name, digest, options)
}

// UpdateManifestProperties - Update properties of a manifest
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - options - ContainerRegistryClientUpdateManifestPropertiesOptions contains the optional parameters for the ContainerRegistryClient.UpdateManifestProperties
//     method.
func (client *ContainerRegistryClient) UpdateManifestProperties(ctx context.Context, name string, reference string, options *ContainerRegistryClientUpdateManifestPropertiesOptions) (ContainerRegistryClientUpdateManifestPropertiesResponse, error) {
	digest, err := client.GetDigest(ctx, name, reference, nil)
	if err != nil {
		return ContainerRegistryClientUpdateManifestPropertiesResponse{}, err
	}
	return client.updateManifestProperties(ctx, name, digest, options)
}

// CalculateDigest - Calculate the digest of a manifest payload
//   - payload - Manifest payload bytes
func CalculateDigest(payload []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(payload))
}
