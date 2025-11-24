// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ContainerProperties represents the properties of a container.
type ContainerProperties struct {
	// ID contains the unique id of the container.
	ID string
	// ETag contains the entity etag of the container.
	ETag *azcore.ETag
	// SelfLink contains the self-link of the container.
	SelfLink string
	// ResourceID contains the resource id of the container.
	ResourceID string
	// LastModified contains the last modified time of the container.
	LastModified time.Time
	// DefaultTimeToLive contains the default time to live in seconds for items in the container.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/time-to-live#time-to-live-configurations
	DefaultTimeToLive *int32
	// AnalyticalStoreTimeToLiveInSeconds contains the default time to live in seconds for analytical store in the container.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/analytical-store-introduction#analytical-ttl
	AnalyticalStoreTimeToLiveInSeconds *int32
	// PartitionKeyDefinition contains the partition key definition of the container.
	PartitionKeyDefinition PartitionKeyDefinition
	// IndexingPolicy contains the indexing definition of the container.
	IndexingPolicy *IndexingPolicy
	// UniqueKeyPolicy contains the unique key policy of the container.
	UniqueKeyPolicy *UniqueKeyPolicy
	// ConflictResolutionPolicy contains the conflict resolution policy of the container.
	ConflictResolutionPolicy *ConflictResolutionPolicy
	// VectorEmbeddingPolicy contains the vector embedding policy of the container.
	// This policy defines how vector embeddings are stored and searched within the container.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/nosql/vector-search
	VectorEmbeddingPolicy *VectorEmbeddingPolicy
	// FullTextPolicy contains the full-text policy of the container.
	// This policy defines how text properties are indexed for full-text search operations.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/gen-ai/full-text-search
	FullTextPolicy *FullTextPolicy
}

// MarshalJSON implements the json.Marshaler interface
func (tp ContainerProperties) MarshalJSON() ([]byte, error) {
	pkDefinition, err := json.Marshal(tp.PartitionKeyDefinition)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"id\":\"%s\"", tp.ID)

	if tp.ResourceID != "" {
		buffer.WriteString(fmt.Sprintf(",\"_rid\":\"%s\"", tp.ResourceID))
	}

	if tp.ETag != nil {
		buffer.WriteString(",\"_etag\":")
		etag, err := json.Marshal(tp.ETag)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	if tp.SelfLink != "" {
		fmt.Fprintf(buffer, ",\"_self\":\"%s\"", tp.SelfLink)
	}

	if !tp.LastModified.IsZero() {
		buffer.WriteString(fmt.Sprintf(",\"_ts\":%v", strconv.FormatInt(tp.LastModified.Unix(), 10)))
	}

	buffer.WriteString(",\"partitionKey\":")
	buffer.Write(pkDefinition)

	if tp.DefaultTimeToLive != nil {
		buffer.WriteString(fmt.Sprintf(",\"defaultTtl\":%v", *tp.DefaultTimeToLive))
	}

	if tp.AnalyticalStoreTimeToLiveInSeconds != nil {
		buffer.WriteString(fmt.Sprintf(",\"analyticalStorageTtl\":%v", *tp.AnalyticalStoreTimeToLiveInSeconds))
	}

	if tp.IndexingPolicy != nil {
		indexingPolicy, err := json.Marshal(tp.IndexingPolicy)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"indexingPolicy\":")
		buffer.Write(indexingPolicy)
	}

	if tp.UniqueKeyPolicy != nil {
		uniquePolicy, err := json.Marshal(tp.UniqueKeyPolicy)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"uniqueKeyPolicy\":")
		buffer.Write(uniquePolicy)
	}

	if tp.ConflictResolutionPolicy != nil {
		conflictPolicy, err := json.Marshal(tp.ConflictResolutionPolicy)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"conflictResolutionPolicy\":")
		buffer.Write(conflictPolicy)
	}

	if tp.VectorEmbeddingPolicy != nil {
		vectorPolicy, err := json.Marshal(tp.VectorEmbeddingPolicy)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"vectorEmbeddingPolicy\":")
		buffer.Write(vectorPolicy)
	}

	if tp.FullTextPolicy != nil {
		fullTextPolicy, err := json.Marshal(tp.FullTextPolicy)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"fullTextPolicy\":")
		buffer.Write(fullTextPolicy)
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (tp *ContainerProperties) UnmarshalJSON(b []byte) error {
	var attributes map[string]json.RawMessage
	err := json.Unmarshal(b, &attributes)
	if err != nil {
		return err
	}

	if id, ok := attributes["id"]; ok {
		if err := json.Unmarshal(id, &tp.ID); err != nil {
			return err
		}
	}

	if etag, ok := attributes["_etag"]; ok {
		if err := json.Unmarshal(etag, &tp.ETag); err != nil {
			return err
		}
	}

	if self, ok := attributes["_self"]; ok {
		if err := json.Unmarshal(self, &tp.SelfLink); err != nil {
			return err
		}
	}

	if rid, ok := attributes["_rid"]; ok {
		if err := json.Unmarshal(rid, &tp.ResourceID); err != nil {
			return err
		}
	}

	if ttl, ok := attributes["defaultTtl"]; ok {
		if err := json.Unmarshal(ttl, &tp.DefaultTimeToLive); err != nil {
			return err
		}
	}

	if analyticalTtl, ok := attributes["analyticalStorageTtl"]; ok {
		if err := json.Unmarshal(analyticalTtl, &tp.AnalyticalStoreTimeToLiveInSeconds); err != nil {
			return err
		}
	}

	if ts, ok := attributes["_ts"]; ok {
		var timestamp int64
		if err := json.Unmarshal(ts, &timestamp); err != nil {
			return err
		}
		tp.LastModified = time.Unix(timestamp, 0)
	}

	if pk, ok := attributes["partitionKey"]; ok {
		if err := json.Unmarshal(pk, &tp.PartitionKeyDefinition); err != nil {
			return err
		}
	}

	if ip, ok := attributes["indexingPolicy"]; ok {
		if err := json.Unmarshal(ip, &tp.IndexingPolicy); err != nil {
			return err
		}
	}

	if up, ok := attributes["uniqueKeyPolicy"]; ok {
		if err := json.Unmarshal(up, &tp.UniqueKeyPolicy); err != nil {
			return err
		}
	}

	if cp, ok := attributes["conflictResolutionPolicy"]; ok {
		if err := json.Unmarshal(cp, &tp.ConflictResolutionPolicy); err != nil {
			return err
		}
	}

	if vp, ok := attributes["vectorEmbeddingPolicy"]; ok {
		if err := json.Unmarshal(vp, &tp.VectorEmbeddingPolicy); err != nil {
			return err
		}
	}

	if fp, ok := attributes["fullTextPolicy"]; ok {
		if err := json.Unmarshal(fp, &tp.FullTextPolicy); err != nil {
			return err
		}
	}

	return nil
}
