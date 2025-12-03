// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armredisenterprise_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v3"
	"github.com/stretchr/testify/suite"
)

type RedisenterpriseOperationsTestSuite struct {
	suite.Suite

	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	subscriptionId string
}

func (testsuite *RedisenterpriseOperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *RedisenterpriseOperationsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestRedisenterpriseOperatopnsTestSuite(t *testing.T) {
	suite.Run(t, new(RedisenterpriseOperationsTestSuite))
}

func (testsuite *RedisenterpriseOperationsTestSuite) TestNewListPager() {
	testsuite.options.APIVersion = "2024-11-01"
	clientFactory, err := armredisenterprise.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
	}
}
