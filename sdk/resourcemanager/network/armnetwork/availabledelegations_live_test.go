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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/stretchr/testify/suite"
)

type AvailableDelegationsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AvailableDelegationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AvailableDelegationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAvailableDelegationsTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableDelegationsTestSuite))
}

// Microsoft.Network/locations/{location}/availableDelegations
func (testsuite *AvailableDelegationsTestSuite) TestAvailableDelegations() {
	var err error
	// From step AvailableDelegations_List
	fmt.Println("Call operation: AvailableDelegations_List")
	availableDelegationsClient, err := armnetwork.NewAvailableDelegationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	availableDelegationsClientNewListPager := availableDelegationsClient.NewListPager(testsuite.location, nil)
	for availableDelegationsClientNewListPager.More() {
		_, err := availableDelegationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailableResourceGroupDelegations_List
	fmt.Println("Call operation: AvailableResourceGroupDelegations_List")
	availableResourceGroupDelegationsClient, err := armnetwork.NewAvailableResourceGroupDelegationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	availableResourceGroupDelegationsClientNewListPager := availableResourceGroupDelegationsClient.NewListPager(testsuite.location, testsuite.resourceGroupName, nil)
	for availableResourceGroupDelegationsClientNewListPager.More() {
		_, err := availableResourceGroupDelegationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
