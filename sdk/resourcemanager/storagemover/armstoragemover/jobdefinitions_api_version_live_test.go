// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

type JobDefinitionAPIVersionScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName      string
	projectName           string
	hourlySourceName      string
	hourlyTargetName      string
	hourlyJobName         string
	agentlessSourceName   string
	agentlessTargetName   string
	agentlessJobName      string
	connectionName        string
	crossTenantMoverName  string
	crossTenantSourceName string
}

func TestJobDefinitionAPIVersionScenarioSuite(t *testing.T) {
	suite.Run(t, new(JobDefinitionAPIVersionScenarioSuite))
}

func (s *JobDefinitionAPIVersionScenarioSuite) SetupSuite() {
	// Hourly schedules and agentless jobs are deployed in eastus2euap.
	s.setupBaseInLocation(scenarioEastUS2EUAPLocation)
	for _, path := range []string{"$..startDate", "$..endDate"} {
		err := recording.AddBodyKeySanitizer(path, "2099-01-01T00:00:00Z", "",
			&recording.RecordingOptions{UseHTTPS: true, TestInstance: s.T()})
		s.Require().NoError(err)
	}
	s.storageMoverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
	s.hourlySourceName = s.generateName("hourlysrc")
	s.hourlyTargetName = s.generateName("hourlytgt")
	s.hourlyJobName = s.generateName("hourlyjob")
	s.agentlessSourceName = s.generateName("agentlesssrc")
	s.agentlessTargetName = s.generateName("agentlesstgt")
	s.agentlessJobName = s.generateName("agentlessjob")
	s.connectionName = s.generateName("conn")
	s.crossTenantMoverName = s.generateName("xtmover")
	s.crossTenantSourceName = s.generateName("xtsource")
	s.createStorageMover(s.storageMoverName, nil, "")
	s.createProject(s.storageMoverName, s.projectName, "")
}

func (s *JobDefinitionAPIVersionScenarioSuite) TearDownSuite() { s.teardownBase() }

func (s *JobDefinitionAPIVersionScenarioSuite) createSmbFileShareEndpoint(name, shareName string, kind armstoragemover.EndpointKind) {
	client, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, name,
		armstoragemover.Endpoint{Properties: &armstoragemover.AzureStorageSmbFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(shareName),
			EndpointKind:             to.Ptr(kind),
		}}, nil)
	s.Require().NoError(err)
}

func (s *JobDefinitionAPIVersionScenarioSuite) TestHourlySchedule() {
	s.createSmbFileShareEndpoint(s.hourlySourceName, "srcshare", armstoragemover.EndpointKindSource)
	s.createSmbFileShareEndpoint(s.hourlyTargetName, "tgtshare", armstoragemover.EndpointKindTarget)
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	startDate := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
	endDate := startDate.AddDate(0, 1, 0)
	// File-share-to-file-share jobs must omit source and target subpaths.
	job := armstoragemover.JobDefinition{Properties: &armstoragemover.JobDefinitionProperties{
		CopyMode:    to.Ptr(armstoragemover.CopyModeAdditive),
		JobType:     to.Ptr(armstoragemover.JobTypeCloudToCloud),
		SourceName:  to.Ptr(s.hourlySourceName),
		TargetName:  to.Ptr(s.hourlyTargetName),
		Description: to.Ptr("Job definition with hourly schedule"),
		Schedule: &armstoragemover.ScheduleInfo{
			Frequency:      to.Ptr(armstoragemover.FrequencyHourly),
			RepeatInterval: to.Ptr("4h"),
			IsActive:       to.Ptr(true),
			ExecutionTime:  &armstoragemover.SchedulerTime{Hour: to.Ptr(int32(9)), Minute: to.Ptr(armstoragemover.MinuteZero)},
			StartDate:      &startDate,
			EndDate:        &endDate,
		},
	}}
	assertSchedule := func(definition armstoragemover.JobDefinition, interval string, active bool) {
		s.Require().NotNil(definition.Properties)
		s.Require().NotNil(definition.Properties.Schedule)
		s.Require().Equal(to.Ptr(armstoragemover.FrequencyHourly), definition.Properties.Schedule.Frequency)
		s.Require().Equal(to.Ptr(interval), definition.Properties.Schedule.RepeatInterval)
		s.Require().Equal(to.Ptr(active), definition.Properties.Schedule.IsActive)
	}

	createResponse, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName, job, nil)
	s.Require().NoError(err)
	s.Require().Equal(to.Ptr(s.hourlyJobName), createResponse.Name)
	assertSchedule(createResponse.JobDefinition, "4h", true)
	getResponse, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName, nil)
	s.Require().NoError(err)
	assertSchedule(getResponse.JobDefinition, "4h", true)

	updateResponse, err := client.Update(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName,
		armstoragemover.JobDefinitionUpdateParameters{Properties: &armstoragemover.JobDefinitionUpdateProperties{
			Schedule: &armstoragemover.ScheduleInfo{RepeatInterval: to.Ptr("6h")},
		}}, nil)
	s.Require().NoError(err)
	assertSchedule(updateResponse.JobDefinition, "6h", true)
	updateResponse, err = client.Update(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName,
		armstoragemover.JobDefinitionUpdateParameters{Properties: &armstoragemover.JobDefinitionUpdateProperties{
			Schedule: &armstoragemover.ScheduleInfo{IsActive: to.Ptr(false)},
		}}, nil)
	s.Require().NoError(err)
	assertSchedule(updateResponse.JobDefinition, "6h", false)

	// The repeat interval must evenly divide 24 hours.
	job.Properties.Schedule.RepeatInterval = to.Ptr("5h")
	_, err = client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName+"bad", job, nil)
	s.expectInvalidRequestPropertyValue(err)
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.hourlyJobName, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)
}

