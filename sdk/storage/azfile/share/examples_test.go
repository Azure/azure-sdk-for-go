//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"log"
	"os"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_share_Client_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	fmt.Println(shareClient.URL())
}

func Example_share_Client_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	shareName := "testshare"
	shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
	handleError(err)

	fmt.Println(shareClient.URL())
}

func Example_share_Client_NewDirectoryClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	dirName := "testdirectory"
	dirClient := shareClient.NewDirectoryClient(dirName)

	fmt.Println(dirClient.URL())
}

func Example_share_Client_NewRootDirectoryClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	dirClient := shareClient.NewRootDirectoryClient()

	fmt.Println(dirClient.URL())
}

func Example_share_Client_CreateSnapshot() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	snapResp, err := shareClient.CreateSnapshot(context.TODO(), nil)
	handleError(err)
	shareSnapshot := *snapResp.Snapshot

	snapshotShareClient, err := shareClient.WithSnapshot(shareSnapshot)
	handleError(err)

	fmt.Println(snapshotShareClient.URL())

	_, err = snapshotShareClient.GetProperties(context.TODO(), nil)
	handleError(err)
}

func Example_share_Client_Create() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	_, err = shareClient.Create(context.TODO(), nil)
	handleError(err)
}

func Example_share_Client_Delete() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	_, err = shareClient.Delete(context.TODO(), nil)
	handleError(err)
}

func Example_share_Client_Restore() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	serviceURL := fmt.Sprintf("https://%s.file.core.windows.net/", accountName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	shareClient := svcClient.NewShareClient(shareName)

	// get share version for restore operation
	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true}, // Include deleted shares in the result
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				_, err = shareClient.Restore(context.TODO(), *s.Version, nil)
				handleError(err)
			}
		}
	}
}

func Example_share_Client_GetProperties() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	_, err = shareClient.GetProperties(context.TODO(), nil)
	handleError(err)
}

func Example_share_Client_SetProperties() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	_, err = shareClient.SetProperties(context.TODO(), &share.SetPropertiesOptions{
		Quota:      to.Ptr(int32(1000)),
		AccessTier: to.Ptr(share.AccessTierHot),
	})
	handleError(err)
}

func Example_share_Client_AccessPolicy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	permission := share.AccessPolicyPermission{Read: true, Write: true, Create: true, Delete: true, List: true}.String()
	permissions := []*share.SignedIdentifier{
		{
			ID: to.Ptr("1"),
			AccessPolicy: &share.AccessPolicy{
				Start:      to.Ptr(time.Now()),
				Expiry:     to.Ptr(time.Now().Add(time.Hour)),
				Permission: &permission,
			},
		}}

	_, err = shareClient.SetAccessPolicy(context.TODO(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	handleError(err)

	resp, err := shareClient.GetAccessPolicy(context.TODO(), nil)
	handleError(err)

	fmt.Println(*resp.SignedIdentifiers[0].AccessPolicy.Permission)
}

func Example_share_Client_CreateGetPermission() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	testSDDL := `O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)`
	createResp, err := shareClient.CreatePermission(context.TODO(), testSDDL, nil)
	handleError(err)

	getResp, err := shareClient.GetPermission(context.TODO(), *createResp.FilePermissionKey, nil)
	handleError(err)
	fmt.Println(*getResp.Permission)
}

func Example_share_Client_SetMetadata() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}
	_, err = shareClient.SetMetadata(context.TODO(), &share.SetMetadataOptions{
		Metadata: md,
	})
	handleError(err)

	resp, err := shareClient.GetProperties(context.TODO(), nil)
	handleError(err)
	for k, v := range resp.Metadata {
		fmt.Printf("%v : %v\n", k, *v)
	}
}

func Example_share_Client_GetStatistics() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	getStats, err := shareClient.GetStatistics(context.Background(), nil)
	handleError(err)
	fmt.Println(*getStats.ShareUsageBytes)
}

func Example_share_Client_GetSASURL() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	handleError(err)

	permissions := sas.SharePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	shareSASURL, err := shareClient.GetSASURL(permissions, expiry, nil)
	handleError(err)

	fmt.Println("SAS URL: ", shareSASURL)

	shareSASClient, err := share.NewClientWithNoCredential(shareSASURL, nil)
	handleError(err)

	var dirs, files []string
	pager := shareSASClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)

		for _, d := range resp.Segment.Directories {
			dirs = append(dirs, *d.Name)
		}
		for _, f := range resp.Segment.Files {
			files = append(files, *f.Name)
		}
	}

	fmt.Println("Directories:")
	for _, d := range dirs {
		fmt.Println(d)
	}

	fmt.Println("Files:")
	for _, f := range files {
		fmt.Println(f)
	}
}

func Example_share_Client_CreateAndGetPermissionOAuth() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	shareName := "testshare"
	shareURL := fmt.Sprintf("https://%s.file.core.windows.net/%s", accountName, shareName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	// FileRequestintent is required if authorization header specifies an OAuth token.
	options := &share.ClientOptions{FileRequestIntent: to.Ptr(share.TokenIntentBackup)}
	shareClient, err := share.NewClient(shareURL, cred, options)
	handleError(err)

	// Create, Delete, GetProperties, SetProperties, etc. operations does not work when share client is created using OAuth credentials
	// Operations supported are: CreatePermission and GetPermission
	// Below GetProperties operation results in an error
	_, err = shareClient.GetProperties(context.TODO(), nil)
	fmt.Println(err.Error())

	testSDDL := `O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)`
	createResp, err := shareClient.CreatePermission(context.TODO(), testSDDL, nil)
	handleError(err)

	getResp, err := shareClient.GetPermission(context.TODO(), *createResp.FilePermissionKey, nil)
	handleError(err)
	fmt.Println(*getResp.Permission)
}
