// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstreamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics/v2"
	"github.com/stretchr/testify/suite"
)

type OutputsTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	jobName            string
	outputName         string
	storageAccountId   string
	storageAccountKey  string
	storageAccountName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *OutputsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.outputName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "outputna", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "streamsc", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *OutputsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestOutputsTestSuite(t *testing.T) {
	suite.Run(t, new(OutputsTestSuite))
}

func (testsuite *OutputsTestSuite) Prepare() {
	var err error
	// From step StreamingJobs_CreateOrReplace
	fmt.Println("Call operation: StreamingJobs_CreateOrReplace")
	streamingJobsClient, err := armstreamanalytics.NewStreamingJobsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	streamingJobsClientCreateOrReplaceResponsePoller, err := streamingJobsClient.BeginCreateOrReplace(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, armstreamanalytics.StreamingJob{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1":      to.Ptr("value1"),
			"key3":      to.Ptr("value3"),
			"randomKey": to.Ptr("randomValue"),
		},
		Properties: &armstreamanalytics.StreamingJobProperties{
			CompatibilityLevel:                 to.Ptr(armstreamanalytics.CompatibilityLevelOne0),
			DataLocale:                         to.Ptr("en-US"),
			EventsLateArrivalMaxDelayInSeconds: to.Ptr[int32](16),
			EventsOutOfOrderMaxDelayInSeconds:  to.Ptr[int32](5),
			EventsOutOfOrderPolicy:             to.Ptr(armstreamanalytics.EventsOutOfOrderPolicyDrop),
			Functions:                          []*armstreamanalytics.Function{},
			Inputs:                             []*armstreamanalytics.Input{},
			OutputErrorPolicy:                  to.Ptr(armstreamanalytics.OutputErrorPolicyDrop),
			Outputs:                            []*armstreamanalytics.Output{},
			SKU: &armstreamanalytics.SKU{
				Name: to.Ptr(armstreamanalytics.SKUNameStandard),
			},
		},
	}, &armstreamanalytics.StreamingJobsClientBeginCreateOrReplaceOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingJobsClientCreateOrReplaceResponsePoller)
	testsuite.Require().NoError(err)

	// From step Create_StorageAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"storageAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
			"storageAccountKey": map[string]any{
				"type":  "string",
				"value": "[listKeys(parameters('storageAccountName'),'2022-05-01').keys[0].value]",
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
					"name": "Standard_RAGRS",
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
	testsuite.storageAccountKey = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountKey"].(map[string]interface{})["value"].(string)
}

// Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}
func (testsuite *OutputsTestSuite) TestOutputs() {
	var err error
	// From step Outputs_CreateOrReplace
	fmt.Println("Call operation: Outputs_CreateOrReplace")
	outputsClient, err := armstreamanalytics.NewOutputsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = outputsClient.CreateOrReplace(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.outputName, armstreamanalytics.Output{
		Properties: &armstreamanalytics.OutputProperties{
			Datasource: &armstreamanalytics.BlobOutputDataSource{
				Type: to.Ptr("Microsoft.Storage/Blob"),
				Properties: &armstreamanalytics.BlobOutputDataSourceProperties{
					Container:   to.Ptr("state"),
					DateFormat:  to.Ptr("yyyy/MM/dd"),
					PathPattern: to.Ptr("{date}/{time}"),
					StorageAccounts: []*armstreamanalytics.StorageAccount{
						{
							AccountKey:  to.Ptr(testsuite.storageAccountKey),
							AccountName: to.Ptr(testsuite.storageAccountName),
						}},
					TimeFormat: to.Ptr("HH"),
				},
			},
			Serialization: &armstreamanalytics.CSVSerialization{
				Type: to.Ptr(armstreamanalytics.EventSerializationTypeCSV),
				Properties: &armstreamanalytics.CSVSerializationProperties{
					Encoding:       to.Ptr(armstreamanalytics.EncodingUTF8),
					FieldDelimiter: to.Ptr(","),
				},
			},
		},
	}, &armstreamanalytics.OutputsClientCreateOrReplaceOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step Outputs_ListByStreamingJob
	fmt.Println("Call operation: Outputs_ListByStreamingJob")
	outputsClientNewListByStreamingJobPager := outputsClient.NewListByStreamingJobPager(testsuite.resourceGroupName, testsuite.jobName, &armstreamanalytics.OutputsClientListByStreamingJobOptions{Select: nil})
	for outputsClientNewListByStreamingJobPager.More() {
		_, err := outputsClientNewListByStreamingJobPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Outputs_Get
	fmt.Println("Call operation: Outputs_Get")
	_, err = outputsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.outputName, nil)
	testsuite.Require().NoError(err)

	// From step Outputs_Update
	fmt.Println("Call operation: Outputs_Update")
	_, err = outputsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.outputName, armstreamanalytics.Output{
		Properties: &armstreamanalytics.OutputProperties{
			Datasource: &armstreamanalytics.BlobOutputDataSource{
				Type: to.Ptr("Microsoft.Storage/Blob"),
				Properties: &armstreamanalytics.BlobOutputDataSourceProperties{
					Container: to.Ptr("differentContainer"),
				},
			},
			Serialization: &armstreamanalytics.CSVSerialization{
				Type: to.Ptr(armstreamanalytics.EventSerializationTypeCSV),
				Properties: &armstreamanalytics.CSVSerializationProperties{
					Encoding:       to.Ptr(armstreamanalytics.EncodingUTF8),
					FieldDelimiter: to.Ptr("|"),
				},
			},
		},
	}, &armstreamanalytics.OutputsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Outputs_Test
	fmt.Println("Call operation: Outputs_Test")
	outputsClientTestResponsePoller, err := outputsClient.BeginTest(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.outputName, &armstreamanalytics.OutputsClientBeginTestOptions{Output: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, outputsClientTestResponsePoller)
	testsuite.Require().NoError(err)

	// From step Outputs_Delete
	fmt.Println("Call operation: Outputs_Delete")
	_, err = outputsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.outputName, nil)
	testsuite.Require().NoError(err)
}