func (s *JobDefinitionAPIVersionScenarioSuite) TestCrossTenantJobValidation() {
	moversClient, err := armstoragemover.NewStorageMoversClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = moversClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.crossTenantMoverName,
		armstoragemover.StorageMover{Location: to.Ptr(scenarioWestCentralUSLocation)}, nil)
	s.Require().NoError(err)
	s.createProject(s.crossTenantMoverName, s.projectName, "")
	endpointsClient, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = endpointsClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.crossTenantMoverName, s.crossTenantSourceName,
		armstoragemover.Endpoint{Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			StorageAccountResourceID:  to.Ptr(s.placeholderStorageAccountResourceID()),
			BlobContainerName:         to.Ptr(scenarioContainerName),
			EndpointKind:              to.Ptr(armstoragemover.EndpointKindSource),
			EnableCrossTenantTransfer: to.Ptr(true),
			AllowedStorageAccounts:    []*string{to.Ptr(scenarioPartnerStorageAccountIDPrefix)},
		}}, nil)
	s.Require().NoError(err)
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	base := armstoragemover.JobDefinitionProperties{
		CopyMode:                    to.Ptr(armstoragemover.CopyModeAdditive),
		JobType:                     to.Ptr(armstoragemover.JobTypeCloudToCloud),
		SourceName:                  to.Ptr(s.crossTenantSourceName),
		TargetName:                  to.Ptr("partner-target"),
		IsCrossTenantJob:            to.Ptr(true),
		CrossTenantEndpointTenantID: to.Ptr(scenarioPartnerTenantID),
		CrossTenantEndpointResourceID: to.Ptr("/subscriptions/" + scenarioPartnerSubscriptionID +
			"/resourceGroups/partner-rg/providers/Microsoft.StorageMover/storageMovers/partner-sm/endpoints/partner-target"),
	}
	tests := []struct {
		name   string
		change func(*armstoragemover.JobDefinitionProperties)
	}{
		{"missing-tenant", func(properties *armstoragemover.JobDefinitionProperties) {
			properties.CrossTenantEndpointTenantID = nil
		}},
		{"missing-endpoint", func(properties *armstoragemover.JobDefinitionProperties) {
			properties.CrossTenantEndpointResourceID = nil
		}},
		{"missing-flag", func(properties *armstoragemover.JobDefinitionProperties) { properties.IsCrossTenantJob = nil }},
		{"not-cloud-to-cloud", func(properties *armstoragemover.JobDefinitionProperties) {
			properties.JobType = to.Ptr(armstoragemover.JobTypeOnPremToCloud)
		}},
		{"same-subscription", func(properties *armstoragemover.JobDefinitionProperties) {
			properties.CrossTenantEndpointResourceID = to.Ptr("/subscriptions/" + s.subscriptionID + "/resourceGroups/" +
				s.resourceGroupName + "/providers/Microsoft.StorageMover/storageMovers/" + s.crossTenantMoverName + "/endpoints/partner-target")
		}},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			properties := base
			tt.change(&properties)
			_, createErr := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.crossTenantMoverName, s.projectName,
				"job-"+tt.name, armstoragemover.JobDefinition{Properties: &properties}, nil)
			s.expectInvalidRequestPropertyValue(createErr)
		})
	}
	pager := client.NewListPager(s.resourceGroupName, s.crossTenantMoverName, s.projectName, nil)
	for pager.More() {
		page, pageErr := pager.NextPage(s.ctx)
		s.Require().NoError(pageErr)
		s.Require().Empty(page.Value, "rejected requests must not create job definitions")
	}
}

