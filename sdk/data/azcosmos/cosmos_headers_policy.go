// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

type headerPolicies struct {
	enableContentResponseOnWrite bool
}

type headerOptionsOverride struct {
	enableContentResponseOnWrite *bool
	partitionKey                 *PartitionKey
	correlatedActivityId         *uuid.UUID
	collectionRID                string
	partitionKeyRangeID          string
	effectivePartitionKey        string
}

func (p *headerPolicies) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if req.OperationValue(&o) {
		enableContentResponseOnWrite := p.enableContentResponseOnWrite

		if resTypeStr := o.resourceType.String(); resTypeStr != "" {
			req.Raw().Header.Set(cosmosHeaderResourceType, resTypeStr)
		}

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

			if o.headerOptionsOverride.collectionRID != "" {
				req.Raw().Header.Set(cosmosHeaderCollectionRid, o.headerOptionsOverride.collectionRID)
			}

			if o.headerOptionsOverride.partitionKeyRangeID != "" {
				req.Raw().Header.Set(cosmosHeaderPartitionKeyRangeId, o.headerOptionsOverride.partitionKeyRangeID)
			}

			if o.headerOptionsOverride.effectivePartitionKey != "" {
				req.Raw().Header.Set(cosmosHeaderEffectivePartitionKey, o.headerOptionsOverride.effectivePartitionKey)
			}
		}

		if o.isWriteOperation && o.resourceType == resourceTypeDocument && !enableContentResponseOnWrite {
			req.Raw().Header.Add(cosmosHeaderPrefer, cosmosHeaderValuesPreferMinimal)
		}
	}

	return req.Next()
}
