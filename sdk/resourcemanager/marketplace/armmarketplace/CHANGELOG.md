# Release History

## 2.0.0 (2026-01-21)
### Breaking Changes

- Function `*PrivateStoreClient.CreateApprovalRequest` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, requestApprovalID string, options *PrivateStoreClientCreateApprovalRequestOptions)` to `(ctx context.Context, privateStoreID string, requestApprovalID string, payload RequestApprovalResource, options *PrivateStoreClientCreateApprovalRequestOptions)`
- Function `*PrivateStoreClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, options *PrivateStoreClientCreateOrUpdateOptions)` to `(ctx context.Context, privateStoreID string, payload PrivateStore, options *PrivateStoreClientCreateOrUpdateOptions)`
- Function `*PrivateStoreClient.UpdateAdminRequestApproval` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, adminRequestApprovalID string, options *PrivateStoreClientUpdateAdminRequestApprovalOptions)` to `(ctx context.Context, privateStoreID string, adminRequestApprovalID string, payload AdminRequestApprovalsResource, options *PrivateStoreClientUpdateAdminRequestApprovalOptions)`
- Function `*PrivateStoreCollectionOfferClient.CreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, offerID string, collectionID string, options *PrivateStoreCollectionOfferClientCreateOrUpdateOptions)` to `(ctx context.Context, privateStoreID string, collectionID string, offerID string, payload Offer, options *PrivateStoreCollectionOfferClientCreateOrUpdateOptions)`
- Function `*PrivateStoreCollectionOfferClient.Delete` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, offerID string, collectionID string, options *PrivateStoreCollectionOfferClientDeleteOptions)` to `(ctx context.Context, privateStoreID string, collectionID string, offerID string, options *PrivateStoreCollectionOfferClientDeleteOptions)`
- Function `*PrivateStoreCollectionOfferClient.Get` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, offerID string, collectionID string, options *PrivateStoreCollectionOfferClientGetOptions)` to `(ctx context.Context, privateStoreID string, collectionID string, offerID string, options *PrivateStoreCollectionOfferClientGetOptions)`
- Function `*PrivateStoreCollectionOfferClient.Post` parameter(s) have been changed from `(ctx context.Context, privateStoreID string, offerID string, collectionID string, options *PrivateStoreCollectionOfferClientPostOptions)` to `(ctx context.Context, privateStoreID string, collectionID string, offerID string, options *PrivateStoreCollectionOfferClientPostOptions)`
- Type of `CollectionsToSubscriptionsMappingResponse.Details` has been changed from `map[string]*CollectionsSubscriptionsMappingDetails` to `*CollectionsToSubscriptionsMappingResponseProperties`
- Type of `PlanDetails.RequestDate` has been changed from `any` to `*time.Time`
- Type of `SystemData.CreatedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Type of `SystemData.LastModifiedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Enum `IdentityType` has been removed
- Operation `*OperationsClient.List` has supported pagination, use `*OperationsClient.NewListPager` instead.
- Operation `*PrivateStoreClient.AdminRequestApprovalsList` has supported pagination, use `*PrivateStoreClient.NewAdminRequestApprovalsListPager` instead.
- Operation `*PrivateStoreClient.GetApprovalRequestsList` has supported pagination, use `*PrivateStoreClient.NewGetApprovalRequestsListPager` instead.
- Operation `*PrivateStoreCollectionClient.List` has supported pagination, use `*PrivateStoreCollectionClient.NewListPager` instead.
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `Resource` has been removed
- Field `Payload` of struct `PrivateStoreClientCreateApprovalRequestOptions` has been removed
- Field `Payload` of struct `PrivateStoreClientCreateOrUpdateOptions` has been removed
- Field `Payload` of struct `PrivateStoreClientUpdateAdminRequestApprovalOptions` has been removed
- Field `Payload` of struct `PrivateStoreCollectionClientCreateOrUpdateOptions` has been removed
- Field `Payload` of struct `PrivateStoreCollectionOfferClientCreateOrUpdateOptions` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `RuleType` with values `RuleTypePrivateProducts`, `RuleTypeTermsAndCondition`
- New function `NewClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.NewOperationsClient() *OperationsClient`
- New function `*Client.NewPrivateStoreClient() *PrivateStoreClient`
- New function `*Client.NewPrivateStoreCollectionClient() *PrivateStoreCollectionClient`
- New function `*Client.NewPrivateStoreCollectionOfferClient() *PrivateStoreCollectionOfferClient`
- New function `*Client.QueryRules(ctx context.Context, privateStoreID string, collectionID string, options *ClientQueryRulesOptions) (ClientQueryRulesResponse, error)`
- New function `*Client.QueryUserRules(ctx context.Context, privateStoreID string, options *ClientQueryUserRulesOptions) (ClientQueryUserRulesResponse, error)`
- New function `*Client.SetCollectionRules(ctx context.Context, privateStoreID string, collectionID string, options *ClientSetCollectionRulesOptions) (ClientSetCollectionRulesResponse, error)`
- New function `*ClientFactory.NewClient() *Client`
- New function `*PrivateStoreClient.AnyExistingOffersInTheCollections(ctx context.Context, privateStoreID string, options *PrivateStoreClientAnyExistingOffersInTheCollectionsOptions) (PrivateStoreClientAnyExistingOffersInTheCollectionsResponse, error)`
- New function `*PrivateStoreClient.QueryUserOffers(ctx context.Context, privateStoreID string, options *PrivateStoreClientQueryUserOffersOptions) (PrivateStoreClientQueryUserOffersResponse, error)`
- New function `*PrivateStoreCollectionClient.ApproveAllItems(ctx context.Context, privateStoreID string, collectionID string, options *PrivateStoreCollectionClientApproveAllItemsOptions) (PrivateStoreCollectionClientApproveAllItemsResponse, error)`
- New function `*PrivateStoreCollectionClient.DisableApproveAllItems(ctx context.Context, privateStoreID string, collectionID string, options *PrivateStoreCollectionClientDisableApproveAllItemsOptions) (PrivateStoreCollectionClientDisableApproveAllItemsResponse, error)`
- New function `*PrivateStoreCollectionOfferClient.ContextsView(ctx context.Context, privateStoreID string, collectionID string, offerID string, options *PrivateStoreCollectionOfferClientContextsViewOptions) (PrivateStoreCollectionOfferClientContextsViewResponse, error)`
- New function `*PrivateStoreCollectionOfferClient.NewListByContextsPager(privateStoreID string, collectionID string, options *PrivateStoreCollectionOfferClientListByContextsOptions) *runtime.Pager[PrivateStoreCollectionOfferClientListByContextsResponse]`
- New function `*PrivateStoreCollectionOfferClient.UpsertOfferWithMultiContext(ctx context.Context, privateStoreID string, collectionID string, offerID string, options *PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextOptions) (PrivateStoreCollectionOfferClientUpsertOfferWithMultiContextResponse, error)`
- New struct `AnyExistingOffersInTheCollectionsResponse`
- New struct `CollectionOffersByAllContextsPayload`
- New struct `CollectionOffersByAllContextsProperties`
- New struct `CollectionOffersByContext`
- New struct `CollectionOffersByContextList`
- New struct `CollectionOffersByContextOffers`
- New struct `CollectionsToSubscriptionsMappingResponseProperties`
- New struct `ContextAndPlansDetails`
- New struct `MultiContextAndPlansPayload`
- New struct `MultiContextAndPlansProperties`
- New struct `QueryUserOffersDetails`
- New struct `QueryUserOffersProperties`
- New struct `QueryUserRulesDetails`
- New struct `QueryUserRulesProperties`
- New struct `Rule`
- New struct `RuleListResponse`
- New struct `SetRulesRequest`
- New field `Icon` in struct `AdminRequestApprovalProperties`
- New field `AppliedRules`, `ApproveAllItems`, `ApproveAllItemsModifiedAt` in struct `CollectionProperties`
- New field `IsStopSell` in struct `OfferProperties`
- New field `IsStopSell` in struct `Plan`
- New field `SubscriptionIDs` in struct `QueryApprovedPlans`
- New field `ID` in struct `SingleOperation`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplace/armmarketplace` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).