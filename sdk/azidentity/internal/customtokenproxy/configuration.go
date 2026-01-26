// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenproxy

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/exported"
)

const (
	EnvAzureKubernetesCAData     = "AZURE_KUBERNETES_CA_DATA"
	EnvAzureKubernetesCAFile     = "AZURE_KUBERNETES_CA_FILE"
	EnvAzureKubernetesSNIName    = "AZURE_KUBERNETES_SNI_NAME"
	EnvAzureKubernetesTokenProxy = "AZURE_KUBERNETES_TOKEN_PROXY"
)

func readOptionsFromEnv() *exported.CustomTokenProxyOptions {
	return &exported.CustomTokenProxyOptions{
		TokenProxy: os.Getenv(EnvAzureKubernetesTokenProxy),
		SNIName:    os.Getenv(EnvAzureKubernetesSNIName),
		CAFile:     os.Getenv(EnvAzureKubernetesCAFile),
		CAData:     os.Getenv(EnvAzureKubernetesCAData),
	}
}

func backfillOptionsFromEnv(opts *exported.CustomTokenProxyOptions) {
	if opts.CAData != "" || opts.CAFile != "" || opts.SNIName != "" || opts.TokenProxy != "" {
		return
	}

	// only backfill if all fields are empty
	*opts = *readOptionsFromEnv()
}

func parseTokenProxyURL(endpoint string) (*url.URL, error) {
	tokenProxy, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom token proxy URL %q: %s", endpoint, err)
	}
	if tokenProxy.Scheme != "https" {
		return nil, fmt.Errorf("custom token endpoint must use https scheme, got %q", tokenProxy.Scheme)
	}
	if tokenProxy.User != nil {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain user info", tokenProxy)
	}
	if tokenProxy.RawQuery != "" {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain a query", tokenProxy)
	}
	if tokenProxy.EscapedFragment() != "" {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain a fragment", tokenProxy)
	}
	if tokenProxy.EscapedPath() == "" {
		// if the path is empty, set it to "/" to avoid stripping the path from req.URL
		tokenProxy.Path = "/"
	}
	return tokenProxy, nil
}

var (
	errCustomEndpointSetWithoutTokenProxy = errors.New(
		"AZURE_KUBERNETES_TOKEN_PROXY is not set but other custom endpoint-related settings are present",
	)
	errCustomEndpointMultipleCASourcesSet = errors.New(
		"only one of AzureKubernetesCAFile or AzureKubernetesCAData can be specified",
	)
)

func noopConfigure(*policy.ClientOptions) {
	// no-op
}

// GetClientOptionsConfigurer returns a function that configures the client options to use the custom token proxy.
func GetClientOptionsConfigurer(opts *exported.CustomTokenProxyOptions) (func(*policy.ClientOptions), error) {
	if opts == nil {
		return noopConfigure, nil
	}

	backfillOptionsFromEnv(opts)

	if opts.TokenProxy == "" {
		// custom token proxy is not set, while other Kubernetes-related environment variables are present,
		// this is likely a configuration issue so erroring out to avoid misconfiguration
		if opts.SNIName != "" || opts.CAFile != "" || opts.CAData != "" {
			return nil, errCustomEndpointSetWithoutTokenProxy
		}

		return noopConfigure, nil
	}

	tokenProxy, err := parseTokenProxyURL(opts.TokenProxy)
	if err != nil {
		return nil, err
	}

	// CAFile and CAData are mutually exclusive, at most one can be set.
	// If none of CAFile or CAData are set, the default system CA pool will be used.
	if opts.CAFile != "" && opts.CAData != "" {
		return nil, errCustomEndpointMultipleCASourcesSet
	}

	// preload the transport
	t := &transport{
		caFile:     opts.CAFile,
		caData:     []byte(opts.CAData),
		sniName:    opts.SNIName,
		tokenProxy: tokenProxy,
	}
	if _, err := t.getTokenTransporter(); err != nil {
		return nil, err
	}

	return func(clientOptions *policy.ClientOptions) {
		clientOptions.Transport = t
	}, nil
}
