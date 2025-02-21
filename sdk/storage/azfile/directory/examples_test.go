//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
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

func Example_client_NewClient() {
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

	dirClient := shareClient.NewDirectoryClient("testDir")
	fmt.Println(dirClient.URL())

}

func Example_directory_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	dirName := "testDirectory"
	dirClient, err := directory.NewClientFromConnectionString(connectionString, shareName, dirName, nil)
	handleError(err)
	fmt.Println(dirClient.URL())
}

func Example_directoryClient_Create() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	dirName := "testDirectory"
	dirClient, err := directory.NewClientFromConnectionString(connectionString, shareName, dirName, nil)
	handleError(err)
	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory created")

	_, err = dirClient.Delete(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory deleted")
}

func Example_directoryClient_SetProperties() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	dirName := "testDirectory"
	dirClient, err := directory.NewClientFromConnectionString(connectionString, shareName, dirName, nil)
	handleError(err)
	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory created")

	creationTime := time.Now().Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := time.Now().Add(10 * time.Minute).Round(time.Millisecond)
	testSDDL := `O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)`

	// Set the custom permissions
	_, err = dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
		},
		FilePermissions: &file.Permissions{
			Permission: &testSDDL,
		},
	})
	handleError(err)
	fmt.Println("Directory properties set")

	_, err = dirClient.GetProperties(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory properties retrieved")

	_, err = dirClient.Delete(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory deleted")
}

func Example_directoryClient_ListFilesAndDirectoriesSegment() {
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	parentDirName := "testParentDirectory"
	childDirName := "testChildDirectory"
	parentDirClient, err := directory.NewClientFromConnectionString(connectionString, shareName, parentDirName, nil)
	handleError(err)
	_, err = parentDirClient.Create(context.Background(), nil)
	handleError(err)
	fmt.Println("Parent directory created")

	childDirClient := parentDirClient.NewSubdirectoryClient(childDirName)
	_, err = childDirClient.Create(context.Background(), nil)
	handleError(err)
	fmt.Println("Child directory created")

	pager := parentDirClient.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err) // if err is not nil, break the loop.
		for _, _dir := range resp.Segment.Directories {
			fmt.Printf("%v", _dir)
		}
	}

	_, err = childDirClient.Delete(context.Background(), nil)
	handleError(err)
	fmt.Println("Child directory deleted")

	_, err = parentDirClient.Delete(context.Background(), nil)
	handleError(err)
	fmt.Println("Parent directory deleted")
}

func Example_directoryClient_SetMetadata() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	shareName := "testShare"
	dirName := "testDirectory"
	dirClient, err := directory.NewClientFromConnectionString(connectionString, shareName, dirName, nil)
	handleError(err)
	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory created")

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	_, err = dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md,
	})
	handleError(err)
	fmt.Println("Directory metadata set")

	_, err = dirClient.Delete(context.Background(), nil)
	handleError(err)
	fmt.Println("Directory deleted")
}

func Example_directoryClient_OAuth() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	dirName := "testDirectory"
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	dirClient, err := directory.NewClient(dirURL, cred, &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)})
	handleError(err)

	_, err = dirClient.Create(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory created")

	_, err = dirClient.GetProperties(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory properties retrieved")

	_, err = dirClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory deleted")
}

func Example_directoryClient_TrailingDot() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	dirName := "testDirectory.." // directory name with trailing dot
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	dirClient, err := directory.NewClient(dirURL, cred, &directory.ClientOptions{
		FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup),
		AllowTrailingDot:  to.Ptr(true),
	})
	handleError(err)

	fmt.Println(dirClient.URL())

	_, err = dirClient.Create(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory created")

	_, err = dirClient.GetProperties(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory properties retrieved")

	_, err = dirClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println("Directory deleted")
}

func Example_directoryClient_Rename() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	shareName := "testShare"
	srcDirName := "testDirectory"
	destDirName := "newDirectory"
	srcDirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + srcDirName

	srcDirClient, err := directory.NewClient(srcDirURL, cred, &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)})
	handleError(err)

	_, err = srcDirClient.Rename(context.TODO(), destDirName, nil)
	handleError(err)
	fmt.Println("Directory renamed")
}

func Example_client_CreateDirectoryNFS() {
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

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	handleError(err)

	dirName := testcommon.GenerateDirectoryName("Dir1")
	dirClient := shareClient.NewDirectoryClient(dirName)

	_, err = dirClient.Create(context.Background(), &directory.CreateOptions{
		Owner:    to.Ptr("123"),
		Group:    to.Ptr("345"),
		FileMode: to.Ptr("7774"),
	})
	handleError(err)
}

func Example_client_SetHTTPHeadersNFS() {
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

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	handleError(err)

	dirName := testcommon.GenerateDirectoryName("Dir1")
	dirClient := shareClient.NewDirectoryClient(dirName)
	_, err = dirClient.Create(context.Background(), nil)
	handleError(err)
	owner := "345"
	group := "123"
	mode := "7777"

	// Set the custom permissions
	_, err = dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(mode),
	})
	handleError(err)
}
