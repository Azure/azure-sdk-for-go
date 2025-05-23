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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/stretchr/testify/suite"
)

type PublicIpPrefixTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	publicIpPrefixName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *PublicIpPrefixTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.publicIpPrefixName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "publicippr", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *PublicIpPrefixTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPublicIpPrefixTestSuite(t *testing.T) {
	suite.Run(t, new(PublicIpPrefixTestSuite))
}

// Microsoft.Network/publicIPPrefixes/{publicIpPrefixName}
func (testsuite *PublicIpPrefixTestSuite) TestPublicIpPrefixes() {
	var err error
	// From step PublicIPPrefixes_CreateOrUpdate
	fmt.Println("Call operation: PublicIPPrefixes_CreateOrUpdate")
	publicIPPrefixesClient, err := armnetwork.NewPublicIPPrefixesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	publicIPPrefixesClientCreateOrUpdateResponsePoller, err := publicIPPrefixesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpPrefixName, armnetwork.PublicIPPrefix{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.PublicIPPrefixPropertiesFormat{
			PrefixLength: to.Ptr[int32](30),
		},
		SKU: &armnetwork.PublicIPPrefixSKU{
			Name: to.Ptr(armnetwork.PublicIPPrefixSKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, publicIPPrefixesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PublicIPPrefixes_ListAll
	fmt.Println("Call operation: PublicIPPrefixes_ListAll")
	publicIPPrefixesClientNewListAllPager := publicIPPrefixesClient.NewListAllPager(nil)
	for publicIPPrefixesClientNewListAllPager.More() {
		_, err := publicIPPrefixesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PublicIPPrefixes_List
	fmt.Println("Call operation: PublicIPPrefixes_List")
	publicIPPrefixesClientNewListPager := publicIPPrefixesClient.NewListPager(testsuite.resourceGroupName, nil)
	for publicIPPrefixesClientNewListPager.More() {
		_, err := publicIPPrefixesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PublicIPPrefixes_Get
	fmt.Println("Call operation: PublicIPPrefixes_Get")
	_, err = publicIPPrefixesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpPrefixName, &armnetwork.PublicIPPrefixesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step PublicIPPrefixes_UpdateTags
	fmt.Println("Call operation: PublicIPPrefixes_UpdateTags")
	_, err = publicIPPrefixesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpPrefixName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step PublicIPPrefixes_Delete
	fmt.Println("Call operation: PublicIPPrefixes_Delete")
	publicIPPrefixesClientDeleteResponsePoller, err := publicIPPrefixesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpPrefixName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, publicIPPrefixesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
