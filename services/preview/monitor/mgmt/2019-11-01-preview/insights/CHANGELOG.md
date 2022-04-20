# Unreleased

## Breaking Changes

### Removed Constants

1. Operator.OperatorNotEquals
1. Sensitivity.SensitivityHigh
1. Sensitivity.SensitivityLow
1. Sensitivity.SensitivityMedium
1. Unit.UnitBitsPerSecond
1. Unit.UnitByteSeconds
1. Unit.UnitBytes
1. Unit.UnitBytesPerSecond
1. Unit.UnitCores
1. Unit.UnitCount
1. Unit.UnitCountPerSecond
1. Unit.UnitMilliCores
1. Unit.UnitMilliSeconds
1. Unit.UnitNanoCores
1. Unit.UnitPercent
1. Unit.UnitSeconds
1. Unit.UnitUnspecified

### Removed Funcs

1. *BaselineResponse.UnmarshalJSON([]byte) error
1. BaselineResponse.MarshalJSON() ([]byte, error)
1. MetricBaselineClient.CalculateBaseline(context.Context, string, TimeSeriesInformation) (CalculateBaselineResponse, error)
1. MetricBaselineClient.CalculateBaselinePreparer(context.Context, string, TimeSeriesInformation) (*http.Request, error)
1. MetricBaselineClient.CalculateBaselineResponder(*http.Response) (CalculateBaselineResponse, error)
1. MetricBaselineClient.CalculateBaselineSender(*http.Request) (*http.Response, error)
1. MetricBaselineClient.Get(context.Context, string, string, string, *string, string, string, ResultType) (BaselineResponse, error)
1. MetricBaselineClient.GetPreparer(context.Context, string, string, string, *string, string, string, ResultType) (*http.Request, error)
1. MetricBaselineClient.GetResponder(*http.Response) (BaselineResponse, error)
1. MetricBaselineClient.GetSender(*http.Request) (*http.Response, error)
1. NewMetricBaselineClient(string) MetricBaselineClient
1. NewMetricBaselineClientWithBaseURI(string, string) MetricBaselineClient
1. PossibleSensitivityValues() []Sensitivity
1. PossibleUnitValues() []Unit

### Struct Changes

#### Removed Structs

1. Baseline
1. BaselineMetadataValue
1. BaselineProperties
1. BaselineResponse
1. CalculateBaselineResponse
1. MetricBaselineClient
1. TimeSeriesInformation

#### Removed Struct Fields

1. TimeSeriesBaseline.Metadata

### Signature Changes

#### Struct Fields

1. DynamicMetricCriteria.TimeAggregation changed type from interface{} to AggregationTypeEnum
1. Metric.Unit changed type from Unit to MetricUnit
1. MetricCriteria.TimeAggregation changed type from interface{} to AggregationTypeEnum
1. MetricDefinition.Unit changed type from Unit to MetricUnit
1. MultiMetricCriteria.TimeAggregation changed type from interface{} to AggregationTypeEnum

## Additive Changes

### New Constants

1. AggregationTypeEnum.AggregationTypeEnumAverage
1. AggregationTypeEnum.AggregationTypeEnumCount
1. AggregationTypeEnum.AggregationTypeEnumMaximum
1. AggregationTypeEnum.AggregationTypeEnumMinimum
1. AggregationTypeEnum.AggregationTypeEnumTotal
1. ConditionalOperator.ConditionalOperatorGreaterThanOrEqual
1. ConditionalOperator.ConditionalOperatorLessThanOrEqual
1. MetricClass.Availability
1. MetricClass.Errors
1. MetricClass.Latency
1. MetricClass.Saturation
1. MetricClass.Transactions
1. MetricStatisticType.MetricStatisticTypeCount
1. MetricUnit.MetricUnitBitsPerSecond
1. MetricUnit.MetricUnitByteSeconds
1. MetricUnit.MetricUnitBytes
1. MetricUnit.MetricUnitBytesPerSecond
1. MetricUnit.MetricUnitCores
1. MetricUnit.MetricUnitCount
1. MetricUnit.MetricUnitCountPerSecond
1. MetricUnit.MetricUnitMilliCores
1. MetricUnit.MetricUnitMilliSeconds
1. MetricUnit.MetricUnitNanoCores
1. MetricUnit.MetricUnitPercent
1. MetricUnit.MetricUnitSeconds
1. MetricUnit.MetricUnitUnspecified
1. NamespaceClassification.Custom
1. NamespaceClassification.Platform
1. NamespaceClassification.Qos
1. ScaleType.ServiceAllowedNextValue

### New Funcs

1. AzureResource.MarshalJSON() ([]byte, error)
1. PossibleAggregationTypeEnumValues() []AggregationTypeEnum
1. PossibleMetricClassValues() []MetricClass
1. PossibleMetricUnitValues() []MetricUnit
1. PossibleNamespaceClassificationValues() []NamespaceClassification

### Struct Changes

#### New Structs

1. AzureResource
1. ErrorContract

#### New Struct Fields

1. ActionGroupResource.Identity
1. ActionGroupResource.Kind
1. ActivityLogAlertResource.Etag
1. ActivityLogAlertResource.Kind
1. AlertRule.Action
1. AlertRule.ProvisioningState
1. AlertRuleResource.Etag
1. AlertRuleResource.Kind
1. AutoscaleSetting.TargetResourceLocation
1. AutoscaleSettingResource.Etag
1. AutoscaleSettingResource.Kind
1. LogProfileResource.Etag
1. LogProfileResource.Kind
1. LogSearchRule.AutoMitigate
1. LogSearchRule.CreatedWithAPIVersion
1. LogSearchRule.DisplayName
1. LogSearchRule.IsLegacyLogAnalyticsRule
1. LogSearchRuleResource.Etag
1. LogSearchRuleResource.Kind
1. Metric.DisplayDescription
1. Metric.ErrorCode
1. Metric.ErrorMessage
1. MetricAlertProperties.IsMigrated
1. MetricAlertPropertiesPatch.IsMigrated
1. MetricAlertResource.Etag
1. MetricAlertResource.Kind
1. MetricDefinition.Category
1. MetricDefinition.DisplayDescription
1. MetricDefinition.MetricClass
1. MetricNamespace.Classification
1. MetricTrigger.DividePerInstance
1. MetricTrigger.MetricResourceLocation
1. Resource.Etag
1. Resource.Kind
1. RuleDataSource.LegacyResourceID
1. RuleDataSource.MetricNamespace
1. RuleDataSource.ResourceLocation
1. RuleManagementEventDataSource.LegacyResourceID
1. RuleManagementEventDataSource.MetricNamespace
1. RuleManagementEventDataSource.ResourceLocation
1. RuleMetricDataSource.LegacyResourceID
1. RuleMetricDataSource.MetricNamespace
1. RuleMetricDataSource.ResourceLocation
1. TimeSeriesBaseline.MetadataValues
