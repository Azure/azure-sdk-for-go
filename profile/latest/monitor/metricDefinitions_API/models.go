package metricdefinitionsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/2016-03-01/metricDefinitions_API"
)

type (
	 ManagementClient = original.ManagementClient
	 MetricDefinitionsClient = original.MetricDefinitionsClient
	 AggregationType = original.AggregationType
	 Unit = original.Unit
	 ErrorResponse = original.ErrorResponse
	 LocalizableString = original.LocalizableString
	 MetricAvailability = original.MetricAvailability
	 MetricDefinition = original.MetricDefinition
	 MetricDefinitionCollection = original.MetricDefinitionCollection
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Average = original.Average
	 Count = original.Count
	 Maximum = original.Maximum
	 Minimum = original.Minimum
	 None = original.None
	 Total = original.Total
	 UnitBytes = original.UnitBytes
	 UnitBytesPerSecond = original.UnitBytesPerSecond
	 UnitCount = original.UnitCount
	 UnitCountPerSecond = original.UnitCountPerSecond
	 UnitMilliSeconds = original.UnitMilliSeconds
	 UnitPercent = original.UnitPercent
	 UnitSeconds = original.UnitSeconds
)
