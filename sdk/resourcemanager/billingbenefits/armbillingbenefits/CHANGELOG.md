# Release History

## 3.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `*RPClient.ValidatePurchase` has been removed
- Operation `*SavingsPlanClient.Update` has been changed to LRO, use `*SavingsPlanClient.BeginUpdate` instead.
- Struct `SavingsPlanPurchaseValidateRequest` has been removed

### Features Added

- New value `CommitmentGrainFullTerm`, `CommitmentGrainUnknown` added to enum type `CommitmentGrain`
- New value `TermP1M` added to enum type `Term`
- New enum type `ApplyDiscountOn` with values `ApplyDiscountOnConsume`, `ApplyDiscountOnPurchase`, `ApplyDiscountOnRenew`
- New enum type `BenefitType` with values `BenefitTypeConditionalCredits`, `BenefitTypeCredits`, `BenefitTypeMACC`, `BenefitTypeSavingsPlan`
- New enum type `ConditionalCreditEntityType` with values `ConditionalCreditEntityTypeContributor`, `ConditionalCreditEntityTypePrimary`
- New enum type `ConditionalCreditStatus` with values `ConditionalCreditStatusActive`, `ConditionalCreditStatusCanceled`, `ConditionalCreditStatusCompleted`, `ConditionalCreditStatusFailed`, `ConditionalCreditStatusPending`, `ConditionalCreditStatusPendingSettlement`, `ConditionalCreditStatusScheduled`, `ConditionalCreditStatusStopped`, `ConditionalCreditStatusUnknown`
- New enum type `ConditionalCreditsProvisioningState` with values `ConditionalCreditsProvisioningStateCanceled`, `ConditionalCreditsProvisioningStateFailed`, `ConditionalCreditsProvisioningStatePending`, `ConditionalCreditsProvisioningStateSucceeded`, `ConditionalCreditsProvisioningStateUnknown`
- New enum type `CreditExpirationPolicy` with values `CreditExpirationPolicyNone`, `CreditExpirationPolicySuspendBillingProfile`
- New enum type `CreditRedemptionPolicy` with values `CreditRedemptionPolicyAutoRedeem`, `CreditRedemptionPolicyManualRedeem`, `CreditRedemptionPolicyNotApplicable`
- New enum type `CreditStatus` with values `CreditStatusActive`, `CreditStatusCanceled`, `CreditStatusExhausted`, `CreditStatusExpired`, `CreditStatusFailed`, `CreditStatusNotStarted`, `CreditStatusPending`, `CreditStatusSucceeded`, `CreditStatusUnknown`
- New enum type `DiscountAppliedScopeType` with values `DiscountAppliedScopeTypeBillingAccount`, `DiscountAppliedScopeTypeBillingProfile`, `DiscountAppliedScopeTypeCustomer`
- New enum type `DiscountCombinationRule` with values `DiscountCombinationRuleBestOf`, `DiscountCombinationRuleStackable`
- New enum type `DiscountEntityType` with values `DiscountEntityTypeAffiliate`, `DiscountEntityTypePrimary`
- New enum type `DiscountProvisioningState` with values `DiscountProvisioningStateCanceled`, `DiscountProvisioningStateFailed`, `DiscountProvisioningStatePending`, `DiscountProvisioningStateSucceeded`, `DiscountProvisioningStateUnknown`
- New enum type `DiscountRuleType` with values `DiscountRuleTypeFixedListPrice`, `DiscountRuleTypeFixedPriceLock`, `DiscountRuleTypePriceCeiling`
- New enum type `DiscountStatus` with values `DiscountStatusActive`, `DiscountStatusCanceled`, `DiscountStatusExpired`, `DiscountStatusFailed`, `DiscountStatusPending`
- New enum type `DiscountType` with values `DiscountTypeCustomPrice`, `DiscountTypeCustomPriceMultiCurrency`, `DiscountTypeProduct`, `DiscountTypeProductFamily`, `DiscountTypeSKU`
- New enum type `EnablementMode` with values `EnablementModeDisabled`, `EnablementModeEnabled`, `EnablementModeUnknown`
- New enum type `FreeServicesStatus` with values `FreeServicesStatusActive`, `FreeServicesStatusCanceled`, `FreeServicesStatusCompleted`, `FreeServicesStatusPending`, `FreeServicesStatusUnknown`
- New enum type `MaccEntityType` with values `MaccEntityTypeContributor`, `MaccEntityTypePrimary`
- New enum type `MaccMilestoneStatus` with values `MaccMilestoneStatusActive`, `MaccMilestoneStatusCanceled`, `MaccMilestoneStatusCompleted`, `MaccMilestoneStatusFailed`, `MaccMilestoneStatusPending`, `MaccMilestoneStatusPendingSettlement`, `MaccMilestoneStatusRemoved`, `MaccMilestoneStatusScheduled`, `MaccMilestoneStatusShortfallCharged`, `MaccMilestoneStatusShortfallWaived`, `MaccMilestoneStatusUnknown`
- New enum type `MaccStatus` with values `MaccStatusActive`, `MaccStatusCanceled`, `MaccStatusCompleted`, `MaccStatusFailed`, `MaccStatusPending`, `MaccStatusPendingSettlement`, `MaccStatusScheduled`, `MaccStatusShortfallCharged`, `MaccStatusShortfallWaived`, `MaccStatusStopped`, `MaccStatusUnknown`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `MilestoneStatus` with values `MilestoneStatusActive`, `MilestoneStatusCanceled`, `MilestoneStatusCompleted`, `MilestoneStatusFailed`, `MilestoneStatusMissed`, `MilestoneStatusPending`, `MilestoneStatusPendingSettlement`, `MilestoneStatusRemoved`, `MilestoneStatusScheduled`, `MilestoneStatusUnknown`
- New enum type `PricingPolicy` with values `PricingPolicyLocked`, `PricingPolicyProtected`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New enum type `ServiceManagedIdentityType` with values `ServiceManagedIdentityTypeNone`, `ServiceManagedIdentityTypeSystemAssigned`, `ServiceManagedIdentityTypeSystemAssignedUserAssigned`, `ServiceManagedIdentityTypeUserAssigned`
- New function `NewApplicableMaccsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ApplicableMaccsClient, error)`
- New function `*ApplicableMaccsClient.NewListPager(billingAccountID string, options *ApplicableMaccsClientListOptions) *runtime.Pager[ApplicableMaccsClientListResponse]`
- New function `*BenefitValidateModel.GetBenefitValidateModel() *BenefitValidateModel`
- New function `*ClientFactory.NewApplicableMaccsClient() *ApplicableMaccsClient`
- New function `*ClientFactory.NewConditionalCreditContributorsClient(subscriptionID string) *ConditionalCreditContributorsClient`
- New function `*ClientFactory.NewConditionalCreditsClient(subscriptionID string) *ConditionalCreditsClient`
- New function `*ClientFactory.NewContributorsClient(subscriptionID string) *ContributorsClient`
- New function `*ClientFactory.NewCreditsClient(subscriptionID string) *CreditsClient`
- New function `*ClientFactory.NewDiscountClient(subscriptionID string) *DiscountClient`
- New function `*ClientFactory.NewDiscountsClient(subscriptionID string) *DiscountsClient`
- New function `*ClientFactory.NewFreeServicesClient(subscriptionID string) *FreeServicesClient`
- New function `*ClientFactory.NewMaccsClient(subscriptionID string) *MaccsClient`
- New function `*ClientFactory.NewSellerResourceClient() *SellerResourceClient`
- New function `*ClientFactory.NewSourcesClient(subscriptionID string) *SourcesClient`
- New function `NewConditionalCreditContributorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConditionalCreditContributorsClient, error)`
- New function `*ConditionalCreditContributorsClient.GetFromPrimary(ctx context.Context, resourceGroupName string, conditionalCreditName string, contributorName string, options *ConditionalCreditContributorsClientGetFromPrimaryOptions) (ConditionalCreditContributorsClientGetFromPrimaryResponse, error)`
- New function `*ConditionalCreditContributorsClient.NewListFromApplicableConditionalCreditPager(billingAccountID string, systemID string, options *ConditionalCreditContributorsClientListFromApplicableConditionalCreditOptions) *runtime.Pager[ConditionalCreditContributorsClientListFromApplicableConditionalCreditResponse]`
- New function `*ConditionalCreditContributorsClient.NewListFromPrimaryPager(resourceGroupName string, conditionalCreditName string, options *ConditionalCreditContributorsClientListFromPrimaryOptions) *runtime.Pager[ConditionalCreditContributorsClientListFromPrimaryResponse]`
- New function `*ConditionalCreditProperties.GetConditionalCreditProperties() *ConditionalCreditProperties`
- New function `NewConditionalCreditsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConditionalCreditsClient, error)`
- New function `*ConditionalCreditsClient.BeginCancel(ctx context.Context, resourceGroupName string, conditionalCreditName string, options *ConditionalCreditsClientBeginCancelOptions) (*runtime.Poller[ConditionalCreditsClientCancelResponse], error)`
- New function `*ConditionalCreditsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, conditionalCreditName string, body ConditionalCredit, options *ConditionalCreditsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConditionalCreditsClientCreateOrUpdateResponse], error)`
- New function `*ConditionalCreditsClient.BeginDelete(ctx context.Context, resourceGroupName string, conditionalCreditName string, options *ConditionalCreditsClientBeginDeleteOptions) (*runtime.Poller[ConditionalCreditsClientDeleteResponse], error)`
- New function `*ConditionalCreditsClient.Get(ctx context.Context, resourceGroupName string, conditionalCreditName string, options *ConditionalCreditsClientGetOptions) (ConditionalCreditsClientGetResponse, error)`
- New function `*ConditionalCreditsClient.NewListByResourceGroupPager(resourceGroupName string, options *ConditionalCreditsClientListByResourceGroupOptions) *runtime.Pager[ConditionalCreditsClientListByResourceGroupResponse]`
- New function `*ConditionalCreditsClient.NewListBySubscriptionPager(options *ConditionalCreditsClientListBySubscriptionOptions) *runtime.Pager[ConditionalCreditsClientListBySubscriptionResponse]`
- New function `*ConditionalCreditsClient.NewScopeListPager(scope string, options *ConditionalCreditsClientScopeListOptions) *runtime.Pager[ConditionalCreditsClientScopeListResponse]`
- New function `*ConditionalCreditsClient.BeginUpdate(ctx context.Context, resourceGroupName string, conditionalCreditName string, body ConditionalCreditPatchRequest, options *ConditionalCreditsClientBeginUpdateOptions) (*runtime.Poller[ConditionalCreditsClientUpdateResponse], error)`
- New function `*ConditionalCreditsValidateModel.GetBenefitValidateModel() *BenefitValidateModel`
- New function `*ContributorConditionalCreditProperties.GetConditionalCreditProperties() *ConditionalCreditProperties`
- New function `NewContributorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContributorsClient, error)`
- New function `*ContributorsClient.GetFromPrimary(ctx context.Context, resourceGroupName string, maccName string, contributorName string, options *ContributorsClientGetFromPrimaryOptions) (ContributorsClientGetFromPrimaryResponse, error)`
- New function `*ContributorsClient.NewListFromApplicableMaccPager(billingAccountID string, systemID string, options *ContributorsClientListFromApplicableMaccOptions) *runtime.Pager[ContributorsClientListFromApplicableMaccResponse]`
- New function `*ContributorsClient.NewListFromPrimaryPager(resourceGroupName string, maccName string, options *ContributorsClientListFromPrimaryOptions) *runtime.Pager[ContributorsClientListFromPrimaryResponse]`
- New function `NewCreditsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CreditsClient, error)`
- New function `*CreditsClient.BeginCancel(ctx context.Context, resourceGroupName string, creditName string, options *CreditsClientBeginCancelOptions) (*runtime.Poller[CreditsClientCancelResponse], error)`
- New function `*CreditsClient.BeginCreate(ctx context.Context, resourceGroupName string, creditName string, body Credit, options *CreditsClientBeginCreateOptions) (*runtime.Poller[CreditsClientCreateResponse], error)`
- New function `*CreditsClient.BeginDelete(ctx context.Context, resourceGroupName string, creditName string, options *CreditsClientBeginDeleteOptions) (*runtime.Poller[CreditsClientDeleteResponse], error)`
- New function `*CreditsClient.Get(ctx context.Context, resourceGroupName string, creditName string, options *CreditsClientGetOptions) (CreditsClientGetResponse, error)`
- New function `*CreditsClient.NewListApplicablePager(scope string, options *CreditsClientListApplicableOptions) *runtime.Pager[CreditsClientListApplicableResponse]`
- New function `*CreditsClient.NewListByResourceGroupPager(resourceGroupName string, options *CreditsClientListByResourceGroupOptions) *runtime.Pager[CreditsClientListByResourceGroupResponse]`
- New function `*CreditsClient.NewListBySubscriptionPager(options *CreditsClientListBySubscriptionOptions) *runtime.Pager[CreditsClientListBySubscriptionResponse]`
- New function `*CreditsClient.BeginUpdate(ctx context.Context, resourceGroupName string, creditName string, body CreditPatchRequest, options *CreditsClientBeginUpdateOptions) (*runtime.Poller[CreditsClientUpdateResponse], error)`
- New function `*CreditsValidateModel.GetBenefitValidateModel() *BenefitValidateModel`
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
- New function `NewFreeServicesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FreeServicesClient, error)`
- New function `*FreeServicesClient.BeginCreate(ctx context.Context, resourceGroupName string, freeServiceName string, body FreeServices, options *FreeServicesClientBeginCreateOptions) (*runtime.Poller[FreeServicesClientCreateResponse], error)`
- New function `*FreeServicesClient.BeginDelete(ctx context.Context, resourceGroupName string, freeServiceName string, options *FreeServicesClientBeginDeleteOptions) (*runtime.Poller[FreeServicesClientDeleteResponse], error)`
- New function `*FreeServicesClient.Get(ctx context.Context, resourceGroupName string, freeServiceName string, options *FreeServicesClientGetOptions) (FreeServicesClientGetResponse, error)`
- New function `*FreeServicesClient.NewListByResourceGroupPager(resourceGroupName string, options *FreeServicesClientListByResourceGroupOptions) *runtime.Pager[FreeServicesClientListByResourceGroupResponse]`
- New function `*FreeServicesClient.NewListBySubscriptionPager(options *FreeServicesClientListBySubscriptionOptions) *runtime.Pager[FreeServicesClientListBySubscriptionResponse]`
- New function `*FreeServicesClient.BeginUpdate(ctx context.Context, resourceGroupName string, freeServiceName string, body FreeServicesPatchRequest, options *FreeServicesClientBeginUpdateOptions) (*runtime.Poller[FreeServicesClientUpdateResponse], error)`
- New function `*MaccValidateModel.GetBenefitValidateModel() *BenefitValidateModel`
- New function `NewMaccsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MaccsClient, error)`
- New function `*MaccsClient.BeginCancel(ctx context.Context, resourceGroupName string, maccName string, options *MaccsClientBeginCancelOptions) (*runtime.Poller[MaccsClientCancelResponse], error)`
- New function `*MaccsClient.BeginChargeShortfall(ctx context.Context, resourceGroupName string, maccName string, body ChargeShortfallRequest, options *MaccsClientBeginChargeShortfallOptions) (*runtime.Poller[MaccsClientChargeShortfallResponse], error)`
- New function `*MaccsClient.BeginCreate(ctx context.Context, resourceGroupName string, maccName string, body Macc, options *MaccsClientBeginCreateOptions) (*runtime.Poller[MaccsClientCreateResponse], error)`
- New function `*MaccsClient.BeginDelete(ctx context.Context, resourceGroupName string, maccName string, options *MaccsClientBeginDeleteOptions) (*runtime.Poller[MaccsClientDeleteResponse], error)`
- New function `*MaccsClient.Get(ctx context.Context, resourceGroupName string, maccName string, options *MaccsClientGetOptions) (MaccsClientGetResponse, error)`
- New function `*MaccsClient.NewListByResourceGroupPager(resourceGroupName string, options *MaccsClientListByResourceGroupOptions) *runtime.Pager[MaccsClientListByResourceGroupResponse]`
- New function `*MaccsClient.NewListBySubscriptionPager(options *MaccsClientListBySubscriptionOptions) *runtime.Pager[MaccsClientListBySubscriptionResponse]`
- New function `*MaccsClient.BeginUpdate(ctx context.Context, resourceGroupName string, maccName string, body MaccPatchRequest, options *MaccsClientBeginUpdateOptions) (*runtime.Poller[MaccsClientUpdateResponse], error)`
- New function `*MaccsClient.BeginWriteOff(ctx context.Context, resourceGroupName string, maccName string, options *MaccsClientBeginWriteOffOptions) (*runtime.Poller[MaccsClientWriteOffResponse], error)`
- New function `*PrimaryConditionalCreditProperties.GetConditionalCreditProperties() *ConditionalCreditProperties`
- New function `*RPClient.Validate(ctx context.Context, body BenefitValidateRequest, options *RPClientValidateOptions) (RPClientValidateResponse, error)`
- New function `*SavingsPlanValidateModel.GetBenefitValidateModel() *BenefitValidateModel`
- New function `NewSellerResourceClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SellerResourceClient, error)`
- New function `*SellerResourceClient.List(ctx context.Context, body SellerResourceListRequest, options *SellerResourceClientListOptions) (SellerResourceClientListResponse, error)`
- New function `NewSourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SourcesClient, error)`
- New function `*SourcesClient.Create(ctx context.Context, resourceGroupName string, creditName string, sourceName string, body CreditSource, options *SourcesClientCreateOptions) (SourcesClientCreateResponse, error)`
- New function `*SourcesClient.Delete(ctx context.Context, resourceGroupName string, creditName string, sourceName string, options *SourcesClientDeleteOptions) (SourcesClientDeleteResponse, error)`
- New function `*SourcesClient.Get(ctx context.Context, resourceGroupName string, creditName string, sourceName string, options *SourcesClientGetOptions) (SourcesClientGetResponse, error)`
- New function `*SourcesClient.NewListByCreditPager(resourceGroupName string, creditName string, options *SourcesClientListByCreditOptions) *runtime.Pager[SourcesClientListByCreditResponse]`
- New function `*SourcesClient.Update(ctx context.Context, resourceGroupName string, creditName string, sourceName string, body CreditSourcePatchRequest, options *SourcesClientUpdateOptions) (SourcesClientUpdateResponse, error)`
- New struct `ApplicableMacc`
- New struct `ApplicableMaccList`
- New struct `AutomaticShortfallSuppressReason`
- New struct `Award`
- New struct `BenefitValidateRequest`
- New struct `BenefitValidateResponse`
- New struct `BenefitValidateResponseProperty`
- New struct `CatalogClaimsItem`
- New struct `ChargeShortfallRequest`
- New struct `ConditionalCredit`
- New struct `ConditionalCreditContributor`
- New struct `ConditionalCreditContributorList`
- New struct `ConditionalCreditList`
- New struct `ConditionalCreditMilestone`
- New struct `ConditionalCreditPatchRequest`
- New struct `ConditionalCreditPatchRequestProperties`
- New struct `ConditionalCreditsValidateModel`
- New struct `ConditionsItem`
- New struct `Contributor`
- New struct `ContributorConditionalCreditMilestone`
- New struct `ContributorConditionalCreditProperties`
- New struct `ContributorList`
- New struct `Credit`
- New struct `CreditBreakdownItem`
- New struct `CreditDimension`
- New struct `CreditPatchProperties`
- New struct `CreditPatchRequest`
- New struct `CreditPolicies`
- New struct `CreditProperties`
- New struct `CreditReason`
- New struct `CreditSource`
- New struct `CreditSourcePatchRequest`
- New struct `CreditSourceProperties`
- New struct `CreditSourcesList`
- New struct `CreditsList`
- New struct `CreditsValidateModel`
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
- New struct `FreeServices`
- New struct `FreeServicesList`
- New struct `FreeServicesPatchRequest`
- New struct `FreeServicesPatchRequestProperties`
- New struct `FreeServicesProperties`
- New struct `Macc`
- New struct `MaccList`
- New struct `MaccMilestone`
- New struct `MaccModelProperties`
- New struct `MaccPatchRequest`
- New struct `MaccPatchRequestProperties`
- New struct `MaccValidateModel`
- New struct `ManagedServiceIdentity`
- New struct `MarketSetPricesItems`
- New struct `Plan`
- New struct `PriceGuaranteeProperties`
- New struct `PrimaryConditionalCreditProperties`
- New struct `SavingsPlanValidateModel`
- New struct `SellerResourceListRequest`
- New struct `SellerResourceListRequestProperties`
- New struct `ServiceManagedIdentity`
- New struct `Shortfall`
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