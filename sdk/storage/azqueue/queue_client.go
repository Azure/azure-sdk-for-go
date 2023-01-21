//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/shared"
)

// QueueClient represents a URL to the Azure Queue Storage service allowing you to manipulate queues.
type QueueClient base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient]

func (q *QueueClient) generated() *generated.QueueClient {
	queue, _, _ := base.InnerClients((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
	return queue
}

//func (q *QueueClient) messagesClient() *generated.MessagesClient {
//	_, messages, _ := base.InnerClients((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
//	return messages
//}
//
//func (q *QueueClient) messagesIDClient() *generated.MessageIDClient {
//	_, _, mID := base.InnerClients((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
//	return mID
//}
//
//func (q *QueueClient) sharedKey() *SharedKeyCredential {
//	return base.SharedKeyComposite((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
//}

// URL returns the URL endpoint used by the ServiceClient object.
func (q *QueueClient) URL() string {
	return q.generated().Endpoint()
}

// NewQueueClient creates an instance of ServiceClient with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.queue.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewQueueClient(queueURL string, cred azcore.TokenCredential, options *ClientOptions) (*QueueClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*QueueClient)(base.NewQueueClient(queueURL, pl, nil)), nil
}

// NewQueueClientWithNoCredential creates an instance of QueueClient with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.queue.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewQueueClientWithNoCredential(queueURL string, options *ClientOptions) (*QueueClient, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*QueueClient)(base.NewQueueClient(queueURL, pl, nil)), nil
}

// NewQueueClientWithSharedKeyCredential creates an instance of ServiceClient with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.queue.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewQueueClientWithSharedKeyCredential(queueURL string, cred *SharedKeyCredential, options *ClientOptions) (*QueueClient, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*QueueClient)(base.NewQueueClient(queueURL, pl, cred)), nil
}

// NewQueueClientFromConnectionString creates an instance of ServiceClient with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewQueueClientFromConnectionString(connectionString string, queueName string, options *ClientOptions) (*QueueClient, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, queueName)
	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewQueueClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewQueueClientWithNoCredential(parsed.ServiceURL, options)
}

// Create creates a new queue within a storage account. If a queue with the same name already exists, the operation fails.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-queue4.
func (q *QueueClient) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	opts := options.format()
	resp, err := q.generated().Create(ctx, opts)
	return resp, err
}

// Delete deletes the specified queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-queue3.
func (q *QueueClient) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts := options.format()
	resp, err := q.generated().Delete(ctx, opts)
	return resp, err
}
