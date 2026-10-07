// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// NewContainerSessionProvider creates a ContainerSessionProvider for the blob service endpoint of
// serviceURL, which may also be a container or blob URL; its path and query are discarded.
// Sessions are created through a pipeline built from options, with sessions disabled so that
// creating a session can't itself require one.
func NewContainerSessionProvider(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*exported.ContainerSessionProvider, error) {
	if cred == nil {
		return nil, errors.New("a token credential is required to create sessions")
	}
	svcURL, err := shared.GetServiceURL(serviceURL)
	if err != nil {
		return nil, err
	}
	var opts ClientOptions
	if options != nil {
		opts = *options
	}
	opts.Session = exported.SessionOptions{Mode: exported.SessionModeDisabled}
	azClient, err := GetAzClient(svcURL, cred, nil, &opts)
	if err != nil {
		return nil, err
	}
	return exported.NewContainerSessionProviderFromClient(generated.NewServiceClient(svcURL, azClient)), nil
}

// newSessionAuthPolicy returns the authentication policy for a token credential client: the
// session policy wrapping bearerTokenPolicy when sessions resolve to enabled, otherwise
// bearerTokenPolicy itself.
//
// The account name signs session requests. It comes from Session.AccountName or, failing that,
// from clientURL. When neither yields one, an explicit SessionModeEnabled is a configuration
// error, while a mode that merely resolves to enabled leaves the client on bearer tokens.
func newSessionAuthPolicy(clientURL string, cred azcore.TokenCredential, bearerTokenPolicy policy.Policy, options *ClientOptions) (policy.Policy, error) {
	sessionOpts := options.Session
	if err := exported.ValidateSessionMode(sessionOpts.Mode); err != nil {
		return nil, err
	}
	if exported.ResolveSessionMode(sessionOpts.Mode) != exported.SessionModeEnabled {
		return bearerTokenPolicy, nil
	}
	explicit := sessionOpts.Mode == exported.SessionModeEnabled

	accountName := sessionOpts.AccountName
	if accountName == "" {
		name, err := shared.GetAccountName(clientURL)
		if err != nil {
			// The URL is reported without its query, which could hold a SAS.
			endpoint := redactedURL(clientURL)
			if explicit {
				log.Writef(exported.EventSession, "Session authentication cannot be enabled: the account name could not be determined from the URL %s. Set ClientOptions.Session.AccountName to use sessions.", endpoint)
				return nil, fmt.Errorf("session authentication is enabled but the account name could not be determined from the URL %s; set ClientOptions.Session.AccountName when using a custom endpoint URL", endpoint)
			}
			log.Writef(exported.EventSession, "Session authentication disabled: the account name could not be determined from the URL %s. Falling back to bearer token authentication. Set ClientOptions.Session.AccountName to use sessions.", endpoint)
			return bearerTokenPolicy, nil
		}
		accountName = name
	}

	provider := sessionOpts.Provider
	if provider == nil {
		p, err := NewContainerSessionProvider(clientURL, cred, options)
		if err != nil {
			return nil, err
		}
		provider = p
	}
	return exported.NewSessionPolicy(accountName, provider, bearerTokenPolicy), nil
}

// redactedURL returns rawURL without its query and fragment, for use in logs and errors.
func redactedURL(rawURL string) string {
	if u, err := shared.ParseURLWithoutQuery(rawURL); err == nil {
		return u
	}
	return "<unparseable URL>"
}
