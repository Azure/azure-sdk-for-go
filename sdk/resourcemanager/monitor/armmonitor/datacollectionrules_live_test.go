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

type DatacollectionrulesTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	dataCollectionRuleName string
	managedClustersName    string
	resourceUri            string
	workspaceId            string
	azureClientId          string
	azureClientSecret      string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *DatacollectionrulesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.dataCollectionRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "datacollectionrulena", 26, false)
	testsuite.managedClustersName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "associationaks", 20, false)
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.azureClientSecret = recording.GetEnvVariable("AZURE_CLIENT_SECRET", "000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DatacollectionrulesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDatacollectionrulesTestSuite(t *testing.T) {
	suite.Run(t, new(DatacollectionrulesTestSuite))
}

func (testsuite *DatacollectionrulesTestSuite) Prepare() {
	var err error
	// From step WorkSpaces_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"workspaceId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.OperationalInsights/workspaces','workspacena')]",
			},
		},
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "workspacena",
				"type":       "Microsoft.OperationalInsights/workspaces",
				"apiVersion": "2021-12-01-preview",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"retentionInDays": float64(30),
					"sku": map[string]interface{}{
						"name": "PerNode",
					},
				},
				"tags": map[string]interface{}{
					"tag1": "value1",
				},
			},
		},
	}
	params := map[string]interface{}{
		"location": map[string]interface{}{"value": testsuite.location},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "WorkSpaces_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.workspaceId = deploymentExtend.Properties.Outputs.(map[string]interface{})["workspaceId"].(map[string]interface{})["value"].(string)

	// From step ContainerService_Create
	template = map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"resourceUri": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.ContainerService/managedClusters', parameters('managedClustersName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"azureClientId": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(azureClientId)",
			},
			"azureClientSecret": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(azureClientSecret)",
			},
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
			"managedClustersName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(managedClustersName)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('managedClustersName')]",
				"type":       "Microsoft.ContainerService/managedClusters",
				"apiVersion": "2022-06-02-preview",
				"identity": map[string]interface{}{
					"type": "SystemAssigned",
				},
				"location": "[parameters('location')]",
				"properties": map[string]interface{}{
					"agentPoolProfiles": []interface{}{
						map[string]interface{}{
							"name": "agentpool",
							"type": "VirtualMachineScaleSets",
							"availabilityZones": []interface{}{
								"1",
								"2",
								"3",
							},
							"count":              float64(1),
							"enableAutoScaling":  true,
							"enableFIPS":         false,
							"enableNodePublicIP": false,
							"kubeletDiskType":    "OS",
							"maxCount":           float64(5),
							"maxPods":            float64(110),
							"minCount":           float64(1),
							"mode":               "System",
							"osDiskSizeGB":       float64(128),
							"osDiskType":         "Managed",
							"osSKU":              "Ubuntu",
							"osType":             "Linux",
							"powerState": map[string]interface{}{
								"code": "Running",
							},
							"vmSize": "Standard_DS2_v2",
						},
						map[string]interface{}{
							"name":               "noodpool2",
							"type":               "VirtualMachineScaleSets",
							"count":              float64(1),
							"enableAutoScaling":  true,
							"enableFIPS":         false,
							"enableNodePublicIP": false,
							"kubeletDiskType":    "OS",
							"maxCount":           float64(10),
							"maxPods":            float64(110),
							"minCount":           float64(1),
							"mode":               "User",
							"osDiskSizeGB":       float64(128),
							"osDiskType":         "Managed",
							"osSKU":              "Ubuntu",
							"osType":             "Linux",
							"powerState": map[string]interface{}{
								"code": "Running",
							},
							"scaleDownMode":   "Delete",
							"upgradeSettings": map[string]interface{}{},
							"vmSize":          "Standard_DS2_v2",
						},
					},
					"dnsPrefix": "[concat(parameters('managedClustersName'), '-dns')]",
					"oidcIssuerProfile": map[string]interface{}{
						"enabled": false,
					},
					"servicePrincipalProfile": map[string]interface{}{
						"clientId": "[parameters('azureClientId')]",
						"secret":   "[parameters('azureClientSecret')]",
					},
					"storageProfile": map[string]interface{}{
						"diskCSIDriver": map[string]interface{}{
							"enabled": true,
							"version": "v1",
						},
						"fileCSIDriver": map[string]interface{}{
							"enabled": true,
						},
						"snapshotController": map[string]interface{}{
							"enabled": true,
						},
					},
					"workloadAutoScalerProfile": map[string]interface{}{},
				},
				"sku": map[string]interface{}{
					"name": "Basic",
					"tier": "Paid",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params = map[string]interface{}{
		"azureClientId":       map[string]interface{}{"value": testsuite.azureClientId},
		"azureClientSecret":   map[string]interface{}{"value": testsuite.azureClientSecret},
		"location":            map[string]interface{}{"value": testsuite.location},
		"managedClustersName": map[string]interface{}{"value": testsuite.managedClustersName},
	}
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "ContainerService_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.resourceUri = deploymentExtend.Properties.Outputs.(map[string]interface{})["resourceUri"].(map[string]interface{})["value"].(string)
}

