//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents_test

import (
	"azsystemevents"
	"testing"

	"github.com/stretchr/testify/require"
)

// AppConfiguration events
func TestConsumeCloudEventAppConfigurationKeyValueDeletedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.AppConfiguration.KeyValueDeleted\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"key\":\"key1\",\"label\":\"label1\",\"etag\":\"etag1\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationKeyValueDeletedEventData](t, events[0].Data)
	require.Equal(t, "key1", *sysEvent.Key)
}

func TestConsumeCloudEventAppConfigurationKeyValueModifiedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.AppConfiguration.KeyValueModified\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"key\":\"key1\",\"label\":\"label1\",\"etag\":\"etag1\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationKeyValueModifiedEventData](t, events[0].Data)
	require.Equal(t, "key1", *sysEvent.Key)
}

func TestConsumeCloudEventAppConfigurationSnapshotCreatedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.AppConfiguration.SnapshotCreated\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"name\":\"Foo\",\"etag\":\"FnUExLaj2moIi4tJX9AXn9sakm0\",\"syncToken\":\"zAJw6V16=Njo1IzUxNjQ2NzM=;sn=5164673\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationSnapshotCreatedEventData](t, events[0].Data)
	require.Equal(t, "Foo", *sysEvent.Name)
	require.Equal(t, "zAJw6V16=Njo1IzUxNjQ2NzM=;sn=5164673", *sysEvent.SyncToken)
}

func TestConsumeCloudEventAppConfigurationSnapshotModifiedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.AppConfiguration.SnapshotModified\",\"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"name\":\"Foo\",\"etag\":\"FnUExLaj2moIi4tJX9AXn9sakm0\",\"syncToken\":\"zAJw6V16=Njo1IzUxNjQ2NzM=;sn=5164673\"}}]"
	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationSnapshotModifiedEventData](t, events[0].Data)
	require.Equal(t, "Foo", *sysEvent.Name)
	require.Equal(t, "zAJw6V16=Njo1IzUxNjQ2NzM=;sn=5164673", *sysEvent.SyncToken)
}

// ContainerRegistry events
func TestConsumeCloudEventContainerRegistryImagePushedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.ContainerRegistry.ImagePushed\",  \"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"eventID\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"testaction\",\"target\":{\"mediaType\":\"test\",\"size\":20,\"digest\":\"digest1\",\"length\":20,\"repository\":\"test\",\"url\":\"url1\",\"tag\":\"test\"},\"request\":{\"id\":\"id\",\"addr\":\"127.0.0.1\",\"host\":\"test\",\"method\":\"method1\",\"useragent\":\"useragent1\"},\"actor\":{\"name\":\"testactor\"},\"source\":{\"addr\":\"127.0.0.1\",\"instanceID\":\"id\"}}}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryImagePushedEventData](t, events[0].Data)
	require.Equal(t, "127.0.0.1", *sysEvent.Request.Addr)
}

func TestConsumeCloudEventContainerRegistryImageDeletedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.ContainerRegistry.ImageDeleted\",  \"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"eventID\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"testaction\",\"target\":{\"mediaType\":\"test\",\"size\":20,\"digest\":\"digest1\",\"length\":20,\"repository\":\"test\",\"url\":\"url1\",\"tag\":\"test\"},\"request\":{\"id\":\"id\",\"addr\":\"127.0.0.1\",\"host\":\"test\",\"method\":\"method1\",\"useragent\":\"useragent1\"},\"actor\":{\"name\":\"testactor\"},\"source\":{\"addr\":\"127.0.0.1\",\"instanceID\":\"id\"}}}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryImageDeletedEventData](t, events[0].Data)
	require.Equal(t, "testactor", *sysEvent.Actor.Name)
}

func TestConsumeCloudEventContainerRegistryChartDeletedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.ContainerRegistry.ChartDeleted\",  \"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"id\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"action1\",\"target\":{\"mediaType\":\"mediatype1\",\"size\":20,\"digest\":\"digest1\",\"repository\":null,\"tag\":null,\"name\":\"name1\",\"version\":null}}}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryChartDeletedEventData](t, events[0].Data)
	require.Equal(t, "mediatype1", *sysEvent.Target.MediaType)
}

