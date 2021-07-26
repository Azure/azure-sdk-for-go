package main

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AuthenticateWithSharedKey() *aztable.TableServiceClient {
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
	client := aztable.NewTableServiceClient(serviceURL, cred, nil)
	return client
}

func AuthenticateWithTokenCredential() *aztable.TableServiceClient {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	client := aztable.NewTableServiceClient(serviceURL, cred, nil)
	return client
}

func AuthenticateWithSharedAccess() {
	// TODO
}

func main() {
	AuthenticateWithTokenCredential()
	AuthenticateWithSharedKey()
	AuthenticateWithSharedAccess()

	// sample_batch.go
	Sample_Batching()

	// sample_authentication.go
	CreateTableClient()
	CreateTableServiceClient()

	// sample_create_delete_table.go
	CreateFromTableServiceClient()
	CreateFromTableClient()
	DeleteFromTableServiceClient()
	DeleteFromTableClient()

	// sample_insert_delete_entities.go
	InsertEntity()
	DeleteEntity()

	// sample_query_table.go
	Sample_QueryTable()

	// sample_query_tables.go
	Sample_QueryTables()

	// sample_update_entities.go
	UpdateEntities()
	MergeEntities()
	UpsertEntity()
}
