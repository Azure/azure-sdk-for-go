//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
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
	_, err = pageBlobClient.UploadPages(context.TODO(), streaming.NopCloser(bytes.NewReader(page)), blob.HTTPRange{
		Offset: 0,
		Count:  0,
	}, nil)
	handleError(err)

	copy(page, "Page 1")
	_, err = pageBlobClient.UploadPages(
		context.TODO(),
		streaming.NopCloser(bytes.NewReader(page)), blob.HTTPRange{
			Count: int64(2 * pageblob.PageBytes),
		}, nil)
	handleError(err)

	pager := pageBlobClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Range: blob.HTTPRange{
			Count: int64(10 * pageblob.PageBytes),
		},
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

	_, err = pageBlobClient.ClearPages(context.TODO(), blob.HTTPRange{Count: 1 * pageblob.PageBytes}, nil)
	handleError(err)

	pager = pageBlobClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Range: blob.HTTPRange{
			Count: int64(10 * pageblob.PageBytes),
		},
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
