//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmediaservices_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type StreamingserviceTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	accountName           string
	assetName             string
	liveEventName         string
	liveOutputName        string
	storageAccountId      string
	storageAccountName    string
	streamingEndpointName string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *StreamingserviceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mediaservices/armmediaservices/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.assetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "assetnam", 14, false)
	testsuite.liveEventName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "liveeven", 14, false)
	testsuite.liveOutputName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "liveoutp", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storageaccount", 20, true)
	testsuite.streamingEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "streamin", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *StreamingserviceTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestStreamingserviceTestSuite(t *testing.T) {
	suite.Run(t, new(StreamingserviceTestSuite))
}

func (testsuite *StreamingserviceTestSuite) Prepare() {
	var err error
	// From step Create_StorageAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"storageAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"storageAccountName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.storageAccountName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]any{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]any{
							"blob": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]any{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []any{},
						"virtualNetworkRules": []any{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]any{
					"name": "Standard_LRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_StorageAccount", &deployment)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step Mediaservices_CreateOrUpdate
	fmt.Println("Call operation: Mediaservices_CreateOrUpdate")
	client, err := armmediaservices.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmediaservices.MediaService{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
		Properties: &armmediaservices.MediaServiceProperties{
			StorageAccounts: []*armmediaservices.StorageAccount{
				{
					Type: to.Ptr(armmediaservices.StorageAccountTypePrimary),
					ID:   to.Ptr(testsuite.storageAccountId),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *StreamingserviceTestSuite) TestGeneratedScenario() {
	var err error
	// From step LiveEvents_Create
	fmt.Println("Call operation: LiveEvents_Create")
	liveEventsClient, err := armmediaservices.NewLiveEventsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	liveEventsClientCreateResponsePoller, err := liveEventsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, armmediaservices.LiveEvent{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armmediaservices.LiveEventProperties{
			Description: to.Ptr("test event 1"),
			Input: &armmediaservices.LiveEventInput{
				AccessControl: &armmediaservices.LiveEventInputAccessControl{
					IP: &armmediaservices.IPAccessControl{
						Allow: []*armmediaservices.IPRange{
							{
								Name:               to.Ptr("AllowAll"),
								Address:            to.Ptr("0.0.0.0"),
								SubnetPrefixLength: to.Ptr[int32](0),
							}},
					},
				},
				KeyFrameIntervalDuration: to.Ptr("PT6S"),
				StreamingProtocol:        to.Ptr(armmediaservices.LiveEventInputProtocolRTMP),
			},
			Preview: &armmediaservices.LiveEventPreview{
				AccessControl: &armmediaservices.LiveEventPreviewAccessControl{
					IP: &armmediaservices.IPAccessControl{
						Allow: []*armmediaservices.IPRange{
							{
								Name:               to.Ptr("AllowAll"),
								Address:            to.Ptr("0.0.0.0"),
								SubnetPrefixLength: to.Ptr[int32](0),
							}},
					},
				},
			},
		},
	}, &armmediaservices.LiveEventsClientBeginCreateOptions{AutoStart: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveEventsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LiveEvents_List
	fmt.Println("Call operation: LiveEvents_List")
	liveEventsClientNewListPager := liveEventsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for liveEventsClientNewListPager.More() {
		_, err := liveEventsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LiveEvents_Get
	fmt.Println("Call operation: LiveEvents_Get")
	_, err = liveEventsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, nil)
	testsuite.Require().NoError(err)

	// From step LiveEvents_Allocate
	fmt.Println("Call operation: LiveEvents_Allocate")
	liveEventsClientAllocateResponsePoller, err := liveEventsClient.BeginAllocate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveEventsClientAllocateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LiveEvents_Stop
	fmt.Println("Call operation: LiveEvents_Stop")
	liveEventsClientStopResponsePoller, err := liveEventsClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, armmediaservices.LiveEventActionInput{
		RemoveOutputsOnStop: to.Ptr(false),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveEventsClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_Create
	fmt.Println("Call operation: StreamingEndpoints_Create")
	streamingEndpointsClient, err := armmediaservices.NewStreamingEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	streamingEndpointsClientCreateResponsePoller, err := streamingEndpointsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, armmediaservices.StreamingEndpoint{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armmediaservices.StreamingEndpointProperties{
			Description: to.Ptr("test event 1"),
			AccessControl: &armmediaservices.StreamingEndpointAccessControl{
				Akamai: &armmediaservices.AkamaiAccessControl{
					AkamaiSignatureHeaderAuthenticationKeyList: []*armmediaservices.AkamaiSignatureHeaderAuthenticationKey{
						{
							Base64Key:  to.Ptr("dGVzdGlkMQ=="),
							Expiration: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2029-12-31T16:00:00-08:00"); return t }()),
							Identifier: to.Ptr("id1"),
						},
						{
							Base64Key:  to.Ptr("dGVzdGlkMQ=="),
							Expiration: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2030-12-31T16:00:00-08:00"); return t }()),
							Identifier: to.Ptr("id2"),
						}},
				},
				IP: &armmediaservices.IPAccessControl{
					Allow: []*armmediaservices.IPRange{
						{
							Name:    to.Ptr("AllowedIp"),
							Address: to.Ptr("192.168.1.1"),
						}},
				},
			},
			AvailabilitySetName: to.Ptr("availableset"),
			CdnEnabled:          to.Ptr(false),
			ScaleUnits:          to.Ptr[int32](1),
		},
	}, &armmediaservices.StreamingEndpointsClientBeginCreateOptions{AutoStart: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingEndpointsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_Get
	fmt.Println("Call operation: StreamingEndpoints_Get")
	_, err = streamingEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_List
	fmt.Println("Call operation: StreamingEndpoints_List")
	streamingEndpointsClientNewListPager := streamingEndpointsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for streamingEndpointsClientNewListPager.More() {
		_, err := streamingEndpointsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StreamingEndpoints_Skus
	fmt.Println("Call operation: StreamingEndpoints_SKUs")
	_, err = streamingEndpointsClient.SKUs(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_Scale
	fmt.Println("Call operation: StreamingEndpoints_Scale")
	streamingEndpointsClientScaleResponsePoller, err := streamingEndpointsClient.BeginScale(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, armmediaservices.StreamingEntityScaleUnit{
		ScaleUnit: to.Ptr[int32](5),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingEndpointsClientScaleResponsePoller)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_Stop
	fmt.Println("Call operation: StreamingEndpoints_Stop")
	streamingEndpointsClientStopResponsePoller, err := streamingEndpointsClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingEndpointsClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Assets_CreateOrUpdate
	fmt.Println("Call operation: Assets_CreateOrUpdate")
	assetsClient, err := armmediaservices.NewAssetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = assetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, armmediaservices.Asset{
		Properties: &armmediaservices.AssetProperties{
			Description:        to.Ptr("A documentary showing the ascent of Mount Logan"),
			StorageAccountName: to.Ptr(testsuite.storageAccountName),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step LiveOutputs_Create
	fmt.Println("Call operation: LiveOutputs_Create")
	liveOutputsClient, err := armmediaservices.NewLiveOutputsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	liveOutputsClientCreateResponsePoller, err := liveOutputsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, testsuite.liveOutputName, armmediaservices.LiveOutput{
		Properties: &armmediaservices.LiveOutputProperties{
			Description:         to.Ptr("test live output 1"),
			ArchiveWindowLength: to.Ptr("PT5M"),
			AssetName:           to.Ptr(testsuite.assetName),
			Hls: &armmediaservices.Hls{
				FragmentsPerTsSegment: to.Ptr[int32](5),
			},
			ManifestName:       to.Ptr("testmanifest"),
			RewindWindowLength: to.Ptr("PT4M"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveOutputsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LiveOutputs_List
	fmt.Println("Call operation: LiveOutputs_List")
	liveOutputsClientNewListPager := liveOutputsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, nil)
	for liveOutputsClientNewListPager.More() {
		_, err := liveOutputsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LiveOutputs_Get
	fmt.Println("Call operation: LiveOutputs_Get")
	_, err = liveOutputsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, testsuite.liveOutputName, nil)
	testsuite.Require().NoError(err)

	// From step LiveOutputs_Delete
	fmt.Println("Call operation: LiveOutputs_Delete")
	liveOutputsClientDeleteResponsePoller, err := liveOutputsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, testsuite.liveOutputName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveOutputsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step StreamingEndpoints_Delete
	fmt.Println("Call operation: StreamingEndpoints_Delete")
	streamingEndpointsClientDeleteResponsePoller, err := streamingEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingEndpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step LiveEvents_Delete
	fmt.Println("Call operation: LiveEvents_Delete")
	liveEventsClientDeleteResponsePoller, err := liveEventsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.liveEventName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, liveEventsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
