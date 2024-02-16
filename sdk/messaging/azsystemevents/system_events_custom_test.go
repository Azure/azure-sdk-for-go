//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents_test

// func TestConsumeMultipleCloudEventsInSameBatch(t *testing.T) {
// 	requestContent := "[" +
// 		"{\"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"source\":\"Subject-0\",\"data\": {    \"api\": \"PutBlockList\",    \"clientRequestId\": \"799304a4-bbc5-45b6-9849-ec2c66be800a\",    \"requestId\": \"602a88ef-0001-00e6-1233-164607000000\",    \"eTag\": \"0x8D4E44A24ABE7F1\",    \"contentType\": \"text/plain\",    \"contentLength\": 447,    \"blobType\": \"BlockBlob\",    \"url\": \"https://myaccount.blob.core.windows.net/testcontainer/file1.txt\",    \"sequencer\": \"00000000000000EB000000000000C65A\"  },\"type\":\"Microsoft.Storage.BlobCreated\",\"specversion\":\"1.0\"}," +
// 		"{\"id\":\"2947780a-356b-c5a5-feb4-f5261fb2f155\",\"source\":\"Subject-1\",\"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },\"type\":\"Microsoft.Storage.BlobDeleted\",\"specversion\":\"1.0\"}," +
// 		"{\"id\":\"cb14e05b-50c6-67dc-cafa-f4bcff3bf520\",\"source\":\"Subject-2\",\"data\": {    \"api\": \"DeleteBlob\",    \"requestId\": \"4c2359fe-001e-00ba-0e04-585868000000\",    \"contentType\": \"text/plain\",    \"blobType\": \"BlockBlob\",    \"url\": \"https://example.blob.core.windows.net/testcontainer/testfile.txt\",    \"sequencer\": \"0000000000000281000000000002F5CA\",    \"storageDiagnostics\": {      \"batchId\": \"b68529f3-68cd-4744-baa4-3c0498ec19f0\"    }  },\"type\":\"Microsoft.Storage.BlobDeleted\",\"specversion\":\"1.0\"}," +
// 		"{\"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"source\":\"Subject-0\",\"data_base64\": \"ZGF0YQ==\",\"type\":\"BinaryDataType\",\"specversion\":\"1.0\"}," +
// 		"{\"id\":\"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",\"source\":\"/contoso/items\",\"subject\": \"\",\"data\": {    \"itemSku\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"itemUri\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"type\": \"Contoso.Items.ItemReceived\",\"specversion\":\"1.0\"}]"

// 	ObjectSerializer camelCaseSerializer = new JsonObjectSerializer(
// 		new JsonSerializerOptions()
// 		{
// 			PropertyNamingPolicy = JsonNamingPolicy.CamelCase
// 		})

// 	events := parseManyCloudEvents(t, requestContent)

// 	require.NotEmpty(t, events)
// 	require.Equal(t, 5, events.Length)
// 	foreach (CloudEvent cloudEvent in events)
// 	{
// 		if (cloudEvent.TryGetSystemEventData(out object eventData))
// 		{
// 			switch (eventData)
// 			{
// 				case StorageBlobCreatedEventData blobCreated:
// 					require.Equal(t, "https://myaccount.blob.core.windows.net/testcontainer/file1.txt", blobCreated.Url)
// 				case StorageBlobDeletedEventData blobDeleted:
// 					require.Equal(t, "https://example.blob.core.windows.net/testcontainer/testfile.txt", *blobDeleted.URL)
// 			}
// 		}
// 		else
// 		{
// 			switch (cloudEvent.Type)
// 			{
// 				case "BinaryDataType":
// 					require.Equal(t, Convert.ToBase64String(cloudEvent.Data.ToArray()), "ZGF0YQ==")
// 					Assert.IsFalse(cloudEvent.TryGetSystemEventData(out var _))
// 				case "Contoso.Items.ItemReceived":
// 					ContosoItemReceivedEventData itemReceived = cloudEvent.Data.ToObject<ContosoItemReceivedEventData>(camelCaseSerializer)
// 					require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", itemReceived.ItemSku)
// 					Assert.IsFalse(cloudEvent.TryGetSystemEventData(out var _))
// 			}
// 		}
// 	}
// }

