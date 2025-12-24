// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armoperationalinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v3"
	"github.com/stretchr/testify/suite"
)

type ClustersTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	clusterName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ClustersTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "oicluster", 15, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ClustersTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestClustersTestSuite(t *testing.T) {
	suite.Run(t, new(ClustersTestSuite))
}

// Microsoft.OperationalInsights/clusters
func (testsuite *ClustersTestSuite) TestCluster() {
	var err error
	// From step Clusters_Create
	fmt.Println("Call operation: Clusters_CreateOrUpdate")
	clustersClient, err := armoperationalinsights.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientCreateOrUpdateResponsePoller, err := clustersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armoperationalinsights.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("val1"),
		},
		Identity: &armoperationalinsights.ManagedServiceIdentity{
			Type: to.Ptr(armoperationalinsights.ManagedServiceIdentityTypeSystemAssigned),
		},
		SKU: &armoperationalinsights.ClusterSKU{
			Name:     to.Ptr(armoperationalinsights.ClusterSKUNameEnumCapacityReservation),
			Capacity: to.Ptr[int64](1000),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_List
	fmt.Println("Call operation: Clusters_List")
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
	clustersClientUpdateResponsePoller, err := clustersClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armoperationalinsights.ClusterPatch{
		Tags: map[string]*string{
			"tag1": to.Ptr("val2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_Delete
	fmt.Println("Call operation: Clusters_Delete")
	clustersClientDeleteResponsePoller, err := clustersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
