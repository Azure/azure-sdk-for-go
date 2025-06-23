//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"

	"github.com/stretchr/testify/require"
)

// MachineLearningServices events
func TestConsumeCloudEventMachineLearningServicesModelRegisteredEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\", \"source\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"type\":\"Microsoft.MachineLearningServices.ModelRegistered\",\"source\":\"models/sklearn_regression_model:3\",\"time\":\"2019-10-17T22:23:57.5350054+00:00\",\"id\":\"3b73ee51-bbf4-480d-9112-cfc23b41bfdb\",\"data\":{\"modelName\":\"sklearn_regression_model\",\"modelVersion\":\"3\",\"modelTags\":{\"area\":\"diabetes\",\"type\":\"regression\"},\"modelProperties\":{\"area\":\"test\"}}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesModelRegisteredEventData](t, events[0].Data)
	require.Equal(t, "sklearn_regression_model", *sysEvent.ModelName)
	require.Equal(t, "3", *sysEvent.ModelVersion)
	require.Equal(t, "regression", sysEvent.ModelTags["type"])
	require.Equal(t, "test", sysEvent.ModelProperties["area"])
}

func TestConsumeCloudEventMachineLearningServicesModelDeployedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\", \"source\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"type\":\"Microsoft.MachineLearningServices.ModelDeployed\",\"subject\":\"endpoints/aciservice1\",\"time\":\"2019-10-23T18:20:08.8824474+00:00\",\"id\":\"40d0b167-be44-477b-9d23-a2befba7cde0\",\"data\":{\"serviceName\":\"aciservice1\",\"serviceComputeType\":\"ACI\",\"serviceTags\":{\"mytag\":\"test tag\"},\"serviceProperties\":{\"myprop\":\"test property\"},\"modelIds\":\"my_first_model:1,my_second_model:1\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesModelDeployedEventData](t, events[0].Data)
	require.Equal(t, "aciservice1", *sysEvent.ServiceName)
	require.Equal(t, 2, len(strings.Split(*sysEvent.ModelIDs, ",")))
}

func TestConsumeCloudEventMachineLearningServicesRunCompletedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\", \"source\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"type\":\"Microsoft.MachineLearningServices.RunCompleted\",\"subject\":\"experiments/0fa9dfaa-cba3-4fa7-b590-23e48548f5c1/runs/AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"time\":\"2019-10-18T19:29:55.8856038+00:00\",\"id\":\"044ac44d-462c-4043-99eb-d9e01dc760ab\",\"data\":{\"experimentId\":\"0fa9dfaa-cba3-4fa7-b590-23e48548f5c1\",\"experimentName\":\"automl-local-regression\",\"runId\":\"AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"runType\":\"automl\",\"RunTags\":{\"experiment_status\":\"ModelSelection\",\"experiment_status_descr\":\"Beginning model selection.\"},\"runProperties\":{\"num_iterations\":\"10\",\"target\":\"local\"}}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesRunCompletedEventData](t, events[0].Data)
	require.Equal(t, "AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc", *sysEvent.RunID)
	require.Equal(t, "automl-local-regression", *sysEvent.ExperimentName)
}

func TestConsumeCloudEventMachineLearningServicesRunStatusChangedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\", \"source\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"type\":\"Microsoft.MachineLearningServices.RunStatusChanged\",\"subject\":\"experiments/0fa9dfaa-cba3-4fa7-b590-23e48548f5c1/runs/AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"time\":\"2020-03-09T23:53:04.4579724Z\",\"id\":\"aa8cd7df-fe28-5d5d-9b40-3342dbc2a887\",\"data\":{\"runStatus\": \"Running\",\"experimentId\":\"0fa9dfaa-cba3-4fa7-b590-23e48548f5c1\",\"experimentName\":\"automl-local-regression\",\"runId\":\"AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"runType\":\"automl\",\"runTags\":{\"experiment_status\":\"ModelSelection\",\"experiment_status_descr\":\"Beginning model selection.\"},\"runProperties\":{\"num_iterations\":\"10\",\"target\":\"local\"}}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesRunStatusChangedEventData](t, events[0].Data)
	require.Equal(t, "AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc", *sysEvent.RunID)
	require.Equal(t, "automl-local-regression", *sysEvent.ExperimentName)
	require.Equal(t, "Running", *sysEvent.RunStatus)
	require.Equal(t, "automl", *sysEvent.RunType)
}

