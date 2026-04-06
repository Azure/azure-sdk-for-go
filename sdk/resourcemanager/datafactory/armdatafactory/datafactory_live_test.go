// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdatafactory_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v10"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type DatafactoryTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	credentialName             string
	dataFlowName               string
	datasetName                string
	factoryName                string
	integrationRuntimeName     string
	linkedServiceName          string
	locationId                 string
	managedPrivateEndpointName string
	managedVirtualNetworkName  string
	nodeName                   string
	pipelineName               string
	runId                      string
	triggerName                string
	location                   string
	resourceGroupName          string
	subscriptionId             string

	datafactoryId string
}

func (testsuite *DatafactoryTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.credentialName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "credenti", 8+6, false)
	testsuite.dataFlowName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "dataflow", 8+6, false)
	testsuite.datasetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "datasetn", 8+6, false)
	testsuite.factoryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "factoryn", 8+6, false)
	testsuite.integrationRuntimeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "integrat", 8+6, false)
	testsuite.linkedServiceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "linkedse", 8+6, false)
	testsuite.locationId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "location", 8+6, false)
	testsuite.managedPrivateEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "managedp", 8+6, false)
	testsuite.managedVirtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "managedv", 8+6, false)
	testsuite.nodeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "nodename", 8+6, false)
	testsuite.pipelineName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "pipeline", 8+6, false)
	testsuite.runId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "runid", 8+6, false)
	testsuite.triggerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "triggern", 8+6, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DatafactoryTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDatafactoryTestSuite(t *testing.T) {
	suite.Run(t, new(DatafactoryTestSuite))
}

