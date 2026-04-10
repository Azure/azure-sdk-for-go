// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcommunication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/communication/armcommunication/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type CommunicationServicesTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	clientFactory            *armcommunication.ClientFactory
	armEndpoint              string
	communicationServiceName string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *CommunicationServicesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	var err error
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.communicationServiceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "communic", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.clientFactory, err = armcommunication.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *CommunicationServicesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestCommunicationServicesTestSuite(t *testing.T) {
	suite.Run(t, new(CommunicationServicesTestSuite))
}

// Microsoft.Communication/operations
func (testsuite *CommunicationServicesTestSuite) TestOperations() {
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient := testsuite.clientFactory.NewOperationsClient()
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Communication/communicationServices/{communicationServiceName}
func (testsuite *CommunicationServicesTestSuite) TestCommunicationServices() {
	var err error
	// From step CommunicationServices_CheckNameAvailability
	fmt.Println("Call operation: CommunicationServices_CheckNameAvailability")
	servicesClient := testsuite.clientFactory.NewServicesClient()
	_, err = servicesClient.CheckNameAvailability(testsuite.ctx, armcommunication.NameAvailabilityParameters{
		Name: to.Ptr("MyCommunicationService"),
		Type: to.Ptr("Microsoft.Communication/CommunicationServices"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_CreateOrUpdate
	fmt.Println("Call operation: CommunicationServices_CreateOrUpdate")
	servicesClientCreateOrUpdateResponsePoller, err := servicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, armcommunication.ServiceResource{
		Location: to.Ptr("Global"),
		Properties: &armcommunication.ServiceProperties{
			DataLocation: to.Ptr("United States"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_ListBySubscription
	fmt.Println("Call operation: CommunicationServices_ListBySubscription")
	servicesClientNewListBySubscriptionPager := servicesClient.NewListBySubscriptionPager(nil)
	for servicesClientNewListBySubscriptionPager.More() {
		_, err := servicesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CommunicationServices_ListByResourceGroup
	fmt.Println("Call operation: CommunicationServices_ListByResourceGroup")
	servicesClientNewListByResourceGroupPager := servicesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for servicesClientNewListByResourceGroupPager.More() {
		_, err := servicesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step CommunicationServices_Get
	fmt.Println("Call operation: CommunicationServices_Get")
	_, err = servicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, nil)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_Update
	fmt.Println("Call operation: CommunicationServices_Update")
	_, err = servicesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, armcommunication.ServiceResourceUpdate{
		Tags: map[string]*string{
			"newTag": to.Ptr("newVal"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_RegenerateKey
	fmt.Println("Call operation: CommunicationServices_RegenerateKey")
	_, err = servicesClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, armcommunication.RegenerateKeyParameters{
		KeyType: to.Ptr(armcommunication.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_ListKeys
	fmt.Println("Call operation: CommunicationServices_ListKeys")
	_, err = servicesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, nil)
	testsuite.Require().NoError(err)

	// From step CommunicationServices_Delete
	fmt.Println("Call operation: CommunicationServices_Delete")
	servicesClientDeleteResponsePoller, err := servicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.communicationServiceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
