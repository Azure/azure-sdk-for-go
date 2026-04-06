// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armchangeanalysis_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/changeanalysis/armchangeanalysis"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ChangeanalysisTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string

	startTime time.Time
	endTime   time.Time
}

func (testsuite *ChangeanalysisTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

	testsuite.startTime, _ = time.Parse(time.RFC3339Nano, "2023-11-25T02:17:41.390Z")
	testsuite.endTime, _ = time.Parse(time.RFC3339Nano, "2023-12-01T02:17:41.390Z")
}

func (testsuite *ChangeanalysisTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestChangeanalysisTestSuite(t *testing.T) {
	suite.Run(t, new(ChangeanalysisTestSuite))
}

// Microsoft.ChangeAnalysis/operations
func (testsuite *ChangeanalysisTestSuite) TestOperations() {
	var err error
	opt := testsuite.options.Clone()
	opt.APIVersion = "2020-04-01-preview"

	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armchangeanalysis.NewOperationsClient(testsuite.cred, opt)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(&armchangeanalysis.OperationsClientListOptions{SkipToken: nil})
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ChangeAnalysis/resourceChanges
func (testsuite *ChangeanalysisTestSuite) TestResourceChanges() {
	var err error
	// From step ResourceChanges_List
	fmt.Println("Call operation: ResourceChanges_List")
	resourceChangesClient, err := armchangeanalysis.NewResourceChangesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceChangesClientNewListPager := resourceChangesClient.NewListPager("subscriptions/"+testsuite.subscriptionId+"/resourceGroups/"+testsuite.resourceGroupName, testsuite.startTime, testsuite.endTime, &armchangeanalysis.ResourceChangesClientListOptions{SkipToken: nil})
	for resourceChangesClientNewListPager.More() {
		_, err := resourceChangesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ChangeAnalysis/changes
func (testsuite *ChangeanalysisTestSuite) TestChanges() {
	var err error
	// From step Changes_ListChangesByResourceGroup
	fmt.Println("Call operation: Changes_ListChangesByResourceGroup")
	changesClient, err := armchangeanalysis.NewChangesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	changesClientNewListChangesByResourceGroupPager := changesClient.NewListChangesByResourceGroupPager(testsuite.resourceGroupName, testsuite.startTime, testsuite.endTime, &armchangeanalysis.ChangesClientListChangesByResourceGroupOptions{SkipToken: nil})
	for changesClientNewListChangesByResourceGroupPager.More() {
		_, err := changesClientNewListChangesByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Changes_ListChangesBySubscription
	fmt.Println("Call operation: Changes_ListChangesBySubscription")
	changesClientNewListChangesBySubscriptionPager := changesClient.NewListChangesBySubscriptionPager(testsuite.startTime, testsuite.endTime, &armchangeanalysis.ChangesClientListChangesBySubscriptionOptions{SkipToken: nil})
	for changesClientNewListChangesBySubscriptionPager.More() {
		_, err := changesClientNewListChangesBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
