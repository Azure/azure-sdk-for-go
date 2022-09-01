//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func handleError(err error) {
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		log.Fatalf("HTTPErrorCode: %s\nStorageErrorCode: %dErrorMessage: %s\n", responseErr.ErrorCode, responseErr.StatusCode, responseErr.Error())
	} else {
		log.Fatal(err.Error())
	}
}

// This example is a quick-starter and demonstrates how to get started using the Azure Blob Storage SDK for Go.
func Example() {
	// Your account name and key can be obtained from the Azure Portal.
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	handleError(err)

	// ===== 1. Create a container =====
	containerName := "testcontainer"
	containerCreateResp, err := client.CreateContainer(context.TODO(), containerName, nil)
	handleError(err)
	fmt.Println(containerCreateResp)

	// ===== 2. Upload and Download a block blob =====
	blobData := "Hello world!"
	blobName := "HelloWorld.txt"
	uploadResp, err := client.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: map[string]string{"Foo": "Bar"},
			Tags:     map[string]string{"Year": "2022"},
		})
	handleError(err)
	fmt.Println(uploadResp)

	// Download the blob's contents and ensure that the download worked properly
	blobDownloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)
	handleError(err)

	// Use the bytes.Buffer object to read the downloaded data.
	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	reader := blobDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	handleError(err)
	if string(downloadData) != blobData {
		log.Fatal("Uploaded data should be same as downloaded data")
	}

	err = reader.Close()
	if err != nil {
		return
	}

	// ===== 3. List blobs =====
	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
	// PageResponse() can be used to iterate over the results of the specific page.
	pager := client.NewListBlobsPager(containerName, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)
		for _, v := range resp.Segment.BlobItems {
			fmt.Println(*v.Name)
		}
	}

	// Delete the blob.
	_, err = client.DeleteBlob(context.TODO(), containerName, blobName, nil)
	handleError(err)

	// Delete the container.
	_, err = client.DeleteContainer(context.TODO(), containerName, nil)
	handleError(err)
}

// blob client ----------------------------------------------------------------------------------------------------

func Example_client_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	fmt.Println(serviceClient.URL())
}

func Example_client_NewClientWithSharedKeyCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	serviceClient, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_client_CreateContainer() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := serviceClient.CreateContainer(context.TODO(), "testcontainer", &azblob.CreateContainerOptions{
		Metadata: map[string]string{"hello": "world"},
	})
	handleError(err)
	fmt.Println(resp)
}

func Example_client_DeleteContainer() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := serviceClient.DeleteContainer(context.TODO(), "testcontainer", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_NewListContainersPager() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	pager := serviceClient.NewListContainersPager(&azblob.ListContainersOptions{
		Include: azblob.ListContainersDetail{Metadata: true, Deleted: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		handleError(err) // if err is not nil, break the loop.
		for _, _container := range resp.ContainerItems {
			fmt.Printf("%v", _container)
		}
	}
}

func Example_client_UploadFile() {
	// Set up file to upload
	fileSize := 8 * 1024 * 1024
	fileName := "test_upload_file.txt"
	fileData := make([]byte, fileSize)
	err := os.WriteFile(fileName, fileData, 0666)
	handleError(err)

	// Open the file to upload
	fileHandler, err := os.Open(fileName)
	handleError(err)

	// close the file after it is no longer required.
	defer func(file *os.File) {
		err = file.Close()
		handleError(err)
	}(fileHandler)

	// delete the local file if required.
	defer func(name string) {
		err = os.Remove(name)
		handleError(err)
	}(fileName)

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Upload the file to a block blob
	_, err = serviceClient.UploadFile(context.TODO(), "testcontainer", "virtual/dir/path/"+fileName, fileHandler,
		&azblob.UploadFileOptions{
			BlockSize:   int64(1024),
			Parallelism: uint16(3),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Println(bytesTransferred)
			},
		})
	handleError(err)
}

func Example_client_DownloadFile() {
	// Set up file to download the blob to
	destFileName := "test_download_file.txt"
	destFile, err := os.Create(destFileName)
	handleError(err)
	defer func(destFile *os.File) {
		err = destFile.Close()
		handleError(err)
	}(destFile)

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Perform download

	_, err = serviceClient.DownloadFile(context.TODO(), "testcontainer", "virtual/dir/path/"+destFileName, destFile,
		&azblob.DownloadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Println(bytesTransferred)
			},
		})

	// Assert download was successful
	handleError(err)
}

