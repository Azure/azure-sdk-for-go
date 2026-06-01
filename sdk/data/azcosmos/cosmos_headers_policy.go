// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

type headerPolicies struct {
	enableContentResponseOnWrite bool
	priorityLevel                *PriorityLevel
	throughputBucket             *int32
}

type headerOptionsOverride struct {
	enableContentResponseOnWrite *bool
	partitionKey                 *PartitionKey
	correlatedActivityId         *uuid.UUID
	priorityLevel                *PriorityLevel
	throughputBucket             *int32
}

func (p *headerPolicies) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if req.OperationValue(&o) {
		enableContentResponseOnWrite := p.enableContentResponseOnWrite
		priorityLevel := p.priorityLevel
		throughputBucket := p.throughputBucket

		if o.headerOptionsOverride != nil {
			if o.headerOptionsOverride.enableContentResponseOnWrite != nil {
				enableContentResponseOnWrite = *o.headerOptionsOverride.enableContentResponseOnWrite
			}

			if o.headerOptionsOverride.partitionKey != nil && len(o.headerOptionsOverride.partitionKey.values) > 0 {
				pkAsString, err := o.headerOptionsOverride.partitionKey.toJsonString()
				if err != nil {
					return nil, err
				}
				req.Raw().Header.Add(cosmosHeaderPartitionKey, string(pkAsString))
			}

			if o.headerOptionsOverride.correlatedActivityId != nil {
				req.Raw().Header.Add(cosmosHeaderCorrelatedActivityId, o.headerOptionsOverride.correlatedActivityId.String())
			}

			if o.headerOptionsOverride.priorityLevel != nil {
				priorityLevel = o.headerOptionsOverride.priorityLevel
			}

			if o.headerOptionsOverride.throughputBucket != nil {
				throughputBucket = o.headerOptionsOverride.throughputBucket
			}
		}

		if o.isWriteOperation && o.resourceType == resourceTypeDocument && !enableContentResponseOnWrite {
			req.Raw().Header.Add(cosmosHeaderPrefer, cosmosHeaderValuesPreferMinimal)
		}

		if priorityLevel != nil {
			req.Raw().Header.Add(cosmosHeaderPriorityLevel, string(*priorityLevel))
		}

		if throughputBucket != nil {
			req.Raw().Header.Add(cosmosHeaderThroughputBucket, strconv.FormatInt(int64(*throughputBucket), 10))
		}
	}

	return req.Next()
}
