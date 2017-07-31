package autoscaleapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/management/2015-04-01/autoscale_API"
)

type (
	 AutoscaleSettingsClient = original.AutoscaleSettingsClient
	 ManagementClient = original.ManagementClient
	 ComparisonOperationType = original.ComparisonOperationType
	 MetricStatisticType = original.MetricStatisticType
	 RecurrenceFrequency = original.RecurrenceFrequency
	 ScaleDirection = original.ScaleDirection
	 ScaleType = original.ScaleType
	 TimeAggregationType = original.TimeAggregationType
	 AutoscaleNotification = original.AutoscaleNotification
	 AutoscaleProfile = original.AutoscaleProfile
	 AutoscaleSetting = original.AutoscaleSetting
	 AutoscaleSettingResource = original.AutoscaleSettingResource
	 AutoscaleSettingResourceCollection = original.AutoscaleSettingResourceCollection
	 AutoscaleSettingResourcePatch = original.AutoscaleSettingResourcePatch
	 EmailNotification = original.EmailNotification
	 ErrorResponse = original.ErrorResponse
	 MetricTrigger = original.MetricTrigger
	 Recurrence = original.Recurrence
	 RecurrentSchedule = original.RecurrentSchedule
	 Resource = original.Resource
	 ScaleAction = original.ScaleAction
	 ScaleCapacity = original.ScaleCapacity
	 ScaleRule = original.ScaleRule
	 TimeWindow = original.TimeWindow
	 WebhookNotification = original.WebhookNotification
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Equals = original.Equals
	 GreaterThan = original.GreaterThan
	 GreaterThanOrEqual = original.GreaterThanOrEqual
	 LessThan = original.LessThan
	 LessThanOrEqual = original.LessThanOrEqual
	 NotEquals = original.NotEquals
	 Average = original.Average
	 Max = original.Max
	 Min = original.Min
	 Sum = original.Sum
	 Day = original.Day
	 Hour = original.Hour
	 Minute = original.Minute
	 Month = original.Month
	 None = original.None
	 Second = original.Second
	 Week = original.Week
	 Year = original.Year
	 ScaleDirectionDecrease = original.ScaleDirectionDecrease
	 ScaleDirectionIncrease = original.ScaleDirectionIncrease
	 ScaleDirectionNone = original.ScaleDirectionNone
	 ChangeCount = original.ChangeCount
	 ExactCount = original.ExactCount
	 PercentChangeCount = original.PercentChangeCount
	 TimeAggregationTypeAverage = original.TimeAggregationTypeAverage
	 TimeAggregationTypeCount = original.TimeAggregationTypeCount
	 TimeAggregationTypeMaximum = original.TimeAggregationTypeMaximum
	 TimeAggregationTypeMinimum = original.TimeAggregationTypeMinimum
	 TimeAggregationTypeTotal = original.TimeAggregationTypeTotal
)
