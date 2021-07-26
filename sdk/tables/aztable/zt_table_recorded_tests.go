// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	recording *recording.Recording
	client    *TableServiceClient
	context   *recording.TestContext
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
var cosmosTestsMap map[string]bool = make(map[string]bool)

func storageURI(accountName string, endpointSuffix string) string {
	return fmt.Sprintf("https://%v.table.%v/", accountName, endpointSuffix)
}

func cosmosURI(accountName string, endpointSuffix string) string {
	return fmt.Sprintf("https://%v.table.%v/", accountName, endpointSuffix)
}

// create the test specific TableClient and wire it up to recordings
func recordedTestSetup(t *testing.T, testName string, endpointType EndpointType, mode recording.RecordMode) {
	var accountName string
	var suffix string
	var cred *SharedKeyCredential
	var secret string
	var uri string
	require := require.New(t)

	// init the test framework
	context := recording.NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { t.Log(msg) }, func() string { return testName })
	r, err := recording.NewRecording(context, mode)
	require.NoError(err)

	if endpointType == StorageEndpoint {
		accountName, err = r.GetRecordedVariable(storageAccountNameEnvVar, recording.Default)
		require.NoError(err)
		suffix = r.GetOptionalRecordedVariable(storageEndpointSuffixEnvVar, DefaultStorageSuffix, recording.Default)
		secret, err = r.GetRecordedVariable(storageAccountKeyEnvVar, recording.Secret_Base64String)
		require.NoError(err)
		cred, err = NewSharedKeyCredential(accountName, secret)
		require.NoError(err)
		uri = storageURI(accountName, suffix)
	} else {
		accountName, err = r.GetRecordedVariable(cosmosAccountNameEnnVar, recording.Default)
		require.NoError(err)
		suffix = r.GetOptionalRecordedVariable(cosmosEndpointSuffixEnvVar, DefaultCosmosSuffix, recording.Default)
		secret, err = r.GetRecordedVariable(cosmosAccountKeyEnvVar, recording.Secret_Base64String)
		require.NoError(err)
		cred, err = NewSharedKeyCredential(accountName, secret)
		require.NoError(err)
		uri = cosmosURI(accountName, suffix)
		cosmosTestsMap[testName] = true
	}

	client, err := NewTableServiceClient(uri, cred, &TableClientOptions{HTTPClient: r, Retry: azcore.RetryOptions{MaxRetries: -1}})
	require.NoError(err)

	clientsMap[testName] = &testContext{client: client, recording: r, context: &context}
}

func recordedTestTeardown(key string) {
	context, ok := clientsMap[key]
	if ok && !(*context.context).IsFailed() {
		err := context.recording.Stop()
		if err != nil {
			fmt.Printf("Error tearing down tests. %v\n", err.Error())
		}
	}
}

func insertNEntities(pk string, n int, client *TableClient) error {
	for i := 0; i < n; i++ {
		e := &map[string]interface{}{
			"PartitionKey": pk,
			"RowKey":       fmt.Sprint(i),
			"Value":        i + 1,
		}
		marshalled, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = client.AddEntity(ctx, marshalled)
		if err != nil {
			return err
		}
	}
	return nil
}

// cleans up the specified tables. If tables is nil, all tables will be deleted
func cleanupTables(context *testContext, tables *[]string) {
	c := context.client
	if tables == nil {
		pager := c.List(nil)
		for pager.NextPage(ctx) {
			for _, t := range pager.PageResponse().TableQueryResponse.Value {
				_, err := c.Delete(ctx, *t.TableName)
				if err != nil {
					fmt.Printf("Error cleaning up tables. %v\n", err.Error())
				}
			}
		}
	} else {
		for _, t := range *tables {
			_, err := c.Delete(ctx, t)
			if err != nil {
				fmt.Printf("There was an error cleaning up tests. %v\n", err.Error())
			}
		}
	}
}

func getTestContext(key string) *testContext {
	return clientsMap[key]
}

func getTableName(context *testContext, prefix ...string) (string, error) {
	if len(prefix) == 0 {
		return context.recording.GenerateAlphaNumericID(tableNamePrefix, 20, true)
	} else {
		return context.recording.GenerateAlphaNumericID(prefix[0], 20, true)
	}
}

type basicTestEntity struct {
	Entity
	Integer int32
	String  string
	Bool    bool
}

func marshalBasicEntity(b basicTestEntity, require *require.Assertions) *[]byte {
	r, e := json.Marshal(b)
	require.NoError(e)
	return &r
}

type complexTestEntity struct {
	Entity
	Integer  int
	String   string
	Bool     bool
	Float    float32
	DateTime time.Time
	Byte     []byte
	// Integer64 int64 // Need to add type hints for ints/floats above 32bits
	// Float64   float64
}

func createSimpleEntity(count int, pk string) basicTestEntity {
	return basicTestEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(count),
		},
		String:  fmt.Sprintf("some string %d", count),
		Integer: int32(count),
		Bool:    true,
	}
}

// Use this for a replaced entity to assert a property (Bool) is removed
func createSimpleEntityNoBool(count int, pk string) map[string]interface{} {
	m := make(map[string]interface{})
	m[partitionKey] = pk
	m[rowKey] = fmt.Sprint(count)
	m["String"] = fmt.Sprintf("some string %d", count)
	m["Integer"] = int32(count)
	return m
}

func createSimpleEntities(count int, pk string) *[]basicTestEntity {
	result := make([]basicTestEntity, count)
	for i := 1; i <= count; i++ {
		result[i-1] = createSimpleEntity(i, pk)
	}
	return &result
}

func createComplexEntity(i int, pk string) complexTestEntity {
	return complexTestEntity{
		Entity: Entity{
			PartitionKey: "partition",
			RowKey:       fmt.Sprint(i),
		},
		Integer:  int(i),
		String:   "someString",
		Bool:     true,
		Float:    3.14159,
		DateTime: time.Date(2021, time.July, 13, 0, 0, 0, 0, time.UTC),
		Byte:     []byte("somebytes"),
		// Integer64: int64(math.Pow(2, 33)),
		// Float64:   math.Pow(2, 33.1),
	}
}

func createComplexEntities(count int, pk string) *[]complexTestEntity {
	result := make([]complexTestEntity, count)

	for i := 1; i <= count; i++ {
		result[i-1] = createComplexEntity(i, pk)
	}
	return &result
}

type simpleEntity struct {
	ETag         string
	PartitionKey string
	RowKey       string
	Timestamp    time.Time
	IntProp      int
	BoolProp     bool
	StringProp   string
}

type complexEntity struct {
	ETag                  string
	PartitionKey          string
	RowKey                string
	Timestamp             time.Time
	SomeBinaryProperty    []byte
	SomeDateProperty      time.Time
	SomeDoubleProperty0   float64
	SomeDoubleProperty1   float64
	SomeGuidProperty      [16]byte `uuid:""`
	SomeInt64Property     int64
	SomeIntProperty       int
	SomeStringProperty    string
	SomePtrStringProperty *string
}
