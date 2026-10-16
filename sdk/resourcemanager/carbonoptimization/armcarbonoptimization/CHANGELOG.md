# Release History

## 2.0.0-beta.1 (2026-10-16)
### Breaking Changes

- Operation `*CarbonServiceClient.QueryCarbonEmissionReports` has supported pagination, use `*CarbonServiceClient.NewQueryCarbonEmissionReportsPager` instead.

### Features Added

- New enum type `BillingAccountCategoryTypeEnum` with values `BillingAccountCategoryTypeEnumLocation`, `BillingAccountCategoryTypeEnumResourceType`, `BillingAccountCategoryTypeEnumSubscription`
- New enum type `BillingAccountOrderByColumnEnum` with values `BillingAccountOrderByColumnEnumItemName`, `BillingAccountOrderByColumnEnumLatestMonthEmissions`, `BillingAccountOrderByColumnEnumMonthOverMonthEmissionsChangeRatio`, `BillingAccountOrderByColumnEnumMonthlyEmissionsChangeValue`, `BillingAccountOrderByColumnEnumPreviousMonthEmissions`
- New function `*BillingAccountItemDetailsReportQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*BillingAccountMonthlySummaryReportQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*BillingAccountOverallSummaryReportQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*BillingAccountQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*BillingAccountTopItemsMonthlySummaryReportQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*BillingAccountTopItemsSummaryReportQueryFilter.GetBillingAccountQueryFilter() *BillingAccountQueryFilter`
- New function `*CarbonServiceClient.NewQueryCarbonEmissionReportsByBillingAccountPager(scope string, queryParameters BillingAccountQueryFilterClassification, options *CarbonServiceClientQueryCarbonEmissionReportsByBillingAccountOptions) *runtime.Pager[CarbonServiceClientQueryCarbonEmissionReportsByBillingAccountResponse]`
- New struct `BillingAccountItemDetailsReportQueryFilter`
- New struct `BillingAccountMonthlySummaryReportQueryFilter`
- New struct `BillingAccountOverallSummaryReportQueryFilter`
- New struct `BillingAccountTopItemsMonthlySummaryReportQueryFilter`
- New struct `BillingAccountTopItemsSummaryReportQueryFilter`


## 1.0.0 (2025-07-21)
### Other Changes

- Release stable version.
- 
## 0.1.0 (2025-05-02)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/carbonoptimization/armcarbonoptimization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
