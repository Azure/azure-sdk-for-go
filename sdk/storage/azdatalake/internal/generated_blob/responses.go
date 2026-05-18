// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// BlobClientSetExpiryResponse contains the response from method BlobClient.SetExpiry.
type BlobClientSetExpiryResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// ETag contains the information returned from the ETag header response.
	ETag *azcore.ETag

	// LastModified contains the information returned from the Last-Modified header response.
	LastModified *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// ServiceClientListFileSystemsSegmentResponse contains the response from method ServiceClient.NewListContainersSegmentPager.
type ServiceClientListFileSystemsSegmentResponse struct {
	// An enumeration of containers
	ListFileSystemsSegmentResponse

	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}
