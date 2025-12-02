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

type CustomIpPrefixTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	customIpPrefixName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *CustomIpPrefixTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.customIpPrefixName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "customippr", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *CustomIpPrefixTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestCustomIpPrefixTestSuite(t *testing.T) {
	suite.Run(t, new(CustomIpPrefixTestSuite))
}

// Microsoft.Network/customIpPrefixes/{customIpPrefixName}
func (testsuite *CustomIpPrefixTestSuite) TestCustomIpPrefixes() {
	var err error
	// From step CustomIPPrefixes_CreateOrUpdate
	fmt.Println("Call operation: CustomIPPrefixes_CreateOrUpdate")
	customIPPrefixesClient, err := armnetwork.NewCustomIPPrefixesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	customIPPrefixesClientCreateOrUpdateResponsePoller, err := customIPPrefixesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.customIpPrefixName, armnetwork.CustomIPPrefix{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.CustomIPPrefixPropertiesFormat{
			Cidr: to.Ptr("0.0.0.0/24"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, customIPPrefixesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step CustomIPPrefixes_ListAll
	fmt.Println("Call operation: CustomIPPrefixes_ListAll")
	customIPPrefixesClientNewListAllPager := customIPPrefixesClient.NewListAllPager(nil)
	for customIPPrefixesClientNewListAllPager.More() {
		_, err := customIPPrefixesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CustomIPPrefixes_List
	fmt.Println("Call operation: CustomIPPrefixes_List")
	customIPPrefixesClientNewListPager := customIPPrefixesClient.NewListPager(testsuite.resourceGroupName, nil)
	for customIPPrefixesClientNewListPager.More() {
		_, err := customIPPrefixesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CustomIPPrefixes_Get
	fmt.Println("Call operation: CustomIPPrefixes_Get")
	_, err = customIPPrefixesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.customIpPrefixName, &armnetwork.CustomIPPrefixesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step CustomIPPrefixes_UpdateTags
	fmt.Println("Call operation: CustomIPPrefixes_UpdateTags")
	_, err = customIPPrefixesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.customIpPrefixName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}
