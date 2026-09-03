# Release History

## 1.0.0 (2026-09-03)
### Breaking Changes

- Field `TargetTenant` of struct `ServiceGroupMemberRelationshipProperties` has been removed

### Features Added

- New function `*ClientFactory.NewContainsRelationshipsClient(subscriptionID string) *ContainsRelationshipsClient`
- New function `*ClientFactory.NewDependencyOfRelationshipsByServiceGroupClient() *DependencyOfRelationshipsByServiceGroupClient`
- New function `NewContainsRelationshipsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContainsRelationshipsClient, error)`
- New function `*ContainsRelationshipsClient.NewListByResourceGroupPager(resourceGroupName string, options *ContainsRelationshipsClientListByResourceGroupOptions) *runtime.Pager[ContainsRelationshipsClientListByResourceGroupResponse]`
- New function `*ContainsRelationshipsClient.NewListBySubscriptionPager(options *ContainsRelationshipsClientListBySubscriptionOptions) *runtime.Pager[ContainsRelationshipsClientListBySubscriptionResponse]`
- New function `NewDependencyOfRelationshipsByServiceGroupClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*DependencyOfRelationshipsByServiceGroupClient, error)`
- New function `*DependencyOfRelationshipsByServiceGroupClient.BeginCreateOrUpdate(ctx context.Context, serviceGroupName string, name string, resource DependencyOfRelationship, options *DependencyOfRelationshipsByServiceGroupClientBeginCreateOrUpdateOptions) (*runtime.Poller[DependencyOfRelationshipsByServiceGroupClientCreateOrUpdateResponse], error)`
- New function `*DependencyOfRelationshipsByServiceGroupClient.BeginDelete(ctx context.Context, serviceGroupName string, name string, options *DependencyOfRelationshipsByServiceGroupClientBeginDeleteOptions) (*runtime.Poller[DependencyOfRelationshipsByServiceGroupClientDeleteResponse], error)`
- New function `*DependencyOfRelationshipsByServiceGroupClient.Get(ctx context.Context, serviceGroupName string, name string, options *DependencyOfRelationshipsByServiceGroupClientGetOptions) (DependencyOfRelationshipsByServiceGroupClientGetResponse, error)`
- New function `*DependencyOfRelationshipsByServiceGroupClient.NewListPager(serviceGroupName string, options *DependencyOfRelationshipsByServiceGroupClientListOptions) *runtime.Pager[DependencyOfRelationshipsByServiceGroupClientListResponse]`
- New function `*DependencyOfRelationshipsClient.NewListByParentPager(resourceURI string, options *DependencyOfRelationshipsClientListByParentOptions) *runtime.Pager[DependencyOfRelationshipsClientListByParentResponse]`
- New function `*ServiceGroupMemberRelationshipsClient.NewListByParentPager(resourceURI string, options *ServiceGroupMemberRelationshipsClientListByParentOptions) *runtime.Pager[ServiceGroupMemberRelationshipsClientListByParentResponse]`
- New struct `ContainsRelationship`
- New struct `ContainsRelationshipListResult`
- New struct `ContainsRelationshipProperties`
- New struct `DependencyOfRelationshipListResult`
- New struct `ServiceGroupMemberRelationshipListResult`
- New field `SourceTenant` in struct `ServiceGroupMemberRelationshipProperties`


## 0.2.0 (2026-08-13)
### Breaking Changes

- Field `TargetTenant` of struct `ServiceGroupMemberRelationshipProperties` has been removed

### Features Added

- New function `*ClientFactory.NewContainsRelationshipsClient(subscriptionID string) *ContainsRelationshipsClient`
- New function `NewContainsRelationshipsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContainsRelationshipsClient, error)`
- New function `*ContainsRelationshipsClient.NewListByResourceGroupPager(resourceGroupName string, options *ContainsRelationshipsClientListByResourceGroupOptions) *runtime.Pager[ContainsRelationshipsClientListByResourceGroupResponse]`
- New function `*ContainsRelationshipsClient.NewListBySubscriptionPager(options *ContainsRelationshipsClientListBySubscriptionOptions) *runtime.Pager[ContainsRelationshipsClientListBySubscriptionResponse]`
- New function `*DependencyOfRelationshipsClient.NewListByParentPager(resourceURI string, options *DependencyOfRelationshipsClientListByParentOptions) *runtime.Pager[DependencyOfRelationshipsClientListByParentResponse]`
- New function `*ServiceGroupMemberRelationshipsClient.NewListByParentPager(resourceURI string, options *ServiceGroupMemberRelationshipsClientListByParentOptions) *runtime.Pager[ServiceGroupMemberRelationshipsClientListByParentResponse]`
- New struct `ContainsRelationship`
- New struct `ContainsRelationshipListResult`
- New struct `ContainsRelationshipProperties`
- New struct `DependencyOfRelationshipListResult`
- New struct `ServiceGroupMemberRelationshipListResult`
- New field `SourceTenant` in struct `ServiceGroupMemberRelationshipProperties`


## 0.1.0 (2026-04-07)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/relationships/armrelationships` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).