//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"io"
	"io/ioutil"
	"log"
	"crypto/rand"
	"os"
	"strings"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"

func generateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
	data := make([]byte, sizeInBytes)
	_len := len(random64BString)
	if sizeInBytes > _len {
		count := sizeInBytes / _len
		if sizeInBytes%_len != 0 {
			count = count + 1
		}
		copy(data[:], strings.Repeat(random64BString, count))
	} else {
		copy(data[:], random64BString)
	}
	return streaming.NopCloser(bytes.NewReader(data)), data
}

func Example_client_NewClient_CreateShare_CreateDir_CreateFile() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	serviceURL := fmt.Sprintf("https://%s.file.core.windows.net/", accountName)
	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	client, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	shareClient := client.NewShareClient("testShare")
	fmt.Println(shareClient.URL())

	dirClient := shareClient.NewDirectoryClient("testDir")
	fmt.Println(dirClient.URL())

	fileClient := dirClient.NewFileClient("testFile")
	fmt.Println(fileClient.URL())

}

func Example_file_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	filePath := "testDir/testFile"
	fileClient, err := file.NewClientFromConnectionString(connectionString, shareName, filePath, nil)
	handleError(err)
	fmt.Println(fileClient.URL())
}

func Example_fileClient_CreateAndDelete() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_GetProperties() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	_, err = fileClient.GetProperties(context.Background(), nil)
	handleError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)

}

func Example_fileClient_SetAndGetMetadata() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	opts := file.SetMetadataOptions{Metadata: map[string]*string{"hello": to.Ptr("world")}}
	_, err = fileClient.SetMetadata(context.Background(), &opts)
	handleError(err)

	get, err := fileClient.GetProperties(context.Background(), nil)
	handleError(err)

	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}
	for k, v := range get.Metadata {
		fmt.Print(k + "=" + *v + "\n")
	}

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_UploadBuffer() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	data := []byte{'h', 'e', 'l', 'l', 'o'}
	err = fileClient.UploadBuffer(context.Background(), data, nil)
	handleError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_UploadStream() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	err = fileClient.UploadStream(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		nil,
	)
	handleError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_UploadAndClearRange() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	fileName := "testFile"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 5, nil)
	handleError(err)

	contentR, _ := generateData(5)

	_, err = fileClient.UploadRange(context.Background(), 0, contentR, nil)
	handleError(err)

	rangeList, err := fileClient.GetRangeList(context.Background(), nil)
	handleError(err)
	fmt.Println(rangeList.Ranges)

	_, err = fileClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(5)}, nil)
	handleError(err)

	rangeList2, err := fileClient.GetRangeList(context.Background(), nil)
	handleError(err)

	fmt.Println(rangeList2.Ranges, 0)
	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_StartCopyFromURL() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	dstFileName := "testFile2"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	dstFileClient := shareClient.NewRootDirectoryClient().NewFileClient(dstFileName)

	contentR, _ := generateData(int(fileSize))

	_, err = srcFileClient.UploadRange(context.Background(), 0, contentR, nil)
	handleError(err)

	_, err = dstFileClient.StartCopyFromURL(context.Background(), srcFileClient.URL(), nil)
	handleError(err)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = dstFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_DownloadStream() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	contentR, _ := generateData(int(fileSize))

	_, err = srcFileClient.UploadRange(context.Background(), 0, contentR, nil)
	handleError(err)

	// validate data copied
	resp, err := srcFileClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 0, Count: fileSize},
	})
	handleError(err)

	content1, err := io.ReadAll(resp.Body)
	handleError(err)
	fmt.Println(content1)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_DownloadBuffer() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	handleError(err)

	err = srcFileClient.UploadBuffer(context.Background(), content, nil)
	handleError(err)

	destBuffer := make([]byte, fileSize)
	_, err = srcFileClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
	})
	handleError(err)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_DownloadFile() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	handleError(err)

	err = srcFileClient.UploadBuffer(context.Background(), content, nil)
	handleError(err)

	destFileName := "file.bin"
	destFile, err := os.Create(destFileName)
	handleError(err)
	defer func(name string) {
		err = os.Remove(name)
		handleError(err)
	}(destFileName)
	defer func(destFile *os.File) {
		err = destFile.Close()
		handleError(err)
	}(destFile)

	_, err = srcFileClient.DownloadFile(context.Background(), destFile, nil)
	handleError(err)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_UploadFile() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	_, content := generateData(int(fileSize))
	err = ioutil.WriteFile(srcFileName, content, 0644)
	handleError(err)
	defer func() {
		err = os.Remove(srcFileName)
		handleError(err)
	}()
	fh, err := os.Open(srcFileName)
	handleError(err)
	defer func(fh *os.File) {
		err := fh.Close()
		handleError(err)
	}(fh)

	err = srcFileClient.UploadFile(context.Background(), fh, nil)

	destFileName := "file.bin"
	destFile, err := os.Create(destFileName)
	handleError(err)
	defer func(name string) {
		err = os.Remove(name)
		handleError(err)
	}(destFileName)
	defer func(destFile *os.File) {
		err = destFile.Close()
		handleError(err)
	}(destFile)

	_, err = srcFileClient.DownloadFile(context.Background(), destFile, nil)
	handleError(err)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_file_ClientGetSASURL() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	srcFileName := "testFile"
	fileSize := int64(5)

	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	permission := sas.FilePermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	options := file.GetSASURLOptions{StartTime: &start}
	sasURL, err := srcFileClient.GetSASURL(permission, expiry, &options)
	handleError(err)
	_ = sasURL

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}
