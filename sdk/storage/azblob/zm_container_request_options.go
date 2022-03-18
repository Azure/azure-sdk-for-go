// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

// ---------------------------------------------------------------------------------------------------------------------

// CreateContainerOptions provides set of configurations for CreateContainer operation
type CreateContainerOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Specifies the encryption scope settings to set on the container.
	cpkScope *ContainerCpkScopeInfo
}

func (o *CreateContainerOptions) pointers() (*containerClientCreateOptions, *ContainerCpkScopeInfo) {
	if o == nil {
		return nil, nil
	}

	basicOptions := containerClientCreateOptions{
		Access:   o.Access,
		Metadata: o.Metadata,
	}

	return &basicOptions, o.cpkScope
}

// ContainerCreateResponse is wrapper around containerClientCreateResponse
type ContainerCreateResponse struct {
	containerClientCreateResponse
}

func toContainerCreateResponse(resp containerClientCreateResponse) ContainerCreateResponse {
	return ContainerCreateResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteContainerOptions provides set of configurations for DeleteContainer operation
type DeleteContainerOptions struct {
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *DeleteContainerOptions) pointers() (*containerClientDeleteOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return nil, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

type ContainerDeleteResponse struct {
	containerClientDeleteResponse
}

func toContainerDeleteResponse(resp containerClientDeleteResponse) ContainerDeleteResponse {
	return ContainerDeleteResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesContainerOptions provides set of configurations for GetPropertiesContainer operation
type GetPropertiesContainerOptions struct {
	containerClientGetPropertiesOptions
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesContainerOptions) pointers() (*containerClientGetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &containerClientGetPropertiesOptions{RequestID: o.RequestID, Timeout: o.Timeout}, o.LeaseAccessConditions
}

type ContainerGetPropertiesResponse struct {
	containerClientGetPropertiesResponse
}

func toContainerGetPropertiesResponse(resp containerClientGetPropertiesResponse) ContainerGetPropertiesResponse {
	return ContainerGetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataContainerOptions provides set of configurations for SetMetadataContainer operation
type SetMetadataContainerOptions struct {
	Metadata                 map[string]string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *SetMetadataContainerOptions) pointers() (*containerClientSetMetadataOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return &containerClientSetMetadataOptions{Metadata: o.Metadata}, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

type ContainerSetMetadataResponse struct {
	containerClientSetMetadataResponse
}

func toContainerSetMetadataResponse(resp containerClientSetMetadataResponse) ContainerSetMetadataResponse {
	return ContainerSetMetadataResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions provides set of configurations for GetAccessPolicy operation
type GetAccessPolicyOptions struct {
	containerClientGetAccessPolicyOptions
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) pointers() (*containerClientGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &containerClientGetAccessPolicyOptions{RequestID: o.RequestID, Timeout: o.Timeout}, o.LeaseAccessConditions
}

type ContainerGetAccessPolicyResponse struct {
	containerClientGetAccessPolicyResponse
}

func toContainerGetAccessPolicyResponse(resp containerClientGetAccessPolicyResponse) ContainerGetAccessPolicyResponse {
	return ContainerGetAccessPolicyResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions provides set of configurations for SetAccessPolicy operation
type SetAccessPolicyOptions struct {
	AccessConditions *ContainerAccessConditions
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType
	// the acls for the container
	ContainerACL []*SignedIdentifier

	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string
	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

func (o *SetAccessPolicyOptions) pointers() (*containerClientSetAccessPolicyOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}
	mac, lac := o.AccessConditions.pointers()
	return &containerClientSetAccessPolicyOptions{
		Access:       o.Access,
		ContainerACL: o.ContainerACL,
		RequestID:    o.RequestID,
		Timeout:      o.Timeout}, lac, mac
}

type ContainerSetAccessPolicyResponse struct {
	containerClientSetAccessPolicyResponse
}

func toContainerSetAccessPolicyResponse(resp containerClientSetAccessPolicyResponse) ContainerSetAccessPolicyResponse {
	return ContainerSetAccessPolicyResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ContainerListBlobFlatSegmentOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListBlobsIncludeItem
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
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

func (o *ContainerListBlobFlatSegmentOptions) pointers() *containerClientListBlobFlatSegmentOptions {
	return &containerClientListBlobFlatSegmentOptions{
		Include:    o.Include,
		Marker:     o.Marker,
		Maxresults: o.Maxresults,
		Prefix:     o.Prefix,
		RequestID:  o.Prefix,
		Timeout:    o.Timeout,
	}
}

type ContainerListBlobFlatSegmentPager struct {
	*containerClientListBlobFlatSegmentPager
}

func toContainerListBlobFlatSegmentPager(resp *containerClientListBlobFlatSegmentPager) *ContainerListBlobFlatSegmentPager {
	return &ContainerListBlobFlatSegmentPager{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ContainerListBlobHierarchySegmentOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListBlobsIncludeItem
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
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

func (o *ContainerListBlobHierarchySegmentOptions) pointers() *containerClientListBlobHierarchySegmentOptions {
	return &containerClientListBlobHierarchySegmentOptions{
		Include:    o.Include,
		Marker:     o.Marker,
		Maxresults: o.Maxresults,
		Prefix:     o.Prefix,
		RequestID:  o.RequestID,
		Timeout:    o.Timeout,
	}
}

type ContainerListBlobHierarchySegmentPager struct {
	*containerClientListBlobHierarchySegmentPager
}

func toContainerListBlobHierarchySegmentPager(resp *containerClientListBlobHierarchySegmentPager) *ContainerListBlobHierarchySegmentPager {
	return &ContainerListBlobHierarchySegmentPager{resp}
}
