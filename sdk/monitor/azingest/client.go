//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest

// this file contains handwritten additions to the generated code

import (
	"bytes"
	"compress/gzip"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest/internal/generated"
)

type Client struct {
	Client *generated.Client
}

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a client to upload logs to Azure Monitor Ingestion.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://monitor.azure.com/" + "/.default"}, nil)
	azcoreClient, err := azcore.NewClient("azingest.Client", version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{Client: &generated.Client{Internal: azcoreClient, Endpoint: endpoint}}, nil
}

// UploadOptions contains the optional parameters for the Client.Upload method.
type UploadOptions struct {
	// gzip
	IsGziped *bool
}

type UploadResponse struct {
}

// might need to make smaller, total needs to be less than 1 mb
const maxUploadBytes int = 1024 * 1024

// fix return type
// Upload
func (client *Client) Upload(ctx context.Context, ruleID string, streamName string, logs []byte, options *UploadOptions) (generated.ClientUploadResponse, error) {

	// if already gzipped, send
	if options != nil && *options.IsGziped {
		return client.Client.Upload(ctx, ruleID, streamName, logs, &generated.ClientUploadOptions{ContentEncoding: to.Ptr("gzip")})
	}

	// check if gzipped bytes are under 1 mb
	// if true, send
	// if not cut in half and call again.

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	length, err := gz.Write(logs)
	if err != nil {
		return generated.ClientUploadResponse{}, err
	}

	// do they gzip or batch first??

	// gzip logs if not already done
	// check if gzipped twice?

	//batching algorithm
	// reference azblob

	// if logs are less than 1 mb, upload in one API call
	if len(logs) <= maxUploadBytes {
		return client.Client.Upload(ctx, ruleID, streamName, logs, to.Ptr(uploadOptions))
	}

	return generated.ClientUploadResponse{}, nil

	// else, batch []byte into smaller chunks then send to Monitor
}
