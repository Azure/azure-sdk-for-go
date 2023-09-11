//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"errors"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const timeFormat = time.RFC3339Nano

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

func (c *ClientOptions) toConnectionOptions() *policy.ClientOptions {
	if c == nil {
		return nil
	}

	return &policy.ClientOptions{
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}
func getDefaultScope(endpoint string) (string, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.New("error parsing endpoint url")
	}

	return url.Scheme + "://" + url.Host + "/.default", nil
}

// NewClient returns a pointer to a Client object tied to an endpointUrl.
func NewClient(endpointUrl string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	tokenScope, err := getDefaultScope(endpointUrl)
	if err != nil {
		return nil, err
	}

	syncTokenPolicy := newSyncTokenPolicy()

	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)

	azcoreClient, err := azcore.NewClient("azappconfig.Client", moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy, syncTokenPolicy}}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{internal: azcoreClient, endpoint: endpointUrl}, nil
}

// NewClientFromConnectionString parses the connection string and returns a pointer to a Client object.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	endpointUrl, credential, secret, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	syncTokenPolicy := newSyncTokenPolicy()
	authPolicy := newHmacAuthenticationPolicy(credential, secret)

	azcoreClient, err := azcore.NewClient("azappconfig.Client", moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy, syncTokenPolicy}}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{internal: azcoreClient, endpoint: endpointUrl}, nil
}

// UpdateSyncToken sets an external synchronization token to ensure service requests receive up-to-date values.
func (c *Client) UpdateSyncToken(token string) {
	c.syncTokenPolicy.addToken(token)
}
