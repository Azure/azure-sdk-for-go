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

type ApimopenidconnectprovidersTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	opid              string
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimopenidconnectprovidersTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.opid, _ = recording.GenerateAlphaNumericID(testsuite.T(), "opid", 10, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceopenid", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimopenidconnectprovidersTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimopenidconnectprovidersTestSuite(t *testing.T) {
	suite.Run(t, new(ApimopenidconnectprovidersTestSuite))
}

func (testsuite *ApimopenidconnectprovidersTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/openidConnectProviders
func (testsuite *ApimopenidconnectprovidersTestSuite) TestOpenidconnectprovider() {
	var err error
	// From step OpenIdConnectProvider_CreateOrUpdate
	fmt.Println("Call operation: OpenIdConnectProvider_CreateOrUpdate")
	openIDConnectProviderClient, err := armapimanagement.NewOpenIDConnectProviderClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = openIDConnectProviderClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, armapimanagement.OpenidConnectProviderContract{
		Properties: &armapimanagement.OpenidConnectProviderContractProperties{
			ClientID:         to.Ptr("oidprovidertemplate3"),
			DisplayName:      to.Ptr("templateoidprovider3"),
			MetadataEndpoint: to.Ptr("https://oidprovider-template3.net"),
			ClientSecret:     to.Ptr("x"),
		},
	}, &armapimanagement.OpenIDConnectProviderClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step OpenIdConnectProvider_GetEntityTag
	fmt.Println("Call operation: OpenIdConnectProvider_GetEntityTag")
	_, err = openIDConnectProviderClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, nil)
	testsuite.Require().NoError(err)

	// From step OpenIdConnectProvider_ListByService
	fmt.Println("Call operation: OpenIdConnectProvider_ListByService")
	openIDConnectProviderClientNewListByServicePager := openIDConnectProviderClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.OpenIDConnectProviderClientListByServiceOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for openIDConnectProviderClientNewListByServicePager.More() {
		_, err := openIDConnectProviderClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step OpenIdConnectProvider_Get
	fmt.Println("Call operation: OpenIdConnectProvider_Get")
	_, err = openIDConnectProviderClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, nil)
	testsuite.Require().NoError(err)

	// From step OpenIdConnectProvider_Update
	fmt.Println("Call operation: OpenIdConnectProvider_Update")
	_, err = openIDConnectProviderClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, "*", armapimanagement.OpenidConnectProviderUpdateContract{
		Properties: &armapimanagement.OpenidConnectProviderUpdateContractProperties{
			ClientSecret: to.Ptr("updatedsecret"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step OpenIdConnectProvider_ListSecrets
	fmt.Println("Call operation: OpenIdConnectProvider_ListSecrets")
	_, err = openIDConnectProviderClient.ListSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, nil)
	testsuite.Require().NoError(err)

	// From step OpenIdConnectProvider_Delete
	fmt.Println("Call operation: OpenIdConnectProvider_Delete")
	_, err = openIDConnectProviderClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.opid, "*", nil)
	testsuite.Require().NoError(err)
}
