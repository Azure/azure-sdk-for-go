package metricsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/2016-09-01/metrics_API"
)

type (
	 ManagementClient = original.ManagementClient
	 MetricsClient = original.MetricsClient
	 Unit = original.Unit
	 ErrorResponse = original.ErrorResponse
	 LocalizableString = original.LocalizableString
	 Metric = original.Metric
	 MetricCollection = original.MetricCollection
	 MetricValue = original.MetricValue
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Bytes = original.Bytes
	 BytesPerSecond = original.BytesPerSecond
	 Count = original.Count
	 CountPerSecond = original.CountPerSecond
	 MilliSeconds = original.MilliSeconds
	 Percent = original.Percent
	 Seconds = original.Seconds
)
