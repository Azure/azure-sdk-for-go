package main

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	// "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

func CreateFromTableServiceClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"
	service := aztable.NewTableServiceClient(serviceURL, cred, nil)
	_, err = service.Create(context.Background(), "fromServiceClient")
	check(err)
}

func CreateFromTableClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"
	client := aztable.NewTableClient("fromTableClient", serviceURL, cred, nil)
	_, err = client.Create(context.Background())
	check(err)
}

func DeleteFromTableServiceClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"
	service := aztable.NewTableServiceClient(serviceURL, cred, nil)
	_, err = service.Delete(context.Background(), "fromServiceClient")
	check(err)
}

func DeleteFromTableClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"
	client := aztable.NewTableClient("fromTableClient", serviceURL, cred, nil)
	_, err = client.Delete(context.Background())
	check(err)
}

func main() {
	CreateFromTableServiceClient()
	CreateFromTableClient()
	DeleteFromTableServiceClient()
	DeleteFromTableClient()
}
