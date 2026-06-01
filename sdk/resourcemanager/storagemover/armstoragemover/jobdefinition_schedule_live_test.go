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

// JobDefinitionScheduleScenarioSuite mirrors .NET JobDefinitionScheduleTests. All schedule dates use
// explicit UTC + Z suffix because the RP returns 500 for `+00:00`-suffixed timestamps in
// api-version 2025-12-01 (see cross-language playbook).
type JobDefinitionScheduleScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	projectName      string
	nfsEndpointName  string
	blobEndpointName string
	weeklyJobDef     string
	dailyJobDef      string
	onetimeJobDef    string
}

func TestJobDefinitionScheduleScenarioSuite(t *testing.T) {
	suite.Run(t, new(JobDefinitionScheduleScenarioSuite))
}

func (s *JobDefinitionScheduleScenarioSuite) SetupSuite() {
	s.setupBase()
	// Schedule dates are computed from time.Now() in scheduleStart/scheduleEnd to avoid hitting the
	// RP's "StartDate must not be in the past" validation. Sanitize them in the recorded body so
	// re-records and playback compare equal regardless of when they ran.
	for _, jp := range []string{"$..startDate", "$..endDate"} {
		if err := recording.AddBodyKeySanitizer(jp, "2099-01-01T00:00:00Z", "", &recording.RecordingOptions{UseHTTPS: true, TestInstance: s.T()}); err != nil {
			s.T().Logf("warning: failed to add body-key sanitizer for %s: %v", jp, err)
		}
	}
	s.storageMoverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
	s.nfsEndpointName = s.generateName("nfsep")
	s.blobEndpointName = s.generateName("blobep")
	s.weeklyJobDef = s.generateName("jobdefw")
	s.dailyJobDef = s.generateName("jobdefd")
	s.onetimeJobDef = s.generateName("jobdefo")
	s.createStorageMover(s.storageMoverName, nil, "")
	s.createProject(s.storageMoverName, s.projectName, "")
	s.createNfsEndpoint(s.storageMoverName, s.nfsEndpointName, "")
	s.createBlobEndpoint(s.storageMoverName, s.blobEndpointName, "")
}

func (s *JobDefinitionScheduleScenarioSuite) TearDownSuite() { s.teardownBase() }

// scheduleStart returns a UTC start date in the near future. The RP rejects schedules with a start
// date in the past, so we anchor to the current day plus one. Recording stability is preserved by a
// body-key sanitizer on `$..startDate` / `$..endDate` registered in SetupSuite.
func (s *JobDefinitionScheduleScenarioSuite) scheduleStart() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
}

// scheduleEnd returns a UTC end date one month after scheduleStart.
func (s *JobDefinitionScheduleScenarioSuite) scheduleEnd() time.Time {
	return s.scheduleStart().AddDate(0, 1, 0)
}

// TestJobDefinitionWeeklySchedule mirrors .NET CreateJobDefinitionWithWeeklyScheduleTest.
func (s *JobDefinitionScheduleScenarioSuite) TestJobDefinitionWeeklySchedule() {
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	startDate := s.scheduleStart()
	endDate := s.scheduleEnd()
	schedule := &armstoragemover.ScheduleInfo{
		Frequency:     to.Ptr(armstoragemover.FrequencyWeekly),
		IsActive:      to.Ptr(true),
		ExecutionTime: &armstoragemover.SchedulerTime{Hour: to.Ptr(int32(2))},
		StartDate:     &startDate,
		EndDate:       &endDate,
		DaysOfWeek:    []*string{to.Ptr("Monday"), to.Ptr("Wednesday"), to.Ptr("Friday")},
	}
	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.weeklyJobDef, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:                to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName:              to.Ptr(s.nfsEndpointName),
			TargetName:              to.Ptr(s.blobEndpointName),
			Description:             to.Ptr("Job definition with weekly schedule"),
			DataIntegrityValidation: to.Ptr(armstoragemover.DataIntegrityValidationSaveVerifyFileMD5),
			Schedule:                schedule,
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.weeklyJobDef, *createResp.Name)
	s.Equal(s.nfsEndpointName, *createResp.Properties.SourceName)
	s.Equal(s.blobEndpointName, *createResp.Properties.TargetName)
	s.Equal(armstoragemover.CopyModeAdditive, *createResp.Properties.CopyMode)
	s.Equal("Job definition with weekly schedule", *createResp.Properties.Description)

	s.Require().NotNil(createResp.Properties.Schedule)
	s.Equal(armstoragemover.FrequencyWeekly, *createResp.Properties.Schedule.Frequency)
	s.True(*createResp.Properties.Schedule.IsActive)
	s.Require().NotNil(createResp.Properties.Schedule.ExecutionTime)
	s.Equal(int32(2), *createResp.Properties.Schedule.ExecutionTime.Hour)
	s.Len(createResp.Properties.Schedule.DaysOfWeek, 3)

	// Verify persistence with a separate Get.
	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.weeklyJobDef, nil)
	s.Require().NoError(err)
	s.Equal(s.weeklyJobDef, *getResp.Name)
	s.Require().NotNil(getResp.Properties.Schedule)
	s.Equal(armstoragemover.FrequencyWeekly, *getResp.Properties.Schedule.Frequency)

	// Clean up.
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.weeklyJobDef, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)

	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.weeklyJobDef, nil)
	s.expectResponseError(err)
}

