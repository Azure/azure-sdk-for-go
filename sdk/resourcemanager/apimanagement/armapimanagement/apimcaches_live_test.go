// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimcachesTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	cacheId			string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimcachesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.cacheId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "cacheid", 13, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicecache", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimcachesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimcachesTestSuite(t *testing.T) {
	suite.Run(t, new(ApimcachesTestSuite))
}

func (testsuite *ApimcachesTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name":	to.Ptr("Contoso"),
			"Test":	to.Ptr("User"),
		},
		Location:	to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail:	to.Ptr("foo@contoso.com"),
			PublisherName:	to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:		to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity:	to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/caches
func (testsuite *ApimcachesTestSuite) TestCache() {
	var err error
	// From step Cache_CreateOrUpdate
	fmt.Println("Call operation: Cache_CreateOrUpdate")
	cacheClient, err := armapimanagement.NewCacheClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = cacheClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.cacheId, armapimanagement.CacheContract{
		Properties: &armapimanagement.CacheContractProperties{
			ConnectionString:	to.Ptr("apim.redis.cache.windows.net:6380,password=xc,ssl=True,abortConnect=False"),
			UseFromLocation:	to.Ptr("default"),
			Description:		to.Ptr("Redis cache instances in West India"),
			ResourceID:		to.Ptr("https://management.azure.com/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cache/redis/" + testsuite.serviceName),
		},
	}, &armapimanagement.CacheClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Cache_GetEntityTag
	fmt.Println("Call operation: Cache_GetEntityTag")
	_, err = cacheClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.cacheId, nil)
	testsuite.Require().NoError(err)

	// From step Cache_Get
	fmt.Println("Call operation: Cache_Get")
	_, err = cacheClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.cacheId, nil)
	testsuite.Require().NoError(err)

	// From step Cache_ListByService
	fmt.Println("Call operation: Cache_ListByService")
	cacheClientNewListByServicePager := cacheClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.CacheClientListByServiceOptions{Top: nil,
		Skip:	nil,
	})
	for cacheClientNewListByServicePager.More() {
		_, err := cacheClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Cache_Update
	fmt.Println("Call operation: Cache_Update")
	_, err = cacheClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.cacheId, "*", armapimanagement.CacheUpdateParameters{
		Properties: &armapimanagement.CacheUpdateProperties{
			Description: to.Ptr("update cache description"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Cache_Delete
	fmt.Println("Call operation: Cache_Delete")
	_, err = cacheClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.cacheId, "*", nil)
	testsuite.Require().NoError(err)
}
