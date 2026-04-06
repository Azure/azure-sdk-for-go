// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armservicefabric_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabric/armservicefabric/v2"
	"github.com/stretchr/testify/suite"
)

type ServicefabricTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	applicationTypeName string
	clusterName         string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *ServicefabricTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.applicationTypeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "applicat", 14, false)
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ServicefabricTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestServicefabricTestSuite(t *testing.T) {
	suite.Run(t, new(ServicefabricTestSuite))
}

func (testsuite *ServicefabricTestSuite) Prepare() {
	var err error
	// From step Clusters_CreateOrUpdate
	fmt.Println("Call operation: Clusters_CreateOrUpdate")
	clustersClient, err := armservicefabric.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientCreateOrUpdateResponsePoller, err := clustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armservicefabric.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armservicefabric.ClusterProperties{
			DiagnosticsStorageAccountConfig: &armservicefabric.DiagnosticsStorageAccountConfig{
				BlobEndpoint:            to.Ptr("https://diag.blob.core.windows.net/"),
				ProtectedAccountKeyName: to.Ptr("StorageAccountKey1"),
				QueueEndpoint:           to.Ptr("https://diag.queue.core.windows.net/"),
				StorageAccountName:      to.Ptr("diag"),
				TableEndpoint:           to.Ptr("https://diag.table.core.windows.net/"),
			},
			FabricSettings: []*armservicefabric.SettingsSectionDescription{
				{
					Name: to.Ptr("UpgradeService"),
					Parameters: []*armservicefabric.SettingsParameterDescription{
						{
							Name:  to.Ptr("AppPollIntervalInSeconds"),
							Value: to.Ptr("60"),
						}},
				}},
			ManagementEndpoint: to.Ptr("http://" + testsuite.clusterName + ".eastus.cloudapp.azure.com:19080"),
			NodeTypes: []*armservicefabric.NodeTypeDescription{
				{
					Name: to.Ptr("nt1vm"),
					ApplicationPorts: &armservicefabric.EndpointRangeDescription{
						EndPort:   to.Ptr[int32](30000),
						StartPort: to.Ptr[int32](20000),
					},
					ClientConnectionEndpointPort: to.Ptr[int32](19000),
					DurabilityLevel:              to.Ptr(armservicefabric.DurabilityLevelBronze),
					EphemeralPorts: &armservicefabric.EndpointRangeDescription{
						EndPort:   to.Ptr[int32](64000),
						StartPort: to.Ptr[int32](49000),
					},
					HTTPGatewayEndpointPort: to.Ptr[int32](19007),
					IsPrimary:               to.Ptr(true),
					VMInstanceCount:         to.Ptr[int32](5),
				}},
			ReliabilityLevel: to.Ptr(armservicefabric.ReliabilityLevelSilver),
			UpgradeMode:      to.Ptr(armservicefabric.UpgradeModeAutomatic),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceFabric/clusters/{clusterName}
func (testsuite *ServicefabricTestSuite) TestClusters() {
	var err error
	// From step Clusters_List
	fmt.Println("Call operation: Clusters_List")
	clustersClient, err := armservicefabric.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientNewListPager := clustersClient.NewListPager(nil)
	for clustersClientNewListPager.More() {
		_, err := clustersClientNewListPager.NextPage(testsuite.ctx)
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
	clustersClientUpdateResponsePoller, err := clustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armservicefabric.ClusterUpdateParameters{
		Tags: map[string]*string{
			"a": to.Ptr("b"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_ListUpgradableVersions
	fmt.Println("Call operation: Clusters_ListUpgradableVersions")
	_, err = clustersClient.ListUpgradableVersions(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, &armservicefabric.ClustersClientListUpgradableVersionsOptions{VersionsDescription: &armservicefabric.UpgradableVersionsDescription{
		TargetVersion: to.Ptr("9.1.1653.9590"),
	},
	})
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceFabric/locations/{location}/clusterVersions/{clusterVersion}
func (testsuite *ServicefabricTestSuite) TestClusterVersions() {
	var err error
	// From step ClusterVersions_List
	fmt.Println("Call operation: ClusterVersions_List")
	clusterVersionsClient, err := armservicefabric.NewClusterVersionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clusterVersionsClient.List(testsuite.ctx, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// From step ClusterVersions_Get
	fmt.Println("Call operation: ClusterVersions_Get")
	_, err = clusterVersionsClient.Get(testsuite.ctx, testsuite.location, "6.1.480.9494", nil)
	testsuite.Require().NoError(err)

	// From step ClusterVersions_ListByEnvironment
	fmt.Println("Call operation: ClusterVersions_ListByEnvironment")
	_, err = clusterVersionsClient.ListByEnvironment(testsuite.ctx, testsuite.location, armservicefabric.ClusterVersionsEnvironmentWindows, nil)
	testsuite.Require().NoError(err)

	// From step ClusterVersions_GetByEnvironment
	fmt.Println("Call operation: ClusterVersions_GetByEnvironment")
	_, err = clusterVersionsClient.GetByEnvironment(testsuite.ctx, testsuite.location, armservicefabric.ClusterVersionsEnvironmentWindows, "6.1.480.9494", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceFabric/clusters/{clusterName}/applicationTypes/{applicationTypeName}
func (testsuite *ServicefabricTestSuite) TestApplicationTypes() {
	var err error
	// From step ApplicationTypes_CreateOrUpdate
	fmt.Println("Call operation: ApplicationTypes_CreateOrUpdate")
	applicationTypesClient, err := armservicefabric.NewApplicationTypesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = applicationTypesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.applicationTypeName, armservicefabric.ApplicationTypeResource{
		Tags: map[string]*string{},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ApplicationTypes_List
	fmt.Println("Call operation: ApplicationTypes_List")
	applicationTypesClientNewListPager := applicationTypesClient.NewListPager(testsuite.resourceGroupName, testsuite.clusterName, nil)
	for applicationTypesClientNewListPager.More() {
		_, err := applicationTypesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationTypes_Get
	fmt.Println("Call operation: ApplicationTypes_Get")
	_, err = applicationTypesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.applicationTypeName, nil)
	testsuite.Require().NoError(err)

	// From step ApplicationTypes_Delete
	fmt.Println("Call operation: ApplicationTypes_Delete")
	applicationTypesClientDeleteResponsePoller, err := applicationTypesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, testsuite.applicationTypeName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationTypesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ServiceFabric/operations
func (testsuite *ServicefabricTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armservicefabric.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ServicefabricTestSuite) Cleanup() {
	var err error
	// From step Clusters_Delete
	fmt.Println("Call operation: Clusters_Delete")
	clustersClient, err := armservicefabric.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clustersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
}
