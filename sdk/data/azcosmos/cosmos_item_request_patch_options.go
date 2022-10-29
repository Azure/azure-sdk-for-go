package azcosmos

// PatchOperationType defines supported values for operation types in Patch Document.
type PatchOperationType string

const (
	// Represents a patch operation Add.
	PatchOperationTypeAdd PatchOperationType = "add"
	// Represents a patch operation Replace.
	PatchOperationTypeReplace PatchOperationType = "replace"
	// Represents a patch operation Remove.
	PatchOperationTypeRemove PatchOperationType = "remove"
	// Represents a patch operation Set.
	PatchOperationTypeSet PatchOperationType = "set"
	// Represents a patch operation Increment.
	PatchOperationTypeIncrement PatchOperationType = "incr"
)

// PatchOperation represents individual patch operation.
type PatchOperation struct {
	Op    PatchOperationType `json:"op"`
	Path  string             `json:"path"`
	Value interface{}        `json:"value,omitempty"`
}

// PatchOptions represents the patch request.
type PatchOptions struct {
	Condition  *string          `json:"condition,omitempty"`
	Operations []PatchOperation `json:"operations"`
}

// CreatePatchOptions returns a struct that represents the patch request.
func CreatePatchOptions() *PatchOptions {
	return &PatchOptions{
		Condition:  nil,
		Operations: make([]PatchOperation, 0),
	}
}

// SetCondition sets condition for the patch request.
func (p *PatchOptions) SetCondition(condition string) {
	p.Condition = &condition
}

// Replace appends a replace operation to the patch request.
func (p *PatchOptions) Replace(path string, value interface{}) {
	p.Operations = append(p.Operations, PatchOperation{
		Op:    PatchOperationTypeReplace,
		Path:  path,
		Value: value,
	})
}

// Replace appends an add operation to the patch request.
func (p *PatchOptions) Add(path string, value interface{}) {
	p.Operations = append(p.Operations, PatchOperation{
		Op:    PatchOperationTypeAdd,
		Path:  path,
		Value: value,
	})
}

// Replace appends a set operation to the patch request.
func (p *PatchOptions) Set(path string, value interface{}) {
	p.Operations = append(p.Operations, PatchOperation{
		Op:    PatchOperationTypeSet,
		Path:  path,
		Value: value,
	})
}

// Replace appends a remove operation to the patch request.
func (p *PatchOptions) Remove(path string) {
	p.Operations = append(p.Operations, PatchOperation{
		Op:   PatchOperationTypeRemove,
		Path: path,
	})
}

// Replace appends an increment operation to the patch request.
func (p *PatchOptions) Increment(path string, value interface{}) {
	p.Operations = append(p.Operations, PatchOperation{
		Op:    PatchOperationTypeIncrement,
		Path:  path,
		Value: value,
	})
}
