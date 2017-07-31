package commerce

import (
	 original "github.com/Azure/azure-sdk-for-go/service/commerce/management/2015-06-01-preview/commerce"
)

type (
	 ManagementClient = original.ManagementClient
	 AggregationGranularity = original.AggregationGranularity
	 ErrorResponse = original.ErrorResponse
	 InfoField = original.InfoField
	 MeterInfo = original.MeterInfo
	 MonetaryCommitment = original.MonetaryCommitment
	 MonetaryCredit = original.MonetaryCredit
	 OfferTermInfo = original.OfferTermInfo
	 RateCardQueryParameters = original.RateCardQueryParameters
	 RecurringCharge = original.RecurringCharge
	 ResourceRateCardInfo = original.ResourceRateCardInfo
	 UsageAggregation = original.UsageAggregation
	 UsageAggregationListResult = original.UsageAggregationListResult
	 UsageSample = original.UsageSample
	 RateCardClient = original.RateCardClient
	 UsageAggregatesClient = original.UsageAggregatesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Daily = original.Daily
	 Hourly = original.Hourly
)
