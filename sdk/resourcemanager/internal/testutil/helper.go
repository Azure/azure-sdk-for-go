//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// GetEnv will return the os env variable and fallback to the given string if env variable not exist.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// CreateResourceGroup will create a resource group with a random generated name: "go-sdk-test-xxx".
// It will return the created resource group entity,
// a delegate function to delete the created resource group which can be used for clean up
// and any error during the creation.
func CreateResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, location string) (*armresources.ResourceGroup, func() (armresources.ResourceGroupsClientDeletePollerResponse, error), error) {
	rand.Seed(time.Now().UnixNano())
	resourceGroupName := fmt.Sprintf("go-sdk-test-%d", rand.Intn(1000))
	rgClient := armresources.NewResourceGroupsClient(subscriptionId, cred, options)
	param := armresources.ResourceGroup{
		Location: to.StringPtr(location),
	}
	resp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)
	if err != nil {
		return nil, nil, err
	}
	return &resp.ResourceGroup, func() (armresources.ResourceGroupsClientDeletePollerResponse, error) {
		return DeleteResourceGroup(ctx, subscriptionId, cred, options, *resp.Name)
	}, nil
}

// DeleteResourceGroup will delete the resource group with the given name.
// It will do the deletion asynchronously and return the poller which can be used to wait for the result.
func DeleteResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, resourceGroupName string) (armresources.ResourceGroupsClientDeletePollerResponse, error) {
	rgClient := armresources.NewResourceGroupsClient(subscriptionId, cred, options)
	return rgClient.BeginDelete(ctx, resourceGroupName, nil)
}
