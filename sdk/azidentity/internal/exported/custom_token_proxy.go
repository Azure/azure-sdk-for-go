// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

// CustomTokenProxyOptions contains optional parameters for custom token proxy configuration.
type CustomTokenProxyOptions struct {
	// CAData specifies the CA certificate data for the Kubernetes cluster.
	// Corresponds to the AZURE_KUBERNETES_CA_DATA environment variable.
	// At most one of CAData or CAFile should be set.
	CAData string

	// CAFile specifies the path to the CA certificate file for the Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_CA_FILE environment variable.
	// At most one of CAData or CAFile should be set.
	CAFile string

	// SNIName specifies the name of the SNI for Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_SNI_NAME environment variable.
	SNIName string

	// TokenProxy specifies the URL of the custom token proxy for the Kubernetes cluster.
	// This field corresponds to the AZURE_KUBERNETES_TOKEN_PROXY environment variable.
	TokenProxy string
}
