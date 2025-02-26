// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azbatch_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/batch/azbatch"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestApplications(t *testing.T) {
	t.Parallel()
	client := record(t)
	for apps := client.NewListApplicationsPager(nil); apps.More(); {
		page, err := apps.NextPage(ctx)
		require.NoError(t, err)
		require.NotNil(t, page)
		for _, app := range page.Value {
			require.NotNil(t, app)
			require.NotNil(t, app.ID)
			ga, err := client.GetApplication(ctx, *app.ID, nil)
			require.NoError(t, err)
			require.NotNil(t, ga)
		}
	}
}

func TestDeallocateNode(t *testing.T) {
	t.Parallel()
	client, poolID := createDefaultPool(t)
	node := firstReadyNode(t, client, poolID)
	dn, err := client.DeallocateNode(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, dn)

	_, err = poll(
		func() azbatch.Node {
			gn, err := client.GetNode(ctx, poolID, *node.ID, nil)
			require.NoError(t, err)
			return gn.Node
		},
		func(n azbatch.Node) bool {
			return n.State != nil && *n.State == azbatch.NodeStateDeallocated
		},
		7*time.Minute,
	)
	require.NoError(t, err)

	sn, err := client.StartNode(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, sn)
}

func TestJob(t *testing.T) {
	t.Parallel()

	client, poolID := createDefaultPool(t)

	jid := randomString(t)
	cj, err := client.CreateJob(ctx, azbatch.CreateJobContent{
		Constraints: &azbatch.JobConstraints{
			MaxWallClockTime: to.Ptr("PT1H"),
		},
		ID:                 to.Ptr(jid),
		JobPreparationTask: &azbatch.JobPreparationTask{CommandLine: to.Ptr("/bin/sh -c 'echo preparing'")},
		JobReleaseTask:     &azbatch.JobReleaseTask{CommandLine: to.Ptr("/bin/sh -c 'echo release'")},
		OnAllTasksComplete: to.Ptr(azbatch.OnAllTasksCompleteNoAction),
		PoolInfo:           &azbatch.PoolInfo{PoolID: to.Ptr(poolID)},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, cj)
	t.Cleanup(func() {
		dj, err := client.DeleteJob(ctx, jid, nil)
		require.NoError(t, err)
		require.NotNil(t, dj)
	})

	disj, err := client.DisableJob(ctx, jid, azbatch.DisableJobContent{
		DisableTasks: to.Ptr(azbatch.DisableJobOptionWait),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, disj)

	ej, err := client.EnableJob(ctx, jid, nil)
	require.NoError(t, err)
	require.NotNil(t, ej)

	gj, err := client.GetJob(ctx, jid, nil)
	require.NoError(t, err)
	require.NotNil(t, gj)

	uj, err := client.UpdateJob(ctx, jid, azbatch.UpdateJobContent{
		Constraints: &azbatch.JobConstraints{
			MaxWallClockTime: to.Ptr("PT2H"),
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, uj)

	for pgr := client.NewListJobsPager(nil); pgr.More(); {
		_, err := pgr.NextPage(ctx)
		require.NoError(t, err)
	}

	rj, err := client.ReplaceJob(ctx, jid, azbatch.Job{
		ID:       to.Ptr(jid + "2"),
		PoolInfo: &azbatch.PoolInfo{PoolID: to.Ptr(poolID)},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, rj)

	for status := client.NewListJobPreparationAndReleaseTaskStatusPager(jid, nil); status.More(); {
		_, err := status.NextPage(ctx)
		require.NoError(t, err)
	}

	tj, err := client.TerminateJob(ctx, jid, nil)
	require.NoError(t, err)
	require.NotNil(t, tj)

}

func TestJobSchedule(t *testing.T) {
	t.Parallel()

	client, poolID := createDefaultPool(t)

	id := randomString(t)
	schedule := azbatch.CreateJobScheduleContent{
		DisplayName: to.Ptr(id),
		ID:          to.Ptr(id),
		JobSpecification: &azbatch.JobSpecification{
			PoolInfo: &azbatch.PoolInfo{PoolID: to.Ptr(poolID)},
		},
		Metadata: []*azbatch.MetadataItem{
			{
				Name:  to.Ptr("key"),
				Value: to.Ptr("value"),
			},
		},
		Schedule: &azbatch.JobScheduleConfiguration{
			RecurrenceInterval: to.Ptr("PT1H"),
		},
	}
	cj, err := client.CreateJobSchedule(ctx, schedule, nil)
	require.NoError(t, err)
	require.NotNil(t, cj)

	rj, err := client.ReplaceJobSchedule(ctx, id, azbatch.JobSchedule{
		ID: to.Ptr(id + "2"),
		JobSpecification: &azbatch.JobSpecification{
			PoolInfo: &azbatch.PoolInfo{PoolID: to.Ptr(poolID)},
		},
		Schedule: &azbatch.JobScheduleConfiguration{
			RecurrenceInterval: to.Ptr("PT2H"),
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, rj)

	gj, err := client.GetJobSchedule(ctx, *schedule.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, gj)

	uj, err := client.UpdateJobSchedule(ctx, *schedule.ID, azbatch.UpdateJobScheduleContent{
		Metadata: []*azbatch.MetadataItem{
			{
				Name:  to.Ptr("key"),
				Value: to.Ptr("value"),
			},
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, uj)

	ex, err := client.JobScheduleExists(ctx, *schedule.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, ex)

	for scheds := client.NewListJobSchedulesPager(nil); scheds.More(); {
		_, err := scheds.NextPage(ctx)
		require.NoError(t, err)
	}

	for jobs := client.NewListJobsFromSchedulePager(*schedule.ID, nil); jobs.More(); {
		_, err := jobs.NextPage(ctx)
		require.NoError(t, err)
	}

	disj, err := client.DisableJobSchedule(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, disj)

	ej, err := client.EnableJobSchedule(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, ej)

	tj, err := client.TerminateJobSchedule(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, tj)

	dj, err := client.DeleteJobSchedule(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, dj)
}

func TestListSupportedImages(t *testing.T) {
	t.Parallel()
	client := record(t)
	for images := client.NewListSupportedImagesPager(nil); images.More(); {
		page, err := images.NextPage(ctx)
		require.NoError(t, err)
		require.NotNil(t, page)
	}
}

func TestNode(t *testing.T) {
	t.Parallel()

	client := record(t)
	pool := defaultPoolContent(t)
	pool.NetworkConfiguration = &azbatch.NetworkConfiguration{
		EndpointConfiguration: &azbatch.PoolEndpointConfiguration{
			InboundNATPools: []*azbatch.InboundNATPool{
				{
					BackendPort:            to.Ptr(int32(22)),
					FrontendPortRangeStart: to.Ptr(int32(1)),
					FrontendPortRangeEnd:   to.Ptr(int32(42)),
					Name:                   to.Ptr("ssh"),
					NetworkSecurityGroupRules: []*azbatch.NetworkSecurityGroupRule{
						{
							Access:              to.Ptr(azbatch.NetworkSecurityGroupRuleAccessDeny),
							Priority:            to.Ptr(int32(150)),
							SourceAddressPrefix: to.Ptr("Internet"),
						},
					},
					Protocol: to.Ptr(azbatch.InboundEndpointProtocolTCP),
				},
			},
		},
	}
	poolID := *pool.ID
	_, err := client.CreatePool(ctx, pool, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = client.DeletePool(ctx, poolID, nil) })

	node := firstReadyNode(t, client, poolID)

	ga, err := client.GetNode(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, ga)

	for counts := client.NewListPoolNodeCountsPager(nil); counts.More(); {
		page, err := counts.NextPage(ctx)
		require.NoError(t, err)
		require.NotNil(t, page)
	}

	rl, err := client.GetNodeRemoteLoginSettings(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, rl)

	// TODO: InstanceViewStatus.Level is defined as a string enum but Batch returns a number
	// for exts := client.NewListNodeExtensionsPager(poolID, *node.ID, nil); exts.More(); {
	// 	page, err := exts.NextPage(ctx)
	// 	require.NotNil(t, page)
	// 	require.NoError(t, err)
	// 	for _, ext := range page.Value {
	// 		require.NotNil(t, ext)
	// 		require.NotNil(t, ext.VMExtension)
	// 		require.NotNil(t, ext.VMExtension.Name)
	// 		ge, err := client.GetNodeExtension(ctx, poolID, *node.ID, *ext.VMExtension.Name, nil)
	// 		require.NoError(t, err)
	// 		require.NotNil(t, ge)
	// 	}
	// }

	sn, err := client.DisableNodeScheduling(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, sn)

	en, err := client.EnableNodeScheduling(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, en)

	ul, err := client.UploadNodeLogs(ctx, poolID, *node.ID, azbatch.UploadNodeLogsContent{
		ContainerURL: to.Ptr("http://localhost"),
		StartTime:    to.Ptr(time.Now().Add(-time.Minute)),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, ul)

	cu, err := client.CreateNodeUser(ctx, poolID, *node.ID, azbatch.CreateNodeUserContent{
		Name:     to.Ptr("username"),
		Password: to.Ptr("password"),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, cu)

	ru, err := client.ReplaceNodeUser(ctx, poolID, *node.ID, "username", azbatch.UpdateNodeUserContent{
		Password: to.Ptr("password2"),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, ru)

	du, err := client.DeleteNodeUser(ctx, poolID, *node.ID, "username", nil)
	require.NoError(t, err)
	require.NotNil(t, du)

	rm, err := client.RemoveNodes(ctx, poolID, azbatch.RemoveNodeContent{
		NodeList: []*string{node.ID},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, rm)
}

func TestNodeFiles(t *testing.T) {
	t.Parallel()

	client, poolID := createDefaultPool(t)
	node := firstReadyNode(t, client, poolID)
	jid := randomString(t)
	cj, err := client.CreateJob(ctx, azbatch.CreateJobContent{
		Constraints: &azbatch.JobConstraints{
			MaxWallClockTime: to.Ptr("PT1H"),
		},
		ID:                 to.Ptr(jid),
		OnAllTasksComplete: to.Ptr(azbatch.OnAllTasksCompleteTerminateJob),
		PoolInfo:           &azbatch.PoolInfo{PoolID: &poolID},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, cj)

	_ = waitForTask(t, client, jid, "/bin/sh -c 'echo done > $AZ_BATCH_NODE_SHARED_DIR/test.txt'")

	var file *azbatch.NodeFile
	files := client.NewListNodeFilesPager(poolID, *node.ID, &azbatch.ListNodeFilesOptions{Recursive: to.Ptr(true)})
	for files.More() {
		p, err := files.NextPage(ctx)
		require.NoError(t, err)
		for _, f := range p.Value {
			if f != nil && f.Name != nil && strings.HasSuffix(*f.Name, "test.txt") {
				file = f
				break
			}
		}
	}
	require.NotNil(t, file, "didn't find test file")

	gf, err := client.GetNodeFile(ctx, poolID, *node.ID, *file.Name, nil)
	require.NoError(t, err)
	require.NotNil(t, gf)

	fp, err := client.GetNodeFileProperties(ctx, poolID, *node.ID, *file.Name, nil)
	require.NoError(t, err)
	require.NotNil(t, fp)

	df, err := client.DeleteNodeFile(ctx, poolID, *node.ID, *file.Name, nil)
	require.NoError(t, err)
	require.NotNil(t, df)
}

func TestPool(t *testing.T) {
	t.Parallel()
	client := record(t)

	pool := defaultPoolContent(t)
	// this test doesn't require a node so, don't allocate one
	pool.TargetDedicatedNodes = to.Ptr(int32(0))
	cp, err := client.CreatePool(ctx, pool, nil)
	require.NoError(t, err)
	require.NotNil(t, cp)
	t.Cleanup(func() {
		dr, err := client.DeletePool(ctx, *pool.ID, nil)
		require.NoError(t, err)
		require.NotNil(t, dr)
	})

	pe, err := client.PoolExists(ctx, *pool.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, pe)

	for nc := client.NewListPoolNodeCountsPager(nil); nc.More(); {
		_, err := nc.NextPage(ctx)
		require.NoError(t, err)
	}

	for pgr := client.NewListPoolsPager(nil); pgr.More(); {
		_, err := pgr.NextPage(ctx)
		require.NoError(t, err)
	}

	up, err := client.UpdatePool(ctx, *pool.ID, azbatch.UpdatePoolContent{
		Metadata: []*azbatch.MetadataItem{
			{
				Name:  to.Ptr("key"),
				Value: to.Ptr("value"),
			},
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, up)

	// TODO: "The missing parameters are: certificateReferences."
	// rpp, err := client.ReplacePoolProperties(ctx, *pool.ID, azbatch.ReplacePoolContent{
	// ApplicationPackageReferences: azcore.NullValue[[]*azbatch.ApplicationPackageReference](),
	// 	Metadata: []*azbatch.MetadataItem{
	// 		{
	// 			Name:  to.Ptr("key2"),
	// 			Value: to.Ptr("value2"),
	// 		},
	// 	},
	// }, nil)
	// require.NoError(t, err)
	// require.NotNil(t, rpp)

	_, err = poll(
		func() azbatch.Pool {
			gp, err := client.GetPool(ctx, *pool.ID, &azbatch.GetPoolOptions{SelectParam: []string{"allocationState"}})
			require.NoError(t, err)
			return gp.Pool
		},
		func(p azbatch.Pool) bool {
			return p.AllocationState != nil && *p.AllocationState == azbatch.AllocationStateSteady
		},
		5*time.Minute,
	)
	require.NoError(t, err)

	ar, err := client.EnablePoolAutoScale(ctx, *pool.ID, azbatch.EnablePoolAutoScaleContent{
		AutoScaleEvaluationInterval: to.Ptr("PT1H"),
		AutoScaleFormula:            to.Ptr("$TargetDedicatedNodes=0"),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, ar)

	eva, err := client.EvaluatePoolAutoScale(ctx, *pool.ID, azbatch.EvaluatePoolAutoScaleContent{
		AutoScaleFormula: to.Ptr("$TargetDedicatedNodes=1"),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, eva)

	dr, err := client.DisablePoolAutoScale(ctx, *pool.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, dr)

	rp, err := client.ResizePool(ctx, *pool.ID, azbatch.ResizePoolContent{
		NodeDeallocationOption: to.Ptr(azbatch.NodeDeallocationOptionRequeue),
		TargetDedicatedNodes:   to.Ptr(*pool.TargetDedicatedNodes + int32(1)),
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, rp)

	sr, err := client.StopPoolResize(ctx, *pool.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, sr)
}

func TestRebootNode(t *testing.T) {
	t.Parallel()
	client, poolID := createDefaultPool(t)
	node := firstReadyNode(t, client, poolID)
	rn, err := client.RebootNode(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, rn)
}

func TestReimageNode(t *testing.T) {
	t.Parallel()
	client, poolID := createDefaultPool(t)
	node := firstReadyNode(t, client, poolID)
	rn, err := client.ReimageNode(ctx, poolID, *node.ID, nil)
	require.NoError(t, err)
	require.NotNil(t, rn)
}

func TestSerDe(t *testing.T) {
	t.Parallel()
	for _, model := range []interface {
		json.Marshaler
		json.Unmarshaler
	}{
		&azbatch.AccountListSupportedImagesResult{},
		&azbatch.AddTaskCollectionResult{},
		&azbatch.AffinityInfo{},
		&azbatch.Application{},
		&azbatch.ApplicationListResult{},
		&azbatch.ApplicationPackageReference{},
		&azbatch.AuthenticationTokenSettings{},
		&azbatch.AutoPoolSpecification{},
		&azbatch.AutoScaleRun{},
		&azbatch.AutoScaleRunError{},
		&azbatch.AutoUserSpecification{},
		&azbatch.AutomaticOSUpgradePolicy{},
		&azbatch.AzureBlobFileSystemConfiguration{},
		&azbatch.AzureFileShareConfiguration{},
		&azbatch.CIFSMountConfiguration{},
		&azbatch.ContainerConfiguration{},
		&azbatch.ContainerHostBindMountEntry{},
		&azbatch.ContainerRegistryReference{},
		&azbatch.CreateJobContent{},
		&azbatch.CreateJobScheduleContent{},
		&azbatch.CreateNodeUserContent{},
		&azbatch.CreatePoolContent{},
		&azbatch.CreateTaskContent{},
		&azbatch.DataDisk{},
		&azbatch.DeallocateNodeContent{},
		&azbatch.DiffDiskSettings{},
		&azbatch.DisableJobContent{},
		&azbatch.DisableNodeSchedulingContent{},
		&azbatch.DiskEncryptionConfiguration{},
		&azbatch.EnablePoolAutoScaleContent{},
		&azbatch.EnvironmentSetting{},
		&azbatch.Error{},
		&azbatch.ErrorDetail{},
		&azbatch.ErrorMessage{},
		&azbatch.EvaluatePoolAutoScaleContent{},
		&azbatch.ExitCodeMapping{},
		&azbatch.ExitCodeRangeMapping{},
		&azbatch.ExitConditions{},
		&azbatch.ExitOptions{},
		&azbatch.FileProperties{},
		&azbatch.HTTPHeader{},
		&azbatch.ImageReference{},
		&azbatch.InboundEndpoint{},
		&azbatch.InboundNATPool{},
		&azbatch.InstanceViewStatus{},
		&azbatch.Job{},
		&azbatch.JobConstraints{},
		&azbatch.JobExecutionInfo{},
		&azbatch.JobListResult{},
		&azbatch.JobManagerTask{},
		&azbatch.JobNetworkConfiguration{},
		&azbatch.JobPreparationAndReleaseTaskStatus{},
		&azbatch.JobPreparationAndReleaseTaskStatusListResult{},
		&azbatch.JobPreparationTask{},
		&azbatch.JobPreparationTaskExecutionInfo{},
		&azbatch.JobReleaseTask{},
		&azbatch.JobReleaseTaskExecutionInfo{},
		&azbatch.JobSchedule{},
		&azbatch.JobScheduleConfiguration{},
		&azbatch.JobScheduleExecutionInfo{},
		&azbatch.JobScheduleListResult{},
		&azbatch.JobScheduleStatistics{},
		&azbatch.JobSchedulingError{},
		&azbatch.JobSpecification{},
		&azbatch.JobStatistics{},
		&azbatch.LinuxUserConfiguration{},
		&azbatch.ListPoolNodeCountsResult{},
		&azbatch.ManagedDisk{},
		&azbatch.MetadataItem{},
		&azbatch.MountConfiguration{},
		&azbatch.MultiInstanceSettings{},
		&azbatch.NFSMountConfiguration{},
		&azbatch.NameValuePair{},
		&azbatch.NetworkConfiguration{},
		&azbatch.NetworkSecurityGroupRule{},
		&azbatch.Node{},
		&azbatch.NodeAgentInfo{},
		&azbatch.NodeCounts{},
		&azbatch.NodeEndpointConfiguration{},
		&azbatch.NodeError{},
		&azbatch.NodeFile{},
		&azbatch.NodeFileListResult{},
		&azbatch.NodeIdentityReference{},
		&azbatch.NodeInfo{},
		&azbatch.NodeListResult{},
		&azbatch.NodePlacementConfiguration{},
		&azbatch.NodeRemoteLoginSettings{},
		&azbatch.NodeVMExtension{},
		&azbatch.NodeVMExtensionListResult{},
		&azbatch.OSDisk{},
		&azbatch.OutputFile{},
		&azbatch.OutputFileBlobContainerDestination{},
		&azbatch.OutputFileDestination{},
		&azbatch.OutputFileUploadConfig{},
		&azbatch.Pool{},
		&azbatch.PoolEndpointConfiguration{},
		&azbatch.PoolIdentity{},
		&azbatch.PoolInfo{},
		&azbatch.PoolListResult{},
		&azbatch.PoolNodeCounts{},
		&azbatch.PoolResourceStatistics{},
		&azbatch.PoolSpecification{},
		&azbatch.PoolStatistics{},
		&azbatch.PoolUsageStatistics{},
		&azbatch.PublicIPAddressConfiguration{},
		&azbatch.RebootNodeContent{},
		&azbatch.RecentJob{},
		&azbatch.ReimageNodeContent{},
		&azbatch.RemoveNodeContent{},
		&azbatch.ReplacePoolContent{},
		&azbatch.ResizeError{},
		&azbatch.ResizePoolContent{},
		&azbatch.ResourceFile{},
		&azbatch.RollingUpgradePolicy{},
		&azbatch.SecurityProfile{},
		&azbatch.ServiceArtifactReference{},
		&azbatch.StartTask{},
		&azbatch.StartTaskInfo{},
		&azbatch.Subtask{},
		&azbatch.SupportedImage{},
		&azbatch.Task{},
		&azbatch.TaskAddResult{},
		&azbatch.TaskConstraints{},
		&azbatch.TaskContainerExecutionInfo{},
		&azbatch.TaskContainerSettings{},
		&azbatch.TaskCounts{},
		&azbatch.TaskCountsResult{},
		&azbatch.TaskDependencies{},
		&azbatch.TaskExecutionInfo{},
		&azbatch.TaskFailureInfo{},
		&azbatch.TaskGroup{},
		&azbatch.TaskIDRange{},
		&azbatch.TaskInfo{},
		&azbatch.TaskListResult{},
		&azbatch.TaskListSubtasksResult{},
		&azbatch.TaskSchedulingPolicy{},
		&azbatch.TaskSlotCounts{},
		&azbatch.TaskStatistics{},
		&azbatch.TerminateJobContent{},
		&azbatch.UEFISettings{},
		&azbatch.UpdateJobContent{},
		&azbatch.UpdateJobScheduleContent{},
		&azbatch.UpdateNodeUserContent{},
		&azbatch.UpdatePoolContent{},
		&azbatch.UpgradePolicy{},
		&azbatch.UploadNodeLogsContent{},
		&azbatch.UploadNodeLogsResult{},
		&azbatch.UserAccount{},
		&azbatch.UserAssignedIdentity{},
		&azbatch.UserIdentity{},
		&azbatch.VMDiskSecurityProfile{},
		&azbatch.VMExtension{},
		&azbatch.VMExtensionInstanceView{},
		&azbatch.VirtualMachineConfiguration{},
		&azbatch.VirtualMachineInfo{},
		&azbatch.WindowsConfiguration{},
		&azbatch.WindowsUserConfiguration{},
	} {
		require.Error(t, model.UnmarshalJSON([]byte{}))
		v := reflect.ValueOf(model).Elem()
		for i := 0; i < v.Type().NumField(); i++ {
			f := v.Field(i)
			switch f.Type().String() {
			case "*bool":
				f.Set(reflect.ValueOf(to.Ptr(true)))
			case "*float32":
				f.Set(reflect.ValueOf(to.Ptr(float32(1))))
			case "*int32":
				f.Set(reflect.ValueOf(to.Ptr(int32(1))))
			case "*int64":
				f.Set(reflect.ValueOf(to.Ptr(int64(1))))
			case "*string":
				f.Set(reflect.ValueOf(to.Ptr("...")))
			}
		}
		b, err := model.MarshalJSON()
		require.NoError(t, err)
		require.NoError(t, model.UnmarshalJSON(b))
	}
}

func TestTask(t *testing.T) {
	t.Parallel()

	client, poolID := createDefaultPool(t)
	jid := randomString(t)
	cj, err := client.CreateJob(ctx, azbatch.CreateJobContent{
		ID:                 to.Ptr(jid),
		OnAllTasksComplete: to.Ptr(azbatch.OnAllTasksCompleteTerminateJob),
		PoolInfo:           &azbatch.PoolInfo{PoolID: to.Ptr(poolID)},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, cj)

	t.Run("Replace", func(t *testing.T) {
		t.Parallel()
		client := record(t)

		tid := randomString(t)
		ct, err := client.CreateTask(ctx, jid, azbatch.CreateTaskContent{
			CommandLine: to.Ptr("/bin/sh -c 'sleep 300'"),
			ID:          to.Ptr(tid),
		}, nil)
		require.NoError(t, err)
		require.NotNil(t, ct)

		jtc, err := client.GetJobTaskCounts(ctx, jid, nil)
		require.NoError(t, err)
		require.NotNil(t, jtc)

		rt, err := client.ReplaceTask(ctx, jid, tid, azbatch.Task{
			Constraints: &azbatch.TaskConstraints{
				MaxTaskRetryCount: to.Ptr(int32(1)),
				MaxWallClockTime:  to.Ptr("PT1H"),
				RetentionTime:     to.Ptr("PT2H"),
			},
		}, nil)
		require.NoError(t, err)
		require.NotNil(t, rt)

		tt, err := client.TerminateTask(ctx, jid, tid, nil)
		require.NoError(t, err)
		require.NotNil(t, tt)

		ret, err := client.ReactivateTask(ctx, jid, tid, nil)
		require.NoError(t, err)
		require.NotNil(t, ret)
	})

	tid := randomString(t)
	ctc, err := client.CreateTaskCollection(ctx, jid, azbatch.TaskGroup{
		Value: []*azbatch.CreateTaskContent{
			{
				CommandLine: to.Ptr("/bin/sh -c 'echo done > $AZ_BATCH_TASK_DIR/task.txt'"),
				ID:          to.Ptr(tid),
			},
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, ctc)

	for pgr := client.NewListTasksPager(jid, nil); pgr.More(); {
		p, err := pgr.NextPage(ctx)
		require.NoError(t, err)
		for _, task := range p.Value {
			if task != nil && task.ID != nil && *task.ID == tid {
				t.Cleanup(func() {
					dt, err := client.DeleteTask(ctx, jid, *task.ID, nil)
					require.NoError(t, err)
					require.NotNil(t, dt)
				})
				break
			}
		}
	}

	for subtasks := client.NewListSubTasksPager(jid, tid, nil); subtasks.More(); {
		_, err := subtasks.NextPage(ctx)
		require.NoError(t, err)
	}

	_, err = poll(
		func() azbatch.Task {
			gt, err := client.GetTask(ctx, jid, tid, nil)
			require.NoError(t, err)
			return gt.Task
		},
		func(task azbatch.Task) bool {
			return task.State != nil && *task.State == azbatch.TaskStateCompleted
		},
		5*time.Minute,
	)
	require.NoError(t, err, "task isn't complete")

	files := client.NewListTaskFilesPager(jid, tid, &azbatch.ListTaskFilesOptions{
		Recursive: to.Ptr(true),
	})
	require.NotNil(t, files)

	var file *azbatch.NodeFile
	for files.More() {
		pg, err := files.NextPage(ctx)
		require.NoError(t, err)
		require.NotNil(t, pg)
		for _, file = range pg.Value {
			if file != nil && file.Name != nil && strings.HasSuffix(*file.Name, "task.txt") {
				props, err := client.GetTaskFileProperties(ctx, jid, tid, *file.Name, nil)
				require.NoError(t, err)
				require.NotNil(t, props)
				break
			}
		}
	}
	require.NotNil(t, file, "didn't find file created by task")

	gtf, err := client.GetTaskFile(ctx, jid, tid, *file.Name, nil)
	require.NoError(t, err)
	require.NotNil(t, gtf)

	dtf, err := client.DeleteTaskFile(ctx, jid, tid, *file.Name, nil)
	require.NoError(t, err)
	require.NotNil(t, dtf)
}
