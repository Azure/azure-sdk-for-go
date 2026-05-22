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

type VirtualMachinesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	labName           string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *VirtualMachinesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.labName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "labna", 11, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualMachinesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualMachinesTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualMachinesTestSuite))
}

func (testsuite *VirtualMachinesTestSuite) Prepare() {
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
}

// Microsoft.LabServices/labs/{labName}/virtualMachines/{virtualMachineName}
func (testsuite *VirtualMachinesTestSuite) TestVirtualMachines() {
	var err error
	var virtualMachineName string
	// From step VirtualMachines_ListByLab
	fmt.Println("Call operation: VirtualMachines_ListByLab")
	virtualMachinesClient, err := armlabservices.NewVirtualMachinesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualMachinesClientNewListByLabPager := virtualMachinesClient.NewListByLabPager(testsuite.resourceGroupName, testsuite.labName, &armlabservices.VirtualMachinesClientListByLabOptions{Filter: nil})
	for virtualMachinesClientNewListByLabPager.More() {
		nextResult, err := virtualMachinesClientNewListByLabPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		virtualMachineName = *nextResult.Value[0].Name
		break
	}

	// From step VirtualMachines_Get
	fmt.Println("Call operation: VirtualMachines_Get")
	_, err = virtualMachinesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Start
	fmt.Println("Call operation: VirtualMachines_Start")
	virtualMachinesClientStartResponsePoller, err := virtualMachinesClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientStartResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Redeploy
	fmt.Println("Call operation: VirtualMachines_Redeploy")
	virtualMachinesClientRedeployResponsePoller, err := virtualMachinesClient.BeginRedeploy(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientRedeployResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Reimage
	fmt.Println("Call operation: VirtualMachines_Reimage")
	virtualMachinesClientReimageResponsePoller, err := virtualMachinesClient.BeginReimage(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientReimageResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_ResetPassword
	fmt.Println("Call operation: VirtualMachines_ResetPassword")
	virtualMachinesClientResetPasswordResponsePoller, err := virtualMachinesClient.BeginResetPassword(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, armlabservices.ResetPasswordBody{
		Password: to.Ptr(testsuite.adminPassword),
		Username: to.Ptr("test-user"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientResetPasswordResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualMachines_Stop
	fmt.Println("Call operation: VirtualMachines_Stop")
	virtualMachinesClientStopResponsePoller, err := virtualMachinesClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, virtualMachineName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualMachinesClientStopResponsePoller)
	testsuite.Require().NoError(err)
}
