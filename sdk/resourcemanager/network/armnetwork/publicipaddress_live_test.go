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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v9"
	"github.com/stretchr/testify/suite"
)

type PublicIpAddressTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	publicIpAddressName string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *PublicIpAddressTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.publicIpAddressName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "publicipad", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *PublicIpAddressTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPublicIpAddressTestSuite(t *testing.T) {
	suite.Run(t, new(PublicIpAddressTestSuite))
}

// Microsoft.Network/publicIPAddresses/{publicIpAddressName}
func (testsuite *PublicIpAddressTestSuite) TestPublicIpAddresses() {
	var err error
	// From step PublicIPAddresses_CreateOrUpdate
	fmt.Println("Call operation: PublicIPAddresses_CreateOrUpdate")
	publicIPAddressesClient, err := armnetwork.NewPublicIPAddressesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	publicIPAddressesClientCreateOrUpdateResponsePoller, err := publicIPAddressesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpAddressName, armnetwork.PublicIPAddress{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, publicIPAddressesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PublicIPAddresses_ListAll
	fmt.Println("Call operation: PublicIPAddresses_ListAll")
	publicIPAddressesClientNewListAllPager := publicIPAddressesClient.NewListAllPager(nil)
	for publicIPAddressesClientNewListAllPager.More() {
		_, err := publicIPAddressesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PublicIPAddresses_List
	fmt.Println("Call operation: PublicIPAddresses_List")
	publicIPAddressesClientNewListPager := publicIPAddressesClient.NewListPager(testsuite.resourceGroupName, nil)
	for publicIPAddressesClientNewListPager.More() {
		_, err := publicIPAddressesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PublicIPAddresses_Get
	fmt.Println("Call operation: PublicIPAddresses_Get")
	_, err = publicIPAddressesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpAddressName, &armnetwork.PublicIPAddressesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step PublicIPAddresses_UpdateTags
	fmt.Println("Call operation: PublicIPAddresses_UpdateTags")
	_, err = publicIPAddressesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpAddressName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PublicIPAddresses_Delete
	fmt.Println("Call operation: PublicIPAddresses_Delete")
	publicIPAddressesClientDeleteResponsePoller, err := publicIPAddressesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpAddressName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, publicIPAddressesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
