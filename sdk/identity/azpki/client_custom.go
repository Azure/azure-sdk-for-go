//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azpki

import (
	"errors"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// Client contains the methods for the Client group.
type Client struct {
	endpoint              string
	certificateFormat     *string
	certificateFileFormat *string
	pl                    runtime.Pipeline
}

// ClientOptions - Acceptable values of ClientOptions.
type ClientOptions struct {
	azcore.ClientOptions
	CertificateFormat     *string
	CertificateFileFormat *string
}

// NewClient creates a new instance of Client with the specified values.
func NewClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	c, ok := options.Cloud.Services[ServiceNameLogs]
	if !ok || c.Audience == "" || c.Endpoint == "" {
		return nil, errors.New("provided Cloud field is missing Azure Identity Certificate Manager endpoint configuration")
	}
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{c.Audience + "/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	client := &Client{
		endpoint:              endpoint,
		certificateFormat:     options.CertificateFormat,
		certificateFileFormat: options.CertificateFileFormat,
		pl:                    pl,
	}
	return client, nil
}
