//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"

	"github.com/stretchr/testify/require"
)

// As close a part as I could get to @JoshLove-msft's tests here:
//   https://github.com/Azure/azure-sdk-for-net/blob/main/sdk/eventgrid/Azure.Messaging.EventGrid/tests/ConsumeEventTests.cs

func TestParsesEventGridEnvelope(t *testing.T) {
	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"mySubject\",  \"data\": {    \"validationCode\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"validationUrl\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"eventType\": \"Microsoft.EventGrid.SubscriptionValidationEvent\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)
	var egEvent = events[0]
	require.Equal(t, "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", *egEvent.Topic)
	require.Equal(t, "2d1781af-3a4c-4d7c-bd0c-e34b19da4e66", *egEvent.ID)
	require.Equal(t, "mySubject", *egEvent.Subject)
	require.Equal(t, string(azsystemevents.TypeSubscriptionValidation), *egEvent.EventType)
	require.Equal(t, mustParseTime(t, "2018-01-25T22:12:19.4556811Z"), *egEvent.EventTime)
	require.Equal(t, "1", *egEvent.DataVersion)
}

func TestParsesEventGridEnvelopeUsingConverter(t *testing.T) {
	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"mySubject\",  \"data\": {    \"validationCode\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"validationUrl\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"eventType\": \"Microsoft.EventGrid.SubscriptionValidationEvent\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)
	var egEvent = events[0]
	require.Equal(t, "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", *egEvent.Topic)
	require.Equal(t, "2d1781af-3a4c-4d7c-bd0c-e34b19da4e66", *egEvent.ID)
	require.Equal(t, "mySubject", *egEvent.Subject)
	require.Equal(t, string(azsystemevents.TypeSubscriptionValidation), *egEvent.EventType)
	require.Equal(t, mustParseTime(t, "2018-01-25T22:12:19.4556811Z"), *egEvent.EventTime)
	require.Equal(t, "1", *egEvent.DataVersion)
}

func TestConsumeEventGridSubscriptionDeletedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {    \"eventSubscriptionId\": \"/subscriptions/id/resourceGroups/rg/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/eventsubscription1\"  },  \"eventType\": \"Microsoft.EventGrid.SubscriptionDeletedEvent\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.SubscriptionDeletedEventData](t, events[0].Data)
	require.Equal(t,
		"/subscriptions/id/resourceGroups/rg/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/eventsubscription1",
		*sysEvent.EventSubscriptionID)
}

func TestConsumeStorageBlobDeletedEventWithExtraProperty(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)
	require.NotEmpty(t, events)

	for _, egEvent := range events {
		require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), *egEvent.EventType)
		sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, events[0].Data)
		require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.URL)
		require.Equal(t, "/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount", *egEvent.Topic)
	}
}

func TestConsumeEventNotWrappedInAnArray(t *testing.T) {
	requestContent := "{  \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}"

	egEvent := parseEvent(t, requestContent)
	require.NotEmpty(t, egEvent)

	require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), *egEvent.EventType)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, egEvent.Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.URL)
}

func TestConsumeEventNotWrappedInAnArrayWithConverter(t *testing.T) {
	requestContent := "{  \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}"

	egEvent := parseEvent(t, requestContent)

	require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), *egEvent.EventType)

	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, egEvent.Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.URL)
}

func TestConsumeMultipleEventsInSameBatch(t *testing.T) {
	requestContent := "[ " +
		"{  \"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/file1.txt\",  \"eventType\": \"Microsoft.Storage.BlobCreated\",  \"eventTime\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"PutBlockList\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"contentType\": \"text/plain\",    \"contentLength\": 447,    \"blobType\": \"BlockBlob\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/file1.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}, " +
		"{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}, " +
		"{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	require.Equal(t, 3, len(events))

	for _, egEvent := range events {
		switch azsystemevents.Type(*egEvent.EventType) {
		case azsystemevents.TypeStorageBlobCreated:
			blobCreated := deserializeSystemEvent[azsystemevents.StorageBlobCreatedEventData](t, egEvent.Data)
			require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/file1.txt", *blobCreated.URL)
		case azsystemevents.TypeStorageBlobDeleted:
			blobDeleted := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, egEvent.Data)
			require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *blobDeleted.URL)
		}
	}
}

func TestConsumeEventUsingBinaryDataExtensionMethod(t *testing.T) {
	messageBody := "{  \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}"

	egEvent := parseEvent(t, messageBody)

	require.NotEmpty(t, egEvent)

	switch azsystemevents.Type(*egEvent.EventType) {
	case azsystemevents.TypeStorageBlobDeleted:
		blobDeleted := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, egEvent.Data)
		require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *blobDeleted.URL)
	}
}

// func TestParseBinaryDataThrowsOnMultipleEgEvents(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"mySubject\",  \"data\": {    \"validationCode\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"validationUrl\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"eventType\": \"Microsoft.EventGrid.SubscriptionValidationEvent\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}, {  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"mySubject\",  \"data\": {    \"validationCode\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"validationUrl\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"eventType\": \"Microsoft.EventGrid.SubscriptionValidationEvent\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]";

// 	Assert.That(() => EventGridEvent.Parse(new BinaryData(requestContent)),
// 		Throws.InstanceOf<ArgumentException>());
// }

func TestConsumeAppConfigurationKeyValueDeletedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.AppConfiguration.KeyValueDeleted\",\"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"key\":\"key1\",\"label\":\"label1\",\"etag\":\"etag1\"}, \"dataVersion\": \"\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationKeyValueDeletedEventData](t, events[0].Data)
	require.Equal(t, "key1", *sysEvent.Key)
}

func TestConsumeAppConfigurationKeyValueModifiedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.AppConfiguration.KeyValueModified\",\"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"key\":\"key1\",\"label\":\"label1\",\"etag\":\"etag1\"}, \"dataVersion\": \"\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.AppConfigurationKeyValueModifiedEventData](t, events[0].Data)
	require.Equal(t, "key1", *sysEvent.Key)
}

func TestConsumeContainerRegistryImagePushedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.ContainerRegistry.ImagePushed\",  \"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"eventID\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"testaction\",\"target\":{\"mediaType\":\"test\",\"size\":20,\"digest\":\"digest1\",\"length\":20,\"repository\":\"test\",\"url\":\"url1\",\"tag\":\"test\"},\"request\":{\"id\":\"id\",\"addr\":\"127.0.0.1\",\"host\":\"test\",\"method\":\"method1\",\"useragent\":\"useragent1\"},\"actor\":{\"name\":\"testactor\"},\"source\":{\"addr\":\"127.0.0.1\",\"instanceID\":\"id\"}},  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryImagePushedEventData](t, events[0].Data)
	require.Equal(t, "127.0.0.1", *sysEvent.Request.Addr)
}

func TestConsumeContainerRegistryImageDeletedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.ContainerRegistry.ImageDeleted\",  \"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"eventID\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"testaction\",\"target\":{\"mediaType\":\"test\",\"size\":20,\"digest\":\"digest1\",\"length\":20,\"repository\":\"test\",\"url\":\"url1\",\"tag\":\"test\"},\"request\":{\"id\":\"id\",\"addr\":\"127.0.0.1\",\"host\":\"test\",\"method\":\"method1\",\"useragent\":\"useragent1\"},\"actor\":{\"name\":\"testactor\"},\"source\":{\"addr\":\"127.0.0.1\",\"instanceID\":\"id\"}},  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryImageDeletedEventData](t, events[0].Data)
	require.Equal(t, "testactor", *sysEvent.Actor.Name)
}

func TestConsumeContainerRegistryChartDeletedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.ContainerRegistry.ChartDeleted\",  \"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"id\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"action1\",\"target\":{\"mediaType\":\"mediatype1\",\"size\":20,\"digest\":\"digest1\",\"repository\":null,\"tag\":null,\"name\":\"name1\",\"version\":null}}, \"dataVersion\":\"\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryChartDeletedEventData](t, events[0].Data)
	require.Equal(t, "mediatype1", *sysEvent.Target.MediaType)
}

func TestConsumeContainerRegistryChartPushedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.ContainerRegistry/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.ContainerRegistry.ChartPushed\",  \"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"id\":\"id\",\"timestamp\":\"2018-06-20T12:00:33.6125843-07:00\",\"action\":\"action1\",\"target\":{\"mediaType\":\"mediatype1\",\"size\":40,\"digest\":\"digest1\",\"repository\":null,\"tag\":null,\"name\":\"name1\",\"version\":null}}, \"dataVersion\":\"\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerRegistryChartPushedEventData](t, events[0].Data)
	require.Equal(t, "mediatype1", *sysEvent.Target.MediaType)
}

func TestConsumeContainerServiceSupportEndedEvent(t *testing.T) {
	requestContent := `
	{
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"eventType": "Microsoft.ContainerService.ClusterSupportEnded",
		"eventTime": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"kubernetesVersion": "1.23.15"
		},
		"dataVersion": "1",
		"metadataVersion": "1"
	}`

	event := parseEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent.KubernetesVersion)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEndedEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent2.KubernetesVersion)
}

func TestConsumeContainerServiceSupportEndingEvent(t *testing.T) {
	requestContent := `
	{
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"eventType": "Microsoft.ContainerService.ClusterSupportEnding",
		"eventTime": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"kubernetesVersion": "1.23.15"
		},
		"dataVersion": "1",
		"metadataVersion": "1"
	}`

	event := parseEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent.KubernetesVersion)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceClusterSupportEndingEventData](t, event.Data)
	require.Equal(t, "1.23.15", *sysEvent2.KubernetesVersion)
}

func TestConsumeContainerServiceNodePoolRollingFailed(t *testing.T) {
	requestContent := `
	{
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"eventType": "Microsoft.ContainerService.NodePoolRollingFailed",
		"eventTime": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"dataVersion": "1",
		"metadataVersion": "1"
	}`

	event := parseEvent(t, requestContent)

	require.NotEmpty(t, event)

	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingFailedEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent2.NodePoolName)
}

