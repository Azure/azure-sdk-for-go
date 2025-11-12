//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetworkcloud_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/networkcloud/armnetworkcloud"
	"github.com/stretchr/testify/suite"
)

type OperationTestSuite struct {
	suite.Suite

	ctx            context.Context
	cred           azcore.TokenCredential
	options        *arm.ClientOptions
	location       string
	subscriptionId string
}

func (testsuite *OperationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
}

func (testsuite *OperationTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestOperationTestSuite(t *testing.T) {
	suite.Run(t, new(OperationTestSuite))
}

func (testsuite *OperationTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	ctx := context.Background()
	clientFactory, err := armnetworkcloud.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		testsuite.Require().NoError(err)
	}
}
