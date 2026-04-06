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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimportalsettingsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimportalsettingsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicesetting", 20, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimportalsettingsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimportalsettingsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimportalsettingsTestSuite))
}

func (testsuite *ApimportalsettingsTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail: to.Ptr("foo@contoso.com"),
			PublisherName:  to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:     to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/portalsettings
func (testsuite *ApimportalsettingsTestSuite) TestPortalsettings() {
	var err error
	// From step PortalSettings_ListByService
	fmt.Println("Call operation: PortalSettings_ListByService")
	portalSettingsClient, err := armapimanagement.NewPortalSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = portalSettingsClient.ListByService(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/portalsettings/signin
func (testsuite *ApimportalsettingsTestSuite) TestSigninsettings() {
	var err error
	// From step SignInSettings_CreateOrUpdate
	fmt.Println("Call operation: SignInSettings_CreateOrUpdate")
	signInSettingsClient, err := armapimanagement.NewSignInSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = signInSettingsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PortalSigninSettings{
		Properties: &armapimanagement.PortalSigninSettingProperties{
			Enabled: to.Ptr(true),
		},
	}, &armapimanagement.SignInSettingsClientCreateOrUpdateOptions{IfMatch: to.Ptr("*")})
	testsuite.Require().NoError(err)

	// From step SignInSettings_GetEntityTag
	fmt.Println("Call operation: SignInSettings_GetEntityTag")
	_, err = signInSettingsClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step SignInSettings_Get
	fmt.Println("Call operation: SignInSettings_Get")
	_, err = signInSettingsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step SignInSettings_Update
	fmt.Println("Call operation: SignInSettings_Update")
	_, err = signInSettingsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "*", armapimanagement.PortalSigninSettings{
		Properties: &armapimanagement.PortalSigninSettingProperties{
			Enabled: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/portalsettings/signup
func (testsuite *ApimportalsettingsTestSuite) TestSignupsettings() {
	var err error
	// From step SignUpSettings_CreateOrUpdate
	fmt.Println("Call operation: SignUpSettings_CreateOrUpdate")
	signUpSettingsClient, err := armapimanagement.NewSignUpSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = signUpSettingsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.PortalSignupSettings{
		Properties: &armapimanagement.PortalSignupSettingsProperties{
			Enabled: to.Ptr(true),
			TermsOfService: &armapimanagement.TermsOfServiceProperties{
				ConsentRequired: to.Ptr(true),
				Enabled:         to.Ptr(true),
				Text:            to.Ptr("Terms of service text."),
			},
		},
	}, &armapimanagement.SignUpSettingsClientCreateOrUpdateOptions{IfMatch: to.Ptr("*")})
	testsuite.Require().NoError(err)

	// From step SignUpSettings_GetEntityTag
	fmt.Println("Call operation: SignUpSettings_GetEntityTag")
	_, err = signUpSettingsClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step SignUpSettings_Get
	fmt.Println("Call operation: SignUpSettings_Get")
	_, err = signUpSettingsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step SignUpSettings_Update
	fmt.Println("Call operation: SignUpSettings_Update")
	_, err = signUpSettingsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "*", armapimanagement.PortalSignupSettings{
		Properties: &armapimanagement.PortalSignupSettingsProperties{
			Enabled: to.Ptr(true),
			TermsOfService: &armapimanagement.TermsOfServiceProperties{
				ConsentRequired: to.Ptr(true),
				Enabled:         to.Ptr(true),
				Text:            to.Ptr("Terms of service text."),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
}