// Microsoft.Insights/dataCollectionRules
func (testsuite *DatacollectionrulesTestSuite) TestDatacollectionrule() {
	associationName := "myAssociation"
	var dataCollectionRuleId string
	var err error
	// From step DataCollectionRules_Create
	fmt.Println("Call operation: DataCollectionRules_Create")
	dataCollectionRulesClient, err := armmonitor.NewDataCollectionRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataCollectionRulesClientCreateResponse, err := dataCollectionRulesClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionRuleName, &armmonitor.DataCollectionRulesClientCreateOptions{
		Body: &armmonitor.DataCollectionRuleResource{
			Location: to.Ptr(testsuite.location),
			Properties: &armmonitor.DataCollectionRuleResourceProperties{
				DataFlows: []*armmonitor.DataFlow{
					{
						Destinations: []*string{
							to.Ptr("centralWorkspace")},
						Streams: []*armmonitor.KnownDataFlowStreams{
							to.Ptr(armmonitor.KnownDataFlowStreamsMicrosoftPerf),
							to.Ptr(armmonitor.KnownDataFlowStreamsMicrosoftSyslog)},
					}},
				DataSources: &armmonitor.DataCollectionRuleDataSources{
					PerformanceCounters: []*armmonitor.PerfCounterDataSource{
						{
							Name: to.Ptr("cloudTeamCoreCounters"),
							CounterSpecifiers: []*string{
								to.Ptr("\\Processor(_Total)\\% Processor Time"),
								to.Ptr("\\Memory\\Committed Bytes"),
								to.Ptr("\\LogicalDisk(_Total)\\Free Megabytes"),
								to.Ptr("\\PhysicalDisk(_Total)\\Avg. Disk Queue Length")},
							SamplingFrequencyInSeconds: to.Ptr[int32](15),
							Streams: []*armmonitor.KnownPerfCounterDataSourceStreams{
								to.Ptr(armmonitor.KnownPerfCounterDataSourceStreamsMicrosoftPerf)},
						},
						{
							Name: to.Ptr("appTeamExtraCounters"),
							CounterSpecifiers: []*string{
								to.Ptr("\\Process(_Total)\\Thread Count")},
							SamplingFrequencyInSeconds: to.Ptr[int32](30),
							Streams: []*armmonitor.KnownPerfCounterDataSourceStreams{
								to.Ptr(armmonitor.KnownPerfCounterDataSourceStreamsMicrosoftPerf)},
						}},
					Syslog: []*armmonitor.SyslogDataSource{
						{
							Name: to.Ptr("cronSyslog"),
							FacilityNames: []*armmonitor.KnownSyslogDataSourceFacilityNames{
								to.Ptr(armmonitor.KnownSyslogDataSourceFacilityNamesCron)},
							LogLevels: []*armmonitor.KnownSyslogDataSourceLogLevels{
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsDebug),
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsCritical),
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsEmergency)},
							Streams: []*armmonitor.KnownSyslogDataSourceStreams{
								to.Ptr(armmonitor.KnownSyslogDataSourceStreamsMicrosoftSyslog)},
						},
						{
							Name: to.Ptr("syslogBase"),
							FacilityNames: []*armmonitor.KnownSyslogDataSourceFacilityNames{
								to.Ptr(armmonitor.KnownSyslogDataSourceFacilityNamesSyslog)},
							LogLevels: []*armmonitor.KnownSyslogDataSourceLogLevels{
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsAlert),
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsCritical),
								to.Ptr(armmonitor.KnownSyslogDataSourceLogLevelsEmergency)},
							Streams: []*armmonitor.KnownSyslogDataSourceStreams{
								to.Ptr(armmonitor.KnownSyslogDataSourceStreamsMicrosoftSyslog)},
						}},
				},
				Destinations: &armmonitor.DataCollectionRuleDestinations{
					LogAnalytics: []*armmonitor.LogAnalyticsDestination{
						{
							Name:                to.Ptr("centralWorkspace"),
							WorkspaceResourceID: to.Ptr(testsuite.workspaceId),
						}},
				},
			},
		},
	})
	testsuite.Require().NoError(err)
	dataCollectionRuleId = *dataCollectionRulesClientCreateResponse.ID

	// From step DataCollectionRules_ListBySubscription
	fmt.Println("Call operation: DataCollectionRules_ListBySubscription")
	dataCollectionRulesClientNewListBySubscriptionPager := dataCollectionRulesClient.NewListBySubscriptionPager(nil)
	for dataCollectionRulesClientNewListBySubscriptionPager.More() {
		_, err := dataCollectionRulesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionRules_ListByResourceGroup
	fmt.Println("Call operation: DataCollectionRules_ListByResourceGroup")
	dataCollectionRulesClientNewListByResourceGroupPager := dataCollectionRulesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for dataCollectionRulesClientNewListByResourceGroupPager.More() {
		_, err := dataCollectionRulesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionRules_Get
	fmt.Println("Call operation: DataCollectionRules_Get")
	_, err = dataCollectionRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionRuleName, nil)
	testsuite.Require().NoError(err)

	// From step DataCollectionRules_Update
	fmt.Println("Call operation: DataCollectionRules_Update")
	_, err = dataCollectionRulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionRuleName, &armmonitor.DataCollectionRulesClientUpdateOptions{
		Body: &armmonitor.ResourceForUpdate{
			Tags: map[string]*string{
				"tag1": to.Ptr("A"),
				"tag2": to.Ptr("B"),
				"tag3": to.Ptr("C"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step DataCollectionRuleAssociations_Create
	fmt.Println("Call operation: DataCollectionRuleAssociations_Create")
	dataCollectionRuleAssociationsClient, err := armmonitor.NewDataCollectionRuleAssociationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dataCollectionRuleAssociationsClient.Create(testsuite.ctx, testsuite.resourceUri, associationName, &armmonitor.DataCollectionRuleAssociationsClientCreateOptions{
		Body: &armmonitor.DataCollectionRuleAssociationProxyOnlyResource{
			Properties: &armmonitor.DataCollectionRuleAssociationProxyOnlyResourceProperties{
				DataCollectionRuleID: to.Ptr(dataCollectionRuleId),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step DataCollectionRuleAssociations_ListByResource
	fmt.Println("Call operation: DataCollectionRuleAssociations_ListByResource")
	dataCollectionRuleAssociationsClientNewListByResourcePager := dataCollectionRuleAssociationsClient.NewListByResourcePager(testsuite.resourceUri, nil)
	for dataCollectionRuleAssociationsClientNewListByResourcePager.More() {
		_, err := dataCollectionRuleAssociationsClientNewListByResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionRuleAssociations_ListByRule
	fmt.Println("Call operation: DataCollectionRuleAssociations_ListByRule")
	dataCollectionRuleAssociationsClientNewListByRulePager := dataCollectionRuleAssociationsClient.NewListByRulePager(testsuite.resourceGroupName, testsuite.dataCollectionRuleName, nil)
	for dataCollectionRuleAssociationsClientNewListByRulePager.More() {
		_, err := dataCollectionRuleAssociationsClientNewListByRulePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DataCollectionRuleAssociations_Get
	fmt.Println("Call operation: DataCollectionRuleAssociations_Get")
	_, err = dataCollectionRuleAssociationsClient.Get(testsuite.ctx, testsuite.resourceUri, associationName, nil)
	testsuite.Require().NoError(err)

	// From step DataCollectionRuleAssociations_Delete
	fmt.Println("Call operation: DataCollectionRuleAssociations_Delete")
	_, err = dataCollectionRuleAssociationsClient.Delete(testsuite.ctx, testsuite.resourceUri, associationName, nil)
	testsuite.Require().NoError(err)

	// From step DataCollectionRules_Delete
	fmt.Println("Call operation: DataCollectionRules_Delete")
	_, err = dataCollectionRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dataCollectionRuleName, nil)
	testsuite.Require().NoError(err)
}
