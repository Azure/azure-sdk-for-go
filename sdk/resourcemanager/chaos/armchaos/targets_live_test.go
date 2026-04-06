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

type TargetsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	cosmosAccountId   string
	cosmosaccountName string
	targetName        string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *TargetsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cosmosaccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "account", 13, true)
	testsuite.targetName = "Microsoft-CosmosDB"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *TargetsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestTargetsTestSuite(t *testing.T) {
	suite.Run(t, new(TargetsTestSuite))
}

func (testsuite *TargetsTestSuite) Prepare() {
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

// Microsoft.Chaos/targets/{targetName}
func (testsuite *TargetsTestSuite) TestTargets() {
	parentProviderNamespace := "Microsoft.DocumentDB"
	parentResourceName := testsuite.cosmosaccountName
	parentResourceType := "databaseAccounts"
	var err error
	// From step Targets_CreateOrUpdate
	fmt.Println("Call operation: Targets_CreateOrUpdate")
	targetsClient, err := armchaos.NewTargetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = targetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, armchaos.Target{
		Properties: map[string]any{
			"identities": []any{
				map[string]any{
					"type":    "CertificateSubjectIssuer",
					"subject": "CN=example.subject",
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Targets_List
	fmt.Println("Call operation: Targets_List")
	targetsClientNewListPager := targetsClient.NewListPager(testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, &armchaos.TargetsClientListOptions{ContinuationToken: nil})
	for targetsClientNewListPager.More() {
		_, err := targetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Targets_Get
	fmt.Println("Call operation: Targets_Get")
	_, err = targetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, nil)
	testsuite.Require().NoError(err)

	// From step Targets_Delete
	fmt.Println("Call operation: Targets_Delete")
	_, err = targetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, nil)
	testsuite.Require().NoError(err)
}
