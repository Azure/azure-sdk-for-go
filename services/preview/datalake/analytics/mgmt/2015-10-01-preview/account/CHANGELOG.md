# Unreleased

## Breaking Changes

### Removed Constants

1. OperationStatus.OperationStatusFailed
1. OperationStatus.OperationStatusInProgress
1. OperationStatus.OperationStatusSucceeded

### Removed Funcs

1. *DataLakeAnalyticsAccountListDataLakeStoreResultIterator.Next() error
1. *DataLakeAnalyticsAccountListDataLakeStoreResultIterator.NextWithContext(context.Context) error
1. *DataLakeAnalyticsAccountListDataLakeStoreResultPage.Next() error
1. *DataLakeAnalyticsAccountListDataLakeStoreResultPage.NextWithContext(context.Context) error
1. *DataLakeAnalyticsAccountListStorageAccountsResultIterator.Next() error
1. *DataLakeAnalyticsAccountListStorageAccountsResultIterator.NextWithContext(context.Context) error
1. *DataLakeAnalyticsAccountListStorageAccountsResultPage.Next() error
1. *DataLakeAnalyticsAccountListStorageAccountsResultPage.NextWithContext(context.Context) error
1. *ListBlobContainersResultIterator.Next() error
1. *ListBlobContainersResultIterator.NextWithContext(context.Context) error
1. *ListBlobContainersResultPage.Next() error
1. *ListBlobContainersResultPage.NextWithContext(context.Context) error
1. *ListSasTokensResultIterator.Next() error
1. *ListSasTokensResultIterator.NextWithContext(context.Context) error
1. *ListSasTokensResultPage.Next() error
1. *ListSasTokensResultPage.NextWithContext(context.Context) error
1. AzureAsyncOperationResult.MarshalJSON() ([]byte, error)
1. BlobContainer.MarshalJSON() ([]byte, error)
1. BlobContainerProperties.MarshalJSON() ([]byte, error)
1. DataLakeAnalyticsAccountListDataLakeStoreResult.IsEmpty() bool
1. DataLakeAnalyticsAccountListDataLakeStoreResult.MarshalJSON() ([]byte, error)
1. DataLakeAnalyticsAccountListDataLakeStoreResultIterator.NotDone() bool
1. DataLakeAnalyticsAccountListDataLakeStoreResultIterator.Response() DataLakeAnalyticsAccountListDataLakeStoreResult
1. DataLakeAnalyticsAccountListDataLakeStoreResultIterator.Value() DataLakeStoreAccountInfo
1. DataLakeAnalyticsAccountListDataLakeStoreResultPage.NotDone() bool
1. DataLakeAnalyticsAccountListDataLakeStoreResultPage.Response() DataLakeAnalyticsAccountListDataLakeStoreResult
1. DataLakeAnalyticsAccountListDataLakeStoreResultPage.Values() []DataLakeStoreAccountInfo
1. DataLakeAnalyticsAccountListStorageAccountsResult.IsEmpty() bool
1. DataLakeAnalyticsAccountListStorageAccountsResult.MarshalJSON() ([]byte, error)
1. DataLakeAnalyticsAccountListStorageAccountsResultIterator.NotDone() bool
1. DataLakeAnalyticsAccountListStorageAccountsResultIterator.Response() DataLakeAnalyticsAccountListStorageAccountsResult
1. DataLakeAnalyticsAccountListStorageAccountsResultIterator.Value() StorageAccountInfo
1. DataLakeAnalyticsAccountListStorageAccountsResultPage.NotDone() bool
1. DataLakeAnalyticsAccountListStorageAccountsResultPage.Response() DataLakeAnalyticsAccountListStorageAccountsResult
1. DataLakeAnalyticsAccountListStorageAccountsResultPage.Values() []StorageAccountInfo
1. Error.MarshalJSON() ([]byte, error)
1. ErrorDetails.MarshalJSON() ([]byte, error)
1. InnerError.MarshalJSON() ([]byte, error)
1. ListBlobContainersResult.IsEmpty() bool
1. ListBlobContainersResult.MarshalJSON() ([]byte, error)
1. ListBlobContainersResultIterator.NotDone() bool
1. ListBlobContainersResultIterator.Response() ListBlobContainersResult
1. ListBlobContainersResultIterator.Value() BlobContainer
1. ListBlobContainersResultPage.NotDone() bool
1. ListBlobContainersResultPage.Response() ListBlobContainersResult
1. ListBlobContainersResultPage.Values() []BlobContainer
1. ListSasTokensResult.IsEmpty() bool
1. ListSasTokensResult.MarshalJSON() ([]byte, error)
1. ListSasTokensResultIterator.NotDone() bool
1. ListSasTokensResultIterator.Response() ListSasTokensResult
1. ListSasTokensResultIterator.Value() SasTokenInfo
1. ListSasTokensResultPage.NotDone() bool
1. ListSasTokensResultPage.Response() ListSasTokensResult
1. ListSasTokensResultPage.Values() []SasTokenInfo
1. NewDataLakeAnalyticsAccountListDataLakeStoreResultIterator(DataLakeAnalyticsAccountListDataLakeStoreResultPage) DataLakeAnalyticsAccountListDataLakeStoreResultIterator
1. NewDataLakeAnalyticsAccountListDataLakeStoreResultPage(DataLakeAnalyticsAccountListDataLakeStoreResult, func(context.Context, DataLakeAnalyticsAccountListDataLakeStoreResult) (DataLakeAnalyticsAccountListDataLakeStoreResult, error)) DataLakeAnalyticsAccountListDataLakeStoreResultPage
1. NewDataLakeAnalyticsAccountListStorageAccountsResultIterator(DataLakeAnalyticsAccountListStorageAccountsResultPage) DataLakeAnalyticsAccountListStorageAccountsResultIterator
1. NewDataLakeAnalyticsAccountListStorageAccountsResultPage(DataLakeAnalyticsAccountListStorageAccountsResult, func(context.Context, DataLakeAnalyticsAccountListStorageAccountsResult) (DataLakeAnalyticsAccountListStorageAccountsResult, error)) DataLakeAnalyticsAccountListStorageAccountsResultPage
1. NewListBlobContainersResultIterator(ListBlobContainersResultPage) ListBlobContainersResultIterator
1. NewListBlobContainersResultPage(ListBlobContainersResult, func(context.Context, ListBlobContainersResult) (ListBlobContainersResult, error)) ListBlobContainersResultPage
1. NewListSasTokensResultIterator(ListSasTokensResultPage) ListSasTokensResultIterator
1. NewListSasTokensResultPage(ListSasTokensResult, func(context.Context, ListSasTokensResult) (ListSasTokensResult, error)) ListSasTokensResultPage
1. PossibleOperationStatusValues() []OperationStatus
1. SasTokenInfo.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AzureAsyncOperationResult
1. BlobContainer
1. BlobContainerProperties
1. DataLakeAnalyticsAccountListDataLakeStoreResult
1. DataLakeAnalyticsAccountListDataLakeStoreResultIterator
1. DataLakeAnalyticsAccountListDataLakeStoreResultPage
1. DataLakeAnalyticsAccountListStorageAccountsResult
1. DataLakeAnalyticsAccountListStorageAccountsResultIterator
1. DataLakeAnalyticsAccountListStorageAccountsResultPage
1. DataLakeStoreAccountInfo
1. Error
1. ErrorDetails
1. InnerError
1. ListBlobContainersResult
1. ListBlobContainersResultIterator
1. ListBlobContainersResultPage
1. ListSasTokensResult
1. ListSasTokensResultIterator
1. ListSasTokensResultPage
1. SasTokenInfo
1. StorageAccountInfo

