// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

type ListContainersSegmentOptions struct {
	Include ListContainersDetail

	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing operation did not return all containers
	// remaining to be listed with the current page. The NextMarker value can be used as the value for the marker parameter in
	// a subsequent call to request the next page of list items. The marker value is opaque to the client.
	Marker *string

	// Specifies the maximum number of containers to return. If the request does not specify max results, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the listing operation crosses a partition boundary,
	// then the service will return a continuation token for retrieving the remainder of the results. For this reason, it is possible
	// that the service will return fewer results than specified by max results, or than the default of 5000.
	MaxResults *int32

	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
}

func (o *ListContainersSegmentOptions) pointers() *ServiceListContainersSegmentOptions {
	if o == nil {
		return nil
	}

	return &ServiceListContainersSegmentOptions{
		Include:    o.Include.pointers(),
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
		Prefix:     o.Prefix,
	}
}

// ListContainersFlatDetail indicates what additional information the service should return with each container.
type ListContainersDetail struct {
	// Tells the service whether to return metadata for each container.
	Metadata bool

	// Tells the service whether to return soft-deleted containers.
	Deleted bool
}

// string produces the Include query parameter's value.
func (d *ListContainersDetail) pointers() *[]ListContainersIncludeType {
	if d.Metadata == false && d.Deleted == false {
		return nil
	}

	items := make([]ListContainersIncludeType, 0, 2)
	// NOTE: Multiple strings MUST be appended in alphabetic order or signing the string for authentication fails!
	if d.Deleted {
		items = append(items, ListContainersIncludeTypeDeleted)
	}
	if d.Metadata {
		items = append(items, ListContainersIncludeTypeMetadata)
	}
	return &items
}
