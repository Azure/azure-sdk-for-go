//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerregistry_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ContainerregistryOperationsTestSuite struct {
	suite.Suite

	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	subscriptionId string
}

func (testsuite *ContainerregistryOperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *ContainerregistryOperationsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestContainerregistryOperationTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerregistryOperationsTestSuite))
}

func (testsuite *ContainerregistryOperationsTestSuite) TestContainerregister() {
	clientFactory, err := armcontainerregistry.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err = pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
	}
}
