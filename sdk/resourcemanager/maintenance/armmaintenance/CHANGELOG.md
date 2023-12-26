# Release History

## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-08-25)
### Features Added

- New value `MaintenanceScopeResource` added to enum type `MaintenanceScope`
- New enum type `RebootOptions` with values `RebootOptionsAlways`, `RebootOptionsIfRequired`, `RebootOptionsNever`
- New enum type `TagOperators` with values `TagOperatorsAll`, `TagOperatorsAny`
- New function `*ClientFactory.NewConfigurationAssignmentsForResourceGroupClient() *ConfigurationAssignmentsForResourceGroupClient`
- New function `*ClientFactory.NewConfigurationAssignmentsForSubscriptionsClient() *ConfigurationAssignmentsForSubscriptionsClient`
- New function `*ClientFactory.NewConfigurationAssignmentsWithinSubscriptionClient() *ConfigurationAssignmentsWithinSubscriptionClient`
- New function `*ConfigurationAssignmentsClient.Get(context.Context, string, string, string, string, string, *ConfigurationAssignmentsClientGetOptions) (ConfigurationAssignmentsClientGetResponse, error)`
- New function `*ConfigurationAssignmentsClient.GetParent(context.Context, string, string, string, string, string, string, string, *ConfigurationAssignmentsClientGetParentOptions) (ConfigurationAssignmentsClientGetParentResponse, error)`
- New function `NewConfigurationAssignmentsForResourceGroupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConfigurationAssignmentsForResourceGroupClient, error)`
- New function `*ConfigurationAssignmentsForResourceGroupClient.CreateOrUpdate(context.Context, string, string, ConfigurationAssignment, *ConfigurationAssignmentsForResourceGroupClientCreateOrUpdateOptions) (ConfigurationAssignmentsForResourceGroupClientCreateOrUpdateResponse, error)`
- New function `*ConfigurationAssignmentsForResourceGroupClient.Delete(context.Context, string, string, *ConfigurationAssignmentsForResourceGroupClientDeleteOptions) (ConfigurationAssignmentsForResourceGroupClientDeleteResponse, error)`
- New function `*ConfigurationAssignmentsForResourceGroupClient.Get(context.Context, string, string, *ConfigurationAssignmentsForResourceGroupClientGetOptions) (ConfigurationAssignmentsForResourceGroupClientGetResponse, error)`
- New function `*ConfigurationAssignmentsForResourceGroupClient.Update(context.Context, string, string, ConfigurationAssignment, *ConfigurationAssignmentsForResourceGroupClientUpdateOptions) (ConfigurationAssignmentsForResourceGroupClientUpdateResponse, error)`
- New function `NewConfigurationAssignmentsForSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConfigurationAssignmentsForSubscriptionsClient, error)`
- New function `*ConfigurationAssignmentsForSubscriptionsClient.CreateOrUpdate(context.Context, string, ConfigurationAssignment, *ConfigurationAssignmentsForSubscriptionsClientCreateOrUpdateOptions) (ConfigurationAssignmentsForSubscriptionsClientCreateOrUpdateResponse, error)`
- New function `*ConfigurationAssignmentsForSubscriptionsClient.Delete(context.Context, string, *ConfigurationAssignmentsForSubscriptionsClientDeleteOptions) (ConfigurationAssignmentsForSubscriptionsClientDeleteResponse, error)`
- New function `*ConfigurationAssignmentsForSubscriptionsClient.Get(context.Context, string, *ConfigurationAssignmentsForSubscriptionsClientGetOptions) (ConfigurationAssignmentsForSubscriptionsClientGetResponse, error)`
- New function `*ConfigurationAssignmentsForSubscriptionsClient.Update(context.Context, string, ConfigurationAssignment, *ConfigurationAssignmentsForSubscriptionsClientUpdateOptions) (ConfigurationAssignmentsForSubscriptionsClientUpdateResponse, error)`
- New function `NewConfigurationAssignmentsWithinSubscriptionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConfigurationAssignmentsWithinSubscriptionClient, error)`
- New function `*ConfigurationAssignmentsWithinSubscriptionClient.NewListPager(*ConfigurationAssignmentsWithinSubscriptionClientListOptions) *runtime.Pager[ConfigurationAssignmentsWithinSubscriptionClientListResponse]`
- New struct `ConfigurationAssignmentFilterProperties`
- New struct `InputLinuxParameters`
- New struct `InputPatchConfiguration`
- New struct `InputWindowsParameters`
- New struct `TagSettingsProperties`
- New field `Filter` in struct `ConfigurationAssignmentProperties`
- New field `InstallPatches` in struct `ConfigurationProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/maintenance/armmaintenance` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).