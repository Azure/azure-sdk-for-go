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

// PartitionKeyRange represents the properties of a partition key range.
type PartitionKeyRange struct {
	// ID contains the unique id of the partition key range.
	ID string `json:"id"`
	// ResourceID contains the resource id of the partition key range.
	ResourceID string `json:"_rid"`
	// ETag contains the entity etag of the partition key range.
	ETag *azcore.ETag `json:"_etag"`
	// MinInclusive contains the minimum inclusive value of the partition key range.
	MinInclusive string `json:"minInclusive"`
	// MaxExclusive contains the maximum exclusive value of the partition key range.
	MaxExclusive string `json:"maxExclusive"`
	// ResourceIDPrefix contains the resource ID prefix of the partition key range.
	ResourceIDPrefix int `json:"ridPrefix"`
	// SelfLink contains the self-link of the partition key range.
	SelfLink string `json:"_self"`
	// ThroughputFraction contains the throughput fraction of the partition key range.
	ThroughputFraction float64 `json:"throughputFraction"`
	// Status contains the status of the partition key range.
	Status string `json:"status"`
	// Parents contains the parent partition key ranges.
	Parents []string `json:"parents"`
	// OwnedArchivalPKRangeIds contains the owned archival partition key range IDs.
	OwnedArchivalPKRangeIds []string `json:"ownedArchivalPKRangeIds"`
	// LastModified contains the last modified time of the partition key range.
	LastModified time.Time `json:"_ts"`
	// LSN contains the LSN of the partition key range.
	LSN int64 `json:"lsn"`
}

// MarshalJSON implements the json.Marshaler interface
func (pkr PartitionKeyRange) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(fmt.Sprintf("\"id\":\"%s\"", pkr.ID))

	if pkr.ResourceID != "" {
		buffer.WriteString(fmt.Sprintf(",\"_rid\":\"%s\"", pkr.ResourceID))
	}

	if pkr.ETag != nil {
		buffer.WriteString(",\"_etag\":")
		etag, err := json.Marshal(pkr.ETag)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	if pkr.MinInclusive != "" {
		buffer.WriteString(fmt.Sprintf(",\"minInclusive\":\"%s\"", pkr.MinInclusive))
	}

	if pkr.MaxExclusive != "" {
		buffer.WriteString(fmt.Sprintf(",\"maxExclusive\":\"%s\"", pkr.MaxExclusive))
	}

	buffer.WriteString(fmt.Sprintf(",\"ridPrefix\":%d", pkr.ResourceIDPrefix))

	if pkr.SelfLink != "" {
		buffer.WriteString(fmt.Sprintf(",\"_self\":\"%s\"", pkr.SelfLink))
	}

	buffer.WriteString(fmt.Sprintf(",\"throughputFraction\":%f", pkr.ThroughputFraction))

	if pkr.Status != "" {
		buffer.WriteString(fmt.Sprintf(",\"status\":\"%s\"", pkr.Status))
	}

	if pkr.Parents != nil {
		parents, err := json.Marshal(pkr.Parents)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"parents\":")
		buffer.Write(parents)
	}

	if pkr.OwnedArchivalPKRangeIds != nil {
		ids, err := json.Marshal(pkr.OwnedArchivalPKRangeIds)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(",\"ownedArchivalPKRangeIds\":")
		buffer.Write(ids)
	}

	if !pkr.LastModified.IsZero() {
		buffer.WriteString(fmt.Sprintf(",\"_ts\":%v", strconv.FormatInt(pkr.LastModified.Unix(), 10)))
	}

	buffer.WriteString(fmt.Sprintf(",\"lsn\":%d", pkr.LSN))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (pkr *PartitionKeyRange) UnmarshalJSON(b []byte) error {
	var attributes map[string]json.RawMessage
	err := json.Unmarshal(b, &attributes)
	if err != nil {
		return err
	}

	if id, ok := attributes["id"]; ok {
		if err := json.Unmarshal(id, &pkr.ID); err != nil {
			return err
		}
	}

	if rid, ok := attributes["_rid"]; ok {
		if err := json.Unmarshal(rid, &pkr.ResourceID); err != nil {
			return err
		}
	}

	if etag, ok := attributes["_etag"]; ok {
		if err := json.Unmarshal(etag, &pkr.ETag); err != nil {
			return err
		}
	}

	if minInclusive, ok := attributes["minInclusive"]; ok {
		if err := json.Unmarshal(minInclusive, &pkr.MinInclusive); err != nil {
			return err
		}
	}

	if maxExclusive, ok := attributes["maxExclusive"]; ok {
		if err := json.Unmarshal(maxExclusive, &pkr.MaxExclusive); err != nil {
			return err
		}
	}

	if ridPrefix, ok := attributes["ridPrefix"]; ok {
		if err := json.Unmarshal(ridPrefix, &pkr.ResourceIDPrefix); err != nil {
			return err
		}
	}

	if self, ok := attributes["_self"]; ok {
		if err := json.Unmarshal(self, &pkr.SelfLink); err != nil {
			return err
		}
	}

	if throughputFraction, ok := attributes["throughputFraction"]; ok {
		if err := json.Unmarshal(throughputFraction, &pkr.ThroughputFraction); err != nil {
			return err
		}
	}

	if status, ok := attributes["status"]; ok {
		if err := json.Unmarshal(status, &pkr.Status); err != nil {
			return err
		}
	}

	if parents, ok := attributes["parents"]; ok {
		if err := json.Unmarshal(parents, &pkr.Parents); err != nil {
			return err
		}
	}

	if ids, ok := attributes["ownedArchivalPKRangeIds"]; ok {
		if err := json.Unmarshal(ids, &pkr.OwnedArchivalPKRangeIds); err != nil {
			return err
		}
	}

	if ts, ok := attributes["_ts"]; ok {
		var timestamp int64
		if err := json.Unmarshal(ts, &timestamp); err != nil {
			return err
		}
		pkr.LastModified = time.Unix(timestamp, 0)
	}

	if lsn, ok := attributes["lsn"]; ok {
		if err := json.Unmarshal(lsn, &pkr.LSN); err != nil {
			return err
		}
	}

	return nil
}
