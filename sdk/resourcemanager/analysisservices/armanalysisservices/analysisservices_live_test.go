// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armanalysisservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/analysisservices/armanalysisservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type AnalysisservicesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	serverName        string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AnalysisservicesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.serverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverna", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AnalysisservicesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAnalysisservicesTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisservicesTestSuite))
}

// Microsoft.AnalysisServices/servers/{serverName}
func (testsuite *AnalysisservicesTestSuite) TestServers() {
	var err error
	// From step Servers_CheckNameAvailability
	fmt.Println("Call operation: Servers_CheckNameAvailability")
	serversClient, err := armanalysisservices.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = serversClient.CheckNameAvailability(testsuite.ctx, testsuite.location, armanalysisservices.CheckServerNameAvailabilityParameters{
		Name: to.Ptr("azsdktest"),
		Type: to.Ptr("Microsoft.AnalysisServices/servers"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Servers_Create
	fmt.Println("Call operation: Servers_Create")
	serversClientCreateResponsePoller, err := serversClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armanalysisservices.Server{
		Location: to.Ptr(testsuite.location),
		SKU: &armanalysisservices.ResourceSKU{
			Name:     to.Ptr("S1"),
			Capacity: to.Ptr[int32](1),
			Tier:     to.Ptr(armanalysisservices.SKUTierStandard),
		},
		Tags: map[string]*string{
			"testKey": to.Ptr("testValue"),
		},
		Properties: &armanalysisservices.ServerProperties{
			AsAdministrators: &armanalysisservices.ServerAdministrators{
				Members: []*string{
					to.Ptr("azsdktest@microsoft.com"),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_List
	fmt.Println("Call operation: Servers_List")
	serversClientNewListPager := serversClient.NewListPager(nil)
	for serversClientNewListPager.More() {
		_, err := serversClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Servers_GetDetails
	fmt.Println("Call operation: Servers_GetDetails")
	_, err = serversClient.GetDetails(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)

	// From step Servers_ListSkusForExisting
	fmt.Println("Call operation: Servers_ListSkusForExisting")
	_, err = serversClient.ListSKUsForExisting(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)

	// From step Servers_ListByResourceGroup
	fmt.Println("Call operation: Servers_ListByResourceGroup")
	serversClientNewListByResourceGroupPager := serversClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for serversClientNewListByResourceGroupPager.More() {
		_, err := serversClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Servers_ListSkusForNew
	fmt.Println("Call operation: Servers_ListSkusForNew")
	_, err = serversClient.ListSKUsForNew(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step Servers_Update
	fmt.Println("Call operation: Servers_Update")
	serversClientUpdateResponsePoller, err := serversClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armanalysisservices.ServerUpdateParameters{
		Tags: map[string]*string{
			"testKey": to.Ptr("testValue"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_Suspend
	fmt.Println("Call operation: Servers_Suspend")
	serversClientSuspendResponsePoller, err := serversClient.BeginSuspend(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientSuspendResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_Resume
	fmt.Println("Call operation: Servers_Resume")
	serversClientResumeResponsePoller, err := serversClient.BeginResume(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientResumeResponsePoller)
	testsuite.Require().NoError(err)

	// From step Servers_Delete
	fmt.Println("Call operation: Servers_Delete")
	serversClientDeleteResponsePoller, err := serversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AnalysisServices/operations
func (testsuite *AnalysisservicesTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armanalysisservices.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
