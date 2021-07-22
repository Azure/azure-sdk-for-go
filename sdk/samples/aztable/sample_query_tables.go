package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

func Sample_QueryTables() {
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
	service := aztable.NewTableServiceClient(serviceURL, cred, nil)

	myTable := "myTableName"
	filter := fmt.Sprintf("TableName ge '%v'", myTable)
	pager := service.Query(&aztable.QueryOptions{Filter: &filter})

	pageCount := 1
	for pager.NextPage(context.Background()) {
		response := pager.PageResponse()
		fmt.Println("There are %d tables in page #%d", len(response.TableQueryResponse.Value), pageCount)
		pageCount += 1
	}

	check(pager.Error())
}

func main() {
	Sample_QueryTables()
}
