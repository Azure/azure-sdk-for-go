package azcosmos

type PartitionKeyBuilder struct {
	components []interface{}
}

// NewPartitionKeyBuilder creates a new PartitionKeyBuilder.
func NewPartitionKeyBuilder() *PartitionKeyBuilder {
	return &PartitionKeyBuilder{}
}

// AppendString appends a string value to the partition key.
func (pkb *PartitionKeyBuilder) AppendString(value string) *PartitionKeyBuilder {
	pkb.components = append(pkb.components, value)
	return pkb
}

// AppendBool appends a boolean value to the partition key.
func (pkb *PartitionKeyBuilder) AppendBool(value bool) *PartitionKeyBuilder {
	pkb.components = append(pkb.components, value)
	return pkb
}

// AppendNumber appends a numeric value to the partition key.
func (pkb *PartitionKeyBuilder) AppendNumber(value float64) *PartitionKeyBuilder {
	pkb.components = append(pkb.components, value)
	return pkb
}

func (pkb *PartitionKeyBuilder) AppendNull() *PartitionKeyBuilder {
	pkb.components = append(pkb.components, nil)
	return pkb
}

// Build creates a PartitionKey from the appended values.
func (pkb *PartitionKeyBuilder) Build() PartitionKey {
	return PartitionKey{
		values: pkb.components,
	}
}