func TestConsumeCloudEventContainerRegistryChartPushedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"source\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"type\": \"Microsoft.ContainerRegistry.ChartPushed\",  \"time\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"id\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"action1\",\"target\":{\"mediaType\":\"mediatype1\",\"size\":40,\"digest\":\"digest1\",\"repository\":null,\"tag\":null,\"name\":\"name1\",\"version\":null}}}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryChartPushedEventData](t, events[0].Data)
	require.Equal(t, "mediatype1", *sysEvent.Target.MediaType)
}

// Container service events
func TestConsumeCloudEventContainerServiceSupportEndedEvent(t *testing.T) {
	requestContent := `
	{
		"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"type": "Microsoft.ContainerService.ClusterSupportEnded",
		"time": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"kubernetesVersion": "1.23.15"
		},
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent.KubernetesVersion)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEndedEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent2.KubernetesVersion)
}

func TestConsumeCloudEventContainerServiceSupportEndingEvent(t *testing.T) {
	requestContent := `
	{
		"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"type": "Microsoft.ContainerService.ClusterSupportEnding",
		"time": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"kubernetesVersion": "1.23.15"
		},
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent.KubernetesVersion)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEndingEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent2.KubernetesVersion)
}

func TestConsumeCloudEventContainerServiceNodePoolRollingFailed(t *testing.T) {
	requestContent := `
	{
		"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"type": "Microsoft.ContainerService.NodePoolRollingFailed",
		"time": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingFailedEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent2.NodePoolName)
}

func TestConsumeCloudEventContainerServiceNodePoolRollingStarted(t *testing.T) {
	requestContent := `
	{
		"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"type": "Microsoft.ContainerService.NodePoolRollingStarted",
		"time": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)
}

