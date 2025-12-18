//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"hash/crc64"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// make sure you create the filesystem before running this example
func Example_file_CreateAndDelete() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fileClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = fileClient.Create(context.Background(), nil)
	handleError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	handleError(err)
}

// This examples shows how to create a file with HTTP Headers, how to read, and how to update the file's HTTP headers.
// make sure you create the filesystem and file before running this example.
func Example_file_HTTPHeaders() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fileClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a directory with HTTP headers
	_, err = fileClient.SetHTTPHeaders(context.TODO(), file.HTTPHeaders{
		ContentType:        to.Ptr("text/html; charset=utf-8"),
		ContentDisposition: to.Ptr("attachment"),
	}, nil)
	handleError(err)

	get, err := fileClient.GetProperties(context.TODO(), nil)
	handleError(err)

	fmt.Println(get.ContentType)
	fmt.Println(get.ContentDisposition)
}

// make sure you create the filesystem before running this example
func Example_file_CreateFileWithExpiryRelativeToNow() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeRelativeToNow,
			ExpiresOn:  strconv.FormatInt((8 * time.Second).Milliseconds(), 10),
		},
	}

	_, err = fClient.Create(context.Background(), createFileOpts)
	handleError(err)

	resp, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*resp.ExpiresOn)

	time.Sleep(time.Second * 10)
	_, err = fClient.GetProperties(context.Background(), nil)
	// we expect datalakeerror.PathNotFound
	handleError(err)
}

// make sure you create the filesystem before running this example
func Example_file_CreateFileWithNeverExpire() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeNeverExpire,
		},
	}

	_, err = fClient.Create(context.Background(), createFileOpts)
	handleError(err)

	resp, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	// should be empty since we never expire
	fmt.Println(*resp.ExpiresOn)
}

// make sure you create the filesystem and file before running this example
func Example_file_Client_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	fileClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = fileClient.SetMetadata(context.TODO(), map[string]*string{"author": to.Ptr("Tamer")}, nil)
	handleError(err)

	// Query the directory's properties and metadata
	get, err := fileClient.GetProperties(context.TODO(), nil)
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
func Example_file_Rename() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fileClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = fileClient.Create(context.Background(), nil)
	handleError(err)

	_, err = fileClient.Rename(context.Background(), "renameFile", nil)
	handleError(err)
}

// set acl on a file
// make sure you create the filesystem and file before running this example
func Example_file_SetACL() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	acl := "user::rwx,group::r-x,other::rwx"
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fileClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	_, err = fileClient.SetAccessControl(context.Background(), &file.SetAccessControlOptions{
		ACL: &acl,
	})
	handleError(err)
}

func getRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

// make sure you create the filesystem before running this example
func Example_file_SetAccessControlIfUnmodifiedSinceTrue() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a directory client
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fileClient, err := directory.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)
	resp, err := fileClient.Create(context.Background(), nil)
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

	_, err = fileClient.SetAccessControl(context.Background(), opts)
	handleError(err)
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

// make sure you create the filesystem before running this example
func Example_file_UploadFileAndDownloadStream() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)
	var fileSize int64 = 10 * 1024

	_, err = fClient.Create(context.Background(), nil)
	handleError(err)

	// create local file
	_, content := generateData(int(fileSize))
	err = os.WriteFile("testFile", content, 0644)
	handleError(err)

	defer func() {
		err = os.Remove("testFile")
		handleError(err)
	}()

	fh, err := os.Open("testFile")
	handleError(err)

	defer func(fh *os.File) {
		err := fh.Close()
		handleError(err)
	}(fh)

	// get md5 hash to compare against after download
	hash := md5.New()
	_, err = io.Copy(hash, fh)
	handleError(err)
	contentMD5 := hash.Sum(nil)

	// upload the file
	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	handleError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	handleError(err)

	data, err := io.ReadAll(dResp.Body)
	handleError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	// compare the hashes
	fmt.Println(downloadedContentMD5, contentMD5)
}

// make sure you create the filesystem before running this example
func Example_file_UploadBufferAndDownloadStream() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)
	var fileSize int64 = 10 * 1024

	_, content := generateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	handleError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	handleError(err)

	data, err := io.ReadAll(dResp.Body)
	handleError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	fmt.Println(downloadedContentMD5, contentMD5)
}

// make sure you create the filesystem before running this example
func Example_file_AppendAndFlushDataWithValidation() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	// validate data using crc64
	opts := &file.AppendDataOptions{
		TransactionalValidation: file.TransferValidationTypeComputeCRC64(),
	}
	putResp, err := fClient.AppendData(context.Background(), 0, rsc, opts)
	handleError(err)
	fmt.Println(putResp.ContentCRC64)
	fmt.Println(binary.LittleEndian.Uint64(putResp.ContentCRC64), contentCRC64)

	// after appending data, flush it
	_, err = fClient.FlushData(context.Background(), int64(contentSize), nil)
	handleError(err)

	// compare content length as well
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println(*gResp2.ContentLength, int64(contentSize))
}

func Example_file_AppendAndFlushDataWithAcquireAndReleaseLease() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a file client
	u := fmt.Sprintf("https://%s.dfs.core.windows.net/fs/file.txt", accountName)
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	fClient, err := file.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	// Acquire lease during append data
	opts := &file.AppendDataOptions{
		LeaseAction:     &file.LeaseActionAcquire,
		LeaseDuration:   to.Ptr(int64(15)),
		ProposedLeaseID: proposedLeaseIDs[1],
	}
	_, err = fClient.AppendData(context.Background(), 0, rsc, opts)
	handleError(err)

	_, err = fClient.FlushData(context.Background(), int64(contentSize), &file.FlushDataOptions{
		LeaseAction: &file.LeaseActionRelease,
		AccessConditions: &path.AccessConditions{
			LeaseAccessConditions: &path.LeaseAccessConditions{LeaseID: proposedLeaseIDs[0]},
		},
	})
	handleError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	handleError(err)
	// Check if the lease is released
	fmt.Println(lease.StateTypeAvailable, *gResp2.LeaseState)

}
