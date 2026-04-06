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

type ServerSecurityAlertPoliciesTestSuite struct {
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

func (testsuite *ServerSecurityAlertPoliciesTestSuite) SetupSuite() {
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

func (testsuite *ServerSecurityAlertPoliciesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestServerSecurityAlertPoliciesTestSuite(t *testing.T) {
	suite.Run(t, new(ServerSecurityAlertPoliciesTestSuite))
}

func (testsuite *ServerSecurityAlertPoliciesTestSuite) Prepare() {
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

// Microsoft.DBforMariaDB/servers/{serverName}/securityAlertPolicies/{securityAlertPolicyName}
func (testsuite *ServerSecurityAlertPoliciesTestSuite) TestServerSecurityAlertPolicies() {
	var err error
	// From step ServerSecurityAlertPolicies_CreateOrUpdate
	fmt.Println("Call operation: ServerSecurityAlertPolicies_CreateOrUpdate")
	serverSecurityAlertPoliciesClient, err := armmariadb.NewServerSecurityAlertPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serverSecurityAlertPoliciesClientCreateOrUpdateResponsePoller, err := serverSecurityAlertPoliciesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.SecurityAlertPolicyNameDefault, armmariadb.ServerSecurityAlertPolicy{
		Properties: &armmariadb.SecurityAlertPolicyProperties{
			State: to.Ptr(armmariadb.ServerSecurityAlertPolicyStateEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serverSecurityAlertPoliciesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ServerSecurityAlertPolicies_ListByServer
	fmt.Println("Call operation: ServerSecurityAlertPolicies_ListByServer")
	serverSecurityAlertPoliciesClientNewListByServerPager := serverSecurityAlertPoliciesClient.NewListByServerPager(testsuite.resourceGroupName, testsuite.serverName, nil)
	for serverSecurityAlertPoliciesClientNewListByServerPager.More() {
		_, err := serverSecurityAlertPoliciesClientNewListByServerPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServerSecurityAlertPolicies_Get
	fmt.Println("Call operation: ServerSecurityAlertPolicies_Get")
	_, err = serverSecurityAlertPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serverName, armmariadb.SecurityAlertPolicyNameDefault, nil)
	testsuite.Require().NoError(err)
}
