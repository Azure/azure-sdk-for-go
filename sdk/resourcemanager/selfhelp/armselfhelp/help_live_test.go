// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armselfhelp_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/selfhelp/armselfhelp/v2"
	"github.com/stretchr/testify/suite"
)

type HelpTestSuite struct {
	suite.Suite

	ctx                     context.Context
	cred                    azcore.TokenCredential
	options                 *arm.ClientOptions
	diagnosticsResourceName string
	location                string
	resourceGroupName       string
	subscriptionId          string
}

func (testsuite *HelpTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.diagnosticsResourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "diagnosticna", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *HelpTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestHelpTestSuite(t *testing.T) {
	suite.Run(t, new(HelpTestSuite))
}

// Microsoft.Help/operations
func (testsuite *HelpTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armselfhelp.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Help/diagnostics/{diagnosticsResourceName}
func (testsuite *HelpTestSuite) TestDiagnostics() {
	var virtualNetworkId string
	var err error
	// From step Create_VirtualNetwork
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"virtualNetworkId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "selfhelpvnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
				},
				"tags": map[string]any{},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetwork", &deployment)
	testsuite.Require().NoError(err)
	virtualNetworkId = deploymentExtend.Properties.Outputs.(map[string]interface{})["virtualNetworkId"].(map[string]interface{})["value"].(string)

	// From step Diagnostics_CheckNameAvailability
	// fmt.Println("Call operation: Diagnostics_CheckNameAvailability")
	diagnosticsClient, err := armselfhelp.NewDiagnosticsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	// _, err = diagnosticsClient.CheckNameAvailability(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, &armselfhelp.DiagnosticsClientCheckNameAvailabilityOptions{CheckNameAvailabilityRequest: &armselfhelp.CheckNameAvailabilityRequest{
	// 	Name: to.Ptr(testsuite.diagnosticsResourceName),
	// 	Type: to.Ptr("Microsoft.Help/diagnostics"),
	// },
	// })
	// testsuite.Require().NoError(err)

	// From step Diagnostics_Create
	fmt.Println("Call operation: Diagnostics_Create")
	diagnosticsClientCreateResponsePoller, err := diagnosticsClient.BeginCreate(testsuite.ctx, virtualNetworkId, testsuite.diagnosticsResourceName, armselfhelp.DiagnosticResource{
		Properties: &armselfhelp.DiagnosticResourceProperties{
			Insights: []*armselfhelp.DiagnosticInvocation{
				{
					SolutionID: to.Ptr("Demo2InsightV2"),
				},
			},
			GlobalParameters: map[string]*string{
				"startTime": to.Ptr("2020-07-01"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, diagnosticsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Diagnostics_Get
	fmt.Println("Call operation: Diagnostics_Get")
	_, err = diagnosticsClient.Get(testsuite.ctx, virtualNetworkId, testsuite.diagnosticsResourceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Help/discoverySolutions
func (testsuite *HelpTestSuite) TestDiscoverySolution() {
	var err error
	// From step DiscoverySolution_List
	fmt.Println("Call operation: DiscoverySolution_List")
	discoverySolutionClient, err := armselfhelp.NewDiscoverySolutionClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	discoverySolutionClientNewListPager := discoverySolutionClient.NewListPager(&armselfhelp.DiscoverySolutionClientListOptions{Filter: nil,
		Skiptoken: nil,
	})
	for discoverySolutionClientNewListPager.More() {
		_, err := discoverySolutionClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
