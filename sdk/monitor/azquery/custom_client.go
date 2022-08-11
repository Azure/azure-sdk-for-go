//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

// this file contains handwritten additions to the generated code

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a client that accesses a monitor.
func NewClient(credential azcore.TokenCredential, options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://api.loganalytics.io/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &Client{pl: pl}
}
