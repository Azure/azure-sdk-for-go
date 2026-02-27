// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimtenantTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimtenantTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicetenant", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimtenantTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimtenantTestSuite(t *testing.T) {
	suite.Run(t, new(ApimtenantTestSuite))
}

func (testsuite *ApimtenantTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name":	to.Ptr("Contoso"),
			"Test":	to.Ptr("User"),
		},
		Location:	to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail:	to.Ptr("foo@contoso.com"),
			PublisherName:	to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:		to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity:	to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TenantAccess_Create
	fmt.Println("Call operation: TenantAccess_Create")
	tenantAccessClient, err := armapimanagement.NewTenantAccessClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tenantAccessClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, "*", armapimanagement.AccessInformationCreateParameters{
		Properties: &armapimanagement.AccessInformationCreateParameterProperties{
			Enabled: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/tenant
func (testsuite *ApimtenantTestSuite) TestTenantaccess() {
	var err error
	// From step TenantAccess_GetEntityTag
	fmt.Println("Call operation: TenantAccess_GetEntityTag")
	tenantAccessClient, err := armapimanagement.NewTenantAccessClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tenantAccessClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccess_ListByService
	fmt.Println("Call operation: TenantAccess_ListByService")
	tenantAccessClientNewListByServicePager := tenantAccessClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.TenantAccessClientListByServiceOptions{Filter: nil})
	for tenantAccessClientNewListByServicePager.More() {
		_, err := tenantAccessClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TenantAccess_Get
	fmt.Println("Call operation: TenantAccess_Get")
	_, err = tenantAccessClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccess_Update
	fmt.Println("Call operation: TenantAccess_Update")
	_, err = tenantAccessClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, "*", armapimanagement.AccessInformationUpdateParameters{
		Properties: &armapimanagement.AccessInformationUpdateParameterProperties{
			Enabled: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccess_RegenerateSecondaryKey
	fmt.Println("Call operation: TenantAccess_RegenerateSecondaryKey")
	_, err = tenantAccessClient.RegenerateSecondaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccess_ListSecrets
	fmt.Println("Call operation: TenantAccess_ListSecrets")
	_, err = tenantAccessClient.ListSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccess_RegeneratePrimaryKey
	fmt.Println("Call operation: TenantAccess_RegeneratePrimaryKey")
	_, err = tenantAccessClient.RegeneratePrimaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/tenant/git
func (testsuite *ApimtenantTestSuite) TestTenantaccessgit() {
	var err error
	// From step TenantAccessGit_RegeneratePrimaryKey
	fmt.Println("Call operation: TenantAccessGit_RegeneratePrimaryKey")
	tenantAccessGitClient, err := armapimanagement.NewTenantAccessGitClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tenantAccessGitClient.RegeneratePrimaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)

	// From step TenantAccessGit_RegenerateSecondaryKey
	fmt.Println("Call operation: TenantAccessGit_RegenerateSecondaryKey")
	_, err = tenantAccessGitClient.RegenerateSecondaryKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.AccessIDNameAccess, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/tenant/configuration
func (testsuite *ApimtenantTestSuite) TestTenantconfiguration() {
	var err error
	// From step TenantConfiguration_Save
	fmt.Println("Call operation: TenantConfiguration_Save")
	tenantConfigurationClient, err := armapimanagement.NewTenantConfigurationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tenantConfigurationClientSaveResponsePoller, err := tenantConfigurationClient.BeginSave(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ConfigurationIDNameConfiguration, armapimanagement.SaveConfigurationParameter{
		Properties: &armapimanagement.SaveConfigurationParameterProperties{
			Branch: to.Ptr("master"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tenantConfigurationClientSaveResponsePoller)
	testsuite.Require().NoError(err)

	// From step TenantConfiguration_GetSyncState
	fmt.Println("Call operation: TenantConfiguration_GetSyncState")
	_, err = tenantConfigurationClient.GetSyncState(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ConfigurationIDNameConfiguration, nil)
	testsuite.Require().NoError(err)

	// From step TenantConfiguration_Validate
	fmt.Println("Call operation: TenantConfiguration_Validate")
	tenantConfigurationClientValidateResponsePoller, err := tenantConfigurationClient.BeginValidate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ConfigurationIDNameConfiguration, armapimanagement.DeployConfigurationParameters{
		Properties: &armapimanagement.DeployConfigurationParameterProperties{
			Branch: to.Ptr("master"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tenantConfigurationClientValidateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TenantConfiguration_Deploy
	fmt.Println("Call operation: TenantConfiguration_Deploy")
	tenantConfigurationClientDeployResponsePoller, err := tenantConfigurationClient.BeginDeploy(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ConfigurationIDNameConfiguration, armapimanagement.DeployConfigurationParameters{
		Properties: &armapimanagement.DeployConfigurationParameterProperties{
			Branch: to.Ptr("master"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tenantConfigurationClientDeployResponsePoller)
	testsuite.Require().NoError(err)
}
