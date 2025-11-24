# Azure Go SDK Breaking Changes Review and Resolution Guide for TypeSpec Migration

This guide helps you identify, review, and resolve breaking changes specific to migrating service specifications from Swagger to TypeSpec. For breaking changes from general TypeSpec specification updates, refer to the [main breaking changes guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/breaking-changes/sdk-breaking-changes-guide.md).

Breaking changes should be resolved through TypeSpec client and/or configuration customizations. Low-impact breaking changes can be reviewed by Go architects as they may be acceptable.

## TypeSpec Configuration Changes

TypeSpec configuration changes should be made in the `tspconfig.yaml` file located in the service's specification directory. This file configures the TypeSpec compiler and client generator options. For example:

```yaml
options:
  "@azure-tools/typespec-go":
    fix-const-stuttering: false
```

For client customizations and `client.tsp` setup instructions, refer to the [main breaking changes guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/breaking-changes/sdk-breaking-changes-guide.md).

## Breaking Changes That Can Be Resolved with TypeSpec Customizations

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

Disable the anti-stuttering rule in the TypeSpec config `tspconfig.yaml` to preserve original enum names:

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

Find the interface and operation using the name from the addition entries. Operation types include:

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

**Reason**: Incorrect TypeSpec conversion for Long-Running Operation (LRO) responses.

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

If the response type name does not have the `<interface name>` prefix and starts with `Client` directly, use the service name as the interface name instead.

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

### 6. Naming Changes from Directives

**Changelog Pattern**:

Paired removal and addition entries showing naming changes for structs:

```md
- Struct `ResourceInfo` has been removed
- New struct `RedisResource`
```

Also, in the legacy configuration for Swagger under the spec folder: `specification/<service>/resource-manager/readme.go.md`, the renaming directives can be found:

```md
directive:

- rename-model:
  from: 'RedisResource'
  to: 'ResourceInfo'
```

**Reason**: Swagger has directives to change the naming.

**Spec Pattern**:

Find the type definition by examining the names from the addition entries in the changelog (pattern: `New xxx '<type name>'`):

```tsp
model RedisResource {
  ...
}
```

**Resolution**:

Use client customization to perform the same renaming as the directives in the legacy configuration:

```tsp
@@clientName(RedisResource, "ResourceInfo", "go");
```

### 7. Type Changed from `string` to Another Enum Type

**Changelog Pattern**:

One property type changes from `*string` to another enum type, along with a newly added enum type:

```md
- Type of `RegistryNameCheckRequest.Type` has been changed from `*string` to `*ResourceType`
- New enum type `ResourceType` with values `ResourceTypeMicrosoftContainerRegistryRegistries`
```

**Reason**: In Swagger, we converted single-value fixed enums to constants, but in TypeSpec we need to ensure the same behavior with client customization.

**Spec Pattern**:

Find the model property and enum using the name from the changelog (pattern: `Type of <model name>.<property name>` has been changed from *string to *<enum name>):

```tsp
model RegistryNameCheckRequest {
  type: ContainerRegistryResourceType;
}

enum ContainerRegistryResourceType {
  `Microsoft.ContainerRegistry/registries`,
}
```

**Resolution**:

Locate the model property and use `@@alternateType` to change the property type back to the constant string:

```tsp
@@alternateType(RegistryNameCheckRequest.type, "Microsoft.ContainerRegistry/registries", "go");
```

## Breaking Changes That Can Be Accepted

All these breaking changes will be released in a new major version, except the last one about unreferenced types.

### 1. Operations List Operation Upgrade

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

**Migration Guide**: Change to use the new `OperationList` operation and related types.

For example:

Previous code:

```go
pager := clientFactory.NewOperationsClient().NewListPager(nil)
for pager.More() {
  page, err := pager.NextPage(ctx)
  if err != nil {
    log.Fatalf("failed to advance page: %v", err)
  }
  for _, v := range page.Value {
    if *v.Origin == "system"{
      // ...
    }
  }
}
```

New code:

```go
pager := clientFactory.NewOperationsClient().NewListPager(nil)
for pager.More() {
  page, err := pager.NextPage(ctx)
  if err != nil {
    log.Fatalf("failed to advance page: %v", err)
  }
  for _, v := range page.Value {
    if *v.Origin == OriginSystem {
      // ...
    }
  }
}
```

### 2. Common Types Upgrade

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

For example:

Previous code:

```go
if *resource.SystemData.CreatedByType == IdentityTypeUser {
  // ...
}
```

New code:

