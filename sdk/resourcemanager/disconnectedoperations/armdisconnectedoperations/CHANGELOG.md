# Release History

## 1.0.0 (2026-02-25)
### Features Added

- New enum type `AutoRenew` with values `AutoRenewDisabled`, `AutoRenewEnabled`
- New enum type `BenefitPlanStatus` with values `BenefitPlanStatusDisabled`, `BenefitPlanStatusEnabled`
- New enum type `BillingStatus` with values `BillingStatusDisabled`, `BillingStatusEnabled`, `BillingStatusStopped`
- New enum type `PricingModel` with values `PricingModelAnnual`, `PricingModelTrial`
- New enum type `SystemReboot` with values `SystemRebootNotRequired`, `SystemRebootRequired`
- New function `*ClientFactory.NewHardwareSettingsClient() *HardwareSettingsClient`
- New function `NewHardwareSettingsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*HardwareSettingsClient, error)`
- New function `*HardwareSettingsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, name string, hardwareSettingName string, resource HardwareSetting, options *HardwareSettingsClientBeginCreateOrUpdateOptions) (*runtime.Poller[HardwareSettingsClientCreateOrUpdateResponse], error)`
- New function `*HardwareSettingsClient.BeginDelete(ctx context.Context, resourceGroupName string, name string, hardwareSettingName string, options *HardwareSettingsClientBeginDeleteOptions) (*runtime.Poller[HardwareSettingsClientDeleteResponse], error)`
- New function `*HardwareSettingsClient.Get(ctx context.Context, resourceGroupName string, name string, hardwareSettingName string, options *HardwareSettingsClientGetOptions) (HardwareSettingsClientGetResponse, error)`
- New function `*HardwareSettingsClient.NewListByParentPager(resourceGroupName string, name string, options *HardwareSettingsClientListByParentOptions) *runtime.Pager[HardwareSettingsClientListByParentResponse]`
- New struct `BenefitPlans`
- New struct `BillingConfiguration`
- New struct `BillingPeriod`
- New struct `HardwareSetting`
- New struct `HardwareSettingListResult`
- New struct `HardwareSettingProperties`
- New struct `ImageUpdateProperties`
- New field `BenefitPlans`, `BillingConfiguration` in struct `DisconnectedOperationDeploymentManifest`
- New field `BenefitPlans`, `BillingConfiguration` in struct `DisconnectedOperationProperties`
- New field `BenefitPlans`, `BillingConfiguration` in struct `DisconnectedOperationUpdateProperties`
- New field `UpdateProperties` in struct `ImageDownloadResult`
- New field `UpdateProperties` in struct `ImageProperties`


## 0.1.0 (2025-09-08)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/disconnectedoperations/armdisconnectedoperations` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).