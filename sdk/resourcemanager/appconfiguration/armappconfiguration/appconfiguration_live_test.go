// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armappconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appconfiguration/armappconfiguration/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeployments"
	"github.com/stretchr/testify/suite"
)

type AppconfigurationTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	armEndpoint          string
	configStoreName      string
	configurationStoreId string
	keyValueName         string
	replicaName          string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *AppconfigurationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.configStoreName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "configst", 14, false)
	testsuite.keyValueName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "keyvalue", 14, false)
	testsuite.replicaName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "replican", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AppconfigurationTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAppconfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(AppconfigurationTestSuite))
}

func (testsuite *AppconfigurationTestSuite) Prepare() {
	var err error
	// From step ConfigurationStores_Create
	fmt.Println("Call operation: ConfigurationStores_Create")
	configurationStoresClient, err := armappconfiguration.NewConfigurationStoresClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationStoresClientCreateResponsePoller, err := configurationStoresClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, armappconfiguration.ConfigurationStore{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"myTag": to.Ptr("myTagValue"),
		},
		SKU: &armappconfiguration.SKU{
			Name: to.Ptr("Standard"),
		},
		Properties: &armappconfiguration.ConfigurationStoreProperties{
			DisableLocalAuth: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var configurationStoresClientCreateResponse *armappconfiguration.ConfigurationStoresClientCreateResponse
	configurationStoresClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, configurationStoresClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.configurationStoreId = *configurationStoresClientCreateResponse.ID
}

// Microsoft.AppConfiguration/operations
func (testsuite *AppconfigurationTestSuite) TestOperations() {
	var err error
	// From step Operations_CheckNameAvailability
	fmt.Println("Call operation: Operations_CheckNameAvailability")
	operationsClient, err := armappconfiguration.NewOperationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.CheckNameAvailability(testsuite.ctx, armappconfiguration.CheckNameAvailabilityParameters{
		Name: to.Ptr("contoso"),
		Type: to.Ptr(armappconfiguration.ConfigurationResourceTypeMicrosoftAppConfigurationConfigurationStores),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Operations_RegionalCheckNameAvailability
	fmt.Println("Call operation: Operations_RegionalCheckNameAvailability")
	_, err = operationsClient.RegionalCheckNameAvailability(testsuite.ctx, testsuite.location, armappconfiguration.CheckNameAvailabilityParameters{
		Name: to.Ptr("contoso"),
		Type: to.Ptr(armappconfiguration.ConfigurationResourceTypeMicrosoftAppConfigurationConfigurationStores),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClientNewListPager := operationsClient.NewListPager(&armappconfiguration.OperationsClientListOptions{SkipToken: nil})
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.AppConfiguration/configurationStores/{configStoreName}
func (testsuite *AppconfigurationTestSuite) TestConfigurationStores() {
	var keyId string
	var err error
	// From step ConfigurationStores_List
	fmt.Println("Call operation: ConfigurationStores_List")
	configurationStoresClient, err := armappconfiguration.NewConfigurationStoresClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationStoresClientNewListPager := configurationStoresClient.NewListPager(&armappconfiguration.ConfigurationStoresClientListOptions{SkipToken: nil})
	for configurationStoresClientNewListPager.More() {
		_, err := configurationStoresClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConfigurationStores_ListByResourceGroup
	fmt.Println("Call operation: ConfigurationStores_ListByResourceGroup")
	configurationStoresClientNewListByResourceGroupPager := configurationStoresClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armappconfiguration.ConfigurationStoresClientListByResourceGroupOptions{SkipToken: nil})
	for configurationStoresClientNewListByResourceGroupPager.More() {
		_, err := configurationStoresClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConfigurationStores_Get
	fmt.Println("Call operation: ConfigurationStores_Get")
	_, err = configurationStoresClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, nil)
	testsuite.Require().NoError(err)

	// From step ConfigurationStores_Update
	fmt.Println("Call operation: ConfigurationStores_Update")
	configurationStoresClientUpdateResponsePoller, err := configurationStoresClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, armappconfiguration.ConfigurationStoreUpdateParameters{
		SKU: &armappconfiguration.SKU{
			Name: to.Ptr("Standard"),
		},
		Tags: map[string]*string{
			"Category": to.Ptr("Marketing"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationStoresClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigurationStores_ListKeys
	fmt.Println("Call operation: ConfigurationStores_ListKeys")
	configurationStoresClientNewListKeysPager := configurationStoresClient.NewListKeysPager(testsuite.resourceGroupName, testsuite.configStoreName, &armappconfiguration.ConfigurationStoresClientListKeysOptions{SkipToken: nil})
	for configurationStoresClientNewListKeysPager.More() {
		nextResult, err := configurationStoresClientNewListKeysPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		if len(nextResult.Value) > 0 {
			keyId = *nextResult.Value[0].ID
			break
		}
	}
	if keyId != "" {
		// From step ConfigurationStores_RegenerateKey
		fmt.Println("Call operation: ConfigurationStores_RegenerateKey")
		_, err = configurationStoresClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, armappconfiguration.RegenerateKeyParameters{
			ID: to.Ptr(keyId),
		}, nil)
		testsuite.Require().NoError(err)
	}
}

// Microsoft.AppConfiguration/configurationStores/{configStoreName}/replicas/{replicaName}
func (testsuite *AppconfigurationTestSuite) TestReplicas() {
	var err error
	// From step Replicas_Create
	fmt.Println("Call operation: Replicas_Create")
	replicasClient, err := armappconfiguration.NewReplicasClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	replicasClientCreateResponsePoller, err := replicasClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, testsuite.replicaName, armappconfiguration.Replica{
		Location: to.Ptr("eastus"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, replicasClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Replicas_ListByConfigurationStore
	fmt.Println("Call operation: Replicas_ListByConfigurationStore")
	replicasClientNewListByConfigurationStorePager := replicasClient.NewListByConfigurationStorePager(testsuite.resourceGroupName, testsuite.configStoreName, &armappconfiguration.ReplicasClientListByConfigurationStoreOptions{SkipToken: nil})
	for replicasClientNewListByConfigurationStorePager.More() {
		_, err := replicasClientNewListByConfigurationStorePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Replicas_Get
	fmt.Println("Call operation: Replicas_Get")
	_, err = replicasClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, testsuite.replicaName, nil)
	testsuite.Require().NoError(err)

	// From step Replicas_Delete
	fmt.Println("Call operation: Replicas_Delete")
	replicasClientDeleteResponsePoller, err := replicasClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, testsuite.replicaName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, replicasClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppConfiguration/configurationStores/{configStoreName}/privateEndpointConnections/{privateEndpointConnectionName}
func (testsuite *AppconfigurationTestSuite) TestPrivateEndpointConnections() {
	var privateEndpointConnectionName string
	var err error
	// From step Create_PrivateEndpoint
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]any{
			"configurationStoreId": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.configurationStoreId,
			},
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"networkInterfaceName": map[string]any{
				"type":         "string",
				"defaultValue": "epappconf-nic",
			},
			"privateEndpointName": map[string]any{
				"type":         "string",
				"defaultValue": "epappconf",
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "epappconfvnet",
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
							"name": "privateEndpointIpConfig",
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
									"configurationStores",
								},
								"privateLinkServiceConnectionState": map[string]any{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('configurationStoreId')]",
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
	deployment := armdeployments.Deployment{
		Properties: &armdeployments.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armdeployments.DeploymentModeIncremental),
		},
	}
	_, err = createDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_PrivateEndpoint", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_ListByConfigurationStore
	fmt.Println("Call operation: PrivateEndpointConnections_ListByConfigurationStore")
	privateEndpointConnectionsClient, err := armappconfiguration.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListByConfigurationStorePager := privateEndpointConnectionsClient.NewListByConfigurationStorePager(testsuite.resourceGroupName, testsuite.configStoreName, nil)
	for privateEndpointConnectionsClientNewListByConfigurationStorePager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListByConfigurationStorePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_CreateOrUpdate
	fmt.Println("Call operation: PrivateEndpointConnections_CreateOrUpdate")
	privateEndpointConnectionsClientCreateOrUpdateResponsePoller, err := privateEndpointConnectionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, privateEndpointConnectionName, armappconfiguration.PrivateEndpointConnection{
		Properties: &armappconfiguration.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armappconfiguration.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Auto-Approved"),
				Status:      to.Ptr(armappconfiguration.ConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByConfigurationStore
	fmt.Println("Call operation: PrivateLinkResources_ListByConfigurationStore")
	privateLinkResourcesClient, err := armappconfiguration.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByConfigurationStorePager := privateLinkResourcesClient.NewListByConfigurationStorePager(testsuite.resourceGroupName, testsuite.configStoreName, nil)
	for privateLinkResourcesClientNewListByConfigurationStorePager.More() {
		_, err := privateLinkResourcesClientNewListByConfigurationStorePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PrivateLinkResources_Get
	fmt.Println("Call operation: PrivateLinkResources_Get")
	_, err = privateLinkResourcesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, "configurationStores", nil)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Delete
	fmt.Println("Call operation: PrivateEndpointConnections_Delete")
	privateEndpointConnectionsClientDeleteResponsePoller, err := privateEndpointConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *AppconfigurationTestSuite) Cleanup() {
	var err error
	// From step ConfigurationStores_Delete
	fmt.Println("Call operation: ConfigurationStores_Delete")
	configurationStoresClient, err := armappconfiguration.NewConfigurationStoresClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationStoresClientDeleteResponsePoller, err := configurationStoresClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.configStoreName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationStoresClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigurationStores_ListDeleted
	fmt.Println("Call operation: ConfigurationStores_ListDeleted")
	configurationStoresClientNewListDeletedPager := configurationStoresClient.NewListDeletedPager(nil)
	for configurationStoresClientNewListDeletedPager.More() {
		_, err := configurationStoresClientNewListDeletedPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConfigurationStores_GetDeleted
	fmt.Println("Call operation: ConfigurationStores_GetDeleted")
	_, err = configurationStoresClient.GetDeleted(testsuite.ctx, testsuite.location, testsuite.configStoreName, nil)
	testsuite.Require().NoError(err)

	// From step ConfigurationStores_PurgeDeleted
	fmt.Println("Call operation: ConfigurationStores_PurgeDeleted")
	configurationStoresClientPurgeDeletedResponsePoller, err := configurationStoresClient.BeginPurgeDeleted(testsuite.ctx, testsuite.location, testsuite.configStoreName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationStoresClientPurgeDeletedResponsePoller)
	testsuite.Require().NoError(err)
}

// CreateDeployment will create a resource using arm template.
// It will return the deployment result entity.
func createDeployment(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, resourceGroupName, deploymentName string, deployment *armdeployments.Deployment) (*armdeployments.DeploymentExtended, error) {
	deployClient, err := armdeployments.NewDeploymentsClient(subscriptionId, cred, options)
	if err != nil {
		return nil, err
	}

	poller, err := deployClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		*deployment,
		&armdeployments.DeploymentsClientBeginCreateOrUpdateOptions{},
	)
	if err != nil {
		return nil, err
	}
	res, err :=  testutil.PollForTest(ctx, poller)
	if err != nil {
		return nil, err
	}
	return &res.DeploymentExtended, nil
}
