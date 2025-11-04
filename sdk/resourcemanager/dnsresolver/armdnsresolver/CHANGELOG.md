# Release History

## 2.0.0-beta.1 (2025-10-27)
### Breaking Changes

- Struct `DomainListResult` has been removed
- Field `DomainListResult` of struct `DomainListsClientListByResourceGroupResponse` has been removed
- Field `DomainListResult` of struct `DomainListsClientListResponse` has been removed

### Features Added

- New enum type `ManagedDomainList` with values `ManagedDomainListAzureDNSThreatIntel`
- New struct `DomainListListResult`
- New field `ManagedDomainLists` in struct `DNSSecurityRulePatchProperties`
- New field `ManagedDomainLists` in struct `DNSSecurityRuleProperties`
- New anonymous field `DomainListListResult` in struct `DomainListsClientListByResourceGroupResponse`
- New anonymous field `DomainListListResult` in struct `DomainListsClientListResponse`


## 1.3.0 (2025-06-12)
### Features Added

- New enum type `Action` with values `ActionDownload`, `ActionUpload`
- New enum type `ActionType` with values `ActionTypeAlert`, `ActionTypeAllow`, `ActionTypeBlock`
- New enum type `DNSSecurityRuleState` with values `DNSSecurityRuleStateDisabled`, `DNSSecurityRuleStateEnabled`
- New function `*ClientFactory.NewDNSSecurityRulesClient() *DNSSecurityRulesClient`
- New function `*ClientFactory.NewDomainListsClient() *DomainListsClient`
- New function `*ClientFactory.NewPoliciesClient() *PoliciesClient`
- New function `*ClientFactory.NewPolicyVirtualNetworkLinksClient() *PolicyVirtualNetworkLinksClient`
- New function `NewDNSSecurityRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DNSSecurityRulesClient, error)`
- New function `*DNSSecurityRulesClient.BeginCreateOrUpdate(context.Context, string, string, string, DNSSecurityRule, *DNSSecurityRulesClientBeginCreateOrUpdateOptions) (*runtime.Poller[DNSSecurityRulesClientCreateOrUpdateResponse], error)`
- New function `*DNSSecurityRulesClient.BeginDelete(context.Context, string, string, string, *DNSSecurityRulesClientBeginDeleteOptions) (*runtime.Poller[DNSSecurityRulesClientDeleteResponse], error)`
- New function `*DNSSecurityRulesClient.Get(context.Context, string, string, string, *DNSSecurityRulesClientGetOptions) (DNSSecurityRulesClientGetResponse, error)`
- New function `*DNSSecurityRulesClient.NewListPager(string, string, *DNSSecurityRulesClientListOptions) *runtime.Pager[DNSSecurityRulesClientListResponse]`
- New function `*DNSSecurityRulesClient.BeginUpdate(context.Context, string, string, string, DNSSecurityRulePatch, *DNSSecurityRulesClientBeginUpdateOptions) (*runtime.Poller[DNSSecurityRulesClientUpdateResponse], error)`
- New function `NewDomainListsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DomainListsClient, error)`
- New function `*DomainListsClient.BeginCreateOrUpdate(context.Context, string, string, DomainList, *DomainListsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DomainListsClientCreateOrUpdateResponse], error)`
- New function `*DomainListsClient.BeginDelete(context.Context, string, string, *DomainListsClientBeginDeleteOptions) (*runtime.Poller[DomainListsClientDeleteResponse], error)`
- New function `*DomainListsClient.Get(context.Context, string, string, *DomainListsClientGetOptions) (DomainListsClientGetResponse, error)`
- New function `*DomainListsClient.NewListByResourceGroupPager(string, *DomainListsClientListByResourceGroupOptions) *runtime.Pager[DomainListsClientListByResourceGroupResponse]`
- New function `*DomainListsClient.NewListPager(*DomainListsClientListOptions) *runtime.Pager[DomainListsClientListResponse]`
- New function `*DomainListsClient.BeginUpdate(context.Context, string, string, DomainListPatch, *DomainListsClientBeginUpdateOptions) (*runtime.Poller[DomainListsClientUpdateResponse], error)`
- New function `*DomainListsClient.BeginBulk(context.Context, string, string, DomainListBulk, *DomainListsClientBeginBulkOptions) (*runtime.Poller[DomainListsClientBulkResponse], error)`
- New function `NewPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PoliciesClient, error)`
- New function `*PoliciesClient.BeginCreateOrUpdate(context.Context, string, string, Policy, *PoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PoliciesClientCreateOrUpdateResponse], error)`
- New function `*PoliciesClient.BeginDelete(context.Context, string, string, *PoliciesClientBeginDeleteOptions) (*runtime.Poller[PoliciesClientDeleteResponse], error)`
- New function `*PoliciesClient.Get(context.Context, string, string, *PoliciesClientGetOptions) (PoliciesClientGetResponse, error)`
- New function `*PoliciesClient.NewListByResourceGroupPager(string, *PoliciesClientListByResourceGroupOptions) *runtime.Pager[PoliciesClientListByResourceGroupResponse]`
- New function `*PoliciesClient.NewListByVirtualNetworkPager(string, string, *PoliciesClientListByVirtualNetworkOptions) *runtime.Pager[PoliciesClientListByVirtualNetworkResponse]`
- New function `*PoliciesClient.NewListPager(*PoliciesClientListOptions) *runtime.Pager[PoliciesClientListResponse]`
- New function `*PoliciesClient.BeginUpdate(context.Context, string, string, PolicyPatch, *PoliciesClientBeginUpdateOptions) (*runtime.Poller[PoliciesClientUpdateResponse], error)`
- New function `NewPolicyVirtualNetworkLinksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PolicyVirtualNetworkLinksClient, error)`
- New function `*PolicyVirtualNetworkLinksClient.BeginCreateOrUpdate(context.Context, string, string, string, PolicyVirtualNetworkLink, *PolicyVirtualNetworkLinksClientBeginCreateOrUpdateOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientCreateOrUpdateResponse], error)`
- New function `*PolicyVirtualNetworkLinksClient.BeginDelete(context.Context, string, string, string, *PolicyVirtualNetworkLinksClientBeginDeleteOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientDeleteResponse], error)`
- New function `*PolicyVirtualNetworkLinksClient.Get(context.Context, string, string, string, *PolicyVirtualNetworkLinksClientGetOptions) (PolicyVirtualNetworkLinksClientGetResponse, error)`
- New function `*PolicyVirtualNetworkLinksClient.NewListPager(string, string, *PolicyVirtualNetworkLinksClientListOptions) *runtime.Pager[PolicyVirtualNetworkLinksClientListResponse]`
- New function `*PolicyVirtualNetworkLinksClient.BeginUpdate(context.Context, string, string, string, PolicyVirtualNetworkLinkPatch, *PolicyVirtualNetworkLinksClientBeginUpdateOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientUpdateResponse], error)`
- New function `PossibleActionValues() []Action`
- New struct `DNSSecurityRule`
- New struct `DNSSecurityRuleAction`
- New struct `DNSSecurityRuleListResult`
- New struct `DNSSecurityRulePatch`
- New struct `DNSSecurityRulePatchProperties`
- New struct `DNSSecurityRuleProperties`
- New struct `DomainList`
- New struct `DomainListBulk`
- New struct `DomainListBulkProperties`
- New struct `DomainListPatch`
- New struct `DomainListPatchProperties`
- New struct `DomainListProperties`
- New struct `DomainListResult`
- New struct `Policy`
- New struct `PolicyListResult`
- New struct `PolicyPatch`
- New struct `PolicyProperties`
- New struct `PolicyVirtualNetworkLink`
- New struct `PolicyVirtualNetworkLinkListResult`
- New struct `PolicyVirtualNetworkLinkPatch`
- New struct `PolicyVirtualNetworkLinkProperties`


