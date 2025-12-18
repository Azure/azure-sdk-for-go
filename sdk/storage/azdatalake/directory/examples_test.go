//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"log"
	"os"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// make sure you create the filesystem before running this example
func Example_directory_CreateAndDelete() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a directory client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)

	_, err = dirClient.Delete(context.Background(), nil)
	handleError(err)
}

// This examples shows how to set a directory HTTP Headers, how to read, and how to update the directory's HTTP headers.
func Example_directory_HTTPHeaders() {
	// make sure you create the filesystem and directory before running this example
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a dir client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a directory with HTTP headers
	_, err = dirClient.SetHTTPHeaders(context.TODO(), directory.HTTPHeaders{
		ContentType:        to.Ptr("text/html; charset=utf-8"),
		ContentDisposition: to.Ptr("attachment"),
	}, nil)
	handleError(err)

	get, err := dirClient.GetProperties(context.TODO(), nil)
	handleError(err)

	fmt.Println(get.ContentType)
	fmt.Println(get.ContentDisposition)
}

// make sure you create the filesystem and directory before running this example
func Example_dir_Client_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	// Create a dir client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = dirClient.SetMetadata(context.TODO(), map[string]*string{"author": to.Ptr("Tamer")}, nil)
	handleError(err)

	// Query the directory's properties and metadata
	get, err := dirClient.GetProperties(context.TODO(), nil)
	handleError(err)

	// Show the directory's metadata
	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}

	for k, v := range get.Metadata {
		fmt.Print(k + "=" + *v + "\n")
	}
}

// make sure you create the filesystem before running this example
func Example_directory_Rename() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a directory client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)

	_, err = dirClient.Rename(context.Background(), "renameDir", nil)
	handleError(err)
}

// for this example make sure to create paths within the dir so you can recursively set the ACL on them
func Example_directory_SetACLRecursive() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a directory client
	acl := "user::rwx,group::r-x,other::rwx"
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = dirClient.SetAccessControlRecursive(context.Background(), acl, &directory.SetAccessControlRecursiveOptions{
		BatchSize: to.Ptr(int32(2)), MaxBatches: to.Ptr(int32(1)), ContinueOnFailure: to.Ptr(true), Marker: nil})
	handleError(err)
}

func getRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

// make sure you create the filesystem before running this example
func Example_directory_SetAccessControlIfUnmodifiedSinceTrue() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a directory client
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/dir1", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	dirClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)
	resp, err := dirClient.Create(context.Background(), nil)
	handleError(err)

	currentTime := getRelativeTimeFromAnchor(resp.Date, 10)
	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	handleError(err)
}
