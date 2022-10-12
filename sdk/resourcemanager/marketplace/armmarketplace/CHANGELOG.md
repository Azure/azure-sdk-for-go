# Release History

## 1.1.0 (2022-10-12)
### Features Added

- New const `RuleTypeTermsAndCondition`
- New const `RuleTypePrivateProducts`
- New type alias `RuleType`
- New function `PossibleRuleTypeValues() []RuleType`
- New function `*RPServiceClient.QueryUserRules(context.Context, string, *RPServiceClientQueryUserRulesOptions) (RPServiceClientQueryUserRulesResponse, error)`
- New function `*PrivateStoreCollectionOfferClient.UpsertOfferWithMultiContext(context.Context, string, string, string, *PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextOptions) (PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextResponse, error)`
- New function `NewRPServiceClient(azcore.TokenCredential, *arm.ClientOptions) (*RPServiceClient, error)`
- New function `*PrivateStoreClient.AnyExistingOffersInTheCollections(context.Context, string, *PrivateStoreClientAnyExistingOffersInTheCollectionsOptions) (PrivateStoreClientAnyExistingOffersInTheCollectionsResponse, error)`
- New function `*PrivateStoreCollectionClient.DisableApproveAllItems(context.Context, string, string, *PrivateStoreCollectionClientDisableApproveAllItemsOptions) (PrivateStoreCollectionClientDisableApproveAllItemsResponse, error)`
- New function `*RPServiceClient.QueryRules(context.Context, string, string, *RPServiceClientQueryRulesOptions) (RPServiceClientQueryRulesResponse, error)`
- New function `*PrivateStoreClient.QueryUserOffers(context.Context, string, *PrivateStoreClientQueryUserOffersOptions) (PrivateStoreClientQueryUserOffersResponse, error)`
- New function `*PrivateStoreCollectionClient.ApproveAllItems(context.Context, string, string, *PrivateStoreCollectionClientApproveAllItemsOptions) (PrivateStoreCollectionClientApproveAllItemsResponse, error)`
- New function `*RPServiceClient.SetCollectionRules(context.Context, string, string, *RPServiceClientSetCollectionRulesOptions) (RPServiceClientSetCollectionRulesResponse, error)`
- New struct `AnyExistingOffersInTheCollectionsResponse`
- New struct `ContextAndPlansDetails`
- New struct `MultiContextAndPlansPayload`
- New struct `MultiContextAndPlansProperties`
- New struct `PrivateStoreClientAnyExistingOffersInTheCollectionsOptions`
- New struct `PrivateStoreClientAnyExistingOffersInTheCollectionsResponse`
- New struct `PrivateStoreClientQueryUserOffersOptions`
- New struct `PrivateStoreClientQueryUserOffersResponse`
- New struct `PrivateStoreCollectionClientApproveAllItemsOptions`
- New struct `PrivateStoreCollectionClientApproveAllItemsResponse`
- New struct `PrivateStoreCollectionClientDisableApproveAllItemsOptions`
- New struct `PrivateStoreCollectionClientDisableApproveAllItemsResponse`
- New struct `PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextOptions`
- New struct `PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextResponse`
- New struct `QueryUserOffersDetails`
- New struct `QueryUserOffersProperties`
- New struct `QueryUserRulesDetails`
- New struct `QueryUserRulesProperties`
- New struct `RPServiceClient`
- New struct `RPServiceClientQueryRulesOptions`
- New struct `RPServiceClientQueryRulesResponse`
- New struct `RPServiceClientQueryUserRulesOptions`
- New struct `RPServiceClientQueryUserRulesResponse`
- New struct `RPServiceClientSetCollectionRulesOptions`
- New struct `RPServiceClientSetCollectionRulesResponse`
- New struct `Rule`
- New struct `RuleListResponse`
- New struct `SetRulesRequest`
- New field `SubscriptionIDs` in struct `QueryApprovedPlans`
- New field `Icon` in struct `AdminRequestApprovalProperties`
- New field `AppliedRules` in struct `CollectionProperties`
- New field `ApproveAllItems` in struct `CollectionProperties`
- New field `ApproveAllItemsModifiedAt` in struct `CollectionProperties`
- New field `ID` in struct `SingleOperation`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplace/armmarketplace` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).