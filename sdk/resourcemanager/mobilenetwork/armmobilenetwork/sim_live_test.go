// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v4"
	"github.com/stretchr/testify/suite"
)

type SimTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	mobileNetworkId   string
	mobileNetworkName string
	simGroupName      string
	simName           string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SimTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.simGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "simgroup", 14, false)
	testsuite.simName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "simname", 13, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SimTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSimTestSuite(t *testing.T) {
	suite.Run(t, new(SimTestSuite))
}

func (testsuite *SimTestSuite) Prepare() {
	var err error
	// From step MobileNetworks_CreateOrUpdate
	fmt.Println("Call operation: MobileNetworks_CreateOrUpdate")
	mobileNetworksClient, err := armmobilenetwork.NewMobileNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mobileNetworksClientCreateOrUpdateResponsePoller, err := mobileNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, armmobilenetwork.MobileNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.PropertiesFormat{
			PublicLandMobileNetworkIdentifier: &armmobilenetwork.PlmnID{
				Mcc: to.Ptr("001"),
				Mnc: to.Ptr("01"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var mobileNetworksClientCreateOrUpdateResponse *armmobilenetwork.MobileNetworksClientCreateOrUpdateResponse
	mobileNetworksClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, mobileNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.mobileNetworkId = *mobileNetworksClientCreateOrUpdateResponse.ID

	// From step SimGroups_CreateOrUpdate
	fmt.Println("Call operation: SimGroups_CreateOrUpdate")
	simGroupsClient, err := armmobilenetwork.NewSimGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	simGroupsClientCreateOrUpdateResponsePoller, err := simGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, armmobilenetwork.SimGroup{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SimGroupPropertiesFormat{
			MobileNetwork: &armmobilenetwork.ResourceID{
				ID: to.Ptr(testsuite.mobileNetworkId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.MobileNetwork/simGroups/{simGroupName}/sims/{simName}
func (testsuite *SimTestSuite) TestSims() {
	var err error
	// From step Sims_CreateOrUpdate
	fmt.Println("Call operation: Sims_CreateOrUpdate")
	simsClient, err := armmobilenetwork.NewSimsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	simsClientCreateOrUpdateResponsePoller, err := simsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, testsuite.simName, armmobilenetwork.Sim{
		Properties: &armmobilenetwork.SimPropertiesFormat{
			DeviceType:                            to.Ptr("Video camera"),
			IntegratedCircuitCardIdentifier:       to.Ptr("8900000000000000000"),
			InternationalMobileSubscriberIdentity: to.Ptr("00000"),
			AuthenticationKey:                     to.Ptr("00000000000000000000000000000000"),
			OperatorKeyCode:                       to.Ptr("00000000000000000000000000000000"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Sims_ListByGroup
	fmt.Println("Call operation: Sims_ListByGroup")
	simsClientNewListByGroupPager := simsClient.NewListByGroupPager(testsuite.resourceGroupName, testsuite.simGroupName, nil)
	for simsClientNewListByGroupPager.More() {
		_, err := simsClientNewListByGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Sims_Get
	fmt.Println("Call operation: Sims_Get")
	_, err = simsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, testsuite.simName, nil)
	testsuite.Require().NoError(err)

	// From step Sims_BulkUpload
	fmt.Println("Call operation: Sims_BulkUpload")
	simsClientBulkUploadResponsePoller, err := simsClient.BeginBulkUpload(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, armmobilenetwork.SimUploadList{
		Sims: []*armmobilenetwork.SimNameAndProperties{
			{
				Name: to.Ptr("testSim"),
				Properties: &armmobilenetwork.SimPropertiesFormat{
					DeviceType:                            to.Ptr("Video camera"),
					IntegratedCircuitCardIdentifier:       to.Ptr("8900000000000000000"),
					InternationalMobileSubscriberIdentity: to.Ptr("00000"),
					AuthenticationKey:                     to.Ptr("00000000000000000000000000000000"),
					OperatorKeyCode:                       to.Ptr("00000000000000000000000000000000"),
				},
			},
			{
				Name: to.Ptr("testSim2"),
				Properties: &armmobilenetwork.SimPropertiesFormat{
					DeviceType:                            to.Ptr("Video camera"),
					IntegratedCircuitCardIdentifier:       to.Ptr("8900000000000000001"),
					InternationalMobileSubscriberIdentity: to.Ptr("00000"),
					AuthenticationKey:                     to.Ptr("00000000000000000000000000000000"),
					OperatorKeyCode:                       to.Ptr("00000000000000000000000000000000"),
				},
			}},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simsClientBulkUploadResponsePoller)
	testsuite.Require().NoError(err)

	// From step Sims_BulkDelete
	fmt.Println("Call operation: Sims_BulkDelete")
	simsClientBulkDeleteResponsePoller, err := simsClient.BeginBulkDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, armmobilenetwork.SimDeleteList{
		Sims: []*string{
			to.Ptr("testSim"),
			to.Ptr("testSim2")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simsClientBulkDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Sims_Delete
	fmt.Println("Call operation: Sims_Delete")
	simsClientDeleteResponsePoller, err := simsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.simGroupName, testsuite.simName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
