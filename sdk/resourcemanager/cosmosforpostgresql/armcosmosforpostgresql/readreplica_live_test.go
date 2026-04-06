// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcosmosforpostgresql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmosforpostgresql/armcosmosforpostgresql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ReadReplicaTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	clusterName        string
	clusterReplicaName string
	adminPassword      string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *ReadReplicaTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.clusterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, true)
	testsuite.clusterReplicaName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "clustern", 14, true)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "")
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *ReadReplicaTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestReadReplicaTestSuite(t *testing.T) {
	suite.Run(t, new(ReadReplicaTestSuite))
}

// Microsoft.DBforPostgreSQL/serverGroupsv2/{clusterName}
func (testsuite *ReadReplicaTestSuite) TestClusters() {
	var clusterId string
	var err error
	// From step Clusters_Create
	fmt.Println("Call operation: Clusters_Create")
	clustersClient, err := armcosmosforpostgresql.NewClustersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clustersClientCreateResponsePoller, err := clustersClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterName, armcosmosforpostgresql.Cluster{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcosmosforpostgresql.ClusterProperties{
			AdministratorLoginPassword:      to.Ptr(testsuite.adminPassword),
			CitusVersion:                    to.Ptr("11.1"),
			CoordinatorEnablePublicIPAccess: to.Ptr(true),
			CoordinatorServerEdition:        to.Ptr("GeneralPurpose"),
			CoordinatorStorageQuotaInMb:     to.Ptr[int32](524288),
			CoordinatorVCores:               to.Ptr[int32](4),
			EnableHa:                        to.Ptr(true),
			EnableShardsOnCoordinator:       to.Ptr(false),
			NodeCount:                       to.Ptr[int32](3),
			NodeEnablePublicIPAccess:        to.Ptr(false),
			NodeServerEdition:               to.Ptr("MemoryOptimized"),
			NodeStorageQuotaInMb:            to.Ptr[int32](524288),
			NodeVCores:                      to.Ptr[int32](8),
			PostgresqlVersion:               to.Ptr("15"),
			PreferredPrimaryZone:            to.Ptr("1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var clustersClientCreateResponse *armcosmosforpostgresql.ClustersClientCreateResponse
	clustersClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	clusterId = *clustersClientCreateResponse.ID

	// The source cluster is not ready for restore yet
	if recording.GetRecordMode() == recording.RecordingMode {
		time.Sleep(5 * time.Minute)
	}

	// From step Clusters_CreateReplica
	fmt.Println("Call operation: Clusters_Create")
	clustersClientCreateResponsePoller, err = clustersClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterReplicaName, armcosmosforpostgresql.Cluster{
		Location: to.Ptr(testsuite.location),
		Properties: &armcosmosforpostgresql.ClusterProperties{
			SourceLocation:   to.Ptr(testsuite.location),
			SourceResourceID: to.Ptr(clusterId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Clusters_PromoteReadReplica
	fmt.Println("Call operation: Clusters_PromoteReadReplica")
	clustersClientPromoteReadReplicaResponsePoller, err := clustersClient.BeginPromoteReadReplica(testsuite.ctx, testsuite.resourceGroupName, testsuite.clusterReplicaName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clustersClientPromoteReadReplicaResponsePoller)
	testsuite.Require().NoError(err)
}
