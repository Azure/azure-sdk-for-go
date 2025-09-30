# Azure Go SDK Breaking Changes Review and Resolution Guide for TypeSpec Migration

The Azure Go SDK generally prohibits breaking changes unless they result from service behavior modifications. This guide helps you identify, review, and resolve breaking changes that may occur in new SDK versions due to migrating of service specifications from Swagger to TypeSpec. For service's TypeSpec specification update scenario, refer this [doc](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/sdk-breaking-changes-guide.md).

Some breaking changes can be accepted as they have low impact on users. Some can be resolved through client customizations or TypeSpec configuration changes:

1. Client Customizations

Client customizations should be implemented in a file named `client.tsp` located in the service's specification directory alongside the main entry point `main.tsp`. This `client.tsp` becomes the new specification entry point, so import `main.tsp` in the `client.tsp` file. **Do not** import `client.tsp` in the `main.tsp` file. **Do not** modify the entry point in `tspconfig.yaml`.

```tsp
import "./main.tsp";
import "@azure-tools/typespec-client-generator-core";

using Azure.ClientGenerator.Core;
using MainNamespaceInMainTsp; // Replace with the actual main namespace in main.tsp

// Add your customizations here
```

2. TypeSpec Configuration Changes

TypeSpec configuration changes should be made in the `tspconfig.yaml` file located in the service's specification directory. This file is used to configure the TypeSpec compiler and client generator options. For example:

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

## 1. Naming Changes with Numbers

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

## 2. Enum Naming Changes from Anti-Stuttering Rules

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

## 3. Operation Naming Changes

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

## 4. Client Organization Changes

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

## 5. Missing Fields in Response Types

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

## 6. Standard Operations List Operation Upgrade

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

**Migration Guide**: Change to use the new `OperationList` operation and the new related types.

## 7. Common Types Upgrade

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

**Migration Guide**: Change to use the new types.

## 8. Removal of Unreferenced Types

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

## 9. Request Body Optionality Changes

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

**Migration Guide**: Update the code to pass the request body as a separate parameter instead of including it in the options struct.

## 10. Naming Changes from Directive

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

Use client customization to do the same renaming as the directives in the legacy config:

```tsp
@@clientName(RedisResource, "ResourceInfo", "go");
```

## 11. Model Naming Changes from Anti-Stuttering Rules

**Changelog Pattern**:

Removal of a `xxxListResult` model, addition of a `xxxListListResult` model and change of related fields:

```md
- Struct `DomainListResult` has been removed
- Field `DomainListResult` of struct `DomainListsClientListByResourceGroupResponse` has been removed
- Field `DomainListResult` of struct `DomainListsClientListResponse` has been removed
- New struct `DomainListListResult`
- New anonymous field `DomainListListResult` in struct `DomainListsClientListByResourceGroupResponse`
- New anonymous field `DomainListListResult` in struct `DomainListsClientListResponse`
```

**Reason**: Swagger has a naming magic to remove stuttering part of the type names. When we migrate to TypeSpec, we want to keep the original names without the magic to avoid confusion.

**Impact**: Low impact since list structs are rarely used directly by users.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Change to use the new structs.

## 12. Type Changes for Enum Values

**Changelog Pattern**:

Removal of enum values and addition of new enum values with the new enum type:

```md
- `ActionTypeEnable`, `ActionTypeOptOut` from enum `ActionType` has been removed
- New enum type `ActionTypeFlag` with values `ActionTypeFlagEnable`, `ActionTypeFlagOptOut`
```

**Reason**: Swagger merges the enum values of enum type with same name. This is incorrect. When migrating to TypeSpec, we fix it with a new enum type.

**Impact**: This corrects the previous SDK behavior.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Update the code to use the new enum types.

## 13. Type Changes for Enum Values

**Changelog Pattern**:

Removal of an enum type and change the refer of this enum type to string:

```md
- Type of `MessageProperties.ContentType` has been changed from `*TranscriptContentType` to `*string`
- Enum `TranscriptContentType` has been removed
- Function `PossibleTranscriptContentTypeValues` has been removed
```

**Reason**: Swagger allows extensible enum without any known value. This is incorrect. When migrating to TypeSpec, we change it to string directly.

**Impact**: This corrects the previous SDK behavior.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Update the code to remove the type casing.