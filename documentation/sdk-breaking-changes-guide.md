# Azure Go SDK Breaking Changes Review and Resolution Guide

This document categorizes common breaking changes in the Azure Go SDK and provides guidance on how to resolve them using client customizations in the spec or by changing `tspconfig.yaml` configuration.
Customizations should be implemented in a file named `client.tsp` under the service's spec directory, which contains the main entry point `main.tsp` of the service.

```tsp
import "./main.tsp";
import "@azure-tools/typespec-client-generator-core";

using Azure.ClientGenerator.Core;

// Add your customizations here
```

## Migration to TypeSpec

These breaking changes occur when migrating from Swagger to TypeSpec for API specifications.

### 1. Naming Changes with Numbers

**Changelog Pattern**:

You can find removals of enum values and additions of new enum values with naming changes from words to numbers.

```md
- `MinuteThirty`, `MinuteZero` from enum `Minute` has been removed
- New value `Minute0`, `Minute30` added to enum type `Minute`
```

**Reason**: Swagger translates names with numbers to words during code generation, but TypeSpec doesn't apply this transformation.

**Spec Pattern**:

You can find the enum and enum values in the spec by examining the names from addition items in the changelog (pattern: `New value '<enum value>' added to enum type '<enum name>'`):

```tsp
enum Minute {
  Minute0 = "0",
  Minute30 = "30"
}
```

**Resolution**:

Rename the enum values to match the names from removal items in the changelog.

```tsp
@@clientName(Minute.Minute0, "Zero", "go");
@@clientName(Minute.Minute30, "Thirty", "go");
```

### 2. Naming Changes of Enum from Stuttering Rules

**Changelog Pattern**:

You can find removals of an enum type and additions of a new enum type with the service name prefix removed. All references to the old enum type are also changed.

```md
- Type of `ConfigurationProperties.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Type of `Update.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Enum `MaintenanceScope` has been removed
- New enum type `Scope` with values `ScopeExtension`, `ScopeHost`, `ScopeInGuestPatch`, `ScopeOSImage`, `ScopeResource`, `ScopeSQLDB`, `ScopeSQLManagedInstance`
```

**Reason**: Discrepancies between stuttering rules for enums in Swagger and TypeSpec can cause enum name changes.

**Spec Pattern**:

You can find the enum by examining the name from removal items in the changelog (pattern: `Enum '<enum name>' has been removed`):

```tsp
union MaintenanceScope {
  string,
  Host: "Host",
  Resource: "Resource",
}
```

**Resolution**:

Change the `fix-const-stuttering` configuration in `tspconfig.yaml` for Go from `true` to `false` to prevent stuttering in enum names:

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

### 3. Naming Changes of Operations

**Changelog Pattern**:

You can see a removal of an operation and addition of a new operation with a similar name for the same client.

```md
- Function `*StorageTaskAssignmentClient.NewListPager` has been removed
- New function `*StorageTaskAssignmentClient.NewStorageTaskAssignmentListPager(string, string, *StorageTaskAssignmentClientStorageTaskAssignmentListOptions) *runtime.Pager[StorageTaskAssignmentClientStorageTaskAssignmentListResponse]`
```

**Reason**: To prevent naming collisions, TypeSpec may generate different operation names than Swagger.

**Spec Pattern**:

You can find the interface and operation in the spec by examining the name from addition items in the changelog (pattern: `New function *<interface name>Client.<operation name>(...)`, `New function *<interface name>Client.New<operation name>Pager(...)` for paging operations, and `New function *<interface name>Client.Begin<operation name>(...)` for long-running operations):

```tsp
interface StorageTaskAssignment {
  op storageTaskAssignmentList(xxx): xxx;
}
```

**Resolution**:

Rename the operation to match the name from removal items in the changelog.

**Note**: For paging operations, the method name in the SDK is `New<OperationName>Pager`. For long-running operations, the method name is `Begin<OperationName>`. When resolving the breaking change, use the name in TypeSpec without the SDK-specific prefix or suffix.

```tsp
@@clientName(StorageTaskAssignment.storageTaskAssignmentList, "list", "go");
```

### 4. Client Organization Changes

**Changelog Pattern**:

You can see operations moving between clients, sometimes along with client removal.

```md
- Function `NewManagementClient` has been removed
- Function `*ManagementClient.BeginRestoreVolume` has been removed
- New function `*VolumesClient.BeginRestoreVolume(context.Context, string, string, string, string, *VolumesClientBeginRestoreVolumeOptions) (*runtime.Poller[VolumesClientRestoreVolumeResponse], error)`
```

**Reason**: Different logic for organizing client operations in TypeSpec.

**Spec Pattern**:

You can find the interface and operation in the spec by examining the name from addition items in the changelog (pattern: `New function *<interface name>Client.<operation name>(...)`):

```tsp
namespace Microsoft.ElasticSan;

interface Volumes {
  @action("restore")
  op restoreVolume is ArmResourceActionAsync<Volume, void, Volume>;
}
```

**Resolution**:

Move the operation to the correct client using `@clientLocation`. The new client name should follow the removal items in the changelog (remove the `Client` suffix).

```tsp
@@clientLocation(Microsoft.ElasticSan.restoreVolume, "Management", "go");
```

### 5. Missing Fields in Response Types

**Changelog Pattern**:

