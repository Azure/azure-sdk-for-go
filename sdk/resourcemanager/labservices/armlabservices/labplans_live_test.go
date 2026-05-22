// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlabservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
	"github.com/stretchr/testify/suite"
)

type LabPlansTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	labPlanName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *LabPlansTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.labPlanName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "labplanna", 15, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *LabPlansTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestLabPlansTestSuite(t *testing.T) {
	suite.Run(t, new(LabPlansTestSuite))
}

// Microsoft.LabServices/labPlans/{labPlanName}
func (testsuite *LabPlansTestSuite) TestLabPlans() {
	var err error
	// From step LabPlans_CreateOrUpdate
	fmt.Println("Call operation: LabPlans_CreateOrUpdate")
	labPlansClient, err := armlabservices.NewLabPlansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	labPlansClientCreateOrUpdateResponsePoller, err := labPlansClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.labPlanName, armlabservices.LabPlan{
		Location: to.Ptr(testsuite.location),
		Properties: &armlabservices.LabPlanProperties{
			AllowedRegions: []*string{
				to.Ptr("eastus"),
				to.Ptr("eastus2")},
			DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
				DisconnectDelay:          to.Ptr("PT15M"),
				IdleDelay:                to.Ptr("PT15M"),
				NoConnectDelay:           to.Ptr("PT15M"),
				ShutdownOnDisconnect:     to.Ptr(armlabservices.EnableStateEnabled),
				ShutdownOnIdle:           to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
				ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
			},
			DefaultConnectionProfile: &armlabservices.ConnectionProfile{
				ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
			},
			SupportInfo: &armlabservices.SupportInfo{
				Email:        to.Ptr("help@contoso.com"),
				Instructions: to.Ptr("Contact support for help."),
				Phone:        to.Ptr("+1-202-555-0123"),
				URL:          to.Ptr("https://help.contoso.com"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labPlansClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LabPlans_ListBySubscription
	fmt.Println("Call operation: LabPlans_ListBySubscription")
	labPlansClientNewListBySubscriptionPager := labPlansClient.NewListBySubscriptionPager(&armlabservices.LabPlansClientListBySubscriptionOptions{Filter: nil})
	for labPlansClientNewListBySubscriptionPager.More() {
		_, err := labPlansClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LabPlans_Get
	fmt.Println("Call operation: LabPlans_Get")
	_, err = labPlansClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.labPlanName, nil)
	testsuite.Require().NoError(err)

	// From step LabPlans_ListByResourceGroup
	fmt.Println("Call operation: LabPlans_ListByResourceGroup")
	labPlansClientNewListByResourceGroupPager := labPlansClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for labPlansClientNewListByResourceGroupPager.More() {
		_, err := labPlansClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LabPlans_Update
	fmt.Println("Call operation: LabPlans_Update")
	labPlansClientUpdateResponsePoller, err := labPlansClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.labPlanName, armlabservices.LabPlanUpdate{
		Properties: &armlabservices.LabPlanUpdateProperties{
			DefaultConnectionProfile: &armlabservices.ConnectionProfile{
				ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labPlansClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LabPlans_Delete
	fmt.Println("Call operation: LabPlans_Delete")
	labPlansClientDeleteResponsePoller, err := labPlansClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.labPlanName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labPlansClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
