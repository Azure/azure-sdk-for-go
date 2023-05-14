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

type (
	// PartitionKeyRangesProperties represents the response from a partition key ranges request.
	PartitionKeyRangesProperties struct {
		// ResourceID contains the resource id of the partition key ranges
		ResourceID string
		// Count contains the number of partition key ranges
		Count int
		// PartitionKeyRanges contains the list of partition key ranges
		PartitionKeyRanges []PartitionKeyRange
	}

	// PartitionKeyRange represents a single partition key range.
	PartitionKeyRange struct {
		// ETag contains the entity etag of the partition key range.
		ETag *azcore.ETag
		// LastModified contains the last modified time of the partition key range information.
		LastModified time.Time
		// SelfLink contains the self-link of the partition key range.
		SelfLink string
		// ResourceID contains the resource id of the partition key range.
		ResourceID string
		// ID contains the partition key range ID.
		ID string
		// MinInclusive contains the minimum value of the partition key range.
		MinInclusive string
		// MaxExclusive contains the maximum value of the partition key range.
		MaxExclusive string
	}
)

// MarshalJSON implements the json.Marshaler interface
func (p *PartitionKeyRangesProperties) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	if p.ResourceID != "" {
		buffer.WriteString(fmt.Sprintf("\"_rid\":\"%s\"", p.ResourceID))
	}

	buffer.WriteString(fmt.Sprintf(",\"_count\":%d", p.Count))

	if p.Count == 0 {
		buffer.WriteString(",\"PartitionKeyRanges\":[]}")
		return buffer.Bytes(), nil
	}

	buffer.WriteString(",\"PartitionKeyRanges\":[")
	for _, r := range p.PartitionKeyRanges {
		buffer.WriteString("{")

		if r.ETag != nil {
			buffer.WriteString("\"_etag\":")
			etag, err := json.Marshal(r.ETag)
			if err != nil {
				return nil, err
			}
			buffer.Write(etag)
		}

		if r.SelfLink != "" {
			buffer.WriteString(fmt.Sprintf(",\"_self\":\"%s\"", r.SelfLink))
		}

		if !r.LastModified.IsZero() {
			buffer.WriteString(fmt.Sprintf(",\"_ts\":%v", strconv.FormatInt(r.LastModified.Unix(), 10)))
		}

		if r.ID != "" {
			buffer.WriteString(fmt.Sprintf(",\"id\":\"%s\"", r.ID))
		}

		buffer.WriteString(fmt.Sprintf(",\"minInclusive\":\"%s\"", r.MinInclusive))
		buffer.WriteString(fmt.Sprintf(",\"maxExclusive\":\"%s\"", r.MaxExclusive))

		buffer.WriteString("}")
	}
	buffer.WriteString("]}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (p *PartitionKeyRangesProperties) UnmarshalJSON(b []byte) error {
	var attributes map[string]json.RawMessage
	err := json.Unmarshal(b, &attributes)
	if err != nil {
		return err
	}

	if id, ok := attributes["_rid"]; ok {
		if err := json.Unmarshal(id, &p.ResourceID); err != nil {
			return err
		}
	}

	if count, ok := attributes["_count"]; ok {
		if err := json.Unmarshal(count, &p.Count); err != nil {
			return err
		}
	}

	if p.Count == 0 {
		return nil
	}

	if _, ok := attributes["PartitionKeyRanges"]; !ok {
		return nil
	}

	var pkRanges []json.RawMessage
	if err := json.Unmarshal(attributes["PartitionKeyRanges"], &pkRanges); err != nil {
		return err
	}

	for _, r := range pkRanges {
		var pkAttributes map[string]json.RawMessage
		if err := json.Unmarshal(r, &pkAttributes); err != nil {
			return err
		}

		var pkr PartitionKeyRange

		if etag, ok := pkAttributes["_etag"]; ok {
			if err := json.Unmarshal(etag, &pkr.ETag); err != nil {
				return err
			}
		}

		if ts, ok := pkAttributes["_ts"]; ok {
			var timestamp int64
			if err := json.Unmarshal(ts, &timestamp); err != nil {
				return err
			}
			pkr.LastModified = time.Unix(timestamp, 0)
		}

		if id, ok := pkAttributes["id"]; ok {
			if err := json.Unmarshal(id, &pkr.ID); err != nil {
				return err
			}
		}

		if self, ok := pkAttributes["_self"]; ok {
			if err := json.Unmarshal(self, &pkr.SelfLink); err != nil {
				return err
			}
		}

		if min, ok := pkAttributes["minInclusive"]; ok {
			if err := json.Unmarshal(min, &pkr.MinInclusive); err != nil {
				return err
			}
		}

		if max, ok := pkAttributes["maxExclusive"]; ok {
			if err := json.Unmarshal(max, &pkr.MaxExclusive); err != nil {
				return err
			}
		}

		p.PartitionKeyRanges = append(p.PartitionKeyRanges, pkr)
	}

	return nil
}
