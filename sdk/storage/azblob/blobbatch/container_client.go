//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blobbatch

import (
	"bytes"
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// ContainerBatchClient represents a URL to the Azure Storage container allowing you to manipulate its blobs.
type ContainerBatchClient struct {
	cnt    *container.Client
	policy policy.Policy
}

// NewContainerBatchClient creates an instance of Client with the specified values.
//   - containerURL - the URL of the container e.g. https://<account>.blob.core.windows.net/container
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewContainerBatchClient(containerURL string, cred azcore.TokenCredential, options *ClientOptions) (*ContainerBatchClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ContainerBatchClient{
		cnt:    (*container.Client)(base.NewContainerClient(containerURL, pl, nil)),
		policy: authPolicy,
	}, nil
}

// NewContainerBatchClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a container or with a shared access signature (SAS) token.
//   - containerURL - the URL of the container e.g. https://<account>.blob.core.windows.net/container?<sas token>
//   - options - client options; pass nil to accept the default values
func NewContainerBatchClientWithNoCredential(containerURL string, options *ClientOptions) (*ContainerBatchClient, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ContainerBatchClient{
		cnt: (*container.Client)(base.NewContainerClient(containerURL, pl, nil)),
	}, nil
}

// NewContainerBatchClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - containerURL - the URL of the container e.g. https://<account>.blob.core.windows.net/container
//   - cred - a SharedKeyCredential created with the matching container's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewContainerBatchClientWithSharedKeyCredential(containerURL string, cred *container.SharedKeyCredential, options *ClientOptions) (*ContainerBatchClient, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ContainerBatchClient{
		cnt:    (*container.Client)(base.NewContainerClient(containerURL, pl, nil)),
		policy: authPolicy,
	}, nil
}

// NewContainerBatchClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - containerName - the name of the container within the storage account
//   - options - client options; pass nil to accept the default values
func NewContainerBatchClientFromConnectionString(connectionString string, containerName string, options *ClientOptions) (*ContainerBatchClient, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, containerName)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewContainerBatchClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewContainerBatchClientWithNoCredential(parsed.ServiceURL, options)
}

func (c *ContainerBatchClient) generated() *generated.ContainerClient {
	return base.InnerClient((*base.Client[generated.ContainerClient])(c.cnt))
}

func (c *ContainerBatchClient) sharedKey() *container.SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ContainerClient])(c.cnt))
}

// URL returns the URL endpoint used by the Client object.
func (c *ContainerBatchClient) URL() string {
	return c.generated().Endpoint()
}

func (c *ContainerBatchClient) SubmitBatch(ctx context.Context, bb *BatchBuilder) (ContainerClientSubmitBatchResponse, error) {
	if bb == nil {
		return ContainerClientSubmitBatchResponse{}, errors.New("batch builder is empty")
	}

	batchID, err := shared.CreateBatchID()
	if err != nil {
		return ContainerClientSubmitBatchResponse{}, err
	}

	batchReq, err := bb.createBatchRequest(ctx, c.policy, to.Ptr(c.URL()), &batchID)
	if err != nil {
		return ContainerClientSubmitBatchResponse{}, err
	}

	reader := bytes.NewReader([]byte(batchReq))
	rsc := streaming.NopCloser(reader)
	multipartContentType := "multipart/mixed; boundary=" + batchID
	resp, err := c.generated().SubmitBatch(ctx, int64(len(batchReq)), multipartContentType, rsc, nil)
	return resp, err
}
