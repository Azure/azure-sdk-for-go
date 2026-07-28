// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions

	// Audience to use when requesting tokens for Azure Active Directory authentication.
	// Only has an effect when credential is of type TokenCredential. The value could be
	// https://storage.azure.com/ (default) or https://<account>.blob.core.windows.net.
	Audience string

	// ExpectContinueBehavior configures the application of the HTTP "Expect: 100-continue"
	// header on operations that include a request body. The default zero-value behavior
	// conditionally applies the header for a short window after the service responds with a
	// throttle/server-error status (429, 500, or 503).
	//
	// Setting the environment variable AZURE_STORAGE_DISABLE_EXPECT_CONTINUE_HEADER to a
	// truthy value disables this behavior entirely, regardless of this setting.
	ExpectContinueBehavior exported.ExpectContinueOptions

	// Session configures session-based authentication behavior.
	Session exported.SessionOptions
}

type Client[T any] struct {
	inner      *T
	credential any
	options    *ClientOptions
}

func InnerClient[T any](client *Client[T]) *T {
	return client.inner
}

func SharedKey[T any](client *Client[T]) *exported.SharedKeyCredential {
	switch cred := client.credential.(type) {
	case *exported.SharedKeyCredential:
		return cred
	default:
		return nil
	}
}

func Credential[T any](client *Client[T]) any {
	return client.credential
}

func GetClientOptions[T any](client *Client[T]) *ClientOptions {
	return client.options
}

func GetAudience(clOpts *ClientOptions) string {
	if clOpts == nil || len(strings.TrimSpace(clOpts.Audience)) == 0 {
		return shared.TokenScope
	} else {
		return strings.TrimRight(clOpts.Audience, "/") + "/.default"
	}
}

// GetAzClient creates an *azcore.Client with the common pipeline configuration used by all blob clients.
// Provide either cred or sharedKey for authentication, or both nil for anonymous/SAS access.
func GetAzClient(cred azcore.TokenCredential, sharedKey *exported.SharedKeyCredential, conOptions *ClientOptions) (*azcore.Client, error) {
	var plOpts runtime.PipelineOptions
	if cred != nil {
		audience := GetAudience(conOptions)
		bearerTokenPolicy := shared.NewStorageChallengePolicy(cred, audience, conOptions.InsecureAllowCredentialWithHTTP)
		var authPolicy policy.Policy
		if conOptions.Session.Mode == exported.SessionModeDisabled || conOptions.Session.Mode == exported.SessionModeDefault {
			authPolicy = bearerTokenPolicy
		} else if conOptions.Session.Mode == exported.SessionModeEnabled {
			// Session Provider
			var provider exported.SessionProvider
			if conOptions.Session.Provider == nil {
				// TODO : Build Session Provider for customer
			} else {
				provider = conOptions.Session.Provider
			}
			authPolicy = exported.NewSessionPolicy(provider, bearerTokenPolicy)
		} else {
			return nil, fmt.Errorf("unsupported session mode %v", conOptions.Session.Mode)
		}
		plOpts.PerRetry = []policy.Policy{authPolicy}
	} else if sharedKey != nil {
		authPolicy := exported.NewSharedKeyCredPolicy(sharedKey)
		plOpts.PerRetry = []policy.Policy{authPolicy}
	}
	if p := NewExpectContinuePolicy(conOptions.ExpectContinueBehavior); p != nil {
		plOpts.PerRetry = append(plOpts.PerRetry, p)
	}
	return azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
}

func NewClient[T any](inner *T) *Client[T] {
	return &Client[T]{inner: inner}
}

func NewServiceClient(containerURL string, azClient *azcore.Client, credential any, options *ClientOptions) *Client[generated.ServiceClient] {
	return &Client[generated.ServiceClient]{
		inner:      generated.NewServiceClient(containerURL, azClient),
		credential: credential,
		options:    options,
	}
}

func NewContainerClient(containerURL string, azClient *azcore.Client, credential any, options *ClientOptions) *Client[generated.ContainerClient] {
	return &Client[generated.ContainerClient]{
		inner:      generated.NewContainerClient(containerURL, azClient),
		credential: credential,
		options:    options,
	}
}

func NewBlobClient(blobURL string, azClient *azcore.Client, credential any, options *ClientOptions) *Client[generated.BlobClient] {
	return &Client[generated.BlobClient]{
		inner:      generated.NewBlobClient(blobURL, azClient),
		credential: credential,
		options:    options,
	}
}

type CompositeClient[T, U any] struct {
	innerT    *T
	innerU    *U
	sharedKey *exported.SharedKeyCredential
}

func InnerClients[T, U any](client *CompositeClient[T, U]) (*Client[T], *U) {
	return &Client[T]{
		inner:      client.innerT,
		credential: client.sharedKey,
	}, client.innerU
}

func NewAppendBlobClient(blobURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential) *CompositeClient[generated.BlobClient, generated.AppendBlobClient] {
	return &CompositeClient[generated.BlobClient, generated.AppendBlobClient]{
		innerT:    generated.NewBlobClient(blobURL, azClient),
		innerU:    generated.NewAppendBlobClient(blobURL, azClient),
		sharedKey: sharedKey,
	}
}

func NewBlockBlobClient(blobURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential) *CompositeClient[generated.BlobClient, generated.BlockBlobClient] {
	return &CompositeClient[generated.BlobClient, generated.BlockBlobClient]{
		innerT:    generated.NewBlobClient(blobURL, azClient),
		innerU:    generated.NewBlockBlobClient(blobURL, azClient),
		sharedKey: sharedKey,
	}
}

func NewPageBlobClient(blobURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential) *CompositeClient[generated.BlobClient, generated.PageBlobClient] {
	return &CompositeClient[generated.BlobClient, generated.PageBlobClient]{
		innerT:    generated.NewBlobClient(blobURL, azClient),
		innerU:    generated.NewPageBlobClient(blobURL, azClient),
		sharedKey: sharedKey,
	}
}

func SharedKeyComposite[T, U any](client *CompositeClient[T, U]) *exported.SharedKeyCredential {
	return client.sharedKey
}
