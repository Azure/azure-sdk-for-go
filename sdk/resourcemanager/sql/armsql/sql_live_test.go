//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armsql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type SqlAccessTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	sqlServersName    string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SqlAccessTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "sqlscenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.sqlServersName = "sqlserverasdsfsdfds-test"

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}

func TestSqlAccessTestSuite(t *testing.T) {
	suite.Run(t, new(SqlAccessTestSuite))
}

func (testsuite *SqlAccessTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func (testsuite *SqlAccessTestSuite) TestCreateServer() {
	sqlClientFactory, err := armsql.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClient := sqlClientFactory.NewServersClient()
	pollerResp, err := serversClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.sqlServersName,
		armsql.Server{
			Location: to.Ptr(ResourceLocation),
			Properties: &armsql.ServerProperties{
				AdministratorLogin:         to.Ptr("dummylogin"),
				AdministratorLoginPassword: to.Ptr("QWE123!@#"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
}

func (testsuite *SqlAccessTestSuite) TestGetServer() {
	sqlClientFactory, err := armsql.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClient := sqlClientFactory.NewServersClient()
	resp, err := serversClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.sqlServersName, nil)
	testsuite.Require().NoError(err)
	fmt.Println("get server:", *resp.Server.ID)
}

func (testsuite *SqlAccessTestSuite) Cleanup() {
	var err error
	// From step CapacityReservationGroups_Delete
	fmt.Println("Call operation: CapacityReservationGroups_Delete")
	sqlClientFactory, err := armsql.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pollerResp, err := sqlClientFactory.NewServersClient().BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.sqlServersName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
}
