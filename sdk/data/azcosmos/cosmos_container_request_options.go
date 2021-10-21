// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ReadContainerOptions includes options for Read
type ReadContainerOptions struct {
	// PopulateQuotaInfo indicates whether to populate quota info in response headers.
	PopulateQuotaInfo bool
}

func (options *ReadContainerOptions) toHeaders() *map[string]string {
	if !options.PopulateQuotaInfo {
		return nil
	}

	headers := make(map[string]string)
	if options.PopulateQuotaInfo {
		headers[cosmosHeaderPopulateQuotaInfo] = "true"
	}
	return &headers
}

// CreateContainerOptions are options for the CreateContainer operation
type CreateContainerOptions struct {
	// ThroughputProperties: Optional throughput configuration of the container
	ThroughputProperties *ThroughputProperties
}

// ReplaceContainerOptions are options for the ReplaceContainer operation
type ReplaceContainerOptions struct{}

func (options *ReplaceContainerOptions) toHeaders() *map[string]string {
	return nil
}

// DeleteContainerOptions are options for the DeleteContainer operation
type DeleteContainerOptions struct{}

func (options *DeleteContainerOptions) toHeaders() *map[string]string {
	return nil
}
