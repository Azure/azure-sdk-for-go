// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armcomputeschedule_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computeschedule/armcomputeschedule"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type ComputeScheduleOperationsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ComputeScheduleOperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "operationTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}

func TTestComputeScheduleOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(ComputeScheduleOperationsTestSuite))
}

func (testsuite *ComputeScheduleOperationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func (testsuite *ComputeScheduleOperationsTestSuite) TestOperationsNewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	ctx := context.Background()
	clientFactory, err := armcomputeschedule.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		testsuite.Require().NoError(err)
	}
}
