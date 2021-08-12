// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
	storageEndpoint EndpointType = "storage"
	cosmosEndpoint  EndpointType = "cosmos"
)

var ctx = context.Background()

func storageURI(accountName string, endpointSuffix string) string {
	return fmt.Sprintf("https://%v.table.%v/", accountName, endpointSuffix)
}

func cosmosURI(accountName string, endpointSuffix string) string {
	return fmt.Sprintf("https://%v.table.%v/", accountName, endpointSuffix)
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
		pager := c.ListTables(nil)
		for pager.NextPage(ctx) {
			for _, t := range pager.PageResponse().TableListResponse.Value {
				_, err := c.DeleteTable(ctx, *t.TableName, nil)
				if err != nil {
					fmt.Printf("Error cleaning up tables. %v\n", err.Error())
				}
			}
		}
	} else {
		for _, t := range *tables {
			_, err := c.DeleteTable(ctx, t, nil)
			if err != nil {
				fmt.Printf("There was an error cleaning up tests. %v\n", err.Error())
			}
		}
	}
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
	}
}

func createComplexEntities(count int, pk string) *[]complexTestEntity {
	result := make([]complexTestEntity, count)

	for i := 1; i <= count; i++ {
		result[i-1] = createComplexEntity(i, pk)
	}
	return &result
}

func createEdmEntity(count int, pk string) EdmEntity {
	return EdmEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(count),
		},
		Properties: map[string]interface{}{
			"Bool":     false,
			"Int32":    int32(1234),
			"Int64":    EdmInt64(123456789012),
			"Double":   1234.1234,
			"String":   "test",
			"Guid":     EdmGuid("4185404a-5818-48c3-b9be-f217df0dba6f"),
			"DateTime": EdmDateTime(time.Date(2013, time.August, 02, 17, 37, 43, 9004348, time.UTC)),
			"Binary":   EdmBinary("SomeBinary"),
		},
	}
}

func requireSameDateTime(t *testing.T, time1, time2 interface{}) {
	t1 := time.Time(time1.(EdmDateTime))
	t2 := time.Time(time2.(EdmDateTime))
	require.Equal(t, t1.Year(), t2.Year())
	require.Equal(t, t1.Month(), t2.Month())
	require.Equal(t, t1.Day(), t2.Day())
	require.Equal(t, t1.Hour(), t2.Hour())
	require.Equal(t, t1.Minute(), t2.Minute())
	require.Equal(t, t1.Second(), t2.Second())
	z1, _ := t1.Zone()
	z2, _ := t2.Zone()
	require.Equal(t, z1, z2)
}
