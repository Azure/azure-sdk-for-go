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

func TestConsumeCloudEventAcsRouterJobQueuedEvent(t *testing.T) {
	requestContent := `{
		"id": "b6d8687a-5a1a-42ae-b8b5-ff7ec338c872",
		"source": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "job/{job-id}/channel/{channel-id}/queue/{queue-id}",
		"data": {
		"jobId": "7f1df17b-570b-4ae5-9cf5-fe6ff64cc712",
		"channelReference": "test-abc",
		"channelId": "FooVoiceChannelId",
		"queueId": "625fec06-ab81-4e60-b780-f364ed96ade1",
		"priority": 1,
		"labels": {
			"Locale": "en-us",
			"Segment": "Enterprise",
			"Token": "FooToken"
		},
		"tags": {
			"Locale": "en-us",
			"Segment": "Enterprise",
			"Token": "FooToken"
		},
		"requestedWorkerSelectors": [
			{
			"key": "string",
			"labelOperator": "Equal",
			"value": 5,
			"ttlSeconds": 1000
			}
		],
		"attachedWorkerSelectors": [
			{
			"key": "string",
			"labelOperator": "Equal",
			"value": 5,
			"ttlSeconds": 1000,
			"state": "active"
			}
		]
		},
		"type": "Microsoft.Communication.RouterJobQueued",
		"specversion": "1.0",
		"time": "2022-02-17T00:55:25.1736293Z"
	}`

	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.AcsRouterJobQueuedEventData](t, event.Data)

	var selectors = sysEvent.AttachedWorkerSelectors
	require.Equal(t, 1, len(selectors))
	require.Equal(t, float32(1000), *selectors[0].TimeToLive)

	require.Equal(t, azsystemevents.AcsRouterLabelOperatorEqual, *selectors[0].LabelOperator)
	// TODO: might have been a field rename?
	//require.Equal(t, azsystemevents.AcsRouterLabelOperatorEqual, selectors[0].Operator)

	require.Equal(t, azsystemevents.AcsRouterWorkerSelectorStateActive, *selectors[0].State)
	// TODO: might have been a field rename?
	//require.Equal(t, azsystemevents.AcsRouterWorkerSelectorStateActive, selectors[0].SelectorState)
}

func TestConsumeCloudEventAcsRouterJobReceivedEvent(t *testing.T) {
	requestContent := `{
		"id": "acdf8fa5-8ab4-4a65-874a-c1d2a4a97f2e",
		"source": "/subscriptions/{subscription-id}/resourceGroups/{group-name}/providers/Microsoft.Communication/communicationServices/{communication-services-resource-name}",
		"subject": "job/{job-id}/channel/{channel-id}",
		"data": {
		"jobId": "7f1df17b-570b-4ae5-9cf5-fe6ff64cc712",
		"channelReference": "test-abc",
		"jobStatus": "PendingClassification",
		"channelId": "FooVoiceChannelId",
		"classificationPolicyId": "test-policy",
		"queueId": "queue-id",
		"priority": 0,
		"labels": {
			"Locale": "en-us",
			"Segment": "Enterprise",
			"Token": "FooToken"
		},
		"tags": {
			"Locale": "en-us",
			"Segment": "Enterprise",
			"Token": "FooToken"
		},
		"requestedWorkerSelectors": [
			{
			"key": "string",
			"labelOperator": "equal",
			"value": 5,
			"ttlSeconds": 36
			}
		],
		"scheduledOn": "2007-03-28T19:13:50+00:00",
		"unavailableForMatching": false
		},
		"type": "Microsoft.Communication.RouterJobReceived",
		"specversion": "1.0",
		"time": "2022-02-17T00:55:25.1736293Z"
	}`

	// TODO: formatting for the time is different - not ISO.
	// "scheduledOn": "3/28/2007 7:13:50 PM +00:00",

	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.AcsRouterJobReceivedEventData](t, event.Data)
	require.Equal(t, azsystemevents.AcsRouterJobStatusPendingClassification, *sysEvent.JobStatus)

	// TODO: don't have a .Status field?
	//require.Equal(t, Azure.Messaging.EventGrid.azsystemevents.AcsRouterJobStatus.PendingClassification, sysEvent.Status)
}

