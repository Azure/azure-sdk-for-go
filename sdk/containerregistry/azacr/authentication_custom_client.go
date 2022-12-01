//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"strings"
)

// AuthenticationClientOptions contains the optional parameters for the NewAuthenticationClient method.
type AuthenticationClientOptions struct {
	azcore.ClientOptions
}

// NewAuthenticationClient creates a new instance of AuthenticationClient with the specified values.
//   - endpoint - Registry login URL
//   - options - Client options, pass nil to accept the default values.
func NewAuthenticationClient(endpoint string, options *AuthenticationClientOptions) *AuthenticationClient {
	if options == nil {
		options = &AuthenticationClientOptions{}
	}

	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	pipeline := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &options.ClientOptions)

	client := &AuthenticationClient{
		endpoint: endpoint,
		pl:       pipeline,
	}
	return client
}
