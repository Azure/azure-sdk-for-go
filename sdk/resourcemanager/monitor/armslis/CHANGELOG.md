# Release History

## 0.1.0 (2026-04-22)

### Features Added

- Initial preview release of `armslis` for managing Service Level Indicator (SLI) resources under the `Microsoft.Monitor` namespace.
- Support for SLI resource CRUD operations: create or update, get, delete, and list.
- SLI evaluation with Availability and Latency categories, supporting both window-based and request-based evaluation types with configurable signal sources, aggregation, and SLO baselines.
- Integration with Azure Monitor Workspace (AMW) accounts for metric emission, with managed identity and alert support.

### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armslis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