func Example_client_NewListBlobsPager() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	pager := serviceClient.NewListBlobsPager("testcontainer", &azblob.ListBlobsOptions{
		Include: []azblob.ListBlobsIncludeItem{azblob.ListBlobsIncludeItemVersions, azblob.ListBlobsIncludeItemDeleted},
	})

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		handleError(err) // if err is not nil, break the loop.
		for _, _blob := range resp.Segment.BlobItems {
			fmt.Printf("%v", _blob.Name)
		}
	}
}

func Example_client_DeleteBlob() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := serviceClient.DeleteBlob(context.TODO(), "testcontainer", "testblob", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_UploadStream() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Set up test blob
	containerName := "testcontainer"
	bufferSize := 8 * 1024 * 1024
	blobName := "test_upload_stream.bin"
	blobData := make([]byte, bufferSize)
	blobContentReader := bytes.NewReader(blobData)

	// Perform UploadStream
	resp, err := serviceClient.UploadStream(context.TODO(), containerName, blobName, blobContentReader,
		&azblob.UploadStreamOptions{
			Metadata: map[string]string{"hello": "world"},
		})
	// Assert that upload was successful
	handleError(err)
	fmt.Println(resp)
}

func Example_client_DownloadStream() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Download the blob
	downloadResponse, err := serviceClient.DownloadStream(ctx, "testcontainer", "test_download_stream.bin", nil)
	handleError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.Body)
	handleError(err)
	fmt.Println(len(actualBlobData))
}

// service client ------------------------------------------------------------------------------------------------------

func Example_service_Client_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_service_Client_NewClientWithSharedKeyCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	serviceClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_service_Client_NewClientWithNoCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	sharedAccessSignature, ok := os.LookupEnv("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE")
	if !ok {
		panic("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, sharedAccessSignature)

	serviceClient, err := service.NewClientWithNoCredential(serviceURL, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_service_Client_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := service.NewClientFromConnectionString(connectionString, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_service_Client_CreateContainer() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)

	_, err = serviceClient.CreateContainer(context.TODO(), "testcontainer", nil)
	handleError(err)

	// ======== 2. Delete a container ========
	defer func(serviceClient1 *service.Client, ctx context.Context, containerName string, options *container.DeleteOptions) {
		_, err = serviceClient1.DeleteContainer(ctx, containerName, options)
		if err != nil {
			log.Fatal(err)
		}
	}(serviceClient, context.TODO(), "testcontainer", nil)
}

func Example_service_Client_DeleteContainer() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)

	_, err = serviceClient.DeleteContainer(context.TODO(), "testcontainer", nil)
	handleError(err)
}

func Example_service_Client_ListContainers() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)

	listContainersOptions := service.ListContainersOptions{
		Include: service.ListContainersDetail{
			Metadata: true, // Include Metadata
			Deleted:  true, // Include deleted containers in the result as well
		},
	}
	pager := serviceClient.NewListContainersPager(&listContainersOptions)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, container := range resp.ContainerItems {
			fmt.Println(*container.Name)
		}
	}
}

func Example_service_Client_GetSASURL() {
	cred, err := azblob.NewSharedKeyCredential("myAccountName", "myAccountKey")
	handleError(err)
	serviceClient, err := service.NewClientWithSharedKeyCredential("https://<myAccountName>.blob.core.windows.net", cred, nil)
	handleError(err)

	resources := service.SASResourceTypes{Service: true}
	permission := service.SASPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	sasURL, err := serviceClient.GetSASURL(resources, permission, service.SASServices{Blob: true}, start, expiry)
	handleError(err)

	serviceURL := fmt.Sprintf("https://<myAccountName>.blob.core.windows.net/?%s", sasURL)
	serviceClientWithSAS, err := service.NewClientWithNoCredential(serviceURL, nil)
	handleError(err)
	_ = serviceClientWithSAS
}

func Example_service_Client_SetProperties() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)
	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)

	enabled := true  // enabling retention period
	days := int32(5) // setting retention period to 5 days
	serviceSetPropertiesResponse, err := serviceClient.SetProperties(context.TODO(), &service.SetPropertiesOptions{
		DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days},
	})

	handleError(err)
	fmt.Println(serviceSetPropertiesResponse)
}