// Health Data Services events
func TestConsumeCloudEventFhirResourceCreatedEvent(t *testing.T) {
	requestContent := "[   { \"source\": \"/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}\", \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"type\":\"Microsoft.HealthcareApis.FhirResourceCreated\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceCreatedEventData](t, events[0].Data)
	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeCloudEventFhirResourceUpdatedEvent(t *testing.T) {
	requestContent := "[   { \"source\": \"/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}\", \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"type\":\"Microsoft.HealthcareApis.FhirResourceUpdated\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	require.NotEmpty(t, events)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceUpdatedEventData](t, events[0].Data)

	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeCloudEventFhirResourceDeletedEvent(t *testing.T) {
	requestContent := "[   { \"source\": \"/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}\", \"subject\":\"{fhir-account}.fhir.azurehealthcareapis.com/Patient/e0a1f743-1a70-451f-830e-e96477163902\",    \"type\":\"Microsoft.HealthcareApis.FhirResourceDeleted\",    \"time\":\"2017-08-16T03:54:38.2696833Z\",    \"id\":\"25b3b0d0-d79b-44d5-9963-440d4e6a9bba\",    \"data\": { \"resourceType\": \"Patient\",  \"resourceFhirAccount\": \"{fhir-account}.fhir.azurehealthcareapis.com\", \"resourceFhirId\": \"e0a1f743-1a70-451f-830e-e96477163902\", \"resourceVersionId\": 1 },   \"specversion\": \"1.0\"  }]"

	events := parseManyCloudEvents(t, requestContent)

	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareFhirResourceDeletedEventData](t, events[0].Data)

	require.NotNil(t, healthEvent)
	require.Equal(t, azsystemevents.HealthcareFhirResourceTypePatient, *healthEvent.FhirResourceType)
	require.Equal(t, "{fhir-account}.fhir.azurehealthcareapis.com", *healthEvent.FhirServiceHostName)
	require.Equal(t, "e0a1f743-1a70-451f-830e-e96477163902", *healthEvent.FhirResourceID)
	require.Equal(t, int64(1), *healthEvent.FhirResourceVersionID)
}

func TestConsumeCloudEventDicomImageCreatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"type": "Microsoft.HealthcareApis.DicomImageCreated",
	"time": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specversion": "1.0"
}`
	event := parseCloudEvent(t, requestContent)

	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageCreatedEventData](t, event.Data)

	require.NotNil(t, healthEvent)
	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

func TestConsumeCloudEventDicomImageUpdatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"type": "Microsoft.HealthcareApis.DicomImageUpdated",
	"time": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specversion": "1.0"
}`
	event := parseCloudEvent(t, requestContent)

	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageUpdatedEventData](t, event.Data)

	require.NotNil(t, healthEvent)
	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

func TestConsumeCloudEventDicomImageDeletedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.HealthcareApis/workspaces/{workspace-name}",
	"subject": "{dicom-account}.dicom.azurehealthcareapis.com/v1/studies/1.2.3.4.3/series/1.2.3.4.3.9423673/instances/1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
	"type": "Microsoft.HealthcareApis.DicomImageDeleted",
	"time": "2022-09-15T01:14:04.5613214Z",
	"id": "d621839d-958b-4142-a638-bb966b4f7dfd",
	"data": {
		"partitionName": "Microsoft.Default",
		"imageStudyInstanceUid": "1.2.3.4.3",
		"imageSeriesInstanceUid": "1.2.3.4.3.9423673",
		"imageSopInstanceUid": "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442",
		"serviceHostName": "{dicom-account}.dicom.azurehealthcareapis.com",
		"sequenceNumber": 1
	},
	"specversion": "1.0"
}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	healthEvent := deserializeSystemEvent[azsystemevents.HealthcareDicomImageDeletedEventData](t, event.Data)
	require.NotNil(t, healthEvent)
	require.Equal(t, "1.2.3.4.3", *healthEvent.ImageStudyInstanceUID)
	require.Equal(t, "1.2.3.4.3.9423673", *healthEvent.ImageSeriesInstanceUID)
	require.Equal(t, "1.3.6.1.4.1.45096.2.296485376.2210.1633373143.864442", *healthEvent.ImageSopInstanceUID)
	require.Equal(t, int64(1), *healthEvent.SequenceNumber)
	require.Equal(t, "Microsoft.Default", *healthEvent.PartitionName)
}

