package main

import (
	"os"
)

func CreateServiceClient() *aztable.NewServiceClient {
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
	return &aztable.NewServiceClient(serviceURL, cred, nil)
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
