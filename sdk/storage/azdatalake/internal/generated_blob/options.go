// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

// ServiceClientListContainersSegmentOptions contains the optional parameters for the ServiceClient.NewListContainersSegmentPager
// method.
type ServiceClientListContainersSegmentOptions struct {
	// Include this parameter to specify that the container's metadata be returned as part of the response body.
	Include []ListContainersIncludeType

	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string

	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	Maxresults *int32

	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://learn.microsoft.com/rest/api/storageservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}
