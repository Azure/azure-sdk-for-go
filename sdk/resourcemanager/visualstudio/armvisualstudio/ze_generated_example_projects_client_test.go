//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armvisualstudio_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/visualstudio/armvisualstudio"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/visualstudio/resource-manager/Microsoft.VisualStudio/preview/2014-04-01-preview/examples/GetProjectResources_List.json
func ExampleProjectsClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armvisualstudio.NewProjectsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.ListByResourceGroup(ctx,
		"<resource-group-name>",
		"<root-resource-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/visualstudio/resource-manager/Microsoft.VisualStudio/preview/2014-04-01-preview/examples/CreateProjectResource.json
func ExampleProjectsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armvisualstudio.NewProjectsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	poller, err := client.BeginCreate(ctx,
		"<resource-group-name>",
		"<root-resource-name>",
		"<resource-name>",
		armvisualstudio.ProjectResource{
			Name:     to.Ptr("<name>"),
			Type:     to.Ptr("<type>"),
			ID:       to.Ptr("<id>"),
			Location: to.Ptr("<location>"),
			Tags:     map[string]*string{},
			Properties: map[string]*string{
				"ProcessTemplateId":    to.Ptr("6B724908-EF14-45CF-84F8-768B5384DA45"),
				"VersionControlOption": to.Ptr("Git"),
			},
		},
		&armvisualstudio.ProjectsClientBeginCreateOptions{Validating: nil,
			ResumeToken: "",
		})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/visualstudio/resource-manager/Microsoft.VisualStudio/preview/2014-04-01-preview/examples/GetProjectResource.json
func ExampleProjectsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armvisualstudio.NewProjectsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<root-resource-name>",
		"<resource-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/visualstudio/resource-manager/Microsoft.VisualStudio/preview/2014-04-01-preview/examples/UpdateProjectResource.json
func ExampleProjectsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armvisualstudio.NewProjectsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<root-resource-name>",
		"<resource-name>",
		armvisualstudio.ProjectResource{
			Name:       to.Ptr("<name>"),
			Type:       to.Ptr("<type>"),
			ID:         to.Ptr("<id>"),
			Location:   to.Ptr("<location>"),
			Tags:       map[string]*string{},
			Properties: map[string]*string{},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/visualstudio/resource-manager/Microsoft.VisualStudio/preview/2014-04-01-preview/examples/GetProjectJobStatus.json
func ExampleProjectsClient_GetJobStatus() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armvisualstudio.NewProjectsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.GetJobStatus(ctx,
		"<resource-group-name>",
		"<root-resource-name>",
		"<resource-name>",
		"<sub-container-name>",
		"<operation>",
		&armvisualstudio.ProjectsClientGetJobStatusOptions{JobID: to.Ptr("<job-id>")})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}
