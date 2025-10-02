# Release History

## 1.0.0 (2025-09-26)
### Breaking Changes

- Type of `SchedulerSKU.Name` has been changed from `*string` to `*SchedulerSKUName`
- Type of `SchedulerSKUUpdate.Name` has been changed from `*string` to `*SchedulerSKUName`

### Features Added

- New enum type `SchedulerSKUName` with values `SchedulerSKUNameConsumption`, `SchedulerSKUNameDedicated`


## 0.2.0 (2025-04-15)
### Features Added

- New enum type `PurgeableOrchestrationState` with values `PurgeableOrchestrationStateCanceled`, `PurgeableOrchestrationStateCompleted`, `PurgeableOrchestrationStateFailed`, `PurgeableOrchestrationStateTerminated`
- New function `*ClientFactory.NewRetentionPoliciesClient() *RetentionPoliciesClient`
- New function `NewRetentionPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RetentionPoliciesClient, error)`
- New function `*RetentionPoliciesClient.BeginCreateOrReplace(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginCreateOrReplaceOptions) (*runtime.Poller[RetentionPoliciesClientCreateOrReplaceResponse], error)`
- New function `*RetentionPoliciesClient.BeginDelete(context.Context, string, string, *RetentionPoliciesClientBeginDeleteOptions) (*runtime.Poller[RetentionPoliciesClientDeleteResponse], error)`
- New function `*RetentionPoliciesClient.Get(context.Context, string, string, *RetentionPoliciesClientGetOptions) (RetentionPoliciesClientGetResponse, error)`
- New function `*RetentionPoliciesClient.NewListBySchedulerPager(string, string, *RetentionPoliciesClientListBySchedulerOptions) *runtime.Pager[RetentionPoliciesClientListBySchedulerResponse]`
- New function `*RetentionPoliciesClient.BeginUpdate(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginUpdateOptions) (*runtime.Poller[RetentionPoliciesClientUpdateResponse], error)`
- New struct `RetentionPolicy`
- New struct `RetentionPolicyDetails`
- New struct `RetentionPolicyListResult`
- New struct `RetentionPolicyProperties`


## 0.1.0 (2025-03-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/durabletask/armdurabletask` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).