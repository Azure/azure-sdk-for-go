package main

import (
	"os"
	// "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateTableServiceClient() *aztable.NewTableServiceClient {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred := aztable.SharedKeyCredential(accountName, accountKey)
	return &aztable.NewTableServiceClient(serviceURL, cred, nil)
}

func CreateTableClient() *aztable.TableClient {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred := aztable.SharedKeyCredential(accountName, accountKey)
	return &aztable.NewTableClient("tableName", serviceURL, cred, nil)
}

func main() {
	CreateTableClient()
	CreateTableServiceClient
}
