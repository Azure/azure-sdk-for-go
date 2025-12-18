//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// This example shows how to perform various lease operations on a share.
// The same lease operations can be performed on individual files as well.
// A lease on a share prevents it from being deleted by others, while a lease on a file
// protects it from both modifications and deletions.
func Example_lease_ShareClient_AcquireLease() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	// Create a unique ID for the lease
	// A lease ID can be any valid GUID string format. To generate UUIDs, consider the github.com/google/uuid package
	leaseID := "36b1a876-cf98-4eb2-a5c3-6d68489658ff"
	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{LeaseID: to.Ptr(leaseID)})
	handleError(err)

	// Now acquire a lease on the share.
	// You can choose to pass an empty string for proposed ID so that the service automatically assigns one for you.
	duration := int32(60)
	acquireLeaseResponse, err := shareLeaseClient.Acquire(context.TODO(), duration, nil)
	handleError(err)
	fmt.Println("The share is leased for delete operations with lease ID", *acquireLeaseResponse.LeaseID)

	// The share cannot be deleted without providing the lease ID.
	_, err = shareClient.Delete(context.TODO(), nil)
	if err == nil {
		log.Fatal("delete should have failed")
	}

	fmt.Println("The share cannot be deleted while there is an active lease")

	// We can release the lease now and the share can be deleted.
	_, err = shareLeaseClient.Release(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease on the share is now released")

	// AcquireLease a lease again to perform other operations.
	// Duration is still 60
	acquireLeaseResponse, err = shareLeaseClient.Acquire(context.TODO(), duration, nil)
	handleError(err)
	fmt.Println("The share is leased again with lease ID", *acquireLeaseResponse.LeaseID)

	// We can change the ID of an existing lease.
	newLeaseID := "6b3e65e5-e1bb-4a3f-8b72-13e9bc9cd3bf"
	changeLeaseResponse, err := shareLeaseClient.Change(context.TODO(), newLeaseID, nil)
	handleError(err)
	fmt.Println("The lease ID was changed to", *changeLeaseResponse.LeaseID)

	// The lease can be renewed.
	renewLeaseResponse, err := shareLeaseClient.Renew(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease was renewed with the same ID", *renewLeaseResponse.LeaseID)

	// Finally, the lease can be broken, and we could prevent others from acquiring a lease for a period of time
	_, err = shareLeaseClient.Break(context.TODO(), &lease.ShareBreakOptions{BreakPeriod: to.Ptr(int32(60))})
	handleError(err)
	fmt.Println("The lease was broken, and nobody can acquire a lease for 60 seconds")
}
