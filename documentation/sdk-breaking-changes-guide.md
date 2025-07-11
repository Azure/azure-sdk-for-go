# Azure Go SDK Breaking Changes Review and Resolution Guide

This document categorizes common breaking changes in the Azure Go SDK and provides guidance on how to resolve them using `client.tsp` customizations or change `tspconfig.yaml` config.

## Table of Contents

- [TypeSpec Migration](#typespec-migration)
- [TypeSpec Update](#typespec-update)
- [Resolution Guidelines](#resolution-guidelines)

## TypeSpec Migration

These breaking changes occur when migrating from Swagger to TypeSpec for API specifications. Most can be resolved using `client.tsp` customizations.

### Naming Changes with Numbers

**Issue**: Swagger translates names with numbers to words during codegen, but TypeSpec doesn't apply this transformation.

**Changelog Pattern**:

```md
- `MinuteThirty`, `MinuteZero` from enum `Minute` has been removed
- New value `Minute0`, `Minute30` added to enum type `Minute`
```

**Spec Pattern**:

```tsp
enum Minute {
  Minute0 = "0",
  Minute30 = "30"
}
```

**Resolution**:

```tsp
@@clientName(Minute.Minute0, "Zero", "go");
@@clientName(Minute.Minute30, "Thirty", "go");
```

### Naming Changes of Enum from Stuttering Rules

**Issue**: Discrepancies between stuttering rules for enum in Swagger and TypeSpec can cause method name changes.

**Changelog Pattern**:

```md
- Type of `ConfigurationProperties.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Type of `Update.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Enum `MaintenanceScope` has been removed
- New enum type `Scope` with values `ScopeExtension`, `ScopeHost`, `ScopeInGuestPatch`, `ScopeOSImage`, `ScopeResource`, `ScopeSQLDB`, `ScopeSQLManagedInstance`
```

**Spec Pattern**:

```tsp
union MaintenanceScope {
  string,
  Host: "Host",
  Resource: "Resource",
}
```

**Resolution**:

Change the `fix-const-stuttering` config in `tspconfig.yaml` for Go from `true` to `false` to prevent stuttering in enum names:

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

### Naming Changes of Operation

**Issue**: In order to prevent collision, TypeSpec could have different operation name from swagger.

**Note**: For paging operations, the method name in SDK is `New<OperationName>Pager`. For long running operations, the method name is `Begin<OperationName>`. When you resolve the breaking change, you should use the name patten in TypeSpec.

**Changelog Pattern**:

```md
- Function `*StorageTaskAssignmentClient.NewListPager` has been removed
- New function `*StorageTaskAssignmentClient.NewStorageTaskAssignmentListPager(string, string, *StorageTaskAssignmentClientStorageTaskAssignmentListOptions) *runtime.Pager[StorageTaskAssignmentClientStorageTaskAssignmentListResponse]`
```

**Spec Pattern**:

```tsp
interface StorageTaskAssignment {
  op storageTaskAssignmentList(xxx): xxx;
}
```

**Resolution**:

```tsp
@@clientName(StorageTaskAssignment.storageTaskAssignmentList, "list", "go");
```

### Client Organization Changes

**Issue**: Different logic for organizing client operations in TypeSpec can move operations between clients.

**Changelog Pattern**:

```md
- Function `NewManagementClient` has been removed
- Function `*ManagementClient.BeginRestoreVolume` has been removed
- New function `*VolumesClient.BeginRestoreVolume(context.Context, string, string, string, string, *VolumesClientBeginRestoreVolumeOptions) (*runtime.Poller[VolumesClientRestoreVolumeResponse], error)`
```

**Spec Pattern**:

```tsp
namespace Microsoft.ElasticSan;

@action("restore")
@tag("Restore Volumes")
op restoreVolume is ArmResourceActionAsync<Volume, void, Volume>;
```

**Resolution**:

```tsp
@@clientLocation(Microsoft.ElasticSan.restoreVolume, "Management", "go");
```

### User Standardard Operations List Operation

**Issue**: Operations list operation is upgraded to use the standardard library's definition.

**Changelog Pattern**:

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

**Impact**: Low impact for users as this operation is rarely used in SDK.

**Resolution**: Generally acceptable breaking changes that don't require `client.tsp` fixes.

### Common Types Upgrade

**Issue**: Common types are upgraded to latest versions during TypeSpec migration.

**Changelog Pattern**:

```md
- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Type of `Error.Error` has been changed from `*ErrorError` to `*ErrorDetail`
- Type of `SystemData.CreatedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Enum `IdentityType` has been removed
- Struct `ErrorError` has been removed
```

**Impact**: Low impact for users as these are common infrastructure types.

**Resolution**: Generally acceptable breaking changes that don't require `client.tsp` fixes.

### Removal of Unreferenced Types

**Issue**: Unreferenced types are removed during TypeSpec migration.

**Changelog Pattern**:

```md
- Struct `TrackedResource` has been removed
- Struct `Resource` has been removed
- Struct `ProxyResource` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `SCConfluentListMetadata` has been removed
```

**Impact**: Low impact for users as these types are typically not directly used.

**Resolution**: Generally acceptable breaking changes that don't require `client.tsp` fixes.

## TypeSpec Update

These breaking changes result from API specification updates in TypeSpec files. Find corresponding changes in `.tsp` file diffs.

### Naming Changes

#### Model/Enum/Union Name Changes

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

#### Property Name Changes

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

#### Operation Name Changes

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

**Impact**: Method names change in generated client.

**Resolution**:

```tsp
@@clientName(b, "a", "go");
```

#### Enum Value Name Changes

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

### Operation Type Changes

#### Long-Running Operation (LRO) Changes

**Changelog Pattern**:

```md
- Operation `*xxx.A` has been changed to LRO, use `*xxx.BeginA` instead.
- Operation `*xxx.BeginB` has been changed to non-LRO, use `*xxx.B` instead.
```

**Impact**: Method name changes and return type changes (direct result ↔ poller).

**Resolution**: Cannot be resolved with `client.tsp`.

#### Paging Operation Changes

**Changelog Pattern**:

```md
- Function `*xxx.NewListAPager` has been removed
- New function `*xxx.A(xxx) (xxx, error)`
```

**Impact**: Method name changes and return type changes (direct result ↔ pager).

**Resolution**: Cannot be resolved with `client.tsp`.

### Type Changes

#### Property Type Changes

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Parameter Type Changes

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Response Type Changes

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

**Resolution**: Cannot be resolved with `client.tsp`.

### Deletions

#### Property Deletion

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Parameter Deletion

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Operation Deletion

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Model Deletion

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

**Resolution**: Cannot be resolved with `client.tsp`.

### Optionality Changes

#### Required Parameter Addition

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Optional to Required Parameter Changes

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

**Resolution**: Cannot be resolved with `client.tsp`.

#### Required to Optional Parameter Changes

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

**Resolution**: Cannot be resolved with `client.tsp`.
