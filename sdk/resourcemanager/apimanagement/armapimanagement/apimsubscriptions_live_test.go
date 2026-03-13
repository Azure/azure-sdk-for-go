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

type ApimsubscriptionsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	productId		string
	serviceName		string
	sid			string
	subproductId		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimsubscriptionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.productId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "subproductid", 18, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicesub", 16, false)
	testsuite.sid, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sid", 9, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimsubscriptionsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimsubscriptionsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimsubscriptionsTestSuite))
}

func (testsuite *ApimsubscriptionsTestSuite) Prepare() {
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

	// From step Product_CreateOrUpdate
	fmt.Println("Call operation: Product_CreateOrUpdate")
	productClient, err := armapimanagement.NewProductClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	productClientCreateOrUpdateResponse, err := productClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.ProductContract{
		Properties: &armapimanagement.ProductContractProperties{
			DisplayName: to.Ptr("Test Template ProductName 4"),
		},
	}, &armapimanagement.ProductClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	testsuite.subproductId = *productClientCreateOrUpdateResponse.ID
}

// Microsoft.ApiManagement/service/subscriptions
func (testsuite *ApimsubscriptionsTestSuite) TestSubscription() {
	var err error
	// From step Subscription_CreateOrUpdate
	fmt.Println("Call operation: Subscription_CreateOrUpdate")
	subscriptionClient, err := armapimanagement.NewSubscriptionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = subscriptionClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, armapimanagement.SubscriptionCreateParameters{
		Properties: &armapimanagement.SubscriptionCreateParameterProperties{
			DisplayName:	to.Ptr(testsuite.sid),
			Scope:		to.Ptr(testsuite.subproductId),
		},
	}, &armapimanagement.SubscriptionClientCreateOrUpdateOptions{Notify: nil,
		IfMatch:	nil,
		AppType:	nil,
	})
	testsuite.Require().NoError(err)

	// From step Subscription_GetEntityTag
	fmt.Println("Call operation: Subscription_GetEntityTag")
	_, err = subscriptionClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, nil)
	testsuite.Require().NoError(err)

	// From step Subscription_List
	fmt.Println("Call operation: Subscription_List")
	subscriptionClientNewListPager := subscriptionClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.SubscriptionClientListOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for subscriptionClientNewListPager.More() {
		_, err := subscriptionClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Subscription_Get
	fmt.Println("Call operation: Subscription_Get")
	_, err = subscriptionClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, nil)
	testsuite.Require().NoError(err)

	// From step Subscription_Update
	fmt.Println("Call operation: Subscription_Update")
	_, err = subscriptionClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, "*", armapimanagement.SubscriptionUpdateParameters{
		Properties: &armapimanagement.SubscriptionUpdateParameterProperties{
			DisplayName: to.Ptr(testsuite.sid),
		},
	}, &armapimanagement.SubscriptionClientUpdateOptions{Notify: nil,
		AppType:	nil,
	})
	testsuite.Require().NoError(err)

	// From step Subscription_RegeneratePrimaryKey
	fmt.Println("Call operation: Subscription_RegeneratePrimaryKey")
	_, err = subscriptionClient.RegeneratePrimaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, nil)
	testsuite.Require().NoError(err)

	// From step Subscription_RegenerateSecondaryKey
	fmt.Println("Call operation: Subscription_RegenerateSecondaryKey")
	_, err = subscriptionClient.RegenerateSecondaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, nil)
	testsuite.Require().NoError(err)

	// From step Subscription_ListSecrets
	fmt.Println("Call operation: Subscription_ListSecrets")
	_, err = subscriptionClient.ListSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, nil)
	testsuite.Require().NoError(err)

	// From step Subscription_Delete
	fmt.Println("Call operation: Subscription_Delete")
	_, err = subscriptionClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.sid, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/skus
func (testsuite *ApimsubscriptionsTestSuite) TestApimanagementskus() {
	var err error
	// From step ApiManagementSkus_List
	fmt.Println("Call operation: ApiManagementSkus_List")
	sKUsClient, err := armapimanagement.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(nil)
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ApiManagement/service/settings
func (testsuite *ApimsubscriptionsTestSuite) TestTenantsettings() {
	var err error
	// From step TenantSettings_ListByService
	fmt.Println("Call operation: TenantSettings_ListByService")
	tenantSettingsClient, err := armapimanagement.NewTenantSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tenantSettingsClientNewListByServicePager := tenantSettingsClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.TenantSettingsClientListByServiceOptions{Filter: nil})
	for tenantSettingsClientNewListByServicePager.More() {
		_, err := tenantSettingsClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TenantSettings_Get
	fmt.Println("Call operation: TenantSettings_Get")
	_, err = tenantSettingsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.SettingsTypeNamePublic, nil)
	testsuite.Require().NoError(err)
}