func Example_service_Client_GetProperties() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	serviceClient, err := service.NewClient(serviceURL, cred, nil)
	handleError(err)

	serviceGetPropertiesResponse, err := serviceClient.GetProperties(context.TODO(), nil)
	handleError(err)

	fmt.Println(serviceGetPropertiesResponse)
}

// container client ----------------------------------------------------------------------------------------------------

func Example_container_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientWithSharedKeyCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, cred, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientWithNoCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	sharedAccessSignature, ok := os.LookupEnv("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE")
	if !ok {
		panic("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s?%s", accountName, containerName, sharedAccessSignature)

	containerClient, err := container.NewClientWithNoCredential(containerURL, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	containerName := "testcontainer"
	containerClient, err := container.NewClientFromConnectionString(connectionString, containerName, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_ClientNewAppendBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	appendBlobClient := containerClient.NewAppendBlobClient("test_append_blob")
	handleError(err)
	fmt.Println(appendBlobClient.URL())
}

func Example_container_ClientNewBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	blobClient := containerClient.NewBlobClient("test_blob")
	handleError(err)
	fmt.Println(blobClient.URL())
}

func Example_container_ClientNewBlockBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	blockBlobClient := containerClient.NewBlockBlobClient("test_block_blob")
	handleError(err)
	fmt.Println(blockBlobClient.URL())
}

func Example_container_ClientNewPageBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	pageBlobClient := containerClient.NewPageBlobClient("test_page_blob")
	handleError(err)
	fmt.Println(pageBlobClient.URL())
}

func Example_container_ClientCreate() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	containerCreateResponse, err := containerClient.Create(context.TODO(), &container.CreateOptions{
		Metadata: map[string]string{"Foo": "Bar"},
	})
	handleError(err)
	fmt.Println(containerCreateResponse)
}

func Example_container_ClientDelete() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	containerDeleteResponse, err := containerClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println(containerDeleteResponse)
}

func Example_container_ClientListBlobsFlat() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemSnapshots, container.ListBlobsIncludeItemVersions},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
}

func Example_container_ClientListBlobsHierarchy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	maxResults := int32(5)
	pager := containerClient.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{
		Include: []container.ListBlobsIncludeItem{
			container.ListBlobsIncludeItemMetadata,
			container.ListBlobsIncludeItemTags,
		},
		MaxResults: &maxResults,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.ListBlobsHierarchySegmentResponse.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
}

func Example_container_ClientGetSASURL() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	permission := container.SASPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	sasURL, err := containerClient.GetSASURL(permission, start, expiry)
	handleError(err)
	_ = sasURL
}

// This example shows how to manipulate a container's permissions.
func Example_container_ClientSetAccessPolicy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	// Create the container
	_, err = containerClient.Create(context.TODO(), nil)
	handleError(err)

	// Upload a simple blob.
	blockBlobClient := containerClient.NewBlockBlobClient("HelloWorld.txt")
	handleError(err)

	_, err = blockBlobClient.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("Hello World!")), nil)
	handleError(err)

	// Attempt to read the blob
	get, err := http.Get(blockBlobClient.URL())
	handleError(err)
	if get.StatusCode == http.StatusNotFound {
		// ChangeLease the blob to be public access blob
		_, err := containerClient.SetAccessPolicy(
			context.TODO(),
			nil,
			&container.SetAccessPolicyOptions{
				Access: to.Ptr(container.PublicAccessTypeBlob),
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		// Now, this works
		get, err = http.Get(blockBlobClient.URL())
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

		fmt.Println("Public access blob data: ", text.String())
	}
}

func Example_container_ClientSetMetadata() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	// Create a container with some metadata, key names are converted to lowercase before being sent to the service.
	// You should always use lowercase letters, especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	handleError(err)
	_, err = containerClient.Create(context.TODO(), &container.CreateOptions{Metadata: map[string]string{"author": "azblob", "app": creatingApp}})
	handleError(err)

	// Query the container's metadata
	containerGetPropertiesResponse, err := containerClient.GetProperties(context.TODO(), nil)
	handleError(err)

	if containerGetPropertiesResponse.Metadata == nil {
		log.Fatal("metadata is empty!")
	}

	for k, v := range containerGetPropertiesResponse.Metadata {
		fmt.Printf("%s=%s\n", k, v)
	}

	// Update the metadata and write it back to the container
	containerGetPropertiesResponse.Metadata["author"] = "Mohit"
	_, err = containerClient.SetMetadata(context.TODO(), &container.SetMetadataOptions{Metadata: containerGetPropertiesResponse.Metadata})
	handleError(err)

	// NOTE: SetMetadata & SetProperties methods update the container's ETag & LastModified properties
}