#### Removed Struct Fields

1. AddDataLakeStoreParameters.Properties
1. AddStorageAccountParameters.Properties
1. DataLakeAnalyticsAccount.Properties

### Signature Changes

#### Funcs

1. Client.AddDataLakeStoreAccount
	- Params
		- From: context.Context, string, string, string, AddDataLakeStoreParameters
		- To: context.Context, string, string, string, *AddDataLakeStoreParameters
1. Client.AddDataLakeStoreAccountPreparer
	- Params
		- From: context.Context, string, string, string, AddDataLakeStoreParameters
		- To: context.Context, string, string, string, *AddDataLakeStoreParameters
1. Client.Create
	- Params
		- From: context.Context, string, string, DataLakeAnalyticsAccount
		- To: context.Context, string, string, CreateDataLakeAnalyticsAccountParameters
1. Client.CreatePreparer
	- Params
		- From: context.Context, string, string, DataLakeAnalyticsAccount
		- To: context.Context, string, string, CreateDataLakeAnalyticsAccountParameters
1. Client.GetDataLakeStoreAccount
	- Returns
		- From: DataLakeStoreAccountInfo, error
		- To: DataLakeStoreAccountInformation, error
1. Client.GetDataLakeStoreAccountResponder
	- Returns
		- From: DataLakeStoreAccountInfo, error
		- To: DataLakeStoreAccountInformation, error
1. Client.GetStorageAccount
	- Returns
		- From: StorageAccountInfo, error
		- To: StorageAccountInformation, error
1. Client.GetStorageAccountResponder
	- Returns
		- From: StorageAccountInfo, error
		- To: StorageAccountInformation, error
1. Client.GetStorageContainer
	- Returns
		- From: BlobContainer, error
		- To: StorageContainer, error
1. Client.GetStorageContainerResponder
	- Returns
		- From: BlobContainer, error
		- To: StorageContainer, error
1. Client.List
	- Params
		- From: context.Context, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, *int32, *int32, string, string, *bool
