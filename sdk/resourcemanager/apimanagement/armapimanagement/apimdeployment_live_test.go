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

type ApimdeploymentTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimdeploymentTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicedeploy", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ApimdeploymentTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimdeploymentTestSuite(t *testing.T) {
	suite.Run(t, new(ApimdeploymentTestSuite))
}

// Microsoft.ApiManagement/service
func (testsuite *ApimdeploymentTestSuite) TestApimanagementservice() {
	var err error
	// From step ApiManagementService_CheckNameAvailability
	fmt.Println("Call operation: ApiManagementService_CheckNameAvailability")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = serviceClient.CheckNameAvailability(testsuite.ctx, armapimanagement.ServiceCheckNameAvailabilityParameters{
		Name: to.Ptr("apimService1"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
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

	// From step ApiManagementService_List
	fmt.Println("Call operation: ApiManagementService_List")
	serviceClientNewListPager := serviceClient.NewListPager(nil)
	for serviceClientNewListPager.More() {
		_, err := serviceClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiManagementService_ListByResourceGroup
	fmt.Println("Call operation: ApiManagementService_ListByResourceGroup")
	serviceClientNewListByResourceGroupPager := serviceClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for serviceClientNewListByResourceGroupPager.More() {
		_, err := serviceClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiManagementService_Get
	fmt.Println("Call operation: ApiManagementService_Get")
	_, err = serviceClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step ApiManagementService_Update
	fmt.Println("Call operation: ApiManagementService_Update")
	serviceClientUpdateResponsePoller, err := serviceClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceUpdateParameters{
		Properties: &armapimanagement.ServiceUpdateProperties{
			CustomProperties: map[string]*string{
				"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10": to.Ptr("false"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApiManagementService_GetDomainOwnershipIdentifier
	fmt.Println("Call operation: ApiManagementService_GetDomainOwnershipIdentifier")
	_, err = serviceClient.GetDomainOwnershipIdentifier(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step ApiManagementService_GetSsoToken
	fmt.Println("Call operation: ApiManagementService_GetSsoToken")
	_, err = serviceClient.GetSsoToken(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step ApiManagementServiceSkus_ListAvailableServiceSkus
	fmt.Println("Call operation: ApiManagementServiceSkus_ListAvailableServiceSkus")
	serviceSKUsClient, err := armapimanagement.NewServiceSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceSKUsClientNewListAvailableServiceSKUsPager := serviceSKUsClient.NewListAvailableServiceSKUsPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for serviceSKUsClientNewListAvailableServiceSKUsPager.More() {
		_, err := serviceSKUsClientNewListAvailableServiceSKUsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiManagementService_Delete
	fmt.Println("Call operation: ApiManagementService_Delete")
	serviceClientDeleteResponsePoller, err := serviceClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step DeletedServices_ListBySubscription
	fmt.Println("Call operation: DeletedServices_ListBySubscription")
	deletedServicesClient, err := armapimanagement.NewDeletedServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	deletedServicesClientNewListBySubscriptionPager := deletedServicesClient.NewListBySubscriptionPager(nil)
	for deletedServicesClientNewListBySubscriptionPager.More() {
		_, err := deletedServicesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DeletedServices_GetByName
	fmt.Println("Call operation: DeletedServices_GetByName")
	_, err = deletedServicesClient.GetByName(testsuite.ctx, testsuite.serviceName, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// From step DeletedServices_Purge
	fmt.Println("Call operation: DeletedServices_Purge")
	deletedServicesClientPurgeResponsePoller, err := deletedServicesClient.BeginPurge(testsuite.ctx, testsuite.serviceName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, deletedServicesClientPurgeResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/operations
func (testsuite *ApimdeploymentTestSuite) TestApimanagementoperations() {
	var err error
	// From step ApiManagementOperations_List
	fmt.Println("Call operation: ApiManagementOperations_List")
	operationsClient, err := armapimanagement.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
