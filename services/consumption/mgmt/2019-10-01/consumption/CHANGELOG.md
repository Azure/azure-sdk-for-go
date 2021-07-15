# Unreleased

## Breaking Changes

### Removed Constants

1. Scope9.Scope9Shared
1. Scope9.Scope9Single

### Removed Funcs

1. PossibleScope11Values() []Scope11
1. PossibleScope9Values() []Scope9

### Signature Changes

#### Const Types

1. Shared changed type from Scope11 to Scope12
1. Single changed type from Scope11 to Scope12

#### Funcs

1. CreditsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. CreditsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. EventsClient.List
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string
1. EventsClient.ListComplete
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string
1. EventsClient.ListPreparer
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string
1. LotsClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. LotsClient.ListComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. LotsClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. ReservationRecommendationDetailsClient.Get
	- Params
		- From: context.Context, string, Scope11, string, Term, LookBackPeriod, string
		- To: context.Context, string, Scope14, string, Term, LookBackPeriod, string
1. ReservationRecommendationDetailsClient.GetPreparer
	- Params
		- From: context.Context, string, Scope11, string, Term, LookBackPeriod, string
		- To: context.Context, string, Scope14, string, Term, LookBackPeriod, string

#### Struct Fields

1. LegacyReservationRecommendationProperties.InstanceFlexibilityRatio changed type from *int32 to *float64
1. ModernReservationRecommendationProperties.InstanceFlexibilityRatio changed type from *int32 to *float64
1. ModernReservationRecommendationProperties.LookBackPeriod changed type from *string to *int32
1. ModernUsageDetailProperties.MeterID changed type from *uuid.UUID to *string

## Additive Changes

### New Constants

1. CultureCode.CsCz
1. CultureCode.DaDk
1. CultureCode.DeDe
1. CultureCode.EnGb
1. CultureCode.EnUs
1. CultureCode.EsEs
1. CultureCode.FrFr
1. CultureCode.HuHu
1. CultureCode.ItIt
1. CultureCode.JaJp
1. CultureCode.KoKr
1. CultureCode.NbNo
1. CultureCode.NlNl
1. CultureCode.PlPl
1. CultureCode.PtBr
1. CultureCode.PtPt
1. CultureCode.RuRu
1. CultureCode.SvSe
1. CultureCode.TrTr
1. CultureCode.ZhCn
1. CultureCode.ZhTw
1. Scope14.Scope14Shared
1. Scope14.Scope14Single

### New Funcs

1. AmountWithExchangeRate.MarshalJSON() ([]byte, error)
1. DownloadProperties.MarshalJSON() ([]byte, error)
1. ForecastSpend.MarshalJSON() ([]byte, error)
1. HighCasedErrorDetails.MarshalJSON() ([]byte, error)
1. PossibleCultureCodeValues() []CultureCode
1. PossibleScope12Values() []Scope12
1. PossibleScope14Values() []Scope14
1. Reseller.MarshalJSON() ([]byte, error)
1. TagProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AmountWithExchangeRate
1. DownloadProperties
1. ForecastSpend
1. HighCasedErrorDetails
1. HighCasedErrorResponse
1. Reseller

#### New Struct Fields

1. Balance.Etag
1. BudgetProperties.ForecastSpend
1. ChargeSummary.Etag
1. CreditBalanceSummary.CurrentBalanceInBillingCurrency
1. CreditBalanceSummary.EstimatedBalanceInBillingCurrency
1. CreditSummary.Etag
1. CreditSummaryProperties.BillingCurrency
1. CreditSummaryProperties.CreditCurrency
1. CreditSummaryProperties.Reseller
1. EventProperties.AdjustmentsInBillingCurrency
1. EventProperties.BillingCurrency
1. EventProperties.ChargesInBillingCurrency
1. EventProperties.ClosedBalanceInBillingCurrency
1. EventProperties.CreditCurrency
1. EventProperties.CreditExpiredInBillingCurrency
1. EventProperties.NewCreditInBillingCurrency
1. EventProperties.Reseller
1. EventSummary.Etag
1. Forecast.Etag
1. LegacyChargeSummary.Etag
1. LegacyReservationRecommendation.Etag
1. LegacyReservationRecommendationProperties.ResourceType
1. LegacyUsageDetail.Etag
1. LotProperties.BillingCurrency
1. LotProperties.ClosedBalanceInBillingCurrency
1. LotProperties.CreditCurrency
1. LotProperties.OriginalAmountInBillingCurrency
1. LotProperties.Reseller
1. LotSummary.Etag
1. ManagementGroupAggregatedCostResult.Etag
1. Marketplace.Etag
1. MarketplaceProperties.AdditionalInfo
1. ModernChargeSummary.Etag
1. ModernReservationRecommendation.ETag
1. ModernReservationRecommendation.Etag
1. ModernReservationRecommendationProperties.Location
1. ModernReservationRecommendationProperties.ResourceType
1. ModernReservationRecommendationProperties.SkuName
1. ModernReservationRecommendationProperties.SubscriptionID
1. ModernUsageDetail.Etag
1. ModernUsageDetailProperties.PayGPrice
1. Notification.Locale
1. Operation.ID
1. OperationDisplay.Description
1. PriceSheetModel.Download
1. PriceSheetResult.Etag
1. ReservationDetail.Etag
1. ReservationRecommendation.Etag
1. ReservationRecommendationDetailsModel.ETag
1. ReservationRecommendationDetailsModel.Etag
1. ReservationRecommendationsListResult.PreviousLink
1. ReservationRecommendationsListResult.TotalCost
1. ReservationSummary.Etag
1. Resource.Etag
1. Tag.Value
1. TagProperties.NextLink
1. TagProperties.PreviousLink
1. UsageDetail.Etag