1. Client.ListByResourceGroup
	- Params
		- From: context.Context, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, *int32, *int32, string, string, *bool
1. Client.ListByResourceGroupComplete
	- Params
		- From: context.Context, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, *int32, *int32, string, string, *bool
1. Client.ListByResourceGroupPreparer
	- Params
		- From: context.Context, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, *int32, *int32, string, string, *bool
1. Client.ListComplete
	- Params
		- From: context.Context, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, *int32, *int32, string, string, *bool
1. Client.ListDataLakeStoreAccounts
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
	- Returns
		- From: DataLakeAnalyticsAccountListDataLakeStoreResultPage, error
		- To: DataLakeStoreAccountInformationListResultPage, error
1. Client.ListDataLakeStoreAccountsComplete
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
	- Returns
		- From: DataLakeAnalyticsAccountListDataLakeStoreResultIterator, error
		- To: DataLakeStoreAccountInformationListResultIterator, error
1. Client.ListDataLakeStoreAccountsPreparer
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
1. Client.ListDataLakeStoreAccountsResponder
	- Returns
		- From: DataLakeAnalyticsAccountListDataLakeStoreResult, error
		- To: DataLakeStoreAccountInformationListResult, error
1. Client.ListPreparer
	- Params
		- From: context.Context, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, *int32, *int32, string, string, *bool
1. Client.ListSasTokens
	- Returns
		- From: ListSasTokensResultPage, error
		- To: SasTokenInformationListResultPage, error
1. Client.ListSasTokensComplete
	- Returns
		- From: ListSasTokensResultIterator, error
		- To: SasTokenInformationListResultIterator, error
1. Client.ListSasTokensResponder
	- Returns
		- From: ListSasTokensResult, error
		- To: SasTokenInformationListResult, error
1. Client.ListStorageAccounts
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
	- Returns
		- From: DataLakeAnalyticsAccountListStorageAccountsResultPage, error
		- To: StorageAccountInformationListResultPage, error
1. Client.ListStorageAccountsComplete
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
	- Returns
		- From: DataLakeAnalyticsAccountListStorageAccountsResultIterator, error
		- To: StorageAccountInformationListResultIterator, error
1. Client.ListStorageAccountsPreparer
	- Params
		- From: context.Context, string, string, string, *int32, *int32, string, string, string, *bool, string, string
		- To: context.Context, string, string, string, *int32, *int32, string, string, *bool
1. Client.ListStorageAccountsResponder
	- Returns
		- From: DataLakeAnalyticsAccountListStorageAccountsResult, error
		- To: StorageAccountInformationListResult, error
1. Client.ListStorageContainers
	- Returns
		- From: ListBlobContainersResultPage, error
		- To: StorageContainerListResultPage, error
1. Client.ListStorageContainersComplete
	- Returns
		- From: ListBlobContainersResultIterator, error
		- To: StorageContainerListResultIterator, error
1. Client.ListStorageContainersResponder
	- Returns
		- From: ListBlobContainersResult, error
		- To: StorageContainerListResult, error
1. Client.Update
	- Params
		- From: context.Context, string, string, DataLakeAnalyticsAccount
		- To: context.Context, string, string, *UpdateDataLakeAnalyticsAccountParameters
1. Client.UpdatePreparer
	- Params
		- From: context.Context, string, string, DataLakeAnalyticsAccount
		- To: context.Context, string, string, *UpdateDataLakeAnalyticsAccountParameters
1. Client.UpdateStorageAccount
	- Params
		- From: context.Context, string, string, string, AddStorageAccountParameters
		- To: context.Context, string, string, string, *UpdateStorageAccountParameters
1. Client.UpdateStorageAccountPreparer
	- Params
		- From: context.Context, string, string, string, AddStorageAccountParameters
		- To: context.Context, string, string, string, *UpdateStorageAccountParameters

#### Struct Fields

1. DataLakeAnalyticsAccountProperties.DataLakeStoreAccounts changed type from *[]DataLakeStoreAccountInfo to *[]DataLakeStoreAccountInformation
1. DataLakeAnalyticsAccountProperties.StorageAccounts changed type from *[]StorageAccountInfo to *[]StorageAccountInformation

## Additive Changes

### New Constants

