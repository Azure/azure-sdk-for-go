// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type LogprofilesTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	storageAccountId   string
	storageAccountName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *LogprofilesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.storageAccountName = "monitorsana"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *LogprofilesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestLogprofilesTestSuite(t *testing.T) {
	suite.Run(t, new(LogprofilesTestSuite))
}

func (testsuite *LogprofilesTestSuite) Prepare() {
	var err error
	// From step StorageAccount_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"storageAccountId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"storageAccountName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(storageAccountName)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "eastus",
				"properties": map[string]interface{}{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]interface{}{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]interface{}{
							"blob": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]interface{}{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []interface{}{},
						"virtualNetworkRules": []interface{}{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]interface{}{
					"name": "Standard_RAGRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params := map[string]interface{}{
		"storageAccountName": map[string]interface{}{"value": testsuite.storageAccountName},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "StorageAccount_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)
}

// Microsoft.Insights/logprofiles
func (testsuite *LogprofilesTestSuite) TestLogprofiles() {
	logProfileName := "logprofilena"
	var err error
	// From step LogProfiles_CreateOrUpdate
	fmt.Println("Call operation: LogProfiles_CreateOrUpdate")
	logProfilesClient, err := armmonitor.NewLogProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = logProfilesClient.CreateOrUpdate(testsuite.ctx, logProfileName, armmonitor.LogProfileResource{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armmonitor.LogProfileProperties{
			Categories: []*string{
				to.Ptr("Write"),
				to.Ptr("Delete"),
				to.Ptr("Action")},
			Locations: []*string{
				to.Ptr("global")},
			RetentionPolicy: &armmonitor.RetentionPolicy{
				Days:    to.Ptr[int32](3),
				Enabled: to.Ptr(true),
			},
			StorageAccountID: to.Ptr(testsuite.storageAccountId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step LogProfiles_List
	fmt.Println("Call operation: LogProfiles_List")
	logProfilesClientNewListPager := logProfilesClient.NewListPager(nil)
	for logProfilesClientNewListPager.More() {
		_, err := logProfilesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LogProfiles_Delete
	fmt.Println("Call operation: LogProfiles_Delete")
	_, err = logProfilesClient.Delete(testsuite.ctx, logProfileName, nil)
	testsuite.Require().NoError(err)
}
