# Release History

## 3.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Type of `PurchaseRequest.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Type of `ReservationOrderAliasRequest.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Type of `ReservationOrderAliasResponse.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Type of `SavingsPlanModel.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Type of `SavingsPlanOrderAliasModel.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Type of `SavingsPlanOrderModel.SKU` has been changed from `*SKU` to `*ResourceSKU`
- Operation `*SavingsPlanClient.Update` has been changed to LRO, use `*SavingsPlanClient.BeginUpdate` instead.

### Features Added

- New value `CommitmentGrainFullTerm`, `CommitmentGrainUnknown` added to enum type `CommitmentGrain`
- New enum type `ApplyDiscountOn` with values `ApplyDiscountOnConsume`, `ApplyDiscountOnPurchase`, `ApplyDiscountOnRenew`
- New enum type `DiscountAppliedScopeType` with values `DiscountAppliedScopeTypeBillingAccount`, `DiscountAppliedScopeTypeBillingProfile`, `DiscountAppliedScopeTypeCustomer`
- New enum type `DiscountCombinationRule` with values `DiscountCombinationRuleBestOf`, `DiscountCombinationRuleStackable`
- New enum type `DiscountEntityType` with values `DiscountEntityTypeAffiliate`, `DiscountEntityTypePrimary`
- New enum type `DiscountProvisioningState` with values `DiscountProvisioningStateCanceled`, `DiscountProvisioningStateFailed`, `DiscountProvisioningStatePending`, `DiscountProvisioningStateSucceeded`, `DiscountProvisioningStateUnknown`
- New enum type `DiscountRuleType` with values `DiscountRuleTypeFixedListPrice`, `DiscountRuleTypeFixedPriceLock`, `DiscountRuleTypePriceCeiling`
- New enum type `DiscountStatus` with values `DiscountStatusActive`, `DiscountStatusCanceled`, `DiscountStatusExpired`, `DiscountStatusFailed`, `DiscountStatusPending`
- New enum type `DiscountType` with values `DiscountTypeCustomPrice`, `DiscountTypeCustomPriceMultiCurrency`, `DiscountTypeProduct`, `DiscountTypeProductFamily`, `DiscountTypeSKU`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `PricingPolicy` with values `PricingPolicyLocked`, `PricingPolicyProtected`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New function `*ClientFactory.NewDiscountClient() *DiscountClient`
- New function `*ClientFactory.NewDiscountsClient() *DiscountsClient`
- New function `NewDiscountClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DiscountClient, error)`
- New function `*DiscountClient.Get(ctx context.Context, resourceGroupName string, discountName string, options *DiscountClientGetOptions) (DiscountClientGetResponse, error)`
- New function `*DiscountClient.BeginUpdate(ctx context.Context, resourceGroupName string, discountName string, body DiscountPatchRequest, options *DiscountClientBeginUpdateOptions) (*runtime.Poller[DiscountClientUpdateResponse], error)`
- New function `*DiscountProperties.GetDiscountProperties() *DiscountProperties`
- New function `*DiscountTypeCustomPrice.GetDiscountTypeCustomPrice() *DiscountTypeCustomPrice`
- New function `*DiscountTypeCustomPrice.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `*DiscountTypeCustomPriceMultiCurrency.GetDiscountTypeCustomPrice() *DiscountTypeCustomPrice`
- New function `*DiscountTypeCustomPriceMultiCurrency.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `*DiscountTypeProduct.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `*DiscountTypeProductFamily.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `*DiscountTypeProductSKU.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `*DiscountTypeProperties.GetDiscountTypeProperties() *DiscountTypeProperties`
- New function `NewDiscountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DiscountsClient, error)`
- New function `*DiscountsClient.BeginCancel(ctx context.Context, resourceGroupName string, discountName string, options *DiscountsClientBeginCancelOptions) (*runtime.Poller[DiscountsClientCancelResponse], error)`
- New function `*DiscountsClient.BeginCreate(ctx context.Context, resourceGroupName string, discountName string, body Discount, options *DiscountsClientBeginCreateOptions) (*runtime.Poller[DiscountsClientCreateResponse], error)`
- New function `*DiscountsClient.BeginDelete(ctx context.Context, resourceGroupName string, discountName string, options *DiscountsClientBeginDeleteOptions) (*runtime.Poller[DiscountsClientDeleteResponse], error)`
- New function `*DiscountsClient.NewResourceGroupListPager(resourceGroupName string, options *DiscountsClientResourceGroupListOptions) *runtime.Pager[DiscountsClientResourceGroupListResponse]`
- New function `*DiscountsClient.NewScopeListPager(scope string, options *DiscountsClientScopeListOptions) *runtime.Pager[DiscountsClientScopeListResponse]`
- New function `*DiscountsClient.NewSubscriptionListPager(options *DiscountsClientSubscriptionListOptions) *runtime.Pager[DiscountsClientSubscriptionListResponse]`
- New function `*EntityTypeAffiliateDiscount.GetDiscountProperties() *DiscountProperties`
- New function `*EntityTypePrimaryDiscount.GetDiscountProperties() *DiscountProperties`
- New struct `CatalogClaimsItem`
- New struct `ConditionsItem`
- New struct `CustomPriceProperties`
- New struct `Discount`
- New struct `DiscountList`
- New struct `DiscountPatchRequest`
- New struct `DiscountPatchRequestProperties`
- New struct `DiscountTypeCustomPriceMultiCurrency`
- New struct `DiscountTypeProduct`
- New struct `DiscountTypeProductFamily`
- New struct `DiscountTypeProductSKU`
- New struct `EntityTypeAffiliateDiscount`
- New struct `EntityTypePrimaryDiscount`
- New struct `ManagedServiceIdentity`
- New struct `MarketSetPricesItems`
- New struct `Plan`
- New struct `PriceGuaranteeProperties`
- New struct `ResourceSKU`
- New struct `UserAssignedIdentity`
- New field `Capacity`, `Family`, `Size`, `Tier` in struct `SKU`
- New field `Renew` in struct `SavingsPlanOrderAliasProperties`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-04-03)
### Breaking Changes

- Function `NewSavingsPlanClient` parameter(s) have been changed from `(*string, azcore.TokenCredential, *arm.ClientOptions)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSavingsPlanOrderClient` parameter(s) have been changed from `(*string, azcore.TokenCredential, *arm.ClientOptions)` to `(azcore.TokenCredential, *arm.ClientOptions)`

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New field `Expand` in struct `SavingsPlanClientGetOptions`
- New field `Expand` in struct `SavingsPlanOrderClientGetOptions`


## 1.0.0 (2022-12-23)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billingbenefits/armbillingbenefits` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).