// APIM
func TestConsumeCloudEventGatewayApiAddedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}",
	"subject": "/gateways/{gateway-name}/apis/example-api",
	"type": "Microsoft.ApiManagement.GatewayAPIAdded",
	"time": "2021-07-02T00:47:47.8536532Z",
	"id": "92c502f2-a966-42a7-a428-d3b319844544",
	"data": {
		"resourceUri": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api"
	},
	"specversion": "1.0"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	apimEvent := deserializeSystemEvent[azsystemevents.APIManagementGatewayAPIAddedEventData](t, event.Data)
	require.Equal(t, "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api", *apimEvent.ResourceURI)
}

func TestConsumeCloudEventGatewayApiRemovedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}",
	"subject": "/gateways/{gateway-name}/apis/example-api",
	"type": "Microsoft.ApiManagement.GatewayAPIRemoved",
	"time": "2021-07-02T00:47:47.8536532Z",
	"id": "92c502f2-a966-42a7-a428-d3b319844544",
	"data": {
		"resourceUri": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api"
	},
	"specversion": "1.0"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	apimEvent := deserializeSystemEvent[azsystemevents.APIManagementGatewayAPIRemovedEventData](t, event.Data)
	require.Equal(t, "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api", *apimEvent.ResourceURI)
}

func TestConsumeCloudEventCertificateAuthorityCreatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}",
	"subject": "/gateways/{gateway-name}/apis/example-api",
	"type": "Microsoft.ApiManagement.GatewayCertificateAuthorityCreated",
	"time": "2021-07-02T00:47:47.8536532Z",
	"id": "92c502f2-a966-42a7-a428-d3b319844544",
	"data": {
		"resourceUri": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api"
	},
	"specversion": "1.0"
	}`
	event := parseCloudEvent(t, requestContent)

	apimEvent := deserializeSystemEvent[azsystemevents.APIManagementGatewayCertificateAuthorityCreatedEventData](t, event.Data)
	require.Equal(t, "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api", *apimEvent.ResourceURI)
}

func TestConsumeCloudEventCertificateAuthorityDeletedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}",
	"subject": "/gateways/{gateway-name}/apis/example-api",
	"type": "Microsoft.ApiManagement.GatewayCertificateAuthorityDeleted",
	"time": "2021-07-02T00:47:47.8536532Z",
	"id": "92c502f2-a966-42a7-a428-d3b319844544",
	"data": {
		"resourceUri": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api"
	},
	"specversion": "1.0"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	apimEvent := deserializeSystemEvent[azsystemevents.APIManagementGatewayCertificateAuthorityDeletedEventData](t, event.Data)
	require.Equal(t, "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api", *apimEvent.ResourceURI)
}

func TestConsumeCloudEventCertificateAuthorityUpdatedEvent(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}",
	"subject": "/gateways/{gateway-name}/apis/example-api",
	"type": "Microsoft.ApiManagement.GatewayCertificateAuthorityUpdated",
	"time": "2021-07-02T00:47:47.8536532Z",
	"id": "92c502f2-a966-42a7-a428-d3b319844544",
	"data": {
		"resourceUri": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api"
	},
	"specversion": "1.0"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	apimEvent := deserializeSystemEvent[azsystemevents.APIManagementGatewayCertificateAuthorityUpdatedEventData](t, event.Data)
	require.Equal(t, "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.ApiManagement/service/{your-APIM-instance}/gateways/{gateway-name}/apis/example-api", *apimEvent.ResourceURI)
}

// DataBox

func TestConsumeCloudEventDataBoxCopyCompleted(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.DataBox/jobs/{your-resource}",
	"subject": "/jobs/{your-resource}",
	"type": "Microsoft.DataBox.CopyCompleted",
	"time": "2022-10-16T02:51:26.4248221Z",
	"id": "759c892a-a628-4e48-a116-2e1d54c555ce",
	"data": {
		"serialNumber": "SampleSerialNumber",
		"stageName": "CopyCompleted",
		"stageTime": "2022-10-12T19:38:08.0218897Z"
	},
	"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	dataBoxEvent := deserializeSystemEvent[azsystemevents.DataBoxCopyCompletedEventData](t, event.Data)

	require.Equal(t, "SampleSerialNumber", *dataBoxEvent.SerialNumber)
	require.Equal(t, azsystemevents.DataBoxStageNameCopyCompleted, *dataBoxEvent.StageName)
	require.Equal(t, mustParseTime(t, "2022-10-12T19:38:08.0218897Z"), *dataBoxEvent.StageTime)
}

func TestConsumeCloudEventDataBoxCopyStarted(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.DataBox/jobs/{your-resource}",
	"subject": "/jobs/{your-resource}",
	"type": "Microsoft.DataBox.CopyStarted",
	"time": "2022-10-16T02:51:26.4248221Z",
	"id": "759c892a-a628-4e48-a116-2e1d54c555ce",
	"data": {
		"serialNumber": "SampleSerialNumber",
		"stageName": "CopyStarted",
		"stageTime": "2022-10-12T19:38:08.0218897Z"
	},
	"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	dataBoxEvent := deserializeSystemEvent[azsystemevents.DataBoxCopyStartedEventData](t, event.Data)

	require.Equal(t, "SampleSerialNumber", *dataBoxEvent.SerialNumber)
	require.Equal(t, azsystemevents.DataBoxStageNameCopyStarted, *dataBoxEvent.StageName)
	require.Equal(t, mustParseTime(t, "2022-10-12T19:38:08.0218897Z"), *dataBoxEvent.StageTime)
}

func TestConsumeCloudEventDataBoxOrderCompleted(t *testing.T) {
	requestContent := `{
	"source": "/subscriptions/{subscription-id}/resourceGroups/{your-rg}/providers/Microsoft.DataBox/jobs/{your-resource}",
	"subject": "/jobs/{your-resource}",
	"type": "Microsoft.DataBox.OrderCompleted",
	"time": "2022-10-16T02:51:26.4248221Z",
	"id": "759c892a-a628-4e48-a116-2e1d54c555ce",
	"data": {
		"serialNumber": "SampleSerialNumber",
		"stageName": "OrderCompleted",
		"stageTime": "2022-10-12T19:38:08.0218897Z"
	},
	"specversion": "1.0"
	}`

	event := parseCloudEvent(t, requestContent)

	dataBoxEvent := deserializeSystemEvent[azsystemevents.DataBoxOrderCompletedEventData](t, event.Data)

	require.Equal(t, "SampleSerialNumber", *dataBoxEvent.SerialNumber)
	require.Equal(t, azsystemevents.DataBoxStageNameOrderCompleted, *dataBoxEvent.StageName)
	require.Equal(t, mustParseTime(t, "2022-10-12T19:38:08.0218897Z"), *dataBoxEvent.StageTime)
}

// Resource Notifications

func TestConsumeCloudEventHealthResourcesAvailiabilityStatusChangedEvent(t *testing.T) {
	requestContent := `{
		"id": "1fb6fa94-d965-4306-abeq-4810f0774e97",
		"source": "/subscriptions/{subscription-id}",
		"subject": "/subscriptions/{subscription-id}/resourceGroups/{rg-name}/providers/Microsoft.Compute/virtualMachines/{vm-name}",
		"data": {
		"resourceInfo": {
			"id": "/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Compute/virtualMachines/{vm-name}/providers/Microsoft.ResourceHealth/availabilityStatuses/{event-id}",
			"name": "{event-id}",
			"type": "Microsoft.ResourceHealth/availabilityStatuses",
			"properties": {
			"targetResourceId": "/subscriptions/{subscription-id}/resourceGroups/{rg-name}/providers/Microsoft.Compute/virtualMachines/{vm-name}",
			"targetResourceType": "Microsoft.Compute/virtualMachines",
			"occurredTime": "2023-07-24T19:20:37.9245071Z",
			"previousAvailabilityState": "Unavailable",
			"availabilityState": "Available"
			}
		},
		"operationalInfo": {
			"resourceEventTime": "2023-07-24T19:20:37.9245071Z"
		},
		"apiVersion": "2023-12-01"
		},
		"type": "Microsoft.ResourceNotifications.HealthResources.AvailabilityStatusChanged",
		"specversion": "1.0",
		"time": "2023-07-24T19:20:37.9245071Z"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	availabilityStatusChangedEventData := deserializeSystemEvent[azsystemevents.ResourceNotificationsHealthResourcesAvailabilityStatusChangedEventData](t, event.Data)

	require.Equal(t, "{event-id}", *availabilityStatusChangedEventData.ResourceDetails.Name)
	require.Equal(t,
		"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Compute/virtualMachines/{vm-name}/providers/Microsoft.ResourceHealth/availabilityStatuses/{event-id}",
		*availabilityStatusChangedEventData.ResourceDetails.ID)
}

