// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
)

func partitionKeyPathCount(body []byte, operation string) (int, error) {
	var metadata struct {
		PartitionKey struct {
			Paths   []string `json:"paths"`
			Kind    string   `json:"kind"`
			Version *int     `json:"version"`
		} `json:"partitionKey"`
	}
	if err := json.Unmarshal(body, &metadata); err != nil {
		return 0, &Error{Code: CodeSerializationFailed, Message: "decoding container partition key definition", cause: err}
	}
	definition := metadata.PartitionKey
	if len(definition.Paths) == 0 {
		return 0, &Error{Code: CodeSerializationFailed, Message: "container has no partition key paths"}
	}
	if definition.Kind != "Hash" && definition.Kind != "MultiHash" {
		return 0, &Error{Code: CodeBadRequest, Message: operation + " requires a Hash or MultiHash partition key definition"}
	}
	if (definition.Kind == "Hash" && len(definition.Paths) != 1) ||
		(definition.Kind == "MultiHash" && (definition.Version == nil || *definition.Version != 2)) ||
		(definition.Version != nil && *definition.Version != 1 && *definition.Version != 2) {
		return 0, &Error{Code: CodeBadRequest, Message: "unsupported " + operation + " partition key definition"}
	}
	return len(definition.Paths), nil
}

func validateCompleteBatchPartitionKey(body []byte, partitionKey PartitionKey) error {
	count, err := partitionKeyPathCount(body, "transactional batch")
	if err != nil {
		return err
	}
	if partitionKey.Len() != count {
		return &Error{Code: CodeBadRequest, Message: fmt.Sprintf(
			"azcosmos: transactional batch requires exactly %d partition key components; got %d",
			count, partitionKey.Len())}
	}
	return nil
}
