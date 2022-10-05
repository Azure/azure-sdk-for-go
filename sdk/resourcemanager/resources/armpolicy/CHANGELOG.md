# Release History

## 0.7.0 (2022-10-05)
### Breaking Changes

- Struct `CloudError` has been removed

### Features Added

- New function `NewVariablesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VariablesClient, error)`
- New function `*VariablesClient.Delete(context.Context, string, *VariablesClientDeleteOptions) (VariablesClientDeleteResponse, error)`
- New function `*VariablesClient.CreateOrUpdate(context.Context, string, Variable, *VariablesClientCreateOrUpdateOptions) (VariablesClientCreateOrUpdateResponse, error)`
- New function `*VariablesClient.GetAtManagementGroup(context.Context, string, string, *VariablesClientGetAtManagementGroupOptions) (VariablesClientGetAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.NewListPager(string, *VariableValuesClientListOptions) *runtime.Pager[VariableValuesClientListResponse]`
- New function `*VariableValuesClient.NewListForManagementGroupPager(string, string, *VariableValuesClientListForManagementGroupOptions) *runtime.Pager[VariableValuesClientListForManagementGroupResponse]`
- New function `*VariableValuesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, string, VariableValue, *VariableValuesClientCreateOrUpdateAtManagementGroupOptions) (VariableValuesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariablesClient.NewListForManagementGroupPager(string, *VariablesClientListForManagementGroupOptions) *runtime.Pager[VariablesClientListForManagementGroupResponse]`
- New function `*VariablesClient.NewListPager(*VariablesClientListOptions) *runtime.Pager[VariablesClientListResponse]`
- New function `*VariablesClient.DeleteAtManagementGroup(context.Context, string, string, *VariablesClientDeleteAtManagementGroupOptions) (VariablesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.DeleteAtManagementGroup(context.Context, string, string, string, *VariableValuesClientDeleteAtManagementGroupOptions) (VariableValuesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.GetAtManagementGroup(context.Context, string, string, string, *VariableValuesClientGetAtManagementGroupOptions) (VariableValuesClientGetAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Get(context.Context, string, string, *VariableValuesClientGetOptions) (VariableValuesClientGetResponse, error)`
- New function `*VariablesClient.Get(context.Context, string, *VariablesClientGetOptions) (VariablesClientGetResponse, error)`
- New function `NewVariableValuesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VariableValuesClient, error)`
- New function `*VariablesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, Variable, *VariablesClientCreateOrUpdateAtManagementGroupOptions) (VariablesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Delete(context.Context, string, string, *VariableValuesClientDeleteOptions) (VariableValuesClientDeleteResponse, error)`
- New function `*VariableValuesClient.CreateOrUpdate(context.Context, string, string, VariableValue, *VariableValuesClientCreateOrUpdateOptions) (VariableValuesClientCreateOrUpdateResponse, error)`
- New struct `Variable`
- New struct `VariableColumn`
- New struct `VariableListResult`
- New struct `VariableProperties`
- New struct `VariableValue`
- New struct `VariableValueColumnValue`
- New struct `VariableValueListResult`
- New struct `VariableValueProperties`
- New struct `VariableValuesClient`
- New struct `VariableValuesClientCreateOrUpdateAtManagementGroupOptions`
- New struct `VariableValuesClientCreateOrUpdateAtManagementGroupResponse`
- New struct `VariableValuesClientCreateOrUpdateOptions`
- New struct `VariableValuesClientCreateOrUpdateResponse`
- New struct `VariableValuesClientDeleteAtManagementGroupOptions`
- New struct `VariableValuesClientDeleteAtManagementGroupResponse`
- New struct `VariableValuesClientDeleteOptions`
- New struct `VariableValuesClientDeleteResponse`
- New struct `VariableValuesClientGetAtManagementGroupOptions`
- New struct `VariableValuesClientGetAtManagementGroupResponse`
- New struct `VariableValuesClientGetOptions`
- New struct `VariableValuesClientGetResponse`
- New struct `VariableValuesClientListForManagementGroupOptions`
- New struct `VariableValuesClientListForManagementGroupResponse`
- New struct `VariableValuesClientListOptions`
- New struct `VariableValuesClientListResponse`
- New struct `VariablesClient`
- New struct `VariablesClientCreateOrUpdateAtManagementGroupOptions`
- New struct `VariablesClientCreateOrUpdateAtManagementGroupResponse`
- New struct `VariablesClientCreateOrUpdateOptions`
- New struct `VariablesClientCreateOrUpdateResponse`
- New struct `VariablesClientDeleteAtManagementGroupOptions`
- New struct `VariablesClientDeleteAtManagementGroupResponse`
- New struct `VariablesClientDeleteOptions`
- New struct `VariablesClientDeleteResponse`
- New struct `VariablesClientGetAtManagementGroupOptions`
- New struct `VariablesClientGetAtManagementGroupResponse`
- New struct `VariablesClientGetOptions`
- New struct `VariablesClientGetResponse`
- New struct `VariablesClientListForManagementGroupOptions`
- New struct `VariablesClientListForManagementGroupResponse`
- New struct `VariablesClientListOptions`
- New struct `VariablesClientListResponse`


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).