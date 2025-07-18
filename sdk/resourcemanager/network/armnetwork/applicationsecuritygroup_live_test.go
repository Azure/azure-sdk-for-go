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

type ApplicationSecurityGroupTestSuite struct {
	suite.Suite

	ctx                          context.Context
	cred                         azcore.TokenCredential
	options                      *arm.ClientOptions
	applicationSecurityGroupName string
	location                     string
	resourceGroupName            string
	subscriptionId               string
}

func (testsuite *ApplicationSecurityGroupTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.applicationSecurityGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "applicatio", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ApplicationSecurityGroupTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApplicationSecurityGroupTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSecurityGroupTestSuite))
}

// Microsoft.Network/applicationSecurityGroups/{applicationSecurityGroupName}
func (testsuite *ApplicationSecurityGroupTestSuite) TestApplicationSecurityGroups() {
	var err error
	// From step ApplicationSecurityGroups_CreateOrUpdate
	fmt.Println("Call operation: ApplicationSecurityGroups_CreateOrUpdate")
	applicationSecurityGroupsClient, err := armnetwork.NewApplicationSecurityGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationSecurityGroupsClientCreateOrUpdateResponsePoller, err := applicationSecurityGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.applicationSecurityGroupName, armnetwork.ApplicationSecurityGroup{
		Location:   to.Ptr(testsuite.location),
		Properties: &armnetwork.ApplicationSecurityGroupPropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationSecurityGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApplicationSecurityGroups_ListAll
	fmt.Println("Call operation: ApplicationSecurityGroups_ListAll")
	applicationSecurityGroupsClientNewListAllPager := applicationSecurityGroupsClient.NewListAllPager(nil)
	for applicationSecurityGroupsClientNewListAllPager.More() {
		_, err := applicationSecurityGroupsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationSecurityGroups_List
	fmt.Println("Call operation: ApplicationSecurityGroups_List")
	applicationSecurityGroupsClientNewListPager := applicationSecurityGroupsClient.NewListPager(testsuite.resourceGroupName, nil)
	for applicationSecurityGroupsClientNewListPager.More() {
		_, err := applicationSecurityGroupsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationSecurityGroups_Get
	fmt.Println("Call operation: ApplicationSecurityGroups_Get")
	_, err = applicationSecurityGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.applicationSecurityGroupName, nil)
	testsuite.Require().NoError(err)

	// From step ApplicationSecurityGroups_UpdateTags
	fmt.Println("Call operation: ApplicationSecurityGroups_UpdateTags")
	_, err = applicationSecurityGroupsClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.applicationSecurityGroupName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ApplicationSecurityGroups_Delete
	fmt.Println("Call operation: ApplicationSecurityGroups_Delete")
	applicationSecurityGroupsClientDeleteResponsePoller, err := applicationSecurityGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.applicationSecurityGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationSecurityGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
