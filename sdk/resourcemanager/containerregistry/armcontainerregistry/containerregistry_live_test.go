//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerregistry_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ContainerregistryTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	exportPipelineName string
	importPipelineName string
	registryName       string
	replicationName    string
	scopeMapId         string
	scopeMapName       string
	tokenId            string
	tokenName          string
	webhookName        string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *ContainerregistryTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.exportPipelineName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "exportpipe", 16, false)
	testsuite.importPipelineName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "importpipe", 16, false)
	testsuite.registryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "registryna", 16, false)
	testsuite.replicationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "replicatio", 16, false)
	testsuite.scopeMapName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "scopemapna", 16, false)
	testsuite.tokenName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tokenname", 19, false)
	testsuite.webhookName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "webhooknam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ContainerregistryTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestContainerregistryTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerregistryTestSuite))
}

func (testsuite *ContainerregistryTestSuite) Prepare() {
	var err error
	// From step Registries_Create
	fmt.Println("Call operation: Registries_Create")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	registriesClientCreateResponsePoller, err := registriesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.Registry{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
		Properties: &armcontainerregistry.RegistryProperties{
			AdminUserEnabled: to.Ptr(true),
		},
		SKU: &armcontainerregistry.SKU{
			Name: to.Ptr(armcontainerregistry.SKUNamePremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ScopeMaps_Create
	fmt.Println("Call operation: ScopeMaps_Create")
	scopeMapsClient, err := armcontainerregistry.NewScopeMapsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scopeMapsClientCreateResponsePoller, err := scopeMapsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.scopeMapName, armcontainerregistry.ScopeMap{
		Properties: &armcontainerregistry.ScopeMapProperties{
			Description: to.Ptr("Developer Scopes"),
			Actions: []*string{
				to.Ptr("repositories/myrepository/content/write"),
				to.Ptr("repositories/myrepository/content/delete")},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var scopeMapsClientCreateResponse *armcontainerregistry.ScopeMapsClientCreateResponse
	scopeMapsClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, scopeMapsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.scopeMapId = *scopeMapsClientCreateResponse.ID

	// From step Tokens_Create
	fmt.Println("Call operation: Tokens_Create")
	tokensClient, err := armcontainerregistry.NewTokensClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tokensClientCreateResponsePoller, err := tokensClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.tokenName, armcontainerregistry.Token{
		Properties: &armcontainerregistry.TokenProperties{
			ScopeMapID: to.Ptr(testsuite.scopeMapId),
			Status:     to.Ptr(armcontainerregistry.TokenStatusDisabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var tokensClientCreateResponse *armcontainerregistry.TokensClientCreateResponse
	tokensClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, tokensClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.tokenId = *tokensClientCreateResponse.ID
}

// Microsoft.ContainerRegistry/registries
func (testsuite *ContainerregistryTestSuite) TestContainerregister() {
	var err error
	// From step Registries_CheckNameAvailability
	fmt.Println("Call operation: Registries_CheckNameAvailability")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = registriesClient.CheckNameAvailability(testsuite.ctx, armcontainerregistry.RegistryNameCheckRequest{
		Name: to.Ptr("myRegistry"),
		Type: to.Ptr("Microsoft.ContainerRegistry/registries"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Registries_List
	fmt.Println("Call operation: Registries_List")
	registriesClientNewListPager := registriesClient.NewListPager(nil)
	for registriesClientNewListPager.More() {
		_, err := registriesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Registries_Get
	fmt.Println("Call operation: Registries_Get")
	_, err = registriesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, nil)
	testsuite.Require().NoError(err)

	// From step Registries_ListByResourceGroup
	fmt.Println("Call operation: Registries_ListByResourceGroup")
	registriesClientNewListByResourceGroupPager := registriesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for registriesClientNewListByResourceGroupPager.More() {
		_, err := registriesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Registries_ListUsages
	fmt.Println("Call operation: Registries_ListUsages")
	_, err = registriesClient.ListUsages(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, nil)
	testsuite.Require().NoError(err)

	// From step Registries_Update
	fmt.Println("Call operation: Registries_Update")
	registriesClientUpdateResponsePoller, err := registriesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.RegistryUpdateParameters{
		Properties: &armcontainerregistry.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: to.Ptr(true),
		},
		SKU: &armcontainerregistry.SKU{
			Name: to.Ptr(armcontainerregistry.SKUNamePremium),
		},
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Registries_ListPrivateLinkResources
	fmt.Println("Call operation: Registries_ListPrivateLinkResources")
	registriesClientNewListPrivateLinkResourcesPager := registriesClient.NewListPrivateLinkResourcesPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for registriesClientNewListPrivateLinkResourcesPager.More() {
		_, err := registriesClientNewListPrivateLinkResourcesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Registries_GetPrivateLinkResource
	fmt.Println("Call operation: Registries_GetPrivateLinkResource")
	_, err = registriesClient.GetPrivateLinkResource(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, "registry", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/webhooks
func (testsuite *ContainerregistryTestSuite) TestWebhooks() {
	var err error
	// From step Webhooks_Create
	fmt.Println("Call operation: Webhooks_Create")
	webhooksClient, err := armcontainerregistry.NewWebhooksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	webhooksClientCreateResponsePoller, err := webhooksClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, armcontainerregistry.WebhookCreateParameters{
		Location: to.Ptr(testsuite.location),
		Properties: &armcontainerregistry.WebhookPropertiesCreateParameters{
			Actions: []*armcontainerregistry.WebhookAction{
				to.Ptr(armcontainerregistry.WebhookActionPush)},
			CustomHeaders: map[string]*string{
				"Authorization": to.Ptr("Basic 000000000000000000000000000000000000000000000000000"),
			},
			ServiceURI: to.Ptr("http://myservice.com"),
			Status:     to.Ptr(armcontainerregistry.WebhookStatusEnabled),
		},
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, webhooksClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Webhooks_List
	fmt.Println("Call operation: Webhooks_List")
	webhooksClientNewListPager := webhooksClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for webhooksClientNewListPager.More() {
		_, err := webhooksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Webhooks_Get
	fmt.Println("Call operation: Webhooks_Get")
	_, err = webhooksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, nil)
	testsuite.Require().NoError(err)

	// From step Webhooks_Update
	fmt.Println("Call operation: Webhooks_Update")
	webhooksClientUpdateResponsePoller, err := webhooksClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, armcontainerregistry.WebhookUpdateParameters{
		Properties: &armcontainerregistry.WebhookPropertiesUpdateParameters{
			Actions: []*armcontainerregistry.WebhookAction{
				to.Ptr(armcontainerregistry.WebhookActionPush)},
			CustomHeaders: map[string]*string{
				"Authorization": to.Ptr("Basic 000000000000000000000000000000000000000000000000000"),
			},
			Scope:      to.Ptr("repository"),
			ServiceURI: to.Ptr("http://myservice.com"),
			Status:     to.Ptr(armcontainerregistry.WebhookStatusEnabled),
		},
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, webhooksClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Webhooks_GetCallbackConfig
	fmt.Println("Call operation: Webhooks_GetCallbackConfig")
	_, err = webhooksClient.GetCallbackConfig(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, nil)
	testsuite.Require().NoError(err)

	// From step Webhooks_Ping
	fmt.Println("Call operation: Webhooks_Ping")
	_, err = webhooksClient.Ping(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, nil)
	testsuite.Require().NoError(err)

	// From step Webhooks_ListEvents
	fmt.Println("Call operation: Webhooks_ListEvents")
	webhooksClientNewListEventsPager := webhooksClient.NewListEventsPager(testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, nil)
	for webhooksClientNewListEventsPager.More() {
		_, err := webhooksClientNewListEventsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Webhooks_Delete
	fmt.Println("Call operation: Webhooks_Delete")
	webhooksClientDeleteResponsePoller, err := webhooksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.webhookName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, webhooksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/replications
func (testsuite *ContainerregistryTestSuite) TestReplications() {
	var err error
	// From step Replications_Create
	fmt.Println("Call operation: Replications_Create")
	replicationsClient, err := armcontainerregistry.NewReplicationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	replicationsClientCreateResponsePoller, err := replicationsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.replicationName, armcontainerregistry.Replication{
		Location: to.Ptr("westus2"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, replicationsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Replications_List
	fmt.Println("Call operation: Replications_List")
	replicationsClientNewListPager := replicationsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for replicationsClientNewListPager.More() {
		_, err := replicationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Replications_Get
	fmt.Println("Call operation: Replications_Get")
	_, err = replicationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.replicationName, nil)
	testsuite.Require().NoError(err)

	// From step Replications_Update
	fmt.Println("Call operation: Replications_Update")
	replicationsClientUpdateResponsePoller, err := replicationsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.replicationName, armcontainerregistry.ReplicationUpdateParameters{
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, replicationsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Replications_Delete
	fmt.Println("Call operation: Replications_Delete")
	replicationsClientDeleteResponsePoller, err := replicationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.replicationName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, replicationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/scopeMaps
func (testsuite *ContainerregistryTestSuite) TestScopemaps() {
	var err error
	// From step ScopeMaps_List
	fmt.Println("Call operation: ScopeMaps_List")
	scopeMapsClient, err := armcontainerregistry.NewScopeMapsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scopeMapsClientNewListPager := scopeMapsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for scopeMapsClientNewListPager.More() {
		_, err := scopeMapsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ScopeMaps_Get
	fmt.Println("Call operation: ScopeMaps_Get")
	_, err = scopeMapsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.scopeMapName, nil)
	testsuite.Require().NoError(err)

	// From step ScopeMaps_Update
	fmt.Println("Call operation: ScopeMaps_Update")
	scopeMapsClientUpdateResponsePoller, err := scopeMapsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.scopeMapName, armcontainerregistry.ScopeMapUpdateParameters{
		Properties: &armcontainerregistry.ScopeMapPropertiesUpdateParameters{
			Description: to.Ptr("Developer Scopes"),
			Actions: []*string{
				to.Ptr("repositories/myrepository/content/write"),
				to.Ptr("repositories/myrepository/content/read")},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, scopeMapsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/tokens
func (testsuite *ContainerregistryTestSuite) TestTokens() {
	var err error
	// From step Tokens_List
	fmt.Println("Call operation: Tokens_List")
	tokensClient, err := armcontainerregistry.NewTokensClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tokensClientNewListPager := tokensClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for tokensClientNewListPager.More() {
		_, err := tokensClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Tokens_Get
	fmt.Println("Call operation: Tokens_Get")
	_, err = tokensClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.tokenName, nil)
	testsuite.Require().NoError(err)

	// From step Tokens_Update
	fmt.Println("Call operation: Tokens_Update")
	tokensClientUpdateResponsePoller, err := tokensClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.tokenName, armcontainerregistry.TokenUpdateParameters{
		Properties: &armcontainerregistry.TokenUpdateProperties{
			ScopeMapID: to.Ptr(testsuite.scopeMapId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tokensClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/generateCredentials
func (testsuite *ContainerregistryTestSuite) TestRegistrycredentials() {
	var err error
	// From step Registries_GenerateCredentials
	fmt.Println("Call operation: Registries_GenerateCredentials")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	registriesClientGenerateCredentialsResponsePoller, err := registriesClient.BeginGenerateCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.GenerateCredentialsParameters{
		TokenID: to.Ptr(testsuite.tokenId),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientGenerateCredentialsResponsePoller)
	testsuite.Require().NoError(err)

	// From step Registries_RegenerateCredential
	fmt.Println("Call operation: Registries_RegenerateCredential")
	_, err = registriesClient.RegenerateCredential(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.RegenerateCredentialParameters{
		Name: to.Ptr(armcontainerregistry.PasswordNamePassword),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Registries_ListCredentials
	fmt.Println("Call operation: Registries_ListCredentials")
	_, err = registriesClient.ListCredentials(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/operations
func (testsuite *ContainerregistryTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armcontainerregistry.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ContainerregistryTestSuite) Cleanup() {
	var err error
	// From step Tokens_Delete
	fmt.Println("Call operation: Tokens_Delete")
	tokensClient, err := armcontainerregistry.NewTokensClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tokensClientDeleteResponsePoller, err := tokensClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.tokenName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tokensClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step ScopeMaps_Delete
	fmt.Println("Call operation: ScopeMaps_Delete")
	scopeMapsClient, err := armcontainerregistry.NewScopeMapsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scopeMapsClientDeleteResponsePoller, err := scopeMapsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.scopeMapName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, scopeMapsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Registries_Delete
	fmt.Println("Call operation: Registries_Delete")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	registriesClientDeleteResponsePoller, err := registriesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
