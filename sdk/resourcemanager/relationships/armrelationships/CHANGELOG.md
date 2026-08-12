# Release History

## 0.2.0 (2026-08-12)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Type of `ServiceGroupMemberRelationship.Properties` has been changed from `*ServiceGroupMemberRelationshipProperties` to `*ServiceGroupMemberRelationshipPropertiesV2`
- Struct `ServiceGroupMemberRelationshipProperties` has been removed

### Features Added

- New function `*ClientFactory.NewContainsRelationshipsClient() *ContainsRelationshipsClient`
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
- New struct `ServiceGroupMemberRelationshipPropertiesV2`


## 0.1.0 (2026-04-07)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/relationships/armrelationships` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).