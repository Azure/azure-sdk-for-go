//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

type StoragemoverTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	agentName         string
	armEndpoint       string
	endpointName      string
	jobDefinitionName string
	projectName       string
	storageMoverName  string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *StoragemoverTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "agentnam", 14, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.endpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "endpoint", 14, false)
	testsuite.jobDefinitionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobdefin", 14, false)
	testsuite.projectName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "projectn", 14, false)
	testsuite.storageMoverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storagem", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *StoragemoverTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestStoragemoverTestSuite(t *testing.T) {
	suite.Run(t, new(StoragemoverTestSuite))
}

func (testsuite *StoragemoverTestSuite) Prepare() {
	var err error
	// From step StorageMovers_CreateOrUpdate
	fmt.Println("Call operation: StorageMovers_CreateOrUpdate")
	storageMoversClient, err := armstoragemover.NewStorageMoversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = storageMoversClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, armstoragemover.StorageMover{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
		Properties: &armstoragemover.Properties{
			Description: to.Ptr("Example Storage Mover Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Projects_CreateOrUpdate
	fmt.Println("Call operation: Projects_CreateOrUpdate")
	projectsClient, err := armstoragemover.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = projectsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, armstoragemover.Project{
		Properties: &armstoragemover.ProjectProperties{
			Description: to.Ptr("Example Project Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/storageMovers/{storageMoverName}
func (testsuite *StoragemoverTestSuite) TestStorageMovers() {
	var err error
	// From step StorageMovers_ListBySubscription
	fmt.Println("Call operation: StorageMovers_ListBySubscription")
	storageMoversClient, err := armstoragemover.NewStorageMoversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storageMoversClientNewListBySubscriptionPager := storageMoversClient.NewListBySubscriptionPager(nil)
	for storageMoversClientNewListBySubscriptionPager.More() {
		_, err := storageMoversClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageMovers_Get
	fmt.Println("Call operation: StorageMovers_Get")
	_, err = storageMoversClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	testsuite.Require().NoError(err)

	// From step StorageMovers_List
	fmt.Println("Call operation: StorageMovers_List")
	storageMoversClientNewListPager := storageMoversClient.NewListPager(testsuite.resourceGroupName, nil)
	for storageMoversClientNewListPager.More() {
		_, err := storageMoversClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StorageMovers_Update
	fmt.Println("Call operation: StorageMovers_Update")
	_, err = storageMoversClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, armstoragemover.UpdateParameters{
		Properties: &armstoragemover.UpdateProperties{
			Description: to.Ptr("Updated Storage Mover Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/storageMovers/{storageMoverName}/projects/{projectName}
func (testsuite *StoragemoverTestSuite) TestProjects() {
	var err error
	// From step Projects_List
	fmt.Println("Call operation: Projects_List")
	projectsClient, err := armstoragemover.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	projectsClientNewListPager := projectsClient.NewListPager(testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	for projectsClientNewListPager.More() {
		_, err := projectsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Projects_Get
	fmt.Println("Call operation: Projects_Get")
	_, err = projectsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)

	// From step Projects_Update
	fmt.Println("Call operation: Projects_Update")
	_, err = projectsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, armstoragemover.ProjectUpdateParameters{
		Properties: &armstoragemover.ProjectUpdateProperties{
			Description: to.Ptr("Updated Example Project Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/storageMovers/{storageMoverName}/projects/{projectName}/jobDefinitions/{jobDefinitionName}
func (testsuite *StoragemoverTestSuite) TestJobDefinitions() {
	var err error
	// From step Endpoints_CreateOrUpdate_Source
	fmt.Println("Call operation: Endpoints_CreateOrUpdate")
	endpointsClient, err := armstoragemover.NewEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = endpointsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, "examples-sourceEndpointName", armstoragemover.Endpoint{
		Properties: &armstoragemover.NfsMountEndpointProperties{
			Description:  to.Ptr("Example NFS Mount Endpoint Description"),
			EndpointType: to.Ptr(armstoragemover.EndpointTypeNfsMount),
			Export:       to.Ptr("examples-exportName"),
			Host:         to.Ptr("0.0.0.0"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_CreateOrUpdate_Target
	fmt.Println("Call operation: Endpoints_CreateOrUpdate")
	_, err = endpointsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, "examples-targetEndpointName", armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			Description:              to.Ptr("Example Azure Storage Blob Container Endpoint Description"),
			EndpointType:             to.Ptr(armstoragemover.EndpointTypeAzureStorageBlobContainer),
			BlobContainerName:        to.Ptr("examples-blobcontainer"),
			StorageAccountResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.Storage/storageAccounts/examplesa"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step JobDefinitions_CreateOrUpdate
	fmt.Println("Call operation: JobDefinitions_CreateOrUpdate")
	jobDefinitionsClient, err := armstoragemover.NewJobDefinitionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = jobDefinitionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, testsuite.jobDefinitionName, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			Description:   to.Ptr("Example Job Definition Description"),
			CopyMode:      to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName:    to.Ptr("examples-sourceEndpointName"),
			SourceSubpath: to.Ptr("/"),
			TargetName:    to.Ptr("examples-targetEndpointName"),
			TargetSubpath: to.Ptr("/"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step JobDefinitions_List
	fmt.Println("Call operation: JobDefinitions_List")
	jobDefinitionsClientNewListPager := jobDefinitionsClient.NewListPager(testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, nil)
	for jobDefinitionsClientNewListPager.More() {
		_, err := jobDefinitionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step JobDefinitions_Get
	fmt.Println("Call operation: JobDefinitions_Get")
	_, err = jobDefinitionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, testsuite.jobDefinitionName, nil)
	testsuite.Require().NoError(err)

	// From step JobDefinitions_Update
	fmt.Println("Call operation: JobDefinitions_Update")
	_, err = jobDefinitionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, testsuite.jobDefinitionName, armstoragemover.JobDefinitionUpdateParameters{
		Properties: &armstoragemover.JobDefinitionUpdateProperties{
			Description: to.Ptr("Updated Job Definition Description"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step JobDefinitions_Delete
	fmt.Println("Call operation: JobDefinitions_Delete")
	jobDefinitionsClientDeleteResponsePoller, err := jobDefinitionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, testsuite.jobDefinitionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, jobDefinitionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/storageMovers/{storageMoverName}/endpoints/{endpointName}
func (testsuite *StoragemoverTestSuite) TestEndpoints() {
	var err error
	// From step Endpoints_CreateOrUpdate
	fmt.Println("Call operation: Endpoints_CreateOrUpdate")
	endpointsClient, err := armstoragemover.NewEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = endpointsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.endpointName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			Description:              to.Ptr("Example Storage Blob Container Endpoint Description"),
			EndpointType:             to.Ptr(armstoragemover.EndpointTypeAzureStorageBlobContainer),
			BlobContainerName:        to.Ptr("examples-blobcontainer"),
			StorageAccountResourceID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Storage/storageAccounts/examplesa"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_List
	fmt.Println("Call operation: Endpoints_List")
	endpointsClientNewListPager := endpointsClient.NewListPager(testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	for endpointsClientNewListPager.More() {
		_, err := endpointsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Endpoints_Get
	fmt.Println("Call operation: Endpoints_Get")
	_, err = endpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_Update
	fmt.Println("Call operation: Endpoints_Update")
	_, err = endpointsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.endpointName, armstoragemover.EndpointBaseUpdateParameters{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointUpdateProperties{
			Description:  to.Ptr("Updated Endpoint Description"),
			EndpointType: to.Ptr(armstoragemover.EndpointTypeAzureStorageBlobContainer),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_Delete
	fmt.Println("Call operation: Endpoints_Delete")
	endpointsClientDeleteResponsePoller, err := endpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StorageMover/operations
func (testsuite *StoragemoverTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armstoragemover.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *StoragemoverTestSuite) Cleanup() {
	var err error
	// From step Projects_Delete
	fmt.Println("Call operation: Projects_Delete")
	projectsClient, err := armstoragemover.NewProjectsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	projectsClientDeleteResponsePoller, err := projectsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, testsuite.projectName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, projectsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageMovers_Delete
	fmt.Println("Call operation: StorageMovers_Delete")
	storageMoversClient, err := armstoragemover.NewStorageMoversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storageMoversClientDeleteResponsePoller, err := storageMoversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageMoverName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storageMoversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