func TestConsumeContainerServiceNodePoolRollingStarted(t *testing.T) {
	requestContent := `
	{
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"eventType": "Microsoft.ContainerService.NodePoolRollingStarted",
		"eventTime": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"dataVersion": "1",
		"metadataVersion": "1"
	}`

	event := parseEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingStartedEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent2.NodePoolName)
}

func TestConsumeContainerServiceNodePoolRollingSucceeded(t *testing.T) {
	requestContent := `
	{
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.ContainerService/managedClusters/{cluster}",
		"subject": "{cluster}",
		"eventType": "Microsoft.ContainerService.NodePoolRollingSucceeded",
		"eventTime": "2023-03-29T18:00:00.0000000Z",
		"id": "1234567890abcdef1234567890abcdef12345678",
		"data": {
			"nodePoolName": "nodepool1"
		},
		"dataVersion": "1",
		"metadataVersion": "1"
	}`

	event := parseEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent.NodePoolName)

	sysEvent2 := deserializeSystemEvent[azsystemevents.ContainerServiceNodePoolRollingSucceededEventData](t, event.Data)
	require.Equal(t, "nodepool1", *sysEvent2.NodePoolName)
}

func TestConsumeIoTHubDeviceCreatedEvent(t *testing.T) {
	requestContent := "[{ \"id\": \"2da5e9b4-4e38-04c1-cc58-9da0b37230c0\", \"topic\": \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\", \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\", \"eventType\": \"Microsoft.Devices.DeviceCreated\", \"eventTime\": \"2018-07-03T23:20:07.6532054Z\",    \"data\": {      \"twin\": {        \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",        \"etag\": \"AAAAAAAAAAE=\",        \"deviceEtag\": null,        \"status\": \"enabled\",        \"statusUpdateTime\": \"0001-01-01T00:00:00\",        \"connectionState\": \"Disconnected\",        \"lastActivityTime\": \"0001-01-01T00:00:00\",        \"cloudToDeviceMessageCount\": 0,        \"authenticationType\": \"sas\",        \"x509Thumbprint\": {          \"primaryThumbprint\": null,          \"secondaryThumbprint\": null        },        \"version\": 2,        \"properties\": {          \"desired\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          },          \"reported\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          }        }      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\"    },    \"dataVersion\": \"\",    \"metadataVersion\": \"1\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IOTHubDeviceCreatedEventData](t, events[0].Data)
	require.Equal(t, "enabled", *sysEvent.Twin.Status)
}

func TestConsumeIoTHubDeviceDeletedEvent(t *testing.T) {
	requestContent := "[  {    \"id\": \"aaaf95c6-ed99-b307-e321-81d8e4f731a6\",    \"topic\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"eventType\": \"Microsoft.Devices.DeviceDeleted\",    \"eventTime\": \"2018-07-03T23:21:33.2753956Z\",    \"data\": {      \"twin\": {        \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",        \"etag\": \"AAAAAAAAAAI=\",        \"deviceEtag\": null,        \"status\": \"enabled\",        \"statusUpdateTime\": \"0001-01-01T00:00:00\",        \"connectionState\": \"Disconnected\",        \"lastActivityTime\": \"0001-01-01T00:00:00\",        \"cloudToDeviceMessageCount\": 0,        \"authenticationType\": \"sas\",        \"x509Thumbprint\": {          \"primaryThumbprint\": null,          \"secondaryThumbprint\": null        },        \"version\": 3,        \"tags\": {          \"testKey\": \"testValue\"        },        \"properties\": {          \"desired\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          },          \"reported\": {            \"$metadata\": {              \"$lastUpdated\": \"2018-07-03T23:20:07.6532054Z\"            },            \"$version\": 1          }        }      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\"    },    \"dataVersion\": \"\",    \"metadataVersion\": \"1\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IOTHubDeviceDeletedEventData](t, events[0].Data)
	require.Equal(t, "AAAAAAAAAAI=", *sysEvent.Twin.Etag)
}

func TestConsumeIoTHubDeviceConnectedEvent(t *testing.T) {
	requestContent := "[  {    \"id\": \"fbfd8ee1-cf78-74c6-dbcf-e1c58638ccbd\",    \"topic\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"eventType\": \"Microsoft.Devices.DeviceConnected\",    \"eventTime\": \"2018-07-03T23:20:11.6921933+00:00\",    \"data\": {      \"deviceConnectionStateEventInfo\": {        \"sequenceNumber\":          \"000000000000000001D4132452F67CE200000002000000000000000000000001\"      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",      \"moduleId\": \"\"    },    \"dataVersion\": \"\",    \"metadataVersion\": \"1\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IOTHubDeviceConnectedEventData](t, events[0].Data)
	require.Equal(t, "EGTESTHUB1", *sysEvent.HubName)
}

func TestConsumeIoTHubDeviceDisconnectedEvent(t *testing.T) {
	requestContent := "[  {    \"id\": \"877f0b10-a086-98ec-27b8-6ae2dfbf5f67\",    \"topic\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"eventType\": \"Microsoft.Devices.DeviceDisconnected\",    \"eventTime\": \"2018-07-03T23:20:52.646434+00:00\",    \"data\": {      \"deviceConnectionStateEventInfo\": {        \"sequenceNumber\":          \"000000000000000001D4132452F67CE200000002000000000000000000000002\"      },      \"hubName\": \"EGTESTHUB1\",      \"deviceId\": \"48e44e11-1437-4907-83b1-4a8d7e89859e\",      \"moduleId\": \"\"    },    \"dataVersion\": \"\",    \"metadataVersion\": \"1\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IOTHubDeviceDisconnectedEventData](t, events[0].Data)
	require.Equal(t, "000000000000000001D4132452F67CE200000002000000000000000000000002", *sysEvent.DeviceConnectionStateEventInfo.SequenceNumber)
}

func TestConsumeIoTHubDeviceTelemetryEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"877f0b10-a086-98ec-27b8-6ae2dfbf5f67\",    \"topic\":      \"/SUBSCRIPTIONS/BDF55CDD-8DAB-4CF4-9B2F-C21E8A780472/RESOURCEGROUPS/EGTESTRG/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/EGTESTHUB1\",    \"subject\": \"devices/48e44e11-1437-4907-83b1-4a8d7e89859e\",    \"eventType\": \"Microsoft.Devices.DeviceTelemetry\",    \"eventTime\": \"2018-07-03T23:20:52.646434+00:00\",    \"data\": { \"body\": { \"Weather\": { \"Temperature\": 900  }, \"Location\": \"USA\"  },  \"properties\": {  \"Status\": \"Active\"  },  \"systemProperties\": { \"iothub-content-type\": \"application/json\", \"iothub-content-encoding\": \"utf-8\"   } }, \"dataVersion\": \"\"}   ]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.IOTHubDeviceTelemetryEventData](t, events[0].Data)
	require.Equal(t, "Active", *sysEvent.Properties["Status"])
}

// Event Hub Events
func TestConsumeEventHubCaptureFileCreatedEvent(t *testing.T) {
	requestContent := "[    {        \"topic\": \"/subscriptions/guid/resourcegroups/rgDataMigrationSample/providers/Microsoft.EventHub/namespaces/tfdatamigratens\",        \"subject\": \"eventhubs/hubdatamigration\",        \"eventType\": \"microsoft.EventHUB.CaptureFileCreated\",        \"eventTime\": \"2017-08-31T19:12:46.0498024Z\",        \"id\": \"14e87d03-6fbf-4bb2-9a21-92bd1281f247\",        \"data\": {            \"fileUrl\": \"https://tf0831datamigrate.blob.core.windows.net/windturbinecapture/tfdatamigratens/hubdatamigration/1/2017/08/31/19/11/45.avro\",            \"fileType\": \"AzureBlockBlob\",            \"partitionId\": \"1\",            \"sizeInBytes\": 249168,            \"eventCount\": 1500,            \"firstSequenceNumber\": 2400,            \"lastSequenceNumber\": 3899,            \"firstEnqueueTime\": \"2017-08-31T19:12:14.674Z\",            \"lastEnqueueTime\": \"2017-08-31T19:12:44.309Z\"        },        \"dataVersion\": \"\",        \"metadataVersion\": \"1\"    }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.EventHubCaptureFileCreatedEventData](t, events[0].Data)
	require.Equal(t, "AzureBlockBlob", *sysEvent.FileType)
}

// MachineLearningServices events
func TestConsumeMachineLearningServicesModelRegisteredEvent(t *testing.T) {
	requestContent := "[{\"topic\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"eventType\":\"Microsoft.MachineLearningServices.ModelRegistered\",\"subject\":\"models/sklearn_regression_model:3\",\"eventTime\":\"2019-10-17T22:23:57.5350054+00:00\",\"id\":\"3b73ee51-bbf4-480d-9112-cfc23b41bfdb\",\"data\":{\"modelName\":\"sklearn_regression_model\",\"modelVersion\":\"3\",\"modelTags\":{\"area\":\"diabetes\",\"type\":\"regression\"},\"modelProperties\":{\"area\":\"test\"}},\"dataVersion\":\"2\",\"metadataVersion\":\"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesModelRegisteredEventData](t, events[0].Data)
	require.Equal(t, "sklearn_regression_model", *sysEvent.ModelName)
	require.Equal(t, "3", *sysEvent.ModelVersion)

	require.Equal(t, "regression", sysEvent.ModelTags.(map[string]any)["type"])
	require.Equal(t, "test", sysEvent.ModelProperties.(map[string]any)["area"])
}

func TestConsumeMachineLearningServicesModelDeployedEvent(t *testing.T) {
	requestContent := "[{\"topic\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"eventType\":\"Microsoft.MachineLearningServices.ModelDeployed\",\"subject\":\"endpoints/aciservice1\",\"eventTime\":\"2019-10-23T18:20:08.8824474+00:00\",\"id\":\"40d0b167-be44-477b-9d23-a2befba7cde0\",\"data\":{\"serviceName\":\"aciservice1\",\"serviceComputeType\":\"ACI\",\"serviceTags\":{\"mytag\":\"test tag\"},\"serviceProperties\":{\"myprop\":\"test property\"},\"modelIds\":\"my_first_model:1,my_second_model:1\"},\"dataVersion\":\"2\",\"metadataVersion\":\"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesModelDeployedEventData](t, events[0].Data)
	require.Equal(t, "aciservice1", *sysEvent.ServiceName)
	sysEvent2 := deserializeSystemEvent[azsystemevents.MachineLearningServicesModelDeployedEventData](t, events[0].Data)
	require.Equal(t, 2, len(strings.Split(*sysEvent2.ModelIDs, ",")))
}

