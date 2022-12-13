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
	"io"
	"strings"
)

// ContainerRegistryBlobClientOptions contains the optional parameters for the NewContainerRegistryBlobClient method.
type ContainerRegistryBlobClientOptions struct {
	azcore.ClientOptions
	// Audience is the audience the client will request for its access tokens.
	// The default will connect to Azure public cloud with value "https://management.core.windows.net/".
	Audience string
}

// NewContainerRegistryBlobClient creates a new instance of ContainerRegistryBlobClient with the specified values.
//   - endpoint - registry login URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewContainerRegistryBlobClient(endpoint string, credential azcore.TokenCredential, options *ContainerRegistryBlobClientOptions) (*ContainerRegistryBlobClient, error) {
	if options == nil {
		options = &ContainerRegistryBlobClientOptions{}
	}

	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	authClient := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{
		options.ClientOptions,
	})
	scope := "https://management.core.windows.net/.default"
	if options.Audience != "" {
		scope = options.Audience + "/.default"
	}
	authPolicy := NewAuthenticationPolicy(
		credential,
		[]string{scope},
		authClient,
		nil,
	)

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &ContainerRegistryBlobClient{
		endpoint,
		pl,
	}, nil
}

// UploadBlob - Upload a blob to this repository.
//   - name - Name of the image (including the namespace)
//   - blob - blob file io
//   - options - UploadBlobOptions contains the optional parameters for the ContainerRegistryBlobClient.UploadBlob method.
func (client *ContainerRegistryBlobClient) UploadBlob(ctx context.Context, name string, blob io.ReadSeekCloser, options *UploadBlobOptions) (UploadBlobResponse, error) {
	startResp, err := client.StartUpload(ctx, name, nil)
	if err != nil {
		return UploadBlobResponse{}, err
	}
	payload, err := io.ReadAll(blob)
	if err != nil {
		return UploadBlobResponse{}, err
	}
	digest := CalculateDigest(payload)
	uploadResp, err := client.UploadChunk(ctx, *startResp.Location, blob, nil)
	if err != nil {
		return UploadBlobResponse{}, err
	}
	completeResp, err := client.CompleteUpload(ctx, digest, *uploadResp.Location, nil)
	if err != nil {
		return UploadBlobResponse{}, err
	}
	if digest != *completeResp.DockerContentDigest {
		return UploadBlobResponse{}, errors.New("digest of blob to upload does not match the digest from the server")
	}
	return UploadBlobResponse{Digest: digest}, nil
}

// CalculateDigest - Calculate the digest of a manifest payload
//   - payload - Manifest payload bytes
func CalculateDigest(payload []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(payload))
}