// Custom event tests
// func TestConsumeCloudEventWithBinaryDataPayload(t *testing.T) {
// 	requestContent := "[{\"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"source\":\"Subject-0\",  \"data_base64\": \"ZGF0YQ==\", \"type\":\"Test.Items.BinaryDataType\",\"specversion\":\"1.0\"}]"

// 	events := parseManyCloudEvents(t, requestContent)
// 	if events[0].Type == "Test.Items.BinaryDataType"{
// 		var eventData = (BinaryData)
// 		require.Equal(t, "data", events[0].Data)
// 	}
// }

// func TestConsumeCloudEventWithCustomEventPayload(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\":\"/contoso/items\",  \"subject\": \"\",  \"data\": {    \"itemSku\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"itemUri\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"type\": \"Contoso.Items.ItemReceived\", \"specversion\": \"1.0\"}]"

// 	ObjectSerializer camelCaseSerializer = new JsonObjectSerializer(
// 		new JsonSerializerOptions()
// 		{
// 			PropertyNamingPolicy = JsonNamingPolicy.CamelCase
// 		})

// 	events := parseManyCloudEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	ContosoItemReceivedEventData eventData = events[0].Data.ToObject<ContosoItemReceivedEventData>(camelCaseSerializer)
// 	require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", eventData.ItemSku)
// }

// func TestConsumeCloudEventWithArrayDataPayload(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\":\"/contoso/items\", \"subject\": \"\",  \"data\": [{    \"itemSku\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"itemUri\": \"https://rp-eastus2.eventgrid.azure.net:553\"  }],  \"type\": \"Contoso.Items.ItemReceived\", \"specversion\": \"1.0\"}]"

// 	ObjectSerializer camelCaseSerializer = new JsonObjectSerializer(
// 		new JsonSerializerOptions()
// 		{
// 			PropertyNamingPolicy = JsonNamingPolicy.CamelCase
// 		})

// 	events := parseManyCloudEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	ContosoItemReceivedEventData[] eventData = events[0].Data.ToObject<ContosoItemReceivedEventData[]>(camelCaseSerializer)
// 	require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", eventData[0].ItemSku)
// }

// // Null data tests
// func TestConsumeCloudEventWithNoData(t *testing.T) {
// 	requestContent := "[{\"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\",\"type\":\"type\",\"source\":\"Subject-0\",\"specversion\":\"1.0\"}]"

// 	events := parseManyCloudEvents(t, requestContent)
// 	var eventData1 = events[0].Data

// 	require.Equal(t, eventData1, null)
// 	require.Equal(t, "type", events[0].Type)
// }

// func TestConsumeCloudEventWithExplicitlyNullData(t *testing.T) {
// 	requestContent := "[{\"id\":\"994bc3f8-c90c-6fc3-9e83-6783db2221d5\", \"type\":\"type\", \"source\":\"Subject-0\", \"data\":null, \"specversion\":\"1.0\"}]"

// 	events := parseManyCloudEvents(t, requestContent)
// 	Assert.IsNull(events[0].Data.ToObjectFromJson<object>())
// 	require.Equal(t, "type", events[0].Type)
// }

// // Primitive/string data tests
// func TestConsumeCloudEventWithBooleanData(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\":\"/contoso/items\",  \"subject\": \"\",  \"data\": true,  \"type\": \"Contoso.Items.ItemReceived\",  \"time\": \"2018-01-25T22:12:19.4556811Z\", \"specversion\": \"1.0\"}]"

// 	events := parseManyCloudEvents(t, requestContent)

// 	require.NotEmpty(t, events)
// 	BinaryData binaryEventData = events[0].Data
// 	bool eventData = binaryEventData.ToObjectFromJson<bool>()
// 	require.True(t, eventData)
// }