func TestConsumeMachineLearningServicesRunCompletedEvent(t *testing.T) {
	requestContent := "[{\"topic\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"eventType\":\"Microsoft.MachineLearningServices.RunCompleted\",\"subject\":\"experiments/0fa9dfaa-cba3-4fa7-b590-23e48548f5c1/runs/AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"eventTime\":\"2019-10-18T19:29:55.8856038+00:00\",\"id\":\"044ac44d-462c-4043-99eb-d9e01dc760ab\",\"data\":{\"experimentId\":\"0fa9dfaa-cba3-4fa7-b590-23e48548f5c1\",\"experimentName\":\"automl-local-regression\",\"runId\":\"AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"runType\":\"automl\",\"RunTags\":{\"experiment_status\":\"ModelSelection\",\"experiment_status_descr\":\"Beginning model selection.\"},\"runProperties\":{\"num_iterations\":\"10\",\"target\":\"local\"}},\"dataVersion\":\"2\",\"metadataVersion\":\"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesRunCompletedEventData](t, events[0].Data)
	require.Equal(t, "AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc", *sysEvent.RunID)
	sysEvent2 := deserializeSystemEvent[azsystemevents.MachineLearningServicesRunCompletedEventData](t, events[0].Data)
	require.Equal(t, "automl-local-regression", *sysEvent2.ExperimentName)
}

func TestConsumeMachineLearningServicesRunStatusChangedEvent(t *testing.T) {
	requestContent := "[{\"topic\":\"/subscriptions/a5fe3bc5-98f0-4c84-affc-a589f54d9b23/resourceGroups/jenns/providers/Microsoft.MachineLearningServices/workspaces/jenns-canary\",\"eventType\":\"Microsoft.MachineLearningServices.RunStatusChanged\",\"subject\":\"experiments/0fa9dfaa-cba3-4fa7-b590-23e48548f5c1/runs/AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"eventTime\":\"2020-03-09T23:53:04.4579724Z\",\"id\":\"aa8cd7df-fe28-5d5d-9b40-3342dbc2a887\",\"data\":{\"runStatus\": \"Running\",\"experimentId\":\"0fa9dfaa-cba3-4fa7-b590-23e48548f5c1\",\"experimentName\":\"automl-local-regression\",\"runId\":\"AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc\",\"runType\":\"automl\",\"runTags\":{\"experiment_status\":\"ModelSelection\",\"experiment_status_descr\":\"Beginning model selection.\"},\"runProperties\":{\"num_iterations\":\"10\",\"target\":\"local\"}},\"dataVersion\":\"2\",\"metadataVersion\":\"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesRunStatusChangedEventData](t, events[0].Data)
	require.Equal(t, "AutoML_ad912b2d-6467-4f32-a616-dbe4af6dd8fc", *sysEvent.RunID)
	require.Equal(t, "automl-local-regression", *sysEvent.ExperimentName)
	require.Equal(t, "Running", *sysEvent.RunStatus)
	require.Equal(t, "automl", *sysEvent.RunType)
}

func TestConsumeMachineLearningServicesDatasetDriftDetectedEvent(t *testing.T) {
	requestContent := "[{\"topic\":\"/subscriptions/60582a10-b9fd-49f1-a546-c4194134bba8/resourceGroups/copetersRG/providers/Microsoft.MachineLearningServices/workspaces/driftDemoWS\",\"eventType\":\"Microsoft.MachineLearningServices.DatasetDriftDetected\",\"subject\":\"datadrift/01d29aa4-e6a4-470a-9ef3-66660d21f8ef/run/01d29aa4-e6a4-470a-9ef3-66660d21f8ef_1571590300380\",\"eventTime\":\"2019-10-20T17:08:08.467191+00:00\",\"id\":\"2684de79-b145-4dcf-ad2e-6a1db798585f\",\"data\":{\"dataDriftId\":\"01d29aa4-e6a4-470a-9ef3-66660d21f8ef\",\"dataDriftName\":\"copetersDriftMonitor3\",\"runId\":\"01d29aa4-e6a4-470a-9ef3-66660d21f8ef_1571590300380\",\"baseDatasetId\":\"3c56d136-0f64-4657-a0e8-5162089a88a3\",\"tarAsSystemEventDatasetId\":\"d7e74d2e-c972-4266-b5fb-6c9c182d2a74\",\"driftCoefficient\":0.8350349068479208,\"startTime\":\"2019-07-04T00:00:00+00:00\",\"endTime\":\"2019-07-05T00:00:00+00:00\"},\"dataVersion\":\"2\",\"metadataVersion\":\"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MachineLearningServicesDatasetDriftDetectedEventData](t, events[0].Data)
	require.Equal(t, "copetersDriftMonitor3", *sysEvent.DataDriftName)
}

// Maps events
func TestConsumeMapsGeofenceEnteredEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.Maps.GeofenceEntered\",\"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}, \"dataVersion\":\"\"}]"
	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceEnteredEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

func TestConsumeMapsGeofenceExitedEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.Maps.GeofenceExited\",\"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}, \"dataVersion\":\"\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceExitedEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

func TestConsumeMapsGeofenceResultEvent(t *testing.T) {
	requestContent := "[{  \"id\": \"56afc886-767b-d359-d59e-0da7877166b2\",  \"topic\": \"/SUBSCRIPTIONS/ID/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.Maps/test1\",  \"subject\": \"test1\",  \"eventType\": \"Microsoft.Maps.GeofenceResult\",\"eventTime\": \"2018-01-02T19:17:44.4383997Z\",  \"data\": {\"expiredGeofenceGeometryId\":[\"id1\",\"id2\"],\"geometries\":[{\"deviceId\":\"id1\",\"distance\":1.0,\"geometryId\":\"gid1\",\"nearestLat\":72.4,\"nearestLon\":100.4,\"udId\":\"id22\"}],\"invalidPeriodGeofenceGeometryId\":[\"id1\",\"id2\"],\"isEventPublished\":true}, \"dataVersion\":\"\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MapsGeofenceResultEventData](t, events[0].Data)
	require.Equal(t, float32(1.0), *sysEvent.Geometries[0].Distance)
}

// Media Services events
func TestConsumeMediaMediaJobStateChangeEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobStateChange\",  \"eventTime\": \"2018-10-12T15:14:20.2412317\",  \"id\": \"341520d0-dac0-4930-97dd-3085538c624f\",  \"data\": {    \"previousState\": \"Scheduled\",    \"state\": \"Processing\",    \"correlationData\": {}  },  \"dataVersion\": \"2.0\",  \"metadataVersion\": \"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobStateChangeEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.State)
}

func TestConsumeMediaJobOutputStateChangeEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputStateChange\",  \"eventTime\": \"2018-10-12T15:14:17.8962704\",  \"id\": \"8d0305c0-28c0-4cc9-b613-776e4dd31e9a\",  \"data\": {    \"previousState\": \"Scheduled\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",      \"error\": {\"code\":\"ServiceError\", \"message\":\"error message\", \"category\":\"Service\", \"retry\":\"DoNotRetry\", \"details\":[{\"code\":\"code\", \"message\":\"Service Error Message\"}]},      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 0,      \"state\": \"Processing\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)
	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputStateChangeEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.PreviousState)
	mediaOutputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *mediaOutputAsset.State)

	require.Equal(t, "output-2ac2fe75-6557-4de5-ab25-5713b74a6901", *mediaOutputAsset.AssetName)

	require.Equal(t, azsystemevents.MediaJobErrorCategoryService, *mediaOutputAsset.Error.Category)
	require.Equal(t, azsystemevents.MediaJobErrorCodeServiceError, *mediaOutputAsset.Error.Code)
}

func TestConsumeMediaJobOutputStateChangeEvent_UnknownError(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputStateChange\",  \"eventTime\": \"2018-10-12T15:14:17.8962704\",  \"id\": \"8d0305c0-28c0-4cc9-b613-776e4dd31e9a\",  \"data\": {    \"previousState\": \"Scheduled\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",      \"error\": {\"code\":\"SomeNewCode\", \"message\":\"error message\", \"category\":\"SomeNewCategory\", \"retry\":\"DoNotRetry\", \"details\":[{\"code\":\"code\", \"message\":\"Service Error Message\"}]},      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 0,      \"state\": \"Processing\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputStateChangeEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.Output.(*azsystemevents.MediaJobOutputAsset).State)

	outputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, "output-2ac2fe75-6557-4de5-ab25-5713b74a6901", *outputAsset.AssetName)

	require.Equal(t, azsystemevents.MediaJobErrorCategory("SomeNewCategory"), *outputAsset.Error.Category)
	require.Equal(t, azsystemevents.MediaJobErrorCode("SomeNewCode"), *outputAsset.Error.Code)
}

func TestConsumeMediaJobScheduledEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobScheduled\",  \"eventTime\": \"2018-10-12T15:14:11.3028183\",  \"id\": \"9b17dbf0-355d-4fb0-9a73-e76b150858c8\",  \"data\": {    \"previousState\": \"Queued\",    \"state\": \"Scheduled\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobScheduledEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateQueued, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.State)
}

func TestConsumeMediaJobProcessingEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobProcessing\",  \"eventTime\": \"2018-10-12T15:14:20.2412317\",  \"id\": \"72162c44-c7f4-437a-9592-48b83cec2d18\",  \"data\": {    \"previousState\": \"Scheduled\",    \"state\": \"Processing\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobProcessingEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.State)
}

func TestConsumeMediaJobCancelingEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",  \"eventType\": \"Microsoft.Media.JobCanceling\",  \"eventTime\": \"2018-10-12T15:41:50.5513295\",  \"id\": \"1f9a488b-abe3-4fca-80b8-aae59bf7f123\",  \"data\": {    \"previousState\": \"Processing\",    \"state\": \"Canceling\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobCancelingEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateCanceling, *sysEvent.State)
}

