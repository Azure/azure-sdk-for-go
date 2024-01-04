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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type CapabilitiesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	capabilityName    string
	cosmosAccountId   string
	cosmosaccountName string
	targetName        string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *CapabilitiesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/chaos/armchaos/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.capabilityName = "Failover-1.0"
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

func (testsuite *CapabilitiesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestCapabilitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CapabilitiesTestSuite))
}

func (testsuite *CapabilitiesTestSuite) Prepare() {
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

// Microsoft.Chaos/targets/{targetName}/capabilities/{capabilityName}
func (testsuite *CapabilitiesTestSuite) TestCapabilities() {
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

	// From step Capabilities_CreateOrUpdate
	fmt.Println("Call operation: Capabilities_CreateOrUpdate")
	capabilitiesClient, err := armchaos.NewCapabilitiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = capabilitiesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, testsuite.capabilityName, armchaos.Capability{
		Properties: &armchaos.CapabilityProperties{},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Capabilities_Get
	fmt.Println("Call operation: Capabilities_Get")
	_, err = capabilitiesClient.Get(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, testsuite.capabilityName, nil)
	testsuite.Require().NoError(err)

	// From step Capabilities_List
	fmt.Println("Call operation: Capabilities_List")
	capabilitiesClientNewListPager := capabilitiesClient.NewListPager(testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, &armchaos.CapabilitiesClientListOptions{ContinuationToken: nil})
	for capabilitiesClientNewListPager.More() {
		_, err := capabilitiesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Capabilities_Delete
	fmt.Println("Call operation: Capabilities_Delete")
	_, err = capabilitiesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, parentProviderNamespace, parentResourceType, parentResourceName, testsuite.targetName, testsuite.capabilityName, nil)
	testsuite.Require().NoError(err)
}
