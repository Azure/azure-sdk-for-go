//go:build go1.18
// +build go1.18

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
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

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

// This examples shows how to create a directory with HTTP Headers, how to read, and how to update the directory's HTTP headers.
func Example_directory_HTTPHeaders() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
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

func Example_dir_Client_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a dir client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/ReadMe.txt", accountName)
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
