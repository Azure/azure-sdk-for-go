// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery/fake"
)

// ExampleWorkspacesServer demonstrates how to use the fake server for testing.
// The fake server allows you to test your code without making actual API calls to Azure.
func ExampleWorkspacesServer() {
	// Create an instance of the fake server for the WorkspacesClient.
	fakeServer := fake.WorkspacesServer{
		// Implement the Get method to return a fake workspace
		Get: func(ctx context.Context, resourceGroupName string, workspaceName string, options *armdiscovery.WorkspacesClientGetOptions) (resp azfake.Responder[armdiscovery.WorkspacesClientGetResponse], errResp azfake.ErrorResponder) {
			// Return a fake workspace response
			resp.SetResponse(http.StatusOK, armdiscovery.WorkspacesClientGetResponse{
				Workspace: armdiscovery.Workspace{
					ID:       to.Ptr("/subscriptions/sub-id/resourceGroups/rg-name/providers/Microsoft.Discovery/workspaces/test-workspace"),
					Name:     to.Ptr(workspaceName),
					Type:     to.Ptr("Microsoft.Discovery/workspaces"),
					Location: to.Ptr("eastus"),
					Properties: &armdiscovery.WorkspaceProperties{
						ProvisioningState: to.Ptr(armdiscovery.ProvisioningStateSucceeded),
					},
				},
			}, nil)
			return
		},

		// Implement the list pager to return fake workspaces
		NewListByResourceGroupPager: func(resourceGroupName string, options *armdiscovery.WorkspacesClientListByResourceGroupOptions) (resp azfake.PagerResponder[armdiscovery.WorkspacesClientListByResourceGroupResponse]) {
			resp.AddPage(http.StatusOK, armdiscovery.WorkspacesClientListByResourceGroupResponse{
				WorkspaceListResult: armdiscovery.WorkspaceListResult{
					Value: []*armdiscovery.Workspace{
						{
							ID:       to.Ptr("/subscriptions/sub-id/resourceGroups/rg-name/providers/Microsoft.Discovery/workspaces/workspace-1"),
							Name:     to.Ptr("workspace-1"),
							Location: to.Ptr("eastus"),
						},
						{
							ID:       to.Ptr("/subscriptions/sub-id/resourceGroups/rg-name/providers/Microsoft.Discovery/workspaces/workspace-2"),
							Name:     to.Ptr("workspace-2"),
							Location: to.Ptr("eastus"),
						},
					},
				},
			}, nil)
			return
		},
	}

	// Create the client, connecting the fake server via client options
	client, err := armdiscovery.NewWorkspacesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewWorkspacesServerTransport(&fakeServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Call the Get API - the response comes from our fake implementation
	resp, err := client.Get(context.TODO(), "fake-resource-group", "test-workspace", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Workspace name:", *resp.Name)
	fmt.Println("Workspace location:", *resp.Location)

	// Call the list pager - the response comes from our fake implementation
	pager := client.NewListByResourceGroupPager("fake-resource-group", nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Number of workspaces:", len(page.Value))
	}

	// Output:
	// Workspace name: test-workspace
	// Workspace location: eastus
	// Number of workspaces: 2
}

// ExampleOperationsServer demonstrates how to use the fake server for Operations.
func ExampleOperationsServer() {
	// Create an instance of the fake server for the OperationsClient
	fakeServer := fake.OperationsServer{
		NewListPager: func(options *armdiscovery.OperationsClientListOptions) (resp azfake.PagerResponder[armdiscovery.OperationsClientListResponse]) {
			resp.AddPage(http.StatusOK, armdiscovery.OperationsClientListResponse{
				OperationListResult: armdiscovery.OperationListResult{
					Value: []*armdiscovery.Operation{
						{
							Name: to.Ptr("Microsoft.Discovery/workspaces/read"),
							Display: &armdiscovery.OperationDisplay{
								Provider:    to.Ptr("Microsoft Discovery"),
								Resource:    to.Ptr("Workspaces"),
								Operation:   to.Ptr("Get Workspace"),
								Description: to.Ptr("Gets a workspace"),
							},
						},
						{
							Name: to.Ptr("Microsoft.Discovery/workspaces/write"),
							Display: &armdiscovery.OperationDisplay{
								Provider:    to.Ptr("Microsoft Discovery"),
								Resource:    to.Ptr("Workspaces"),
								Operation:   to.Ptr("Create Workspace"),
								Description: to.Ptr("Creates or updates a workspace"),
							},
						},
					},
				},
			}, nil)
			return
		},
	}

	// Create the client with the fake server
	client, err := armdiscovery.NewOperationsClient(&azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewOperationsServerTransport(&fakeServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Call the list pager
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, op := range page.Value {
			fmt.Println("Operation:", *op.Name)
		}
	}

	// Output:
	// Operation: Microsoft.Discovery/workspaces/read
	// Operation: Microsoft.Discovery/workspaces/write
}