// blob client ---------------------------------------------------------------------------------------------------

// ExampleBlockBlobClient shows how to upload data (in blocks) to a blob.
// A block blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of a block blob is slightly more than 190 TiB (4000 MiB X 50,000 blocks).
func Example_blockblob_Client() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_block_blob.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	blockBlobClient, err := blockblob.NewClient(blobURL, cred, nil)
	handleError(err)

	// NOTE: The blockID must be <= 64 bytes and ALL blockIDs for the block must be the same length
	blockIDBinaryToBase64 := func(blockID []byte) string { return base64.StdEncoding.EncodeToString(blockID) }
	blockIDBase64ToBinary := func(blockID string) []byte { _binary, _ := base64.StdEncoding.DecodeString(blockID); return _binary }

	// These helper functions convert an int block ID to a base-64 string and vice versa
	blockIDIntToBase64 := func(blockID int) string {
		binaryBlockID := &[4]byte{} // All block IDs are 4 bytes long
		binary.LittleEndian.PutUint32(binaryBlockID[:], uint32(blockID))
		return blockIDBinaryToBase64(binaryBlockID[:])
	}
	blockIDBase64ToInt := func(blockID string) int {
		blockIDBase64ToBinary(blockID)
		return int(binary.LittleEndian.Uint32(blockIDBase64ToBinary(blockID)))
	}

	// Upload 4 blocks to the blob (these blocks are tiny; they can be up to 100MB each)
	words := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(words)) // The collection of block IDs (base 64 strings)

	// Upload each block sequentially (one after the other)
	for index, word := range words {
		// This example uses the index as the block ID; convert the index/ID into a base-64 encoded string as required by the service.
		// NOTE: Over the lifetime of a blob, all block IDs (before base 64 encoding) must be the same length (this example uses 4 byte block IDs).
		base64BlockIDs[index] = blockIDIntToBase64(index)

		// Upload a block to this blob specifying the Block ID and its content (up to 100MB); this block is uncommitted.
		_, err := blockBlobClient.StageBlock(context.TODO(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// After all the blocks are uploaded, atomically commit them to the blob.
	_, err = blockBlobClient.CommitBlockList(context.TODO(), base64BlockIDs, nil)
	handleError(err)

	// For the blob, show each block (ID and size) that is a committed part of it.
	getBlock, err := blockBlobClient.GetBlockList(context.TODO(), blockblob.BlockListTypeAll, nil)
	handleError(err)
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fmt.Printf("Block ID=%d, Size=%d\n", blockIDBase64ToInt(*block.Name), block.Size)
	}

	// Download the blob in its entirety; download operations do not take blocks into account.
	blobDownloadResponse, err := blockBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)

	blobData := &bytes.Buffer{}
	reader := blobDownloadResponse.Body
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(blobData)
}

// ExampleAppendBlobClient shows how to append data (in blocks) to an append blob.
// An append blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of an append blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func Example_appendblob_Client() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_append_blob.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	appendBlobClient, err := appendblob.NewClient(blobURL, cred, nil)
	handleError(err)

	_, err = appendBlobClient.Create(context.TODO(), nil)
	handleError(err)

	for i := 0; i < 5; i++ { // Append 5 blocks to the append blob
		_, err := appendBlobClient.AppendBlock(context.TODO(), streaming.NopCloser(strings.NewReader(fmt.Sprintf("Appending block #%d\n", i))), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Download the entire append blob's contents and read into a bytes.Buffer.
	get, err := appendBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	b := bytes.Buffer{}
	reader := get.Body
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(b.String())
}

// ExamplePageBlobClient shows how to manipulate a page blob with PageBlobClient.
// A page blob is a collection of 512-byte pages optimized for random read and write operations.
// The maximum size for a page blob is 8 TB.
func Example_pageblob_Client() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_page_blob.vhd"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	pageBlobClient, err := pageblob.NewClient(blobURL, cred, nil)
	handleError(err)

	_, err = pageBlobClient.Create(context.TODO(), pageblob.PageBytes*4, nil)
	handleError(err)

	page := make([]byte, pageblob.PageBytes)
	copy(page, "Page 0")
	_, err = pageBlobClient.UploadPages(context.TODO(), streaming.NopCloser(bytes.NewReader(page)), nil)
	handleError(err)

	copy(page, "Page 1")
	_, err = pageBlobClient.UploadPages(
		context.TODO(),
		streaming.NopCloser(bytes.NewReader(page)),
		&pageblob.UploadPagesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(2 * pageblob.PageBytes))},
	)
	handleError(err)

	pager := pageBlobClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(10 * pageblob.PageBytes)),
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, pr := range resp.PageList.PageRange {
			fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
		}
	}

	_, err = pageBlobClient.ClearPages(context.TODO(), 0, 1*pageblob.PageBytes, nil)
	handleError(err)

	pager = pageBlobClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(10 * pageblob.PageBytes)),
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, pr := range resp.PageList.PageRange {
			fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
		}
	}

	get, err := pageBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	blobData := &bytes.Buffer{}
	reader := get.Body
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(blobData.String())
}

