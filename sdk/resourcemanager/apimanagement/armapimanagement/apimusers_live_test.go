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

type ApimusersTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	serviceName		string
	userId			string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimusersTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceuser", 17, false)
	testsuite.userId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "userid", 12, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimusersTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimusersTestSuite(t *testing.T) {
	suite.Run(t, new(ApimusersTestSuite))
}

func (testsuite *ApimusersTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/users
func (testsuite *ApimusersTestSuite) TestUser() {
	var err error
	// From step User_CreateOrUpdate
	fmt.Println("Call operation: User_CreateOrUpdate")
	userClient, err := armapimanagement.NewUserClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = userClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, armapimanagement.UserCreateParameters{
		Properties: &armapimanagement.UserCreateParameterProperties{
			Confirmation:	to.Ptr(armapimanagement.ConfirmationSignup),
			Email:		to.Ptr("foobar@outlook.com"),
			FirstName:	to.Ptr("foo"),
			LastName:	to.Ptr("bar"),
		},
	}, &armapimanagement.UserClientCreateOrUpdateOptions{Notify: nil,
		IfMatch:	nil,
	})
	testsuite.Require().NoError(err)

	// From step User_GetEntityTag
	fmt.Println("Call operation: User_GetEntityTag")
	_, err = userClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, nil)
	testsuite.Require().NoError(err)

	// From step User_ListByService
	fmt.Println("Call operation: User_ListByService")
	userClientNewListByServicePager := userClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.UserClientListByServiceOptions{Filter: nil,
		Top:		nil,
		Skip:		nil,
		ExpandGroups:	nil,
	})
	for userClientNewListByServicePager.More() {
		_, err := userClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step User_Get
	fmt.Println("Call operation: User_Get")
	_, err = userClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, nil)
	testsuite.Require().NoError(err)

	// From step User_Update
	fmt.Println("Call operation: User_Update")
	_, err = userClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, "*", armapimanagement.UserUpdateParameters{
		Properties: &armapimanagement.UserUpdateParametersProperties{
			Email:		to.Ptr("foobar@outlook.com"),
			FirstName:	to.Ptr("foo"),
			LastName:	to.Ptr("bar"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step User_GenerateSsoUrl
	fmt.Println("Call operation: User_GenerateSsoUrl")
	_, err = userClient.GenerateSsoURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, nil)
	testsuite.Require().NoError(err)

	// From step UserGroup_List
	fmt.Println("Call operation: UserGroup_List")
	userGroupClient, err := armapimanagement.NewUserGroupClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	userGroupClientNewListPager := userGroupClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, &armapimanagement.UserGroupClientListOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for userGroupClientNewListPager.More() {
		_, err := userGroupClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step UserSubscription_List
	fmt.Println("Call operation: UserSubscription_List")
	userSubscriptionClient, err := armapimanagement.NewUserSubscriptionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	userSubscriptionClientNewListPager := userSubscriptionClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, &armapimanagement.UserSubscriptionClientListOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for userSubscriptionClientNewListPager.More() {
		_, err := userSubscriptionClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step UserIdentities_List
	fmt.Println("Call operation: UserIdentities_List")
	userIdentitiesClient, err := armapimanagement.NewUserIdentitiesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	userIdentitiesClientNewListPager := userIdentitiesClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, nil)
	for userIdentitiesClientNewListPager.More() {
		_, err := userIdentitiesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step UserConfirmationPassword_Send
	fmt.Println("Call operation: UserConfirmationPassword_Send")
	userConfirmationPasswordClient, err := armapimanagement.NewUserConfirmationPasswordClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = userConfirmationPasswordClient.Send(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.userId, &armapimanagement.UserConfirmationPasswordClientSendOptions{AppType: nil})
	testsuite.Require().NoError(err)
}
