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

type ExpressRouteCircuitTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	circuitName       string
	connectionName    string
	peeringName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ExpressRouteCircuitTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.circuitName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "circuitnam", 16, false)
	testsuite.connectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "connerc", 13, false)
	testsuite.peeringName = "AzurePrivatePeering"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ExpressRouteCircuitTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestExpressRouteCircuitTestSuite(t *testing.T) {
	suite.Run(t, new(ExpressRouteCircuitTestSuite))
}

func (testsuite *ExpressRouteCircuitTestSuite) Prepare() {
	var err error
	// From step ExpressRouteCircuits_CreateOrUpdate
	fmt.Println("Call operation: ExpressRouteCircuits_CreateOrUpdate")
	expressRouteCircuitsClient, err := armnetwork.NewExpressRouteCircuitsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitsClientCreateOrUpdateResponsePoller, err := expressRouteCircuitsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, armnetwork.ExpressRouteCircuit{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.ExpressRouteCircuitPropertiesFormat{
			AllowClassicOperations: to.Ptr(false),
			Authorizations:         []*armnetwork.ExpressRouteCircuitAuthorization{},
			Peerings:               []*armnetwork.ExpressRouteCircuitPeering{},
			ServiceProviderProperties: &armnetwork.ExpressRouteCircuitServiceProviderProperties{
				BandwidthInMbps:     to.Ptr[int32](200),
				PeeringLocation:     to.Ptr("Silicon Valley"),
				ServiceProviderName: to.Ptr("Equinix"),
			},
		},
		SKU: &armnetwork.ExpressRouteCircuitSKU{
			Name:   to.Ptr("Standard_MeteredData"),
			Family: to.Ptr(armnetwork.ExpressRouteCircuitSKUFamilyMeteredData),
			Tier:   to.Ptr(armnetwork.ExpressRouteCircuitSKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuitPeerings_CreateOrUpdate
	fmt.Println("Call operation: ExpressRouteCircuitPeerings_CreateOrUpdate")
	expressRouteCircuitPeeringsClient, err := armnetwork.NewExpressRouteCircuitPeeringsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitPeeringsClientCreateOrUpdateResponsePoller, err := expressRouteCircuitPeeringsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, testsuite.peeringName, armnetwork.ExpressRouteCircuitPeering{
		Properties: &armnetwork.ExpressRouteCircuitPeeringPropertiesFormat{
			PeerASN:                    to.Ptr[int64](200),
			PrimaryPeerAddressPrefix:   to.Ptr("192.168.16.252/30"),
			SecondaryPeerAddressPrefix: to.Ptr("192.168.18.252/30"),
			VlanID:                     to.Ptr[int32](200),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitPeeringsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/expressRouteCircuits/{circuitName}
func (testsuite *ExpressRouteCircuitTestSuite) TestExpressRouteCircuits() {
	var err error
	// From step ExpressRouteCircuits_ListAll
	fmt.Println("Call operation: ExpressRouteCircuits_ListAll")
	expressRouteCircuitsClient, err := armnetwork.NewExpressRouteCircuitsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitsClientNewListAllPager := expressRouteCircuitsClient.NewListAllPager(nil)
	for expressRouteCircuitsClientNewListAllPager.More() {
		_, err := expressRouteCircuitsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ExpressRouteCircuits_List
	fmt.Println("Call operation: ExpressRouteCircuits_List")
	expressRouteCircuitsClientNewListPager := expressRouteCircuitsClient.NewListPager(testsuite.resourceGroupName, nil)
	for expressRouteCircuitsClientNewListPager.More() {
		_, err := expressRouteCircuitsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ExpressRouteCircuits_Get
	fmt.Println("Call operation: ExpressRouteCircuits_Get")
	_, err = expressRouteCircuitsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, nil)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuits_GetStats
	fmt.Println("Call operation: ExpressRouteCircuits_GetStats")
	_, err = expressRouteCircuitsClient.GetStats(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, nil)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuits_UpdateTags
	fmt.Println("Call operation: ExpressRouteCircuits_UpdateTags")
	_, err = expressRouteCircuitsClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuits_GetPeeringStats
	fmt.Println("Call operation: ExpressRouteCircuits_GetPeeringStats")
	_, err = expressRouteCircuitsClient.GetPeeringStats(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, testsuite.peeringName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/expressRouteCircuits/{circuitName}/peerings/{peeringName}
func (testsuite *ExpressRouteCircuitTestSuite) TestExpressRouteCircuitPeerings() {
	var err error
	// From step ExpressRouteCircuitPeerings_List
	fmt.Println("Call operation: ExpressRouteCircuitPeerings_List")
	expressRouteCircuitPeeringsClient, err := armnetwork.NewExpressRouteCircuitPeeringsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitPeeringsClientNewListPager := expressRouteCircuitPeeringsClient.NewListPager(testsuite.resourceGroupName, testsuite.circuitName, nil)
	for expressRouteCircuitPeeringsClientNewListPager.More() {
		_, err := expressRouteCircuitPeeringsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ExpressRouteCircuitPeerings_Get
	fmt.Println("Call operation: ExpressRouteCircuitPeerings_Get")
	_, err = expressRouteCircuitPeeringsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, testsuite.peeringName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/expressRouteCircuits/{circuitName}/authorizations/{authorizationName}
func (testsuite *ExpressRouteCircuitTestSuite) TestExpressRouteCircuitAuthorizations() {
	authorizationName := "ercauthorization"
	var err error
	// From step ExpressRouteCircuitAuthorizations_CreateOrUpdate
	fmt.Println("Call operation: ExpressRouteCircuitAuthorizations_CreateOrUpdate")
	expressRouteCircuitAuthorizationsClient, err := armnetwork.NewExpressRouteCircuitAuthorizationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitAuthorizationsClientCreateOrUpdateResponsePoller, err := expressRouteCircuitAuthorizationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, authorizationName, armnetwork.ExpressRouteCircuitAuthorization{
		Properties: &armnetwork.AuthorizationPropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitAuthorizationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuitAuthorizations_List
	fmt.Println("Call operation: ExpressRouteCircuitAuthorizations_List")
	expressRouteCircuitAuthorizationsClientNewListPager := expressRouteCircuitAuthorizationsClient.NewListPager(testsuite.resourceGroupName, testsuite.circuitName, nil)
	for expressRouteCircuitAuthorizationsClientNewListPager.More() {
		_, err := expressRouteCircuitAuthorizationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ExpressRouteCircuitAuthorizations_Get
	fmt.Println("Call operation: ExpressRouteCircuitAuthorizations_Get")
	_, err = expressRouteCircuitAuthorizationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, authorizationName, nil)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuitAuthorizations_Delete
	fmt.Println("Call operation: ExpressRouteCircuitAuthorizations_Delete")
	expressRouteCircuitAuthorizationsClientDeleteResponsePoller, err := expressRouteCircuitAuthorizationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, authorizationName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitAuthorizationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/expressRouteServiceProviders
func (testsuite *ExpressRouteCircuitTestSuite) TestExpressRouteServiceProviders() {
	var err error
	// From step ExpressRouteServiceProviders_List
	fmt.Println("Call operation: ExpressRouteServiceProviders_List")
	expressRouteServiceProvidersClient, err := armnetwork.NewExpressRouteServiceProvidersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteServiceProvidersClientNewListPager := expressRouteServiceProvidersClient.NewListPager(nil)
	for expressRouteServiceProvidersClientNewListPager.More() {
		_, err := expressRouteServiceProvidersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ExpressRouteCircuitTestSuite) Cleanup() {
	var err error
	// From step ExpressRouteCircuitPeerings_Delete
	fmt.Println("Call operation: ExpressRouteCircuitPeerings_Delete")
	expressRouteCircuitPeeringsClient, err := armnetwork.NewExpressRouteCircuitPeeringsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitPeeringsClientDeleteResponsePoller, err := expressRouteCircuitPeeringsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, testsuite.peeringName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitPeeringsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step ExpressRouteCircuits_Delete
	fmt.Println("Call operation: ExpressRouteCircuits_Delete")
	expressRouteCircuitsClient, err := armnetwork.NewExpressRouteCircuitsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	expressRouteCircuitsClientDeleteResponsePoller, err := expressRouteCircuitsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.circuitName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, expressRouteCircuitsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
