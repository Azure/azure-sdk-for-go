//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"io"
	"reflect"
)

// BlobClientOptions contains the optional parameters for the NewBlobClient method.
type BlobClientOptions struct {
	azcore.ClientOptions
}

// NewBlobClient creates a new instance of BlobClient with the specified values.
//   - endpoint - registry login URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewBlobClient(endpoint string, credential azcore.TokenCredential, options *BlobClientOptions) (*BlobClient, error) {
	if options == nil {
		options = &BlobClientOptions{}
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
	return &BlobClient{
		endpoint,
		pl,
	}, nil
}

// CalculateDigest - Calculate the digest of a payload
//   - payload - Payload io
func CalculateDigest(payload io.ReadSeekCloser) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, payload); err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}