func TestConsumeMediaJobFinishedEvent(t *testing.T) {
	requestContent := "[{ \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-298338bb-f8d1-4d0f-9fde-544e0ac4d983\",  \"eventType\": \"Microsoft.Media.JobFinished\",  \"eventTime\": \"2018-10-01T20:58:26.7886175\",  \"id\": \"83f8464d-be94-48e5-b67b-46c6199fe28e\",  \"data\": {    \"outputs\": [      {        \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",        \"assetName\": \"output-298338bb-f8d1-4d0f-9fde-544e0ac4d983\",       \"label\": \"VideoAnalyzerPreset_0\",        \"progress\": 100,        \"state\": \"Finished\"      }    ],    \"previousState\": \"Processing\",    \"state\": \"Finished\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\" }]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobFinishedEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateFinished, *sysEvent.State)
	require.Equal(t, 1, len(sysEvent.Outputs))

	outputAsset := sysEvent.Outputs[0].(*azsystemevents.MediaJobOutputAsset)

	require.Equal(t, azsystemevents.MediaJobStateFinished, *outputAsset.State)
	require.Nil(t, outputAsset.Error)
	require.Equal(t, int64(100), *outputAsset.Progress)
	require.Equal(t, "output-298338bb-f8d1-4d0f-9fde-544e0ac4d983", *outputAsset.AssetName)
}

func TestConsumeMediaJobCanceledEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",  \"eventType\": \"Microsoft.Media.JobCanceled\",  \"eventTime\": \"2018-10-12T15:42:05.6519929\",  \"id\": \"3fef7871-f916-4980-8a45-e79a2675808b\",  \"data\": {    \"outputs\": [      {        \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",        \"assetName\": \"output-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",        \"error\": {\"code\":\"ServiceError\", \"message\":\"error message\", \"category\":\"Service\", \"retry\":\"DoNotRetry\", \"details\":[{\"code\":\"code\", \"message\":\"Service Error Message\"}]},      \"label\": \"VideoAnalyzerPreset_0\",        \"progress\": 83,        \"state\": \"Canceled\"      }    ],    \"previousState\": \"Canceling\",    \"state\": \"Canceled\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)
	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobCanceledEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateCanceling, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateCanceled, *sysEvent.State)
	require.Equal(t, 1, len(sysEvent.Outputs))

	outputAsset := sysEvent.Outputs[0].(*azsystemevents.MediaJobOutputAsset)

	require.Equal(t, azsystemevents.MediaJobStateCanceled, *outputAsset.State)
	require.NotEqual(t, int64(100), *outputAsset.Progress)
	require.Equal(t, "output-7a8215f9-0f8d-48a6-82ed-1ead772bc221", *outputAsset.AssetName)
}

func TestConsumeMediaJobErroredEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobErrored\",  \"eventTime\": \"2018-10-12T15:29:20.9954767\",  \"id\": \"2749e9cf-4095-4723-9bc5-df8e15289135\",  \"data\": {    \"outputs\": [      {        \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",        \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",        \"error\": {          \"category\": \"Service\",          \"code\": \"ServiceError\",          \"details\": [            {              \"code\": \"Internal\",              \"message\": \"Internal error in initializing the task for processing\"            }          ],          \"message\": \"Fatal service error, please contact support.\",          \"retry\": \"DoNotRetry\"        },        \"label\": \"VideoAnalyzerPreset_0\",        \"progress\": 83,        \"state\": \"Error\"      }    ],    \"previousState\": \"Processing\",    \"state\": \"Error\",    \"correlationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobErroredEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateError, *sysEvent.State)
	require.Equal(t, 1, len(sysEvent.Outputs))

	outputAsset := sysEvent.Outputs[0].(*azsystemevents.MediaJobOutputAsset)

	require.Equal(t, azsystemevents.MediaJobStateError, *outputAsset.State)
	require.NotEmpty(t, *outputAsset.Error)
	require.Equal(t, azsystemevents.MediaJobErrorCategoryService, *outputAsset.Error.Category)
	require.Equal(t, azsystemevents.MediaJobErrorCodeServiceError, *outputAsset.Error.Code)
}

func TestConsumeMediaJobOutputCanceledEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",  \"eventType\": \"Microsoft.Media.JobOutputCanceled\",  \"eventTime\": \"2018-10-12T15:42:04.949555\",  \"id\": \"9297cda2-4a50-4622-a679-c3785d27d512\",  \"data\": {    \"previousState\": \"Canceling\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",      \"error\": {\"code\":\"ServiceError\", \"message\":\"error message\", \"category\":\"Service\", \"retry\":\"DoNotRetry\", \"details\":[{\"code\":\"code\", \"message\":\"Service Error Message\"}]},      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 83,      \"state\": \"Canceled\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputCanceledEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateCanceling, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateCanceled, *sysEvent.Output.(*azsystemevents.MediaJobOutputAsset).State)
}

func TestConsumeMediaJobOutputCancelingEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",  \"eventType\": \"Microsoft.Media.JobOutputCanceling\",  \"eventTime\": \"2018-10-12T15:42:04.949555\",  \"id\": \"9297cda2-4a50-4622-a679-c3785d27d512\",  \"data\": {    \"previousState\": \"Processing\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-7a8215f9-0f8d-48a6-82ed-1ead772bc221\",      \"error\": {        \"category\": \"Service\",        \"code\": \"ServiceError\",        \"details\": [          {            \"code\": \"Internal\",            \"message\": \"Internal error in initializing the task for processing\"          }        ],        \"message\": \"Fatal service error, please contact support.\",        \"retry\": \"DoNotRetry\"      },      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 83,      \"state\": \"Canceling\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputCancelingEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)
	require.Equal(t, azsystemevents.MediaJobStateCanceling, *sysEvent.Output.(*azsystemevents.MediaJobOutputAsset).State)
}

func TestConsumeMediaJobOutputErroredEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputErrored\",  \"eventTime\": \"2018-10-12T15:29:20.2621252\",  \"id\": \"bc9e6342-f081-49c2-a579-92f506a622c2\",  \"data\": {    \"previousState\": \"Processing\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",      \"error\": {        \"category\": \"Service\",        \"code\": \"ServiceError\",        \"details\": [          {            \"code\": \"Internal\",            \"message\": \"Internal error in initializing the task for processing\"          }        ],        \"message\": \"Fatal service error, please contact support.\",        \"retry\": \"DoNotRetry\"      },      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 83,      \"state\": \"Error\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputErroredEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)

	outputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, azsystemevents.MediaJobStateError, *outputAsset.State)
	require.NotEmpty(t, outputAsset.Error)
	require.Equal(t, azsystemevents.MediaJobErrorCategoryService, *outputAsset.Error.Category)
	require.Equal(t, azsystemevents.MediaJobErrorCodeServiceError, *outputAsset.Error.Code)
}

func TestConsumeMediaJobOutputFinishedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputFinished\",  \"eventTime\": \"2018-10-12T15:29:20.2621252\",  \"id\": \"bc9e6342-f081-49c2-a579-92f506a622c2\",  \"data\": {    \"previousState\": \"Processing\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",            \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 100,      \"state\": \"Finished\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputFinishedEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *sysEvent.PreviousState)

	outputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, azsystemevents.MediaJobStateFinished, *outputAsset.State)
	require.Equal(t, int64(100), *outputAsset.Progress)
	require.Equal(t, "output-2ac2fe75-6557-4de5-ab25-5713b74a6901", *outputAsset.AssetName)
}

func TestConsumeMediaJobOutputProcessingEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputProcessing\",  \"eventTime\": \"2018-10-12T15:14:17.8962704\",  \"id\": \"d48eeb0b-2bfa-4265-a2f8-624654c3781c\",  \"data\": {    \"previousState\": \"Scheduled\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",      \"error\": {        \"category\": \"Service\",        \"code\": \"ServiceError\",        \"details\": [          {            \"code\": \"Internal\",            \"message\": \"Internal error in initializing the task for processing\"          }        ],        \"message\": \"Fatal service error, please contact support.\",        \"retry\": \"DoNotRetry\"      },      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 0,      \"state\": \"Processing\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputProcessingEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *sysEvent.PreviousState)

	outputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, azsystemevents.MediaJobStateProcessing, *outputAsset.State)
}

func TestConsumeMediaJobOutputScheduledEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6901\",  \"eventType\": \"Microsoft.Media.JobOutputScheduled\",  \"eventTime\": \"2018-10-12T15:14:11.2244618\",  \"id\": \"635ca6ea-5306-4590-b2e1-22f172759336\",  \"data\": {    \"previousState\": \"Queued\",    \"output\": {      \"@odata.type\": \"#Microsoft.Media.JobOutputAsset\",      \"assetName\": \"output-2ac2fe75-6557-4de5-ab25-5713b74a6901\",      \"error\": {        \"category\": \"Service\",        \"code\": \"ServiceError\",        \"details\": [          {            \"code\": \"Internal\",            \"message\": \"Internal error in initializing the task for processing\"          }        ],        \"message\": \"Fatal service error, please contact support.\",        \"retry\": \"DoNotRetry\"      },      \"label\": \"VideoAnalyzerPreset_0\",      \"progress\": 0,      \"state\": \"Scheduled\"    },    \"jobCorrelationData\": {}  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputScheduledEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.MediaJobStateQueued, *sysEvent.PreviousState)

	outputAsset := sysEvent.Output.(*azsystemevents.MediaJobOutputAsset)
	require.Equal(t, azsystemevents.MediaJobStateScheduled, *outputAsset.State)
}

func TestConsumeMediaJobOutputProgressEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"transforms/VideoAnalyzerTransform/jobs/job-2ac2fe75-6557-4de5-ab25-5713b74a6981\",  \"eventType\": \"Microsoft.Media.JobOutputProgress\",  \"eventTime\": \"2018-10-12T15:14:11.2244618\",  \"id\": \"635ca6ea-5306-4590-b2e1-22f172759336\",  \"data\": {    \"jobCorrelationData\": {    \"Field1\": \"test1\",    \"Field2\": \"test2\" },    \"label\": \"TestLabel\",    \"progress\": 50 },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaJobOutputProgressEventData](t, events[0].Data)
	require.Equal(t, "TestLabel", *sysEvent.Label)
	require.Equal(t, int64(50), *sysEvent.Progress)

	require.Equal(t, "test1", *sysEvent.JobCorrelationData["Field1"])
	require.Equal(t, "test2", *sysEvent.JobCorrelationData["Field2"])
}

