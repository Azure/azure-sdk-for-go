//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armpostgresqlflexibleservers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type postgresqlflexibleserversOperationsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *postgresqlflexibleserversOperationsTestSuite) SetupSuite() {
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
	testsuite.Prepare()
}

func TestPostgresqlflexibleserversOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(postgresqlflexibleserversOperationsTestSuite))
}

func (testsuite *postgresqlflexibleserversOperationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func (testsuite *postgresqlflexibleserversOperationsTestSuite) TestOperationsNewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = clientFactory.NewOperationsClient().List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *postgresqlflexibleserversOperationsTestSuite) Prepare() {
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	// new client factory

	fmt.Println("subscriptionId", testsuite.subscriptionId, "groupName", testsuite.resourceGroupName, "location", testsuite.location)
	clientFactory, err := armresources.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	client := clientFactory.NewResourceGroupsClient()

	testsuite.Require().NoError(err)
	// check whether create new group successfully
	res, err := client.CheckExistence(testsuite.ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	if !res.Success {
		_, err = client.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, armresources.ResourceGroup{
			Location: to.Ptr(testsuite.location),
		}, nil)
		testsuite.Require().NoError(err)
	}

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}
