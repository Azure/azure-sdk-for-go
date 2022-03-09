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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func CreateResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, location string) (*armresources.ResourceGroup, func() error, error) {
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
	return &resp.ResourceGroup, func() error {
		err := DeleteResourceGroup(ctx, subscriptionId, cred, options, *resp.Name)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

func DeleteResourceGroup(ctx context.Context, subscriptionId string, cred azcore.TokenCredential, options *arm.ClientOptions, resourceGroupName string) error {
	rgClient := armresources.NewResourceGroupsClient(subscriptionId, cred, options)
	poller, err := rgClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = poller.Poller.Poll(ctx)
			if err != nil {
				return err
			}
			if poller.Poller.Done() {
				_, err = poller.Poller.FinalResponse(ctx)
				if err != nil {
					return err
				}
				break
			}
		}
	} else {
		if _, err = poller.PollUntilDone(ctx, 10*time.Second); err != nil {
			return err
		}
	}
	return nil
}
