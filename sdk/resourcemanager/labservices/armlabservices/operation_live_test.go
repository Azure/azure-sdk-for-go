// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlabservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
	"github.com/stretchr/testify/suite"
)

type OperationTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *OperationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *OperationTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestOperationTestSuite(t *testing.T) {
	suite.Run(t, new(OperationTestSuite))
}

// Microsoft.LabServices/operations
func (testsuite *OperationTestSuite) TestOperation() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armlabservices.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.LabServices/skus
func (testsuite *OperationTestSuite) TestSkus() {
	var err error
	// From step Skus_List
	fmt.Println("Call operation: Skus_List")
	sKUsClient, err := armlabservices.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(&armlabservices.SKUsClientListOptions{Filter: nil})
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.LabServices/locations/{location}/usages
func (testsuite *OperationTestSuite) TestUsages() {
	var err error
	// From step Usages_ListByLocation
	fmt.Println("Call operation: Usages_ListByLocation")
	usagesClient, err := armlabservices.NewUsagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	usagesClientNewListByLocationPager := usagesClient.NewListByLocationPager(testsuite.location, &armlabservices.UsagesClientListByLocationOptions{Filter: nil})
	for usagesClientNewListByLocationPager.More() {
		_, err := usagesClientNewListByLocationPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
