// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ContainerRequestOptions includes options for operations against a database.
type ContainerRequestOptions struct {
	PopulateQuotaInfo bool
}

func (options *ContainerRequestOptions) toHeaders() *map[string]string {
	if !options.PopulateQuotaInfo {
		return nil
	}

	headers := make(map[string]string)
	if options.PopulateQuotaInfo {
		headers[cosmosHeaderPopulateQuotaInfo] = "true"
	}
	return &headers
}
