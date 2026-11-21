// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/sanitizer"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

// ClientOptions contains configuration settings for a client's pipeline.
type ClientOptions = armpolicy.ClientOptions

// Client is a HTTP client for use with ARM endpoints.  It consists of an endpoint, pipeline, and tracing provider.
type Client struct {
	ep string
	pl runtime.Pipeline
	tr tracing.Tracer
}

// NewClient creates a new Client instance with the provided values.
// This client is intended to be used with Azure Resource Manager endpoints.
//   - moduleName - the fully qualified name of the module where the client is defined; used by the telemetry policy and tracing provider.
//   - moduleVersion - the semantic version of the module; used by the telemetry policy and tracing provider.
//   - cred - the TokenCredential used to authenticate the request
//   - options - optional client configurations; pass nil to accept the default values
func NewClient(moduleName, moduleVersion string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	if !options.Telemetry.Disabled {
		if err := shared.ValidateModVer(moduleVersion); err != nil {
			return nil, err
		}
	}

	cc, err := resolveCloud(options.Cloud)
	if err != nil {
		return nil, err
	}
	// NewPipeline returns an error when cc has no Resource Manager configuration
	ep := cc.Services[cloud.ResourceManager].Endpoint
	// give the pipeline the resolved cloud so its token audience always matches the
	// endpoint. Copy the options to avoid modifying the caller's value
	plOpts := *options
	plOpts.Cloud = cc
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, cred, runtime.PipelineOptions{}, &plOpts)
	if err != nil {
		return nil, err
	}

	tr := options.TracingProvider.NewTracer(moduleName, moduleVersion)
	tr = sanitizer.WrapTracer(tr)
	return &Client{ep: ep, pl: pl, tr: tr}, nil
}

// resolveCloud returns the cloud configuration a client should use. Precedence is:
//  1. c, when it isn't the zero value i.e., the application specified a cloud
//  2. the well-known cloud named by the environment variable AZURE_ENVIRONMENT
//  3. default: Azure Public Cloud
//
// It returns an error when the environment variable has an unrecognized value.
func resolveCloud(c cloud.Configuration) (cloud.Configuration, error) {
	if !reflect.ValueOf(c).IsZero() {
		return c, nil
	}
	// names are case insensitive and include those accepted by earlier Azure SDKs and the Azure CLI
	name := os.Getenv(shared.EnvAzureEnvironment)
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "azurepublic", "azurepubliccloud", "azurecloud":
		return cloud.AzurePublic, nil
	case "azurechina", "azurechinacloud":
		return cloud.AzureChina, nil
	case "azuregovernment", "azureusgovernment", "azureusgovernmentcloud":
		return cloud.AzureGovernment, nil
	default:
		return cloud.Configuration{}, fmt.Errorf(
			"invalid value %q for environment variable %s. Valid values are AzureChina, AzureGovernment, and AzurePublic",
			name, shared.EnvAzureEnvironment,
		)
	}
}

// Endpoint returns the service's base URL for this client.
func (c *Client) Endpoint() string {
	return c.ep
}

// Pipeline returns the pipeline for this client.
func (c *Client) Pipeline() runtime.Pipeline {
	return c.pl
}

// Tracer returns the tracer for this client.
func (c *Client) Tracer() tracing.Tracer {
	return c.tr
}