func TestConsumeCloudEventMachineLearningServicesDatasetDriftDetectedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\", \"source\":\"/subscriptions/60582a10-b9fd-49f1-a546-c4194134bba8/resourceGroups/copetersRG/providers/Microsoft.MachineLearningServices/workspaces/driftDemoWS\",\"type\":\"Microsoft.MachineLearningServices.DatasetDriftDetected\",\"subject\":\"datadrift/01d29aa4-e6a4-470a-9ef3-66660d21f8ef/run/01d29aa4-e6a4-470a-9ef3-66660d21f8ef_1571590300380\",\"time\":\"2019-10-20T17:08:08.467191+00:00\",\"id\":\"2684de79-b145-4dcf-ad2e-6a1db798585f\",\"data\":{\"dataDriftId\":\"01d29aa4-e6a4-470a-9ef3-66660d21f8ef\",\"dataDriftName\":\"copetersDriftMonitor3\",\"runId\":\"01d29aa4-e6a4-470a-9ef3-66660d21f8ef_1571590300380\",\"baseDatasetId\":\"3c56d136-0f64-4657-a0e8-5162089a88a3\",\"tarAsSystemEventDatasetId\":\"d7e74d2e-c972-4266-b5fb-6c9c182d2a74\",\"driftCoefficient\":0.8350349068479208,\"startTime\":\"2019-07-04T00:00:00+00:00\",\"endTime\":\"2019-07-05T00:00:00+00:00\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesDatasetDriftDetectedEventData](t, events[0].Data)
	require.Equal(t, "copetersDriftMonitor3", *sysEvent.DataDriftName)
}

// Maps events
func TestConsumeCloudEventMapsGeofenceEnteredEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.Maps.GeofenceEntered\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceEnteredEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

func TestConsumeCloudEventMapsGeofenceExitedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.Maps.GeofenceExited\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceExitedEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

func TestConsumeCloudEventMapsGeofenceResultEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.Maps.GeofenceResult\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceResultEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

// Resource Manager (Azure Subscription/Resource Group) events
func TestConsumeCloudEventResourceWriteSuccessEvent(t *testing.T) {
	requestContent := `[   {     "source":"/subscriptions/{subscription-id}",     "subject":"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "type":"Microsoft.Resources.ResourceWriteSuccess",    "time":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": {        "authorization":{},        "claims":{},        "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",        "httpRequest":{},        "resourceProvider":"Microsoft.EventGrid",        "resourceUri":"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",        "status":"Succeeded",        "subscriptionId":"{subscription-id}",        "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },    "specversion": "1.0"  }]`

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceWriteSuccessEventData](t, events[0].Data)

	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceWriteFailureEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceWriteFailure\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },    \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.ResourceWriteFailureEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceWriteCancelEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceWriteCancel\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },    \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceWriteCancelEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceDeleteSuccessEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceDeleteSuccess\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceDeleteSuccessEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceDeleteFailureEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceDeleteFailure\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceDeleteFailureEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceDeleteCancelEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceDeleteCancel\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceDeleteCancelEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceActionSuccessEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceActionSuccess\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceActionSuccessEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceActionFailureEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceActionFailure\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceActionFailureEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

func TestConsumeCloudEventResourceActionCancelEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.Resources.ResourceActionCancel\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": {        \"authorization\":{},        \"claims\":{},        \"correlationId\":\"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6\",        \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",        \"operationName\":\"Microsoft.EventGrid/eventSubscriptions/write\",        \"status\":\"Succeeded\",        \"subscriptionId\":\"{subscription-id}\",        \"tenantId\":\"72f988bf-86f1-41af-91ab-2d7cd011db47\"        },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceActionCancelEventData](t, events[0].Data)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", *sysEvent.TenantID)
}

