# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. CustomDomainsClient.DisableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.DisableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error
1. LogAnalyticsClient.GetLogAnalyticsMetrics
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string, []string, []string
		- To: context.Context, string, string, []LogMetric, date.Time, date.Time, LogMetricsGranularity, []string, []string, []LogMetricsGroupBy, []string, []string
1. LogAnalyticsClient.GetLogAnalyticsMetricsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string, []string, []string
		- To: context.Context, string, string, []LogMetric, date.Time, date.Time, LogMetricsGranularity, []string, []string, []LogMetricsGroupBy, []string, []string
1. LogAnalyticsClient.GetLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
		- To: context.Context, string, string, []LogRanking, []LogRankingMetric, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
		- To: context.Context, string, string, []LogRanking, []LogRankingMetric, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetWafLogAnalyticsMetrics
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, WafGranularity, []WafAction, []WafRankingGroupBy, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsMetricsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, WafGranularity, []WafAction, []WafRankingGroupBy, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, int32, []WafRankingType, []WafAction, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, int32, []WafRankingType, []WafAction, []WafRuleType
1. RuleSetsClient.Create
	- Params
		- From: context.Context, string, string, string, RuleSet
		- To: context.Context, string, string, string
1. RuleSetsClient.CreatePreparer
	- Params
		- From: context.Context, string, string, string, RuleSet
		- To: context.Context, string, string, string

## Additive Changes

### New Constants

1. LogMetric.ClientRequestBandwidth
1. LogMetric.ClientRequestCount
1. LogMetric.ClientRequestTraffic
1. LogMetric.OriginRequestBandwidth
1. LogMetric.OriginRequestTraffic
1. LogMetric.TotalLatency
1. LogMetricsGranularity.LogMetricsGranularityP1D
1. LogMetricsGranularity.LogMetricsGranularityPT1H
1. LogMetricsGranularity.LogMetricsGranularityPT5M
1. LogMetricsGroupBy.LogMetricsGroupByCacheStatus
1. LogMetricsGroupBy.LogMetricsGroupByCountry
1. LogMetricsGroupBy.LogMetricsGroupByCustomDomain
1. LogMetricsGroupBy.LogMetricsGroupByHTTPStatusCode
1. LogMetricsGroupBy.LogMetricsGroupByProtocol
1. LogRanking.LogRankingBrowser
1. LogRanking.LogRankingCountryOrRegion
1. LogRanking.LogRankingReferrer
1. LogRanking.LogRankingURL
1. LogRanking.LogRankingUserAgent
1. LogRankingMetric.LogRankingMetricClientRequestCount
1. LogRankingMetric.LogRankingMetricClientRequestTraffic
1. LogRankingMetric.LogRankingMetricErrorCount
1. LogRankingMetric.LogRankingMetricHitCount
1. LogRankingMetric.LogRankingMetricMissCount
1. LogRankingMetric.LogRankingMetricUserErrorCount
1. WafAction.WafActionAllow
1. WafAction.WafActionBlock
1. WafAction.WafActionLog
1. WafAction.WafActionRedirect
1. WafGranularity.WafGranularityP1D
1. WafGranularity.WafGranularityPT1H
1. WafGranularity.WafGranularityPT5M
1. WafMetric.WafMetricClientRequestCount
1. WafRankingGroupBy.WafRankingGroupByCustomDomain
1. WafRankingGroupBy.WafRankingGroupByHTTPStatusCode
1. WafRankingType.WafRankingTypeAction
1. WafRankingType.WafRankingTypeClientIP
1. WafRankingType.WafRankingTypeCountry
1. WafRankingType.WafRankingTypeRuleGroup
1. WafRankingType.WafRankingTypeRuleID
1. WafRankingType.WafRankingTypeRuleType
1. WafRankingType.WafRankingTypeURL
1. WafRankingType.WafRankingTypeUserAgent
1. WafRuleType.Bot
1. WafRuleType.Custom
1. WafRuleType.Managed

### New Funcs

1. *CustomDomainsDisableCustomHTTPSFuture.UnmarshalJSON([]byte) error
1. *CustomDomainsEnableCustomHTTPSFuture.UnmarshalJSON([]byte) error
1. PossibleLogMetricValues() []LogMetric
1. PossibleLogMetricsGranularityValues() []LogMetricsGranularity
1. PossibleLogMetricsGroupByValues() []LogMetricsGroupBy
1. PossibleLogRankingMetricValues() []LogRankingMetric
1. PossibleLogRankingValues() []LogRanking
1. PossibleWafActionValues() []WafAction
1. PossibleWafGranularityValues() []WafGranularity
1. PossibleWafMetricValues() []WafMetric
1. PossibleWafRankingGroupByValues() []WafRankingGroupBy
1. PossibleWafRankingTypeValues() []WafRankingType
1. PossibleWafRuleTypeValues() []WafRuleType

### Struct Changes

#### New Structs

1. CustomDomainsDisableCustomHTTPSFuture
1. CustomDomainsEnableCustomHTTPSFuture

#### New Struct Fields

1. ManagedRuleSetDefinition.SystemData
1. Resource.SystemData
