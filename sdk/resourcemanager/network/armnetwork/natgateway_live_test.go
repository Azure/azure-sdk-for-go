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

type NatGatewayTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	natGatewayName    string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *NatGatewayTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.natGatewayName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "natgateway", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *NatGatewayTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNatGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(NatGatewayTestSuite))
}

// Microsoft.Network/natGateways/{natGatewayName}
func (testsuite *NatGatewayTestSuite) TestNatGateways() {
	var err error
	// From step NatGateways_CreateOrUpdate
	fmt.Println("Call operation: NatGateways_CreateOrUpdate")
	natGatewaysClient, err := armnetwork.NewNatGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	natGatewaysClientCreateOrUpdateResponsePoller, err := natGatewaysClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.natGatewayName, armnetwork.NatGateway{
		Location: to.Ptr(testsuite.location),
		SKU: &armnetwork.NatGatewaySKU{
			Name: to.Ptr(armnetwork.NatGatewaySKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, natGatewaysClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step NatGateways_ListAll
	fmt.Println("Call operation: NatGateways_ListAll")
	natGatewaysClientNewListAllPager := natGatewaysClient.NewListAllPager(nil)
	for natGatewaysClientNewListAllPager.More() {
		_, err := natGatewaysClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NatGateways_Get
	fmt.Println("Call operation: NatGateways_Get")
	_, err = natGatewaysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.natGatewayName, &armnetwork.NatGatewaysClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step NatGateways_List
	fmt.Println("Call operation: NatGateways_List")
	natGatewaysClientNewListPager := natGatewaysClient.NewListPager(testsuite.resourceGroupName, nil)
	for natGatewaysClientNewListPager.More() {
		_, err := natGatewaysClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NatGateways_UpdateTags
	fmt.Println("Call operation: NatGateways_UpdateTags")
	_, err = natGatewaysClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.natGatewayName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NatGateways_Delete
	fmt.Println("Call operation: NatGateways_Delete")
	natGatewaysClientDeleteResponsePoller, err := natGatewaysClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.natGatewayName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, natGatewaysClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
