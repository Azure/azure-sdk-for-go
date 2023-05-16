//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"io"
	"reflect"
	"strings"
)

// ClientOptions contains the optional parameters for the NewClient method.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a new instance of Client with the specified values.
//   - endpoint - registry login URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	c, ok := options.Cloud.Services[ServiceName]
	if !ok || c.Audience == "" {
		return nil, errors.New("provided Cloud field is missing Azure Container Registry configuration")
	}

	authClient := newAuthenticationClient(endpoint, &authenticationClientOptions{
		options.ClientOptions,
	})
	authPolicy := newAuthenticationPolicy(
		credential,
		[]string{c.Audience + "/.default"},
		authClient,
		nil,
	)

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &Client{
		endpoint,
		pl,
	}, nil
}

func extractNextLink(value string) string {
	return value[1:strings.Index(value, ">")]
}

// GetManifest - Get the manifest identified by name and reference where reference can be a tag or digest.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-07-01
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - options - ClientGetManifestOptions contains the optional parameters for the Client.GetManifest method.
func (client *Client) GetManifest(ctx context.Context, name string, reference string, options *ClientGetManifestOptions) (ClientGetManifestResponse, error) {
	resp, err := client.getManifest(ctx, name, reference, options)
	if err != nil {
		return resp, err
	}
	payload, err := io.ReadAll(resp.ManifestData)
	_ = resp.ManifestData.Close()
	if err != nil {
		return ClientGetManifestResponse{}, err
	}
	payloadDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(payload))
	if strings.HasPrefix(reference, "sha256:") {
		if reference != payloadDigest {
			return ClientGetManifestResponse{}, fmt.Errorf("retrieved manifest digest %s does not match required digest %s", payloadDigest, reference)
		}
	}
	if *resp.DockerContentDigest != payloadDigest {
		return ClientGetManifestResponse{}, fmt.Errorf("retrieved manifest digest %s does not match server-computed digest %s", payloadDigest, *resp.DockerContentDigest)
	}
	resp.ManifestData = io.NopCloser(bytes.NewReader(payload))
	return resp, err
}

// UploadManifest - Put the manifest identified by name and reference where reference can be a tag or digest.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-07-01
//   - name - Name of the image (including the namespace)
//   - reference - A tag or a digest, pointing to a specific image
//   - contentType - Upload file type
//   - manifestData - Manifest body, can take v1 or v2 values depending on accept header
//   - options - ClientUploadManifestOptions contains the optional parameters for the Client.UploadManifest method.
func (client *Client) UploadManifest(ctx context.Context, name string, reference string, contentType ContentType, manifestData io.ReadSeekCloser, options *ClientUploadManifestOptions) (ClientUploadManifestResponse, error) {
	_, err := manifestData.Seek(0, io.SeekStart)
	if err != nil {
		return ClientUploadManifestResponse{}, err
	}
	payload, err := io.ReadAll(manifestData)
	if err != nil {
		return ClientUploadManifestResponse{}, err
	}
	clientDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(payload))
	_, err = manifestData.Seek(0, io.SeekStart)
	if err != nil {
		return ClientUploadManifestResponse{}, err
	}
	if strings.HasPrefix(reference, "sha256:") {
		if reference != clientDigest {
			return ClientUploadManifestResponse{}, fmt.Errorf("client-computed manifest digest %s does not match required digest %s", clientDigest, reference)
		}
	}
	resp, err := client.uploadManifest(ctx, name, reference, contentType, manifestData, options)
	if err != nil {
		return resp, err
	}
	if *resp.DockerContentDigest != clientDigest {
		return ClientUploadManifestResponse{}, fmt.Errorf("client-computed manifest digest %s does not match server-computed digest %s", clientDigest, *resp.DockerContentDigest)
	}
	return resp, err
}
