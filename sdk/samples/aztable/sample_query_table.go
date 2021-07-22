package main

import (
	"context"
	"fmt"
	"os"
	// "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

func Sample_QueryTable() {
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
	client := aztable.NewTableClient("myTable", serviceURL, cred, nil)

	filter := fmt.Sprintf("PartitionKey eq '%v' or PartitionKey eq '%v'", "pk001", "pk002")
	pager := client.Query(&aztable.QueryOptions{Filter: &filter})

	pageCount := 1
	for pager.NextPage(context.Background()) {
		response := pager.PageResponse()
		fmt.Println("There are %d entities in page #%d", len(response.TableEntityQueryResponse.Value), pageCount)
		pageCount += 1
	}
	check(pager.Error())

	// To list all entities in a table, provide nil to Query()
	listPager := client.Query(nil)
	pageCount = 1
	for listPager.NextPage(context.Background()) {
		response := listPager.PageResponse()
		fmt.Println("There are %d entities in page #%d", len(response.TableEntityQueryResponse.Value), pageCount)
		pageCount += 1
	}
	check(pager.Error())
}

func main() {
	Sample_QueryTable()
}
