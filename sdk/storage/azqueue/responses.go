//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
)

// CreateQueueResponse contains the response from method queue.Client.Create.
type CreateQueueResponse = generated.QueueClientCreateResponse

// DeleteQueueResponse contains the response from method queue.Client.Delete
type DeleteQueueResponse = generated.QueueClientDeleteResponse

// ListQueuesResponse contains the response from method Client.ListQueuesSegment.
type ListQueuesResponse = generated.ServiceClientListQueuesSegmentResponse

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.ServiceClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method Client.SetProperties.
type SetPropertiesResponse = generated.ServiceClientSetPropertiesResponse

// GetStatisticsResponse contains the response from method Client.GetStatistics.
type GetStatisticsResponse = generated.ServiceClientGetStatisticsResponse
