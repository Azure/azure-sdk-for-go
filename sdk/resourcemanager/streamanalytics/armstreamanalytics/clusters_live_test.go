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

type ClustersTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	clusterName         string
	privateEndpointName string
	storageAccountId    string
	storageAccountKey   string
	storageAccountName  string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *ClustersTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, false)
	testsuite.privateEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "privatee", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "streamsc", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ClustersTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestClustersTestSuite(t *testing.T) {
	suite.Run(t, new(ClustersTestSuite))
}

func (testsuite *ClustersTestSuite) Prepare() {
	var err error
	// From step Clusters_CreateOrUpdate
	fmt.Println("Call operation: Clusters_CreateOrUpdate")
	clustersClient, err := armstreamanalytics.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientCreateOrUpdateResponsePoller, err := clustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armstreamanalytics.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
		SKU: &armstreamanalytics.ClusterSKU{
			Name:     to.Ptr(armstreamanalytics.ClusterSKUNameDefault),
			Capacity: to.Ptr[int32](36),
		},
	}, &armstreamanalytics.ClustersClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateOrUpdateResponsePoller)
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

// Microsoft.StreamAnalytics/clusters/{clusterName}
func (testsuite *ClustersTestSuite) TestClusters() {
	var err error
	// From step Clusters_ListBySubscription
	fmt.Println("Call operation: Clusters_ListBySubscription")
	clustersClient, err := armstreamanalytics.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientNewListBySubscriptionPager := clustersClient.NewListBySubscriptionPager(nil)
	for clustersClientNewListBySubscriptionPager.More() {
		_, err := clustersClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_ListByResourceGroup
	fmt.Println("Call operation: Clusters_ListByResourceGroup")
	clustersClientNewListByResourceGroupPager := clustersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for clustersClientNewListByResourceGroupPager.More() {
		_, err := clustersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Clusters_Get
	fmt.Println("Call operation: Clusters_Get")
	_, err = clustersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)

	// From step Clusters_Update
	fmt.Println("Call operation: Clusters_Update")
	clustersClientUpdateResponsePoller, err := clustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armstreamanalytics.Cluster{}, &armstreamanalytics.ClustersClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_ListStreamingJobs
	fmt.Println("Call operation: Clusters_ListStreamingJobs")
	clustersClientNewListStreamingJobsPager := clustersClient.NewListStreamingJobsPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for clustersClientNewListStreamingJobsPager.More() {
		_, err := clustersClientNewListStreamingJobsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.StreamAnalytics/clusters/{clusterName}/privateEndpoints/{privateEndpointName}
func (testsuite *ClustersTestSuite) TestPrivateEndpoint() {
	var err error
	// From step PrivateEndpoints_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpoints_CreateOrUpdate")
	privateEndpointsClient, err := armstreamanalytics.NewPrivateEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = privateEndpointsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.privateEndpointName, armstreamanalytics.PrivateEndpoint{
		Properties: &armstreamanalytics.PrivateEndpointProperties{
			ManualPrivateLinkServiceConnections: []*armstreamanalytics.PrivateLinkServiceConnection{
				{
					Properties: &armstreamanalytics.PrivateLinkServiceConnectionProperties{
						GroupIDs: []*string{
							to.Ptr("blob")},
						PrivateLinkServiceID: to.Ptr(testsuite.storageAccountId),
					},
				}},
		},
	}, &armstreamanalytics.PrivateEndpointsClientCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step PrivateEndpoints_ListByCluster
	fmt.Println("Call operation: PrivateEndpoints_ListByCluster")
	privateEndpointsClientNewListByClusterPager := privateEndpointsClient.NewListByClusterPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for privateEndpointsClientNewListByClusterPager.More() {
		_, err := privateEndpointsClientNewListByClusterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateEndpoints_Get
	fmt.Println("Call operation: PrivateEndpoints_Get")
	_, err = privateEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.privateEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpoints_Delete
	fmt.Println("Call operation: PrivateEndpoints_Delete")
	privateEndpointsClientDeleteResponsePoller, err := privateEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.privateEndpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
