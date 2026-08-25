// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
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
func GetAzClient(serviceURL string, cred azcore.TokenCredential, sharedKey *exported.SharedKeyCredential, conOptions *ClientOptions) (*azcore.Client, error) {
	var plOpts runtime.PipelineOptions

	// A session is signed with a session key obtained via a token credential, so asking for one
	// without a token credential is a configuration error rather than something to silently ignore.
	if conOptions.Session.Mode == exported.SessionModeEnabled && cred == nil {
		log.Writef(exported.EventSession, "session authentication cannot be enabled: no token credential was provided; session-based authentication requires a TokenCredential.")
		return nil, errors.New("session mode is enabled but no token credential was provided; session-based authentication requires a TokenCredential")
	}

	if cred != nil {
		audience := GetAudience(conOptions)
		bearerTokenPolicy := shared.NewStorageChallengePolicy(cred, audience, conOptions.InsecureAllowCredentialWithHTTP)
		var authPolicy policy.Policy
		switch conOptions.Session.Mode {
		case exported.SessionModeDefault, exported.SessionModeDisabled:
			authPolicy = bearerTokenPolicy
		case exported.SessionModeEnabled:
			// Session Provider
			var provider exported.SessionProvider
			var accountName string
			if conOptions.Session.Provider == nil {
				svcURL, err := shared.GetServiceURL(serviceURL)
				if err != nil {
					// In the future, we will make session default enabled. When the caller didn't
					// explicitly ask for sessions, a setup failure isn't a configuration error,
					// so warn and authenticate with bearer tokens instead.
					if conOptions.Session.Mode == exported.SessionModeDefault {
						log.Writef(exported.EventSession, "session authentication disabled: the service URL could not be determined: %v. Falling back to bearer token authentication.", err)
						authPolicy = bearerTokenPolicy
						break
					}
					log.Writef(exported.EventSession, "session authentication cannot be enabled: the service URL could not be determined: %v.", err)
					return nil, fmt.Errorf("session mode is enabled but service URL could not be determined: %w", err)
				}
				p, err := NewContainerSessionProvider(cred, svcURL, conOptions)
				if err != nil {
					if conOptions.Session.Mode == exported.SessionModeDefault {
						log.Writef(exported.EventSession, "session authentication disabled: the default session provider could not be created: %v. Falling back to bearer token authentication.", err)
						authPolicy = bearerTokenPolicy
						break
					}
					log.Writef(exported.EventSession, "session authentication cannot be enabled: the default session provider could not be created: %v.", err)
					return nil, fmt.Errorf("failed to create default session provider: %w", err)
				}
				provider = p
			} else {
				provider = conOptions.Session.Provider
			}
			if conOptions.Session.AccountName == "" {
				name, err := shared.GetAccountName(serviceURL)
				if err != nil {
					if conOptions.Session.Mode == exported.SessionModeDefault {
						log.Writef(exported.EventSession, "session authentication disabled: the account name could not be determined from the URL: %v. Falling back to bearer token authentication. Set ClientOptions.Session.AccountName to use sessions.", err)
						authPolicy = bearerTokenPolicy
						break
					}
					log.Writef(exported.EventSession, "session authentication cannot be enabled: the account name could not be determined from the URL: %v. Set ClientOptions.Session.AccountName to use sessions.", err)
					return nil, fmt.Errorf("session mode is enabled but account name could not be determined from URL: %w. Please explicitly pass in options.Session.AccountName", err)
				}
				accountName = name
			} else {
				accountName = conOptions.Session.AccountName
			}
			authPolicy = exported.NewSessionPolicy(accountName, provider, bearerTokenPolicy)
		default:
			log.Writef(exported.EventSession, "session authentication cannot be enabled: unsupported session mode %v.", conOptions.Session.Mode)
			return nil, fmt.Errorf("unsupported session mode %v", conOptions.Session.Mode)
		}
		plOpts.PerRetry = []policy.Policy{authPolicy}
	} else if sharedKey != nil {
		authPolicy := exported.NewSharedKeyCredPolicy(sharedKey)
		plOpts.PerRetry = []policy.Policy{authPolicy}
	}
	// The layout policy is registered as a per-call policy because its rewrite must happen
	// exactly once. It moves the account host into the Host header and replaces the URL host
	// with the layout endpoint; those mutations are made on the request itself, so they
	// naturally persist across retries and every attempt is still sent to the layout endpoint.
	// Running it per-retry would re-apply the rewrite to an already-rewritten request, copying
	// the layout host into the Host header and losing the original account host.
	plOpts.PerCall = []policy.Policy{shared.NewLayoutPolicy()}
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
