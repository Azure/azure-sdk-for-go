package armconnectedcache_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedcache/armconnectedcache"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

// C:\Users\v-liujudy\go\pkg\mod\github.com\!azure\azure-sdk-for-go\sdk\resourcemanager
const (
	SubscriptionID   = "faa080af-c1d8-40ad-9cce-e1a450ca5b57"
	ResourceLocation = "eastus2"
)

type AppConnectedCacheTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AppConnectedCacheTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.Prepare()
}

func TestConnectedCacheTestSuite(t *testing.T) {
	suite.Run(t, new(AppConnectedCacheTestSuite))
}

func (testsuite *AppConnectedCacheTestSuite) TestClear() {
	clientFactory, err := armresources.NewClientFactory(testsuite.subscriptionId, testsuite.cred, nil)
	testsuite.Require().NoError(err)
	resourceGroupsClientDeleteResponse, err := clientFactory.NewResourceGroupsClient().BeginDelete(testsuite.ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	time.Sleep(time.Second * 2)
	_, err = resourceGroupsClientDeleteResponse.Poll(testsuite.ctx)
	testsuite.Require().NoError(err)
	fmt.Println("delete resource group successfully")
}

func (testsuite *AppConnectedCacheTestSuite) TestGetIspCustomersClient() {
	// make sure that the group has been created
	ResourceName := recording.GetEnvVariable("RESOURCE_NAME", "scenarioTestTempClient")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armconnectedcache.NewClientFactory(testsuite.subscriptionId, cred, nil)

	_, err1 := clientFactory.NewIspCustomersClient().Get(ctx, testsuite.resourceGroupName, ResourceName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err1)
	}
	fmt.Println("TestGetIspCustomersClient over")
}

func (testsuite *AppConnectedCacheTestSuite) Prepare() {
	resoureGroupName := os.Getenv("RESOURCE_GROUP_NAME")
	subsriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	// new client factory
	clientFactory, err := armresources.NewClientFactory(subsriptionId, cred, nil)
	testsuite.Require().NoError(err)
	client := clientFactory.NewResourceGroupsClient()
	ctx := context.Background()

	_, err = client.CreateOrUpdate(ctx, resoureGroupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	testsuite.Require().NoError(err)

	// check whether create new group successfully
	_, err = client.CheckExistence(ctx, resoureGroupName, nil)
	testsuite.Require().NoError(err)

	fmt.Println("create new resource group ", resoureGroupName, " of ", subsriptionId, "successfully")
}
