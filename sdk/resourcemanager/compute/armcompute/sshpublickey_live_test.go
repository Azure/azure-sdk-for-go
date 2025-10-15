//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type SshPublicKeyTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	sshPublicKeyName  string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SshPublicKeyTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.sshPublicKeyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sshpublick", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *SshPublicKeyTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSshPublicKeyTestSuite(t *testing.T) {
	suite.Run(t, new(SshPublicKeyTestSuite))
}

// Microsoft.Compute/sshPublicKeys/{sshPublicKeyName}
func (testsuite *SshPublicKeyTestSuite) TestSshPublicKeys() {
	var err error
	// From step SshPublicKeys_Create
	fmt.Println("Call operation: SSHPublicKeys_Create")
	sSHPublicKeysClient, err := armcompute.NewSSHPublicKeysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = sSHPublicKeysClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.sshPublicKeyName, armcompute.SSHPublicKeyResource{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)

	// From step SshPublicKeys_ListBySubscription
	fmt.Println("Call operation: SSHPublicKeys_ListBySubscription")
	sSHPublicKeysClientNewListBySubscriptionPager := sSHPublicKeysClient.NewListBySubscriptionPager(nil)
	for sSHPublicKeysClientNewListBySubscriptionPager.More() {
		_, err := sSHPublicKeysClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SshPublicKeys_ListByResourceGroup
	fmt.Println("Call operation: SSHPublicKeys_ListByResourceGroup")
	sSHPublicKeysClientNewListByResourceGroupPager := sSHPublicKeysClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for sSHPublicKeysClientNewListByResourceGroupPager.More() {
		_, err := sSHPublicKeysClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SshPublicKeys_Get
	fmt.Println("Call operation: SSHPublicKeys_Get")
	_, err = sSHPublicKeysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.sshPublicKeyName, nil)
	testsuite.Require().NoError(err)

	// From step SshPublicKeys_Update
	fmt.Println("Call operation: SSHPublicKeys_Update")
	_, err = sSHPublicKeysClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.sshPublicKeyName, armcompute.SSHPublicKeyUpdateResource{
		Tags: map[string]*string{
			"key2854": to.Ptr("a"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SshPublicKeys_GenerateKeyPair
	fmt.Println("Call operation: SSHPublicKeys_GenerateKeyPair")
	_, err = sSHPublicKeysClient.GenerateKeyPair(testsuite.ctx, testsuite.resourceGroupName, testsuite.sshPublicKeyName, nil)
	testsuite.Require().NoError(err)

	// From step SshPublicKeys_Delete
	fmt.Println("Call operation: SSHPublicKeys_Delete")
	_, err = sSHPublicKeysClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.sshPublicKeyName, nil)
	testsuite.Require().NoError(err)
}
