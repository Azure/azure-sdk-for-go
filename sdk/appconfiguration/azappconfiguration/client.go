//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"errors"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal/generated"
)

type Client struct {
	appConfigClient *generated.AzureAppConfigurationClient
	syncTokenPolicy *internal.SyncTokenPolicy
}

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

func NewClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := options.toConnectionOptions()

	tokenScope, err := getDefaultScope(endpoint)
	if err != nil {
		return nil, err
	}

	syncTokenPolicy := internal.NewSyncTokenPolicy()
	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil),
		syncTokenPolicy,
	)

	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)
	return &Client{
		appConfigClient: generated.NewAzureAppConfigurationClient(endpoint, nil, pl),
		syncTokenPolicy: syncTokenPolicy,
	}, nil
}

func (c *Client) UpdateSyncToken(token *string) {
	c.syncTokenPolicy.AddToken(*token)
}

const timeFormat = time.RFC3339Nano
