# Azure Go SDK Breaking Changes Categories

## 1. Naming Changes

### Model/Type Naming

- **TypeSpec**: Model, enum, union name changes
- **Swagger**: Definition name changes in `#/definitions/`
- **Impact**: Type references become invalid

### Property/Parameter Naming

- **TypeSpec**: Property name changes in models, parameter name changes in operations
- **Swagger**: Property name changes in schemas, parameter name changes in paths
- **Impact**: Field access patterns break

### Operation Naming

- **TypeSpec**: Operation name changes
- **Swagger**: `operationId` changes
- **Impact**: Method names change in generated client

### Enum Value Naming

- **TypeSpec**: Enum member name changes
- **Swagger**: Enum value changes in schemas
- **Impact**: Constant values become unavailable

## 2. Operation Type Changes

### Long-Running Operation (LRO) Changes

- **TypeSpec**: Switch to `Async` related template
- **Swagger**: Addition of `x-ms-long-running-operation: true`
- **Impact**: Return type changes from direct result to poller

### Paging Operation Changes

- **TypeSpec**: Switch to `List` related templates
- **Swagger**: Addition of `x-ms-pageable` extension
- **Impact**: Return type changes from direct result to pager

## 3. Type Changes

### Property Type Changes

- **TypeSpec**: Property type reference changes
- **Swagger**: Property schema type changes
- **Impact**: Field types become incompatible

### Parameter Type Changes

- **TypeSpec**: Parameter type reference changes
- **Swagger**: Parameter schema type changes
- **Impact**: Method signatures change

### Response Type Changes

- **TypeSpec**: Response model changes
- **Swagger**: Response schema changes
- **Impact**: Return types become incompatible

## 4. Structural Changes

### Property/Parameter Deletion

- **TypeSpec**: Removal of model properties or operation parameters
- **Swagger**: Removal from schemas or parameter lists
- **Impact**: Fields/parameters no longer available

### Operation Deletion

- **TypeSpec**: Removal of operations from interfaces
- **Swagger**: Removal of paths or HTTP methods
- **Impact**: Client methods no longer available

### Model Deletion

- **TypeSpec**: Removal of model definitions
- **Swagger**: Removal from definitions section
- **Impact**: Types no longer available

## 5. Optionality Changes

### Required Parameter Addition

- **TypeSpec**: Addition of required parameters to operations
- **Swagger**: Addition of `required: true` parameters
- **Impact**: Method signatures require additional parameters

### Optional to Required Parameter Changes

- **TypeSpec**: Removal of `?` from parameter of operations
- **Swagger**: Change `required: true` for parameters
- **Impact**: Previously optional fields become mandatory

### Required to Optional Parameter Changes

- **TypeSpec**: Add `?` for parameter of operations
- **Swagger**: Remove `required: true` for parameters
- **Impact**: Previously mandatory fields become optional