func TestConsumeMediaLiveEventEncoderConnectedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventEncoderConnected\",  \"eventTime\": \"2018-10-12T15:52:04.2013501\",  \"id\": \"3d1f5b26-c466-47e7-927b-900985e0c5d5\",  \"data\": {    \"ingestUrl\": \"rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59\",    \"streamId\": \"Mystream1\",    \"encoderIp\": \"<ip address>\",    \"encoderPort\": \"3557\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventEncoderConnectedEventData](t, events[0].Data)

	require.Equal(t, "rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59", *sysEvent.IngestURL)
	require.Equal(t, "Mystream1", *sysEvent.StreamID)
	require.Equal(t, "<ip address>", *sysEvent.EncoderIP)
	require.Equal(t, "3557", *sysEvent.EncoderPort)
}

func TestConsumeMediaLiveEventConnectionRejectedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventConnectionRejected\",  \"eventTime\": \"2018-10-12T15:52:04.2013501\",  \"id\": \"3d1f5b26-c466-47e7-927b-900985e0c5d5\",  \"data\": {    \"ingestUrl\": \"rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59\",    \"streamId\": \"Mystream1\",    \"encoderIp\": \"<ip address>\",    \"encoderPort\": \"3557\",    \"resultCode\": \"MPE_INGEST_CODEC_NOT_SUPPORTED\"   },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)

	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventConnectionRejectedEventData](t, events[0].Data)
	require.Equal(t, "rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59", *sysEvent.IngestURL)
	require.Equal(t, "Mystream1", *sysEvent.StreamID)
	require.Equal(t, "<ip address>", *sysEvent.EncoderIP)
	require.Equal(t, "3557", *sysEvent.EncoderPort)
}

func TestConsumeMediaLiveEventEncoderDisconnectedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventEncoderDisconnected\",  \"eventTime\": \"2018-10-12T15:52:19.8982128\",  \"id\": \"e4b55140-42d2-4c24-b08e-9aa12f1587fc\",  \"data\": {    \"ingestUrl\": \"rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59\",    \"streamId\": \"Mystream1\",    \"encoderIp\": \"<ip address>\",    \"encoderPort\": \"3557\",    \"resultCode\": \"MPE_CLIENT_TERMINATED_SESSION\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventEncoderDisconnectedEventData](t, events[0].Data)

	require.Equal(t, "MPE_CLIENT_TERMINATED_SESSION", *sysEvent.ResultCode)
	require.Equal(t, "rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59", *sysEvent.IngestURL)
	require.Equal(t, "Mystream1", *sysEvent.StreamID)
	require.Equal(t, "<ip address>", *sysEvent.EncoderIP)
	require.Equal(t, "3557", *sysEvent.EncoderPort)
}

func TestConsumeMediaLiveEventIncomingStreamReceivedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIncomingStreamReceived\",  \"eventTime\": \"2018-10-12T15:52:16.5726463Z\",  \"id\": \"eb688fa1-5a19-4703-8aeb-6a65a09790da\",  \"data\": {    \"ingestUrl\": \"rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59\",    \"trackType\": \"audio\",    \"trackName\": \"audio_160000\",    \"bitrate\": 160000,    \"encoderIp\": \"<ip address>\",    \"encoderPort\": \"3557\",    \"timestamp\": \"66\",    \"duration\": \"1950\",    \"timescale\": \"1000\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventIncomingStreamReceivedEventData](t, events[0].Data)

	require.Equal(t, "rtmp://liveevent-ec9d26a8.channel.media.azure.net:1935/live/cb5540b10a5646218c1328be95050c59", *sysEvent.IngestURL)
	require.Equal(t, "<ip address>", *sysEvent.EncoderIP)
	require.Equal(t, "3557", *sysEvent.EncoderPort)
	require.Equal(t, "audio", *sysEvent.TrackType)
	require.Equal(t, "audio_160000", *sysEvent.TrackName)
	require.Equal(t, int64(160000), *sysEvent.Bitrate)
	require.Equal(t, "66", *sysEvent.Timestamp)
	require.Equal(t, "1950", *sysEvent.Duration)
	require.Equal(t, "1000", *sysEvent.Timescale)
}

func TestConsumeMediaLiveEventIncomingStreamsOutOfSyncEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIncomingStreamsOutOfSync\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"minLastTimestamp\": \"10999\",    \"typeOfStreamWithMinLastTimestamp\": \"video\",    \"maxLastTimestamp\": \"100999\",    \"typeOfStreamWithMaxLastTimestamp\": \"audio\",    \"timescaleOfMinLastTimestamp\": \"1000\",  \"timescaleOfMaxLastTimestamp\": \"1000\"    },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventIncomingStreamsOutOfSyncEventData](t, events[0].Data)

	require.Equal(t, "10999", *sysEvent.MinLastTimestamp)
	require.Equal(t, "video", *sysEvent.TypeOfStreamWithMinLastTimestamp)
	require.Equal(t, "100999", *sysEvent.MaxLastTimestamp)
	require.Equal(t, "audio", *sysEvent.TypeOfStreamWithMaxLastTimestamp)
	require.Equal(t, "1000", *sysEvent.TimescaleOfMinLastTimestamp)
	require.Equal(t, "1000", *sysEvent.TimescaleOfMaxLastTimestamp)
}

func TestConsumeMediaLiveEventIncomingVideoStreamsOutOfSyncEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIncomingVideoStreamsOutOfSync\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"firstTimestamp\": \"10999\",    \"firstDuration\": \"2000\",    \"secondTimestamp\": \"100999\",    \"secondDuration\": \"2000\",    \"timescale\": \"1000\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventIncomingVideoStreamsOutOfSyncEventData](t, events[0].Data)

	require.Equal(t, "10999", *sysEvent.FirstTimestamp)
	require.Equal(t, "2000", *sysEvent.FirstDuration)
	require.Equal(t, "100999", *sysEvent.SecondTimestamp)
	require.Equal(t, "2000", *sysEvent.SecondDuration)
	require.Equal(t, "1000", *sysEvent.Timescale)
}

func TestConsumeMediaLiveEventIncomingDataChunkDroppedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIncomingDataChunkDropped\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"timestamp\": \"8999\",    \"trackType\": \"video\",    \"trackName\": \"video1\",    \"bitrate\": 2500000,    \"timescale\": \"1000\",    \"resultCode\": \"FragmentDrop_OverlapTimestamp\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventIncomingDataChunkDroppedEventData](t, events[0].Data)

	require.Equal(t, "8999", *sysEvent.Timestamp)
	require.Equal(t, "video", *sysEvent.TrackType)
	require.Equal(t, "video1", *sysEvent.TrackName)
	require.Equal(t, int64(2500000), *sysEvent.Bitrate)
	require.Equal(t, "1000", *sysEvent.Timescale)
	require.Equal(t, "FragmentDrop_OverlapTimestamp", *sysEvent.ResultCode)
}

func TestConsumeMediaLiveEventIngestHeartbeatEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIngestHeartbeat\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"trackType\": \"video\",    \"trackName\": \"video\",    \"bitrate\": 2500000,    \"incomingBitrate\": 500726,    \"lastTimestamp\": \"11999\",    \"timescale\": \"1000\",    \"overlapCount\": 0,    \"discontinuityCount\": 0,    \"nonincreasingCount\": 0,    \"unexpectedBitrate\": true,    \"state\": \"Running\",    \"healthy\": false,  \"lastFragmentArrivalTime\": \"2021-05-14T23:50:00.00\", \"ingestDriftValue\": \"0\" },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventIngestHeartbeatEventData](t, events[0].Data)
	require.Equal(t, "video", *sysEvent.TrackType)
	require.Equal(t, "video", *sysEvent.TrackName)
	require.Equal(t, int64(2500000), *sysEvent.Bitrate)
	require.Equal(t, int64(500726), *sysEvent.IncomingBitrate)
	require.Equal(t, "11999", *sysEvent.LastTimestamp)
	require.Equal(t, "1000", *sysEvent.Timescale)
	require.Equal(t, int64(0), *sysEvent.OverlapCount)
	require.Equal(t, int64(0), *sysEvent.DiscontinuityCount)
	require.Equal(t, int64(0), *sysEvent.NonincreasingCount)
	require.True(t, *sysEvent.UnexpectedBitrate)
	require.Equal(t, "Running", *sysEvent.State)
	require.False(t, *sysEvent.Healthy)
	require.Equal(t, "0", *sysEvent.IngestDriftValue)
	require.Equal(t, "2021-05-14T23:50:00.00", *sysEvent.LastFragmentArrivalTime)

	requestContent = "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventIngestHeartbeat\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"trackType\": \"video\",    \"trackName\": \"video\",    \"bitrate\": 2500000,    \"incomingBitrate\": 500726,    \"lastTimestamp\": \"11999\",    \"timescale\": \"1000\",    \"overlapCount\": 0,    \"discontinuityCount\": 0,    \"nonincreasingCount\": 0,    \"unexpectedBitrate\": true,    \"state\": \"Running\",    \"healthy\": false,  \"lastFragmentArrivalTime\": \"2021-05-14T23:50:00.00\", \"ingestDriftValue\": \"n/a\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events2 := parseManyEvents(t, requestContent)
	sysEvent2 := deserializeSystemEvent[azsystemevents.MediaLiveEventIngestHeartbeatEventData](t, events2[0].Data)
	require.Nil(t, sysEvent2.IngestDriftValue)
}

