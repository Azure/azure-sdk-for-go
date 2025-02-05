//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"io"
	"log"
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
			count++
		}
		copy(data, strings.Repeat(random64BString, count))
	} else {
		copy(data, random64BString)
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

	// you can also use AbortCopy to abort copying
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
	_, err = srcFileClient.DownloadBuffer(context.Background(), destBuffer, nil)
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
	err = os.WriteFile(srcFileName, content, 0644)
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

func Example_fileClient_Resize() {
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

	resp1, err := srcFileClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*resp1.ContentLength)

	_, err = srcFileClient.Resize(context.Background(), 6, nil)
	handleError(err)

	resp1, err = srcFileClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*resp1.ContentLength)

	_, err = srcFileClient.Delete(context.Background(), nil)
	handleError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	handleError(err)
}

func Example_fileClient_UploadRangeFromURL() {
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

	shareName := "testShare"
	srcFileName := "testFile"
	dstFileName := "testFile2"
	fileSize := int64(5)

	shareClient := client.NewShareClient(shareName)
	_, err = shareClient.Create(context.Background(), nil)
	handleError(err)

	srcFileClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	contentR, _ := generateData(int(fileSize))

	_, err = srcFileClient.UploadRange(context.Background(), 0, contentR, nil)
	handleError(err)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = srcFileClient.UploadRange(context.Background(), 0, rsc, nil)
	handleError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	handleError(err)

	srcFileSAS := srcFileClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := shareClient.NewRootDirectoryClient().NewFileClient(dstFileName)
	_, err = destFClient.Create(context.Background(), fileSize, nil)
	handleError(err)

	_, err = destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), nil)
	handleError(err)
}

func Example_fileClient_OAuth() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	fileName := "testFile"
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	fileClient, err := file.NewClient(fileURL, cred, &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)})
	handleError(err)

	_, err = fileClient.Create(context.TODO(), 2048, nil)
	handleError(err)
	fmt.Println("File created")

	_, err = fileClient.GetProperties(context.TODO(), nil)
	handleError(err)
	fmt.Println("File properties retrieved")

	_, err = fileClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println("File deleted")
}

func Example_fileClient_TrailingDot() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	fileName := "testFile.." // file name with trailing dot
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	fileClient, err := file.NewClient(fileURL, cred, &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		AllowTrailingDot:  to.Ptr(true),
	})
	handleError(err)

	_, err = fileClient.Create(context.TODO(), 2048, nil)
	handleError(err)
	fmt.Println("File created")

	_, err = fileClient.GetProperties(context.TODO(), nil)
	handleError(err)
	fmt.Println("File properties retrieved")

	_, err = fileClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println("File deleted")
}

func Example_fileClient_Rename() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	srcFileName := "testFile"
	destFileName := "newFile"
	srcFileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + srcFileName

	srcFileClient, err := file.NewClient(srcFileURL, cred, &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)})
	handleError(err)

	_, err = srcFileClient.Rename(context.TODO(), destFileName, nil)
	handleError(err)
	fmt.Println("File renamed")
}

func Example_fileClient_CopyFileUsingSourceProperties() {
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

	_, err = dstFileClient.StartCopyFromURL(context.Background(), srcFileClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:       file.SourceCopyFileCreationTime{},
			LastWriteTime:      file.SourceCopyFileLastWriteTime{},
			ChangeTime:         file.SourceCopyFileChangeTime{},
			Attributes:         file.SourceCopyFileAttributes{},
			PermissionCopyMode: to.Ptr(file.PermissionCopyModeTypeSource),
		},
	})
	handleError(err)
	fmt.Println("File copied")
}

func Example_fileClient_CopyFileUsingDestinationProperties() {
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

	destCreationTime := time.Now().Add(5 * time.Minute)
	destLastWriteTIme := time.Now().Add(6 * time.Minute)
	destChangeTime := time.Now().Add(7 * time.Minute)
	_, err = dstFileClient.StartCopyFromURL(context.Background(), srcFileClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:  file.DestinationCopyFileCreationTime(destCreationTime),
			LastWriteTime: file.DestinationCopyFileLastWriteTime(destLastWriteTIme),
			ChangeTime:    file.DestinationCopyFileChangeTime(destChangeTime),
			Attributes:    file.DestinationCopyFileAttributes{ReadOnly: true},
		},
	})
	handleError(err)
	fmt.Println("File copied")
}

func Example_fileClient_CreateNFSShare() {

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")

	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareName := "testShare"
	shareURL := "https://" + accountName + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	fileMode := "7777"

	options := &share.ClientOptions{}
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	handleError(err)

	fClient := premiumShareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(shareName))
	_, err = fClient.Create(context.Background(), 1024, &file.CreateOptions{
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(fileMode),
	})
	handleError(err)
	fmt.Println("NFS Share created with given properties")
}

func Example_fileClient_SetHttpHeadersNFS() {

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")

	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareName := "testShare"
	shareURL := "https://" + accountName + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	fileMode := "7777"

	options := &share.ClientOptions{}
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	handleError(err)
	fClient := premiumShareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(shareName))
	_, err = fClient.Create(context.Background(), 0, nil)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	opts := &file.SetHTTPHeadersOptions{
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(fileMode),
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), opts)

	handleError(err)
	fmt.Println("Properties set on NFS share")
}

func Example_fileClient_NFS_CreateHardLink() {

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")

	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareName := "testShare"
	shareURL := "https://" + accountName + ".file.core.windows.net/" + shareName

	options := &share.ClientOptions{}
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	handleError(err)
	directoryName := testcommon.GenerateDirectoryName("dirName")
	directoryClient := premiumShareClient.NewRootDirectoryClient().NewSubdirectoryClient(directoryName)
	_, err = directoryClient.Create(context.Background(), nil)

	// Create a source file
	sourceFileName := testcommon.GenerateFileName("file1")
	sourceFileClient := directoryClient.NewFileClient(sourceFileName)
	_, err = sourceFileClient.Create(context.Background(), int64(1024), nil)

	// Create a hard link to the source file
	hardLinkFileName := testcommon.GenerateFileName("file2")
	hardLinkFileClient := directoryClient.NewFileClient(hardLinkFileName)

	targetFilePath := fmt.Sprintf("/%s/%s", directoryName, sourceFileName)
	_, err = hardLinkFileClient.CreateHardLink(context.Background(), &file.CreateHardLinkOptions{
		TargetFile: targetFilePath,
	})

	handleError(err)
	fmt.Println("Hard Link created - Link count = 2")
}
