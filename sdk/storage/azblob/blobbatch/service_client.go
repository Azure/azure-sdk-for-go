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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// ServiceBatchClient represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
// It includes auth policy as a parameters which is set when the client is created.
// Executing the auth policy returns the Authorization header which needs to be set in the sub-request.
type ServiceBatchClient struct {
	svc    *service.Client
	policy policy.Policy
}

// NewServiceBatchClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.blob.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewServiceBatchClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*ServiceBatchClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ServiceBatchClient{
		svc:    (*service.Client)(base.NewServiceClient(serviceURL, pl, nil)),
		policy: authPolicy,
	}, nil
}

// NewServiceBatchClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.blob.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewServiceBatchClientWithNoCredential(serviceURL string, options *ClientOptions) (*ServiceBatchClient, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ServiceBatchClient{
		svc: (*service.Client)(base.NewServiceClient(serviceURL, pl, nil)),
	}, nil
}

// NewServiceBatchClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.blob.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewServiceBatchClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*ServiceBatchClient, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return &ServiceBatchClient{
		svc:    (*service.Client)(base.NewServiceClient(serviceURL, pl, nil)),
		policy: authPolicy,
	}, nil
}

// NewServiceBatchClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewServiceBatchClientFromConnectionString(connectionString string, options *ClientOptions) (*ServiceBatchClient, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewServiceBatchClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewServiceBatchClientWithNoCredential(parsed.ServiceURL, options)
}

func (s *ServiceBatchClient) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s.svc))
}

func (s *ServiceBatchClient) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s.svc))
}

// URL returns the URL endpoint used by the Client object.
func (s *ServiceBatchClient) URL() string {
	return s.generated().Endpoint()
}

// SubmitBatch operation allows multiple API calls to be embedded into a single HTTP request.
//   - BatchBuilder - contains the list of operations to be submitted. It supports up to 256 sub-requests in a single batch.
func (s *ServiceBatchClient) SubmitBatch(ctx context.Context, bb *BatchBuilder) (ServiceClientSubmitBatchResponse, error) {
	if bb == nil {
		return ServiceClientSubmitBatchResponse{}, errors.New("batch builder is empty")
	}

	batchID, err := shared.CreateBatchID()
	if err != nil {
		return ServiceClientSubmitBatchResponse{}, err
	}

	// create the request body
	batchReq, err := bb.createBatchRequest(ctx, s.policy, to.Ptr(s.URL()), &batchID)
	if err != nil {
		return ServiceClientSubmitBatchResponse{}, err
	}

	reader := bytes.NewReader([]byte(batchReq))
	rsc := streaming.NopCloser(reader)
	multipartContentType := "multipart/mixed; boundary=" + batchID
	resp, err := s.generated().SubmitBatch(ctx, int64(len(batchReq)), multipartContentType, rsc, nil)

	// TODO: parse the response body to map individual operations to their responses
	return resp, err
}
