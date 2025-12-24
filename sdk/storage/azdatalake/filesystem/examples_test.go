// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_fs_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)
	fmt.Println(fsClient.DFSURL())
}

func Example_fs_NewClientWithSharedKeyCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fsClient, err := filesystem.NewClientWithSharedKeyCredential(fsURL, cred, nil)
	handleError(err)
	fmt.Println(fsClient.DFSURL())
}

func Example_fs_NewClientWithNoCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	sharedAccessSignature, ok := os.LookupEnv("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE")
	if !ok {
		panic("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s?%s", accountName, fsName, sharedAccessSignature)

	fsClient, err := filesystem.NewClientWithNoCredential(fsURL, nil)
	handleError(err)
	fmt.Println(fsClient.DFSURL())
}

func Example_fs_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	fsName := "testfs"
	fsClient, err := filesystem.NewClientFromConnectionString(connectionString, fsName, nil)
	handleError(err)
	fmt.Println(fsClient.DFSURL())
}

func Example_fs_ClientNewFileClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	fileClient := fsClient.NewFileClient("test_File")
	handleError(err)
	fmt.Println(fileClient.DFSURL())
}

func Example_fs_ClientCreate() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	fsCreateResponse, err := fsClient.Create(context.TODO(), &filesystem.CreateOptions{
		Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
	})
	handleError(err)
	fmt.Println(fsCreateResponse)
}

func Example_fs_ClientDelete() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	fsDeleteResponse, err := fsClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println(fsDeleteResponse)
}

func Example_fs_ClientListPaths() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	pager := fsClient.NewListPathsPager(true, nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range resp.Paths {
			fmt.Println(*path.Name)
		}
	}
}

func Example_fs_ClientListDirectoryPaths() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	pager := fsClient.NewListDirectoryPathsPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range resp.Segment.PathItems {
			fmt.Println(*path.Name)
		}
	}
}

func Example_fs_ClientListDeletedPaths() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	pager := fsClient.NewListDeletedPathsPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range resp.Segment.PathItems {
			fmt.Println(*path.Name)
		}
	}
}

func Example_fs_ClientGetSASURL() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	permission := sas.FileSystemPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	options := filesystem.GetSASURLOptions{StartTime: &start}
	sasURL, err := fsClient.GetSASURL(permission, expiry, &options)
	handleError(err)
	_ = sasURL
}

// This example shows how to manipulate a fs's permissions.
func Example_fs_ClientSetAccessPolicy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	// Create the fs
	_, err = fsClient.Create(context.TODO(), nil)
	handleError(err)

	// Upload a simple File.
	fileClient := fsClient.NewFileClient("HelloWorld.txt")
	handleError(err)

	err = fileClient.UploadStream(context.TODO(), streaming.NopCloser(strings.NewReader("Hello World!")), nil)
	handleError(err)

	// Attempt to read the File
	get, err := http.Get(fileClient.DFSURL())
	handleError(err)
	if get.StatusCode == http.StatusNotFound {
		// ChangeLease the File to be public access File
		_, err := fsClient.SetAccessPolicy(
			context.TODO(),
			&filesystem.SetAccessPolicyOptions{
				Access: to.Ptr(filesystem.File),
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		// Now, this works
		get, err = http.Get(fileClient.DFSURL())
		if err != nil {
			log.Fatal(err)
		}
		var text bytes.Buffer
		_, err = text.ReadFrom(get.Body)
		if err != nil {
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(get.Body)

		fmt.Println("Public access File data: ", text.String())
	}
}

func Example_fs_ClientSetMetadata() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	// Create a fs with some metadata, key names are converted to lowercase before being sent to the service.
	// You should always use lowercase letters, especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	handleError(err)
	_, err = fsClient.Create(context.TODO(), &filesystem.CreateOptions{Metadata: map[string]*string{"author": to.Ptr("azFile"), "app": to.Ptr(creatingApp)}})
	handleError(err)

	// Query the fs's metadata
	fsGetPropertiesResponse, err := fsClient.GetProperties(context.TODO(), nil)
	handleError(err)

	if fsGetPropertiesResponse.Metadata == nil {
		log.Fatal("metadata is empty!")
	}

	for k, v := range fsGetPropertiesResponse.Metadata {
		fmt.Printf("%s=%s\n", k, *v)
	}

	// Update the metadata and write it back to the fs
	fsGetPropertiesResponse.Metadata["author"] = to.Ptr("Mohit")
	_, err = fsClient.SetMetadata(context.TODO(), &filesystem.SetMetadataOptions{Metadata: fsGetPropertiesResponse.Metadata})
	handleError(err)
}

func Example_fs_ClientCreateFile() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)
	filePath := "testFile"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	fsCreateResponse, err := fsClient.Create(context.TODO(), &filesystem.CreateOptions{
		Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
	})
	handleError(err)
	fmt.Println(fsCreateResponse)

	createFileOptions := &filesystem.CreateFileOptions{
		Umask: to.Ptr("0000"),
		ACL:   to.Ptr("user::rwx,group::r-x,other::rwx"),
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeAbsolute,
			ExpiresOn:  time.Now().Add(20 * time.Second).UTC().Format(http.TimeFormat),
		},
		LeaseDuration:   to.Ptr(int64(15)),
		ProposedLeaseID: to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"),
	}
	resp, err := fsClient.CreateFile(context.Background(), filePath, createFileOptions)
	handleError(err)
	fmt.Println(resp)
}

func Example_fs_ClientCreateDirectory() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	fsName := "testfs"
	fsURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s", accountName, fsName)
	dirPath := "testDir"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	handleError(err)

	fsCreateResponse, err := fsClient.Create(context.TODO(), &filesystem.CreateOptions{
		Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
	})
	handleError(err)
	fmt.Println(fsCreateResponse)

	options := &filesystem.CreateDirectoryOptions{
		Umask:           to.Ptr("0000"),
		ACL:             to.Ptr("user::rwx,group::r-x,other::rwx"),
		LeaseDuration:   to.Ptr(int64(15)),
		ProposedLeaseID: to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"),
	}
	resp, err := fsClient.CreateDirectory(context.Background(), dirPath, options)
	handleError(err)
	fmt.Println(resp)
}
