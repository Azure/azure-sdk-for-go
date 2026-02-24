// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armpostgresql_test

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql/v2"
	"github.com/stretchr/testify/suite"
)

type PostgresqlOperationsTestSuite struct {
	suite.Suite

	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	subscriptionId string
}

func (testsuite *PostgresqlOperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *PostgresqlOperationsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestPostgresqlOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresqlOperationsTestSuite))
}

// Microsoft.DBforPostgreSQL/servers/{serverName}
func (testsuite *PostgresqlOperationsTestSuite) TestNewListPager() {
	clientFactory, err := armpostgresql.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
	}
}
