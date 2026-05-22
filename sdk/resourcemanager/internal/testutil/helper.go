// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeployments"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
)

// CreateResourceGroup will create a resource group with a random generated name: "go-sdk-test-xxx".
// It will return the created resource group entity,
// a delegate function to delete the created resource group which can be used for clean up
// and any error during the creation.
func CreateResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, location string) (*armresources.ResourceGroup, func() (*runtime.Poller[armresources.ResourceGroupsClientDeleteResponse], error), error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	resourceGroupName := fmt.Sprintf("go-sdk-test-%d", r.Intn(1000))
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, options)
	if err != nil {
		return nil, nil, err
	}
	param := armresources.ResourceGroup{
		Location: to.Ptr(location),
	}
	resp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)
	if err != nil {
		return nil, nil, err
	}
	return &resp.ResourceGroup, func() (*runtime.Poller[armresources.ResourceGroupsClientDeleteResponse], error) {
		return DeleteResourceGroup(ctx, subscriptionId, cred, options, *resp.Name)
	}, nil
}

// DeleteResourceGroup will delete the resource group with the given name.
// It will do the deletion asynchronously and return the poller which can be used to wait for the result.
func DeleteResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, resourceGroupName string) (*runtime.Poller[armresources.ResourceGroupsClientDeleteResponse], error) {
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, options)
	if err != nil {
		return nil, err
	}
	return rgClient.BeginDelete(ctx, resourceGroupName, nil)
}

// CreateDeployment will create a resource using arm template.
// It will return the deployment result entity.
func CreateDeployment(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, resourceGroupName, deploymentName string, deployment *armdeployments.Deployment) (*armdeployments.DeploymentExtended, error) {
	deployClient, err := armdeployments.NewDeploymentsClient(subscriptionId, cred, options)
	if err != nil {
		return nil, err
	}
	poller, err := deployClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		*deployment,
		&armdeployments.DeploymentsClientBeginCreateOrUpdateOptions{},
	)
	if err != nil {
		return nil, err
	}
	res, err := PollForTest(ctx, poller)
	if err != nil {
		return nil, err
	}
	return &res.DeploymentExtended, nil
}

// PollForTest will poll result according to the recording mode:
// Playback: customer poll loop until get result
// Others: use original poll until done
func PollForTest[T any](ctx context.Context, poller *runtime.Poller[T]) (*T, error) {
	pollOptions := runtime.PollUntilDoneOptions{
		Frequency: 0, // Pass zero to accept the default value (30s).
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		pollOptions.Frequency = time.Millisecond // If playback, do not wait
	}
	res, err := poller.PollUntilDone(ctx, &pollOptions)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
