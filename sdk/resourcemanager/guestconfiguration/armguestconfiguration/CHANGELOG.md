# Release History

## 2.0.0 (2026-03-27)
### Breaking Changes

- Function `*AssignmentReportsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, reportID string, vmName string, options *AssignmentReportsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, reportID string, options *AssignmentReportsClientGetOptions)`
- Function `*AssignmentReportsClient.List` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, vmName string, options *AssignmentReportsClientListOptions)` to `(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *AssignmentReportsClientListOptions)`
- Function `*AssignmentsClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, guestConfigurationAssignmentName string, resourceGroupName string, vmName string, parameters Assignment, options *AssignmentsClientCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, parameters Assignment, options *AssignmentsClientCreateOrUpdateOptions)`
- Function `*AssignmentsClient.Delete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, vmName string, options *AssignmentsClientDeleteOptions)` to `(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *AssignmentsClientDeleteOptions)`
- Function `*AssignmentsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, vmName string, options *AssignmentsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *AssignmentsClientGetOptions)`
- Function `*HCRPAssignmentReportsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, reportID string, machineName string, options *HCRPAssignmentReportsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, machineName string, guestConfigurationAssignmentName string, reportID string, options *HCRPAssignmentReportsClientGetOptions)`
- Function `*HCRPAssignmentReportsClient.List` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, machineName string, options *HCRPAssignmentReportsClientListOptions)` to `(ctx context.Context, resourceGroupName string, machineName string, guestConfigurationAssignmentName string, options *HCRPAssignmentReportsClientListOptions)`
- Function `*HCRPAssignmentsClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, guestConfigurationAssignmentName string, resourceGroupName string, machineName string, parameters Assignment, options *HCRPAssignmentsClientCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, machineName string, guestConfigurationAssignmentName string, parameters Assignment, options *HCRPAssignmentsClientCreateOrUpdateOptions)`
- Function `*HCRPAssignmentsClient.Delete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, machineName string, options *HCRPAssignmentsClientDeleteOptions)` to `(ctx context.Context, resourceGroupName string, machineName string, guestConfigurationAssignmentName string, options *HCRPAssignmentsClientDeleteOptions)`
- Function `*HCRPAssignmentsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, machineName string, options *HCRPAssignmentsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, machineName string, guestConfigurationAssignmentName string, options *HCRPAssignmentsClientGetOptions)`
- Struct `Resource` has been removed

### Features Added

- New function `*AssignmentsVMSSClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, vmssName string, name string, parameters Assignment, options *AssignmentsVMSSClientCreateOrUpdateOptions) (AssignmentsVMSSClientCreateOrUpdateResponse, error)`
- New function `*ClientFactory.NewConnectedVMwarevSphereAssignmentsClient() *ConnectedVMwarevSphereAssignmentsClient`
- New function `*ClientFactory.NewConnectedVMwarevSphereAssignmentsReportsClient() *ConnectedVMwarevSphereAssignmentsReportsClient`
- New function `NewConnectedVMwarevSphereAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConnectedVMwarevSphereAssignmentsClient, error)`
- New function `*ConnectedVMwarevSphereAssignmentsClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, parameters Assignment, options *ConnectedVMwarevSphereAssignmentsClientCreateOrUpdateOptions) (ConnectedVMwarevSphereAssignmentsClientCreateOrUpdateResponse, error)`
- New function `*ConnectedVMwarevSphereAssignmentsClient.Delete(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *ConnectedVMwarevSphereAssignmentsClientDeleteOptions) (ConnectedVMwarevSphereAssignmentsClientDeleteResponse, error)`
- New function `*ConnectedVMwarevSphereAssignmentsClient.Get(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *ConnectedVMwarevSphereAssignmentsClientGetOptions) (ConnectedVMwarevSphereAssignmentsClientGetResponse, error)`
- New function `*ConnectedVMwarevSphereAssignmentsClient.NewListPager(resourceGroupName string, vmName string, options *ConnectedVMwarevSphereAssignmentsClientListOptions) *runtime.Pager[ConnectedVMwarevSphereAssignmentsClientListResponse]`
- New function `NewConnectedVMwarevSphereAssignmentsReportsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConnectedVMwarevSphereAssignmentsReportsClient, error)`
- New function `*ConnectedVMwarevSphereAssignmentsReportsClient.Get(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, reportID string, options *ConnectedVMwarevSphereAssignmentsReportsClientGetOptions) (ConnectedVMwarevSphereAssignmentsReportsClientGetResponse, error)`
- New function `*ConnectedVMwarevSphereAssignmentsReportsClient.List(ctx context.Context, resourceGroupName string, vmName string, guestConfigurationAssignmentName string, options *ConnectedVMwarevSphereAssignmentsReportsClientListOptions) (ConnectedVMwarevSphereAssignmentsReportsClientListResponse, error)`
- New field `NextLink` in struct `AssignmentList`
- New field `NextLink` in struct `AssignmentReportList`
- New field `ContentManagedIdentity` in struct `Navigation`
- New field `NextLink` in struct `OperationList`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/guestconfiguration/armguestconfiguration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).