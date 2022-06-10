# Unreleased

## Breaking Changes

### Removed Funcs

1. *Operation.UnmarshalJSON([]byte) error
1. *Subnet.UnmarshalJSON([]byte) error
1. NewSubnetsClient(string) SubnetsClient
1. NewSubnetsClientWithBaseURI(string, string) SubnetsClient
1. PossibleStatusValues() []Status
1. Subnet.MarshalJSON() ([]byte, error)
1. SubnetsClient.CreateOrUpdate(context.Context, string, string, string) (Subnet, error)
1. SubnetsClient.CreateOrUpdatePreparer(context.Context, string, string, string) (*http.Request, error)
1. SubnetsClient.CreateOrUpdateResponder(*http.Response) (Subnet, error)
1. SubnetsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. SubnetsClient.Get(context.Context, string, string, string) (Subnet, error)
1. SubnetsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. SubnetsClient.GetResponder(*http.Response) (Subnet, error)
1. SubnetsClient.GetSender(*http.Request) (*http.Response, error)
1. SubnetsClient.ListByEnterprisePolicy(context.Context, string, string) (SubnetListResult, error)
1. SubnetsClient.ListByEnterprisePolicyPreparer(context.Context, string, string) (*http.Request, error)
1. SubnetsClient.ListByEnterprisePolicyResponder(*http.Response) (SubnetListResult, error)
1. SubnetsClient.ListByEnterprisePolicySender(*http.Request) (*http.Response, error)

### Struct Changes

#### Removed Structs

1. ErrorResponseBody
1. OperationList
1. OperationProperties
1. Subnet
1. SubnetEndpointProperty
1. SubnetListResult
1. SubnetsClient

#### Removed Struct Fields

1. KeyVaultProperties.Status
1. Operation.*OperationProperties
1. PrivateEndpointConnection.Location
1. PrivateLinkResource.Location
1. PrivateLinkResourceProperties.SystemData
1. Properties.SystemData
1. PropertiesEncryption.KeyVaultProperties
1. PropertiesLockbox.Status
1. ProxyResource.Location
1. Resource.Location
1. SubnetProperties.Status
1. SubnetProperties.Subnet
1. SubnetProperties.SystemData

### Signature Changes

#### Const Types

1. Disabled changed type from Status to State
1. Enabled changed type from Status to State
1. NotConfigured changed type from Status to State

#### Funcs

1. EnterprisePoliciesClient.ListByResourceGroup
	- Returns
		- From: EnterprisePolicyList, error
		- To: EnterprisePolicyListPage, error
1. EnterprisePoliciesClient.ListBySubscription
	- Returns
		- From: EnterprisePolicyList, error
		- To: EnterprisePolicyListPage, error
1. EnterprisePoliciesClient.Update
	- Params
		- From: context.Context, string, string, EnterprisePolicy
		- To: context.Context, string, string, PatchEnterprisePolicy
1. EnterprisePoliciesClient.UpdatePreparer
	- Params
		- From: context.Context, string, string, EnterprisePolicy
		- To: context.Context, string, string, PatchEnterprisePolicy
1. OperationsClient.List
	- Returns
		- From: OperationList, error
		- To: OperationListResultPage, error
1. OperationsClient.ListResponder
	- Returns
		- From: OperationList, error
		- To: OperationListResult, error

#### Struct Fields

1. ErrorResponse.Error changed type from *ErrorResponseBody to *ErrorDetail
1. Operation.IsDataAction changed type from *string to *bool

## Additive Changes

### New Constants

1. ActionType.Internal
1. EnterprisePolicyKind.EnterprisePolicyKindEncryption
1. EnterprisePolicyKind.EnterprisePolicyKindLockbox
1. EnterprisePolicyKind.EnterprisePolicyKindNetworkInjection
1. EnterprisePolicyKind.EnterprisePolicyKindPrivateEndpoint
1. Origin.OriginSystem
1. Origin.OriginUser
1. Origin.OriginUsersystem

### New Funcs

