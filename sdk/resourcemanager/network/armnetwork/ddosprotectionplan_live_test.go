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

type DdosProtectionPlanTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	ddosProtectionPlanName string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *DdosProtectionPlanTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.ddosProtectionPlanName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "ddosprotec", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DdosProtectionPlanTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDdosProtectionPlanTestSuite(t *testing.T) {
	suite.Run(t, new(DdosProtectionPlanTestSuite))
}

// Microsoft.Network/ddosProtectionPlans/{ddosProtectionPlanName}
func (testsuite *DdosProtectionPlanTestSuite) TestDdosProtectionPlans() {
	var err error
	// From step DdosProtectionPlans_CreateOrUpdate
	fmt.Println("Call operation: DdosProtectionPlans_CreateOrUpdate")
	ddosProtectionPlansClient, err := armnetwork.NewDdosProtectionPlansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ddosProtectionPlansClientCreateOrUpdateResponsePoller, err := ddosProtectionPlansClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.ddosProtectionPlanName, armnetwork.DdosProtectionPlan{
		Location:   to.Ptr(testsuite.location),
		Properties: &armnetwork.DdosProtectionPlanPropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, ddosProtectionPlansClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DdosProtectionPlans_List
	fmt.Println("Call operation: DdosProtectionPlans_List")
	ddosProtectionPlansClientNewListPager := ddosProtectionPlansClient.NewListPager(nil)
	for ddosProtectionPlansClientNewListPager.More() {
		_, err := ddosProtectionPlansClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DdosProtectionPlans_ListByResourceGroup
	fmt.Println("Call operation: DdosProtectionPlans_ListByResourceGroup")
	ddosProtectionPlansClientNewListByResourceGroupPager := ddosProtectionPlansClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for ddosProtectionPlansClientNewListByResourceGroupPager.More() {
		_, err := ddosProtectionPlansClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DdosProtectionPlans_Get
	fmt.Println("Call operation: DdosProtectionPlans_Get")
	_, err = ddosProtectionPlansClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.ddosProtectionPlanName, nil)
	testsuite.Require().NoError(err)

	// From step DdosProtectionPlans_UpdateTags
	fmt.Println("Call operation: DdosProtectionPlans_UpdateTags")
	_, err = ddosProtectionPlansClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.ddosProtectionPlanName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step DdosProtectionPlans_Delete
	fmt.Println("Call operation: DdosProtectionPlans_Delete")
	ddosProtectionPlansClientDeleteResponsePoller, err := ddosProtectionPlansClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.ddosProtectionPlanName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, ddosProtectionPlansClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
