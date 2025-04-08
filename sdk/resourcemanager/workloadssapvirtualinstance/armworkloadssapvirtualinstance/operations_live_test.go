//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package armworkloadssapvirtualinstance_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloadssapvirtualinstance/armworkloadssapvirtualinstance"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type WorkloadssapvirtualinstanceOperationsTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	subscriptionId    string
}

func (testsuite *WorkloadssapvirtualinstanceOperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func TestWorkloadssapvirtualinstanceOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(WorkloadssapvirtualinstanceOperationsTestSuite))
}

func (testsuite *WorkloadssapvirtualinstanceOperationsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func (testsuite *WorkloadssapvirtualinstanceOperationsTestSuite) TestOperationsNewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	testsuite.Require().NoError(err)
	clientFactory, err := armworkloadssapvirtualinstance.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
	}
}
