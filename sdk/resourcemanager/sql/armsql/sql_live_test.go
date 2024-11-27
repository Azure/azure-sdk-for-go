//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armsql_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
	DatabaseName     = "test01"
	PathToPackage    = "sdk/resourcemanager/sql/armsql/testdata"
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
	testutil.StartRecording(testsuite.T(), PathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "sqlscenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.sqlServersName = "sql-test-" + strconv.Itoa(rand.Intn(1000))
	clientFactory, err := armresources.NewClientFactory(testsuite.subscriptionId, testsuite.cred, nil)
	client := clientFactory.NewResourceGroupsClient()
	ctx := context.Background()

	testsuite.resourceGroupName = "sqltestgroup-" + strconv.Itoa(rand.Intn(1000))
	_, err = client.CreateOrUpdate(ctx, testsuite.resourceGroupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	testsuite.Require().NoError(err)

	fmt.Println("create new resource group ", testsuite.resourceGroupName, " of ", testsuite.subscriptionId, "successfully")
}

func (testsuite *SqlAccessTestSuite) TestCreateServer() {
	fmt.Println("start to create sql server~")
	// create new msql resource under group
	sqlClientFactory, err := armsql.NewClientFactory(testsuite.subscriptionId, testsuite.cred, nil)
	testsuite.Require().NoError(err)
	serverclient := sqlClientFactory.NewServersClient()

	ctx := context.Background()
	server, err := testsuite.createServer(testsuite.ctx, serverclient)
	testsuite.Require().NoError(err)
	fmt.Println("create serverId", *server.ID)

	server, err = testsuite.getServer(ctx, serverclient)
	testsuite.Require().NoError(err)
	fmt.Println("get server:", *server.ID)

	// create database
	databasesClient := sqlClientFactory.NewDatabasesClient()
	testsuite.Require().NotNil(databasesClient)
	database, err := testsuite.createDatabase(ctx, databasesClient)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(database)

	fmt.Println("database:", *database.ID)
}

func TestSqlAccessTestSuite(t *testing.T) {
	suite.Run(t, new(SqlAccessTestSuite))
}

func (testsuite *SqlAccessTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func (testsuite *SqlAccessTestSuite) createServer(ctx context.Context, serversClient *armsql.ServersClient) (*armsql.Server, error) {

	pollerResp, err := serversClient.BeginCreateOrUpdate(
		ctx,
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
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func (testsuite *SqlAccessTestSuite) getServer(ctx context.Context, serversClient *armsql.ServersClient) (*armsql.Server, error) {

	resp, err := serversClient.Get(ctx, testsuite.resourceGroupName, testsuite.sqlServersName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func (testsuite *SqlAccessTestSuite) createDatabase(ctx context.Context, databasesClient *armsql.DatabasesClient) (*armsql.Database, error) {

	pollerResp, err := databasesClient.BeginCreateOrUpdate(
		ctx,
		testsuite.resourceGroupName,
		testsuite.sqlServersName,
		DatabaseName,
		armsql.Database{
			Location: to.Ptr(ResourceLocation),
			Properties: &armsql.DatabaseProperties{
				ReadScale: to.Ptr(armsql.DatabaseReadScaleDisabled),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Database, nil
}