func TestConsumeCloudEventContainerServiceNodePoolRollingSucceeded(t *testing.T) {
	requestContent := `
	{
		"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"type": "Microsoft.ContainerService.NodePoolRollingSucceeded",
		"time": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingSucceededEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent2.NodePoolName)
}

// IoTHub Device events
func TestConsumeCloudEventIoTHubDeviceCreatedEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",  \"id\": \"2da5e9b4-4e38-04c1-cc58-9da0b37230c0\", \"source\": \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\", \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\", \"type\": \"Microsoft.Devices.DeviceCreated\", \"time\": \"2018-07-03T23:20:07.6532054Z\",    \"data\": {      \"twin\": {        \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",        \"etag\": \"AAAAAAAAAAE=\",        \"deviceEtag\": null,        \"status\": \"enabled\",        \"statusUpdateTime\": \"0001-01-01T00:00:00\",        \"connectionState\": \"Disconnected\",        \"lastActivityTime\": \"0001-01-01T00:00:00\",        \"cloudToDeviceMessageCount\": 0,        \"authenticationType\": \"sas\",        \"x509Thumbprint\": {          \"primaryThumbprint\": null,          \"secondaryThumbprint\": null        },        \"version\": 2,        \"properties\": {          \"desired\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          },          \"reported\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          }        }      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\"    }}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IotHubDeviceCreatedEventData](t, events[0].Data)
	require.Equal(t, "enabled", *sysEvent.Twin.Status)
}

func TestConsumeCloudEventIoTHubDeviceDeletedEvent(t *testing.T) {
	requestContent := "[  {\"specversion\": \"1.0\",     \"id\": \"aaaf95c6-ed99-b307-e321-81d8e4f731a6\",    \"source\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"type\": \"Microsoft.Devices.DeviceDeleted\",    \"time\": \"2018-07-03T23:21:33.2753956Z\",    \"data\": {      \"twin\": {        \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",        \"etag\": \"AAAAAAAAAAI=\",        \"deviceEtag\": null,        \"status\": \"enabled\",        \"statusUpdateTime\": \"0001-01-01T00:00:00\",        \"connectionState\": \"Disconnected\",        \"lastActivityTime\": \"0001-01-01T00:00:00\",        \"cloudToDeviceMessageCount\": 0,        \"authenticationType\": \"sas\",        \"x509Thumbprint\": {          \"primaryThumbprint\": null,          \"secondaryThumbprint\": null        },        \"version\": 3,        \"tags\": {          \"testKey\": \"testValue\"        },        \"properties\": {          \"desired\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          },          \"reported\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          }        }      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\"    }}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IotHubDeviceDeletedEventData](t, events[0].Data)
	require.Equal(t, "AAAAAAAAAAI=", *sysEvent.Twin.Etag)
}

func TestConsumeCloudEventIoTHubDeviceConnectedEvent(t *testing.T) {
	requestContent := "[  {\"specversion\": \"1.0\",     \"id\": \"fbfd8ee1-cf78-74c6-dbcf-e1c58638ccbd\",    \"source\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"type\": \"Microsoft.Devices.DeviceConnected\",    \"time\": \"2018-07-03T23:20:11.6921933+00:00\",    \"data\": {      \"deviceConnectionStateEventInfo\": {        \"sequenceNumber\":          \"000000000000000001D4132452F67CE200000002000000000000000000000001\"      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",      \"moduleId\": \"\"    }}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IotHubDeviceConnectedEventData](t, events[0].Data)
	require.Equal(t, "EGTESTHUB1", *sysEvent.HubName)
}

func TestConsumeCloudEventIoTHubDeviceDisconnectedEvent(t *testing.T) {
	requestContent := "[  { \"specversion\": \"1.0\",    \"id\": \"877f0b10-a086-98ec-27b8-6ae2dfbf5f67\",    \"source\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"type\": \"Microsoft.Devices.DeviceDisconnected\",    \"time\": \"2018-07-03T23:20:52.646434+00:00\",    \"data\": {      \"deviceConnectionStateEventInfo\": {        \"sequenceNumber\":          \"000000000000000001D4132452F67CE200000002000000000000000000000002\"      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",      \"moduleId\": \"\"    }}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IotHubDeviceDisconnectedEventData](t, events[0].Data)
	require.Equal(t, "000000000000000001D4132452F67CE200000002000000000000000000000002", *sysEvent.DeviceConnectionStateEventInfo.SequenceNumber)
}

func TestConsumeCloudEventIoTHubDeviceTelemetryEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"877f0b10-a086-98ec-27b8-6ae2dfbf5f67\",    \"source\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"type\": \"Microsoft.Devices.DeviceTelemetry\",    \"time\": \"2018-07-03T23:20:52.646434+00:00\",    \"data\": { \"body\": { \"Weather\": { \"Temperature\": 900  }, \"Location\": \"USA\"  },  \"properties\": {  \"Status\": \"Active\"  },  \"systemProperties\": { \"iothub-content-type\": \"application/json\", \"iothub-content-encoding\": \"utf-8\"   } }}   ]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IotHubDeviceTelemetryEventData](t, events[0].Data)
	require.Equal(t, "Active", *sysEvent.Properties["Status"])
}

// EventGrid events
func TestConsumeCloudEventEventGridSubscriptionValidationEvent(t *testing.T) {
	requestContent := "[{\"specversion\": \"1.0\",   \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {    \"validationCode\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"validationUrl\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"type\": \"Microsoft.EventGrid.SubscriptionValidationEvent\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.SubscriptionValidationEventData](t, events[0].Data)
	require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", *sysEvent.ValidationCode)
}

func TestConsumeCloudEventEventGridSubscriptionDeletedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {    \"eventSubscriptionId\": \"/subscriptions/id/resourceGroups/rg/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/eventsubscription1\"  },  \"type\": \"Microsoft.EventGrid.SubscriptionDeletedEvent\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.SubscriptionDeletedEventData](t, events[0].Data)
	require.Equal(t, "/subscriptions/id/resourceGroups/rg/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/eventsubscription1", *sysEvent.EventSubscriptionID)
}

func TestConsumeCloudEventEventGridMqttClientCreatedOrUpdatedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {  \"createdOn\": \"2023-07-29T01:14:34.2048108Z\", \"updatedOn\": \"2023-07-29T01:14:34.2048108Z\",\"namespaceName\": \"myns\",\"clientName\": \"client1\",\"clientAuthenticationName\": \"client1\",\"state\": \"Enabled\",\"attributes\": {\"attribute1\": \"value1\"}  },  \"type\": \"Microsoft.EventGrid.MQTTClientCreatedOrUpdated\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventGridMQTTClientCreatedOrUpdatedEventData](t, events[0].Data)
	require.Equal(t, "client1", *sysEvent.ClientName)
	require.Equal(t, "myns", *sysEvent.NamespaceName)
	require.Equal(t, "client1", *sysEvent.ClientAuthenticationName)
}

func TestConsumeCloudEventEventGridMqttClientDeletedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {  \"namespaceName\": \"myns\",\"clientName\": \"client1\",\"clientAuthenticationName\": \"client1\" },  \"type\": \"Microsoft.EventGrid.MQTTClientDeleted\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventGridMQTTClientDeletedEventData](t, events[0].Data)
	require.Equal(t, "client1", *sysEvent.ClientName)
	require.Equal(t, "myns", *sysEvent.NamespaceName)
	require.Equal(t, "client1", *sysEvent.ClientAuthenticationName)
}

func TestConsumeCloudEventEventGridMqttClientSessionConnectedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {  \"namespaceName\": \"myns\",\"clientSessionName\": \"session\",\"clientAuthenticationName\": \"client1\", \"sequenceNumber\": 1 },  \"type\": \"Microsoft.EventGrid.MQTTClientSessionConnected\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventGridMQTTClientSessionConnectedEventData](t, events[0].Data)
	require.Equal(t, "session", *sysEvent.ClientSessionName)
	require.Equal(t, "myns", *sysEvent.NamespaceName)
	require.Equal(t, "client1", *sysEvent.ClientAuthenticationName)
	require.Equal(t, int64(1), *sysEvent.SequenceNumber)
}

func TestConsumeCloudEventEventGridMqttClientSessionDisconnectedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {  \"namespaceName\": \"myns\",\"clientSessionName\": \"session\",\"clientAuthenticationName\": \"client1\", \"sequenceNumber\": 1, \"disconnectionReason\": \"ClientInitiatedDisconnect\" },  \"type\": \"Microsoft.EventGrid.MQTTClientSessionDisconnected\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventGridMQTTClientSessionDisconnectedEventData](t, events[0].Data)
	require.Equal(t, "session", *sysEvent.ClientSessionName)
	require.Equal(t, "myns", *sysEvent.NamespaceName)
	require.Equal(t, "client1", *sysEvent.ClientAuthenticationName)
	require.Equal(t, int64(1), *sysEvent.SequenceNumber)
	require.Equal(t, azsystemevents.EventGridMQTTClientDisconnectionReasonClientInitiatedDisconnect, *sysEvent.DisconnectionReason)
}

// Event Hub Events
func TestConsumeCloudEventEventHubCaptureFileCreatedEvent(t *testing.T) {
	requestContent := "[    {        \"source\": \"/subscriptions/guid/resourcegroups/rgDataMigrationSample/providers/Microsoft.EventHub/namespaces/tfdatamigratens\",        \"subject\": \"eventhubs/hubdatamigration\",        \"type\": \"microsoft.EventHUB.CaptureFileCreated\",        \"time\": \"2017-08-31T19:12:46.0498024Z\",        \"id\": \"14e87d03-6fbf-4bb2-9a21-92bd1281f247\",        \"data\": {            \"fileUrl\": \"https://tf0831datamigrate.blob.core.windows.net/windturbinecapture/tfdatamigratens/hubdatamigration/1/2017/08/31/19/11/45.avro\",            \"fileType\": \"AzureBlockBlob\",            \"partitionId\": \"1\",            \"sizeInBytes\": 249168,            \"eventCount\": 1500,            \"firstSequenceNumber\": 2400,            \"lastSequenceNumber\": 3899,            \"firstEnqueueTime\": \"2017-08-31T19:12:14.674Z\",            \"lastEnqueueTime\": \"2017-08-31T19:12:44.309Z\"        },  \"specversion\": \"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventHubCaptureFileCreatedEventData](t, events[0].Data)
	require.Equal(t, "AzureBlockBlob", *sysEvent.FileType)
	require.Equal(t, "https://tf0831datamigrate.blob.core.windows.net/windturbinecapture/tfdatamigratens/hubdatamigration/1/2017/08/31/19/11/45.avro", *sysEvent.Fileurl)
}
