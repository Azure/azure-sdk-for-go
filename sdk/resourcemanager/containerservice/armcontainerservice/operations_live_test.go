//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerservice_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type OperationTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	azureClientId     string
	azureClientSecret string
	subscriptionId    string
}

func (testsuite *OperationTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.azureClientId = recording.GetEnvVariable("AZURE_CLIENT_ID", "000000000000")
	testsuite.azureClientSecret = recording.GetEnvVariable("AZURE_CLIENT_SECRET", "000000000000")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	testutil.StartRecording(testsuite.T(), pathToPackage)

}

func (testsuite *OperationTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestOperationTestSuite(t *testing.T) {
	suite.Run(t, new(OperationTestSuite))
}

func (testsuite *OperationTestSuite) TestOperationsNewListPager() {
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		testsuite.Require().NoError(err)
	}
}
