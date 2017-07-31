package usagemetricsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/2014-04-01/usageMetrics_API"
)

type (
	 ManagementClient = original.ManagementClient
	 ErrorResponse = original.ErrorResponse
	 LocalizableString = original.LocalizableString
	 UsageMetric = original.UsageMetric
	 UsageMetricCollection = original.UsageMetricCollection
	 UsageMetricsClient = original.UsageMetricsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