// ---------------------------------------------------------------------------------------------------------------------

// This example shows how to copy a large stream in blocks (chunks) to a block blob.
func Example_blockblob_Client_UploadFile() {
	file, err := os.Open("BigFile.bin") // Open the file we want to upload
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	fileSize, err := file.Stat() // Get the size of the file (stream)
	if err != nil {
		log.Fatal(err)
	}

	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blockBlobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := blockBlobClient.UploadFile(context.TODO(), file,
		&blockblob.UploadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Printf("Uploaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	_ = response // Avoid compiler's "declared and not used" error

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(destFile *os.File) {
		_ = destFile.Close()

	}(destFile)

	// Perform download
	_, err = blockBlobClient.DownloadFile(context.TODO(), destFile,
		&blob.DownloadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			}})

	if err != nil {
		log.Fatal(err)
	}
}

// This example shows how to perform various lease operations on a container.
// The same lease operations can be performed on individual blobs as well.
// A lease on a container prevents it from being deleted by others, while a lease on a blob
// protects it from both modifications and deletions.
func Example_lease_ContainerClient_AcquireLease() {
	// From the Azure portal, get your Storage account's name and account key.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// Create an containerClient object that wraps the container's URL and a default pipeline.
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)
	handleError(err)

	// Create a unique ID for the lease
	// A lease ID can be any valid GUID string format. To generate UUIDs, consider the github.com/google/uuid package
	leaseID := "36b1a876-cf98-4eb2-a5c3-6d68489658ff"
	containerLeaseClient, err := lease.NewContainerClient(containerClient,
		&lease.ContainerClientOptions{LeaseID: to.Ptr(leaseID)})
	handleError(err)

	// Now acquire a lease on the container.
	// You can choose to pass an empty string for proposed ID so that the service automatically assigns one for you.
	duration := int32(60)
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(context.TODO(),
		&lease.ContainerAcquireOptions{Duration: &duration})
	handleError(err)
	fmt.Println("The container is leased for delete operations with lease ID", *acquireLeaseResponse.LeaseID)

	// The container cannot be deleted without providing the lease ID.
	_, err = containerClient.Delete(context.TODO(), nil)
	if err == nil {
		log.Fatal("delete should have failed")
	}

	fmt.Println("The container cannot be deleted while there is an active lease")

	// _, err = containerClient.Delete(context.TODO(), &container.DeleteOptions{
	// 	AccessConditions: &container.AccessConditions{
	// 		LeaseAccessConditions: &container.LeaseAccessConditions{LeaseID: acquireLeaseResponse.LeaseID},
	// 	},
	// })

	// We can release the lease now and the container can be deleted.
	_, err = containerLeaseClient.ReleaseLease(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease on the container is now released")

	// AcquireLease a lease again to perform other operations.
	// Duration is still 60
	acquireLeaseResponse, err = containerLeaseClient.AcquireLease(context.TODO(),
		&lease.ContainerAcquireOptions{Duration: &duration})
	handleError(err)
	fmt.Println("The container is leased again with lease ID", *acquireLeaseResponse.LeaseID)

	// We can change the ID of an existing lease.
	newLeaseID := "6b3e65e5-e1bb-4a3f-8b72-13e9bc9cd3bf"
	changeLeaseResponse, err := containerLeaseClient.ChangeLease(context.TODO(),
		&lease.ContainerChangeOptions{ProposedLeaseID: to.Ptr(newLeaseID)})
	handleError(err)
	fmt.Println("The lease ID was changed to", *changeLeaseResponse.LeaseID)

	// The lease can be renewed.
	renewLeaseResponse, err := containerLeaseClient.RenewLease(context.TODO(), nil)
	handleError(err)
	fmt.Println("The lease was renewed with the same ID", *renewLeaseResponse.LeaseID)

	// Finally, the lease can be broken, and we could prevent others from acquiring a lease for a period of time
	_, err = containerLeaseClient.BreakLease(context.TODO(), &lease.ContainerBreakOptions{BreakPeriod: to.Ptr(int32(60))})
	handleError(err)
	fmt.Println("The lease was broken, and nobody can acquire a lease for 60 seconds")
}

