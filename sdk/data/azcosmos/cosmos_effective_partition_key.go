// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// EffectivePartitionKey defines the contract for computing effective partition key strings.
type CosmosEffectivePartitionKey interface {
	// ComputeEffectivePartitionKey takes a JSON string representing either a single value (e.g. "\"abc\"")
	// or an array of values (e.g. "[\"abc\", 5]") and returns the hex EPK string.
	// Rules :
	//   * Empty JSON string => error
	//   * Top-level "Infinity" string stands for the boundary value (no array brackets)
	//   * Arrays cannot be nested; non-empty objects are invalid; only the empty object represents Undefined
	//   * MultiHash requires version V2
	ComputeEffectivePartitionKey(partitionKeyJSON string, version int, kind PartitionKeyKind) (string, error)
}