You can see the removal of fields in response structs named with the `xxxResponse` pattern.

```md
- Field `CacheAccessPolicyAssignment` of struct `AccessPolicyAssignmentClientCreateUpdateResponse` has been removed
```

**Reason**: Incorrect TypeSpec conversion for LRO operations.

**Spec Pattern**:

You can find the interface and operation in the spec by examining the name from removal items in the changelog (pattern: `Field 'xxx' of struct *<interface name>Client<operation name>Response`):

```tsp
@armResourceOperations
interface RedisCacheAccessPolicies {
  createUpdate is ArmResourceCreateOrReplaceAsync<
    RedisCacheAccessPolicy,
    LroHeaders = ArmLroLocationHeader & Azure.Core.Foundations.RetryAfterHeader
  >;
}
```

**Resolution**:

Locate the operation and add the `FinalResult` parameter in `ArmLroLocationHeader`, `ArmAsyncOperationHeader`, or `ArmCombinedLroHeaders` with the correct type:

```tsp
@armResourceOperations
interface RedisCacheAccessPolicies {
  createUpdate is ArmResourceCreateOrReplaceAsync<
    RedisCacheAccessPolicy,
    LroHeaders = ArmLroLocationHeader<FinalResult = CacheAccessPolicyAssignment> & Azure.Core.Foundations.RetryAfterHeader
  >;
}
```

### 6. Use Standard Operations List Operation

**Changelog Pattern**:

You can see a batch of changes related to the `Operation` type and its related fields, sometimes along with changes to the `OperationList` operation.

```md
- Type of `Operation.Display` has been changed from `*OperationInfo` to `*OperationDisplay`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Struct `OperationInfo` has been removed
- Field `Properties` of struct `Operation` has been removed
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New struct `OperationDisplay`
- New field `ActionType` in struct `Operation`
```

**Reason**: The operations list operation is upgraded to use the standard library's definition.

**Impact**: Low impact for users as this operation is rarely used in the SDK.

**Resolution**: Accept these breaking changes.

### 7. Common Types Upgrade

**Changelog Pattern**:

You can see a batch of changes related to common types, such as `SystemData`, `Error`, and `IdentityType`.

```md
- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Type of `Error.Error` has been changed from `*ErrorError` to `*ErrorDetail`
- Type of `SystemData.CreatedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Enum `IdentityType` has been removed
- Struct `ErrorError` has been removed
```

**Reason**: Common types are upgraded to the latest versions during TypeSpec migration.

**Impact**: Low impact for users as these are common infrastructure types.

**Resolution**: Accept these breaking changes.

### 8. Removal of Unreferenced Types

**Changelog Pattern**:

You can see a batch of removals of unreferenced types, which are typically not used in the SDK.

```md
- Struct `TrackedResource` has been removed
- Struct `Resource` has been removed
- Struct `ProxyResource` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `SCConfluentListMetadata` has been removed
```

**Reason**: Unreferenced types are removed during TypeSpec migration.

**Impact**: Low impact for users as these types are typically not directly used.

**Resolution**: Accept these breaking changes.

## TypeSpec Update

These breaking changes result from API specification updates in TypeSpec files.

### 1. Model/Enum/Union Name Changes

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

### 3. Operation Name Changes

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

### 4. Enum Value Name Changes

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

### 5. Long-Running Operation (LRO) Changes

**Changelog Pattern**:

```md
- Operation `*xxx.A` has been changed to LRO, use `*xxx.BeginA` instead.
- Operation `*xxx.BeginB` has been changed to non-LRO, use `*xxx.B` instead.
```

**Impact**: Method name changes and return type changes (direct result ↔ poller).

**Resolution**: Cannot be resolved.

### 6. Paging Operation Changes

**Changelog Pattern**:

```md
- Function `*xxx.NewListAPager` has been removed
- New function `*xxx.A(xxx) (xxx, error)`
```

**Impact**: Method name changes and return type changes (direct result ↔ pager).

**Resolution**: Cannot be resolved.

### 7. Property Type Changes

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

**Resolution**: Cannot be resolved.

### 8. Parameter Type Changes

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

**Resolution**: Cannot be resolved.

### 9. Response Type Changes

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

**Resolution**: Cannot be resolved.

### 10. Property Deletion

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

**Impact**: Fields no longer available, causing compilation errors.

**Resolution**: Cannot be resolved.

### 11. Parameter Deletion

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

**Impact**: Parameters no longer available, requiring method signature updates.

**Resolution**: Cannot be resolved.

### 12. Operation Deletion

**Changelog Pattern**:

```md
- Function `*xxx.Test` has been removed
```

**Spec Pattern**:

```tsp
@removed(Versions.v2)
op test(): void;
```

**Impact**: Client methods no longer available, requiring alternative implementation.

**Resolution**: Cannot be resolved.

### 13. Model Deletion

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

**Impact**: Types no longer available, requiring alternative type usage.

**Resolution**: Cannot be resolved.

### 14. Required Parameter Addition

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

**Resolution**: Cannot be resolved.

### 15. Optional to Required Parameter Changes

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

**Impact**: Previously optional fields become mandatory, requiring parameter updates.

**Resolution**: Cannot be resolved.

### 16. Required to Optional Parameter Changes

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

**Impact**: Previously mandatory fields become optional, potentially affecting validation logic.

**Resolution**: Cannot be resolved.
