# Release History

## 0.2.0 (2026-05-23)
### Breaking Changes

- `SamplingTypeAvg` from enum `SamplingType` has been removed and replaced by `SamplingTypeAverage`
- Wire values for existing enum `SamplingType` values have been re-cased: `SamplingTypeMax` is now `Max` (previously `max`), `SamplingTypeMin` is now `Min` (previously `min`), and `SamplingTypeSum` is now `Sum` (previously `sum`)
- Wire values for existing enum `ConditionOperator` values have changed: `ConditionOperatorEqual` is now `eq` (previously `==`), `ConditionOperatorNotEqual` is now `ne` (previously `!=`), `ConditionOperatorGreaterThan` is now `gt` (previously `>`), `ConditionOperatorGreaterThanOrEqual` is now `gte` (previously `>=`), `ConditionOperatorLessThan` is now `lt` (previously `<`), `ConditionOperatorLessThanOrEqual` is now `lte` (previously `<=`), `ConditionOperatorIn` is now `in` (previously `@in`), `ConditionOperatorNotIn` is now `notin` (previously `!in`), `ConditionOperatorNotContains` is now `notcontains` (previously `!contains`), and `ConditionOperatorNotStartsWith` is now `notstartswith` (previously `!startswith`)
- Wire values for existing enum `WindowUptimeCriteriaComparator` values have changed: `WindowUptimeCriteriaComparatorGreaterThan` is now `gt` (previously `>`), `WindowUptimeCriteriaComparatorGreaterThanOrEqual` is now `gte` (previously `>=`), `WindowUptimeCriteriaComparatorLessThan` is now `lt` (previously `<`), and `WindowUptimeCriteriaComparatorLessThanOrEqual` is now `lte` (previously `<=`)

### Features Added

- New values `SamplingTypeAverage` and `SamplingTypeCount` added to enum type `SamplingType`


## 0.1.0 (2026-04-22)

### Features Added

- Initial preview release of `armslis` for managing Service Level Indicator (SLI) resources under the `Microsoft.Monitor` namespace.
- Support for SLI resource CRUD operations: create or update, get, delete, and list.
- SLI evaluation with Availability and Latency categories, supporting both window-based and request-based evaluation types with configurable signal sources, aggregation, and SLO baselines.
- Integration with Azure Monitor Workspace (AMW) accounts for metric emission, with managed identity and alert support.

### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armslis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
