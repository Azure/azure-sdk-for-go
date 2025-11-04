//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armchaos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type ExperimentsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	cosmosAccountId   string
	cosmosaccountName string
	experimentName    string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ExperimentsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cosmosaccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "account", 13, true)
	testsuite.experimentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "experiment", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ExperimentsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestExperimentsTestSuite(t *testing.T) {
	suite.Run(t, new(ExperimentsTestSuite))
}

func (testsuite *ExperimentsTestSuite) Prepare() {
	var err error
	// From step Create_CosmosAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"cosmosAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.DocumentDB/databaseAccounts', parameters('cosmosaccountName'))]",
			},
		},
		"parameters": map[string]any{
			"cosmosaccountName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.cosmosaccountName,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('cosmosaccountName')]",
				"type":       "Microsoft.DocumentDB/databaseAccounts",
				"apiVersion": "2023-09-15",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"databaseAccountOfferType": "Standard",
					"locations": []any{
						map[string]any{
							"failoverPriority":  "0",
							"is_zone_redundant": false,
							"locationName":      "westus",
						},
						map[string]any{
							"failoverPriority": "1",
							"locationName":     "eastus",
						},
					},
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_CosmosAccount", &deployment)
	testsuite.Require().NoError(err)
	testsuite.cosmosAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["cosmosAccountId"].(map[string]interface{})["value"].(string)
}

// Microsoft.Chaos/experiments/{experimentName}
func (testsuite *ExperimentsTestSuite) TestExperiments() {
	var err error
	// From step Experiments_CreateOrUpdate
	fmt.Println("Call operation: Experiments_CreateOrUpdate")
	experimentsClient, err := armchaos.NewExperimentsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	experimentsClientCreateOrUpdateResponsePoller, err := experimentsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.experimentName, armchaos.Experiment{
		Location: to.Ptr(testsuite.location),
		Identity: &armchaos.ManagedServiceIdentity{
			Type: to.Ptr(armchaos.ManagedServiceIdentityTypeSystemAssigned),
		},
		Properties: &armchaos.ExperimentProperties{
			Selectors: []armchaos.TargetSelectorClassification{
				&armchaos.TargetListSelector{
					Type: to.Ptr(armchaos.SelectorTypeList),
					ID:   to.Ptr("selector1"),
					Targets: []*armchaos.TargetReference{
						{
							Type: to.Ptr(armchaos.TargetReferenceTypeChaosTarget),
							ID:   to.Ptr(testsuite.cosmosAccountId),
						}},
				}},
			Steps: []*armchaos.ExperimentStep{
				{
					Name: to.Ptr("step1"),
					Branches: []*armchaos.ExperimentBranch{
						{
							Name: to.Ptr("branch1"),
							Actions: []armchaos.ExperimentActionClassification{
								&armchaos.ContinuousAction{
									Name:     to.Ptr("urn:csci:microsoft:virtualMachine:shutdown/1.0"),
									Type:     to.Ptr(armchaos.ExperimentActionTypeContinuous),
									Duration: to.Ptr("PT10M"),
									Parameters: []*armchaos.KeyValuePair{
										{
											Key:   to.Ptr("abruptShutdown"),
											Value: to.Ptr("false"),
										}},
									SelectorID: to.Ptr("selector1"),
								}},
						}},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, experimentsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Experiments_Get
	fmt.Println("Call operation: Experiments_Get")
	_, err = experimentsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.experimentName, nil)
	testsuite.Require().NoError(err)

	// From step Experiments_ListAll
	fmt.Println("Call operation: Experiments_ListAll")
	experimentsClientNewListAllPager := experimentsClient.NewListAllPager(&armchaos.ExperimentsClientListAllOptions{Running: nil,
		ContinuationToken: nil,
	})
	for experimentsClientNewListAllPager.More() {
		_, err := experimentsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Experiments_List
	fmt.Println("Call operation: Experiments_List")
	experimentsClientNewListPager := experimentsClient.NewListPager(testsuite.resourceGroupName, &armchaos.ExperimentsClientListOptions{Running: nil,
		ContinuationToken: nil,
	})
	for experimentsClientNewListPager.More() {
		_, err := experimentsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Experiments_ListAllExecutions
	fmt.Println("Call operation: Experiments_ListAllExecutions")
	experimentsClientNewListAllExecutionsPager := experimentsClient.NewListAllExecutionsPager(testsuite.resourceGroupName, testsuite.experimentName, nil)
	for experimentsClientNewListAllExecutionsPager.More() {
		_, err := experimentsClientNewListAllExecutionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Experiments_Delete
	fmt.Println("Call operation: Experiments_Delete")
	experimentsClientDeleteResponsePoller, err := experimentsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.experimentName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, experimentsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
