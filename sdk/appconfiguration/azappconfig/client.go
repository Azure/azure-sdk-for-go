//go:build go1.16
// +build go1.16

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

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

// Client is the struct for interacting with an Azure App Configuration instance.
type Client struct {
	appConfigClient *generated.AzureAppConfigurationClient
	syncTokenPolicy *syncTokenPolicy
}

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

// NewClient returns a pointer to a Client object affinitized to an endpointUrl.
func NewClient(endpointUrl string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := options.toConnectionOptions()

	tokenScope, err := getDefaultScope(endpointUrl)
	if err != nil {
		return nil, err
	}

	syncTokenPolicy := newSyncTokenPolicy()
	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil),
		syncTokenPolicy,
	)

	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)
	return &Client{
		appConfigClient: generated.NewAzureAppConfigurationClient(endpointUrl, nil, pl),
		syncTokenPolicy: syncTokenPolicy,
	}, nil
}

// NewClientFromConnectionString parses the connection string and returns a pointer to a Client object.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := options.toConnectionOptions()

	endpoint, credential, secret, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	syncTokenPolicy := newSyncTokenPolicy()
	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		newHmacAuthenticationPolicy(credential, secret),
		syncTokenPolicy,
	)

	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)
	return &Client{
		appConfigClient: generated.NewAzureAppConfigurationClient(endpoint, nil, pl),
		syncTokenPolicy: syncTokenPolicy,
	}, nil
}

// UpdateSyncToken sets an external synchronization token to ensure service requests receive up-to-date values.
func (c *Client) UpdateSyncToken(token string) {
	c.syncTokenPolicy.addToken(token)
}

const timeFormat = time.RFC3339Nano
