// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armconnectedcache_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedcache/armconnectedcache"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

const (
	// ResourceLocation = "eastus2"
	ResourceLocation = "westus"
)

type ConnectedCacheTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ConnectedCacheTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	fmt.Println("testsuite.resourceGroupName:", testsuite.resourceGroupName)
	testsuite.Prepare()
}

func (testsuite *ConnectedCacheTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	time.Sleep(time.Second * 3)
	testutil.StopRecording(testsuite.T())
}

func TestConnectedCacheTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectedCacheTestSuite))
}

func (testsuite *ConnectedCacheTestSuite) Cleanup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armconnectedcache.NewClientFactory(recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000"), cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	IspCustomersClientDeleteResponsePoller, err := clientFactory.NewIspCustomersClient().BeginDelete(testsuite.ctx, testsuite.resourceGroupName, "MccRPTest2", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, IspCustomersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *ConnectedCacheTestSuite) TestIspCustomersCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armconnectedcache.NewClientFactory(recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000"), cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewIspCustomersClient().BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, "MccRPTest2", armconnectedcache.IspCustomerResource{
		Location: to.Ptr(testsuite.location),
		Properties: &armconnectedcache.CustomerProperty{
			Error: &armconnectedcache.ErrorDetail{},
		},
		Tags: map[string]*string{
			"key1878": to.Ptr("warz"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = testutil.PollForTest(testsuite.ctx, poller)
	testsuite.Require().NoError(err)
	fmt.Println("finish TestCreateIspCustomersClient====")
}

func (testsuite *ConnectedCacheTestSuite) TestIspCustomersGet() {
	// make sure that the group has been created
	ResourceName := recording.GetEnvVariable("RESOURCE_NAME", "MccRPTest2")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armconnectedcache.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)

	_, err1 := clientFactory.NewIspCustomersClient().Get(testsuite.ctx, testsuite.resourceGroupName, ResourceName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err1)
	}
	fmt.Println("TestGetIspCustomersClient over")
}

func (testsuite *ConnectedCacheTestSuite) Prepare() {
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	// new client factory

	fmt.Println("subscriptionId", testsuite.subscriptionId, "groupName", testsuite.resourceGroupName, "location", testsuite.location)
	clientFactory, err := armresources.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	client := clientFactory.NewResourceGroupsClient()
	ctx := context.Background()

	testsuite.Require().NoError(err)
	// check whether create new group successfully
	res, err := client.CheckExistence(ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	if !res.Success {
		_, err = client.CreateOrUpdate(ctx, testsuite.resourceGroupName, armresources.ResourceGroup{
			Location: to.Ptr(testsuite.location),
		}, nil)
		testsuite.Require().NoError(err)
	}

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}
