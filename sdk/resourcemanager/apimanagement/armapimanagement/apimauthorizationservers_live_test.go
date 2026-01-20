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

type ApimauthorizationserversTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	authsid			string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimauthorizationserversTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.authsid, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authsid", 13, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceauth", 17, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimauthorizationserversTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimauthorizationserversTestSuite(t *testing.T) {
	suite.Run(t, new(ApimauthorizationserversTestSuite))
}

func (testsuite *ApimauthorizationserversTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/authorizationServers
func (testsuite *ApimauthorizationserversTestSuite) TestAuthorizationserver() {
	var err error
	// From step AuthorizationServer_CreateOrUpdate
	fmt.Println("Call operation: AuthorizationServer_CreateOrUpdate")
	authorizationServerClient, err := armapimanagement.NewAuthorizationServerClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = authorizationServerClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, armapimanagement.AuthorizationServerContract{
		Properties: &armapimanagement.AuthorizationServerContractProperties{
			Description:	to.Ptr("test server"),
			AuthorizationMethods: []*armapimanagement.AuthorizationMethod{
				to.Ptr(armapimanagement.AuthorizationMethodGET)},
			BearerTokenSendingMethods: []*armapimanagement.BearerTokenSendingMethod{
				to.Ptr(armapimanagement.BearerTokenSendingMethodAuthorizationHeader)},
			DefaultScope:			to.Ptr("read write"),
			ResourceOwnerPassword:		to.Ptr("pwd"),
			ResourceOwnerUsername:		to.Ptr("un"),
			SupportState:			to.Ptr(true),
			TokenEndpoint:			to.Ptr("https://www.contoso.com/oauth2/token"),
			AuthorizationEndpoint:		to.Ptr("https://www.contoso.com/oauth2/auth"),
			ClientID:			to.Ptr("1"),
			ClientRegistrationEndpoint:	to.Ptr("https://www.contoso.com/apps"),
			ClientSecret:			to.Ptr("2"),
			DisplayName:			to.Ptr("test2"),
			GrantTypes: []*armapimanagement.GrantType{
				to.Ptr(armapimanagement.GrantTypeAuthorizationCode),
				to.Ptr(armapimanagement.GrantTypeImplicit)},
		},
	}, &armapimanagement.AuthorizationServerClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step AuthorizationServer_GetEntityTag
	fmt.Println("Call operation: AuthorizationServer_GetEntityTag")
	_, err = authorizationServerClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, nil)
	testsuite.Require().NoError(err)

	// From step AuthorizationServer_ListByService
	authorizationServerClientNewListByServicePager := authorizationServerClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.AuthorizationServerClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for authorizationServerClientNewListByServicePager.More() {
		_, err := authorizationServerClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AuthorizationServer_Get
	fmt.Println("Call operation: AuthorizationServer_Get")
	_, err = authorizationServerClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, nil)
	testsuite.Require().NoError(err)

	// From step AuthorizationServer_Update
	fmt.Println("Call operation: AuthorizationServer_Update")
	_, err = authorizationServerClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, "*", armapimanagement.AuthorizationServerUpdateContract{
		Properties: &armapimanagement.AuthorizationServerUpdateContractProperties{
			ClientID:	to.Ptr("update"),
			ClientSecret:	to.Ptr("updated"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AuthorizationServer_ListSecrets
	fmt.Println("Call operation: AuthorizationServer_ListSecrets")
	_, err = authorizationServerClient.ListSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, nil)
	testsuite.Require().NoError(err)

	// From step AuthorizationServer_Delete
	fmt.Println("Call operation: AuthorizationServer_Delete")
	_, err = authorizationServerClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.authsid, "*", nil)
	testsuite.Require().NoError(err)
}
