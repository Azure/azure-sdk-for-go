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

type ApimidentityproviderTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	azureClientId     string
	azureClientSecret string
	azureTenantId     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimidentityproviderTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceidentity", 21, false)
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.azureClientSecret = recording.GetEnvVariable("AZURE_CLIENT_SECRET", "000000000000")
	testsuite.azureTenantId = recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimidentityproviderTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimidentityproviderTestSuite(t *testing.T) {
	suite.Run(t, new(ApimidentityproviderTestSuite))
}

func (testsuite *ApimidentityproviderTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/identityProviders
func (testsuite *ApimidentityproviderTestSuite) TestIdentityprovider() {
	var err error
	// From step IdentityProvider_CreateOrUpdate
	fmt.Println("Call operation: IdentityProvider_CreateOrUpdate")
	identityProviderClient, err := armapimanagement.NewIdentityProviderClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = identityProviderClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, armapimanagement.IdentityProviderCreateContract{
		Properties: &armapimanagement.IdentityProviderCreateContractProperties{
			AllowedTenants: []*string{
				to.Ptr(testsuite.azureTenantId)},
			ClientID:     to.Ptr(testsuite.azureClientId),
			ClientSecret: to.Ptr(testsuite.azureClientSecret),
		},
	}, &armapimanagement.IdentityProviderClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step IdentityProvider_GetEntityTag
	fmt.Println("Call operation: IdentityProvider_GetEntityTag")
	_, err = identityProviderClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, nil)
	testsuite.Require().NoError(err)

	// From step IdentityProvider_ListByService
	fmt.Println("Call operation: IdentityProvider_ListByService")
	identityProviderClientNewListByServicePager := identityProviderClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for identityProviderClientNewListByServicePager.More() {
		_, err := identityProviderClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IdentityProvider_Get
	fmt.Println("Call operation: IdentityProvider_Get")
	_, err = identityProviderClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, nil)
	testsuite.Require().NoError(err)

	// From step IdentityProvider_Update
	fmt.Println("Call operation: IdentityProvider_Update")
	_, err = identityProviderClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, "*", armapimanagement.IdentityProviderUpdateParameters{
		Properties: &armapimanagement.IdentityProviderUpdateProperties{
			ClientID:     to.Ptr(testsuite.azureClientId),
			ClientSecret: to.Ptr(testsuite.azureClientSecret),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IdentityProvider_ListSecrets
	fmt.Println("Call operation: IdentityProvider_ListSecrets")
	_, err = identityProviderClient.ListSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, nil)
	testsuite.Require().NoError(err)

	// From step IdentityProvider_Delete
	fmt.Println("Call operation: IdentityProvider_Delete")
	_, err = identityProviderClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.IdentityProviderTypeAAD, "*", nil)
	testsuite.Require().NoError(err)
}
