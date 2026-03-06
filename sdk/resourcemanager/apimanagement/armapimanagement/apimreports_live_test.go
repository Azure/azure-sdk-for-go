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

type ApimreportsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimreportsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicereports", 20, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimreportsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimreportsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimreportsTestSuite))
}

func (testsuite *ApimreportsTestSuite) Prepare() {
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

// Microsoft.ApiManagement/service/reports
func (testsuite *ApimreportsTestSuite) TestReports() {
	var err error
	// From step Reports_ListByApi
	fmt.Println("Call operation: Reports_ListByApi")
	reportsClient, err := armapimanagement.NewReportsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	reportsClientNewListByAPIPager := reportsClient.NewListByAPIPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByAPIOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListByAPIPager.More() {
		_, err := reportsClientNewListByAPIPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListByGeo
	fmt.Println("Call operation: Reports_ListByGeo")
	reportsClientNewListByGeoPager := reportsClient.NewListByGeoPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByGeoOptions{Top: nil,
		Skip:	nil,
	})
	for reportsClientNewListByGeoPager.More() {
		_, err := reportsClientNewListByGeoPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListByOperation
	fmt.Println("Call operation: Reports_ListByOperation")
	reportsClientNewListByOperationPager := reportsClient.NewListByOperationPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByOperationOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListByOperationPager.More() {
		_, err := reportsClientNewListByOperationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListByProduct
	fmt.Println("Call operation: Reports_ListByProduct")
	reportsClientNewListByProductPager := reportsClient.NewListByProductPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByProductOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListByProductPager.More() {
		_, err := reportsClientNewListByProductPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListBySubscription
	fmt.Println("Call operation: Reports_ListBySubscription")
	reportsClientNewListBySubscriptionPager := reportsClient.NewListBySubscriptionPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListBySubscriptionOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListBySubscriptionPager.More() {
		_, err := reportsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListByTime
	fmt.Println("Call operation: Reports_ListByTime")
	reportsClientNewListByTimePager := reportsClient.NewListByTimePager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", "PT15M", &armapimanagement.ReportsClientListByTimeOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListByTimePager.More() {
		_, err := reportsClientNewListByTimePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Reports_ListByUser
	fmt.Println("Call operation: Reports_ListByUser")
	reportsClientNewListByUserPager := reportsClient.NewListByUserPager(testsuite.resourceGroupName, testsuite.serviceName, "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByUserOptions{Top: nil,
		Skip:		nil,
		Orderby:	nil,
	})
	for reportsClientNewListByUserPager.More() {
		_, err := reportsClientNewListByUserPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ApiManagement/service/regions
func (testsuite *ApimreportsTestSuite) TestRegion() {
	var err error
	// From step Region_ListByService
	fmt.Println("Call operation: Region_ListByService")
	regionClient, err := armapimanagement.NewRegionClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	regionClientNewListByServicePager := regionClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for regionClientNewListByServicePager.More() {
		_, err := regionClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
