//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armvideoanalyzer_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/videoanalyzer/armvideoanalyzer"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-list.json
func ExamplePipelineJobsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager("testrg",
		"testaccount2",
		&armvideoanalyzer.PipelineJobsClientListOptions{Filter: nil,
			Top: to.Ptr[int32](2),
		})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-get-by-name.json
func ExamplePipelineJobsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"testrg",
		"testaccount2",
		"pipelineJob1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-create.json
func ExamplePipelineJobsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"testrg",
		"testaccount2",
		"pipelineJob1",
		armvideoanalyzer.PipelineJob{
			Properties: &armvideoanalyzer.PipelineJobProperties{
				Description: to.Ptr("Pipeline Job 1 Dsecription"),
				Parameters: []*armvideoanalyzer.ParameterDefinition{
					{
						Name:  to.Ptr("timesequences"),
						Value: to.Ptr("[[\"2020-10-05T03:30:00Z\", \"2020-10-05T04:30:00Z\"]]"),
					},
					{
						Name:  to.Ptr("videoSourceName"),
						Value: to.Ptr("camera001"),
					}},
				TopologyName: to.Ptr("pipelinetopology1"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-delete.json
func ExamplePipelineJobsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"testrg",
		"testaccount2",
		"pipelineJob1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-patch.json
func ExamplePipelineJobsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"testrg",
		"testaccount2",
		"pipelineJob1",
		armvideoanalyzer.PipelineJobUpdate{
			Properties: &armvideoanalyzer.PipelineJobPropertiesUpdate{
				Description: to.Ptr("Pipeline Job 1 description"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/videoanalyzer/resource-manager/Microsoft.Media/preview/2021-11-01-preview/examples/pipeline-job-cancel.json
func ExamplePipelineJobsClient_BeginCancel() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armvideoanalyzer.NewPipelineJobsClient("591e76c3-3e97-44db-879c-3e2b12961b62", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCancel(ctx,
		"testrg",
		"testaccount2",
		"pipelineJob1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
