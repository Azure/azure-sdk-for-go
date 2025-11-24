// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenproxy

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	AzureKubernetesCAData  = "AZURE_KUBERNETES_CA_DATA"
	AzureKubernetesCAFile  = "AZURE_KUBERNETES_CA_FILE"
	AzureKubernetesSNIName = "AZURE_KUBERNETES_SNI_NAME"

	AzureKubernetesTokenProxy = "AZURE_KUBERNETES_TOKEN_PROXY"
)

// Options contains optional parameters for custom token proxy configuration.
type Options struct {
	// AzureKubernetesCAData specifies the CA certificate data for the Kubernetes cluster.
	// Corresponds to the AZURE_KUBERNETES_CA_DATA environment variable.
	// At most one of AzureKubernetesCAData or AzureKubernetesCAFile should be set.
	AzureKubernetesCAData string

	// AzureKubernetesCAFile specifies the path to the CA certificate file for the Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_CA_FILE environment variable.
	// At most one of AzureKubernetesCAData or AzureKubernetesCAFile should be set.
	AzureKubernetesCAFile string

	// AzureKubernetesSNIName specifies the name of the SNI for Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_SNI_NAME environment variable.
	AzureKubernetesSNIName string

	// AzureKubernetesTokenProxy specifies the URL of the custom token proxy for the Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_TOKEN_PROXY environment variable.
	AzureKubernetesTokenProxy string
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

func (o *Options) defaults() {
	if o.AzureKubernetesTokenProxy == "" {
		o.AzureKubernetesTokenProxy = os.Getenv(AzureKubernetesTokenProxy)
	}
	if o.AzureKubernetesSNIName == "" {
		o.AzureKubernetesSNIName = os.Getenv(AzureKubernetesSNIName)
	}
	if o.AzureKubernetesCAFile == "" {
		o.AzureKubernetesCAFile = os.Getenv(AzureKubernetesCAFile)
	}
	if o.AzureKubernetesCAData == "" {
		o.AzureKubernetesCAData = os.Getenv(AzureKubernetesCAData)
	}
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

// Apply returns a function that configures the client options to use the custom token proxy.
func Apply(opts *Options) (func(*policy.ClientOptions), error) {
	if opts == nil {
		return noopConfigure, nil
	}

	opts.defaults()

	if opts.AzureKubernetesTokenProxy == "" {
		// custom token proxy is not set, while other Kubernetes-related environment variables are present,
		// this is likely a configuration issue so erroring out to avoid misconfiguration
		if opts.AzureKubernetesSNIName != "" || opts.AzureKubernetesCAFile != "" || opts.AzureKubernetesCAData != "" {
			return nil, errCustomEndpointSetWithoutTokenProxy
		}

		return noopConfigure, nil
	}

	tokenProxy, err := parseTokenProxyURL(opts.AzureKubernetesTokenProxy)
	if err != nil {
		return nil, err
	}

	// CAFile and CAData are mutually exclusive, at most one can be set.
	// If none of CAFile or CAData are set, the default system CA pool will be used.
	if opts.AzureKubernetesCAFile != "" && opts.AzureKubernetesCAData != "" {
		return nil, errCustomEndpointMultipleCASourcesSet
	}

	// preload the transport
	t := &transport{
		caFile:     opts.AzureKubernetesCAFile,
		caData:     []byte(opts.AzureKubernetesCAData),
		sniName:    opts.AzureKubernetesSNIName,
		tokenProxy: tokenProxy,
	}
	if _, err := t.getTokenTransporter(); err != nil {
		return nil, err
	}

	return func(clientOptions *policy.ClientOptions) {
		clientOptions.Transport = t
	}, nil
}
