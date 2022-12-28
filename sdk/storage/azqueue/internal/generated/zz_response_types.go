//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package generated

import "time"

// MessageIDClientDeleteResponse contains the response from method MessageIDClient.Delete.
type MessageIDClientDeleteResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// MessageIDClientUpdateResponse contains the response from method MessageIDClient.Update.
type MessageIDClientUpdateResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// PopReceipt contains the information returned from the x-ms-popreceipt header response.
	PopReceipt *string

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// TimeNextVisible contains the information returned from the x-ms-time-next-visible header response.
	TimeNextVisible *time.Time

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// MessagesClientClearResponse contains the response from method MessagesClient.Clear.
type MessagesClientClearResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// MessagesClientDequeueResponse contains the response from method MessagesClient.Dequeue.
type MessagesClientDequeueResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// The object returned when calling Get Messages on a Queue
	QueueMessagesList []*DequeuedMessageItem `xml:"QueueMessage"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// MessagesClientEnqueueResponse contains the response from method MessagesClient.Enqueue.
type MessagesClientEnqueueResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// The object returned when calling Put Message on a Queue
	QueueMessagesList []*EnqueuedMessage `xml:"QueueMessage"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// MessagesClientPeekResponse contains the response from method MessagesClient.Peek.
type MessagesClientPeekResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// The object returned when calling Peek Messages on a Queue
	QueueMessagesList []*PeekedMessageItem `xml:"QueueMessage"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// QueueClientCreateResponse contains the response from method QueueClient.Create.
type QueueClientCreateResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// QueueClientDeleteResponse contains the response from method QueueClient.Delete.
type QueueClientDeleteResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// QueueClientGetAccessPolicyResponse contains the response from method QueueClient.GetAccessPolicy.
type QueueClientGetAccessPolicyResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// a collection of signed identifiers
	SignedIdentifiers []*SignedIdentifier `xml:"SignedIdentifier"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// QueueClientGetPropertiesResponse contains the response from method QueueClient.GetProperties.
type QueueClientGetPropertiesResponse struct {
	// ApproximateMessagesCount contains the information returned from the x-ms-approximate-messages-count header response.
	ApproximateMessagesCount *int32

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// Metadata contains the information returned from the x-ms-meta header response.
	Metadata map[string]string

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// QueueClientSetAccessPolicyResponse contains the response from method QueueClient.SetAccessPolicy.
type QueueClientSetAccessPolicyResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// QueueClientSetMetadataResponse contains the response from method QueueClient.SetMetadata.
type QueueClientSetMetadataResponse struct {
	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// ServiceClientGetPropertiesResponse contains the response from method ServiceClient.GetProperties.
type ServiceClientGetPropertiesResponse struct {
	StorageServiceProperties
	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// ServiceClientGetStatisticsResponse contains the response from method ServiceClient.GetStatistics.
type ServiceClientGetStatisticsResponse struct {
	StorageServiceStats
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// ServiceClientListQueuesSegmentResponse contains the response from method ServiceClient.ListQueuesSegment.
type ServiceClientListQueuesSegmentResponse struct {
	ListQueuesSegmentResponse
	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// ServiceClientSetPropertiesResponse contains the response from method ServiceClient.SetProperties.
type ServiceClientSetPropertiesResponse struct {
	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