1. AADObjectType.Group
1. AADObjectType.ServicePrincipal
1. AADObjectType.User
1. DebugDataAccessLevel.All
1. DebugDataAccessLevel.Customer
1. DebugDataAccessLevel.None
1. FirewallAllowAzureIpsState.Disabled
1. FirewallAllowAzureIpsState.Enabled
1. FirewallState.FirewallStateDisabled
1. FirewallState.FirewallStateEnabled
1. NestedResourceProvisioningState.NestedResourceProvisioningStateCanceled
1. NestedResourceProvisioningState.NestedResourceProvisioningStateFailed
1. NestedResourceProvisioningState.NestedResourceProvisioningStateSucceeded
1. OperationOrigin.OperationOriginSystem
1. OperationOrigin.OperationOriginUser
1. OperationOrigin.OperationOriginUsersystem
1. SubscriptionState.SubscriptionStateDeleted
1. SubscriptionState.SubscriptionStateRegistered
1. SubscriptionState.SubscriptionStateSuspended
1. SubscriptionState.SubscriptionStateUnregistered
1. SubscriptionState.SubscriptionStateWarned
1. TierType.Commitment100000AUHours
1. TierType.Commitment10000AUHours
1. TierType.Commitment1000AUHours
1. TierType.Commitment100AUHours
1. TierType.Commitment500000AUHours
1. TierType.Commitment50000AUHours
1. TierType.Commitment5000AUHours
1. TierType.Commitment500AUHours
1. TierType.Consumption
1. VirtualNetworkRuleState.VirtualNetworkRuleStateActive
1. VirtualNetworkRuleState.VirtualNetworkRuleStateFailed
1. VirtualNetworkRuleState.VirtualNetworkRuleStateNetworkSourceDeleted

### New Funcs

