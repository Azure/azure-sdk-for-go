//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	computefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4/fake"
)

// myFakeAvailabilitySetsServer provides fake implementations for armcompute.AvailabilitySetsClient.
// it wraps the generated fake server so consumers only have to implement the methods they care about.
// by default, the generated fake methods return an error with message "method Foo not implemented".
var myFakeAvailabilitySetsServer = computefake.AvailabilitySetsServer{
	// Get is a custom fake implementation. It responds to the armcompute.AvailabilitySetsClient.Get() method.
	Get: func(ctx context.Context, resourceGroupName string, availabilitySetName string, options *armcompute.AvailabilitySetsClientGetOptions) (resp azfake.Responder[armcompute.AvailabilitySetsClientGetResponse], err azfake.ErrorResponder) {
		// create a custom AvailabilitySetsClientGetResponse
		resp.Set(armcompute.AvailabilitySetsClientGetResponse{
			AvailabilitySet: armcompute.AvailabilitySet{
				ID: to.Ptr("this-should-be-a-resource-ID"),
			},
		})
		return
	},
}

func ExampleAvailabilitySetsServer_Get() {
	client, err := armcompute.NewAvailabilitySetsClient("subscriptionID", azfake.NewTokenCredential(), &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: computefake.NewAvailabilitySetsServerTransport(&myFakeAvailabilitySetsServer),
		},
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.Get(context.TODO(), "resourceGroupName", "availabilitySetName", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(*resp.ID)

	// output:
	// this-should-be-a-resource-ID
}
