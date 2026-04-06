// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armfeatures_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armfeatures"
	"github.com/stretchr/testify/suite"
)

type FeaturesClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *FeaturesClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), pathToPackage)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *FeaturesClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestFeaturesClient(t *testing.T) {
	suite.Run(t, new(FeaturesClientTestSuite))
}

func (testsuite *FeaturesClientTestSuite) TestFeaturesCRUD() {
	featureClient, err := armfeatures.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	// list
	fmt.Println("Call operation: Features_List")
	pager := featureClient.NewListPager("Microsoft.Network", nil)
	testsuite.Require().True(pager.More())

	// list all
	fmt.Println("Call operation: Features_ListAll")
	listAll := featureClient.NewListAllPager(nil)
	testsuite.Require().True(listAll.More())

	// list operation
	fmt.Println("Call operation: ListOperations")
	featuresClient, err := armfeatures.NewFeatureClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	listOperations := featuresClient.NewListOperationsPager(nil)
	testsuite.Require().True(listOperations.More())
}
