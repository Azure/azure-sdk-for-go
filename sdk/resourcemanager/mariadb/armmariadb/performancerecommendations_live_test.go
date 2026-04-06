// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mariadb/armmariadb"
	"github.com/stretchr/testify/suite"
)

type PerformanceRecommendationsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	serverName        string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *PerformanceRecommendationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.serverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serverna", 14, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *PerformanceRecommendationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPerformanceRecommendationsTestSuite(t *testing.T) {
	suite.Run(t, new(PerformanceRecommendationsTestSuite))
}

func (testsuite *PerformanceRecommendationsTestSuite) Prepare() {
	var err error
	// From step Servers_Create
	fmt.Println("Call operation: Servers_Create")
	serversClient, err := armmariadb.NewServersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serversClientCreateResponsePoller, err := serversClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.ServerForCreate{
		Location: to.Ptr(testsuite.location),
		Properties: &armmariadb.ServerPropertiesForDefaultCreate{
			CreateMode:        to.Ptr(armmariadb.CreateModeDefault),
			MinimalTLSVersion: to.Ptr(armmariadb.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:    to.Ptr(armmariadb.SSLEnforcementEnumEnabled),
			StorageProfile: &armmariadb.StorageProfile{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armmariadb.GeoRedundantBackupEnabled),
				StorageMB:           to.Ptr[int32](128000),
			},
			AdministratorLogin:         to.Ptr("cloudsa"),
			AdministratorLoginPassword: to.Ptr(testsuite.adminPassword),
		},
		SKU: &armmariadb.SKU{
			Name:     to.Ptr("GP_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armmariadb.SKUTierGeneralPurpose),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serversClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/advisors/{advisorName}
func (testsuite *PerformanceRecommendationsTestSuite) TestAdvisors() {
	var err error
	// From step Advisors_ListByServer
	fmt.Println("Call operation: Advisors_ListByServer")
	advisorsClient, err := armmariadb.NewAdvisorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	advisorsClientNewListByServerPager := advisorsClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for advisorsClientNewListByServerPager.More() {
		_, err := advisorsClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Advisors_Get
	fmt.Println("Call operation: Advisors_Get")
	_, err = advisorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "Index", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/advisors/{advisorName}/createRecommendedActionSession
func (testsuite *PerformanceRecommendationsTestSuite) TestCreateRecommendedActionSession() {
	var err error
	// From step CreateRecommendedActionSession
	fmt.Println("Call operation: CreateRecommendedActionSession")
	managementClient, err := armmariadb.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managementClientCreateRecommendedActionSessionResponsePoller, err := managementClient.BeginCreateRecommendedActionSession(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, "Index", "someDatabaseName", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managementClientCreateRecommendedActionSessionResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DBforMariaDB/servers/{serverName}/advisors/{advisorName}/recommendedActions/{recommendedActionName}
func (testsuite *PerformanceRecommendationsTestSuite) TestRecommendedActions() {
	var err error
	// From step RecommendedActions_ListByServer
	fmt.Println("Call operation: RecommendedActions_ListByServer")
	recommendedActionsClient, err := armmariadb.NewRecommendedActionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	recommendedActionsClientNewListByServerPager := recommendedActionsClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, "Index", &armmariadb.RecommendedActionsClientListByServerOptions{SessionID: nil})
	for recommendedActionsClientNewListByServerPager.More() {
		_, err := recommendedActionsClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
