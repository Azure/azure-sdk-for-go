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

// DatabaseProperties represents the properties of a database.
type DatabaseProperties struct {
	// ID contains the unique id of the database.
	ID string `json:"id"`
	// ETag contains the entity etag of the database
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink contains the self-link of the database
	SelfLink string `json:"_self,omitempty"`
	// ResourceID contains the resource id of the database
	ResourceID string `json:"_rid,omitempty"`
	// LastModified contains the last modified time of the database
	LastModified time.Time `json:"_ts,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface
func (tp DatabaseProperties) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	fmt.Fprintf(buffer, "\"id\":\"%s\"", tp.ID)

	if tp.ResourceID != "" {
		fmt.Fprintf(buffer, ",\"_rid\":\"%s\"", tp.ResourceID)
	}

	if tp.ETag != nil {
		fmt.Fprint(buffer, ",\"_etag\":")
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
		fmt.Fprintf(buffer, ",\"_ts\":%v", strconv.FormatInt(tp.LastModified.Unix(), 10))
	}

	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (tp *DatabaseProperties) UnmarshalJSON(b []byte) error {
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

	if ts, ok := attributes["_ts"]; ok {
		var timestamp int64
		if err := json.Unmarshal(ts, &timestamp); err != nil {
			return err
		}
		tp.LastModified = time.Unix(timestamp, 0)
	}

	return nil
}