// ServiceBus events
func TestConsumeCloudEventServiceBusActiveMessagesAvailableWithNoListenersEvent(t *testing.T) {
	requestContent := "[{  \"source\": \"/subscriptions/id/resourcegroups/rg/providers/Microsoft.ServiceBus/namespaces/testns1\",  \"subject\": \"topics/topic1/subscriptions/sub1\",  \"type\": \"Microsoft.ServiceBus.ActiveMessagesAvailableWithNoListeners\",  \"time\": \"2018-02-14T05:12:53.4133526Z\",  \"id\": \"dede87b0-3656-419c-acaf-70c95ddc60f5\",  \"data\": {    \"namespaceName\": \"testns1\",    \"requestUri\": \"https://testns1.servicebus.windows.net/t1/subscriptions/sub1/messages/head\",    \"entityType\": \"subscriber\",    \"queueName\": \"queue1\",    \"topicName\": \"topic1\",    \"subscriptionName\": \"sub1\"  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ServiceBusActiveMessagesAvailableWithNoListenersEventData](t, events[0].Data)
	require.Equal(t, "testns1", *sysEvent.NamespaceName)
}

func TestConsumeCloudEventServiceBusDeadletterMessagesAvailableWithNoListenersEvent(t *testing.T) {
	requestContent := "[{  \"source\": \"/subscriptions/id/resourcegroups/rg/providers/Microsoft.ServiceBus/namespaces/testns1\",  \"subject\": \"topics/topic1/subscriptions/sub1\",  \"type\": \"Microsoft.ServiceBus.DeadletterMessagesAvailableWithNoListeners\",  \"time\": \"2018-02-14T05:12:53.4133526Z\",  \"id\": \"dede87b0-3656-419c-acaf-70c95ddc60f5\",  \"data\": {    \"namespaceName\": \"testns1\",    \"requestUri\": \"https://testns1.servicebus.windows.net/t1/subscriptions/sub1/messages/head\",    \"entityType\": \"subscriber\",    \"queueName\": \"queue1\",    \"topicName\": \"topic1\",    \"subscriptionName\": \"sub1\"  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ServiceBusDeadletterMessagesAvailableWithNoListenersEventData](t, events[0].Data)
	require.Equal(t, "testns1", *sysEvent.NamespaceName)
}

// Storage events
func TestConsumeCloudEventStorageBlobCreatedEvent(t *testing.T) {
	requestContent := "[ {  \"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/file1.txt\",  \"type\": \"Microsoft.Storage.BlobCreated\",  \"time\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"PutBlockList\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"contentType\": \"text/plain\",    \"contentLength\": 447,    \"blobType\": \"BlockBlob\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/file1.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\" },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobCreatedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/file1.txt", *sysEvent.URL)
}

func TestConsumeCloudEventStorageBlobDeletedEvent(t *testing.T) {
	requestContent := "[{   \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"type\": \"Microsoft.Storage.BlobDeleted\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.URL)
}

func TestConsumeCloudEventStorageBlobRenamedEvent(t *testing.T) {
	requestContent := "[ {  \"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"type\": \"Microsoft.Storage.BlobRenamed\",  \"time\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"RenameFile\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"destinationUrl\": \"https://myaccount.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobRenamedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.DestinationURL)
}

func TestConsumeCloudEventStorageDirectoryCreatedEvent(t *testing.T) {
	requestContent := "[ {  \"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"type\": \"Microsoft.Storage.DirectoryCreated\",  \"time\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"CreateDirectory\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryCreatedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeCloudEventStorageDirectoryDeletedEvent(t *testing.T) {
	requestContent := "[{   \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"type\": \"Microsoft.Storage.DirectoryDeleted\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },   \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryDeletedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeCloudEventStorageDirectoryRenamedEvent(t *testing.T) {
	requestContent := "[{   \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"type\": \"Microsoft.Storage.DirectoryRenamed\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"destinationUrl\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryRenamedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.DestinationURL)
}

