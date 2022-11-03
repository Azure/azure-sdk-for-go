package sas_test

import (
	"context"
	"fmt"
	"log"
	"os"
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
	clientID, ok := os.LookupEnv("AZURE_STORAGE_SPN_CLIENT_ID")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	tenantID, ok := os.LookupEnv("AZURE_STORAGE_SPN_TENANT_ID")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	clientSecret, ok := os.LookupEnv("AZURE_STORAGE_SPN_CLIENT_SECRET")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"

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

	// Create Account Signature Values with desired permissions and sign with user delegation credential
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
