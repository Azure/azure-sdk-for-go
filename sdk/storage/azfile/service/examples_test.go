// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"log"
	"os"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_service_Client_NewClient() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	fmt.Println(svcClient.URL())
}

func Example_service_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	svcClient, err := service.NewClientFromConnectionString(connectionString, nil)
	handleError(err)

	fmt.Println(svcClient.URL())
}

func Example_service_Client_NewShareClient() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	shareName := "testShare"
	shareClient := svcClient.NewShareClient(shareName)

	fmt.Println(shareClient.URL())
}

func Example_service_Client_CreateShare() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	shareName := "testShare"
	_, err = svcClient.CreateShare(context.TODO(), shareName, nil)
	handleError(err)
	fmt.Println("Share created")
}

func Example_service_Client_DeleteShare() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	shareName := "testShare"
	_, err = svcClient.DeleteShare(context.TODO(), shareName, nil)
	handleError(err)
	fmt.Println("Share deleted")
}

func Example_service_Client_RestoreShare() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	// get share version for restore operation
	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true}, // Include deleted shares in the result
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				_, err = svcClient.RestoreShare(context.TODO(), *s.Name, *s.Version, nil)
				handleError(err)
			}
		}
	}
}

func Example_service_Client_GetProperties() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	_, err = svcClient.GetProperties(context.TODO(), nil)
	handleError(err)
}

func Example_service_Client_SetProperties() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	setPropertiesOpts := service.SetPropertiesOptions{
		HourMetrics: &service.Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(true),
			RetentionPolicy: &service.RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(2)),
			},
		},
		MinuteMetrics: &service.Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(false),
			RetentionPolicy: &service.RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(2)),
			},
		},
		CORS: []*service.CORSRule{
			{
				AllowedOrigins:  to.Ptr("*"),
				AllowedMethods:  to.Ptr("PUT"),
				AllowedHeaders:  to.Ptr("x-ms-client-request-id"),
				ExposedHeaders:  to.Ptr("x-ms-*"),
				MaxAgeInSeconds: to.Ptr(int32(2)),
			},
		},
	}
	_, err = svcClient.SetProperties(context.TODO(), &setPropertiesOpts)
	handleError(err)
}

func Example_service_Client_ListShares() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	pager := svcClient.NewListSharesPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, s := range resp.Shares {
			fmt.Println(*s.Name)
		}
	}
}

func Example_service_Client_GetSASURL() {
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

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)
	sasUrl, err := svcClient.GetSASURL(resources, permissions, expiry, nil)
	handleError(err)

	fmt.Println("SAS URL: ", sasUrl)

	svcSASClient, err := service.NewClientWithNoCredential(sasUrl, nil)
	handleError(err)

	_, err = svcSASClient.GetProperties(context.TODO(), nil)
	handleError(err)
}
