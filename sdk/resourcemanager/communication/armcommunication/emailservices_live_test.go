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

type EmailServicesTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	clientFactory     *armcommunication.ClientFactory
	armEndpoint       string
	domainName        string
	emailServiceName  string
	senderUsername    string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *EmailServicesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	var err error
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.domainName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "domainna", 14, true)
	testsuite.domainName = fmt.Sprintf("%s.com", testsuite.domainName)
	testsuite.emailServiceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "emailser", 14, false)
	testsuite.senderUsername, _ = recording.GenerateAlphaNumericID(testsuite.T(), "senderus", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.clientFactory, err = armcommunication.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *EmailServicesTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestEmailServicesTestSuite(t *testing.T) {
	suite.Run(t, new(EmailServicesTestSuite))
}

func (testsuite *EmailServicesTestSuite) Prepare() {
	var err error
	// From step EmailServices_CreateOrUpdate
	fmt.Println("Call operation: EmailServices_CreateOrUpdate")
	emailServicesClient := testsuite.clientFactory.NewEmailServicesClient()
	emailServicesClientCreateOrUpdateResponsePoller, err := emailServicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, armcommunication.EmailServiceResource{
		Location: to.Ptr("Global"),
		Properties: &armcommunication.EmailServiceProperties{
			DataLocation: to.Ptr("United States"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, emailServicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Domains_CreateOrUpdate
	fmt.Println("Call operation: Domains_CreateOrUpdate")
	domainsClient := testsuite.clientFactory.NewDomainsClient()
	domainsClientCreateOrUpdateResponsePoller, err := domainsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, armcommunication.DomainResource{
		Location: to.Ptr("Global"),
		Properties: &armcommunication.DomainProperties{
			DomainManagement: to.Ptr(armcommunication.DomainManagementCustomerManaged),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Communication/emailServices/{emailServiceName}
func (testsuite *EmailServicesTestSuite) TestEmailServices() {
	var err error
	// From step EmailServices_ListBySubscription
	fmt.Println("Call operation: EmailServices_ListBySubscription")
	emailServicesClient := testsuite.clientFactory.NewEmailServicesClient()
	emailServicesClientNewListBySubscriptionPager := emailServicesClient.NewListBySubscriptionPager(nil)
	for emailServicesClientNewListBySubscriptionPager.More() {
		_, err := emailServicesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EmailServices_Get
	fmt.Println("Call operation: EmailServices_Get")
	_, err = emailServicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, nil)
	testsuite.Require().NoError(err)

	// From step EmailServices_ListByResourceGroup
	fmt.Println("Call operation: EmailServices_ListByResourceGroup")
	emailServicesClientNewListByResourceGroupPager := emailServicesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for emailServicesClientNewListByResourceGroupPager.More() {
		_, err := emailServicesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EmailServices_Update
	fmt.Println("Call operation: EmailServices_Update")
	emailServicesClientUpdateResponsePoller, err := emailServicesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, armcommunication.EmailServiceResourceUpdate{
		Tags: map[string]*string{
			"newTag": to.Ptr("newVal"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, emailServicesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Communication/emailServices/{emailServiceName}/domains/{domainName}
func (testsuite *EmailServicesTestSuite) TestDomains() {
	var err error
	// From step Domains_ListByEmailServiceResource
	fmt.Println("Call operation: Domains_ListByEmailServiceResource")
	domainsClient := testsuite.clientFactory.NewDomainsClient()
	domainsClientNewListByEmailServiceResourcePager := domainsClient.NewListByEmailServiceResourcePager(testsuite.resourceGroupName, testsuite.emailServiceName, nil)
	for domainsClientNewListByEmailServiceResourcePager.More() {
		_, err := domainsClientNewListByEmailServiceResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Domains_Get
	fmt.Println("Call operation: Domains_Get")
	_, err = domainsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, nil)
	testsuite.Require().NoError(err)

	// From step Domains_Update
	fmt.Println("Call operation: Domains_Update")
	domainsClientUpdateResponsePoller, err := domainsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, armcommunication.UpdateDomainRequestParameters{
		Properties: &armcommunication.UpdateDomainProperties{
			UserEngagementTracking: to.Ptr(armcommunication.UserEngagementTrackingEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Domains_InitiateVerification
	fmt.Println("Call operation: Domains_InitiateVerification")
	domainsClientInitiateVerificationResponsePoller, err := domainsClient.BeginInitiateVerification(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, armcommunication.VerificationParameter{
		VerificationType: to.Ptr(armcommunication.VerificationTypeSPF),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientInitiateVerificationResponsePoller)
	testsuite.Require().NoError(err)

	// From step Domains_CancelVerification
	fmt.Println("Call operation: Domains_CancelVerification")
	domainsClientCancelVerificationResponsePoller, err := domainsClient.BeginCancelVerification(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, armcommunication.VerificationParameter{
		VerificationType: to.Ptr(armcommunication.VerificationTypeSPF),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientCancelVerificationResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Communication/emailServices/{emailServiceName}/domains/{domainName}/senderUsernames/{senderUsername}
func (testsuite *EmailServicesTestSuite) TestSenderUsernames() {
	var err error
	// From step SenderUsernames_CreateOrUpdate
	fmt.Println("Call operation: SenderUsernames_CreateOrUpdate")
	senderUsernamesClient := testsuite.clientFactory.NewSenderUsernamesClient()
	_, err = senderUsernamesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, testsuite.senderUsername, armcommunication.SenderUsernameResource{
		Properties: &armcommunication.SenderUsernameProperties{
			DisplayName: to.Ptr("Contoso News Alerts"),
			Username:    to.Ptr(testsuite.senderUsername),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SenderUsernames_ListByDomains
	fmt.Println("Call operation: SenderUsernames_ListByDomains")
	senderUsernamesClientNewListByDomainsPager := senderUsernamesClient.NewListByDomainsPager(testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, nil)
	for senderUsernamesClientNewListByDomainsPager.More() {
		_, err := senderUsernamesClientNewListByDomainsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SenderUsernames_Get
	fmt.Println("Call operation: SenderUsernames_Get")
	_, err = senderUsernamesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, testsuite.senderUsername, nil)
	testsuite.Require().NoError(err)

	// From step SenderUsernames_Delete
	fmt.Println("Call operation: SenderUsernames_Delete")
	_, err = senderUsernamesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, testsuite.senderUsername, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *EmailServicesTestSuite) Cleanup() {
	var err error
	// From step Domains_Delete
	fmt.Println("Call operation: Domains_Delete")
	domainsClient := testsuite.clientFactory.NewDomainsClient()
	domainsClientDeleteResponsePoller, err := domainsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, testsuite.domainName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step EmailServices_Delete
	fmt.Println("Call operation: EmailServices_Delete")
	emailServicesClient := testsuite.clientFactory.NewEmailServicesClient()
	emailServicesClientDeleteResponsePoller, err := emailServicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.emailServiceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, emailServicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
