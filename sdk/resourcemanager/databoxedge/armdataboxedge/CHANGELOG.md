# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `*AddonsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, roleName string, addonName string, resourceGroupName string, addon AddonClassification, options *AddonsClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, roleName string, addonName string, resourceGroupName string, resource AddonClassification, options *AddonsClientBeginCreateOrUpdateOptions)`
- Function `*BandwidthSchedulesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, name string, resourceGroupName string, parameters BandwidthSchedule, options *BandwidthSchedulesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, name string, resourceGroupName string, resource BandwidthSchedule, options *BandwidthSchedulesClientBeginCreateOrUpdateOptions)`
- Function `*ContainersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, storageAccountName string, containerName string, resourceGroupName string, containerParam Container, options *ContainersClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, storageAccountName string, containerName string, resourceGroupName string, resource Container, options *ContainersClientBeginCreateOrUpdateOptions)`
- Function `*DevicesClient.BeginCreateOrUpdateSecuritySettings` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, securitySettings SecuritySettings, options *DevicesClientBeginCreateOrUpdateSecuritySettingsOptions)` to `(ctx context.Context, deviceName string, resourceGroupName string, body SecuritySettings, options *DevicesClientBeginCreateOrUpdateSecuritySettingsOptions)`
- Function `*DevicesClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, dataBoxEdgeDevice Device, options *DevicesClientCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, dataBoxEdgeDevice Device, options *DevicesClientCreateOrUpdateOptions)`
- Function `*DevicesClient.Update` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, parameters DevicePatch, options *DevicesClientUpdateOptions)` to `(ctx context.Context, deviceName string, resourceGroupName string, properties DevicePatch, options *DevicesClientUpdateOptions)`
- Function `*DevicesClient.UpdateExtendedInformation` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, parameters DeviceExtendedInfoPatch, options *DevicesClientUpdateExtendedInformationOptions)` to `(ctx context.Context, deviceName string, resourceGroupName string, body DeviceExtendedInfoPatch, options *DevicesClientUpdateExtendedInformationOptions)`
- Function `*DevicesClient.UploadCertificate` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, parameters UploadCertificateRequest, options *DevicesClientUploadCertificateOptions)` to `(ctx context.Context, deviceName string, resourceGroupName string, body UploadCertificateRequest, options *DevicesClientUploadCertificateOptions)`
- Function `*DiagnosticSettingsClient.BeginUpdateDiagnosticProactiveLogCollectionSettings` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, proactiveLogCollectionSettings DiagnosticProactiveLogCollectionSettings, options *DiagnosticSettingsClientBeginUpdateDiagnosticProactiveLogCollectionSettingsOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, proactiveLogCollectionSettings DiagnosticProactiveLogCollectionSettings, options *DiagnosticSettingsClientBeginUpdateDiagnosticProactiveLogCollectionSettingsOptions)`
- Function `*DiagnosticSettingsClient.BeginUpdateDiagnosticRemoteSupportSettings` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, diagnosticRemoteSupportSettings DiagnosticRemoteSupportSettings, options *DiagnosticSettingsClientBeginUpdateDiagnosticRemoteSupportSettingsOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, diagnosticRemoteSupportSettings DiagnosticRemoteSupportSettings, options *DiagnosticSettingsClientBeginUpdateDiagnosticRemoteSupportSettingsOptions)`
- Function `*DiagnosticSettingsClient.GetDiagnosticProactiveLogCollectionSettings` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, options *DiagnosticSettingsClientGetDiagnosticProactiveLogCollectionSettingsOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, options *DiagnosticSettingsClientGetDiagnosticProactiveLogCollectionSettingsOptions)`
- Function `*DiagnosticSettingsClient.GetDiagnosticRemoteSupportSettings` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, options *DiagnosticSettingsClientGetDiagnosticRemoteSupportSettingsOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, options *DiagnosticSettingsClientGetDiagnosticRemoteSupportSettingsOptions)`
- Function `*OrdersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, order Order, options *OrdersClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, resourceGroupName string, resource Order, options *OrdersClientBeginCreateOrUpdateOptions)`
- Function `*RolesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, name string, resourceGroupName string, role RoleClassification, options *RolesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, name string, resourceGroupName string, resource RoleClassification, options *RolesClientBeginCreateOrUpdateOptions)`
- Function `*SharesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, name string, resourceGroupName string, share Share, options *SharesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, name string, resourceGroupName string, resource Share, options *SharesClientBeginCreateOrUpdateOptions)`
- Function `*StorageAccountCredentialsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, name string, resourceGroupName string, storageAccountCredential StorageAccountCredential, options *StorageAccountCredentialsClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, name string, resourceGroupName string, resource StorageAccountCredential, options *StorageAccountCredentialsClientBeginCreateOrUpdateOptions)`
- Function `*SupportPackagesClient.BeginTriggerSupportPackage` parameter(s) have been changed from `(ctx context.Context, deviceName string, resourceGroupName string, triggerSupportPackageRequest TriggerSupportPackageRequest, options *SupportPackagesClientBeginTriggerSupportPackageOptions)` to `(ctx context.Context, resourceGroupName string, deviceName string, triggerSupportPackageRequest TriggerSupportPackageRequest, options *SupportPackagesClientBeginTriggerSupportPackageOptions)`
- Function `*TriggersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, deviceName string, name string, resourceGroupName string, trigger TriggerClassification, options *TriggersClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, deviceName string, name string, resourceGroupName string, resource TriggerClassification, options *TriggersClientBeginCreateOrUpdateOptions)`
- Struct `ARMBaseModel` has been removed
- Struct `MoveRequest` has been removed

### Features Added

- New field `KubernetesWorkloadProfile` in struct `DeviceProperties`
- New field `SystemData` in struct `Job`
- New field `IPRange` in struct `LoadBalancerConfig`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databoxedge/armdataboxedge` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).