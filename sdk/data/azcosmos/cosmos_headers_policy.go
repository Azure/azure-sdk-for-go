// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type headerPolicies struct {
	enableContentResponseOnWrite bool
}

type headerOptionsOverride struct {
	enableContentResponseOnWrite *bool
	partitionKey                 *PartitionKey
}

func (p *headerPolicies) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if req.OperationValue(&o) {
		enableContentResponseOnWrite := p.enableContentResponseOnWrite

		if o.headerOptionsOverride != nil {
			if o.headerOptionsOverride.enableContentResponseOnWrite != nil {
				enableContentResponseOnWrite = *o.headerOptionsOverride.enableContentResponseOnWrite
			}

			if o.headerOptionsOverride.partitionKey != nil {
				pkAsString, err := o.headerOptionsOverride.partitionKey.toJsonString()
				if err != nil {
					return nil, err
				}
				req.Raw().Header.Add(cosmosHeaderPartitionKey, string(pkAsString))
			}
		}

		if o.isWriteOperation && !enableContentResponseOnWrite {
			req.Raw().Header.Add(cosmosHeaderPrefer, cosmosHeaderValuesPreferMinimal)
		}
	}

	return req.Next()
}