```go
if *resource.SystemData.CreatedByType == CreatedByTypeUser {
  // ...
}
```

### 3. Request Body Optionality Changes

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

For example:

Previous code:

```go
client.CreateOrUpdate(ctx, &armdatadog.MarketplaceAgreementsClientCreateOrUpdateOptions{
  Body: &armdatadog.AgreementResource{
    Properties: &armdatadog.AgreementProperties{
      PlanID:   to.Ptr("plan-id"),
      Product:  to.Ptr("product"),
      Publisher: to.Ptr("publisher"),
      Terms:    to.Ptr("terms"),
    },
  },
})
```

New code:

```go
client.CreateOrUpdate(ctx, armdatadog.AgreementResource{
  Properties: &armdatadog.AgreementProperties{
    PlanID:   to.Ptr("plan-id"),
    Product:  to.Ptr("product"),
    Publisher: to.Ptr("publisher"),
    Terms:    to.Ptr("terms"),
  },
}, nil)
```

### 4. Model Naming Changes from Anti-Stuttering Rules

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

**Reason**: Swagger has naming logic to remove the stuttering part of type names. When we migrate to TypeSpec, we want to keep the original names without this logic to avoid confusion.

**Impact**: Low impact since list structs are rarely used directly by users.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Change to use the new structs.

For example:

Previous code:

```go
pager := clientFactory.NewDomainListsClient().NewListByResourceGroupPager("rg1", nil)
for pager.More() {
  page, err := pager.NextPage(ctx)
  if err != nil {
    log.Fatalf("failed to advance page: %v", err)
  }
  for _, v := range page.DomainListResult.Value {
    // ...
  }
}
```

New code:

```go
pager := clientFactory.NewDomainListsClient().NewListByResourceGroupPager("rg1", nil)
for pager.More() {
  page, err := pager.NextPage(ctx)
  if err != nil {
    log.Fatalf("failed to advance page: %v", err)
  }
  for _, v := range page.DomainListListResult.Value {
    // ...
  }
}
```

### 5. Enum Splitting

**Changelog Pattern**:

Removal of enum values and addition of new enum values with the new enum type:

```md
- `ActionTypeEnable`, `ActionTypeOptOut` from enum `ActionType` has been removed
- New enum type `ActionTypeFlag` with values `ActionTypeFlagEnable`, `ActionTypeFlagOptOut`
```

**Reason**: Swagger merges the enum values of enum types with the same name. This is incorrect. When migrating to TypeSpec, we fix this with a new enum type.

**Impact**: This corrects the previous SDK behavior.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Update the code to use the new enum types.

For example:

Previous code:

```go
if *resource.ActionType == ActionTypeEnable {
  // ...
}
```

New code:

```go
if *resource.ActionType == ActionTypeFlagEnable {
  // ...
}
```

### 6. Type Changes for Enum Values

**Changelog Pattern**:

Removal of an enum type and change the refer of this enum type to string:

```md
- Type of `MessageProperties.ContentType` has been changed from `*TranscriptContentType` to `*string`
- Enum `TranscriptContentType` has been removed
- Function `PossibleTranscriptContentTypeValues` has been removed
```

**Reason**: Swagger allows extensible enums without any known values. This is incorrect. When migrating to TypeSpec, we change this to string directly.

**Impact**: This corrects the previous SDK behavior.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Update the code to remove type casting.

For example:

Previous code:

```go
if *message.ContentType == TranscriptContentTypePlainText {
  // ...
}
```

New code:

```go
if *message.ContentType == "PlainText" {
  // ...
}
```

### 7. Change type from `string` to `azore.ETag`

**Changelog Pattern**:

Type change for ETag fields from `*string` to `*azcore.ETag`:

```md
- Type of `PrivateEndpointConnection.Etag` has been changed from `*string` to `*azcore.ETag`
```

**Reason**: When migrating to TypeSpec, we change ETag fields from `*string` to `azcore.ETag`.

**Impact**: Low impact since underlaying type of `azure.Etag` is `string`.

**Resolution**: Accept these breaking changes.

**Migration Guide**: Update the code to remove type casting.

For example:

Previous code:

```go
privateEndpointConnection.Etag = to.Ptr("*")
```

New code:

```go
privateEndpointConnection.Etag = to.Ptr(azcore.ETag("*"))
```

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

**Impact**: No impact since these types are typically not used directly by users.

**Resolution**: Accept these breaking changes.

**Caution**: Since these types are unreferenced, their removal should not affect existing code. We will release these changes without a new major version.
