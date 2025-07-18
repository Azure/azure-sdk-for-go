//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type RouteTableTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	routeTableName    string
	routeName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *RouteTableTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.routeTableName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "routetable", 16, false)
	testsuite.routeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "routename", 15, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *RouteTableTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestRouteTableTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTableTestSuite))
}

func (testsuite *RouteTableTestSuite) Prepare() {
	var err error
	// From step RouteTables_CreateOrUpdate
	fmt.Println("Call operation: RouteTables_CreateOrUpdate")
	routeTablesClient, err := armnetwork.NewRouteTablesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routeTablesClientCreateOrUpdateResponsePoller, err := routeTablesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, armnetwork.RouteTable{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routeTablesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/routeTables/{routeTableName}
func (testsuite *RouteTableTestSuite) TestRouteTables() {
	var err error
	// From step RouteTables_ListAll
	fmt.Println("Call operation: RouteTables_ListAll")
	routeTablesClient, err := armnetwork.NewRouteTablesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routeTablesClientNewListAllPager := routeTablesClient.NewListAllPager(nil)
	for routeTablesClientNewListAllPager.More() {
		_, err := routeTablesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RouteTables_List
	fmt.Println("Call operation: RouteTables_List")
	routeTablesClientNewListPager := routeTablesClient.NewListPager(testsuite.resourceGroupName, nil)
	for routeTablesClientNewListPager.More() {
		_, err := routeTablesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RouteTables_Get
	fmt.Println("Call operation: RouteTables_Get")
	_, err = routeTablesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, &armnetwork.RouteTablesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step RouteTables_UpdateTags
	fmt.Println("Call operation: RouteTables_UpdateTags")
	_, err = routeTablesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/routeTables/{routeTableName}/routes/{routeName}
func (testsuite *RouteTableTestSuite) TestRoutes() {
	var err error
	// From step Routes_CreateOrUpdate
	fmt.Println("Call operation: Routes_CreateOrUpdate")
	routesClient, err := armnetwork.NewRoutesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routesClientCreateOrUpdateResponsePoller, err := routesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, testsuite.routeName, armnetwork.Route{
		Properties: &armnetwork.RoutePropertiesFormat{
			AddressPrefix: to.Ptr("10.0.3.0/24"),
			NextHopType:   to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Routes_List
	fmt.Println("Call operation: Routes_List")
	routesClientNewListPager := routesClient.NewListPager(testsuite.resourceGroupName, testsuite.routeTableName, nil)
	for routesClientNewListPager.More() {
		_, err := routesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Routes_Get
	fmt.Println("Call operation: Routes_Get")
	_, err = routesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, testsuite.routeName, nil)
	testsuite.Require().NoError(err)

	// From step Routes_Delete
	fmt.Println("Call operation: Routes_Delete")
	routesClientDeleteResponsePoller, err := routesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, testsuite.routeName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *RouteTableTestSuite) Cleanup() {
	var err error
	// From step RouteTables_Delete
	fmt.Println("Call operation: RouteTables_Delete")
	routeTablesClient, err := armnetwork.NewRouteTablesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routeTablesClientDeleteResponsePoller, err := routeTablesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.routeTableName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routeTablesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
