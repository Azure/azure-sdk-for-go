//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

func handleError(err error) {
	if err != nil {
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
			Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
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
	pager := client.NewListBlobsFlatPager(containerName, nil)
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

func Example_client_NewClient() {
	// this example uses Azure Active Directory (AAD) to authenticate with Azure Blob Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	fmt.Println(client.URL())
}

func Example_client_NewClientWithSharedKeyCredential() {
	// this example uses a shared key to authenticate with Azure Blob Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// shared key authentication requires the storage account name and access key
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	serviceClient, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_client_NewClientFromConnectionString() {
	// this example uses a connection string to authenticate with Azure Blob Storage
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_client_anonymous_NewClientWithNoCredential() {
	// this example uses anonymous access to access a public blob
	serviceClient, err := azblob.NewClientWithNoCredential("https://azurestoragesamples.blob.core.windows.net/samples/cloud.jpg", nil)
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateContainer(context.TODO(), "testcontainer", &azblob.CreateContainerOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.DeleteContainer(context.TODO(), "testcontainer", nil)
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	pager := client.NewListContainersPager(&azblob.ListContainersOptions{
		Include: azblob.ListContainersInclude{Metadata: true, Deleted: true},
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Upload the file to a block blob
	_, err = client.UploadFile(context.TODO(), "testcontainer", "virtual/dir/path/"+fileName, fileHandler,
		&azblob.UploadFileOptions{
			BlockSize:   int64(1024),
			Concurrency: uint16(3),
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Perform download

	_, err = client.DownloadFile(context.TODO(), "testcontainer", "virtual/dir/path/"+destFileName, destFile,
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	pager := client.NewListBlobsFlatPager("testcontainer", &azblob.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Deleted: true, Versions: true},
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.DeleteBlob(context.TODO(), "testcontainer", "testblob", nil)
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Set up test blob
	containerName := "testcontainer"
	bufferSize := 8 * 1024 * 1024
	blobName := "test_upload_stream.bin"
	blobData := make([]byte, bufferSize)
	blobContentReader := bytes.NewReader(blobData)

	// Perform UploadStream
	resp, err := client.UploadStream(context.TODO(), containerName, blobName, blobContentReader,
		&azblob.UploadStreamOptions{
			Metadata: map[string]*string{"hello": to.Ptr("world")},
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

	client, err := azblob.NewClient(serviceURL, cred, nil)
	handleError(err)

	// Download the blob
	downloadResponse, err := client.DownloadStream(ctx, "testcontainer", "test_download_stream.bin", nil)
	handleError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.Body)
	handleError(err)
	fmt.Println(len(actualBlobData))
}

// ---------------------------------------------------------------------------------------------------------------------

func ExampleResponseError() {
	contClient, err := container.NewClientWithNoCredential("https://myaccount.blob.core.windows.net/mycontainer", nil)
	handleError(err)
	_, err = contClient.Create(context.TODO(), nil)
	handleError(err)
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
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(azcore.ETagAny)},
			},
		}))
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
