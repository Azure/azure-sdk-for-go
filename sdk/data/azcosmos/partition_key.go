// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// partitionKeyKind identifies which member of a [partitionKeyComponent] holds its value. The
// values mirror the driver's component kinds so the binding can pass them through directly.
type partitionKeyKind int32

const (
	partitionKeyKindString    partitionKeyKind = 0
	partitionKeyKindNumber    partitionKeyKind = 1
	partitionKeyKindBool      partitionKeyKind = 2
	partitionKeyKindNull      partitionKeyKind = 3
	partitionKeyKindUndefined partitionKeyKind = 4
)

// partitionKeyComponent is one component of a partition key value. Null and undefined components
// carry no value.
type partitionKeyComponent struct {
	kind        partitionKeyKind
	stringValue string
	numberValue float64
	boolValue   bool
}

// PartitionKey is a logical partition key value.
//
// A partition key has one component per path in the container's partition key definition, so a
// container with a hierarchical partition key takes a value with one component per level. Build
// those by appending:
//
//	pk := azcosmos.NewPartitionKeyString("Contoso").
//		AppendString("Redmond").
//		AppendNumber(98052)
//
// Cosmos DB distinguishes a component that is explicitly JSON null from one whose value is
// missing, and they route differently, so [PartitionKey.AppendNull] and
// [PartitionKey.AppendUndefined] are separate operations.
//
// The zero value is a partition key with no components, ready to be appended to. Values are
// immutable: every append returns a new key and leaves the receiver untouched, so a partially
// built key can be shared and extended in more than one direction.
type PartitionKey struct {
	components []partitionKeyComponent
}

// NewPartitionKey creates a partition key with no components, ready to be appended to. It is
// equivalent to the zero value.
func NewPartitionKey() PartitionKey {
	return PartitionKey{}
}

// NewPartitionKeyString creates a partition key with a single string component.
func NewPartitionKeyString(value string) PartitionKey {
	return NewPartitionKey().AppendString(value)
}

// NewPartitionKeyNumber creates a partition key with a single numeric component.
func NewPartitionKeyNumber(value float64) PartitionKey {
	return NewPartitionKey().AppendNumber(value)
}

// NewPartitionKeyBool creates a partition key with a single boolean component.
func NewPartitionKeyBool(value bool) PartitionKey {
	return NewPartitionKey().AppendBool(value)
}

// NewPartitionKeyNull creates a partition key with a single component that is explicitly JSON
// null. This is not the same as [NewPartitionKeyUndefined]; the two route differently.
func NewPartitionKeyNull() PartitionKey {
	return NewPartitionKey().AppendNull()
}

// NewPartitionKeyUndefined creates a partition key with a single component whose value is missing
// from the item. This is not the same as [NewPartitionKeyNull]; the two route differently.
func NewPartitionKeyUndefined() PartitionKey {
	return NewPartitionKey().AppendUndefined()
}

// AppendString returns a copy of the partition key with a string component appended.
func (pk PartitionKey) AppendString(value string) PartitionKey {
	return pk.append(partitionKeyComponent{kind: partitionKeyKindString, stringValue: value})
}

// AppendNumber returns a copy of the partition key with a numeric component appended.
func (pk PartitionKey) AppendNumber(value float64) PartitionKey {
	return pk.append(partitionKeyComponent{kind: partitionKeyKindNumber, numberValue: value})
}

// AppendBool returns a copy of the partition key with a boolean component appended.
func (pk PartitionKey) AppendBool(value bool) PartitionKey {
	return pk.append(partitionKeyComponent{kind: partitionKeyKindBool, boolValue: value})
}

// AppendNull returns a copy of the partition key with a component appended that is explicitly
// JSON null. This is not the same as [PartitionKey.AppendUndefined]; the two route differently.
func (pk PartitionKey) AppendNull() PartitionKey {
	return pk.append(partitionKeyComponent{kind: partitionKeyKindNull})
}

// AppendUndefined returns a copy of the partition key with a component appended whose value is
// missing from the item. This is not the same as [PartitionKey.AppendNull]; the two route
// differently.
func (pk PartitionKey) AppendUndefined() PartitionKey {
	return pk.append(partitionKeyComponent{kind: partitionKeyKindUndefined})
}

// Len reports how many components the partition key has. A value for a container with a
// hierarchical partition key has one component per level.
func (pk PartitionKey) Len() int {
	return len(pk.components)
}

// append copies the component slice rather than extending it in place. Appending in place would
// let two keys built from a common prefix share a backing array and overwrite each other's
// components, which is a trap for a type callers reasonably treat as a value.
func (pk PartitionKey) append(component partitionKeyComponent) PartitionKey {
	components := make([]partitionKeyComponent, len(pk.components), len(pk.components)+1)
	copy(components, pk.components)
	return PartitionKey{components: append(components, component)}
}
