//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_userDelegationSAS() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	tenantID, ok := os.LookupEnv("AZURE_TENANT_ID")
	if !ok {
		panic("AZURE_TENANT_ID could not be found")
	}
	clientID, ok := os.LookupEnv("AZURE_CLIENT_ID")
	if !ok {
		panic("AZURE_CLIENT_ID could not be found")
	}
	clientSecret, ok := os.LookupEnv("AZURE_CLIENT_SECRET")
	if !ok {
		panic("AZURE_CLIENT_SECRET could not be found")
	}
	const containerName = "testcontainer"

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	handleError(err)

	svcClient, err := service.NewClient(
		fmt.Sprintf("https://%s.blob.core.windows.net/", accountName),
		cred,
		&service.ClientOptions{},
	)
	handleError(err)

	// Set current and past time and create key
	now := time.Now().UTC().Add(-10 * time.Second)
	expiry := now.Add(48 * time.Hour)
	info := service.KeyInfo{
		Start:  to.Ptr(now.UTC().Format(sas.TimeFormat)),
		Expiry: to.Ptr(expiry.UTC().Format(sas.TimeFormat)),
	}

	udc, err := svcClient.GetUserDelegationCredential(context.Background(), info, nil)
	handleError(err)

	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC().Add(time.Second * -10),
		ExpiryTime:    time.Now().UTC().Add(15 * time.Minute),
		Permissions:   to.Ptr(sas.ContainerPermissions{Read: true, List: true}).String(),
		ContainerName: containerName,
	}.SignWithUserDelegation(udc)
	handleError(err)

	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, sasQueryParams.Encode())

	// This URL can be used to authenticate requests now
	azClient, err := azblob.NewClientWithNoCredential(sasURL, nil)
	handleError(err)

	// list blobs in container
	pager := azClient.NewListBlobsFlatPager(containerName, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, b := range resp.Segment.BlobItems {
			fmt.Println(*b.Name)
		}
	}

	// User Delegation SAS doesn't support operations like creation, deletion or listing of containers
	// For more details, see https://docs.microsoft.com/rest/api/storageservices/create-user-delegation-sas#specify-permissions
	_, err = azClient.CreateContainer(context.Background(), "newcontainer", nil)
	if err != nil {
		fmt.Println("Containers can't be created using User Delegation SAS")
	}

	_, err = azClient.DeleteContainer(context.Background(), containerName, nil)
	if err != nil {
		fmt.Println("Containers can't be deleted using User Delegation SAS")
	}
}

func Example_serviceSAS() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	const containerName = "testContainer"

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true, Create: true, Write: true, Tag: true}).String(),
		ContainerName: containerName,
	}.SignWithSharedKey(credential)
	handleError(err)

	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, sasQueryParams.Encode())
	fmt.Println(sasURL)

	// This URL can be used to authenticate requests now
	azClient, err := azblob.NewClientWithNoCredential(sasURL, nil)
	handleError(err)

	const blobData, blobName = "test data", "testBlob"
	uploadResp, err := azClient.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
			Tags:     map[string]string{"Year": "2022"},
		})
	handleError(err)
	fmt.Println(uploadResp)

	blobDownloadResponse, err := azClient.DownloadStream(context.TODO(), containerName, blobName, nil)
	handleError(err)

	reader := blobDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	handleError(err)
	fmt.Println(string(downloadData))
	if string(downloadData) != blobData {
		log.Fatal("Uploaded data should be same as downloaded data")
	}

	err = reader.Close()
	if err != nil {
		return
	}
}