func TestConsumeCloudEventStorageAsyncOperationInitiatedEvent(t *testing.T) {
	requestContent := "[{   \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"type\": \"Microsoft.Storage.AsyncOperationInitiated\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"specversion\": \"1.0\"}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageAsyncOperationInitiatedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeCloudEventStorageBlobTierChangedEvent(t *testing.T) {
	requestContent := "[{   \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"type\": \"Microsoft.Storage.BlobTierChanged\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"specversion\": \"1.0\"}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobTierChangedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

// App Service events
func TestConsumeCloudEventWebAppUpdatedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.AppUpdated\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebAppUpdatedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebBackupOperationStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.BackupOperationStarted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebBackupOperationCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.BackupOperationCompleted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebBackupOperationFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.BackupOperationFailed\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)
	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebRestoreOperationStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.RestoreOperationStarted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebRestoreOperationCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.RestoreOperationCompleted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebRestoreOperationFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.RestoreOperationFailed\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebSlotSwapStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.SlotSwapStarted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebSlotSwapCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"specversion\": \"1.0\", \"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.SlotSwapCompleted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"}}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebSlotSwapFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.SlotSwapFailed\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},   \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebSlotSwapWithPreviewStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.SlotSwapWithPreviewStarted\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapWithPreviewStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebSlotSwapWithPreviewCancelledEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"type\": \"Microsoft.Web.SlotSwapWithPreviewCancelled\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},   \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapWithPreviewCancelledEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeCloudEventWebAppServicePlanUpdatedEvent(t *testing.T) {
	planName := "testPlan01"
	requestContent := "[{\"source\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/serverfarms/testPlan01\", \"subject\": \"/Microsoft.Web/serverfarms/testPlan01\",\"type\": \"Microsoft.Web.AppServicePlanUpdated\", \"time\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appServicePlanEventTypeDetail\": { \"stampKind\": \"Public\",\"action\": \"Updated\",\"status\": \"Started\" },\"name\": \"testPlan01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"specversion\": \"1.0\",\"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebAppServicePlanUpdatedEventData](t, events[0].Data)
	require.Equal(t, planName, *sysEvent.Name)
}

// Policy Insights
func TestConsumeCloudEventPolicyInsightsPolicyStateChangedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.PolicyInsights.PolicyStateChanged\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateChangedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

func TestConsumeCloudEventPolicyInsightsPolicyStateCreatedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.PolicyInsights.PolicyStateCreated\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateCreatedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

func TestConsumeCloudEventPolicyInsightsPolicyStateDeletedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"type\":\"Microsoft.PolicyInsights.PolicyStateDeleted\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":{},        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateDeletedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

// Communication events
func TestConsumeCloudEventAcsRecordingFileStatusUpdatedEventData(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/recording/call/{call-id}/recordingId/{recording-id}\",    \"type\":\"Microsoft.Communication.RecordingFileStatusUpdated\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"recordingStorageInfo\": { \"recordingChunks\": [ { \"documentId\": \"0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8\", \"index\": 0, \"endReason\": \"SessionEnded\", \"contentLocation\": \"https://storage.asm.skype.com/v1/objects/0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8/content/video\", \"metadataLocation\": \"https://storage.asm.skype.com/v1/objects/0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8/content/acsmetadata\" }]}, \"recordingChannelType\": \"Mixed\", \"recordingContentType\": \"Audio\", \"recordingFormatType\": \"Mp3\"},   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ACSRecordingFileStatusUpdatedEventData](t, events[0].Data)

	// back compat
	require.Equal(t, azsystemevents.ACSRecordingChannelTypeMixed, *sysEvent.RecordingChannelType)
	require.Equal(t, azsystemevents.ACSRecordingContentTypeAudio, *sysEvent.RecordingContentType)
	require.Equal(t, azsystemevents.ACSRecordingFormatTypeMP3, *sysEvent.RecordingFormatType)
}