func ExampleResponseError() {
	contClient, err := container.NewClientWithNoCredential("https://myaccount.blob.core.windows.net/mycontainer", nil)
	handleError(err)
	_, err = contClient.Create(context.TODO(), nil)
	handleError(err)
}

// This example demonstrates splitting a URL into its parts so you can examine and modify the URL in an Azure Storage fluent way.
func ExampleParseURL() {
	// Here is an example of a blob snapshot.
	u := "https://myaccount.blob.core.windows.net/mycontainter/ReadMe.txt?" +
		"snapshot=2011-03-09T01:42:34Z&" +
		"sv=2015-02-21&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&" +
		"spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=92836758923659283652983562=="

	// Breaking the URL down into it's parts by conversion to URLParts
	parts, _ := azblob.ParseURL(u)

	// The URLParts allows access to individual portions of a Blob URL
	fmt.Printf("Host: %s\nContainerName: %s\nBlobName: %s\nSnapshot: %s\n", parts.Host, parts.ContainerName, parts.BlobName, parts.Snapshot)
	fmt.Printf("Version: %s\nResource: %s\nStartTime: %s\nExpiryTime: %s\nPermissions: %s\n", parts.SAS.Version(), parts.SAS.Resource(), parts.SAS.StartTime(), parts.SAS.ExpiryTime(), parts.SAS.Permissions())
}

// This example shows how to create and use an Azure Storage account Shared Access Signature (SAS).
func Example_service_SASSignatureValues_Sign() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	sasQueryParams, err := service.SASSignatureValues{
		Protocol:      service.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   to.Ptr(service.SASPermissions{Read: true, List: true}).String(),
		Services:      to.Ptr(service.SASServices{Blob: true}).String(),
		ResourceTypes: to.Ptr(service.SASResourceTypes{Container: true, Object: true}).String(),
	}.Sign(credential)
	handleError(err)

	queryParams := sasQueryParams.Encode()
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, queryParams)

	// This URL can be used to authenticate requests now
	serviceClient, err := service.NewClientWithNoCredential(sasURL, nil)
	handleError(err)

	// You can also break a blob URL up into it's constituent parts
	blobURLParts, _ := azblob.ParseURL(serviceClient.URL())
	fmt.Printf("SAS expiry time = %s\n", blobURLParts.SAS.ExpiryTime())
}

// This example shows how to perform operations on blob conditionally.
func Example_blob_AccessConditions() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blockBlob, err := blockblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/Data.txt", accountName), credential, nil)
	handleError(err)

	// This function displays the results of an operation
	showResult := func(response *blob.DownloadStreamResponse, err error) {
		if err != nil {
			log.Fatalf("Failure: %s\n", err.Error())
		} else {
			err := response.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			// The client must close the response body when finished with it
			fmt.Printf("Success: %v\n", response)
		}

		// Close the response
		if err != nil {
			return
		}
		fmt.Printf("Success: %v\n", response)
	}

	showResultUpload := func(response blockblob.UploadResponse, err error) {
		if err != nil {
			log.Fatalf("Failure: %s\n", err.Error())
		}
		fmt.Printf("Success: %v\n", response)
	}

	// Create the blob
	upload, err := blockBlob.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("Text-1")), nil)
	showResultUpload(upload, err)

	// Download blob content if the blob has been modified since we uploaded it (fails):
	downloadResp, err := blockBlob.DownloadStream(
		context.TODO(),
		&azblob.DownloadStreamOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{
					IfModifiedSince: upload.LastModified,
				},
			},
		},
	)
	showResult(&downloadResp, err)

	// Download blob content if the blob hasn't been modified in the last 24 hours (fails):
	downloadResp, err = blockBlob.DownloadStream(
		context.TODO(),
		&azblob.DownloadStreamOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{
					IfUnmodifiedSince: to.Ptr(time.Now().UTC().Add(time.Hour * -24))},
			},
		},
	)
	showResult(&downloadResp, err)

	// Upload new content if the blob hasn't changed since the version identified by ETag (succeeds):
	showResultUpload(blockBlob.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Text-2")),
		&blockblob.UploadOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: upload.ETag},
			},
		},
	))

	// Download content if it has changed since the version identified by ETag (fails):
	downloadResp, err = blockBlob.DownloadStream(
		context.TODO(),
		&azblob.DownloadStreamOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: upload.ETag}},
		})
	showResult(&downloadResp, err)

	// Upload content if the blob doesn't already exist (fails):
	showResultUpload(blockBlob.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Text-3")),
		&blockblob.UploadOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(string(azcore.ETagAny))},
			},
		}))
}

