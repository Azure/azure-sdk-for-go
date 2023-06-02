//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClientFromSharedKey creates a [Client] using a shared key.
func NewClientFromSharedKey(key string, options *ClientOptions) (*Client, error) {
	// TODO: I believe we're supposed to allow for dynamically updating the key at any time as well.
	azc, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			&skpolicy{Key: key},
		},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azc,
	}, nil
}

// TODO: remove in favor of a common policy instead?
type skpolicy struct {
	Key string
}

func (p *skpolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Add("Authorization", "SharedAccessKey "+p.Key)
	return req.Next()
}