func TestConsumeMediaLiveEventChannelArchiveHeartbeatEvent(t *testing.T) {
	requestContent := `[
	{
		"topic": "/subscriptions/<subscription-id>/resourceGroups/<rg-name>/providers/Microsoft.Media/mediaservices/<account-name>",
		"subject": "liveEvent/mle1",
		"eventType": "Microsoft.Media.LiveEventChannelArchiveHeartbeat",
		"eventTime": "2021-05-14T23:50:00.324",
		"id": "7f450938-491f-41e1-b06f-c6cd3965d786",
		"data": {
			"channelLatencyMs": "10",
			"latencyResultCode": "S_OK"
		},
		"dataVersion": "1.0",
		"metadataVersion": "1"
	}]`

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventChannelArchiveHeartbeatEventData](t, events[0].Data)
	require.Equal(t, "10", *sysEvent.ChannelLatencyMS)
	require.Equal(t, "S_OK", *sysEvent.LatencyResultCode)

	requestContent = `[
	{
		"topic": "/subscriptions/<subscription-id>/resourceGroups/<rg-name>/providers/Microsoft.Media/mediaservices/<account-name>",
		"subject": "liveEvent/mle1",
		"eventType": "Microsoft.Media.LiveEventChannelArchiveHeartbeat",
		"eventTime": "2021-05-14T23:50:00.324",
		"id": "7f450938-491f-41e1-b06f-c6cd3965d786",
		"data": {
			"channelLatencyMs": "n/a",
			"latencyResultCode": "S_OK"
		},
		"dataVersion": "1.0",
		"metadataVersion": "1"
	}]`

	events2 := parseManyEvents(t, requestContent)

	sysEvent2 := deserializeSystemEvent[azsystemevents.MediaLiveEventChannelArchiveHeartbeatEventData](t, events2[0].Data)

	// n/a should be translated to null ChannelLatency
	require.Nil(t, sysEvent2.ChannelLatencyMS)
	require.Equal(t, "S_OK", *sysEvent2.LatencyResultCode)
}

func TestConsumeMediaLiveEventTrackDiscontinuityDetectedEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/{subscription id}/resourceGroups/{resource group}/providers/Microsoft.Media/mediaservices/{account name}\",  \"subject\": \"liveEvent/liveevent-ec9d26a8\",  \"eventType\": \"Microsoft.Media.LiveEventTrackDiscontinuityDetected\",  \"eventTime\": \"2018-10-12T15:52:37.3710102\",  \"id\": \"d84727e2-d9c0-4a21-a66b-8d23f06b3e06\",  \"data\": {    \"trackType\": \"video\",    \"trackName\": \"video\",    \"bitrate\": 2500000,    \"previousTimestamp\": \"10999\",    \"newTimestamp\": \"14999\",    \"timescale\": \"1000\",    \"discontinuityGap\": \"4000\"  },  \"dataVersion\": \"1.0\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.MediaLiveEventTrackDiscontinuityDetectedEventData](t, events[0].Data)
	require.Equal(t, "video", *sysEvent.TrackType)
	require.Equal(t, "video", *sysEvent.TrackName)
	require.Equal(t, int64(2500000), *sysEvent.Bitrate)
	require.Equal(t, "10999", *sysEvent.PreviousTimestamp)
	require.Equal(t, "14999", *sysEvent.NewTimestamp)
	require.Equal(t, "1000", *sysEvent.Timescale)
	require.Equal(t, "4000", *sysEvent.DiscontinuityGap)
}

// Resource Manager (Azure Subscription/Resource Group) events

const Authorization = `{"scope":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Web/sites/function/host/default","action":"Microsoft.Web/sites/host/listKeys/action","evidence":{"role":"Azure EventGrid Service BuiltIn Role","roleAssignmentScope":"/subscriptions/sub","roleAssignmentId":"rid","roleDefinitionId":"rd","principalId":"principal","principalType":"ServicePrincipal"}}`
const Claims = `{"aud":"https://management.core.windows.net","iat":"16303066","nbf":"16303066","exp":"16303066"}`
const HttpRequest = `{"clientRequestId":"","clientIpAddress":"ip","method":"POST","url":"https://management.azure.com/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Web/sites/function/host/default/listKeys?api-version=2018-11-01"}`

func TestConsumeResourceWriteSuccessEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceWriteSuccess",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	eventData := deserializeSystemEvent[azsystemevents.ResourceWriteSuccessEventData](t, events[0].Data)
	assertResourceEventData(t, eventData)
}

func TestConsumeResourceWriteFailureEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceWriteFailure",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceWriteFailureEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceWriteCancelEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceWriteCancel",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceWriteCancelEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceDeleteSuccessEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceDeleteSuccess",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceDeleteSuccessEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceDeleteFailureEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceDeleteFailure",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceDeleteFailureEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceDeleteCancelEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceDeleteCancel",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceDeleteCancelEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceActionSuccessEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceActionSuccess",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceActionSuccessEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

func TestConsumeResourceActionFailureEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceActionFailure",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)
	eventData := deserializeSystemEvent[azsystemevents.ResourceActionFailureEventData](t, events[0].Data)
	assertResourceEventData(t, eventData)
}

func TestConsumeResourceActionCancelEvent(t *testing.T) {
	requestContent := fmt.Sprintf(`[{"topic":"/subscriptions/subscription-id", "subject":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",    "eventType":"Microsoft.Resources.ResourceActionCancel",    "eventTime":"2017-08-16T03:54:38.2696833Z",    "id":"25b3b0d0-d79b-44d5-9963-440d4e6a9bba",    "data": { "authorization":%s,   "claims":%s,  "correlationId":"54ef1e39-6a82-44b3-abc1-bdeb6ce4d3c6",  "httpRequest":%s,   "resourceProvider":"Microsoft.EventGrid",  "resourceUri":"/subscriptions/subscription-id/resourceGroups/resource-group/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501",        "operationName":"Microsoft.EventGrid/eventSubscriptions/write",    "status":"Succeeded",   "subscriptionId":"subscription-id",  "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47"        },      "dataVersion": "",    "metadataVersion": "1"  }]`, Authorization, Claims, HttpRequest)

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	eventData := deserializeSystemEvent[azsystemevents.ResourceActionCancelEventData](t, events[0].Data)

	assertResourceEventData(t, eventData)
}

// Using dynamic to avoid duplicating the test cases for each event. The events don't share a common base type but they all have the
// properties being tested below.
func assertResourceEventData(t *testing.T, rawEventData any) {
	jsonBytes, err := json.Marshal(rawEventData)
	require.NoError(t, err)

	var eventData *struct {
		TenantID      string
		Authorization struct {
			Scope    string
			Action   string
			Evidence struct {
				Role                string
				RoleAssignmentScope string
				PrincipalType       string
			}
		} `json:"authorization"`
		Claims struct {
			Aud string
		}
		HttpRequest struct {
			Method string
			URL    string
		}
	} = nil

	err = json.Unmarshal(jsonBytes, &eventData)
	require.NoError(t, err)

	require.NotEmpty(t, eventData)
	require.Equal(t, "72f988bf-86f1-41af-91ab-2d7cd011db47", eventData.TenantID)

	//var authorizationJson = JsonDocument.Parse(eventData.Authorization).RootElement

	require.Equal(t, "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Web/sites/function/host/default",
		eventData.Authorization.Scope)
	require.Equal(t, "Microsoft.Web/sites/host/listKeys/action", eventData.Authorization.Action)
	require.Equal(t, "Azure EventGrid Service BuiltIn Role", eventData.Authorization.Evidence.Role)
	require.Equal(t, "/subscriptions/sub", eventData.Authorization.Evidence.RoleAssignmentScope)
	require.Equal(t, "ServicePrincipal", eventData.Authorization.Evidence.PrincipalType)

	require.Equal(t, "https://management.core.windows.net", eventData.Claims.Aud)
	require.Equal(t, "POST", eventData.HttpRequest.Method)
	require.Equal(t, "https://management.azure.com/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Web/sites/function/host/default/listKeys?api-version=2018-11-01", eventData.HttpRequest.URL)
}

// ServiceBus events
func TestConsumeServiceBusActiveMessagesAvailableWithNoListenersEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/id/resourcegroups/rg/providers/Microsoft.ServiceBus/namespaces/testns1\",  \"subject\": \"topics/topic1/subscriptions/sub1\",  \"eventType\": \"Microsoft.ServiceBus.ActiveMessagesAvailableWithNoListeners\",  \"eventTime\": \"2018-02-14T05:12:53.4133526Z\",  \"id\": \"dede87b0-3656-419c-acaf-70c95ddc60f5\",  \"data\": {    \"namespaceName\": \"testns1\",    \"requestUri\": \"https://testns1.servicebus.windows.net/t1/subscriptions/sub1/messages/head\",    \"entityType\": \"subscriber\",    \"queueName\": \"queue1\",    \"topicName\": \"topic1\",    \"subscriptionName\": \"sub1\"  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ServiceBusActiveMessagesAvailableWithNoListenersEventData](t, events[0].Data)
	require.Equal(t, "testns1", *sysEvent.NamespaceName)
}

func TestConsumeServiceBusDeadletterMessagesAvailableWithNoListenersEvent(t *testing.T) {
	requestContent := "[{  \"topic\": \"/subscriptions/id/resourcegroups/rg/providers/Microsoft.ServiceBus/namespaces/testns1\",  \"subject\": \"topics/topic1/subscriptions/sub1\",  \"eventType\": \"Microsoft.ServiceBus.DeadletterMessagesAvailableWithNoListeners\",  \"eventTime\": \"2018-02-14T05:12:53.4133526Z\",  \"id\": \"dede87b0-3656-419c-acaf-70c95ddc60f5\",  \"data\": {    \"namespaceName\": \"testns1\",    \"requestUri\": \"https://testns1.servicebus.windows.net/t1/subscriptions/sub1/messages/head\",    \"entityType\": \"subscriber\",    \"queueName\": \"queue1\",    \"topicName\": \"topic1\",    \"subscriptionName\": \"sub1\"  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ServiceBusDeadletterMessagesAvailableWithNoListenersEventData](t, events[0].Data)
	require.Equal(t, "testns1", *sysEvent.NamespaceName)
}

// Storage events
func TestConsumeStorageBlobCreatedEvent(t *testing.T) {
	requestContent := "[ {  \"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/file1.txt\",  \"eventType\": \"Microsoft.Storage.BlobCreated\",  \"eventTime\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"PutBlockList\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"contentType\": \"text/plain\",    \"contentLength\": 447,    \"blobType\": \"BlockBlob\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/file1.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\" },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobCreatedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/file1.txt", *sysEvent.URL)
}