func (testsuite *DatafactoryTestSuite) Prepare() {
	var err error
	// From step Factories_CreateOrUpdate
	fmt.Println("Call operation: Factories_CreateOrUpdate")
	factoriesClient, err := armdatafactory.NewFactoriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	factoriesClientCreateOrUpdateResponse, err := factoriesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.Factory{
		Location: to.Ptr(testsuite.location),
	}, &armdatafactory.FactoriesClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	testsuite.datafactoryId = *factoriesClientCreateOrUpdateResponse.ID

	// From step LinkedServices_CreateOrUpdate
	fmt.Println("Call operation: LinkedServices_CreateOrUpdate")
	linkedServicesClient, err := armdatafactory.NewLinkedServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = linkedServicesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.linkedServiceName, armdatafactory.LinkedServiceResource{
		Properties: &armdatafactory.AzureStorageLinkedService{
			Type: to.Ptr("AzureStorage"),
			TypeProperties: &armdatafactory.AzureStorageLinkedServiceTypeProperties{
				ConnectionString: map[string]any{
					"type":  "SecureString",
					"value": "DefaultEndpointsProtocol=https;AccountName=examplestorageaccount;AccountKey=<storage key>",
				},
			},
		},
	}, &armdatafactory.LinkedServicesClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Datasets_CreateOrUpdate
	fmt.Println("Call operation: Datasets_CreateOrUpdate")
	datasetsClient, err := armdatafactory.NewDatasetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = datasetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.datasetName, armdatafactory.DatasetResource{
		Properties: &armdatafactory.AzureBlobDataset{
			Type: to.Ptr("AzureBlob"),
			LinkedServiceName: &armdatafactory.LinkedServiceReference{
				Type:          to.Ptr(armdatafactory.LinkedServiceReferenceTypeLinkedServiceReference),
				ReferenceName: to.Ptr(testsuite.linkedServiceName),
			},
			Parameters: map[string]*armdatafactory.ParameterSpecification{
				"MyFileName": {
					Type: to.Ptr(armdatafactory.ParameterTypeString),
				},
				"MyFolderPath": {
					Type: to.Ptr(armdatafactory.ParameterTypeString),
				},
			},
			TypeProperties: &armdatafactory.AzureBlobDatasetTypeProperties{
				Format: &armdatafactory.TextFormat{
					Type: to.Ptr("TextFormat"),
				},
				FileName: map[string]any{
					"type":  "Expression",
					"value": "@dataset().MyFileName",
				},
				FolderPath: map[string]any{
					"type":  "Expression",
					"value": "@dataset().MyFolderPath",
				},
			},
		},
	}, &armdatafactory.DatasetsClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Pipelines_CreateOrUpdate
	fmt.Println("Call operation: Pipelines_CreateOrUpdate")
	pipelinesClient, err := armdatafactory.NewPipelinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = pipelinesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.pipelineName, armdatafactory.PipelineResource{
		Properties: &armdatafactory.Pipeline{
			Activities: []armdatafactory.ActivityClassification{
				&armdatafactory.ForEachActivity{
					Name: to.Ptr("ExampleForeachActivity"),
					Type: to.Ptr("ForEach"),
					TypeProperties: &armdatafactory.ForEachActivityTypeProperties{
						Activities: []armdatafactory.ActivityClassification{
							&armdatafactory.CopyActivity{
								Name: to.Ptr("ExampleCopyActivity"),
								Type: to.Ptr("Copy"),
								Inputs: []*armdatafactory.DatasetReference{
									{
										Type: to.Ptr(armdatafactory.DatasetReferenceTypeDatasetReference),
										Parameters: map[string]any{
											"MyFileName":   "examplecontainer.csv",
											"MyFolderPath": "examplecontainer",
										},
										ReferenceName: to.Ptr(testsuite.datasetName),
									}},
								Outputs: []*armdatafactory.DatasetReference{
									{
										Type: to.Ptr(armdatafactory.DatasetReferenceTypeDatasetReference),
										Parameters: map[string]any{
											"MyFileName": map[string]any{
												"type":  "Expression",
												"value": "@item()",
											},
											"MyFolderPath": "examplecontainer",
										},
										ReferenceName: to.Ptr(testsuite.datasetName),
									}},
								TypeProperties: &armdatafactory.CopyActivityTypeProperties{
									DataIntegrationUnits: float64(32),
									Sink: &armdatafactory.BlobSink{
										Type: to.Ptr("BlobSink"),
									},
									Source: &armdatafactory.BlobSource{
										Type: to.Ptr("BlobSource"),
									},
								},
							}},
						IsSequential: to.Ptr(true),
						Items: &armdatafactory.Expression{
							Type:  to.Ptr(armdatafactory.ExpressionTypeExpression),
							Value: to.Ptr("@pipeline().parameters.OutputBlobNameList"),
						},
					},
				}},
			Parameters: map[string]*armdatafactory.ParameterSpecification{
				"JobId": {
					Type: to.Ptr(armdatafactory.ParameterTypeString),
				},
				"OutputBlobNameList": {
					Type: to.Ptr(armdatafactory.ParameterTypeArray),
				},
			},
			Policy: &armdatafactory.PipelinePolicy{
				ElapsedTimeMetric: &armdatafactory.PipelineElapsedTimeMetricPolicy{
					Duration: "0.00:10:00",
				},
			},
			RunDimensions: map[string]any{
				"JobId": map[string]any{
					"type":  "Expression",
					"value": "@pipeline().parameters.JobId",
				},
			},
			Variables: map[string]*armdatafactory.VariableSpecification{
				"TestVariableArray": {
					Type: to.Ptr(armdatafactory.VariableTypeArray),
				},
			},
		},
	}, &armdatafactory.PipelinesClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/operations
func (testsuite *DatafactoryTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armdatafactory.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DataFactory/factories/{factoryName}
func (testsuite *DatafactoryTestSuite) TestFactories() {
	var err error
	// From step Factories_List
	fmt.Println("Call operation: Factories_List")
	factoriesClient, err := armdatafactory.NewFactoriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	factoriesClientNewListPager := factoriesClient.NewListPager(nil)
	for factoriesClientNewListPager.More() {
		_, err := factoriesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Factories_ListByResourceGroup
	fmt.Println("Call operation: Factories_ListByResourceGroup")
	factoriesClientNewListByResourceGroupPager := factoriesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for factoriesClientNewListByResourceGroupPager.More() {
		_, err := factoriesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Factories_Get
	fmt.Println("Call operation: Factories_Get")
	_, err = factoriesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, &armdatafactory.FactoriesClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step Factories_Update
	fmt.Println("Call operation: Factories_Update")
	_, err = factoriesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.FactoryUpdateParameters{
		Tags: map[string]*string{
			"exampleTag": to.Ptr("exampleValue"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Factories_GetDataPlaneAccess
	fmt.Println("Call operation: Factories_GetDataPlaneAccess")
	_, err = factoriesClient.GetDataPlaneAccess(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.UserAccessPolicy{
		AccessResourcePath: to.Ptr(""),
		Permissions:        to.Ptr("r"),
		ProfileName:        to.Ptr("DefaultProfile"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/locations/{locationId}/getFeatureValue
func (testsuite *DatafactoryTestSuite) TestExposureControl() {
	var err error
	// From step ExposureControl_GetFeatureValueByFactory
	fmt.Println("Call operation: ExposureControl_GetFeatureValueByFactory")
	exposureControlClient, err := armdatafactory.NewExposureControlClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = exposureControlClient.GetFeatureValueByFactory(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.ExposureControlRequest{
		FeatureName: to.Ptr("ADFIntegrationRuntimeSharingRbac"),
		FeatureType: to.Ptr("Feature"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step ExposureControl_QueryFeatureValuesByFactory
	fmt.Println("Call operation: ExposureControl_QueryFeatureValuesByFactory")
	_, err = exposureControlClient.QueryFeatureValuesByFactory(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.ExposureControlBatchRequest{
		ExposureControlRequests: []*armdatafactory.ExposureControlRequest{
			{
				FeatureName: to.Ptr("ADFIntegrationRuntimeSharingRbac"),
				FeatureType: to.Ptr("Feature"),
			},
			{
				FeatureName: to.Ptr("ADFSampleFeature"),
				FeatureType: to.Ptr("Feature"),
			}},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/integrationRuntimes/{integrationRuntimeName}
func (testsuite *DatafactoryTestSuite) TestIntegrationRuntimes() {
	var err error
	// From step IntegrationRuntimes_CreateOrUpdate
	fmt.Println("Call operation: IntegrationRuntimes_CreateOrUpdate")
	integrationRuntimesClient, err := armdatafactory.NewIntegrationRuntimesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationRuntimesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, armdatafactory.IntegrationRuntimeResource{
		Properties: &armdatafactory.SelfHostedIntegrationRuntime{
			Type:        to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
			Description: to.Ptr("A selfhosted integration runtime"),
		},
	}, &armdatafactory.IntegrationRuntimesClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_ListByFactory
	fmt.Println("Call operation: IntegrationRuntimes_ListByFactory")
	integrationRuntimesClientNewListByFactoryPager := integrationRuntimesClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for integrationRuntimesClientNewListByFactoryPager.More() {
		_, err := integrationRuntimesClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationRuntimes_Get
	fmt.Println("Call operation: IntegrationRuntimes_Get")
	_, err = integrationRuntimesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, &armdatafactory.IntegrationRuntimesClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_Update
	fmt.Println("Call operation: IntegrationRuntimes_Update")
	_, err = integrationRuntimesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, armdatafactory.UpdateIntegrationRuntimeRequest{
		AutoUpdate:        to.Ptr(armdatafactory.IntegrationRuntimeAutoUpdateOff),
		UpdateDelayOffset: to.Ptr("\"PT3H\""),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_GetStatus
	fmt.Println("Call operation: IntegrationRuntimes_GetStatus")
	_, err = integrationRuntimesClient.GetStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_Upgrade
	fmt.Println("Call operation: IntegrationRuntimes_Upgrade")
	_, err = integrationRuntimesClient.Upgrade(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_RemoveLinks
	fmt.Println("Call operation: IntegrationRuntimes_RemoveLinks")
	_, err = integrationRuntimesClient.RemoveLinks(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, armdatafactory.LinkedIntegrationRuntimeRequest{
		LinkedFactoryName: to.Ptr(testsuite.factoryName + "-linked"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_SyncCredentials
	fmt.Println("Call operation: IntegrationRuntimes_SyncCredentials")
	_, err = integrationRuntimesClient.SyncCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_ListAuthKeys
	fmt.Println("Call operation: IntegrationRuntimes_ListAuthKeys")
	_, err = integrationRuntimesClient.ListAuthKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_RegenerateAuthKey
	fmt.Println("Call operation: IntegrationRuntimes_RegenerateAuthKey")
	_, err = integrationRuntimesClient.RegenerateAuthKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, armdatafactory.IntegrationRuntimeRegenerateKeyParameters{
		KeyName: to.Ptr(armdatafactory.IntegrationRuntimeAuthKeyNameAuthKey2),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_GetMonitoringData
	fmt.Println("Call operation: IntegrationRuntimes_GetMonitoringData")
	_, err = integrationRuntimesClient.GetMonitoringData(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationRuntimes_Delete
	fmt.Println("Call operation: IntegrationRuntimes_Delete")
	_, err = integrationRuntimesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.integrationRuntimeName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/linkedservices/{linkedServiceName}
func (testsuite *DatafactoryTestSuite) TestLinkedServices() {
	var err error
	// From step LinkedServices_ListByFactory
	fmt.Println("Call operation: LinkedServices_ListByFactory")
	linkedServicesClient, err := armdatafactory.NewLinkedServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	linkedServicesClientNewListByFactoryPager := linkedServicesClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for linkedServicesClientNewListByFactoryPager.More() {
		_, err := linkedServicesClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LinkedServices_Get
	fmt.Println("Call operation: LinkedServices_Get")
	_, err = linkedServicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.linkedServiceName, &armdatafactory.LinkedServicesClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/datasets/{datasetName}
func (testsuite *DatafactoryTestSuite) TestDatasets() {
	var err error
	// From step Datasets_ListByFactory
	fmt.Println("Call operation: Datasets_ListByFactory")
	datasetsClient, err := armdatafactory.NewDatasetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	datasetsClientNewListByFactoryPager := datasetsClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for datasetsClientNewListByFactoryPager.More() {
		_, err := datasetsClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Datasets_Get
	fmt.Println("Call operation: Datasets_Get")
	_, err = datasetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.datasetName, &armdatafactory.DatasetsClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/pipelines/{pipelineName}
func (testsuite *DatafactoryTestSuite) TestPipelines() {
	// runId := testsuite.runId
	var err error
	// From step Pipelines_ListByFactory
	fmt.Println("Call operation: Pipelines_ListByFactory")
	pipelinesClient, err := armdatafactory.NewPipelinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pipelinesClientNewListByFactoryPager := pipelinesClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for pipelinesClientNewListByFactoryPager.More() {
		_, err := pipelinesClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Pipelines_Get
	fmt.Println("Call operation: Pipelines_Get")
	_, err = pipelinesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.pipelineName, &armdatafactory.PipelinesClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step Pipelines_CreateRun
	fmt.Println("Call operation: Pipelines_CreateRun")
	pipelinesClientCreateRunResponse, err := pipelinesClient.CreateRun(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.pipelineName, &armdatafactory.PipelinesClientCreateRunOptions{ReferencePipelineRunID: nil,
		IsRecovery:        nil,
		StartActivityName: nil,
		StartFromFailure:  nil,
		Parameters: map[string]any{
			"OutputBlobNameList": []any{
				"exampleoutput.csv",
			},
		},
	})
	testsuite.Require().NoError(err)
	testsuite.runId = *pipelinesClientCreateRunResponse.RunID

	// From step PipelineRuns_Get
	fmt.Println("Call operation: PipelineRuns_Get")
	pipelineRunsClient, err := armdatafactory.NewPipelineRunsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = pipelineRunsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.runId, nil)
	testsuite.Require().NoError(err)

	// From step ActivityRuns_QueryByPipelineRun
	fmt.Println("Call operation: ActivityRuns_QueryByPipelineRun")
	activityRunsClient, err := armdatafactory.NewActivityRunsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = activityRunsClient.QueryByPipelineRun(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.runId, armdatafactory.RunFilterParameters{
		LastUpdatedAfter:  to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:36:44.3345758Z"); return t }()),
		LastUpdatedBefore: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:49:48.3686473Z"); return t }()),
	}, nil)
	testsuite.Require().NoError(err)

	// From step PipelineRuns_QueryByFactory
	fmt.Println("Call operation: PipelineRuns_QueryByFactory")
	_, err = pipelineRunsClient.QueryByFactory(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.RunFilterParameters{
		Filters: []*armdatafactory.RunQueryFilter{
			{
				Operand:  to.Ptr(armdatafactory.RunQueryFilterOperandPipelineName),
				Operator: to.Ptr(armdatafactory.RunQueryFilterOperatorEquals),
				Values: []*string{
					to.Ptr("examplePipeline")},
			}},
		LastUpdatedAfter:  to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:36:44.3345758Z"); return t }()),
		LastUpdatedBefore: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:49:48.3686473Z"); return t }()),
	}, nil)
	testsuite.Require().NoError(err)

	// From step PipelineRuns_Cancel
	fmt.Println("Call operation: PipelineRuns_Cancel")
	_, err = pipelineRunsClient.Cancel(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.runId, &armdatafactory.PipelineRunsClientCancelOptions{IsRecursive: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/triggers/{triggerName}
func (testsuite *DatafactoryTestSuite) TestTriggers() {
	var err error
	// From step Triggers_CreateOrUpdate
	fmt.Println("Call operation: Triggers_CreateOrUpdate")
	triggersClient, err := armdatafactory.NewTriggersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = triggersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, armdatafactory.TriggerResource{
		Properties: &armdatafactory.ScheduleTrigger{
			Type: to.Ptr("ScheduleTrigger"),
			Pipelines: []*armdatafactory.TriggerPipelineReference{
				{
					Parameters: map[string]any{
						"OutputBlobNameList": []any{
							"exampleoutput.csv",
						},
					},
					PipelineReference: &armdatafactory.PipelineReference{
						Type:          to.Ptr(armdatafactory.PipelineReferenceTypePipelineReference),
						ReferenceName: to.Ptr(testsuite.pipelineName),
					},
				}},
			TypeProperties: &armdatafactory.ScheduleTriggerTypeProperties{
				Recurrence: &armdatafactory.ScheduleTriggerRecurrence{
					EndTime:   to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:55:13.8441801Z"); return t }()),
					Frequency: to.Ptr(armdatafactory.RecurrenceFrequencyMinute),
					Interval:  to.Ptr[int32](4),
					StartTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-16T00:39:13.8441801Z"); return t }()),
					TimeZone:  to.Ptr("UTC"),
				},
			},
		},
	}, &armdatafactory.TriggersClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Triggers_ListByFactory
	fmt.Println("Call operation: Triggers_ListByFactory")
	triggersClientNewListByFactoryPager := triggersClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for triggersClientNewListByFactoryPager.More() {
		_, err := triggersClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Triggers_Get
	fmt.Println("Call operation: Triggers_Get")
	_, err = triggersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, &armdatafactory.TriggersClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step Triggers_GetEventSubscriptionStatus
	fmt.Println("Call operation: Triggers_GetEventSubscriptionStatus")
	_, err = triggersClient.GetEventSubscriptionStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)

	// From step Triggers_SubscribeToEvents
	fmt.Println("Call operation: Triggers_SubscribeToEvents")
	triggersClientSubscribeToEventsResponsePoller, err := triggersClient.BeginSubscribeToEvents(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, triggersClientSubscribeToEventsResponsePoller)
	testsuite.Require().NoError(err)

	// From step Triggers_Start
	fmt.Println("Call operation: Triggers_Start")
	triggersClientStartResponsePoller, err := triggersClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, triggersClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step Triggers_QueryByFactory
	fmt.Println("Call operation: Triggers_QueryByFactory")
	_, err = triggersClient.QueryByFactory(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.TriggerFilterParameters{
		ParentTriggerName: to.Ptr(testsuite.triggerName),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Triggers_UnsubscribeFromEvents
	fmt.Println("Call operation: Triggers_UnsubscribeFromEvents")
	triggersClientUnsubscribeFromEventsResponsePoller, err := triggersClient.BeginUnsubscribeFromEvents(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, triggersClientUnsubscribeFromEventsResponsePoller)
	testsuite.Require().NoError(err)

	// From step Triggers_Stop
	fmt.Println("Call operation: Triggers_Stop")
	triggersClientStopResponsePoller, err := triggersClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, triggersClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Triggers_Delete
	fmt.Println("Call operation: Triggers_Delete")
	_, err = triggersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.triggerName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/dataflows/{dataFlowName}
func (testsuite *DatafactoryTestSuite) TestDataFlows() {
	var err error
	// From step DataFlows_CreateOrUpdate
	fmt.Println("Call operation: DataFlows_CreateOrUpdate")
	dataFlowsClient, err := armdatafactory.NewDataFlowsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dataFlowsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.dataFlowName, armdatafactory.DataFlowResource{
		Properties: &armdatafactory.MappingDataFlow{
			Type:        to.Ptr("MappingDataFlow"),
			Description: to.Ptr("Sample demo data flow to convert currencies showing usage of union, derive and conditional split transformation."),
			TypeProperties: &armdatafactory.MappingDataFlowTypeProperties{
				ScriptLines: []*string{
					to.Ptr("source(output("),
					to.Ptr("PreviousConversionRate as double,"),
					to.Ptr("Country as string,"),
					to.Ptr("DateTime1 as string,"),
					to.Ptr("CurrentConversionRate as double"),
					to.Ptr("),"),
					to.Ptr("allowSchemaDrift: false,"),
					to.Ptr("validateSchema: false) ~> USDCurrency"),
					to.Ptr("source(output("),
					to.Ptr("PreviousConversionRate as double,"),
					to.Ptr("Country as string,"),
					to.Ptr("DateTime1 as string,"),
					to.Ptr("CurrentConversionRate as double"),
					to.Ptr("),"),
					to.Ptr("allowSchemaDrift: true,"),
					to.Ptr("validateSchema: false) ~> CADSource"),
					to.Ptr("USDCurrency, CADSource union(byName: true)~> Union"),
					to.Ptr("Union derive(NewCurrencyRate = round(CurrentConversionRate*1.25)) ~> NewCurrencyColumn"),
					to.Ptr("NewCurrencyColumn split(Country == 'USD',"),
					to.Ptr("Country == 'CAD',disjoint: false) ~> ConditionalSplit1@(USD, CAD)"),
					to.Ptr("ConditionalSplit1@USD sink(saveMode:'overwrite' ) ~> USDSink"),
					to.Ptr("ConditionalSplit1@CAD sink(saveMode:'overwrite' ) ~> CADSink")},
				Sinks: []*armdatafactory.DataFlowSink{
					{
						Name: to.Ptr("USDSink"),
						Dataset: &armdatafactory.DatasetReference{
							Type:          to.Ptr(armdatafactory.DatasetReferenceTypeDatasetReference),
							ReferenceName: to.Ptr(testsuite.datasetName),
						},
					}},
			},
		},
	}, &armdatafactory.DataFlowsClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step DataFlows_ListByFactory
	fmt.Println("Call operation: DataFlows_ListByFactory")
	dataFlowsClientNewListByFactoryPager := dataFlowsClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for dataFlowsClientNewListByFactoryPager.More() {
		_, err := dataFlowsClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataFlows_Get
	fmt.Println("Call operation: DataFlows_Get")
	_, err = dataFlowsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.dataFlowName, &armdatafactory.DataFlowsClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step DataFlows_Delete
	fmt.Println("Call operation: DataFlows_Delete")
	_, err = dataFlowsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.dataFlowName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/createDataFlowDebugSession
func (testsuite *DatafactoryTestSuite) TestDataFlowDebugSession() {
	var err error
	// From step DataFlowDebugSession_Create
	fmt.Println("Call operation: DataFlowDebugSession_Create")
	dataFlowDebugSessionClient, err := armdatafactory.NewDataFlowDebugSessionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataFlowDebugSessionClientCreateResponsePoller, err := dataFlowDebugSessionClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.CreateDataFlowDebugSessionRequest{
		IntegrationRuntime: &armdatafactory.IntegrationRuntimeDebugResource{
			Name: to.Ptr("ir1"),
			Properties: &armdatafactory.ManagedIntegrationRuntime{
				Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeManaged),
				TypeProperties: &armdatafactory.ManagedIntegrationRuntimeTypeProperties{
					ComputeProperties: &armdatafactory.IntegrationRuntimeComputeProperties{
						DataFlowProperties: &armdatafactory.IntegrationRuntimeDataFlowProperties{
							ComputeType: to.Ptr(armdatafactory.DataFlowComputeTypeGeneral),
							CoreCount:   to.Ptr[int32](48),
							TimeToLive:  to.Ptr[int32](10),
						},
						Location: to.Ptr("AutoResolve"),
					},
				},
			},
		},
		TimeToLive: to.Ptr[int32](60),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dataFlowDebugSessionClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DataFlowDebugSession_QueryByFactory
	fmt.Println("Call operation: DataFlowDebugSession_QueryByFactory")
	dataFlowDebugSessionClientNewQueryByFactoryPager := dataFlowDebugSessionClient.NewQueryByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for dataFlowDebugSessionClientNewQueryByFactoryPager.More() {
		_, err := dataFlowDebugSessionClientNewQueryByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataFlowDebugSession_Delete
	fmt.Println("Call operation: DataFlowDebugSession_Delete")
	_, err = dataFlowDebugSessionClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, armdatafactory.DeleteDataFlowDebugSessionRequest{
		SessionID: to.Ptr("91fb57e0-8292-47be-89ff-c8f2d2bb2a7e"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/managedVirtualNetworks/{managedVirtualNetworkName}
func (testsuite *DatafactoryTestSuite) TestManagedVirtualNetworks() {
	var err error
	// From step ManagedVirtualNetworks_CreateOrUpdate
	fmt.Println("Call operation: ManagedVirtualNetworks_CreateOrUpdate")
	managedVirtualNetworksClient, err := armdatafactory.NewManagedVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managedVirtualNetworksClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, armdatafactory.ManagedVirtualNetworkResource{
		Properties: &armdatafactory.ManagedVirtualNetwork{},
	}, &armdatafactory.ManagedVirtualNetworksClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ManagedVirtualNetworks_ListByFactory
	fmt.Println("Call operation: ManagedVirtualNetworks_ListByFactory")
	managedVirtualNetworksClientNewListByFactoryPager := managedVirtualNetworksClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for managedVirtualNetworksClientNewListByFactoryPager.More() {
		_, err := managedVirtualNetworksClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedVirtualNetworks_Get
	fmt.Println("Call operation: ManagedVirtualNetworks_Get")
	_, err = managedVirtualNetworksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, &armdatafactory.ManagedVirtualNetworksClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_CreateOrUpdate
	fmt.Println("Call operation: ManagedPrivateEndpoints_CreateOrUpdate")
	managedPrivateEndpointsClient, err := armdatafactory.NewManagedPrivateEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managedPrivateEndpointsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, testsuite.managedPrivateEndpointName, armdatafactory.ManagedPrivateEndpointResource{
		Properties: &armdatafactory.ManagedPrivateEndpoint{
			Fqdns:                 []*string{},
			GroupID:               to.Ptr("blob"),
			PrivateLinkResourceID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Storage/storageAccounts/exampleBlobStorage"),
		},
	}, &armdatafactory.ManagedPrivateEndpointsClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_ListByFactory
	fmt.Println("Call operation: ManagedPrivateEndpoints_ListByFactory")
	managedPrivateEndpointsClientNewListByFactoryPager := managedPrivateEndpointsClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, nil)
	for managedPrivateEndpointsClientNewListByFactoryPager.More() {
		_, err := managedPrivateEndpointsClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ManagedPrivateEndpoints_Get
	fmt.Println("Call operation: ManagedPrivateEndpoints_Get")
	_, err = managedPrivateEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, testsuite.managedPrivateEndpointName, &armdatafactory.ManagedPrivateEndpointsClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step ManagedPrivateEndpoints_Delete
	fmt.Println("Call operation: ManagedPrivateEndpoints_Delete")
	_, err = managedPrivateEndpointsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.managedVirtualNetworkName, testsuite.managedPrivateEndpointName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *DatafactoryTestSuite) TestPrivateEndpointConnection() {
	// datafactoryId := testutil.GetEnv("DATAFACTORY_ID", "")
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"datafactoryId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.datafactoryId,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointdatafactory-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "endpointdatafactory",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "datafactoryvnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2020-11-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection": false,
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix":                     "10.0.0.0/24",
								"delegations":                       []any{},
								"privateEndpointNetworkPolicies":    "Disabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
					},
					"virtualNetworkPeerings": []any{},
				},
			},
			map[string]any{
				"name":       "[parameters('networkInterfaceName')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"dnsSettings": map[string]any{
						"dnsServers": []any{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []any{
						map[string]any{
							"name": "privateEndpointIpConfig.ab24488f-044e-43f0-b9d1-af1f04071719",
							"properties": map[string]any{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]any{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
							},
						},
					},
				},
			},
			map[string]any{
				"name":       "[parameters('privateEndpointName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]any{
					"customDnsConfigs":                    []any{},
					"manualPrivateLinkServiceConnections": []any{},
					"privateLinkServiceConnections": []any{
						map[string]any{
							"name": "[parameters('privateEndpointName')]",
							"properties": map[string]any{
								"groupIds": []any{
									"dataFactory",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('datafactoryId')]",
							},
						},
					},
					"subnet": map[string]any{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
			map[string]any{
				"name":       "[concat(parameters('virtualNetworksName'), '/default')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []any{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]any{
					"addressPrefix":                     "10.0.0.0/24",
					"delegations":                       []any{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
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
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	// From step privateEndPointConnections_ListByFactory
	fmt.Println("Call operation: privateEndPointConnections_ListByFactory")
	privateEndPointConnectionsClient, err := armdatafactory.NewPrivateEndPointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndPointConnectionsClientNewListByFactoryPager := privateEndPointConnectionsClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for privateEndPointConnectionsClientNewListByFactoryPager.More() {
		nextResult, err := privateEndPointConnectionsClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnection_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpointConnection_CreateOrUpdate")
	privateEndpointConnectionClient, err := armdatafactory.NewPrivateEndpointConnectionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateEndpointConnectionClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, privateEndpointConnectionName, armdatafactory.PrivateLinkConnectionApprovalRequestResource{
		Properties: &armdatafactory.PrivateLinkConnectionApprovalRequest{
			PrivateLinkServiceConnectionState: &armdatafactory.PrivateLinkConnectionState{
				Description:     to.Ptr("Approved by admin."),
				ActionsRequired: to.Ptr(""),
				Status:          to.Ptr("Rejected"),
			},
		},
	}, &armdatafactory.PrivateEndpointConnectionClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_Get
	fmt.Println("Call operation: PrivateEndpointConnection_Get")
	_, err = privateEndpointConnectionClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, privateEndpointConnectionName, &armdatafactory.PrivateEndpointConnectionClientGetOptions{IfNoneMatch: nil})
	testsuite.Require().NoError(err)

	// From step privateLinkResources_Get
	fmt.Println("Call operation: privateLinkResources_Get")
	privateLinkResourcesClient, err := armdatafactory.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnection_Delete
	fmt.Println("Call operation: PrivateEndpointConnection_Delete")
	_, err = privateEndpointConnectionClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataFactory/factories/{factoryName}/globalParameters/{globalParameterName}
func (testsuite *DatafactoryTestSuite) TestGlobalParameters() {
	var err error
	// From step GlobalParameters_CreateOrUpdate
	fmt.Println("Call operation: GlobalParameters_CreateOrUpdate")
	globalParametersClient, err := armdatafactory.NewGlobalParametersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = globalParametersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, "default", armdatafactory.GlobalParameterResource{
		Properties: map[string]*armdatafactory.GlobalParameterSpecification{
			"waitTime": {
				Type:  to.Ptr(armdatafactory.GlobalParameterTypeInt),
				Value: float64(5),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step GlobalParameters_ListByFactory
	fmt.Println("Call operation: GlobalParameters_ListByFactory")
	globalParametersClientNewListByFactoryPager := globalParametersClient.NewListByFactoryPager(testsuite.resourceGroupName, testsuite.factoryName, nil)
	for globalParametersClientNewListByFactoryPager.More() {
		_, err := globalParametersClientNewListByFactoryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GlobalParameters_Get
	fmt.Println("Call operation: GlobalParameters_Get")
	_, err = globalParametersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, "default", nil)
	testsuite.Require().NoError(err)

	// From step GlobalParameters_Delete
	fmt.Println("Call operation: GlobalParameters_Delete")
	_, err = globalParametersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, "default", nil)
	testsuite.Require().NoError(err)
}

func (testsuite *DatafactoryTestSuite) Cleanup() {
	var err error
	// From step Pipelines_Delete
	fmt.Println("Call operation: Pipelines_Delete")
	pipelinesClient, err := armdatafactory.NewPipelinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = pipelinesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.pipelineName, nil)
	testsuite.Require().NoError(err)

	// From step Datasets_Delete
	fmt.Println("Call operation: Datasets_Delete")
	datasetsClient, err := armdatafactory.NewDatasetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = datasetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.datasetName, nil)
	testsuite.Require().NoError(err)

	// From step LinkedServices_Delete
	fmt.Println("Call operation: LinkedServices_Delete")
	linkedServicesClient, err := armdatafactory.NewLinkedServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = linkedServicesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, testsuite.linkedServiceName, nil)
	testsuite.Require().NoError(err)

	// From step Factories_Delete
	fmt.Println("Call operation: Factories_Delete")
	factoriesClient, err := armdatafactory.NewFactoriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = factoriesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.factoryName, nil)
	testsuite.Require().NoError(err)
}
