//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmediaservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type EncodingTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	accountName        string
	assetName          string
	jobName            string
	storageAccountId   string
	storageAccountName string
	transformName      string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *EncodingTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mediaservices/armmediaservices/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.assetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "assetnam", 14, false)
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storageaccount", 20, true)
	testsuite.transformName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "transfor", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *EncodingTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestEncodingTestSuite(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	suite.Run(t, new(EncodingTestSuite))
}

func (testsuite *EncodingTestSuite) Prepare() {
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

func (testsuite *EncodingTestSuite) TestEncodingTests() {
	var err error
	// From step Transforms_CreateOrUpdate
	fmt.Println("Call operation: Transforms_CreateOrUpdate")
	transformsClient, err := armmediaservices.NewTransformsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = transformsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, armmediaservices.Transform{
		Properties: &armmediaservices.TransformProperties{
			Description: to.Ptr("Example Transform to illustrate create and update."),
			Outputs: []*armmediaservices.TransformOutput{
				{
					Preset: &armmediaservices.BuiltInStandardEncoderPreset{
						ODataType:  to.Ptr("#Microsoft.Media.BuiltInStandardEncoderPreset"),
						PresetName: to.Ptr(armmediaservices.EncoderNamedPresetAdaptiveStreaming),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Transforms_List
	fmt.Println("Call operation: Transforms_List")
	transformsClientNewListPager := transformsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armmediaservices.TransformsClientListOptions{Filter: nil,
		Orderby: nil,
	})
	for transformsClientNewListPager.More() {
		_, err := transformsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Transforms_Get
	fmt.Println("Call operation: Transforms_Get")
	_, err = transformsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, nil)
	testsuite.Require().NoError(err)

	// From step Transforms_Update
	fmt.Println("Call operation: Transforms_Update")
	_, err = transformsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, armmediaservices.Transform{
		Properties: &armmediaservices.TransformProperties{
			Description: to.Ptr("Example transform to illustrate update."),
			Outputs: []*armmediaservices.TransformOutput{
				{
					Preset: &armmediaservices.BuiltInStandardEncoderPreset{
						ODataType:  to.Ptr("#Microsoft.Media.BuiltInStandardEncoderPreset"),
						PresetName: to.Ptr(armmediaservices.EncoderNamedPresetH264MultipleBitrate720P),
					},
					RelativePriority: to.Ptr(armmediaservices.PriorityHigh),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Assets_CreateOrUpdate_InputAsset
	fmt.Println("Call operation: Assets_CreateOrUpdate")
	assetsClient, err := armmediaservices.NewAssetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = assetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.jobName+"-InputAsset", armmediaservices.Asset{
		Properties: &armmediaservices.AssetProperties{
			Description:        to.Ptr("A documentary showing the ascent of Mount Logan"),
			StorageAccountName: to.Ptr(testsuite.storageAccountName),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Assets_CreateOrUpdate_OutputAsset
	fmt.Println("Call operation: Assets_CreateOrUpdate")
	_, err = assetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.jobName+"-OutputAsset", armmediaservices.Asset{
		Properties: &armmediaservices.AssetProperties{
			Description:        to.Ptr("A documentary showing the ascent of Mount Logan"),
			StorageAccountName: to.Ptr(testsuite.storageAccountName),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Jobs_Create
	fmt.Println("Call operation: Jobs_Create")
	jobsClient, err := armmediaservices.NewJobsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = jobsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, testsuite.jobName, armmediaservices.Job{
		Properties: &armmediaservices.JobProperties{
			CorrelationData: map[string]*string{
				"Key 2": to.Ptr("Value 2"),
				"key1":  to.Ptr("value1"),
			},
			Input: &armmediaservices.JobInputAsset{
				ODataType: to.Ptr("#Microsoft.Media.JobInputAsset"),
				AssetName: to.Ptr(testsuite.jobName + "-InputAsset"),
			},
			Outputs: []armmediaservices.JobOutputClassification{
				&armmediaservices.JobOutputAsset{
					ODataType: to.Ptr("#Microsoft.Media.JobOutputAsset"),
					AssetName: to.Ptr(testsuite.jobName + "-OutputAsset"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Jobs_List
	fmt.Println("Call operation: Jobs_List")
	jobsClientNewListPager := jobsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, &armmediaservices.JobsClientListOptions{Filter: nil,
		Orderby: nil,
	})
	for jobsClientNewListPager.More() {
		_, err := jobsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Jobs_Get
	fmt.Println("Call operation: Jobs_Get")
	_, err = jobsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, testsuite.jobName, nil)
	testsuite.Require().NoError(err)

	// From step Jobs_Update
	fmt.Println("Call operation: Jobs_Update")
	_, err = jobsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, testsuite.jobName, armmediaservices.Job{
		Properties: &armmediaservices.JobProperties{
			Description: to.Ptr("Example job to illustrate update."),
			Input: &armmediaservices.JobInputAsset{
				ODataType: to.Ptr("#Microsoft.Media.JobInputAsset"),
				AssetName: to.Ptr(testsuite.jobName + "-InputAsset"),
			},
			Outputs: []armmediaservices.JobOutputClassification{
				&armmediaservices.JobOutputAsset{
					ODataType: to.Ptr("#Microsoft.Media.JobOutputAsset"),
					AssetName: to.Ptr(testsuite.jobName + "-OutputAsset"),
				}},
			Priority: to.Ptr(armmediaservices.PriorityHigh),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Jobs_CancelJob
	fmt.Println("Call operation: Jobs_CancelJob")
	_, err = jobsClient.CancelJob(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, testsuite.jobName, nil)
	testsuite.Require().NoError(err)

	// From step Transforms_Delete
	fmt.Println("Call operation: Transforms_Delete")
	_, err = transformsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.transformName, nil)
	testsuite.Require().NoError(err)
}