// TestJobDefinitionDailyScheduleAndPreservePermissions mirrors .NET
// CreateJobDefinitionWithDailyScheduleAndPreservePermissionsTest.
func (s *JobDefinitionScheduleScenarioSuite) TestJobDefinitionDailyScheduleAndPreservePermissions() {
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	startDate := s.scheduleStart()
	endDate := s.scheduleEnd()
	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.dailyJobDef, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:                to.Ptr(armstoragemover.CopyModeMirror),
			SourceName:              to.Ptr(s.nfsEndpointName),
			TargetName:              to.Ptr(s.blobEndpointName),
			Description:             to.Ptr("Job definition with daily schedule"),
			DataIntegrityValidation: to.Ptr(armstoragemover.DataIntegrityValidationNone),
			PreservePermissions:     to.Ptr(true),
			Schedule: &armstoragemover.ScheduleInfo{
				Frequency:     to.Ptr(armstoragemover.FrequencyDaily),
				IsActive:      to.Ptr(true),
				ExecutionTime: &armstoragemover.SchedulerTime{Hour: to.Ptr(int32(0))},
				StartDate:     &startDate,
				EndDate:       &endDate,
			},
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.dailyJobDef, *createResp.Name)
	s.Equal(armstoragemover.CopyModeMirror, *createResp.Properties.CopyMode)
	s.Require().NotNil(createResp.Properties.Schedule)
	s.Equal(armstoragemover.FrequencyDaily, *createResp.Properties.Schedule.Frequency)
	s.True(*createResp.Properties.Schedule.IsActive)

	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.dailyJobDef, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)
}

// TestJobDefinitionOnetimeSchedule mirrors .NET CreateJobDefinitionWithOnetimeScheduleTest. Onetime
// frequency only requires StartDate; EndDate is omitted.
func (s *JobDefinitionScheduleScenarioSuite) TestJobDefinitionOnetimeSchedule() {
	client, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)

	startDate := s.scheduleStart()
	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.onetimeJobDef, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:    to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName:  to.Ptr(s.nfsEndpointName),
			TargetName:  to.Ptr(s.blobEndpointName),
			Description: to.Ptr("Job definition with one-time schedule"),
			Schedule: &armstoragemover.ScheduleInfo{
				Frequency:     to.Ptr(armstoragemover.FrequencyOnetime),
				IsActive:      to.Ptr(true),
				ExecutionTime: &armstoragemover.SchedulerTime{Hour: to.Ptr(int32(10))},
				StartDate:     &startDate,
			},
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.onetimeJobDef, *createResp.Name)
	s.Require().NotNil(createResp.Properties.Schedule)
	s.Equal(armstoragemover.FrequencyOnetime, *createResp.Properties.Schedule.Frequency)
	s.True(*createResp.Properties.Schedule.IsActive)

	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.onetimeJobDef, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)
}