func TestConsumeStorageBlobDeletedEvent(t *testing.T) {
	requestContent := "[{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.URL)
}

func TestConsumeStorageBlobRenamedEvent(t *testing.T) {
	requestContent := "[ {  \"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"eventType\": \"Microsoft.Storage.BlobRenamed\",  \"eventTime\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"RenameFile\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"destinationUrl\": \"https://myaccount.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobRenamedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/testfile.txt", *sysEvent.DestinationURL)
}

func TestConsumeStorageDirectoryCreatedEvent(t *testing.T) {
	requestContent := "[ {  \"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.DirectoryCreated\",  \"eventTime\": \"2017-08-16T01:57:26.005121Z\",  \"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",  \"data\": {    \"api\": \"CreateDirectory\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },  \"dataVersion\": \"2\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryCreatedEventData](t, events[0].Data)
	require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeStorageDirectoryDeletedEvent(t *testing.T) {
	requestContent := "[{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\", \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.DirectoryDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryDeletedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
	require.Nil(t, sysEvent.Recursive)
}

func TestConsumeStorageDirectoryDeletedEvent_Recursive(t *testing.T) {
	requestContent := "[{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",   \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.DirectoryDeleted\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": { \"recursive\":\"true\",   \"api\": \"DeleteDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryDeletedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
	require.Equal(t, "true", *sysEvent.Recursive)
}

func TestConsumeStorageDirectoryRenamedEvent(t *testing.T) {
	requestContent := "[{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.DirectoryRenamed\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"destinationUrl\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"1\",  \"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageDirectoryRenamedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.DestinationURL)
}

func TestConsumeStorageAsyncOperationInitiatedEvent(t *testing.T) {
	requestContent := "[{    \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.AsyncOperationInitiated\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"1.0\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageAsyncOperationInitiatedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeStorageBlobTierChangedEvent(t *testing.T) {
	requestContent := "[{   \"topic\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testDir\",  \"eventType\": \"Microsoft.Storage.BlobTierChanged\",  \"eventTime\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"RenameDirectory\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testDir\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },  \"dataVersion\": \"1.0\"}]"
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageBlobTierChangedEventData](t, events[0].Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testDir", *sysEvent.URL)
}

func TestConsumeStorageTaskQueuedEvent(t *testing.T) {
	requestContent := `[{
	"topic": "/subscriptions/c86a9c18-8373-41fa-92d4-1d7bdc16977b/resourceGroups/shulin-rg/providers/Microsoft.Storage/storageAccounts/shulinstcanest2",
	"subject": "DataManagement/StorageTasks",
	"eventType": "Microsoft.Storage.StorageTaskQueued",
	"id": "7fddaf06-24e8-4d57-9b66-5b7ab920a626",
	"data": {
		"queuedDateTime": "2023-03-23T16:43:50Z",
		"taskExecutionId": "deletetest-2023-03-23T16:42:33.8658256Z_2023-03-23T16:42:58.8983000Z"
	},
	"dataVersion": "1.0",
	"metadataVersion": "1",
	"eventTime": "2023-03-23T16:43:50Z"
}]`
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageTaskQueuedEventData](t, events[0].Data)
	require.Equal(t, "deletetest-2023-03-23T16:42:33.8658256Z_2023-03-23T16:42:58.8983000Z", *sysEvent.TaskExecutionID)
	require.Equal(t, mustParseTime(t, "2023-03-23T16:43:50Z"), *sysEvent.QueuedDateTime)
}

func TestConsumeStorageTaskCompletedEvent(t *testing.T) {
	requestContent := `[{
	"topic": "/subscriptions/c86a9c18-8373-41fa-92d4-1d7bdc16977b/resourceGroups/shulin-rg/providers/Microsoft.Storage/storageAccounts/shulinstcanest2",
	"subject": "DataManagement/StorageTasks",
	"eventType": "Microsoft.Storage.StorageTaskCompleted",
	"id": "7fddaf06-24e8-4d57-9b66-5b7ab920a626",
	"data": {
		"status": "Succeeded",
		"completedDateTime": "2023-03-23T16:52:58Z",
		"taskExecutionId": "deletetest-2023-03-23T16:42:33.8658256Z_2023-03-23T16:42:58.8983000Z",
		"taskName": "delete123",
		"summaryReportBlobUrl": "https://shulinstcanest2.blob.core.windows.net/report/delete123_deletetest_2023-03-23T16:43:50/SummaryReport.json"
	},
	"dataVersion": "1.0",
	"metadataVersion": "1",
	"eventTime": "2023-03-23T16:43:50Z"
}]`
	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.StorageTaskCompletedEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.StorageTaskCompletedStatusSucceeded, *sysEvent.Status)
	require.Equal(t, mustParseTime(t, "2023-03-23T16:52:58Z"), *sysEvent.CompletedDateTime)
	require.Equal(t, "deletetest-2023-03-23T16:42:33.8658256Z_2023-03-23T16:42:58.8983000Z", *sysEvent.TaskExecutionID)
	require.Equal(t, "delete123", *sysEvent.TaskName)
	require.Equal(t, "https://shulinstcanest2.blob.core.windows.net/report/delete123_deletetest_2023-03-23T16:43:50/SummaryReport.json", *sysEvent.SummaryReportBlobURL)
}

