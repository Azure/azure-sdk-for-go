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

type ApimnotificationsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	email			string
	serviceName		string
	userId			string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimnotificationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.email, _ = recording.GenerateAlphaNumericID(testsuite.T(), "email", 11, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicenotifi", 19, false)
	testsuite.userId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "userid", 12, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimnotificationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimnotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimnotificationsTestSuite))
}

func (testsuite *ApimnotificationsTestSuite) Prepare() {
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

	// From step Notification_CreateOrUpdate
	fmt.Println("Call operation: Notification_CreateOrUpdate")
	notificationClient, err := armapimanagement.NewNotificationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = notificationClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, &armapimanagement.NotificationClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/notifications
func (testsuite *ApimnotificationsTestSuite) TestNotification() {
	var err error
	// From step Notification_ListByService
	fmt.Println("Call operation: Notification_ListByService")
	notificationClient, err := armapimanagement.NewNotificationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	notificationClientNewListByServicePager := notificationClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.NotificationClientListByServiceOptions{Top: nil,
		Skip:	nil,
	})
	for notificationClientNewListByServicePager.More() {
		_, err := notificationClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Notification_Get
	fmt.Println("Call operation: Notification_Get")
	_, err = notificationClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/notifications/recipientUsers
func (testsuite *ApimnotificationsTestSuite) TestNotificationrecipientemail() {
	var err error
	// From step NotificationRecipientEmail_CreateOrUpdate
	fmt.Println("Call operation: NotificationRecipientEmail_CreateOrUpdate")
	notificationRecipientEmailClient, err := armapimanagement.NewNotificationRecipientEmailClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = notificationRecipientEmailClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, testsuite.email, nil)
	testsuite.Require().NoError(err)

	// From step NotificationRecipientEmail_CheckEntityExists
	fmt.Println("Call operation: NotificationRecipientEmail_CheckEntityExists")
	_, err = notificationRecipientEmailClient.CheckEntityExists(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, testsuite.email, nil)
	testsuite.Require().NoError(err)

	// From step NotificationRecipientEmail_ListByNotification
	fmt.Println("Call operation: NotificationRecipientEmail_ListByNotification")
	_, err = notificationRecipientEmailClient.ListByNotification(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, nil)
	testsuite.Require().NoError(err)

	// From step NotificationRecipientEmail_Delete
	fmt.Println("Call operation: NotificationRecipientEmail_Delete")
	_, err = notificationRecipientEmailClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.NotificationNameRequestPublisherNotificationMessage, testsuite.email, nil)
	testsuite.Require().NoError(err)
}
