// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// CosmosContainerRequestOptions includes options for operations against a database.
type CosmosContainerRequestOptions struct {
	PopulateQuotaInfo bool
}

func (options *CosmosContainerRequestOptions) toHeaders() *map[string]string {
	if !options.PopulateQuotaInfo {
		return nil
	}

	headers := make(map[string]string)
	if options.PopulateQuotaInfo {
		headers[cosmosHeaderPopulateQuotaInfo] = "true"
	}
	return &headers
}
