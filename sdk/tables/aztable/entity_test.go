// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
)

// type edmEntity struct {
// 	Entity
// 	BigInt   EdmInt64    `json:"BigInt"`
// 	Guid     EdmGuid     `json:"Guid"`
// 	Time     EdmDateTime `json:"Time"`
// 	SmallInt int         `json:"SmallInt"`
// 	Bool     bool        `json:"Bool"`
// 	Bytes    []byte      `json:"Bytes"`
// }

// "odata.type":"account.Customers",
// "odata.id":https://myaccount.table.core.windows.net/Customers(PartitionKey='Customer03',RowKey='Name'),
// "odata.etag":"W/\"0x5B168C7B6E589D2\"",
// "odata.editlink":"Customers(PartitionKey='Customer03',RowKey='Name')",
// "PartitionKey":"partitionkey",
// "RowKey":"rowkey",
// "Timestamp":"2013-08-09T18:55:48.3402073Z",
// "Bool":false,
// "Int32":1234,
// Int64@odata.type:"Edm.Int64",
// "Int64":"123456789012",
// "Double":1234.1234,
// "String":"test",
// Guid@odata.type:"Edm.Guid",
// "Guid":"4185404a-5818-48c3-b9be-f217df0dba6f",
// DateTime@odata.type:"Edm.DateTime",
// "DateTime":"2013-08-02T17:37:43.9004348Z",
// Binary@odata.type:"Edm.Binary",
// "Binary":"AQIDBA=="

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

func (s *tableClientLiveTests) TestEdmMarshalling() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	edmEntity := createEdmEntity(1, "partition")

	marshalled, err := json.Marshal(edmEntity)
	require.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	require.Nil(err)

	fullMetadata := &QueryOptions{
		Format: OdataMetadataFormatApplicationJSONOdataFullmetadata.ToPtr(),
	}

	resp, err := client.GetEntity(ctx, "partition", fmt.Sprint(1), fullMetadata)
	require.Nil(err)
	var receivedEntity EdmEntity
	err = json.Unmarshal(resp.Value, &receivedEntity)
	require.Nil(err)

	require.Equal(edmEntity.PartitionKey, receivedEntity.PartitionKey)
	require.Equal(edmEntity.RowKey, receivedEntity.RowKey)
	require.Equal(edmEntity.Properties["Bool"], receivedEntity.Properties["Bool"])
	require.Equal(edmEntity.Properties["Int32"], receivedEntity.Properties["Int32"])
	fmt.Println(receivedEntity.Properties["Int64"])
	// require.Equal(edmEntity.Properties["Int64"], receivedEntity.Properties["Int64"])
	require.Equal(edmEntity.Properties["Double"], receivedEntity.Properties["Double"])
	require.Equal(edmEntity.Properties["String"], receivedEntity.Properties["String"])
	require.Equal(edmEntity.Properties["Guid"], receivedEntity.Properties["Guid"])
	require.Equal(edmEntity.Properties["DateTime"], receivedEntity.Properties["DateTime"])
	require.Equal(edmEntity.Properties["Binary"], receivedEntity.Properties["Binary"])
}