1. *AddDataLakeStoreParameters.UnmarshalJSON([]byte) error
1. *AddDataLakeStoreWithAccountParameters.UnmarshalJSON([]byte) error
1. *AddStorageAccountParameters.UnmarshalJSON([]byte) error
1. *AddStorageAccountWithAccountParameters.UnmarshalJSON([]byte) error
1. *ComputePolicy.UnmarshalJSON([]byte) error
1. *ComputePolicyListResultIterator.Next() error
1. *ComputePolicyListResultIterator.NextWithContext(context.Context) error
1. *ComputePolicyListResultPage.Next() error
1. *ComputePolicyListResultPage.NextWithContext(context.Context) error
1. *CreateComputePolicyWithAccountParameters.UnmarshalJSON([]byte) error
1. *CreateDataLakeAnalyticsAccountParameters.UnmarshalJSON([]byte) error
1. *CreateFirewallRuleWithAccountParameters.UnmarshalJSON([]byte) error
1. *CreateOrUpdateComputePolicyParameters.UnmarshalJSON([]byte) error
1. *CreateOrUpdateFirewallRuleParameters.UnmarshalJSON([]byte) error
1. *DataLakeAnalyticsAccount.UnmarshalJSON([]byte) error
1. *DataLakeAnalyticsAccountBasic.UnmarshalJSON([]byte) error
1. *DataLakeStoreAccountInformation.UnmarshalJSON([]byte) error
1. *DataLakeStoreAccountInformationListResultIterator.Next() error
1. *DataLakeStoreAccountInformationListResultIterator.NextWithContext(context.Context) error
1. *DataLakeStoreAccountInformationListResultPage.Next() error
1. *DataLakeStoreAccountInformationListResultPage.NextWithContext(context.Context) error
1. *FirewallRule.UnmarshalJSON([]byte) error
1. *FirewallRuleListResultIterator.Next() error
1. *FirewallRuleListResultIterator.NextWithContext(context.Context) error
1. *FirewallRuleListResultPage.Next() error
1. *FirewallRuleListResultPage.NextWithContext(context.Context) error
1. *HiveMetastore.UnmarshalJSON([]byte) error
1. *SasTokenInformationListResultIterator.Next() error
1. *SasTokenInformationListResultIterator.NextWithContext(context.Context) error
1. *SasTokenInformationListResultPage.Next() error
1. *SasTokenInformationListResultPage.NextWithContext(context.Context) error
1. *StorageAccountInformation.UnmarshalJSON([]byte) error
1. *StorageAccountInformationListResultIterator.Next() error
1. *StorageAccountInformationListResultIterator.NextWithContext(context.Context) error
1. *StorageAccountInformationListResultPage.Next() error
1. *StorageAccountInformationListResultPage.NextWithContext(context.Context) error
1. *StorageContainer.UnmarshalJSON([]byte) error
1. *StorageContainerListResultIterator.Next() error
1. *StorageContainerListResultIterator.NextWithContext(context.Context) error
1. *StorageContainerListResultPage.Next() error
1. *StorageContainerListResultPage.NextWithContext(context.Context) error
1. *UpdateComputePolicyParameters.UnmarshalJSON([]byte) error
1. *UpdateComputePolicyWithAccountParameters.UnmarshalJSON([]byte) error
1. *UpdateDataLakeAnalyticsAccountParameters.UnmarshalJSON([]byte) error
1. *UpdateDataLakeStoreWithAccountParameters.UnmarshalJSON([]byte) error
1. *UpdateFirewallRuleParameters.UnmarshalJSON([]byte) error
1. *UpdateFirewallRuleWithAccountParameters.UnmarshalJSON([]byte) error
1. *UpdateStorageAccountParameters.UnmarshalJSON([]byte) error
1. *UpdateStorageAccountWithAccountParameters.UnmarshalJSON([]byte) error
1. *VirtualNetworkRule.UnmarshalJSON([]byte) error
1. AccountsClient.CheckNameAvailability(context.Context, string, CheckNameAvailabilityParameters) (NameAvailabilityInformation, error)
1. AccountsClient.CheckNameAvailabilityPreparer(context.Context, string, CheckNameAvailabilityParameters) (*http.Request, error)
1. AccountsClient.CheckNameAvailabilityResponder(*http.Response) (NameAvailabilityInformation, error)
1. AccountsClient.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)
1. AddDataLakeStoreParameters.MarshalJSON() ([]byte, error)
1. AddDataLakeStoreWithAccountParameters.MarshalJSON() ([]byte, error)
1. AddStorageAccountParameters.MarshalJSON() ([]byte, error)
1. AddStorageAccountWithAccountParameters.MarshalJSON() ([]byte, error)
1. CapabilityInformation.MarshalJSON() ([]byte, error)
1. ComputePoliciesClient.CreateOrUpdate(context.Context, string, string, string, CreateOrUpdateComputePolicyParameters) (ComputePolicy, error)
1. ComputePoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, string, CreateOrUpdateComputePolicyParameters) (*http.Request, error)
1. ComputePoliciesClient.CreateOrUpdateResponder(*http.Response) (ComputePolicy, error)
1. ComputePoliciesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ComputePoliciesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ComputePoliciesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ComputePoliciesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ComputePoliciesClient.DeleteSender(*http.Request) (*http.Response, error)
1. ComputePoliciesClient.Get(context.Context, string, string, string) (ComputePolicy, error)
1. ComputePoliciesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ComputePoliciesClient.GetResponder(*http.Response) (ComputePolicy, error)
1. ComputePoliciesClient.GetSender(*http.Request) (*http.Response, error)
1. ComputePoliciesClient.ListByAccount(context.Context, string, string) (ComputePolicyListResultPage, error)
1. ComputePoliciesClient.ListByAccountComplete(context.Context, string, string) (ComputePolicyListResultIterator, error)
1. ComputePoliciesClient.ListByAccountPreparer(context.Context, string, string) (*http.Request, error)
1. ComputePoliciesClient.ListByAccountResponder(*http.Response) (ComputePolicyListResult, error)
1. ComputePoliciesClient.ListByAccountSender(*http.Request) (*http.Response, error)
1. ComputePoliciesClient.Update(context.Context, string, string, string, *UpdateComputePolicyParameters) (ComputePolicy, error)
1. ComputePoliciesClient.UpdatePreparer(context.Context, string, string, string, *UpdateComputePolicyParameters) (*http.Request, error)
1. ComputePoliciesClient.UpdateResponder(*http.Response) (ComputePolicy, error)
1. ComputePoliciesClient.UpdateSender(*http.Request) (*http.Response, error)
1. ComputePolicy.MarshalJSON() ([]byte, error)
1. ComputePolicyListResult.IsEmpty() bool
1. ComputePolicyListResult.MarshalJSON() ([]byte, error)
1. ComputePolicyListResultIterator.NotDone() bool
1. ComputePolicyListResultIterator.Response() ComputePolicyListResult
1. ComputePolicyListResultIterator.Value() ComputePolicy
1. ComputePolicyListResultPage.NotDone() bool
1. ComputePolicyListResultPage.Response() ComputePolicyListResult
1. ComputePolicyListResultPage.Values() []ComputePolicy
1. ComputePolicyProperties.MarshalJSON() ([]byte, error)
1. CreateComputePolicyWithAccountParameters.MarshalJSON() ([]byte, error)
1. CreateDataLakeAnalyticsAccountParameters.MarshalJSON() ([]byte, error)
1. CreateFirewallRuleWithAccountParameters.MarshalJSON() ([]byte, error)
1. CreateOrUpdateComputePolicyParameters.MarshalJSON() ([]byte, error)
1. CreateOrUpdateFirewallRuleParameters.MarshalJSON() ([]byte, error)
1. DataLakeAnalyticsAccountBasic.MarshalJSON() ([]byte, error)
1. DataLakeAnalyticsAccountPropertiesBasic.MarshalJSON() ([]byte, error)
1. DataLakeStoreAccountInfoProperties.MarshalJSON() ([]byte, error)
1. DataLakeStoreAccountInformation.MarshalJSON() ([]byte, error)
1. DataLakeStoreAccountInformationListResult.IsEmpty() bool
1. DataLakeStoreAccountInformationListResult.MarshalJSON() ([]byte, error)
1. DataLakeStoreAccountInformationListResultIterator.NotDone() bool
1. DataLakeStoreAccountInformationListResultIterator.Response() DataLakeStoreAccountInformationListResult
1. DataLakeStoreAccountInformationListResultIterator.Value() DataLakeStoreAccountInformation
1. DataLakeStoreAccountInformationListResultPage.NotDone() bool
1. DataLakeStoreAccountInformationListResultPage.Response() DataLakeStoreAccountInformationListResult
1. DataLakeStoreAccountInformationListResultPage.Values() []DataLakeStoreAccountInformation
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. FirewallRule.MarshalJSON() ([]byte, error)
1. FirewallRuleListResult.IsEmpty() bool
1. FirewallRuleListResult.MarshalJSON() ([]byte, error)
1. FirewallRuleListResultIterator.NotDone() bool
1. FirewallRuleListResultIterator.Response() FirewallRuleListResult
1. FirewallRuleListResultIterator.Value() FirewallRule
1. FirewallRuleListResultPage.NotDone() bool
1. FirewallRuleListResultPage.Response() FirewallRuleListResult
1. FirewallRuleListResultPage.Values() []FirewallRule
1. FirewallRuleProperties.MarshalJSON() ([]byte, error)
1. FirewallRulesClient.CreateOrUpdate(context.Context, string, string, string, CreateOrUpdateFirewallRuleParameters) (FirewallRule, error)
1. FirewallRulesClient.CreateOrUpdatePreparer(context.Context, string, string, string, CreateOrUpdateFirewallRuleParameters) (*http.Request, error)
1. FirewallRulesClient.CreateOrUpdateResponder(*http.Response) (FirewallRule, error)
1. FirewallRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. FirewallRulesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. FirewallRulesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. FirewallRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. FirewallRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. FirewallRulesClient.Get(context.Context, string, string, string) (FirewallRule, error)
1. FirewallRulesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. FirewallRulesClient.GetResponder(*http.Response) (FirewallRule, error)
1. FirewallRulesClient.GetSender(*http.Request) (*http.Response, error)
1. FirewallRulesClient.ListByAccount(context.Context, string, string) (FirewallRuleListResultPage, error)
1. FirewallRulesClient.ListByAccountComplete(context.Context, string, string) (FirewallRuleListResultIterator, error)
1. FirewallRulesClient.ListByAccountPreparer(context.Context, string, string) (*http.Request, error)
1. FirewallRulesClient.ListByAccountResponder(*http.Response) (FirewallRuleListResult, error)
1. FirewallRulesClient.ListByAccountSender(*http.Request) (*http.Response, error)
1. FirewallRulesClient.Update(context.Context, string, string, string, *UpdateFirewallRuleParameters) (FirewallRule, error)
1. FirewallRulesClient.UpdatePreparer(context.Context, string, string, string, *UpdateFirewallRuleParameters) (*http.Request, error)
1. FirewallRulesClient.UpdateResponder(*http.Response) (FirewallRule, error)
1. FirewallRulesClient.UpdateSender(*http.Request) (*http.Response, error)
1. HiveMetastore.MarshalJSON() ([]byte, error)
1. HiveMetastoreListResult.MarshalJSON() ([]byte, error)
1. HiveMetastoreProperties.MarshalJSON() ([]byte, error)
1. LocationsClient.GetCapability(context.Context, string) (CapabilityInformation, error)
1. LocationsClient.GetCapabilityPreparer(context.Context, string) (*http.Request, error)
1. LocationsClient.GetCapabilityResponder(*http.Response) (CapabilityInformation, error)
1. LocationsClient.GetCapabilitySender(*http.Request) (*http.Response, error)
1. NameAvailabilityInformation.MarshalJSON() ([]byte, error)
1. NewAccountsClient(string) AccountsClient
1. NewAccountsClientWithBaseURI(string, string) AccountsClient
1. NewComputePoliciesClient(string) ComputePoliciesClient
1. NewComputePoliciesClientWithBaseURI(string, string) ComputePoliciesClient
1. NewComputePolicyListResultIterator(ComputePolicyListResultPage) ComputePolicyListResultIterator
1. NewComputePolicyListResultPage(ComputePolicyListResult, func(context.Context, ComputePolicyListResult) (ComputePolicyListResult, error)) ComputePolicyListResultPage
1. NewDataLakeStoreAccountInformationListResultIterator(DataLakeStoreAccountInformationListResultPage) DataLakeStoreAccountInformationListResultIterator
1. NewDataLakeStoreAccountInformationListResultPage(DataLakeStoreAccountInformationListResult, func(context.Context, DataLakeStoreAccountInformationListResult) (DataLakeStoreAccountInformationListResult, error)) DataLakeStoreAccountInformationListResultPage
1. NewFirewallRuleListResultIterator(FirewallRuleListResultPage) FirewallRuleListResultIterator
1. NewFirewallRuleListResultPage(FirewallRuleListResult, func(context.Context, FirewallRuleListResult) (FirewallRuleListResult, error)) FirewallRuleListResultPage
1. NewFirewallRulesClient(string) FirewallRulesClient
1. NewFirewallRulesClientWithBaseURI(string, string) FirewallRulesClient
1. NewLocationsClient(string) LocationsClient
1. NewLocationsClientWithBaseURI(string, string) LocationsClient
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. NewSasTokenInformationListResultIterator(SasTokenInformationListResultPage) SasTokenInformationListResultIterator
1. NewSasTokenInformationListResultPage(SasTokenInformationListResult, func(context.Context, SasTokenInformationListResult) (SasTokenInformationListResult, error)) SasTokenInformationListResultPage
1. NewStorageAccountInformationListResultIterator(StorageAccountInformationListResultPage) StorageAccountInformationListResultIterator
1. NewStorageAccountInformationListResultPage(StorageAccountInformationListResult, func(context.Context, StorageAccountInformationListResult) (StorageAccountInformationListResult, error)) StorageAccountInformationListResultPage
1. NewStorageContainerListResultIterator(StorageContainerListResultPage) StorageContainerListResultIterator
1. NewStorageContainerListResultPage(StorageContainerListResult, func(context.Context, StorageContainerListResult) (StorageContainerListResult, error)) StorageContainerListResultPage
1. Operation.MarshalJSON() ([]byte, error)
1. OperationDisplay.MarshalJSON() ([]byte, error)
1. OperationListResult.MarshalJSON() ([]byte, error)
1. OperationsClient.List(context.Context) (OperationListResult, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (OperationListResult, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PossibleAADObjectTypeValues() []AADObjectType
1. PossibleDebugDataAccessLevelValues() []DebugDataAccessLevel
1. PossibleFirewallAllowAzureIpsStateValues() []FirewallAllowAzureIpsState
1. PossibleFirewallStateValues() []FirewallState
1. PossibleNestedResourceProvisioningStateValues() []NestedResourceProvisioningState
1. PossibleOperationOriginValues() []OperationOrigin
1. PossibleSubscriptionStateValues() []SubscriptionState
1. PossibleTierTypeValues() []TierType
1. PossibleVirtualNetworkRuleStateValues() []VirtualNetworkRuleState
1. Resource.MarshalJSON() ([]byte, error)
1. SasTokenInformation.MarshalJSON() ([]byte, error)
1. SasTokenInformationListResult.IsEmpty() bool
1. SasTokenInformationListResult.MarshalJSON() ([]byte, error)
1. SasTokenInformationListResultIterator.NotDone() bool
1. SasTokenInformationListResultIterator.Response() SasTokenInformationListResult
1. SasTokenInformationListResultIterator.Value() SasTokenInformation
1. SasTokenInformationListResultPage.NotDone() bool
1. SasTokenInformationListResultPage.Response() SasTokenInformationListResult
1. SasTokenInformationListResultPage.Values() []SasTokenInformation
1. StorageAccountInformation.MarshalJSON() ([]byte, error)
1. StorageAccountInformationListResult.IsEmpty() bool
1. StorageAccountInformationListResult.MarshalJSON() ([]byte, error)
1. StorageAccountInformationListResultIterator.NotDone() bool
1. StorageAccountInformationListResultIterator.Response() StorageAccountInformationListResult
1. StorageAccountInformationListResultIterator.Value() StorageAccountInformation
1. StorageAccountInformationListResultPage.NotDone() bool
1. StorageAccountInformationListResultPage.Response() StorageAccountInformationListResult
1. StorageAccountInformationListResultPage.Values() []StorageAccountInformation
1. StorageAccountInformationProperties.MarshalJSON() ([]byte, error)
1. StorageContainer.MarshalJSON() ([]byte, error)
1. StorageContainerListResult.IsEmpty() bool
1. StorageContainerListResult.MarshalJSON() ([]byte, error)
1. StorageContainerListResultIterator.NotDone() bool
1. StorageContainerListResultIterator.Response() StorageContainerListResult
1. StorageContainerListResultIterator.Value() StorageContainer
1. StorageContainerListResultPage.NotDone() bool
1. StorageContainerListResultPage.Response() StorageContainerListResult
1. StorageContainerListResultPage.Values() []StorageContainer
1. StorageContainerProperties.MarshalJSON() ([]byte, error)
1. SubResource.MarshalJSON() ([]byte, error)
1. UpdateComputePolicyParameters.MarshalJSON() ([]byte, error)
1. UpdateComputePolicyWithAccountParameters.MarshalJSON() ([]byte, error)
1. UpdateDataLakeAnalyticsAccountParameters.MarshalJSON() ([]byte, error)
1. UpdateDataLakeStoreWithAccountParameters.MarshalJSON() ([]byte, error)
1. UpdateFirewallRuleParameters.MarshalJSON() ([]byte, error)
1. UpdateFirewallRuleWithAccountParameters.MarshalJSON() ([]byte, error)
1. UpdateStorageAccountParameters.MarshalJSON() ([]byte, error)
1. UpdateStorageAccountWithAccountParameters.MarshalJSON() ([]byte, error)
1. VirtualNetworkRule.MarshalJSON() ([]byte, error)
1. VirtualNetworkRuleListResult.MarshalJSON() ([]byte, error)
1. VirtualNetworkRuleProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AccountsClient
1. AddDataLakeStoreProperties
1. AddDataLakeStoreWithAccountParameters
1. AddStorageAccountWithAccountParameters
1. CapabilityInformation
1. CheckNameAvailabilityParameters
1. ComputePoliciesClient
1. ComputePolicy
1. ComputePolicyListResult
1. ComputePolicyListResultIterator
1. ComputePolicyListResultPage
1. ComputePolicyProperties
1. CreateComputePolicyWithAccountParameters
1. CreateDataLakeAnalyticsAccountParameters
1. CreateDataLakeAnalyticsAccountProperties
1. CreateFirewallRuleWithAccountParameters
1. CreateOrUpdateComputePolicyParameters
1. CreateOrUpdateComputePolicyProperties
1. CreateOrUpdateFirewallRuleParameters
1. CreateOrUpdateFirewallRuleProperties
1. DataLakeAnalyticsAccountBasic
1. DataLakeAnalyticsAccountPropertiesBasic
1. DataLakeStoreAccountInformation
1. DataLakeStoreAccountInformationListResult
1. DataLakeStoreAccountInformationListResultIterator
1. DataLakeStoreAccountInformationListResultPage
1. ErrorAdditionalInfo
1. ErrorDetail
1. ErrorResponse
1. FirewallRule
1. FirewallRuleListResult
1. FirewallRuleListResultIterator
1. FirewallRuleListResultPage
1. FirewallRuleProperties
1. FirewallRulesClient
1. HiveMetastore
1. HiveMetastoreListResult
1. HiveMetastoreProperties
1. LocationsClient
1. NameAvailabilityInformation
1. Operation
1. OperationDisplay
1. OperationListResult
1. OperationMetaLogSpecification
1. OperationMetaMetricAvailabilitiesSpecification
1. OperationMetaMetricSpecification
1. OperationMetaPropertyInfo
1. OperationMetaServiceSpecification
1. OperationsClient
1. Resource
1. SasTokenInformation
1. SasTokenInformationListResult
1. SasTokenInformationListResultIterator
1. SasTokenInformationListResultPage
1. StorageAccountInformation
1. StorageAccountInformationListResult
1. StorageAccountInformationListResultIterator
1. StorageAccountInformationListResultPage
1. StorageAccountInformationProperties
1. StorageContainer
1. StorageContainerListResult
1. StorageContainerListResultIterator
1. StorageContainerListResultPage
1. StorageContainerProperties
1. SubResource
1. UpdateComputePolicyParameters
1. UpdateComputePolicyProperties
1. UpdateComputePolicyWithAccountParameters
1. UpdateDataLakeAnalyticsAccountParameters
1. UpdateDataLakeAnalyticsAccountProperties
1. UpdateDataLakeStoreProperties
1. UpdateDataLakeStoreWithAccountParameters
1. UpdateFirewallRuleParameters
1. UpdateFirewallRuleProperties
1. UpdateFirewallRuleWithAccountParameters
1. UpdateStorageAccountParameters
1. UpdateStorageAccountProperties
1. UpdateStorageAccountWithAccountParameters
1. VirtualNetworkRule
1. VirtualNetworkRuleListResult
1. VirtualNetworkRuleProperties

#### New Struct Fields

1. AddDataLakeStoreParameters.*AddDataLakeStoreProperties
1. AddStorageAccountParameters.*StorageAccountProperties
1. DataLakeAnalyticsAccount.*DataLakeAnalyticsAccountProperties
1. DataLakeAnalyticsAccountProperties.AccountID
1. DataLakeAnalyticsAccountProperties.ComputePolicies
1. DataLakeAnalyticsAccountProperties.CurrentTier
1. DataLakeAnalyticsAccountProperties.DebugDataAccessLevel
1. DataLakeAnalyticsAccountProperties.FirewallAllowAzureIps
1. DataLakeAnalyticsAccountProperties.FirewallRules
1. DataLakeAnalyticsAccountProperties.FirewallState
1. DataLakeAnalyticsAccountProperties.HierarchicalQueueState
1. DataLakeAnalyticsAccountProperties.HiveMetastores
1. DataLakeAnalyticsAccountProperties.MaxDegreeOfParallelismPerJob
1. DataLakeAnalyticsAccountProperties.MinPriorityPerJob
1. DataLakeAnalyticsAccountProperties.NewTier
1. DataLakeAnalyticsAccountProperties.PublicDataLakeStoreAccounts
1. DataLakeAnalyticsAccountProperties.QueryStoreRetention
1. DataLakeAnalyticsAccountProperties.SystemMaxDegreeOfParallelism
1. DataLakeAnalyticsAccountProperties.SystemMaxJobCount
1. DataLakeAnalyticsAccountProperties.VirtualNetworkRules
