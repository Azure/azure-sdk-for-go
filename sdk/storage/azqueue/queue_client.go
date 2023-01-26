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

func (q *QueueClient) messagesClient() *generated.MessagesClient {
	_, messages, _ := base.InnerClients((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
	return messages
}

func (q *QueueClient) messagesIDClient() *generated.MessageIDClient {
	_, _, mID := base.InnerClients((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
	return mID
}

func (q *QueueClient) sharedKey() *SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.QueueClient, generated.MessagesClient, generated.MessageIDClient])(q))
}

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

// SetMetadata sets the metadata for the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-queue-metadata.
func (q *QueueClient) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts := options.format()
	resp, err := q.generated().SetMetadata(ctx, opts)
	return resp, err
}

// GetProperties gets properties including metadata of a queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-queue-metadata.
func (q *QueueClient) GetProperties(ctx context.Context, options *GetQueuePropertiesOptions) (GetQueuePropertiesResponse, error) {
	opts := options.format()
	resp, err := q.generated().GetProperties(ctx, opts)
	return resp, err
}

// GetAccessPolicy returns the queue's access policy.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-queue-acl.
func (q *QueueClient) GetAccessPolicy(ctx context.Context, o *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	options := o.format()
	resp, err := q.generated().GetAccessPolicy(ctx, options)
	return resp, err
}

// SetAccessPolicy sets the queue's permissions.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-queue-acl.
func (q *QueueClient) SetAccessPolicy(ctx context.Context, o *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	opts, acl, err := o.format()
	if err != nil {
		return SetAccessPolicyResponse{}, err
	}
	resp, err := q.generated().SetAccessPolicy(ctx, acl, opts)
	return resp, err
}

// EnqueueMessage adds a message to the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/put-message.
func (q *QueueClient) EnqueueMessage(ctx context.Context, content string, o *EnqueueMessageOptions) (EnqueueMessagesResponse, error) {
	opts := o.format()
	message := generated.QueueMessage{MessageText: &content}
	resp, err := q.messagesClient().Enqueue(ctx, message, opts)
	return resp, err
}

// DequeueMessage removes one message from the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-messages.
func (q *QueueClient) DequeueMessage(ctx context.Context, o *DequeueMessageOptions) (DequeueMessagesResponse, error) {
	opts := o.format()
	resp, err := q.messagesClient().Dequeue(ctx, opts)
	return resp, err
}

// UpdateMessage updates a message from the queue with the given ID.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/update-message.
func (q *QueueClient) UpdateMessage(ctx context.Context, messageID string, content string, o *UpdateMessageOptions) (UpdateMessageResponse, error) {
	opts := o.format()
	message := generated.QueueMessage{MessageText: &content}
	resp, err := q.messagesIDClient().Update(ctx, messageID, message, opts)
	return resp, err
}

// DeleteMessage deletes message from queue with the given ID.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-message2.
func (q *QueueClient) DeleteMessage(ctx context.Context, messageID string, o *DeleteMessageOptions) (DeleteMessageResponse, error) {
	opts := o.format()
	resp, err := q.messagesIDClient().Delete(ctx, messageID, opts)
	return resp, err
}

// PeekMessage peeks the first message from the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/peek-messages.
func (q *QueueClient) PeekMessage(ctx context.Context, o *PeekMessageOptions) (PeekMessagesResponse, error) {
	opts := o.format()
	resp, err := q.messagesClient().Peek(ctx, opts)
	return resp, err
}

// DequeueMessages removes one or more messages from the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-messages.
func (q *QueueClient) DequeueMessages(ctx context.Context, o *DequeueMessagesOptions) (DequeueMessagesResponse, error) {
	opts := o.format()
	resp, err := q.messagesClient().Dequeue(ctx, opts)
	return resp, err
}

// PeekMessages peeks one or more messages from the queue
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/peek-messages.
func (q *QueueClient) PeekMessages(ctx context.Context, o *PeekMessagesOptions) (PeekMessagesResponse, error) {
	opts := o.format()
	resp, err := q.messagesClient().Peek(ctx, opts)
	return resp, err
}

// ClearMessages deletes all messages from the queue.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/clear-messages.
func (q *QueueClient) ClearMessages(ctx context.Context, o *ClearMessagesOptions) (ClearMessagesResponse, error) {
	opts := o.format()
	resp, err := q.messagesClient().Clear(ctx, opts)
	return resp, err
}
