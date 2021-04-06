// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"

	chk "gopkg.in/check.v1"
)

type tablesRecordedTests struct{}

type testContext struct {
	recording *testframework.Recording
	client    *TableServiceClient
}

const (
	storageAccountNameEnvVar    = "TABLES_STORAGE_ACCOUNT_NAME"
	cosmosAccountNameEnnVar     = "TABLES_COSMOS_ACCOUNT_NAME"
	storageEndpointSuffixEnvVar = "STORAGE_ENDPOINT_SUFFIX"
	cosmosEndpointSuffixEnvVar  = "COSMOS_TABLES_ENDPOINT_SUFFIX"
	storageAccountKeyEnvVar     = "TABLES_PRIMARY_STORAGE_ACCOUNT_KEY"
	cosmosAccountKeyEnvVar      = "TABLES_PRIMARY_COSMOS_ACCOUNT_KEY"
	tableNamePrefix             = "gotable"
	DefaultStorageSuffix        = "core.windows.net"
	DefaultCosmosSuffix         = "cosmos.azure.com"
)

type EndpointType string

const (
	StorageEndpoint EndpointType = "storage"
	CosmosEndpoint  EndpointType = "cosmos"
)

var ctx = context.Background()
var clientsMap map[string]*testContext = make(map[string]*testContext)

func storageURI(accountName string, endpointSuffix string) string {
	return "https://" + accountName + ".table." + endpointSuffix
}

func cosmosURI(accountName string, endpointSuffix string) string {
	return "https://" + accountName + ".table." + endpointSuffix
}

// create the test specific TableClient and wire it up to recordings
func recordedTestSetup(c *chk.C, testName string, endpointType EndpointType, mode testframework.RecordMode) {
	var accountName string
	var suffix string
	var cred *SharedKeyCredential
	var uri string

	// init the test framework
	context := testframework.NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
	recording, err := testframework.NewRecording(context, mode, string(endpointType))
	c.Assert(err, chk.IsNil)

	if endpointType == StorageEndpoint {
		accountName = recording.GetRecordedVariable(storageAccountNameEnvVar)
		suffix = recording.GetOptionalRecordedVariable(storageEndpointSuffixEnvVar, DefaultStorageSuffix)
		cred, _ = NewSharedKeyCredential(accountName, recording.GetRecordedVariable(storageAccountKeyEnvVar, testframework.Secret_Base64String))
		uri = storageURI(accountName, suffix)
	} else {
		accountName = recording.GetRecordedVariable(cosmosAccountNameEnnVar)
		suffix = recording.GetOptionalRecordedVariable(cosmosEndpointSuffixEnvVar, DefaultCosmosSuffix)
		cred, _ = NewSharedKeyCredential(accountName, recording.GetRecordedVariable(cosmosAccountKeyEnvVar, testframework.Secret_Base64String))
		uri = cosmosURI(accountName, suffix)
	}

	client, err := NewTableServiceClient(uri, cred, &TableClientOptions{HTTPClient: recording, Retry: azcore.RetryOptions{MaxRetries: -1}})
	c.Assert(err, chk.IsNil)
	clientsMap[testName] = &testContext{client: client, recording: recording}
}

func recordedTestTeardown(key string) {
	clientsMap[key].recording.Stop()
}

// cleans up the specified tables. If tables is nil, all tables will be deleted
func cleanupTables(context *testContext, tables *[]string) {
	c := context.client
	if tables == nil {
		pager := c.QueryTables(QueryOptions{})
		for pager.NextPage(ctx) {
			for _, t := range *(pager.PageResponse().TableQueryResponse.Value) {
				c.Delete(ctx, *t.TableName)
			}
		}
	} else {
		for _, t := range *tables {
			c.Delete(ctx, t)
		}
	}
}

func getTestContext(key string) *testContext {
	return clientsMap[key]
}

func getTableName(context *testContext, prefix ...string) string {
	if len(prefix) == 0 {
		return context.recording.GenerateAlphaNumericId(tableNamePrefix, 20, true)
	} else {
		return context.recording.GenerateAlphaNumericId(prefix[0], 20, true)
	}
}

func testKey(c *chk.C, ep EndpointType) string {
	return c.TestName() + string(ep)
}

func createSimpleEntities(count int, pk string) *[]map[string]interface{} {
	result := make([]map[string]interface{}, count)

	for i := 1; i <= count; i++ {
		var e = map[string]interface{}{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(i),
			"StringProp": fmt.Sprintf("some string %d", i),
			"IntProp":    i,
			"BoolProp":   true,
		}
		result[i-1] = e
	}
	return &result
}
