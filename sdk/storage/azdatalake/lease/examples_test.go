//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/lease"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// This example shows how to perform various lease operations on a filesystem.
// The same lease operations can be performed on individual files as well.
// A lease on a filesystem prevents it from being deleted by others, while a lease on a file
// protects it from both modifications and deletions.
func Example_lease_FileSystemClient_AcquireLease() {
	// From the Azure portal, get your Storage account's name and account key.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// Create an fsClient object that wraps the filesystem's URL and a default pipeline.
	filesystemURL := fmt.Sprintf("https://%s.dfs.core.windows.net/myfs", accountName)
	fsClient, err := filesystem.NewClientWithSharedKeyCredential(filesystemURL, credential, nil)
	handleError(err)

	// Create a unique ID for the lease
	// A lease ID can be any valid GUID string format. To generate UUIDs, consider the github.com/google/uuid package
	leaseID := "36b1a876-cf98-4eb2-a5c3-6d68489658ff"
	filesystemLeaseClient, err := lease.NewFileSystemClient(fsClient,
		&lease.FileSystemClientOptions{LeaseID: to.Ptr(leaseID)})
	handleError(err)

	// Now acquire a lease on the filesystem.
	// You can choose to pass an empty string for proposed ID so that the service automatically assigns one for you.
	duration := int32(60)
	acquireLeaseResponse, err := filesystemLeaseClient.AcquireLease(context.TODO(), duration, nil)
	handleError(err)
	fmt.Println("The filesystem is leased for delete operations with lease ID", *acquireLeaseResponse.LeaseID)

	// The filesystem cannot be deleted without providing the lease ID.
	_, err = fsClient.Delete(context.TODO(), nil)
	if err == nil {
		log.Fatal("delete should have failed")
	}

	fmt.Println("The filesystem cannot be deleted while there is an active lease")

	// We can release the lease now and the filesystem can be deleted.
	_, err = filesystemLeaseClient.ReleaseLease(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease on the filesystem is now released")

	// AcquireLease a lease again to perform other operations.
	// Duration is still 60
	acquireLeaseResponse, err = filesystemLeaseClient.AcquireLease(context.TODO(), duration, nil)
	handleError(err)
	fmt.Println("The filesystem is leased again with lease ID", *acquireLeaseResponse.LeaseID)

	// We can change the ID of an existing lease.
	newLeaseID := "6b3e65e5-e1bb-4a3f-8b72-13e9bc9cd3bf"
	changeLeaseResponse, err := filesystemLeaseClient.ChangeLease(context.TODO(), newLeaseID, nil)
	handleError(err)
	fmt.Println("The lease ID was changed to", *changeLeaseResponse.LeaseID)

	// The lease can be renewed.
	renewLeaseResponse, err := filesystemLeaseClient.RenewLease(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease was renewed with the same ID", *renewLeaseResponse.LeaseID)

	// Finally, the lease can be broken, and we could prevent others from acquiring a lease for a period of time
	_, err = filesystemLeaseClient.BreakLease(context.TODO(), &lease.FileSystemBreakOptions{BreakPeriod: to.Ptr(int32(60))})
	handleError(err)
	fmt.Println("The lease was broken, and nobody can acquire a lease for 60 seconds")
}
