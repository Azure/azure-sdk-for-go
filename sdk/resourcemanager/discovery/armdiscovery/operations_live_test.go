// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type OperationsTestSuite struct {
	suite.Suite
	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	location       string
	subscriptionId string
}

func (testsuite *OperationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *OperationsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func SkipTestOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(OperationsTestSuite))
}

// Test listing available API operations
func (testsuite *OperationsTestSuite) TestOperationsNewListPager() {
	fmt.Println("Call operation: Operations_List")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	testsuite.Require().True(pager.More())

	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		if len(result.Value) > 0 {
			testsuite.Require().NotNil(result.Value[0].Name)
		}
		break // Just verify first page
	}
}