## 1.3.0-beta.1 (2024-10-23)
### Features Added

- New enum type `ActionType` with values `ActionTypeAlert`, `ActionTypeAllow`, `ActionTypeBlock`
- New enum type `BlockResponseCode` with values `BlockResponseCodeSERVFAIL`
- New enum type `DNSSecurityRuleState` with values `DNSSecurityRuleStateDisabled`, `DNSSecurityRuleStateEnabled`
- New function `*ClientFactory.NewDNSSecurityRulesClient() *DNSSecurityRulesClient`
- New function `*ClientFactory.NewDomainListsClient() *DomainListsClient`
- New function `*ClientFactory.NewPoliciesClient() *PoliciesClient`
- New function `*ClientFactory.NewPolicyVirtualNetworkLinksClient() *PolicyVirtualNetworkLinksClient`
- New function `NewDNSSecurityRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DNSSecurityRulesClient, error)`
- New function `*DNSSecurityRulesClient.BeginCreateOrUpdate(context.Context, string, string, string, DNSSecurityRule, *DNSSecurityRulesClientBeginCreateOrUpdateOptions) (*runtime.Poller[DNSSecurityRulesClientCreateOrUpdateResponse], error)`
- New function `*DNSSecurityRulesClient.BeginDelete(context.Context, string, string, string, *DNSSecurityRulesClientBeginDeleteOptions) (*runtime.Poller[DNSSecurityRulesClientDeleteResponse], error)`
- New function `*DNSSecurityRulesClient.Get(context.Context, string, string, string, *DNSSecurityRulesClientGetOptions) (DNSSecurityRulesClientGetResponse, error)`
- New function `*DNSSecurityRulesClient.NewListPager(string, string, *DNSSecurityRulesClientListOptions) *runtime.Pager[DNSSecurityRulesClientListResponse]`
- New function `*DNSSecurityRulesClient.BeginUpdate(context.Context, string, string, string, DNSSecurityRulePatch, *DNSSecurityRulesClientBeginUpdateOptions) (*runtime.Poller[DNSSecurityRulesClientUpdateResponse], error)`
- New function `NewDomainListsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DomainListsClient, error)`
- New function `*DomainListsClient.BeginCreateOrUpdate(context.Context, string, string, DomainList, *DomainListsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DomainListsClientCreateOrUpdateResponse], error)`
- New function `*DomainListsClient.BeginDelete(context.Context, string, string, *DomainListsClientBeginDeleteOptions) (*runtime.Poller[DomainListsClientDeleteResponse], error)`
- New function `*DomainListsClient.Get(context.Context, string, string, *DomainListsClientGetOptions) (DomainListsClientGetResponse, error)`
- New function `*DomainListsClient.NewListByResourceGroupPager(string, *DomainListsClientListByResourceGroupOptions) *runtime.Pager[DomainListsClientListByResourceGroupResponse]`
- New function `*DomainListsClient.NewListPager(*DomainListsClientListOptions) *runtime.Pager[DomainListsClientListResponse]`
- New function `*DomainListsClient.BeginUpdate(context.Context, string, string, DomainListPatch, *DomainListsClientBeginUpdateOptions) (*runtime.Poller[DomainListsClientUpdateResponse], error)`
- New function `NewPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PoliciesClient, error)`
- New function `*PoliciesClient.BeginCreateOrUpdate(context.Context, string, string, Policy, *PoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PoliciesClientCreateOrUpdateResponse], error)`
- New function `*PoliciesClient.BeginDelete(context.Context, string, string, *PoliciesClientBeginDeleteOptions) (*runtime.Poller[PoliciesClientDeleteResponse], error)`
- New function `*PoliciesClient.Get(context.Context, string, string, *PoliciesClientGetOptions) (PoliciesClientGetResponse, error)`
- New function `*PoliciesClient.NewListByResourceGroupPager(string, *PoliciesClientListByResourceGroupOptions) *runtime.Pager[PoliciesClientListByResourceGroupResponse]`
- New function `*PoliciesClient.NewListByVirtualNetworkPager(string, string, *PoliciesClientListByVirtualNetworkOptions) *runtime.Pager[PoliciesClientListByVirtualNetworkResponse]`
- New function `*PoliciesClient.NewListPager(*PoliciesClientListOptions) *runtime.Pager[PoliciesClientListResponse]`
- New function `*PoliciesClient.BeginUpdate(context.Context, string, string, PolicyPatch, *PoliciesClientBeginUpdateOptions) (*runtime.Poller[PoliciesClientUpdateResponse], error)`
- New function `NewPolicyVirtualNetworkLinksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PolicyVirtualNetworkLinksClient, error)`
- New function `*PolicyVirtualNetworkLinksClient.BeginCreateOrUpdate(context.Context, string, string, string, PolicyVirtualNetworkLink, *PolicyVirtualNetworkLinksClientBeginCreateOrUpdateOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientCreateOrUpdateResponse], error)`
- New function `*PolicyVirtualNetworkLinksClient.BeginDelete(context.Context, string, string, string, *PolicyVirtualNetworkLinksClientBeginDeleteOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientDeleteResponse], error)`
- New function `*PolicyVirtualNetworkLinksClient.Get(context.Context, string, string, string, *PolicyVirtualNetworkLinksClientGetOptions) (PolicyVirtualNetworkLinksClientGetResponse, error)`
- New function `*PolicyVirtualNetworkLinksClient.NewListPager(string, string, *PolicyVirtualNetworkLinksClientListOptions) *runtime.Pager[PolicyVirtualNetworkLinksClientListResponse]`
- New function `*PolicyVirtualNetworkLinksClient.BeginUpdate(context.Context, string, string, string, PolicyVirtualNetworkLinkPatch, *PolicyVirtualNetworkLinksClientBeginUpdateOptions) (*runtime.Poller[PolicyVirtualNetworkLinksClientUpdateResponse], error)`
- New struct `DNSSecurityRule`
- New struct `DNSSecurityRuleAction`
- New struct `DNSSecurityRuleListResult`
- New struct `DNSSecurityRulePatch`
- New struct `DNSSecurityRulePatchProperties`
- New struct `DNSSecurityRuleProperties`
- New struct `DomainList`
- New struct `DomainListPatch`
- New struct `DomainListPatchProperties`
- New struct `DomainListProperties`
- New struct `DomainListResult`
- New struct `Policy`
- New struct `PolicyListResult`
- New struct `PolicyPatch`
- New struct `PolicyProperties`
- New struct `PolicyVirtualNetworkLink`
- New struct `PolicyVirtualNetworkLinkListResult`
- New struct `PolicyVirtualNetworkLinkPatch`
- New struct `PolicyVirtualNetworkLinkProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-09-15)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `TrackedResource` has been removed

### Features Added

- New field `DNSResolverOutboundEndpoints` in struct `DNSForwardingRulesetPatch`


## 0.4.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dnsresolver/armdnsresolver` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.4.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).