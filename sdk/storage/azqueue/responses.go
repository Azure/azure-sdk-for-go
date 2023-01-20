//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
)

// CreateQueueResponse contains the response from method queue.ServiceClient.Create.
type CreateQueueResponse = generated.QueueClientCreateResponse

// DeleteQueueResponse contains the response from method queue.ServiceClient.Delete
type DeleteQueueResponse = generated.QueueClientDeleteResponse

// ListQueuesResponse contains the response from method ServiceClient.ListQueuesSegment.
type ListQueuesResponse = generated.ServiceClientListQueuesSegmentResponse

// GetPropertiesResponse contains the response from method ServiceClient.GetProperties.
type GetPropertiesResponse = generated.ServiceClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method ServiceClient.SetProperties.
type SetPropertiesResponse = generated.ServiceClientSetPropertiesResponse

// GetStatisticsResponse contains the response from method ServiceClient.GetStatistics.
type GetStatisticsResponse = generated.ServiceClientGetStatisticsResponse

//------------------------------------------ QUEUES -------------------------------------------------------------------

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.QueueClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.QueueClientDeleteResponse
