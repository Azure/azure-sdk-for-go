// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armresources_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2019-05-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

const (
	resourceGroupName = "sampleRG"
	location          = "westus"
)

var (
	// The subscription ID where the resource group exists
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
)

// returns a credential that can be used to authenticate with Azure Active Directory
func getCredential() azcore.Credential {
	// NewEnvironmentCredential() will read various environment vars
	// to obtain a credential.  see the documentation for more info.
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		panic(err)
	}
	return cred
}

func getResourceGroupsOperations() armresources.ResourceGroupsOperations {
	return armresources.NewResourceGroupsClient(
		armresources.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

func ExampleResourceGroupsOperations_CreateOrUpdate() {
	rgClient := getResourceGroupsOperations()
	rgResp, err := rgClient.CreateOrUpdate(
		context.Background(),
		resourceGroupName,
		armresources.ResourceGroup{
			Name:     to.StringPtr(resourceGroupName),
			Location: to.StringPtr(location),
		},
		nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*rgResp.ResourceGroup.Name)
	// Output:
	// sampleRG
}

func ExampleResourceGroupsOperations_BeginDelete() {
	rgClient := getResourceGroupsOperations()
	rgResp, err := rgClient.BeginDelete(context.Background(), resourceGroupName, nil)
	if err != nil {
		panic(err)
	}
	// the following demonstrates the recommended way to manually handle polling
	poller := rgResp.Poller
	for {
		resp, err := poller.Poll(context.Background())
		if err != nil {
			panic(err)
		}
		if poller.Done() {
			break
		}
		if delay := azcore.RetryAfter(resp); delay > 0 {
			time.Sleep(delay)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	res, err := poller.FinalResponse(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)
	// Output:
	// 200
}
