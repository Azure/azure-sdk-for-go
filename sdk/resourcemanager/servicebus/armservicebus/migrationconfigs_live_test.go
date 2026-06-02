// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armservicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"github.com/stretchr/testify/suite"
)

type MigrationconfigsTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	namespaceName       string
	namespaceNameSecond string
	secondNamespaceId   string
	postMigrationName   string
	location            string
	resourceGroupName   string
	subscriptionId      string
}

func (testsuite *MigrationconfigsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespac", 14, false)
	testsuite.namespaceNameSecond, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespacsecond", 20, false)
	testsuite.postMigrationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "postmigrationna", 21, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MigrationconfigsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMigrationconfigsTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationconfigsTestSuite))
}

func (testsuite *MigrationconfigsTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armservicebus.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.SBNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNameStandard),
			Tier: to.Ptr(armservicebus.SKUTierStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Namespaces_CreateOrUpdate_Second
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClientCreateOrUpdateResponsePoller, err = namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceNameSecond, armservicebus.SBNamespace{
		Location: to.Ptr("westus"),
		SKU: &armservicebus.SBSKU{
			Name: to.Ptr(armservicebus.SKUNamePremium),
			Tier: to.Ptr(armservicebus.SKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var namespacesClientCreateOrUpdateResponse *armservicebus.NamespacesClientCreateOrUpdateResponse
	namespacesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.secondNamespaceId = *namespacesClientCreateOrUpdateResponse.ID
}

// Microsoft.ServiceBus/namespaces/{namespaceName}/migrationConfigurations/{configName}
func (testsuite *MigrationconfigsTestSuite) TestMigrationConfigs() {
	var err error
	// From step MigrationConfigs_CreateAndStartMigration
	fmt.Println("Call operation: MigrationConfigs_CreateAndStartMigration")
	migrationConfigsClient, err := armservicebus.NewMigrationConfigsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	migrationConfigsClientCreateAndStartMigrationResponsePoller, err := migrationConfigsClient.BeginCreateAndStartMigration(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.MigrationConfigurationNameDefault, armservicebus.MigrationConfigProperties{
		Properties: &armservicebus.MigrationConfigPropertiesProperties{
			PostMigrationName: to.Ptr(testsuite.postMigrationName),
			TargetNamespace:   to.Ptr(testsuite.secondNamespaceId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, migrationConfigsClientCreateAndStartMigrationResponsePoller)
	testsuite.Require().NoError(err)

	// From step MigrationConfigs_List
	fmt.Println("Call operation: MigrationConfigs_List")
	migrationConfigsClientNewListPager := migrationConfigsClient.NewListPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for migrationConfigsClientNewListPager.More() {
		_, err := migrationConfigsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MigrationConfigs_Get
	fmt.Println("Call operation: MigrationConfigs_Get")
	_, err = migrationConfigsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.MigrationConfigurationNameDefault, nil)
	testsuite.Require().NoError(err)

	// From step MigrationConfigs_Revert
	fmt.Println("Call operation: MigrationConfigs_Revert")
	_, err = migrationConfigsClient.Revert(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.MigrationConfigurationNameDefault, nil)
	testsuite.Require().NoError(err)

	// From step MigrationConfigs_CompleteMigration
	fmt.Println("Call operation: MigrationConfigs_CompleteMigration")
	_, err = migrationConfigsClient.CompleteMigration(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armservicebus.MigrationConfigurationNameDefault, nil)
	testsuite.Require().NoError(err)
}
