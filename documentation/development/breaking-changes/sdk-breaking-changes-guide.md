# Azure Go SDK Breaking Changes Review and Resolution Guide

The Azure Go SDK generally prohibits breaking changes unless they result from service behavior modifications. This guide helps you identify, review, and resolve breaking changes that may occur in new SDK versions due to service's TypeSpec specification update. For migration of service specifications from Swagger to TypeSpec, refer this [doc](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/breaking-changes/sdk-breaking-changes-guide-migration.md).

Breaking changes can be resolved by client Customizations:

Client customizations should be implemented in a file named `client.tsp` located in the service's specification directory alongside the main entry point `main.tsp`. This `client.tsp` becomes the new specification entry point, so import `main.tsp` in the `client.tsp` file. **Do not** import `client.tsp` in the `main.tsp` file.

```tsp
import "./main.tsp";
import "@azure-tools/typespec-client-generator-core";

using Azure.ClientGenerator.Core;

// Add your customizations here
```

## 1. Model/Enum/Union Name Changes

**Changelog Pattern**:

```md
- Struct `A` has been removed
- New struct `B`
```

**Spec Pattern**:

```tsp
@renamedFrom(Versions.v2, "A")
model B {
  prop: string
}
```

**Impact**: Type references become invalid.

**Resolution**:

```tsp
@@clientName(B, "A", "go");
```

### 2. Property Name Changes

**Changelog Pattern**:

```md
- Field `A` of struct `Test` has been removed
- New field `B` in struct `Test`
```

**Spec Pattern**:

```tsp
model Test {
  @renamedFrom(Versions.v2, "a")
  b: string
}
```

**Impact**: Field access patterns break.

**Resolution**:

```tsp
@@clientName(Test.b, "a", "go");
```

## 3. Operation Name Changes

**Changelog Pattern**:

```md
- Function 'A' has been removed
- New function '*xxx.B(xxx) *xxx'
```

**Spec Pattern**:

```tsp
@renamedFrom(Versions.v2, "a")
op b(): void;
```

**Impact**: Method names change in the generated client.

**Resolution**:

```tsp
@@clientName(b, "a", "go");
```

## 4. Enum Value Name Changes

**Changelog Pattern**:

```md
- `A` from enum `Test` has been removed
- New value `B` added to enum type `Test`
```

**Spec Pattern**:

```tsp
enum Test {
  @renamedFrom(Versions.v1, "a")
  b: "b",
}
```

**Impact**: Constant values become unavailable.

**Resolution**:

```tsp
@@clientName(Test.b, "a", "go");
```

## 5. Long-Running Operation (LRO) Changes

**Changelog Pattern**:

```md
- Operation `*xxx.A` has been changed to LRO, use `*xxx.BeginA` instead.
- Operation `*xxx.BeginB` has been changed to non-LRO, use `*xxx.B` instead.
```

**Impact**: Method names and return types change (direct result ↔ poller).

**Resolution**: Cannot be resolved through client customizations.

## 6. Paging Operation Changes

**Changelog Pattern**:

```md
- Function `*xxx.NewListAPager` has been removed
- New function `*xxx.A(xxx) (xxx, error)`
```

**Impact**: Method names and return types change (direct result ↔ pager).

**Resolution**: Cannot be resolved through client customizations.

## 7. Property Type Changes

**Changelog Pattern**:

```md
- Type of `Test.Prop` has been changed from `*string` to `*int32`
```

**Spec Pattern**:

```tsp
model Test {
  @typeChangedFrom(Versions.v2, "string")
  prop: int32
}
```

**Impact**: Field types become incompatible, requiring type casting or conversion.

**Resolution**: Cannot be resolved through client customizations.

## 8. Parameter Type Changes

**Changelog Pattern**:

```md
- Function `*xxx.Test` parameter(s) have been changed from `(string)` to `(int32)`
```

**Spec Pattern**:

```tsp
op test(
  @typeChangedFrom(Versions.v2, "string")
  prop: int32
): void;
```

**Impact**: Method signatures change, requiring parameter type updates.

**Resolution**: Cannot be resolved through client customizations.

## 9. Response Type Changes

**Changelog Pattern**:

```md
- Type of `xxxTestResponse.Result` has been changed from `*string` to `*int32`
```

**Spec Pattern**:

```tsp
op test(): {
  @returnTypeChangedFrom(Versions.v2, "string")
  @body result: int32
};
```

**Impact**: Return types become incompatible, requiring response handling updates.

**Resolution**: Cannot be resolved through client customizations.

## 10. Property Deletion

**Changelog Pattern**:

```md
- Field `DeletedProp` of struct `Test` has been removed
```

**Spec Pattern**:

```tsp
model Test {
  @removed(Versions.v2)
  deletedProp: string;

  remainingProp: string;
}
```

**Impact**: Fields are no longer available, causing compilation errors.

**Resolution**: Cannot be resolved through client customizations.

## 11. Parameter Deletion

**Changelog Pattern**:

```md
- Function `*xxx.Test` parameter(s) have been changed from `(string)` to `()`
```

**Spec Pattern**:

```tsp
op test(
  @removed(Versions.v2)
  deletedParam: string
): void;
```

**Impact**: Parameters are no longer available, requiring method signature updates.

**Resolution**: Cannot be resolved through client customizations.

## 12. Operation Deletion

**Changelog Pattern**:

```md
- Function `*xxx.Test` has been removed
```

**Spec Pattern**:

```tsp
@removed(Versions.v2)
op test(): void;
```

**Impact**: Client methods are no longer available, requiring alternative implementation.

**Resolution**: Cannot be resolved through client customizations.

## 13. Model Deletion

**Changelog Pattern**:

```md
- Struct `Test` has been removed
```

**Spec Pattern**:

```tsp
@removed(Versions.v2)
model Test {
  prop: string
}
```

**Impact**: Types are no longer available, requiring alternative type usage.

**Resolution**: Cannot be resolved through client customizations.

## 14. Required Parameter Addition

**Changelog Pattern**:

```md
- Function `*xxx.Test` parameter(s) have been changed from `(string)` to `(string, int32)`
```

**Spec Pattern**:

```tsp
op test(
  existingParam: string,
  @added(Versions.v2)
  newParam: string
): void;
```

**Impact**: Method signatures require additional parameters, breaking existing calls.

**Resolution**: Cannot be resolved through client customizations.

## 15. Optional to Required Parameter Changes

**Changelog Pattern**:

```md
- Function `*xxx.Test` parameter(s) have been changed from `()` to `(string)`
- Field `Param` of struct `xxxTestOptions` has been removed
```

**Spec Pattern**:

```tsp
op test(
  @madeRequired(Versions.v2)
  param: string
): void;
```

**Impact**: Previously optional parameters become mandatory, requiring parameter updates.

**Resolution**: Cannot be resolved through client customizations.

## 16. Required to Optional Parameter Changes

**Changelog Pattern**:

```md
- Function `*xxx.Test` parameter(s) have been changed from `(string)` to `()`
- New field `Param` in struct `xxxTestOptions`
```

**Spec Pattern**:

```tsp
op test(
  @madeOptional(Versions.v2)
  param?: string
): void;
```

**Impact**: Previously mandatory parameters become optional, potentially affecting validation logic.

**Resolution**: Cannot be resolved through client customizations.
