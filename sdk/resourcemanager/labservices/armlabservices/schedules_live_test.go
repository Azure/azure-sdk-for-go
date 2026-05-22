// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlabservices_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
	"github.com/stretchr/testify/suite"
)

type SchedulesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	labName           string
	scheduleName      string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SchedulesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.labName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "labna", 11, false)
	testsuite.scheduleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "schedulena", 16, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SchedulesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSchedulesTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulesTestSuite))
}

func (testsuite *SchedulesTestSuite) Prepare() {
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

// Microsoft.LabServices/labs/{labName}/schedules/{scheduleName}
func (testsuite *SchedulesTestSuite) TestSchedules() {
	var err error
	// From step Schedules_CreateOrUpdate
	fmt.Println("Call operation: Schedules_CreateOrUpdate")
	schedulesClient, err := armlabservices.NewSchedulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = schedulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, testsuite.scheduleName, armlabservices.Schedule{
		Properties: &armlabservices.ScheduleProperties{
			Notes: to.Ptr("Schedule 1 for students"),
			RecurrencePattern: &armlabservices.RecurrencePattern{
				ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t }()),
				Frequency:      to.Ptr(armlabservices.RecurrenceFrequencyDaily),
				Interval:       to.Ptr[int32](2),
			},
			StartAt:    to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t }()),
			StopAt:     to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t }()),
			TimeZoneID: to.Ptr("America/Los_Angeles"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Schedules_ListByLab
	fmt.Println("Call operation: Schedules_ListByLab")
	schedulesClientNewListByLabPager := schedulesClient.NewListByLabPager(testsuite.resourceGroupName, testsuite.labName, &armlabservices.SchedulesClientListByLabOptions{Filter: nil})
	for schedulesClientNewListByLabPager.More() {
		_, err := schedulesClientNewListByLabPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Schedules_Get
	fmt.Println("Call operation: Schedules_Get")
	_, err = schedulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, testsuite.scheduleName, nil)
	testsuite.Require().NoError(err)

	// From step Schedules_Update
	fmt.Println("Call operation: Schedules_Update")
	_, err = schedulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, testsuite.scheduleName, armlabservices.ScheduleUpdate{
		Properties: &armlabservices.ScheduleUpdateProperties{
			RecurrencePattern: &armlabservices.RecurrencePattern{
				ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t }()),
				Frequency:      to.Ptr(armlabservices.RecurrenceFrequencyDaily),
				Interval:       to.Ptr[int32](2),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Schedules_Delete
	fmt.Println("Call operation: Schedules_Delete")
	schedulesClientDeleteResponsePoller, err := schedulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.labName, testsuite.scheduleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, schedulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
