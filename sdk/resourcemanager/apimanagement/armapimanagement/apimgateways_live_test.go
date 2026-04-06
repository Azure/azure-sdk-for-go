// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimgatewaysTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	gatewayId         string
	hcId              string
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimgatewaysTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.gatewayId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "gatewayid", 15, false)
	testsuite.hcId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "hcid", 10, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicegateway", 20, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimgatewaysTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimgatewaysTestSuite(t *testing.T) {
	suite.Run(t, new(ApimgatewaysTestSuite))
}

func (testsuite *ApimgatewaysTestSuite) Prepare() {
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
			Name:     to.Ptr(armapimanagement.SKUTypePremium),
			Capacity: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/gateways
func (testsuite *ApimgatewaysTestSuite) TestGateway() {
	var err error
	// From step Gateway_CreateOrUpdate
	fmt.Println("Call operation: Gateway_CreateOrUpdate")
	gatewayClient, err := armapimanagement.NewGatewayClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = gatewayClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, armapimanagement.GatewayContract{
		Properties: &armapimanagement.GatewayContractProperties{
			Description: to.Ptr("my gateway 1"),
			LocationData: &armapimanagement.ResourceLocationDataContract{
				Name: to.Ptr("my location"),
			},
		},
	}, &armapimanagement.GatewayClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Gateway_GetEntityTag
	fmt.Println("Call operation: Gateway_GetEntityTag")
	_, err = gatewayClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_ListByService
	fmt.Println("Call operation: Gateway_ListByService")
	gatewayClientNewListByServicePager := gatewayClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.GatewayClientListByServiceOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for gatewayClientNewListByServicePager.More() {
		_, err := gatewayClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Gateway_Get
	fmt.Println("Call operation: Gateway_Get")
	_, err = gatewayClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_Update
	fmt.Println("Call operation: Gateway_Update")
	_, err = gatewayClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, "*", armapimanagement.GatewayContract{
		Properties: &armapimanagement.GatewayContractProperties{
			Description: to.Ptr("my gateway 1"),
			LocationData: &armapimanagement.ResourceLocationDataContract{
				Name: to.Ptr("my location"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_ListKeys
	fmt.Println("Call operation: Gateway_ListKeys")
	_, err = gatewayClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_RegenerateKey
	fmt.Println("Call operation: Gateway_RegenerateKey")
	_, err = gatewayClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, armapimanagement.GatewayKeyRegenerationRequestContract{
		KeyType: to.Ptr(armapimanagement.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_GenerateToken
	fmt.Println("Call operation: Gateway_GenerateToken")
	_, err = gatewayClient.GenerateToken(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, armapimanagement.GatewayTokenRequestContract{
		Expiry:  to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-04-21T00:44:24.2845269Z"); return t }()),
		KeyType: to.Ptr(armapimanagement.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Gateway_Delete
	fmt.Println("Call operation: Gateway_Delete")
	_, err = gatewayClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.gatewayId, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/outboundNetworkDependenciesEndpoints
func (testsuite *ApimgatewaysTestSuite) TestOutboundnetworkdependenciesendpoints() {
	var err error
	// From step OutboundNetworkDependenciesEndpoints_ListByService
	fmt.Println("Call operation: OutboundNetworkDependenciesEndpoints_ListByService")
	outboundNetworkDependenciesEndpointsClient, err := armapimanagement.NewOutboundNetworkDependenciesEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = outboundNetworkDependenciesEndpointsClient.ListByService(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/networkstatus
func (testsuite *ApimgatewaysTestSuite) TestNetworkstatus() {
	var err error
	// From step NetworkStatus_ListByLocation
	fmt.Println("Call operation: NetworkStatus_ListByLocation")
	networkStatusClient, err := armapimanagement.NewNetworkStatusClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = networkStatusClient.ListByLocation(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// From step NetworkStatus_ListByService
	fmt.Println("Call operation: NetworkStatus_ListByService")
	_, err = networkStatusClient.ListByService(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
}
