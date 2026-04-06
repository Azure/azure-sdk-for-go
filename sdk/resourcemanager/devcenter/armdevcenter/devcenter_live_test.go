// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdevcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devcenter/armdevcenter/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type DevcenterTestSuite struct {
	suite.Suite

	ctx                           context.Context
	cred                          azcore.TokenCredential
	options                       *arm.ClientOptions
	attachedNetworkConnectionName string
	catalogName                   string
	devCenterName                 string
	devcenterId                   string
	environmentTypeName           string
	galleryName                   string
	imageName                     string
	networkConnectionId           string
	networkConnectionName         string
	projectName                   string
	subnetId                      string
	location                      string
	resourceGroupName             string
	subscriptionId                string
}

func (testsuite *DevcenterTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.attachedNetworkConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attached", 14, false)
	testsuite.catalogName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "catalogn", 14, false)
	testsuite.devCenterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "devcente", 14, false)
	testsuite.environmentTypeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "environm", 14, false)
	testsuite.galleryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "galleryn", 14, false)
	testsuite.imageName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "imagenam", 14, false)
	testsuite.networkConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkc", 14, false)
	testsuite.projectName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "projectn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DevcenterTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDevcenterTestSuite(t *testing.T) {
	suite.Run(t, new(DevcenterTestSuite))
}