// App Service events
func TestConsumeWebAppUpdatedEvent(t *testing.T) {
	siteName := "testSite01"

	requestContent := `[{"topic": "/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01", "subject": "/Microsoft.Web/sites/testSite01","eventType": "Microsoft.Web.AppUpdated", "eventTime": "2017-08-16T01:57:26.005121Z","id": "602a88ef-0001-00e6-1233-1646070610ea","data": { "appEventTypeDetail": { "action": "Restarted"},"name": "testSite01","clientRequestId": "ce636635-2b81-4981-a9d4-cec28fb5b014","correlationRequestId": "61baa426-c91f-4e58-b9c6-d3852c4d88d","requestId": "0a4d5b5e-7147-482f-8e21-4219aaacf62a","address": "/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01","verb": "POST"},"dataVersion": "2","metadataVersion": "1"}]`

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebAppUpdatedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebBackupOperationStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.BackupOperationStarted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebBackupOperationCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.BackupOperationCompleted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebBackupOperationFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.BackupOperationFailed\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebBackupOperationFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebRestoreOperationStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.RestoreOperationStarted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebRestoreOperationCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.RestoreOperationCompleted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebRestoreOperationFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.RestoreOperationFailed\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebRestoreOperationFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebSlotSwapStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.SlotSwapStarted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebSlotSwapCompletedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.SlotSwapCompleted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapCompletedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebSlotSwapFailedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.SlotSwapFailed\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapFailedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebSlotSwapWithPreviewStartedEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.SlotSwapWithPreviewStarted\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapWithPreviewStartedEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebSlotSwapWithPreviewCancelledEvent(t *testing.T) {
	siteName := "testSite01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/sites/testSite01\", \"subject\": \"/Microsoft.Web/sites/testSite01\",\"eventType\": \"Microsoft.Web.SlotSwapWithPreviewCancelled\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appEventTypeDetail\": { \"action\": \"Restarted\"},\"name\": \"testSite01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebSlotSwapWithPreviewCancelledEventData](t, events[0].Data)
	require.Equal(t, siteName, *sysEvent.Name)
}

func TestConsumeWebAppServicePlanUpdatedEvent(t *testing.T) {
	planName := "testPlan01"
	requestContent := "[{\"topic\": \"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Web/serverfarms/testPlan01\", \"subject\": \"/Microsoft.Web/serverfarms/testPlan01\",\"eventType\": \"Microsoft.Web.AppServicePlanUpdated\", \"eventTime\": \"2017-08-16T01:57:26.005121Z\",\"id\": \"602a88ef-0001-00e6-1233-1646070610ea\",\"data\": { \"appServicePlanEventTypeDetail\": { \"stampKind\": \"Public\",\"action\": \"Updated\",\"status\": \"Started\" },\"name\": \"testPlan01\",\"clientRequestId\": \"ce636635-2b81-4981-a9d4-cec28fb5b014\",\"correlationRequestId\": \"61baa426-c91f-4e58-b9c6-d3852c4d88d\",\"requestId\": \"0a4d5b5e-7147-482f-8e21-4219aaacf62a\",\"address\": \"/subscriptions/ef90e930-9d7f-4a60-8a99-748e0eea69de/resourcegroups/egcanarytest/providers/Microsoft.Web/sites/egtestapp/restart?api-version=2016-03-01\",\"verb\": \"POST\"},\"dataVersion\": \"2\",\"metadataVersion\": \"1\"}]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.WebAppServicePlanUpdatedEventData](t, events[0].Data)
	require.Equal(t, planName, *sysEvent.Name)
}

// Policy Insights
func TestConsumePolicyInsightsPolicyStateChangedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"eventType\":\"Microsoft.PolicyInsights.PolicyStateChanged\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":\"{request-operation}\",        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateChangedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

func TestConsumePolicyInsightsPolicyStateCreatedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"eventType\":\"Microsoft.PolicyInsights.PolicyStateCreated\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":\"{request-operation}\",        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateCreatedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

func TestConsumePolicyInsightsPolicyStateDeletedEvent(t *testing.T) {
	requestContent := "[   {     \"source\":\"/subscriptions/{subscription-id}\",     \"subject\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"eventType\":\"Microsoft.PolicyInsights.PolicyStateDeleted\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"timestamp\":\"2017-08-16T03:54:38.2696833Z\",  \"policyDefinitionId\":\"4c2359fe-001e-00ba-0e04-585868000000\",       \"httpRequest\":\"{request-operation}\",        \"resourceProvider\":\"Microsoft.EventGrid\",        \"resourceUri\":\"/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.EventGrid/eventSubscriptions/LogicAppdd584bdf-8347-49c9-b9a9-d1f980783501\",    \"subscriptionId\":\"{subscription-id}\"   },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.PolicyInsightsPolicyStateDeletedEventData](t, events[0].Data)
	require.Equal(t, "4c2359fe-001e-00ba-0e04-585868000000", *sysEvent.PolicyDefinitionID)
}

// Communication events
func TestConsumeAcsRecordingFileStatusUpdatedEventData(t *testing.T) {
	requestContent := "[   {      \"subject\":\"/recording/call/{call-id}/recordingId/{recording-id}\",    \"eventType\":\"Microsoft.Communication.RecordingFileStatusUpdated\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"recordingStorageInfo\": { \"recordingChunks\": [ { \"documentId\": \"0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8\", \"index\": 0, \"endReason\": \"SessionEnded\", \"contentLocation\": \"https://storage.asm.skype.com/v1/objects/0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8/content/video\", \"metadataLocation\": \"https://storage.asm.skype.com/v1/objects/0-eus-d12-801b3f3fc462fe8a01e6810cbff729b8/content/acsmetadata\" }]}, \"recordingChannelType\": \"Mixed\", \"recordingContentType\": \"Audio\", \"recordingFormatType\": \"Mp3\"},   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	sysEvent := deserializeSystemEvent[azsystemevents.ACSRecordingFileStatusUpdatedEventData](t, events[0].Data)

	require.Equal(t, azsystemevents.RecordingChannelKindMixed, *sysEvent.RecordingChannelKind)
	require.Equal(t, azsystemevents.RecordingContentTypeAudio, *sysEvent.RecordingContentType)
	require.Equal(t, azsystemevents.RecordingFormatTypeMp3, *sysEvent.RecordingFormatType)
}

func TestConsumeAcsEmailDeliveryReportReceivedEvent(t *testing.T) {
	// TODO: the enum value here for 'status' used to be 'delivered'. I'm not sure if this is just
	// a test artifact, or if the value really isn't cased like the enum value is ('Delivered').
	requestContent := `{
		"id": "5f04f77c-2a6a-43bd-9b74-576a64c01f9e",
		"source": "source",
		"specversion": "1.0",
		"type": "type",
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "sender/test2@contoso.org/message/950850f5-bcdf-4315-b77a-6447cf56fac9",
		"data": {
			"sender": "test2@contoso.org",
			"recipient": "test1@contoso.com",
			"messageId": "950850f5-bcdf-4315-b77a-6447cf56fac9",
			"status": "Delivered",
			"deliveryAttemptTimeStamp": "2023-02-09T19:46:12.2480265+00:00",
			"deliveryStatusDetails": {
				"statusMessage": "DestinationMailboxFull"
			}
		},
		"eventType": "Microsoft.Communication.EmailDeliveryReportReceived",
		"dataVersion": "1.0",
		"metadataVersion": "1",
		"eventTime": "2023-02-09T19:46:12.2478002Z"
	}`

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	emailEvent := deserializeSystemEvent[azsystemevents.ACSEmailDeliveryReportReceivedEventData](t, event.Data)
	require.Equal(t, "test2@contoso.org", *emailEvent.Sender)
	require.Equal(t, "test1@contoso.com", *emailEvent.Recipient)
	require.Equal(t, azsystemevents.ACSEmailDeliveryReportStatusDelivered, *emailEvent.Status)
	require.Equal(t, "DestinationMailboxFull", *emailEvent.DeliveryStatusDetails.StatusMessage)
	require.Equal(t, mustParseTime(t, "2023-02-09T19:46:12.2480265+00:00"), *emailEvent.DeliveryAttemptTimestamp)
}

func TestConsumeAcsIncomingCallEvent(t *testing.T) {
	requestContent := `{
		"id": "e80026e7-e298-46ba-bc42-dab0eda92581",
		"topic": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
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
		"eventType": "Microsoft.Communication.IncomingCall",
		"dataVersion": "1.0",
		"metadataVersion": "1",
		"eventTime": "2023-04-04T17:18:42.5542219Z"
	}`

	event := parseEvent(t, requestContent)

	incomingCallEvent := deserializeSystemEvent[azsystemevents.ACSIncomingCallEventData](t, event.Data)
	require.Equal(t, "{recipient-id}", *incomingCallEvent.ToCommunicationIdentifier.CommunicationUser.ID)
	require.Equal(t, "{caller-id}", *incomingCallEvent.FromCommunicationIdentifier.CommunicationUser.ID)
	require.Equal(t, "VOIP Caller", *incomingCallEvent.CallerDisplayName)
	require.Equal(t, "616d617a6f6e5f6368696;encoding=hex", *incomingCallEvent.CustomContext.SipHeaders["userToUser"])
	require.Equal(t, "35567842", *incomingCallEvent.CustomContext.SipHeaders["X-MS-Custom-myheader1"])
	require.Equal(t, "customsipheadervalue", *incomingCallEvent.CustomContext.SipHeaders["X-MS-Custom-myheader2"])
	require.Equal(t, "customValue", *incomingCallEvent.CustomContext.VoipHeaders["customHeader"])
	require.Equal(t, "{incoming-call-contextValue}", *incomingCallEvent.IncomingCallContext)
	require.Equal(t, "correlationId", *incomingCallEvent.CorrelationID)
}

// Health Data Services events
func TestConsumeFhirResourceCreatedEvent(t *testing.T) {
	requestContent := "[   {  \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"eventType\":\"Microsoft.HealthcareApis.FhirResourceCreated\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceCreatedEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeFhirResourceUpdatedEvent(t *testing.T) {
	requestContent := "[   {  \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"eventType\":\"Microsoft.HealthcareApis.FhirResourceUpdated\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceUpdatedEventData](t, events[0].Data)

	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeFhirResourceDeletedEvent(t *testing.T) {
	requestContent := "[   {  \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"eventType\":\"Microsoft.HealthcareApis.FhirResourceDeleted\",    \"eventTime\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"dataVersion\": \"1.0\"  }]"

	events := parseManyEvents(t, requestContent)

	require.NotEmpty(t, events)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceDeletedEventData](t, events[0].Data)

	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeDicomImageCreatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"eventType": "Microsoft.HealthcareApis.DicomImageCreated",
	"dataVersion": "1",
	"metadataVersion": "1",
	"eventTime": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specVersion": "1.0"
}`
	event := parseEvent(t, requestContent)

	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageCreatedEventData](t, event.Data)

	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

func TestConsumeDicomImageUpdatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"eventType": "Microsoft.HealthcareApis.DicomImageUpdated",
	"dataVersion": "1",
	"metadataVersion": "1",
	"eventTime": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specVersion": "1.0"
}`
	event := parseEvent(t, requestContent)

	require.NotEmpty(t, event)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageUpdatedEventData](t, event.Data)

	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

func TestConsumeDicomImageDeletedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"eventType": "Microsoft.HealthcareApis.DicomImageDeleted",
	"dataVersion": "1",
	"metadataVersion": "1",
	"eventTime": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specVersion": "1.0"
}`
	event := parseEvent(t, requestContent)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageDeletedEventData](t, event.Data)

	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

// CloudEvent tests

// Miscellaneous tests

func TestParsesCloudEventEnvelope(t *testing.T) {
	requestContent := "[{\"key\": \"value\",  \"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"source\":\"Subject-0\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",   \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  }, \"type\":\"Microsoft.Storage.BlobDeleted\",\"specversion\":\"1.0\", \"dataschema\":\"1.0\", \"subject\":\"subject\", \"datacontenttype\": \"text/plain\", \"time\": \"2017-08-16T01:57:26.005121Z\"}]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	var cloudEvent = events[0]

	require.Equal(t, "994bc3f8-c90c-6fc3-9e83-6783db2221d5", cloudEvent.ID)
	require.Equal(t, "Subject-0", cloudEvent.Source)
	require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), cloudEvent.Type)
	require.Equal(t, "text/plain", *cloudEvent.DataContentType)
	require.Equal(t, "subject", *cloudEvent.Subject)
	require.Equal(t, "1.0", *cloudEvent.DataSchema)
	require.Equal(t, mustParseTime(t, "2017-08-16T01:57:26.005121Z"), *cloudEvent.Time)
	require.Equal(t, "value", cloudEvent.Extensions["key"])
}

func TestConsumeCloudEventsWithAdditionalProperties(t *testing.T) {
	requestContent := "[{\"key\": \"value\",  \"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"source\":\"Subject-0\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  }, \"type\":\"Microsoft.Storage.BlobDeleted\",\"specversion\":\"1.0\"}]"

	events := parseManyCloudEvents(t, requestContent)
	require.NotEmpty(t, events)

	if events[0].Type == string(azsystemevents.TypeStorageBlobDeleted) {
		eventData := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, events[0].Data)
		require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *eventData.URL)
	}

	require.Equal(t, "value", events[0].Extensions["key"])
}

func TestConsumeCloudEventUsingBinaryDataExtensionMethod(t *testing.T) {
	messageBody := "{  \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\", \"specversion\": \"1.0\", \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"type\": \"Microsoft.Storage.BlobDeleted\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  }}"

	cloudEvent := parseCloudEvent(t, messageBody)

	require.NotEmpty(t, cloudEvent)
	require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), cloudEvent.Type)
	blobDeleted := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, cloudEvent.Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *blobDeleted.URL)
}

func TestConsumeCloudEventNotWrappedInAnArray(t *testing.T) {
	requestContent := "{  \"source\": \"/subscriptions/id/resourceGroups/Storage/providers/Microsoft.Storage/storageAccounts/xstoretestaccount\",  \"subject\": \"/blobServices/default/containers/testcontainer/blobs/testfile.txt\",  \"specversion\": \"1.0\", \"type\": \"Microsoft.Storage.BlobDeleted\",  \"time\": \"2017-11-07T20:09:22.5674003Z\",  \"id\": \"4c2359fe-001e-00ba-0e04-58586806d298\",  \"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",   \"brandNewProperty\": \"0000000000000281000000000002F5CA\", \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  }}"

	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	require.Equal(t, string(azsystemevents.TypeStorageBlobDeleted), event.Type)
	eventData := deserializeSystemEvent[azsystemevents.StorageBlobDeletedEventData](t, event.Data)
	require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *eventData.URL)
}