// func TestConsumeCloudEventWithStringData(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"source\":\"/contoso/items\",  \"subject\": \"\",  \"data\": \"stringdata\",  \"type\": \"Contoso.Items.ItemReceived\",  \"time\": \"2018-01-25T22:12:19.4556811Z\",  \"specversion\": \"1.0\"}]"
// 	events := parseManyCloudEvents(t, requestContent)

// 	require.NotEmpty(t, events)
// 	BinaryData binaryEventData = events[0].Data
// 	string eventData = binaryEventData.ToObjectFromJson<string>()
// 	require.Equal(t, "stringdata", eventData)
// }

// //
// // Custom event tests
// //

// func TestConsumeCustomEvents(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": {    \"itemSku\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"itemUri\": \"https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d\"  },  \"eventType\": \"Contoso.Items.ItemReceived\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]";

// 	events := parseManyEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	for _, egEvent := range events {
// 	{
// 		if (egEvent.EventType == "Contoso.Items.ItemReceived")
// 		{
// 			ContosoItemReceivedEventData eventData = egEvent.Data.ToObject<ContosoItemReceivedEventData>(new JsonObjectSerializer(
// 				new JsonSerializerOptions()
// 				{
// 					PropertyNamingPolicy = JsonNamingPolicy.CamelCase
// 				}));
// 			require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", eventData.ItemSku);
// 		}
// 	}
// }

// func TestConsumeCustomEventWithArrayData(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": [{    \"itemSku\": \"512d38b6-c7b8-40c8-89fe-f46f9e9622b6\",    \"itemUri\": \"https://rp-eastus2.eventgrid.azure.net:553\"  }],  \"eventType\": \"Contoso.Items.ItemReceived\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]";

// 	events := parseManyEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	for _, egEvent := range events {
// 	{
// 		if (egEvent.EventType == "Contoso.Items.ItemReceived")
// 		{
// 			ContosoItemReceivedEventData[] eventData = egEvent.Data.ToObject<ContosoItemReceivedEventData[]>(new JsonObjectSerializer(
// 				new JsonSerializerOptions()
// 				{
// 					PropertyNamingPolicy = JsonNamingPolicy.CamelCase
// 				}));
// 			require.Equal(t, "512d38b6-c7b8-40c8-89fe-f46f9e9622b6", eventData[0].ItemSku);
// 		}
// 	}
// }

// //
// // Primitive/string data tests
// //
// func TestConsumeCustomEventWithBooleanData(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": true,  \"eventType\": \"Contoso.Items.ItemReceived\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]";

// 	events := parseManyEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	for _, egEvent := range events {
// 	{
// 		if (egEvent.EventType == "Contoso.Items.ItemReceived")
// 		{
// 			BinaryData binaryEventData = egEvent.Data;
// 			require.True(t, binaryEventData.ToObjectFromJson<bool>());
// 		}
// 	}
// }

// func TestConsumeCustomEventWithStringData(t *testing.T) {
// 	requestContent := "[{  \"id\": \"2d1781af-3a4c-4d7c-bd0c-e34b19da4e66\",  \"topic\": \"/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\",  \"subject\": \"\",  \"data\": \"stringdata\",  \"eventType\": \"Contoso.Items.ItemReceived\",  \"eventTime\": \"2018-01-25T22:12:19.4556811Z\",  \"metadataVersion\": \"1\",  \"dataVersion\": \"1\"}]";

// 	events := parseManyEvents(t, requestContent)

// 	require.NotEmpty(t, events)

// 	for _, egEvent := range events {
// 	{
// 		if (egEvent.EventType == "Contoso.Items.ItemReceived")
// 		{
// 			BinaryData binaryEventData = egEvent.Data;
// 			require.Equal(t, "stringdata", binaryEventData.ToObjectFromJson<string>());
// 		}
// 	}
// }
// #endregion