func (testsuite *DevcenterTestSuite) Prepare() {
	var err error
	// From step CheckNameAvailability_Execute
	fmt.Println("Call operation: CheckNameAvailability_Execute")
	checkNameAvailabilityClient, err := armdevcenter.NewCheckNameAvailabilityClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = checkNameAvailabilityClient.Execute(testsuite.ctx, armdevcenter.CheckNameAvailabilityRequest{
		Name: to.Ptr(testsuite.devCenterName),
		Type: to.Ptr("Microsoft.DevCenter/devcenters"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step DevCenters_CreateOrUpdate
	fmt.Println("Call operation: DevCenters_CreateOrUpdate")
	devCentersClient, err := armdevcenter.NewDevCentersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	devCentersClientCreateOrUpdateResponsePoller, err := devCentersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, armdevcenter.DevCenter{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"CostCode": to.Ptr("12345"),
		},
		Properties: &armdevcenter.Properties{},
	}, nil)
	testsuite.Require().NoError(err)
	var devCentersClientCreateOrUpdateResponse *armdevcenter.DevCentersClientCreateOrUpdateResponse
	devCentersClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, devCentersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.devcenterId = *devCentersClientCreateOrUpdateResponse.ID

	// From step Projects_CreateOrUpdate
	fmt.Println("Call operation: Projects_CreateOrUpdate")
	projectsClient, err := armdevcenter.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	projectsClientCreateOrUpdateResponsePoller, err := projectsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.projectName, armdevcenter.Project{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"CostCenter": to.Ptr("R&D"),
		},
		Properties: &armdevcenter.ProjectProperties{
			Description: to.Ptr("This is my first project."),
			DevCenterID: to.Ptr(testsuite.devcenterId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, projectsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Create_Subnet
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "devcenter-vnet",
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
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
							},
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_Subnet", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step NetworkConnections_CreateOrUpdate
	fmt.Println("Call operation: NetworkConnections_CreateOrUpdate")
	networkConnectionsClient, err := armdevcenter.NewNetworkConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	networkConnectionsClientCreateOrUpdateResponsePoller, err := networkConnectionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, armdevcenter.NetworkConnection{
		Location: to.Ptr(testsuite.location),
		Properties: &armdevcenter.NetworkProperties{
			DomainName:                  to.Ptr("mydomaincontroller.local"),
			DomainPassword:              to.Ptr("Password value for user"),
			DomainUsername:              to.Ptr("testuser@mydomaincontroller.local"),
			SubnetID:                    to.Ptr(testsuite.subnetId),
			DomainJoinType:              to.Ptr(armdevcenter.DomainJoinTypeHybridAzureADJoin),
			NetworkingResourceGroupName: to.Ptr("NetworkInterfaces"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var networkConnectionsClientCreateOrUpdateResponse *armdevcenter.NetworkConnectionsClientCreateOrUpdateResponse
	networkConnectionsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, networkConnectionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.networkConnectionId = *networkConnectionsClientCreateOrUpdateResponse.ID

	// From step AttachedNetworks_CreateOrUpdate
	fmt.Println("Call operation: AttachedNetworks_CreateOrUpdate")
	attachedNetworksClient, err := armdevcenter.NewAttachedNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	attachedNetworksClientCreateOrUpdateResponsePoller, err := attachedNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.attachedNetworkConnectionName, armdevcenter.AttachedNetworkConnection{
		Properties: &armdevcenter.AttachedNetworkConnectionProperties{
			NetworkConnectionID: to.Ptr(testsuite.networkConnectionId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/devcenters/{devCenterName}
func (testsuite *DevcenterTestSuite) TestDevCenters() {
	var err error
	// From step DevCenters_ListBySubscription
	fmt.Println("Call operation: DevCenters_ListBySubscription")
	devCentersClient, err := armdevcenter.NewDevCentersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	devCentersClientNewListBySubscriptionPager := devCentersClient.NewListBySubscriptionPager(&armdevcenter.DevCentersClientListBySubscriptionOptions{Top: nil})
	for devCentersClientNewListBySubscriptionPager.More() {
		_, err := devCentersClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DevCenters_ListByResourceGroup
	fmt.Println("Call operation: DevCenters_ListByResourceGroup")
	devCentersClientNewListByResourceGroupPager := devCentersClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdevcenter.DevCentersClientListByResourceGroupOptions{Top: nil})
	for devCentersClientNewListByResourceGroupPager.More() {
		_, err := devCentersClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DevCenters_Get
	fmt.Println("Call operation: DevCenters_Get")
	_, err = devCentersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, nil)
	testsuite.Require().NoError(err)

	// From step DevCenters_Update
	fmt.Println("Call operation: DevCenters_Update")
	devCentersClientUpdateResponsePoller, err := devCentersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, armdevcenter.Update{
		Tags: map[string]*string{
			"CostCode": to.Ptr("12345"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, devCentersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/projects/{projectName}
func (testsuite *DevcenterTestSuite) TestProjects() {
	var err error
	// From step Projects_ListBySubscription
	fmt.Println("Call operation: Projects_ListBySubscription")
	projectsClient, err := armdevcenter.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	projectsClientNewListBySubscriptionPager := projectsClient.NewListBySubscriptionPager(&armdevcenter.ProjectsClientListBySubscriptionOptions{Top: nil})
	for projectsClientNewListBySubscriptionPager.More() {
		_, err := projectsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Projects_ListByResourceGroup
	fmt.Println("Call operation: Projects_ListByResourceGroup")
	projectsClientNewListByResourceGroupPager := projectsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdevcenter.ProjectsClientListByResourceGroupOptions{Top: nil})
	for projectsClientNewListByResourceGroupPager.More() {
		_, err := projectsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Projects_Get
	fmt.Println("Call operation: Projects_Get")
	_, err = projectsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)

	// From step Projects_Update
	fmt.Println("Call operation: Projects_Update")
	projectsClientUpdateResponsePoller, err := projectsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.projectName, armdevcenter.ProjectUpdate{
		Tags: map[string]*string{
			"CostCenter": to.Ptr("R&D"),
		},
		Properties: &armdevcenter.ProjectUpdateProperties{
			Description: to.Ptr("This is my first project."),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, projectsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/devcenters/{devCenterName}/catalogs/{catalogName}
func (testsuite *DevcenterTestSuite) TestCatalogs() {
	var err error
	// From step Catalogs_CreateOrUpdate
	fmt.Println("Call operation: Catalogs_CreateOrUpdate")
	catalogsClient, err := armdevcenter.NewCatalogsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	catalogsClientCreateOrUpdateResponsePoller, err := catalogsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.catalogName, armdevcenter.Catalog{
		Properties: &armdevcenter.CatalogProperties{
			GitHub: &armdevcenter.GitCatalog{
				Path:             to.Ptr("/templates"),
				Branch:           to.Ptr("main"),
				SecretIdentifier: to.Ptr("https://" + testsuite.devCenterName + "kv.vault.azure.net/secrets/CentralRepoPat"),
				URI:              to.Ptr("https://github.com/" + testsuite.devCenterName + "/centralrepo-fake.git"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, catalogsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Catalogs_ListByDevCenter
	fmt.Println("Call operation: Catalogs_ListByDevCenter")
	catalogsClientNewListByDevCenterPager := catalogsClient.NewListByDevCenterPager(testsuite.resourceGroupName, testsuite.devCenterName, &armdevcenter.CatalogsClientListByDevCenterOptions{Top: nil})
	for catalogsClientNewListByDevCenterPager.More() {
		_, err := catalogsClientNewListByDevCenterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Catalogs_Get
	fmt.Println("Call operation: Catalogs_Get")
	_, err = catalogsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.catalogName, nil)
	testsuite.Require().NoError(err)

	// From step Catalogs_Update
	fmt.Println("Call operation: Catalogs_Update")
	catalogsClientUpdateResponsePoller, err := catalogsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.catalogName, armdevcenter.CatalogUpdate{
		Properties: &armdevcenter.CatalogUpdateProperties{
			GitHub: &armdevcenter.GitCatalog{
				Path: to.Ptr("/environments"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, catalogsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Catalogs_Delete
	fmt.Println("Call operation: Catalogs_Delete")
	catalogsClientDeleteResponsePoller, err := catalogsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.catalogName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, catalogsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/devcenters/{devCenterName}/environmentTypes/{environmentTypeName}
func (testsuite *DevcenterTestSuite) TestEnvironmentTypes() {
	var err error
	// From step EnvironmentTypes_CreateOrUpdate
	fmt.Println("Call operation: EnvironmentTypes_CreateOrUpdate")
	environmentTypesClient, err := armdevcenter.NewEnvironmentTypesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = environmentTypesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.environmentTypeName, armdevcenter.EnvironmentType{
		Tags: map[string]*string{
			"Owner": to.Ptr("superuser"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step EnvironmentTypes_ListByDevCenter
	fmt.Println("Call operation: EnvironmentTypes_ListByDevCenter")
	environmentTypesClientNewListByDevCenterPager := environmentTypesClient.NewListByDevCenterPager(testsuite.resourceGroupName, testsuite.devCenterName, &armdevcenter.EnvironmentTypesClientListByDevCenterOptions{Top: nil})
	for environmentTypesClientNewListByDevCenterPager.More() {
		_, err := environmentTypesClientNewListByDevCenterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EnvironmentTypes_Get
	fmt.Println("Call operation: EnvironmentTypes_Get")
	_, err = environmentTypesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.environmentTypeName, nil)
	testsuite.Require().NoError(err)

	// From step EnvironmentTypes_Update
	fmt.Println("Call operation: EnvironmentTypes_Update")
	_, err = environmentTypesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.environmentTypeName, armdevcenter.EnvironmentTypeUpdate{
		Tags: map[string]*string{
			"Owner": to.Ptr("superuser"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step EnvironmentTypes_Delete
	fmt.Println("Call operation: EnvironmentTypes_Delete")
	_, err = environmentTypesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.environmentTypeName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/operations
func (testsuite *DevcenterTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armdevcenter.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Usages_ListByLocation
	fmt.Println("Call operation: Usages_ListByLocation")
	usagesClient, err := armdevcenter.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListByLocationPager := usagesClient.NewListByLocationPager(testsuite.location, nil)
	for usagesClientNewListByLocationPager.More() {
		_, err := usagesClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DevCenter/skus
func (testsuite *DevcenterTestSuite) TestSkus() {
	var err error
	// From step Skus_ListBySubscription
	fmt.Println("Call operation: SKUs_ListBySubscription")
	sKUsClient, err := armdevcenter.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListBySubscriptionPager := sKUsClient.NewListBySubscriptionPager(&armdevcenter.SKUsClientListBySubscriptionOptions{Top: nil})
	for sKUsClientNewListBySubscriptionPager.More() {
		_, err := sKUsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.DevCenter/networkConnections/{networkConnectionName}
func (testsuite *DevcenterTestSuite) TestNetworkConnections() {
	var err error
	// From step NetworkConnections_ListBySubscription
	fmt.Println("Call operation: NetworkConnections_ListBySubscription")
	networkConnectionsClient, err := armdevcenter.NewNetworkConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	networkConnectionsClientNewListBySubscriptionPager := networkConnectionsClient.NewListBySubscriptionPager(&armdevcenter.NetworkConnectionsClientListBySubscriptionOptions{Top: nil})
	for networkConnectionsClientNewListBySubscriptionPager.More() {
		_, err := networkConnectionsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkConnections_ListHealthDetails
	fmt.Println("Call operation: NetworkConnections_ListHealthDetails")
	networkConnectionsClientNewListHealthDetailsPager := networkConnectionsClient.NewListHealthDetailsPager(testsuite.resourceGroupName, testsuite.networkConnectionName, &armdevcenter.NetworkConnectionsClientListHealthDetailsOptions{Top: nil})
	for networkConnectionsClientNewListHealthDetailsPager.More() {
		_, err := networkConnectionsClientNewListHealthDetailsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkConnections_ListOutboundNetworkDependenciesEndpoints
	fmt.Println("Call operation: NetworkConnections_ListOutboundNetworkDependenciesEndpoints")
	networkConnectionsClientNewListOutboundNetworkDependenciesEndpointsPager := networkConnectionsClient.NewListOutboundNetworkDependenciesEndpointsPager(testsuite.resourceGroupName, testsuite.networkConnectionName, &armdevcenter.NetworkConnectionsClientListOutboundNetworkDependenciesEndpointsOptions{Top: nil})
	for networkConnectionsClientNewListOutboundNetworkDependenciesEndpointsPager.More() {
		_, err := networkConnectionsClientNewListOutboundNetworkDependenciesEndpointsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkConnections_GetHealthDetails
	fmt.Println("Call operation: NetworkConnections_GetHealthDetails")
	_, err = networkConnectionsClient.GetHealthDetails(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step NetworkConnections_ListByResourceGroup
	fmt.Println("Call operation: NetworkConnections_ListByResourceGroup")
	networkConnectionsClientNewListByResourceGroupPager := networkConnectionsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdevcenter.NetworkConnectionsClientListByResourceGroupOptions{Top: nil})
	for networkConnectionsClientNewListByResourceGroupPager.More() {
		_, err := networkConnectionsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkConnections_Get
	fmt.Println("Call operation: NetworkConnections_Get")
	_, err = networkConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step NetworkConnections_Update
	fmt.Println("Call operation: NetworkConnections_Update")
	networkConnectionsClientUpdateResponsePoller, err := networkConnectionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, armdevcenter.NetworkConnectionUpdate{
		Properties: &armdevcenter.NetworkConnectionUpdateProperties{
			DomainPassword: to.Ptr("New Password value for user"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, networkConnectionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step NetworkConnections_RunHealthChecks
	fmt.Println("Call operation: NetworkConnections_RunHealthChecks")
	networkConnectionsClientRunHealthChecksResponsePoller, err := networkConnectionsClient.BeginRunHealthChecks(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, networkConnectionsClientRunHealthChecksResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DevCenter/devcenters/{devCenterName}/attachednetworks/{attachedNetworkConnectionName}
func (testsuite *DevcenterTestSuite) TestAttachedNetworks() {
	var err error
	// From step AttachedNetworks_ListByDevCenter
	fmt.Println("Call operation: AttachedNetworks_ListByDevCenter")
	attachedNetworksClient, err := armdevcenter.NewAttachedNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	attachedNetworksClientNewListByDevCenterPager := attachedNetworksClient.NewListByDevCenterPager(testsuite.resourceGroupName, testsuite.devCenterName, &armdevcenter.AttachedNetworksClientListByDevCenterOptions{Top: nil})
	for attachedNetworksClientNewListByDevCenterPager.More() {
		_, err := attachedNetworksClientNewListByDevCenterPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AttachedNetworks_ListByProject
	fmt.Println("Call operation: AttachedNetworks_ListByProject")
	attachedNetworksClientNewListByProjectPager := attachedNetworksClient.NewListByProjectPager(testsuite.resourceGroupName, testsuite.projectName, &armdevcenter.AttachedNetworksClientListByProjectOptions{Top: nil})
	for attachedNetworksClientNewListByProjectPager.More() {
		_, err := attachedNetworksClientNewListByProjectPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AttachedNetworks_GetByDevCenter
	fmt.Println("Call operation: AttachedNetworks_GetByDevCenter")
	_, err = attachedNetworksClient.GetByDevCenter(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.attachedNetworkConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step AttachedNetworks_GetByProject
	fmt.Println("Call operation: AttachedNetworks_GetByProject")
	_, err = attachedNetworksClient.GetByProject(testsuite.ctx, testsuite.resourceGroupName, testsuite.projectName, testsuite.attachedNetworkConnectionName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *DevcenterTestSuite) Cleanup() {
	var err error
	// From step AttachedNetworks_Delete
	fmt.Println("Call operation: AttachedNetworks_Delete")
	attachedNetworksClient, err := armdevcenter.NewAttachedNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	attachedNetworksClientDeleteResponsePoller, err := attachedNetworksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, testsuite.attachedNetworkConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, attachedNetworksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step NetworkConnections_Delete
	fmt.Println("Call operation: NetworkConnections_Delete")
	networkConnectionsClient, err := armdevcenter.NewNetworkConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	networkConnectionsClientDeleteResponsePoller, err := networkConnectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkConnectionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, networkConnectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Projects_Delete
	fmt.Println("Call operation: Projects_Delete")
	projectsClient, err := armdevcenter.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	projectsClientDeleteResponsePoller, err := projectsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, projectsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step DevCenters_Delete
	fmt.Println("Call operation: DevCenters_Delete")
	devCentersClient, err := armdevcenter.NewDevCentersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	devCentersClientDeleteResponsePoller, err := devCentersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.devCenterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, devCentersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
