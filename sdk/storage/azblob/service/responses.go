//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// CreateContainerResponse contains the response from method ContainerClient.Create.
type CreateContainerResponse = generated.ContainerClientCreateResponse

// DeleteContainerResponse contains the response from method ContainerClient.Delete.
type DeleteContainerResponse = generated.ContainerClientDeleteResponse

// GetAccountInfoResponse contains the response from method ServiceClient.GetAccountInfo.
type GetAccountInfoResponse = generated.ServiceClientGetAccountInfoResponse

// ListContainersResponse contains the response from method ServiceClient.ListContainersSegment.
type ListContainersResponse = generated.ServiceClientListContainersSegmentResponse

// GetPropertiesResponse contains the response from method ServiceClient.GetProperties.
type GetPropertiesResponse = generated.ServiceClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method ServiceClient.SetProperties.
type SetPropertiesResponse = generated.ServiceClientSetPropertiesResponse

// GetStatisticsResponse contains the response from method ServiceClient.GetStatistics.
type GetStatisticsResponse = generated.ServiceClientGetStatisticsResponse

type FilterBlobsResponse = generated.ServiceClientFilterBlobsResponse
