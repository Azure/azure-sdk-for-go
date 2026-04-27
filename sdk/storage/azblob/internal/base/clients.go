// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// SessionMode contains the possible values for session-based authentication modes.
type SessionMode = exported.SessionMode

const (
	// SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault SessionMode = exported.SessionModeDefault
	// SessionModeOff explicitly disables session-based authentication.
	SessionModeOff SessionMode = exported.SessionModeOff
	// SessionModeSingleSpecifiedContainer enables session-based authentication for a single container.
	SessionModeSingleSpecifiedContainer SessionMode = exported.SessionModeSingleSpecifiedContainer
)

// PossibleSessionModeValues returns a slice of possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return exported.PossibleSessionModeValues()
}

// SessionOptions contains the optional parameters for session-based authentication.
type SessionOptions = exported.SessionOptions

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions

	// Audience to use when requesting tokens for Azure Active Directory authentication.
	// Only has an effect when credential is of type TokenCredential. The value could be
	// https://storage.azure.com/ (default) or https://<account>.blob.core.windows.net.
	Audience string

	// SessionOptions configures session-based authentication behavior.
	// Only has an effect when credential is of type TokenCredential.
	SessionOptions SessionOptions
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

func GetAzClient(storageURL string, cred azcore.TokenCredential, options *ClientOptions) (*azcore.Client, error) {
	audience := GetAudience(options)
	conOptions := shared.GetClientOptions(options)
	authPolicy := shared.NewStorageChallengePolicy(cred, audience, conOptions.InsecureAllowCredentialWithHTTP)
	oauthPlOpts := runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}
	oauthAzClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, oauthPlOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}
	if options == nil || options.SessionOptions.Mode == SessionModeOff || options.SessionOptions.Mode == SessionModeDefault {
		return oauthAzClient, nil
	}
	oAuthServiceClient, err := getServiceClient(storageURL, oauthAzClient, &cred, conOptions)
	if err != nil {
		return nil, err
	}
	sessionPolicy, err := exported.NewSessionPolicy(options.SessionOptions, authPolicy, oAuthServiceClient)
	if err != nil {
		return nil, err
	}

	sessionPlOpts := runtime.PipelineOptions{PerRetry: []policy.Policy{sessionPolicy}}
	sessionAzClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, sessionPlOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}
	return sessionAzClient, nil
}

func getServiceClient(storageURL string, azClient *azcore.Client, cred any, conOptions *ClientOptions) (*generated.ServiceClient, error) {
	serviceURL, err := shared.GetServiceURL(storageURL)
	if err != nil {
		return nil, err
	}
	return InnerClient(NewServiceClient(serviceURL, azClient, &cred, conOptions)), nil
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
