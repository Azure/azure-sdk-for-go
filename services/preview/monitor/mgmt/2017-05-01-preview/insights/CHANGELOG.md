# Unreleased

## Breaking Changes

### Removed Constants

1. Sensitivity.High
1. Sensitivity.Low
1. Sensitivity.Medium

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

### Struct Changes

#### Removed Structs

1. Baseline
1. BaselineMetadataValue
1. BaselineProperties
1. BaselineResponse
1. CalculateBaselineResponse
1. MetricBaselineClient
1. TimeSeriesInformation

## Additive Changes

### New Constants

1. MetricStatisticType.MetricStatisticTypeCount
1. ScaleType.ServiceAllowedNextValue
1. Unit.UnitBitsPerSecond
1. Unit.UnitCores
1. Unit.UnitMilliCores
1. Unit.UnitNanoCores

### Struct Changes

#### New Struct Fields

1. AlertRule.Action
1. AlertRule.ProvisioningState
1. AutoscaleSetting.TargetResourceLocation
1. Metric.DisplayDescription
1. Metric.ErrorCode
1. Metric.ErrorMessage
1. MetricDefinition.Category
1. MetricDefinition.DisplayDescription
1. MetricTrigger.DividePerInstance
1. MetricTrigger.MetricResourceLocation
1. RuleDataSource.LegacyResourceID
1. RuleDataSource.MetricNamespace
1. RuleDataSource.ResourceLocation
1. RuleManagementEventDataSource.LegacyResourceID
1. RuleManagementEventDataSource.MetricNamespace
1. RuleManagementEventDataSource.ResourceLocation
1. RuleMetricDataSource.LegacyResourceID
1. RuleMetricDataSource.MetricNamespace
1. RuleMetricDataSource.ResourceLocation
