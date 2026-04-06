// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlogic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic"
	"github.com/stretchr/testify/suite"
)

type WorkflowsTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	integrationAccountId   string
	integrationAccountName string
	workflowName           string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *WorkflowsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.integrationAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "integrat", 14, false)
	testsuite.workflowName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "workflow", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *WorkflowsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWorkflowsTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowsTestSuite))
}

func (testsuite *WorkflowsTestSuite) Prepare() {
	var err error
	// From step IntegrationAccounts_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccounts_CreateOrUpdate")
	integrationAccountsClient, err := armlogic.NewIntegrationAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	integrationAccountsClientCreateOrUpdateResponse, err := integrationAccountsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, armlogic.IntegrationAccount{
		Location:   to.Ptr(testsuite.location),
		Properties: &armlogic.IntegrationAccountProperties{},
		SKU: &armlogic.IntegrationAccountSKU{
			Name: to.Ptr(armlogic.IntegrationAccountSKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.integrationAccountId = *integrationAccountsClientCreateOrUpdateResponse.ID
}

// Microsoft.Logic/workflows/{workflowName}
func (testsuite *WorkflowsTestSuite) TestWorkflows() {
	var err error
	// From step Workflows_CreateOrUpdate
	fmt.Println("Call operation: Workflows_CreateOrUpdate")
	workflowsClient, err := armlogic.NewWorkflowsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = workflowsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, armlogic.Workflow{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armlogic.WorkflowProperties{
			Definition: map[string]any{
				"$schema": "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
				"actions": map[string]any{
					"Find_pet_by_ID": map[string]any{
						"type": "ApiConnection",
						"inputs": map[string]any{
							"path":   "/pet/@{encodeURIComponent('1')}",
							"method": "get",
							"host": map[string]any{
								"connection": map[string]any{
									"name": "@parameters('$connections')['test-custom-connector']['connectionId']",
								},
							},
						},
						"runAfter": map[string]any{},
					},
				},
				"contentVersion": "1.0.0.0",
				"outputs":        map[string]any{},
				"parameters": map[string]any{
					"$connections": map[string]any{
						"type":         "Object",
						"defaultValue": map[string]any{},
					},
				},
				"triggers": map[string]any{
					"manual": map[string]any{
						"type": "Request",
						"inputs": map[string]any{
							"schema": map[string]any{},
						},
						"kind": "Http",
					},
				},
			},
			IntegrationAccount: &armlogic.ResourceReference{
				ID: to.Ptr(testsuite.integrationAccountId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_ListBySubscription
	fmt.Println("Call operation: Workflows_ListBySubscription")
	workflowsClientNewListBySubscriptionPager := workflowsClient.NewListBySubscriptionPager(&armlogic.WorkflowsClientListBySubscriptionOptions{Top: nil,
		Filter: nil,
	})
	for workflowsClientNewListBySubscriptionPager.More() {
		_, err := workflowsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Workflows_ListByResourceGroup
	fmt.Println("Call operation: Workflows_ListByResourceGroup")
	workflowsClientNewListByResourceGroupPager := workflowsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armlogic.WorkflowsClientListByResourceGroupOptions{Top: nil,
		Filter: nil,
	})
	for workflowsClientNewListByResourceGroupPager.More() {
		_, err := workflowsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Workflows_Get
	fmt.Println("Call operation: Workflows_Get")
	_, err = workflowsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_Disable
	fmt.Println("Call operation: Workflows_Disable")
	_, err = workflowsClient.Disable(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_Enable
	fmt.Println("Call operation: Workflows_Enable")
	_, err = workflowsClient.Enable(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_ListSwagger
	fmt.Println("Call operation: Workflows_ListSwagger")
	_, err = workflowsClient.ListSwagger(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_ValidateByLocation
	fmt.Println("Call operation: Workflows_ValidateByLocation")
	_, err = workflowsClient.ValidateByLocation(testsuite.ctx, testsuite.resourceGroupName, testsuite.location, testsuite.workflowName, armlogic.Workflow{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armlogic.WorkflowProperties{
			Definition: map[string]any{
				"$schema":        "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
				"actions":        map[string]any{},
				"contentVersion": "1.0.0.0",
				"outputs":        map[string]any{},
				"parameters":     map[string]any{},
				"triggers":       map[string]any{},
			},
			IntegrationAccount: &armlogic.ResourceReference{
				ID: to.Ptr(testsuite.integrationAccountId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_ValidateByResourceGroup
	fmt.Println("Call operation: Workflows_ValidateByResourceGroup")
	_, err = workflowsClient.ValidateByResourceGroup(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, armlogic.Workflow{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armlogic.WorkflowProperties{
			Definition: map[string]any{
				"$schema":        "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
				"actions":        map[string]any{},
				"contentVersion": "1.0.0.0",
				"outputs":        map[string]any{},
				"parameters":     map[string]any{},
				"triggers":       map[string]any{},
			},
			IntegrationAccount: &armlogic.ResourceReference{
				ID: to.Ptr(testsuite.integrationAccountId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_RegenerateAccessKey
	fmt.Println("Call operation: Workflows_RegenerateAccessKey")
	_, err = workflowsClient.RegenerateAccessKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, armlogic.RegenerateActionParameter{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_ListCallbackUrl
	fmt.Println("Call operation: Workflows_ListCallbackURL")
	_, err = workflowsClient.ListCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Workflows_Delete
	fmt.Println("Call operation: Workflows_Delete")
	_, err = workflowsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.workflowName, nil)
	testsuite.Require().NoError(err)
}
