// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-03-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

const (
	ipName = "sampleIP"
)

func getPublicIPClient() armnetwork.PublicIPAddressesOperations {
	return armnetwork.NewPublicIPAddressesClient(
		armnetwork.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

func ExamplePublicIPAddressesOperations_BeginCreateOrUpdate() {
	ipClient := getPublicIPClient()
	ipResp, err := ipClient.BeginCreateOrUpdate(
		context.Background(),
		resourceGroupName,
		ipName,
		armnetwork.PublicIPAddress{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(ipName),
				Location: to.StringPtr(location),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
			},
		},
		nil)
	if err != nil {
		panic(err)
	}
	result, err := ipResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	if result.PublicIPAddress != nil {
		fmt.Println(*result.PublicIPAddress.Name)
	}
	// Output:
	// sampleIP
}

func ExamplePublicIPAddressesOperations_Get() {
	ipClient := getPublicIPClient()
	ipResp, err := ipClient.Get(context.Background(), resourceGroupName, ipName, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*ipResp.PublicIPAddress.Name)
	// Output:
	// sampleIP
}

func ExamplePublicIPAddressesOperations_BeginDelete() {
	ipClient := getPublicIPClient()
	ipResp, err := ipClient.BeginDelete(context.Background(), resourceGroupName, ipName, nil)
	if err != nil {
		panic(err)
	}
	result, err := ipResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	if result != nil {
		fmt.Println(result.StatusCode)
	}
	// Output:
	// 200
}