func TestConsumeCloudEventAcsEmailDeliveryReportReceivedEvent(t *testing.T) {
	requestContent := `{
		"id": "5f04f77c-2a6a-43bd-9b74-576a64c01f9e",
		"source": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "sender/test2@contoso.org/message/950850f5-bcdf-4315-b77a-6447cf56fac9",
		"data": {
			"sender": "test2@contoso.org",
			"recipient": "test1@contoso.com",
			"messageId": "950850f5-bcdf-4315-b77a-6447cf56fac9",
			"status": "Delivered",
			"deliveryStatusDetails": {
				"statusMessage": "DestinationMailboxFull"
			},
			"deliveryAttemptTimestamp": "2023-02-09T19:46:12.2480265+00:00"
		},
		"type": "Microsoft.Communication.EmailDeliveryReportReceived",
		"time": "2023-02-09T19:46:12.2478002Z",
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ACSEmailDeliveryReportReceivedEventData](t, event.Data)
	require.Equal(t, "test2@contoso.org", *sysEvent.Sender)
	require.Equal(t, "test1@contoso.com", *sysEvent.Recipient)
	require.Equal(t, azsystemevents.ACSEmailDeliveryReportStatusDelivered, *sysEvent.Status)
	require.Equal(t, "DestinationMailboxFull", *sysEvent.DeliveryStatusDetails.StatusMessage)
	require.Equal(t, mustParseTime(t, "2023-02-09T19:46:12.2480265+00:00"), *sysEvent.DeliveryAttemptTimestamp)
}

func TestConsumeCloudEventAcsIncomingCallEvent(t *testing.T) {
	requestContent := `{
		"id": "e80026e7-e298-46ba-bc42-dab0eda92581",
		"source": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "/caller/{caller-id}/recipient/{recipient-id}",
		"data": {
			"to": {
				"kind": "communicationUser",
				"rawId": "{recipient-id}",
				"communicationUser": {
					"id": "{recipient-id}"
				}
			},
			"from": {
				"kind": "communicationUser",
				"rawId": "{caller-id}",
				"communicationUser": {
					"id": "{caller-id}"
				}
			},
			"serverCallId": "{server-call-id}",
			"callerDisplayName": "VOIP Caller",
			"customContext": {
				"sipHeaders": {
					"userToUser": "616d617a6f6e5f6368696;encoding=hex",
					"X-MS-Custom-myheader1": "35567842",
					"X-MS-Custom-myheader2": "customsipheadervalue"
				},
				"voipHeaders": {
					"customHeader": "customValue"
				}
			},
			"incomingCallContext": "{incoming-call-contextValue}",
			"correlationId": "correlationId"
		},
		"type": "Microsoft.Communication.IncomingCall",
		"specversion": "1.0",
		"time": "2023-04-04T17:18:42.5542219Z"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ACSIncomingCallEventData](t, event.Data)

	require.Equal(t, "{recipient-id}", *sysEvent.ToCommunicationIdentifier.CommunicationUser.ID)
	require.Equal(t, "{caller-id}", *sysEvent.FromCommunicationIdentifier.CommunicationUser.ID)
	require.Equal(t, "VOIP Caller", *sysEvent.CallerDisplayName)
	require.Equal(t, "616d617a6f6e5f6368696;encoding=hex", *sysEvent.CustomContext.SipHeaders["userToUser"])
	require.Equal(t, "35567842", *sysEvent.CustomContext.SipHeaders["X-MS-Custom-myheader1"])
	require.Equal(t, "customsipheadervalue", *sysEvent.CustomContext.SipHeaders["X-MS-Custom-myheader2"])
	require.Equal(t, "customValue", *sysEvent.CustomContext.VoipHeaders["customHeader"])
	require.Equal(t, "{incoming-call-contextValue}", *sysEvent.IncomingCallContext)
	require.Equal(t, "correlationId", *sysEvent.CorrelationID)
}

func TestConsumeCloudEventAcsRouterJobClassificationFailedEvent(t *testing.T) {
	requestContent := `{
		"id": "e80026e7-e298-46ba-bc42-dab0eda92581",
		"source": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "job/{job-id}/channel/{channel-id}/classificationpolicy/{classificationpolicy-id}",
			"data": {
			"errors": [
				{
				"code": "Failure",
				"message": "Classification failed due to <reason>",
				"target": null,
				"innererror": {
								"code": "InnerFailure",
								"message": "Classification failed due to <reason>",
								"target": null},
				"details": null
				}
			],
			"jobId": "7f1df17b-570b-4ae5-9cf5-fe6ff64cc712",
			"channelReference": "test-abc",
			"channelId": "FooVoiceChannelId",
			"classificationPolicyId": "test-policy",
			"labels": {
				"Locale": "en-us",
				"Segment": "Enterprise",
				"Token": "FooToken"
			},
			"tags": {
				"Locale": "en-us",
				"Segment": "Enterprise",
				"Token": "FooToken"
			}
			},
			"type": "Microsoft.Communication.RouterJobClassificationFailed",
			"specversion": "1.0",
			"time": "2022-02-17T00:55:25.1736293Z"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)

	sysEvent := deserializeSystemEvent[azsystemevents.ACSRouterJobClassificationFailedEventData](t, event.Data)

	var errors = sysEvent.Errors
	require.Equal(t, 1, len(errors))
	require.Equal(t, "Failure", errors[0].Code)
	require.Equal(t, "Code: Failure\n"+
		"Message: Classification failed due to <reason>\n"+
		"InnerError:\n"+
		"  Code: InnerFailure\n"+
		"  Message: Classification failed due to <reason>\n", (*errors[0]).Error())
}