1. *Account.UnmarshalJSON([]byte) error
1. *AccountListIterator.Next() error
1. *AccountListIterator.NextWithContext(context.Context) error
1. *AccountListPage.Next() error
1. *AccountListPage.NextWithContext(context.Context) error
1. *EnterprisePolicyListIterator.Next() error
1. *EnterprisePolicyListIterator.NextWithContext(context.Context) error
1. *EnterprisePolicyListPage.Next() error
1. *EnterprisePolicyListPage.NextWithContext(context.Context) error
1. *OperationListResultIterator.Next() error
1. *OperationListResultIterator.NextWithContext(context.Context) error
1. *OperationListResultPage.Next() error
1. *OperationListResultPage.NextWithContext(context.Context) error
1. *PatchAccount.UnmarshalJSON([]byte) error
1. *PatchEnterprisePolicy.UnmarshalJSON([]byte) error
1. Account.MarshalJSON() ([]byte, error)
1. AccountList.IsEmpty() bool
1. AccountListIterator.NotDone() bool
1. AccountListIterator.Response() AccountList
1. AccountListIterator.Value() Account
1. AccountListPage.NotDone() bool
1. AccountListPage.Response() AccountList
1. AccountListPage.Values() []Account
1. AccountsClient.CreateOrUpdate(context.Context, string, string, Account) (Account, error)
1. AccountsClient.CreateOrUpdatePreparer(context.Context, string, string, Account) (*http.Request, error)
1. AccountsClient.CreateOrUpdateResponder(*http.Response) (Account, error)
1. AccountsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. AccountsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. AccountsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. AccountsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. AccountsClient.DeleteSender(*http.Request) (*http.Response, error)
1. AccountsClient.Get(context.Context, string, string) (Account, error)
1. AccountsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. AccountsClient.GetResponder(*http.Response) (Account, error)
1. AccountsClient.GetSender(*http.Request) (*http.Response, error)
1. AccountsClient.ListByResourceGroup(context.Context, string) (AccountListPage, error)
1. AccountsClient.ListByResourceGroupComplete(context.Context, string) (AccountListIterator, error)
1. AccountsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. AccountsClient.ListByResourceGroupResponder(*http.Response) (AccountList, error)
1. AccountsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. AccountsClient.ListBySubscription(context.Context) (AccountListPage, error)
1. AccountsClient.ListBySubscriptionComplete(context.Context) (AccountListIterator, error)
1. AccountsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. AccountsClient.ListBySubscriptionResponder(*http.Response) (AccountList, error)
1. AccountsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. AccountsClient.Update(context.Context, string, string, PatchAccount) (Account, error)
1. AccountsClient.UpdatePreparer(context.Context, string, string, PatchAccount) (*http.Request, error)
1. AccountsClient.UpdateResponder(*http.Response) (Account, error)
1. AccountsClient.UpdateSender(*http.Request) (*http.Response, error)
1. AzureEntityResource.MarshalJSON() ([]byte, error)
1. EnterprisePoliciesClient.ListByResourceGroupComplete(context.Context, string) (EnterprisePolicyListIterator, error)
1. EnterprisePoliciesClient.ListBySubscriptionComplete(context.Context) (EnterprisePolicyListIterator, error)
1. EnterprisePolicyList.IsEmpty() bool
1. EnterprisePolicyListIterator.NotDone() bool
1. EnterprisePolicyListIterator.Response() EnterprisePolicyList
1. EnterprisePolicyListIterator.Value() EnterprisePolicy
1. EnterprisePolicyListPage.NotDone() bool
1. EnterprisePolicyListPage.Response() EnterprisePolicyList
1. EnterprisePolicyListPage.Values() []EnterprisePolicy
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. NewAccountListIterator(AccountListPage) AccountListIterator
1. NewAccountListPage(AccountList, func(context.Context, AccountList) (AccountList, error)) AccountListPage
1. NewAccountsClient(string) AccountsClient
1. NewAccountsClientWithBaseURI(string, string) AccountsClient
1. NewEnterprisePolicyListIterator(EnterprisePolicyListPage) EnterprisePolicyListIterator
1. NewEnterprisePolicyListPage(EnterprisePolicyList, func(context.Context, EnterprisePolicyList) (EnterprisePolicyList, error)) EnterprisePolicyListPage
1. NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator
1. NewOperationListResultPage(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage
1. OperationDisplay.MarshalJSON() ([]byte, error)
1. OperationListResult.IsEmpty() bool
1. OperationListResult.MarshalJSON() ([]byte, error)
1. OperationListResultIterator.NotDone() bool
1. OperationListResultIterator.Response() OperationListResult
1. OperationListResultIterator.Value() Operation
1. OperationListResultPage.NotDone() bool
1. OperationListResultPage.Response() OperationListResult
1. OperationListResultPage.Values() []Operation
1. OperationsClient.ListComplete(context.Context) (OperationListResultIterator, error)
1. PatchAccount.MarshalJSON() ([]byte, error)
1. PatchEnterprisePolicy.MarshalJSON() ([]byte, error)
1. PatchTrackedResource.MarshalJSON() ([]byte, error)
1. PossibleActionTypeValues() []ActionType
1. PossibleEnterprisePolicyKindValues() []EnterprisePolicyKind
1. PossibleOriginValues() []Origin
1. PossibleStateValues() []State

### Struct Changes

#### New Structs

1. Account
1. AccountList
1. AccountListIterator
1. AccountListPage
1. AccountProperties
1. AccountsClient
1. AzureEntityResource
1. EnterprisePolicyListIterator
1. EnterprisePolicyListPage
1. ErrorAdditionalInfo
1. ErrorDetail
1. OperationListResult
1. OperationListResultIterator
1. OperationListResultPage
1. PatchAccount
1. PatchEnterprisePolicy
1. PatchTrackedResource
1. PropertiesNetworkInjection
1. VirtualNetworkProperties
1. VirtualNetworkPropertiesList

#### New Struct Fields

1. EnterprisePolicy.Kind
1. EnterprisePolicy.SystemData
1. EnterprisePolicyList.NextLink
1. Operation.ActionType
1. Operation.Origin
1. PrivateEndpointConnection.SystemData
1. PrivateLinkResourceProperties.RequiredZoneNames
1. Properties.NetworkInjection
1. PropertiesEncryption.KeyVault
1. PropertiesEncryption.State
1. PropertiesLockbox.State
1. SubnetProperties.Name
