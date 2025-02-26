// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azbatch_test

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/batch/azbatch"
)

var client *azbatch.Client

// A pool is a collection of compute nodes (virtual machines) that run portions of your application's workload.
// A node's size determines the number of CPU cores, memory capacity, and local storage allocated to the node. See
// [Nodes and pools in Azure Batch] for more information.
//
// A job is a collection of tasks that manages how those tasks run on a pool's nodes. A task runs one or more programs
// or scripts on a node. New tasks are immediately assigned to a node for execution or queued until the pool has an
// available node. See [Jobs and tasks in Azure Batch] for more information.
//
// [Nodes and pools in Azure Batch]: https://learn.microsoft.com/azure/batch/nodes-and-pools
// [Jobs and tasks in Azure Batch]: https://learn.microsoft.com/azure/batch/jobs-and-tasks
func Example_package() {
	// create a pool with two dedicated nodes
	poolID := "HelloWorldPool"
	content := azbatch.CreatePoolContent{
		ID:                   to.Ptr(poolID),
		TargetDedicatedNodes: to.Ptr(int32(2)),
		VirtualMachineConfiguration: &azbatch.VirtualMachineConfiguration{
			DataDisks: []*azbatch.DataDisk{
				{
					DiskSizeGB:        to.Ptr(int32(1)),
					LogicalUnitNumber: to.Ptr(int32(1)),
				},
			},
			ImageReference: &azbatch.ImageReference{
				Offer:     to.Ptr("0001-com-ubuntu-server-jammy"),
				Publisher: to.Ptr("canonical"),
				SKU:       to.Ptr("22_04-lts"),
			},
			NodeAgentSKUID: to.Ptr("batch.node.ubuntu 22.04"),
		},
		VMSize: to.Ptr("Standard_A1_v2"),
	}
	_, err := client.CreatePool(context.TODO(), content, nil)
	if err != nil {
		// TODO: handle error
	}

	// create a job to run tasks in the pool
	jobID := "HelloWorldJob"
	jobContent := azbatch.CreateJobContent{
		ID: to.Ptr(jobID),
		PoolInfo: &azbatch.PoolInfo{
			PoolID: to.Ptr("HelloWorldPool"),
		},
	}
	_, err = client.CreateJob(context.TODO(), jobContent, nil)
	if err != nil {
		// TODO: handle error
	}
	// create a task to run as soon as the pool has an available node
	taskContent := azbatch.CreateTaskContent{
		ID:          to.Ptr("HelloWorldTask"),
		CommandLine: to.Ptr("echo Hello, world!"),
	}
	_, err = client.CreateTask(context.TODO(), jobID, taskContent, nil)
	if err != nil {
		// TODO: handle error
	}
}

// Each task has a working directory under which it can create files and directories. A task can use this directory to
// store the program run by the task, the data it processes, and its output. A task's files and directories are owned by
// the task's user.
//
// A portion of the node's file system is available to all tasks running on that node as a root directory located on the
// node's temporary storage drive. Tasks can access this root directory by referencing the AZ_BATCH_NODE_ROOT_DIR
// environment variable. For more information see [Files and directories in Azure Batch].
//
// [Files and directories in Azure Batch]: https://learn.microsoft.com/azure/batch/files-and-directories
func Example_taskOutputFile() {
	completedTasks := client.NewListTasksPager("TODO: job ID", &azbatch.ListTasksOptions{
		Filter: to.Ptr(fmt.Sprintf("state eq '%s'", azbatch.TaskStateCompleted)),
	})
	for completedTasks.More() {
		page, err := completedTasks.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, task := range page.Value {
			file := "stdout.txt"
			if *task.ExecutionInfo.ExitCode != 0 {
				file = "stderr.txt"
			}
			fc, err := client.GetTaskFile(context.TODO(), "TODO: job ID", *task.ID, file, nil)
			if err != nil {
				// TODO: handle error
			}
			fmt.Println(io.ReadAll(fc.Body))
		}
	}
}

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	client, err = azbatch.NewClient("https://TODO.batch.azure.com", cred, nil)
	if err != nil {
		// TODO: handle error
	}
	_ = client
}

func ExampleClient_NewListJobsPager() {
	for jobs := client.NewListJobsPager(nil); jobs.More(); {
		page, err := jobs.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, job := range page.Value {
			fmt.Println(*job.ID)
		}
	}
}

func ExampleClient_NewListJobSchedulesPager() {
	for schedules := client.NewListJobSchedulesPager(nil); schedules.More(); {
		page, err := schedules.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, schedule := range page.Value {
			fmt.Println(*schedule.ID)
		}
	}
}

func ExampleClient_NewListJobsFromSchedulePager() {
	for jobs := client.NewListJobsFromSchedulePager("TODO: schedule ID", nil); jobs.More(); {
		page, err := jobs.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, job := range page.Value {
			fmt.Println(*job.ID)
		}
	}
}

func ExampleClient_NewListNodeExtensionsPager() {
	for extensions := client.NewListNodeExtensionsPager("TODO: pool ID", "TODO: node ID", nil); extensions.More(); {
		page, err := extensions.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, extension := range page.Value {
			fmt.Println(*extension.VMExtension.Name)
		}
	}
}

func ExampleClient_NewListNodeFilesPager() {
	for files := client.NewListNodeFilesPager("TODO: pool ID", "TODO: node ID", nil); files.More(); {
		page, err := files.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, file := range page.Value {
			fmt.Println(*file.Name)
		}
	}
}

func ExampleClient_NewListNodesPager() {
	for nodes := client.NewListNodesPager("TODO: pool ID", nil); nodes.More(); {
		page, err := nodes.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, node := range page.Value {
			fmt.Println(*node.ID)
		}
	}
}

func ExampleClient_NewListPoolNodeCountsPager() {
	for counts := client.NewListPoolNodeCountsPager(nil); counts.More(); {
		page, err := counts.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, count := range page.Value {
			fmt.Println(*count.Dedicated)
		}
	}
}

func ExampleClient_NewListPoolsPager() {
	for pools := client.NewListPoolsPager(nil); pools.More(); {
		page, err := pools.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, pool := range page.Value {
			fmt.Println(*pool.ID)
		}
	}
}

func ExampleClient_NewListSubTasksPager() {
	for subtasks := client.NewListSubTasksPager("TODO: job ID", "TODO: task ID", nil); subtasks.More(); {
		page, err := subtasks.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, subtask := range page.Value {
			fmt.Println(*subtask.State)
		}
	}
}

func ExampleClient_NewListSupportedImagesPager() {
	for images := client.NewListSupportedImagesPager(nil); images.More(); {
		page, err := images.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, image := range page.Value {
			fmt.Println(*image.ImageReference.Offer)
		}
	}
}

func ExampleClient_NewListTaskFilesPager() {
	for files := client.NewListTaskFilesPager("TODO: job ID", "TODO: task ID", nil); files.More(); {
		page, err := files.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, file := range page.Value {
			fmt.Println(*file.Name)
		}
	}
}

func ExampleClient_NewListTasksPager() {
	for tasks := client.NewListTasksPager("TODO: job ID", nil); tasks.More(); {
		page, err := tasks.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, task := range page.Value {
			fmt.Println(*task.ID)
		}
	}
}
