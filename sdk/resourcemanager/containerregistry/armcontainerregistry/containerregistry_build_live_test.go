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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/stretchr/testify/suite"
)

type ContainerregistryBuildTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	agentPoolName     string
	registryName      string
	taskName          string
	taskRunName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ContainerregistryBuildTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/containerregistry/armcontainerregistry/testdata")
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agentPoolName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "agentpooln", 16, false)
	testsuite.registryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "registryna2", 17, false)
	testsuite.taskName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "taskname", 14, false)
	testsuite.taskRunName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "taskrunnam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ContainerregistryBuildTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestContainerregistryBuildTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerregistryBuildTestSuite))
}

func (testsuite *ContainerregistryBuildTestSuite) Prepare() {
	var err error
	// From step Registries_Create2
	fmt.Println("Call operation: Registries_Create")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	registriesClientCreateResponsePoller, err := registriesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, armcontainerregistry.Registry{
		Location: to.Ptr(testsuite.location),
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
}

// Microsoft.ContainerRegistry/registries/agentPools
func (testsuite *ContainerregistryBuildTestSuite) TestAgentpools() {
	var err error
	// From step AgentPools_Create
	fmt.Println("Call operation: AgentPools_Create")
	agentPoolsClient, err := armcontainerregistry.NewAgentPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	agentPoolsClientCreateResponsePoller, err := agentPoolsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, armcontainerregistry.AgentPool{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key": to.Ptr("value"),
		},
		Properties: &armcontainerregistry.AgentPoolProperties{
			Count: to.Ptr[int32](1),
			OS:    to.Ptr(armcontainerregistry.OSLinux),
			Tier:  to.Ptr("S1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AgentPools_List
	fmt.Println("Call operation: AgentPools_List")
	agentPoolsClientNewListPager := agentPoolsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for agentPoolsClientNewListPager.More() {
		_, err := agentPoolsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AgentPools_Get
	fmt.Println("Call operation: AgentPools_Get")
	_, err = agentPoolsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Update
	fmt.Println("Call operation: AgentPools_Update")
	agentPoolsClientUpdateResponsePoller, err := agentPoolsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, armcontainerregistry.AgentPoolUpdateParameters{
		Properties: &armcontainerregistry.AgentPoolPropertiesUpdateParameters{
			Count: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AgentPools_GetQueueStatus
	fmt.Println("Call operation: AgentPools_GetQueueStatus")
	_, err = agentPoolsClient.GetQueueStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)

	// From step AgentPools_Delete
	fmt.Println("Call operation: AgentPools_Delete")
	agentPoolsClientDeleteResponsePoller, err := agentPoolsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.agentPoolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, agentPoolsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ContainerRegistry/registries/tasks
func (testsuite *ContainerregistryBuildTestSuite) TestTasks() {
	var runId string
	var err error
	// From step Tasks_Create
	fmt.Println("Call operation: Tasks_Create")
	tasksClient, err := armcontainerregistry.NewTasksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tasksClientCreateResponsePoller, err := tasksClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskName, armcontainerregistry.Task{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"testkey": to.Ptr("value"),
		},
		Properties: &armcontainerregistry.TaskProperties{
			AgentConfiguration: &armcontainerregistry.AgentProperties{
				CPU: to.Ptr[int32](2),
			},
			IsSystemTask: to.Ptr(false),
			Platform: &armcontainerregistry.PlatformProperties{
				Architecture: to.Ptr(armcontainerregistry.ArchitectureAmd64),
				OS:           to.Ptr(armcontainerregistry.OSLinux),
			},
			Status: to.Ptr(armcontainerregistry.TaskStatusEnabled),
			Step: &armcontainerregistry.DockerBuildStep{
				Type:           to.Ptr(armcontainerregistry.StepTypeDocker),
				ContextPath:    to.Ptr("https://github.com/SteveLasker/node-helloworld"),
				DockerFilePath: to.Ptr("DockerFile"),
				ImageNames: []*string{
					to.Ptr("testtask:v1")},
				IsPushEnabled: to.Ptr(true),
				NoCache:       to.Ptr(false),
			},
			Trigger: &armcontainerregistry.TriggerProperties{
				BaseImageTrigger: &armcontainerregistry.BaseImageTrigger{
					Name:                     to.Ptr("myBaseImageTrigger"),
					BaseImageTriggerType:     to.Ptr(armcontainerregistry.BaseImageTriggerTypeRuntime),
					Status:                   to.Ptr(armcontainerregistry.TriggerStatusEnabled),
					UpdateTriggerPayloadType: to.Ptr(armcontainerregistry.UpdateTriggerPayloadTypeDefault),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tasksClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Tasks_List
	fmt.Println("Call operation: Tasks_List")
	tasksClientNewListPager := tasksClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for tasksClientNewListPager.More() {
		_, err := tasksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Tasks_Get
	fmt.Println("Call operation: Tasks_Get")
	_, err = tasksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskName, nil)
	testsuite.Require().NoError(err)

	// From step Tasks_Update
	fmt.Println("Call operation: Tasks_Update")
	tasksClientUpdateResponsePoller, err := tasksClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskName, armcontainerregistry.TaskUpdateParameters{
		Tags: map[string]*string{
			"testkey": to.Ptr("value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tasksClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Tasks_GetDetails
	fmt.Println("Call operation: Tasks_GetDetails")
	_, err = tasksClient.GetDetails(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskName, nil)
	testsuite.Require().NoError(err)

	// From step TaskRuns_Create
	fmt.Println("Call operation: TaskRuns_Create")
	taskRunsClient, err := armcontainerregistry.NewTaskRunsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	taskRunsClientCreateResponsePoller, err := taskRunsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskRunName, armcontainerregistry.TaskRun{
		Properties: &armcontainerregistry.TaskRunProperties{
			ForceUpdateTag: to.Ptr("test"),
			RunRequest: &armcontainerregistry.EncodedTaskRunRequest{
				Type:                 to.Ptr("EncodedTaskRunRequest"),
				Credentials:          &armcontainerregistry.Credentials{},
				EncodedTaskContent:   to.Ptr("c3RlcHM6IAogIC0gY21kOiB7eyAuVmFsdWVzLmNvbW1hbmQgfX0K"),
				EncodedValuesContent: to.Ptr("Y29tbWFuZDogYmFzaCBlY2hvIHt7LlJ1bi5SZWdpc3RyeX19Cg=="),
				Platform: &armcontainerregistry.PlatformProperties{
					Architecture: to.Ptr(armcontainerregistry.ArchitectureAmd64),
					OS:           to.Ptr(armcontainerregistry.OSLinux),
				},
				Values: []*armcontainerregistry.SetValue{},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var taskRunsClientCreateResponse *armcontainerregistry.TaskRunsClientCreateResponse
	taskRunsClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, taskRunsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	runId = *taskRunsClientCreateResponse.Properties.RunResult.Properties.RunID

	// From step TaskRuns_List
	fmt.Println("Call operation: TaskRuns_List")
	taskRunsClientNewListPager := taskRunsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, nil)
	for taskRunsClientNewListPager.More() {
		_, err := taskRunsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TaskRuns_Get
	fmt.Println("Call operation: TaskRuns_Get")
	_, err = taskRunsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskRunName, nil)
	testsuite.Require().NoError(err)

	// From step TaskRuns_Update
	fmt.Println("Call operation: TaskRuns_Update")
	taskRunsClientUpdateResponsePoller, err := taskRunsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskRunName, armcontainerregistry.TaskRunUpdateParameters{
		Properties: &armcontainerregistry.TaskRunPropertiesUpdateParameters{
			ForceUpdateTag: to.Ptr("test"),
			RunRequest: &armcontainerregistry.EncodedTaskRunRequest{
				Type:                 to.Ptr("EncodedTaskRunRequest"),
				IsArchiveEnabled:     to.Ptr(true),
				Credentials:          &armcontainerregistry.Credentials{},
				EncodedTaskContent:   to.Ptr("c3RlcHM6IAogIC0gY21kOiB7eyAuVmFsdWVzLmNvbW1hbmQgfX0K"),
				EncodedValuesContent: to.Ptr("Y29tbWFuZDogYmFzaCBlY2hvIHt7LlJ1bi5SZWdpc3RyeX19Cg=="),
				Platform: &armcontainerregistry.PlatformProperties{
					Architecture: to.Ptr(armcontainerregistry.ArchitectureAmd64),
					OS:           to.Ptr(armcontainerregistry.OSLinux),
				},
				Values: []*armcontainerregistry.SetValue{},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, taskRunsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TaskRuns_GetDetails
	fmt.Println("Call operation: TaskRuns_GetDetails")
	_, err = taskRunsClient.GetDetails(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskRunName, nil)
	testsuite.Require().NoError(err)

	// From step Registries_GetBuildSourceUploadUrl
	fmt.Println("Call operation: Registries_GetBuildSourceUploadUrl")
	registriesClient, err := armcontainerregistry.NewRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = registriesClient.GetBuildSourceUploadURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, nil)
	testsuite.Require().NoError(err)

	// From step Runs_List
	fmt.Println("Call operation: Runs_List")
	runsClient, err := armcontainerregistry.NewRunsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	runsClientNewListPager := runsClient.NewListPager(testsuite.resourceGroupName, testsuite.registryName, &armcontainerregistry.RunsClientListOptions{Filter: nil,
		Top: to.Ptr[int32](10),
	})
	for runsClientNewListPager.More() {
		_, err := runsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Runs_Get
	fmt.Println("Call operation: Runs_Get")
	_, err = runsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, runId, nil)
	testsuite.Require().NoError(err)

	// From step Runs_Cancel
	fmt.Println("Call operation: Runs_Cancel")
	runsClientCancelResponsePoller, err := runsClient.BeginCancel(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, runId, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, runsClientCancelResponsePoller)
	testsuite.Require().NoError(err)

	// From step Runs_GetLogSasUrl
	fmt.Println("Call operation: Runs_GetLogSasUrl")
	_, err = runsClient.GetLogSasURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, runId, nil)
	testsuite.Require().NoError(err)

	// From step Registries_ScheduleRun
	fmt.Println("Call operation: Registries_ScheduleRun")
	registriesClientScheduleRunResponsePoller, err := registriesClient.BeginScheduleRun(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, &armcontainerregistry.DockerBuildRequest{
		Type:             to.Ptr("DockerBuildRequest"),
		IsArchiveEnabled: to.Ptr(true),
		AgentConfiguration: &armcontainerregistry.AgentProperties{
			CPU: to.Ptr[int32](2),
		},
		Arguments: []*armcontainerregistry.Argument{
			{
				Name:     to.Ptr("mytestargument"),
				IsSecret: to.Ptr(false),
				Value:    to.Ptr("mytestvalue"),
			},
			{
				Name:     to.Ptr("mysecrettestargument"),
				IsSecret: to.Ptr(true),
				Value:    to.Ptr("mysecrettestvalue"),
			}},
		DockerFilePath: to.Ptr("DockerFile"),
		ImageNames: []*string{
			to.Ptr("azurerest:testtag")},
		IsPushEnabled: to.Ptr(true),
		NoCache:       to.Ptr(true),
		Platform: &armcontainerregistry.PlatformProperties{
			Architecture: to.Ptr(armcontainerregistry.ArchitectureAmd64),
			OS:           to.Ptr(armcontainerregistry.OSLinux),
		},
		SourceLocation: to.Ptr("https://myaccount.blob.core.windows.net/sascontainer/source.zip?sv=2015-04-05&st=2015-04-29T22%3A18%3A26Z&se=2015-04-30T02%3A23%3A26Z&sr=b&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https&sig=Z%2FRHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk%3D"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, registriesClientScheduleRunResponsePoller)
	testsuite.Require().NoError(err)

	// From step TaskRuns_Delete
	fmt.Println("Call operation: TaskRuns_Delete")
	taskRunsClientDeleteResponsePoller, err := taskRunsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskRunName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, taskRunsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Runs_Update
	fmt.Println("Call operation: Runs_Update")
	runsClientUpdateResponsePoller, err := runsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, runId, armcontainerregistry.RunUpdateParameters{
		IsArchiveEnabled: to.Ptr(true),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, runsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Tasks_Delete
	fmt.Println("Call operation: Tasks_Delete")
	tasksClientDeleteResponsePoller, err := tasksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.registryName, testsuite.taskName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tasksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