func (s *JobDefinitionAPIVersionScenarioSuite) TestCreateAgentlessJobDefinition() {
	endpointsClient, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = endpointsClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.agentlessSourceName,
		armstoragemover.Endpoint{Properties: &armstoragemover.SmbMountEndpointProperties{
			Host:      to.Ptr(scenarioSmbHost),
			ShareName: to.Ptr(scenarioSmbShareName),
			Credentials: &armstoragemover.AzureKeyVaultSmbCredentials{
				UsernameURI: to.Ptr(scenarioKeyVaultUsernameURI),
				PasswordURI: to.Ptr(scenarioKeyVaultPasswordURI),
			},
		}}, nil)
	s.Require().NoError(err)
	s.createSmbFileShareEndpoint(s.agentlessTargetName, scenarioFileShareName, armstoragemover.EndpointKindTarget)
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	job := armstoragemover.JobDefinition{Properties: &armstoragemover.JobDefinitionProperties{
		CopyMode:      to.Ptr(armstoragemover.CopyModeAdditive),
		JobType:       to.Ptr(armstoragemover.JobTypeOnPremToCloudAgentLess),
		SourceName:    to.Ptr(s.agentlessSourceName),
		TargetName:    to.Ptr(s.agentlessTargetName),
		SourceSubpath: to.Ptr("/"),
		TargetSubpath: to.Ptr("/"),
	}}
	_, err = client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.agentlessJobName, job, nil)
	s.expectInvalidRequestPropertyValue(err)

	connectionsClient, err := armstoragemover.NewConnectionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	connection, err := connectionsClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName,
		armstoragemover.Connection{Properties: &armstoragemover.ConnectionProperties{
			PrivateLinkServiceID: to.Ptr(scenarioPrivateLinkServiceIDPrefix + "E2E-VPC-PLS-E2ETest.Test.ECY"),
			Description:          to.Ptr("Connection for agentless job"),
		}}, nil)
	s.Require().NoError(err)
	defer func() {
		poller, deleteErr := connectionsClient.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
		s.Require().NoError(deleteErr)
		_, deleteErr = testutil.PollForTest(s.ctx, poller)
		s.Require().NoError(deleteErr)
	}()
	s.Require().NotNil(connection.ID)
	s.Require().NotNil(connection.Properties)
	s.Require().NotNil(connection.Properties.ConnectionStatus)
	s.Require().Contains([]armstoragemover.ConnectionStatus{
		armstoragemover.ConnectionStatusPending, armstoragemover.ConnectionStatusApproved,
	}, *connection.Properties.ConnectionStatus)

	job.Properties.Connections = []*string{connection.ID}
	createResponse, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.agentlessJobName, job, nil)
	s.Require().NoError(err)
	s.Require().Equal(to.Ptr(s.agentlessJobName), createResponse.Name)
	s.Require().NotNil(createResponse.Properties)
	s.Require().Equal(to.Ptr(armstoragemover.JobTypeOnPremToCloudAgentLess), createResponse.Properties.JobType)
	s.Require().Equal(job.Properties.Connections, createResponse.Properties.Connections)
	if createResponse.Properties.AgentName != nil {
		s.Require().Empty(*createResponse.Properties.AgentName)
	}
	getResponse, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.agentlessJobName, nil)
	s.Require().NoError(err)
	s.Require().NotNil(getResponse.Properties)
	s.Require().Equal(to.Ptr(armstoragemover.JobTypeOnPremToCloudAgentLess), getResponse.Properties.JobType)
	s.Require().Equal(job.Properties.Connections, getResponse.Properties.Connections)
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.agentlessJobName, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)
}
