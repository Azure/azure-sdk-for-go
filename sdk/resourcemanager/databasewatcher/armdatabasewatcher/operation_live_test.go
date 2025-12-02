// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armdatabasewatcher_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databasewatcher/armdatabasewatcher"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type DatabaseWatcherTestSuite struct {
	suite.Suite
	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	location       string
	subscriptionId string
}

func (testsuite *DatabaseWatcherTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func TestDatabaseWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseWatcherTestSuite))
}

func (testsuite *DatabaseWatcherTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func (testsuite *DatabaseWatcherTestSuite) TestOperationGet() {
	clientFactory, err := armdatabasewatcher.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
	}
}