// This example shows how to create a blob with metadata, read blob metadata, and update a blob's read-only properties and metadata.
func Example_blob_Client_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blob with metadata (string key/value pairs)
	// Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	handleError(err)
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&blockblob.UploadOptions{Metadata: map[string]string{"author": "Jeffrey", "app": creatingApp}},
	)
	handleError(err)

	// Query the blob's properties and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	handleError(err)

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Show the blob's metadata
	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}

	for k, v := range get.Metadata {
		fmt.Print(k + "=" + v + "\n")
	}

	// Update the blob's metadata and write it back to the blob
	get.Metadata["editor"] = "Grant"
	_, err = blobClient.SetMetadata(context.TODO(), get.Metadata, nil)
	handleError(err)
}

// This examples shows how to create a blob with HTTP Headers, how to read, and how to update the blob's HTTP headers.
func Example_blob_HTTPHeaders() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blob with HTTP headers
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&blockblob.UploadOptions{HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType:        to.Ptr("text/html; charset=utf-8"),
			BlobContentDisposition: to.Ptr("attachment"),
		}},
	)
	handleError(err)

	// GetMetadata returns the blob's properties, HTTP headers, and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	handleError(err)

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Shows some of the blob's HTTP Headers
	httpHeaders := blob.ParseHTTPHeaders(get)
	fmt.Println(httpHeaders.BlobContentType, httpHeaders.BlobContentDisposition)

	// Update the blob's HTTP Headers and write them back to the blob
	httpHeaders.BlobContentType = to.Ptr("text/plain")
	_, err = blobClient.SetHTTPHeaders(context.TODO(), httpHeaders, nil)
	handleError(err)
}

// This example show how to create a blob, take a snapshot of it, update the base blob,
// read from the blob snapshot, list blobs with their snapshots, and delete blob snapshots.
func Example_blob_Client_CreateSnapshot() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	containerClient, err := container.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blockBlobClient object to a blob in the container.
	baseBlobClient := containerClient.NewBlockBlobClient("Original.txt")

	// Create the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(streaming.NopCloser(strings.NewReader("Some text"))), nil)
	handleError(err)

	// Create a snapshot of the original blob & save its timestamp:
	createSnapshot, err := baseBlobClient.CreateSnapshot(context.TODO(), nil)
	handleError(err)
	snapshot := *createSnapshot.Snapshot

	// Modify the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("New text")), nil)
	handleError(err)

	// Download the modified blob:
	get, err := baseBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	b := bytes.Buffer{}
	reader := get.Body
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(b.String())

	// Show snapshot blob via original blob URI & snapshot time:
	snapshotBlobClient, _ := baseBlobClient.WithSnapshot(snapshot)
	get, err = snapshotBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	b.Reset()
	reader = get.Body
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(b.String())

	// FYI: You can get the base blob URL from one of its snapshot by passing "" to WithSnapshot:
	baseBlobClient, _ = snapshotBlobClient.WithSnapshot("")

	// Show all blobs in the container with their snapshots:
	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	pager := containerClient.NewListBlobsFlatPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.Segment.BlobItems {
			// Process the blobs returned
			snapTime := "N/A"
			if blob.Snapshot != nil {
				snapTime = *blob.Snapshot
			}
			fmt.Printf("Blob name: %s, Snapshot: %s\n", *blob.Name, snapTime)
		}
	}

	// Promote read-only snapshot to writable base blob:
	_, err = baseBlobClient.StartCopyFromURL(context.TODO(), snapshotBlobClient.URL(), nil)
	handleError(err)

	// When calling Delete on a base blob:
	// DeleteSnapshotsOptionOnly deletes all the base blob's snapshots but not the base blob itself
	// DeleteSnapshotsOptionInclude deletes the base blob & all its snapshots.
	// DeleteSnapshotOptionNone produces an error if the base blob has any snapshots.
	_, err = baseBlobClient.Delete(context.TODO(), &blob.DeleteOptions{DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeInclude)})
	handleError(err)
}

