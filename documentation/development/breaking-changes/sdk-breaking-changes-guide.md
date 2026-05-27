# Azure Go SDK Breaking Changes Review and Resolution Guide

The Azure Go SDK generally prohibits breaking changes unless they result from service behavior modifications. This guide helps you identify, review, and resolve breaking changes that may occur in new SDK versions due to service's TypeSpec specification update. For migration of service specifications from Swagger to TypeSpec, refer to this [doc](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/breaking-changes/sdk-breaking-changes-guide-migration.md).

Some breaking changes can be resolved through client customizations. You should follow the guidelines below to review and resolve breaking changes.

Client customizations should be implemented in a file named `client.tsp` located in the service's specification directory alongside the main entry point `main.tsp`. This `client.tsp` becomes the new specification entry point, so import `main.tsp` in the `client.tsp` file. **Do not** import `client.tsp` in the `main.tsp` file. **Do not** modify the entry point in `tspconfig.yaml`.

```tsp
import "./main.tsp";
import "@azure-tools/typespec-client-generator-core";

using Azure.ClientGenerator.Core;
using MainNamespaceInMainTsp; // Replace with the actual main namespace in main.tsp

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

**Breaking**: An exported type `A` is renamed to `B`, invalidating any references to `A`.

**Reason**: The type was renamed in the TypeSpec specification using `@renamedFrom`, so the generated Go type name changes accordingly.

**Resolution**:

Use client customization to restore the original type name:

```tsp
@@clientName(B, "A", "go");
```

## 2. Property Name Changes

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

**Breaking**: The exported field `A` on struct `Test` is renamed to `B`, breaking field access patterns.

**Reason**: The property was renamed in the TypeSpec specification using `@renamedFrom`, so the generated Go field name changes accordingly.

**Resolution**:

Use client customization to restore the original property name:

```tsp
@@clientName(Test.b, "a", "go");
```

## 3. Operation Name Changes

**Changelog Pattern**:

```md
- Function `*xxx.A` has been removed
- New function `*xxx.B(xxx) *xxx`
```

**Spec Pattern**:

```tsp
@renamedFrom(Versions.v2, "a")
op b(): void;
```

**Breaking**: The client method `A` is renamed to `B`, breaking any callers of the original method.

**Reason**: The operation was renamed in the TypeSpec specification using `@renamedFrom`, so the generated client method name changes accordingly.

**Resolution**:

Use client customization to restore the original operation name:

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

**Breaking**: The exported enum value `A` is renamed to `B`, breaking references to the original constant.

**Reason**: The enum value was renamed in the TypeSpec specification using `@renamedFrom`, so the generated Go constant name changes accordingly.

**Resolution**:

Use client customization to restore the original enum value name:

```tsp
@@clientName(Test.b, "a", "go");
```

## 5. Long-Running Operation (LRO) Changes

**Changelog Pattern**:

```md
- Operation `*xxx.A` has been changed to LRO, use `*xxx.BeginA` instead.
- Operation `*xxx.BeginB` has been changed to non-LRO, use `*xxx.B` instead.
```

**Breaking**: An operation switches between synchronous and long-running form, which changes both the method name (`A` ↔ `BeginA`) and the return type (direct result ↔ poller).

**Reason**: The service's TypeSpec operation template (e.g., switching to/from an `*Async` ARM template) was changed to mark the operation as long-running or non-long-running.

**Resolution**: Cannot be resolved through client customizations.

## 6. Paging Operation Changes

**Changelog Pattern**:

```md
- Function `*xxx.NewListAPager` has been removed
- New function `*xxx.A(xxx) (xxx, error)`
```

**Breaking**: An operation switches between paged and non-paged form, which changes both the method name (`NewListAPager` ↔ `A`) and the return type (pager ↔ direct result).

**Reason**: The service's TypeSpec operation template was changed to mark the operation as paged or non-paged (e.g., adding or removing `@list` / a paged template).

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

**Breaking**: The type of an exported field `Test.Prop` changes from `*string` to `*int32`, requiring callers to update their code with type conversions or different value handling.

**Reason**: The property's type was modified in the TypeSpec specification using `@typeChangedFrom`.

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

**Breaking**: The type of a parameter on `*xxx.Test` changes from `string` to `int32`, breaking the method signature for callers.

**Reason**: The operation parameter's type was modified in the TypeSpec specification using `@typeChangedFrom`.

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

**Breaking**: The type of a response field changes from `*string` to `*int32`, breaking callers that read the response.

**Reason**: The operation's return type was modified in the TypeSpec specification using `@returnTypeChangedFrom`.

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

**Breaking**: The exported field `DeletedProp` on struct `Test` is no longer available, causing compilation errors for any code that reads or writes it.

**Reason**: The property was removed in the TypeSpec specification using `@removed`.

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

**Breaking**: A parameter is removed from `*xxx.Test`, breaking the method signature for callers.

**Reason**: The operation parameter was removed in the TypeSpec specification using `@removed`.

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

**Breaking**: The client method `*xxx.Test` is no longer available, breaking any caller that invoked it.

**Reason**: The operation was removed in the TypeSpec specification using `@removed`.

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

**Breaking**: The exported type `Test` is no longer available, breaking any code that references the type.

**Reason**: The model was removed in the TypeSpec specification using `@removed`.

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

**Breaking**: A new required parameter is added to `*xxx.Test`, breaking existing call sites which must now supply the additional argument.

**Reason**: A new required parameter was added to the operation in the TypeSpec specification using `@added`.

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

**Breaking**: A previously optional parameter on `*xxx.Test` becomes required: the parameter moves out of the options struct and into the method signature, breaking existing callers.

**Reason**: The parameter was changed from optional to required in the TypeSpec specification using `@madeRequired`.

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

**Breaking**: A previously required parameter on `*xxx.Test` becomes optional: the parameter moves out of the method signature and into the options struct, breaking existing callers.

**Reason**: The parameter was changed from required to optional in the TypeSpec specification using `@madeOptional`.

**Resolution**: Cannot be resolved through client customizations.
