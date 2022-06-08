# Release History

## 1.1.0 (2022-06-08)
### Features Added

- New function `*PrivateStoreCollectionClient.ApproveAllItems(context.Context, string, string, *PrivateStoreCollectionClientApproveAllItemsOptions) (PrivateStoreCollectionClientApproveAllItemsResponse, error)`
- New function `*PrivateStoreClient.AnyExistingOffersInTheCollections(context.Context, string, *PrivateStoreClientAnyExistingOffersInTheCollectionsOptions) (PrivateStoreClientAnyExistingOffersInTheCollectionsResponse, error)`
- New function `*PrivateStoreCollectionClient.DisableApproveAllItems(context.Context, string, string, *PrivateStoreCollectionClientDisableApproveAllItemsOptions) (PrivateStoreCollectionClientDisableApproveAllItemsResponse, error)`
- New function `*PrivateStoreCollectionOfferClient.UpsertOfferWithMultiContext(context.Context, string, string, string, *PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextOptions) (PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextResponse, error)`
- New function `*PrivateStoreClient.QueryUserOffers(context.Context, string, *PrivateStoreClientQueryUserOffersOptions) (PrivateStoreClientQueryUserOffersResponse, error)`
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
- New field `SubscriptionIDs` in struct `QueryApprovedPlans`
- New field `AllItemsApproved` in struct `CollectionProperties`
- New field `AllItemsApprovedModifiedAt` in struct `CollectionProperties`
- New field `Icon` in struct `AdminRequestApprovalProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplace/armmarketplace` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).