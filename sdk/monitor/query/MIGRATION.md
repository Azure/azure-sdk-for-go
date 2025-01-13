# Guide to migrate from `azquery` to `azlogs`, `azmetrics`, and `armmonitor`

With the service team's creation of a new data plane metrics endpoint, the Go SDK team decided to reimagine our Monitor story. After careful deliberation, we decided to restructure Monitor Query to more accurately reflect the [REST API][rest_api].

Logs largely remained the same; the code just moved from `azquery` to [`azlogs`][azlogs].

Metrics has a larger split. Users wanting to use the existing ARM APIs should use the [`armmonitor`][armmonitor] module. `armmonitor` contains all the metrics functionality from `azquery` and more.

[`azmetrics`][azmetrics] contains a [new data plane metrics feature][azmetrics_blog] that wasn't previously available in `azquery`. This allows users to query data from one or multiple resource IDs using a data plane endpoint. `azmetrics` will expand in the future when the service adds more features.

## Name changes

| Old module   | Old method name |New module | New method name | 
| ----------- | ----------- | --- | --- |
| `azquery` | `LogsClient.QueryWorkspace` | `azlogs` | `Client.QueryWorkspace` |
|  | `LogsClient.QueryResource` |  | `Client.QueryResource` |
| | `LogsClient.QueryBatch` | | N/A |
| | `MetricsClient.QueryResource` | `armmonitor` | `MetricsClient.List` |
| | `MetricsClient.NewListDefinitionsPager`  | | `MetricDefinitionsClient.NewListPager` |
| | `MetricsClient.NewListNamespacesPager`  | | `MetricNamespacesClient.NewListPager` |

## Query Logs

### `azquery`

```go
import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func main() {
	// create the logs client
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}

	// execute the logs query
	res, err := client.QueryWorkspace(context.TODO(), "workspaceID",
		azquery.Body{
			Query:    to.Ptr("<kusto query>"),
			Timespan: to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
		},
		nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}
}
```

### `azlogs`

The logs code for `azlogs` and `azquery` are very similiar. 

```go
import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs"
)

func main() {
	// create the client
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azlogs.NewClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}

	// execute the logs query
	res, err := client.QueryWorkspace(
		context.TODO(),
		"<workspaceID>",
		azlogs.QueryBody{
			Query:    to.Ptr("<kusto query>"), // example Kusto query
			Timespan: to.Ptr(azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
		},
		nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}
}
```

## Query Metrics

### `azquery`

```go
import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func main() {
	// create the metrics client
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}

	// execute the metrics query
	res, err := client.QueryResource(context.TODO(), "<resourceID>",
		&azquery.MetricsClientQueryResourceOptions{
			Timespan:        to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
			Interval:        to.Ptr("PT1M"),
			MetricNames:     nil,
			Aggregation:     to.SliceOfPtrs(azquery.AggregationTypeAverage, azquery.AggregationTypeCount),
			Top:             to.Ptr[int32](3),
			OrderBy:         to.Ptr("Average asc"),
			Filter:          to.Ptr("BlobType eq '*'"),
			ResultType:      nil,
			MetricNamespace: to.Ptr("Microsoft.Storage/storageAccounts/blobServices"),
		})
	if err != nil {
		//TODO: handle error
	}
	_ = res
}
```

### `armmonitor`

The code in `armmonitor` is closer to the REST API; therefore, there are some name changes. Additionally, `azquery.MetricsClient` is separated into [different clients](#name-changes).

```go
import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
)

func main() {
	// create the client
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := armmonitor.NewMetricsClient("<subscription-id>", cred, nil)
	if err != nil {
		//TODO: handle error
	}

	// execute the query
	res, err := client.List(context.Background(), "<resourceID>", &armmonitor.MetricsClientListOptions{
		Timespan:            to.Ptr("2021-04-20T09:00:00.000Z/2021-04-20T14:00:00.000Z"),
		Interval:            to.Ptr("PT6H"),
		Metricnames:         to.Ptr("BlobCount,BlobCapacity"),
		Aggregation:         to.Ptr("average,minimum,maximum"),
		Top:                 to.Ptr[int32](5),
		Orderby:             to.Ptr("average asc"),
		Filter:              to.Ptr("Tier eq '*'"),
		ResultType:          nil,
		Metricnamespace:     to.Ptr("Microsoft.Storage/storageAccounts/blobServices"),
		AutoAdjustTimegrain: to.Ptr(true),
		ValidateDimensions:  to.Ptr(false),
	})
	if err != nil {
		//TODO: handle partial error
	}
	_ = res
}
```

<!-- LINKS -->
[armmonitor]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor
[azlogs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs
[azmetrics]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics
[azmetrics_blog]: https://devblogs.microsoft.com/azure-sdk/multi-resource-metrics-query-support-in-the-azure-monitor-query-libraries/
[github_issues]: https://github.com/Azure/azure-sdk-for-go/issues
[rest_api]: https://learn.microsoft.com/rest/api/monitor/