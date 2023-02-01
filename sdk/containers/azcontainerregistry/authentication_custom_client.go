//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// authenticationClientOptions contains the optional parameters for the newAuthenticationClient method.
type authenticationClientOptions struct {
	azcore.ClientOptions
}

// newAuthenticationClient creates a new instance of AuthenticationClient with the specified values.
//   - endpoint - Registry login URL
//   - options - Client options, pass nil to accept the default values.
func newAuthenticationClient(endpoint string, options *authenticationClientOptions) *authenticationClient {
	if options == nil {
		options = &authenticationClientOptions{}
	}

	pipeline := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &options.ClientOptions)

	client := &authenticationClient{
		endpoint: endpoint,
		pl:       pipeline,
	}
	return client
}
