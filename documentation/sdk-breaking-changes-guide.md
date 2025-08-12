# Azure Go SDK Breaking Changes Review and Resolution Guide

The Azure Go SDK generally prohibits breaking changes unless they result from service behavior modifications. This guide helps you identify, review, and resolve breaking changes that may occur in new SDK versions due to:

- Updates to a service's TypeSpec specifications
- Migration of service specifications from Swagger to TypeSpec

Breaking changes can often be resolved by:

1. Client Customizations

Client customizations should be implemented in a file named `client.tsp` located in the service's specification directory alongside the main entry point `main.tsp`. This `client.tsp` becomes the new specification entry point, so import `main.tsp` in the `client.tsp` file. **Do not** import `client.tsp` in the `main.tsp` file.

```tsp
import "./main.tsp";
import "@azure-tools/typespec-client-generator-core";

using Azure.ClientGenerator.Core;

// Add your customizations here
```

2. TypeSpec Configuration Changes

TypeSpec configuration changes should be made in the `tspconfig.yaml` file located in the service's specification directory. This file is used to configure the TypeSpec compiler and client generator options.

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

## Migration from Swagger to TypeSpec

The following breaking changes commonly occur when migrating API specifications from Swagger to TypeSpec.

### 1. Naming Changes with Numbers

**Changelog Pattern**:

Paired removal and addition entries showing naming changes from words to numbers:

```md
- `MinuteThirty`, `MinuteZero` from enum `Minute` has been removed
- New value `Minute0`, `Minute30` added to enum type `Minute`
```

**Reason**: Swagger automatically converts numeric names to words during code generation, while TypeSpec preserves the original naming. This affects all type names, including enums, models, and operations.

**Spec Pattern**:

Find the type definition by examining the names from the addition entries in the changelog (pattern: `New xxx '<type name>'`):

```tsp
enum Minute {
  Minute0 = "0",
  Minute30 = "30"
}
```

**Resolution**:

Use client customization to restore the original names from the removal entries:

```tsp
@@clientName(Minute.Minute0, "Zero", "go");
@@clientName(Minute.Minute30, "Thirty", "go");
```

### 2. Enum Naming Changes from Anti-Stuttering Rules

**Changelog Pattern**:

Removal of an enum type and addition of a new enum type with the service name prefix removed, along with updates to all references:

```md
- Type of `ConfigurationProperties.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Type of `Update.MaintenanceScope` has been changed from `*MaintenanceScope` to `*Scope`
- Enum `MaintenanceScope` has been removed
- New enum type `Scope` with values `ScopeExtension`, `ScopeHost`, `ScopeInGuestPatch`, `ScopeOSImage`, `ScopeResource`, `ScopeSQLDB`, `ScopeSQLManagedInstance`
```

**Reason**: Differences in enum anti-stuttering rules between Swagger and TypeSpec can cause enum name changes.

**Spec Pattern**:

Find the enum using the name from the removal entries (pattern: `Enum '<enum name>' has been removed`):

```tsp
union MaintenanceScope {
  string,
  Host: "Host",
  Resource: "Resource",
}
```

**Resolution**:

Disable the anti-stuttering rule in TypeSpec config `tspconfig.yaml` to preserve original enum names:

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

### 3. Operation Naming Changes

**Changelog Pattern**:

Removal of an operation and addition of a similarly named operation for the same client:

```md
- Function `*StorageTaskAssignmentClient.NewListPager` has been removed
- New function `*StorageTaskAssignmentClient.NewStorageTaskAssignmentListPager(string, string, *StorageTaskAssignmentClientStorageTaskAssignmentListOptions) *runtime.Pager[StorageTaskAssignmentClientStorageTaskAssignmentListResponse]`
```

**Reason**: TypeSpec may generate different operation names than Swagger to avoid naming collisions.

**Spec Pattern**:

Locate the interface and operation using the name from the addition entries. Operation types include:

- Regular operations: `New function *<interface name>Client.<operation name>(...)`
- Paging operations: `New function *<interface name>Client.New<operation name>Pager(...)`
- Long-running operations: `New function *<interface name>Client.Begin<operation name>(...)`

```tsp
interface StorageTaskAssignment {
  op storageTaskAssignmentList(xxx): xxx;
}
```

**Resolution**:

Use client naming to restore the original operation name from the removal entries:

**Note**: For paging operations, the SDK method name is `New<OperationName>Pager`. For long-running operations, it's `Begin<OperationName>`. When resolving breaking changes, use only the TypeSpec operation name without these SDK-specific prefixes or suffixes.

```tsp
@@clientName(StorageTaskAssignment.storageTaskAssignmentList, "list", "go");
```

### 4. Client Organization Changes

**Changelog Pattern**:

Operations moving between clients, sometimes accompanied by client removal:

```md
- Function `NewManagementClient` has been removed
- Function `*ManagementClient.BeginRestoreVolume` has been removed
- New function `*VolumesClient.BeginRestoreVolume(context.Context, string, string, string, string, *VolumesClientBeginRestoreVolumeOptions) (*runtime.Poller[VolumesClientRestoreVolumeResponse], error)`
```

**Reason**: TypeSpec uses different logic for organizing client operations than Swagger.

**Spec Pattern**:

Find the interface and operation using the name from the addition entries (pattern: `New function *<interface name>Client.<operation name>(...)`):

```tsp
namespace Microsoft.ElasticSan;