func Example_progressUploadDownload() {
	// Create a credentials object with your Azure Storage Account name and key.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// From the Azure portal, get your Storage account blob service URL endpoint.
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)

	// Create an serviceClient object that wraps the service URL and a request pipeline to making requests.
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)
	handleError(err)

	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	blobClient := containerClient.NewBlockBlobClient("Data.bin")

	// requestBody is the stream of data to write
	requestBody := streaming.NopCloser(strings.NewReader("Some text to write"))

	// Wrap the request body in a RequestBodyProgress and pass a callback function for progress reporting.
	requestProgress := streaming.NewRequestProgress(streaming.NopCloser(requestBody), func(bytesTransferred int64) {
		fmt.Printf("Wrote %d of %d bytes.", bytesTransferred, requestBody)
	})
	_, err = blobClient.Upload(context.TODO(), requestProgress, &blockblob.UploadOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType:        to.Ptr("text/html; charset=utf-8"),
			BlobContentDisposition: to.Ptr("attachment"),
		},
	})
	handleError(err)

	// Here's how to read the blob's data with progress reporting:
	get, err := blobClient.DownloadStream(context.TODO(), nil)
	handleError(err)

	// Wrap the response body in a ResponseBodyProgress and pass a callback function for progress reporting.
	responseBody := streaming.NewResponseProgress(
		get.Body,
		func(bytesTransferred int64) {
			fmt.Printf("Read %d of %d bytes.", bytesTransferred, *get.ContentLength)
		},
	)

	downloadedData := &bytes.Buffer{}
	_, err = downloadedData.ReadFrom(responseBody)
	if err != nil {
		return
	}
	err = responseBody.Close()
	if err != nil {
		return
	}
	fmt.Printf("Downloaded data: %s\n", downloadedData.String())
}

// This example shows how to copy a source document on the Internet to a blob.
func Example_blob_Client_StartCopyFromURL() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a containerClient object to a container where we'll create a blob and its snapshot.
	// Create a blockBlobClient object to a blob in the container.
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/CopiedBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blob.NewClientWithSharedKeyCredential(blobURL, credential, nil)
	handleError(err)

	src := "https://cdn2.auth0.com/docs/media/addons/azure_blob.svg"
	startCopy, err := blobClient.StartCopyFromURL(context.TODO(), src, nil)
	handleError(err)

	copyID := *startCopy.CopyID
	copyStatus := *startCopy.CopyStatus
	for copyStatus == blob.CopyStatusTypePending {
		time.Sleep(time.Second * 2)
		getMetadata, err := blobClient.GetProperties(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		copyStatus = *getMetadata.CopyStatus
	}
	fmt.Printf("Copy from %s to %s: ID=%s, Status=%s\n", src, blobClient.URL(), copyID, copyStatus)
}

// This example shows how to download a large stream with intelligent retries. Specifically, if
// the connection fails while reading, continuing to read from this stream initiates a new
// GetBlob call passing a range that starts from the last byte successfully read before the failure.
func Example_blob_Client_Download() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blobClient object to a blob in the container (we assume the container & blob already exist).
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blob.NewClientWithSharedKeyCredential(blobURL, credential, nil)
	handleError(err)

	contentLength := int64(0) // Used for progress reporting to report the total number of bytes being downloaded.

	// Download returns an intelligent retryable stream around a blob; it returns an io.ReadCloser.
	dr, err := blobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	rs := dr.Body

	// NewResponseBodyProgress wraps the GetRetryStream with progress reporting; it returns an io.ReadCloser.
	stream := streaming.NewResponseProgress(
		rs,
		func(bytesTransferred int64) {
			fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
		},
	)
	defer func(stream io.ReadCloser) {
		err := stream.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stream) // The client must close the response body when finished with it

	file, err := os.Create("BigFile.bin") // Create the file to hold the downloaded blob contents.
	handleError(err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	written, err := io.Copy(file, stream) // Write to the file by reading from the blob (with intelligent retries).
	handleError(err)
	fmt.Printf("Wrote %d bytes.\n", written)
}