func TestConsumeCloudEventResourceDeletedEvent(t *testing.T) {
	requestContent := `{
		"id": "d4611260-d179-4f86-b196-3a9d4128be2d",
		"source": "/subscriptions/{subscription-id}",
		"subject": "/subscriptions/{subscription-id}/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
		"data": {
		"resourceInfo": {
			"id": "/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
			"name": "storageAccount-name",
			"type": "Microsoft.Storage/storageAccounts"
		},
		"operationalInfo": {
			"resourceEventTime": "2023-07-28T20:11:36.6347858Z"
		}
		},
		"type": "Microsoft.ResourceNotifications.Resources.Deleted",
		"specversion": "1.0",
		"time": "2023-07-28T20:11:36.6347858Z"
	}`
	event := parseCloudEvent(t, requestContent)

	require.NotEmpty(t, event)
	sysEvent := deserializeSystemEvent[azsystemevents.ResourceNotificationsResourceManagementDeletedEventData](t, event.Data)

	require.Equal(t,
		"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
		*sysEvent.ResourceDetails.ID)
}

func TestConsumeCloudEventResourceCreatedOrUpdatedEvent(t *testing.T) {
	requestContent := `{
		"id": "4eef929a-a65c-47dd-93e2-46b8c17c6c17",
		"source": "/subscriptions/{subscription-id}",
		"subject": "/subscriptions/{subscription-id}/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
		"data": {
		"resourceInfo": {
			"tags": {},
			"id": "/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
			"name": "StorageAccount-name",
			"type": "Microsoft.Storage/storageAccounts",
			"location": "eastus",
			"properties": {
			"privateEndpointConnections": [],
			"minimumTlsVersion": "TLS1_2",
			"allowBlobPublicAccess": 1,
			"allowSharedKeyAccess": 1,
			"networkAcls": {
				"bypass": "AzureServices",
				"virtualNetworkRules": [],
				"ipRules": [],
				"defaultAction": "Allow"
			},
			"supportsHttpsTrafficOnly": 1,
			"encryption": {
				"requireInfrastructureEncryption": 0,
				"services": {
				"file": {
					"keyType": "Account",
					"enabled": 1,
					"lastEnabledTime": "2023-07-28T20:12:50.6380308Z"
				},
				"blob": {
					"keyType": "Account",
					"enabled": 1,
					"lastEnabledTime": "2023-07-28T20:12:50.6380308Z"
				}
				},
				"keySource": "Microsoft.Storage"
			},
			"accessTier": "Hot",
			"provisioningState": "Succeeded",
			"creationTime": "2023-07-28T20:12:50.4661564Z",
			"primaryEndpoints": {
				"dfs": "https://{storageAccount-name}.dfs.core.windows.net/",
				"web": "https://{storageAccount-name}.z13.web.core.windows.net/",
				"blob": "https://{storageAccount-name}.blob.core.windows.net/",
				"queue": "https://{storageAccount-name}.queue.core.windows.net/",
				"table": "https://{storageAccount-name}.table.core.windows.net/",
				"file": "https://{storageAccount-name}.file.core.windows.net/"
			},
			"primaryLocation": "eastus",
			"statusOfPrimary": "available",
			"secondaryLocation": "westus",
			"statusOfSecondary": "available",
			"secondaryEndpoints": {
				"dfs": "https://{storageAccount-name} -secondary.dfs.core.windows.net/",
				"web": "https://{storageAccount-name}-secondary.z13.web.core.windows.net/",
				"blob": "https://{storageAccount-name}-secondary.blob.core.windows.net/",
				"queue": "https://{storageAccount-name}-secondary.queue.core.windows.net/",
				"table": "https://{storageAccount-name}-secondary.table.core.windows.net/"
			}
			}
		},
		"operationalInfo": {
			"resourceEventTime": "2023-07-28T20:13:10.8418063Z"
		},
		"apiVersion": "2019-06-01"
		},
		"type": "Microsoft.ResourceNotifications.Resources.CreatedOrUpdated",
		"specversion": "1.0",
		"time": "2023-07-28T20:13:10.8418063Z"
	}`
	event := parseCloudEvent(t, requestContent)

	sysEvent := deserializeSystemEvent[azsystemevents.ResourceNotificationsResourceManagementCreatedOrUpdatedEventData](t, event.Data)

	require.Equal(t,
		"/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/{rg-name}/providers/Microsoft.Storage/storageAccounts/{storageAccount-name}",
		*sysEvent.ResourceDetails.ID)
}