interface Volumes {
  @action("restore")
  op restoreVolume is ArmResourceActionAsync<Volume, void, Volume>;
}
```

**Resolution**:

Move the operation to the correct client using `@@clientLocation`. Use the client name from the removal entries (removing the `Client` suffix):

```tsp
@@clientLocation(Microsoft.ElasticSan.restoreVolume, "Management", "go");
```

### 5. Missing Fields in Response Types

**Changelog Pattern**:

Removal of fields in response structures with the `xxxResponse` naming pattern:

```md
- Field `CacheAccessPolicyAssignment` of struct `AccessPolicyAssignmentClientCreateUpdateResponse` has been removed
```

**Reason**: Incorrect TypeSpec conversion for long-running operation (LRO) responses.

**Spec Pattern**:

Find the interface and operation using the name from the removal entries (pattern: `Field 'xxx' of struct *<interface name>Client<operation name>Response`):

```tsp
@armResourceOperations
interface RedisCacheAccessPolicies {
  createUpdate is ArmResourceCreateOrReplaceAsync<
    RedisCacheAccessPolicy,
    LroHeaders = ArmLroLocationHeader & Azure.Core.Foundations.RetryAfterHeader
  >;
}
```

If the response type name does not have the `<interface name>` prefix and start with `Client` directly, use the service name as the interface name instead.

**Resolution**:

Locate the operation and add the `FinalResult` parameter to the appropriate LRO header (`ArmLroLocationHeader`, `ArmAsyncOperationHeader`, or `ArmCombinedLroHeaders`) with the correct type:

```tsp
@armResourceOperations
interface RedisCacheAccessPolicies {
  createUpdate is ArmResourceCreateOrReplaceAsync<
    RedisCacheAccessPolicy,
    LroHeaders = ArmLroLocationHeader<FinalResult = CacheAccessPolicyAssignment> & Azure.Core.Foundations.RetryAfterHeader
  >;
}
```

### 6. Standard Operations List Operation Upgrade

**Changelog Pattern**:

Multiple changes related to the `Operation` type and its fields, sometimes including changes to the `OperationList` operation:

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

**Reason**: The operations list operation is upgraded to use the standard library definition.

**Impact**: Low impact since this operation is rarely used in the SDK.

**Resolution**: Accept these breaking changes.

### 7. Common Types Upgrade

**Changelog Pattern**:

Multiple changes related to common infrastructure types such as `SystemData`, `Error`, and `IdentityType`:

```md
- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Type of `Error.Error` has been changed from `*ErrorError` to `*ErrorDetail`
- Type of `SystemData.CreatedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Enum `IdentityType` has been removed
- Struct `ErrorError` has been removed
```

**Reason**: Common types are upgraded to their latest versions during TypeSpec migration.

**Impact**: Low impact since these are common infrastructure types rarely used directly by users.

**Resolution**: Accept these breaking changes.

### 8. Removal of Unreferenced Types

**Changelog Pattern**:

Multiple removals of unreferenced types that are typically not used in the SDK:

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

**Impact**: Low impact since these types are typically not used directly by users.

**Resolution**: Accept these breaking changes.

### 9. Request Body Optionality Changes

**Changelog Pattern**:

An additional parameter is added to an operation, and a corresponding field is removed from the operation's options struct:

```md
- Function `*MarketplaceAgreementsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, *MarketplaceAgreementsClientCreateOrUpdateOptions)` to `(context.Context, AgreementResource, *MarketplaceAgreementsClientCreateOrUpdateOptions)`
- `*MonitorsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, *MonitorsClientBeginCreateOptions)` to `(context.Context, string, string, MonitorResource, *MonitorsClientBeginCreateOptions)`
- Field `Body` of struct `MarketplaceAgreementsClientCreateOrUpdateOptions` has been removed
- Field `Body` of struct `MonitorsClientBeginCreateOptions` has been removed
```

**Reason**: For PUT and PATCH operations, the request body is always treated as required in TypeSpec.

**Impact**: This corrects the previous SDK behavior.

**Resolution**: Accept these breaking changes.

### 10. Naming Changes from Directive

**Changelog Pattern**:

Paired removal and addition entries showing naming changes for structs:

```md
- Struct `ResourceInfo` has been removed
- New struct `RedisResource`
```

Also, in the legacy config for swagger under the spec folder: `specification/<service>/resource-manager/readme.go.md`, the renaming directives could be found:

```md
directive:

- rename-model:
  from: 'RedisResource'
  to: 'ResourceInfo'
```

**Reason**: Swagger has directive ways to change the naming.

**Spec Pattern**:

Find the type definition by examining the names from the addition entries in the changelog (pattern: `New xxx '<type name>'`):

```tsp
model RedisResource {
  ...
}
```

**Resolution**:

Use client customization to do the same renaming as the directive in the legacy config:

```tsp
@@clientName(RedisResource, "ResourceInfo", "go");
```

## TypeSpec Specification Updates

The following breaking changes result from updates to existing TypeSpec specifications.

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

**Impact**: Method names and return types change (direct result ↔ poller).

**Resolution**: Cannot be resolved through client customizations.

### 6. Paging Operation Changes

**Changelog Pattern**:

```md
- Function `*xxx.NewListAPager` has been removed
- New function `*xxx.A(xxx) (xxx, error)`
```

**Impact**: Method names and return types change (direct result ↔ pager).

**Resolution**: Cannot be resolved through client customizations.

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

**Resolution**: Cannot be resolved through client customizations.

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

**Resolution**: Cannot be resolved through client customizations.

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

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Fields are no longer available, causing compilation errors.

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Parameters are no longer available, requiring method signature updates.

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Client methods are no longer available, requiring alternative implementation.

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Types are no longer available, requiring alternative type usage.

**Resolution**: Cannot be resolved through client customizations.

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

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Previously optional parameters become mandatory, requiring parameter updates.

**Resolution**: Cannot be resolved through client customizations.

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

**Impact**: Previously mandatory parameters become optional, potentially affecting validation logic.

**Resolution**: Cannot be resolved through client customizations.
