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

type LabsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	labName           string
	labPlanId         string
	labPlanName       string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *LabsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.labName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "labna", 11, false)
	testsuite.labPlanName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "labplanna", 15, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *LabsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestLabsTestSuite(t *testing.T) {
	suite.Run(t, new(LabsTestSuite))
}

func (testsuite *LabsTestSuite) Prepare() {
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
	var labPlansClientCreateOrUpdateResponse *armlabservices.LabPlansClientCreateOrUpdateResponse
	labPlansClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, labPlansClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.labPlanId = *labPlansClientCreateOrUpdateResponse.ID
}

// Microsoft.LabServices/labs/{labName}
func (testsuite *LabsTestSuite) TestLabs() {
	var err error
	// From step Labs_CreateOrUpdate
	fmt.Println("Call operation: Labs_CreateOrUpdate")
	labsClient, err := armlabservices.NewLabsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	labsClientCreateOrUpdateResponsePoller, err := labsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, armlabservices.Lab{
		Location: to.Ptr(testsuite.location),
		Properties: &armlabservices.LabProperties{
			Description: to.Ptr("This is a test lab."),
			AutoShutdownProfile: &armlabservices.AutoShutdownProfile{
				DisconnectDelay:          to.Ptr("PT15M"),
				IdleDelay:                to.Ptr("PT15M"),
				NoConnectDelay:           to.Ptr("PT15M"),
				ShutdownOnDisconnect:     to.Ptr(armlabservices.EnableStateEnabled),
				ShutdownOnIdle:           to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
				ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
			},
			ConnectionProfile: &armlabservices.ConnectionProfile{
				ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
			},
			LabPlanID: to.Ptr(testsuite.labPlanId),
			SecurityProfile: &armlabservices.SecurityProfile{
				OpenAccess: to.Ptr(armlabservices.EnableStateDisabled),
			},
			Title: to.Ptr("Test Lab"),
			VirtualMachineProfile: &armlabservices.VirtualMachineProfile{
				AdditionalCapabilities: &armlabservices.VirtualMachineAdditionalCapabilities{
					InstallGpuDrivers: to.Ptr(armlabservices.EnableStateDisabled),
				},
				AdminUser: &armlabservices.Credentials{
					Password: to.Ptr(testsuite.adminPassword),
					Username: to.Ptr("test-user"),
				},
				CreateOption: to.Ptr(armlabservices.CreateOptionTemplateVM),
				ImageReference: &armlabservices.ImageReference{
					Offer:     to.Ptr("0001-com-ubuntu-server-focal"),
					Publisher: to.Ptr("canonical"),
					SKU:       to.Ptr("20_04-lts"),
					Version:   to.Ptr("latest"),
				},
				SKU: &armlabservices.SKU{
					Name:     to.Ptr("Standard_Fsv2_2_4GB_64_S_SSD"),
					Capacity: to.Ptr[int32](0),
				},
				UsageQuota:        to.Ptr("2"),
				UseSharedPassword: to.Ptr(armlabservices.EnableStateDisabled),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Labs_ListBySubscription
	fmt.Println("Call operation: Labs_ListBySubscription")
	labsClientNewListBySubscriptionPager := labsClient.NewListBySubscriptionPager(&armlabservices.LabsClientListBySubscriptionOptions{Filter: nil})
	for labsClientNewListBySubscriptionPager.More() {
		_, err := labsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Labs_Get
	fmt.Println("Call operation: Labs_Get")
	_, err = labsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, nil)
	testsuite.Require().NoError(err)

	// From step Labs_ListByResourceGroup
	fmt.Println("Call operation: Labs_ListByResourceGroup")
	labsClientNewListByResourceGroupPager := labsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for labsClientNewListByResourceGroupPager.More() {
		_, err := labsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Labs_Update
	fmt.Println("Call operation: Labs_Update")
	labsClientUpdateResponsePoller, err := labsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, armlabservices.LabUpdate{
		Properties: &armlabservices.LabUpdateProperties{
			SecurityProfile: &armlabservices.SecurityProfile{
				OpenAccess: to.Ptr(armlabservices.EnableStateEnabled),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Labs_Delete
	fmt.Println("Call operation: Labs_Delete")
	labsClientDeleteResponsePoller, err := labsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, labsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
