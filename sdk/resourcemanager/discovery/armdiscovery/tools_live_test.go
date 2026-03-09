// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ToolsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
	toolName          string
}

func (testsuite *ToolsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.toolName = "test-tool"
}

func (testsuite *ToolsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ToolsTestSuite))
}

// Test listing tools by subscription
func (testsuite *ToolsTestSuite) TestToolsListBySubscription() {
	fmt.Println("Call operation: Tools_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewToolsClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test listing tools by resource group
func (testsuite *ToolsTestSuite) TestToolsListByResourceGroup() {
	fmt.Println("Call operation: Tools_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewToolsClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a tool
func (testsuite *ToolsTestSuite) TestToolsGet() {
	fmt.Println("Call operation: Tools_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewToolsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.toolName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a tool
func (testsuite *ToolsTestSuite) TestToolsCreateOrUpdate() {
	fmt.Println("Call operation: Tools_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	toolsClient := clientFactory.NewToolsClient()
	poller, err := toolsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.toolName,
		armdiscovery.Tool{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.ToolProperties{
				Version: to.Ptr("1.0.0"),
				DefinitionContent: map[string]any{
					"name":        "molpredictor",
					"description": "Molecular property prediction for single SMILES strings.",
					"version":     "1.0.0",
					"category":    "cheminformatics",
					"license":     "MIT",
					"infra": []any{
						map[string]any{
							"name":       "worker",
							"infra_type": "container",
							"image": map[string]any{
								"acr": "demodiscoveryacr.azurecr.io/molpredictor:latest",
							},
							"compute": map[string]any{
								"min_resources": map[string]any{
									"cpu": "1", "ram": "1Gi", "storage": "32", "gpu": "0",
								},
								"max_resources": map[string]any{
									"cpu": "2", "ram": "1Gi", "storage": "64", "gpu": "0",
								},
								"recommended_sku": []any{"Standard_D4s_v6"},
								"pool_type":       "static",
								"pool_size":       1,
							},
						},
					},
					"actions": []any{
						map[string]any{
							"name":        "predict",
							"description": "Predict molecular properties for SMILES strings.",
							"input_schema": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"action": map[string]any{
										"type":        "string",
										"description": "The property to predict. Must be one of [log_p, boiling_point, solubility, density, critical_point]",
									},
								},
								"required": []any{"action"},
							},
							"command":    "python molpredictor.py --action {{ action }}",
							"infra_node": "worker",
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	tool, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(tool.ID)
	fmt.Println("Created tool:", *tool.Name)
}

// Test updating a tool
func (testsuite *ToolsTestSuite) TestToolsUpdate() {
	fmt.Println("Call operation: Tools_Update")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewToolsClient().BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.toolName,
		armdiscovery.Tool{
			Tags: map[string]*string{
				"updated": to.Ptr("true"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	result, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(result.ID)
	fmt.Println("Updated tool:", *result.Name)
}

// Test deleting a tool
func (testsuite *ToolsTestSuite) TestToolsDelete() {
	fmt.Println("Call operation: Tools_Delete")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewToolsClient().BeginDelete(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.toolName,
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	fmt.Println("Deleted tool:", testsuite.toolName)
}
