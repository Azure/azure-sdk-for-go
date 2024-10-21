// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type endpointType string

const (
	storageEndpoint                endpointType = "storage"
	storageTokenCredentialEndpoint endpointType = "storage_tc"
	cosmosEndpoint                 endpointType = "cosmos"
	cosmosTokenCredentialEndpoint  endpointType = "cosmos_tc"
)

var ctx = context.Background()

func storageURI(accountName string) string {
	return fmt.Sprintf("https://%v.table.core.windows.net/", accountName)
}

func cosmosURI(accountName string) string {
	return fmt.Sprintf("https://%v.table.cosmos.azure.com/", accountName)
}

func insertNEntities(pk string, n int, client *Client) error {
	for i := 0; i < n; i++ {
		e := &map[string]any{
			"PartitionKey": pk,
			"RowKey":       fmt.Sprint(i),
			"Value":        i + 1,
		}
		marshalled, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = client.AddEntity(ctx, marshalled, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

type basicTestEntity struct {
	Entity
	Integer int32
	String  string
	Bool    bool
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
	return createSimpleEntityWithRowKey(count, pk, fmt.Sprint(count))
}

func createSimpleEntityWithRowKey(count int, pk string, rk string) basicTestEntity {
	return basicTestEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       rk,
		},
		String:  fmt.Sprintf("some string %d", count),
		Integer: int32(count),
		Bool:    true,
	}
}

// Use this for a replaced entity to assert a property (Bool) is removed
func createSimpleEntityNoBool(count int, pk string) map[string]any {
	m := make(map[string]any)
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
			PartitionKey: pk,
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

func createComplexEntities(count int, pk string) []complexTestEntity {
	result := make([]complexTestEntity, count)

	for i := 1; i <= count; i++ {
		result[i-1] = createComplexEntity(i, pk)
	}
	return result
}

func createEdmEntity(count int, pk string) EDMEntity {
	return EDMEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(count),
		},
		Properties: map[string]any{
			"Bool":     false,
			"Int32":    int32(1234),
			"Int64":    EDMInt64(123456789012),
			"Double":   1234.1234,
			"String":   "test",
			"Guid":     EDMGUID("4185404a-5818-48c3-b9be-f217df0dba6f"),
			"DateTime": EDMDateTime(time.Date(2013, time.August, 02, 17, 37, 43, 9004348, time.UTC)),
			"Binary":   EDMBinary("SomeBinary"),
		},
	}
}

func requireSameDateTime(t *testing.T, time1, time2 any) {
	t1 := time.Time(time1.(EDMDateTime))
	t2 := time.Time(time2.(EDMDateTime))
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
