// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/require"
)

func (s *tableClientLiveTests) TestAddBasicEntity() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	basicEntity := basicTestEntity{
		Entity: Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Integer: 10,
		String:  "abcdef",
		Bool:    true,
	}

	marshalled, err := json.Marshal(basicEntity)
	require.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	require.Nil(err)

	resp, err := client.GetEntity(ctx, "pk001", "rk001", nil)
	require.Nil(err)

	receivedEntity := basicTestEntity{}
	err = json.Unmarshal(resp.Value, &receivedEntity)
	require.Nil(err)
	require.Equal(receivedEntity.PartitionKey, "pk001")
	require.Equal(receivedEntity.RowKey, "rk001")

	queryString := "PartitionKey eq 'pk001'"
	listOptions := ListOptions{Filter: &queryString}
	pager := client.List(&listOptions)
	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, e := range resp.TableEntityQueryResponse.Value {
			err = json.Unmarshal(e, &receivedEntity)
			require.NoError(err)
			require.Equal(receivedEntity.PartitionKey, "pk001")
			require.Equal(receivedEntity.RowKey, "rk001")
			count += 1
		}
	}

	require.Equal(count, 1)
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
	require.Equal(edmEntity.Properties["Int64"], receivedEntity.Properties["Int64"])
	require.Equal(edmEntity.Properties["Double"], receivedEntity.Properties["Double"])
	require.Equal(edmEntity.Properties["String"], receivedEntity.Properties["String"])
	require.Equal(edmEntity.Properties["Guid"], receivedEntity.Properties["Guid"])
	require.Equal(edmEntity.Properties["Binary"], receivedEntity.Properties["Binary"])
	requireSameDateTime(require, edmEntity.Properties["DateTime"], receivedEntity.Properties["DateTime"])

	// Unmarshal to raw json
	var received2 map[string]json.RawMessage
	err = json.Unmarshal(resp.Value, &received2)
	require.Nil(err)

	// Unmarshal to plain map
	var received3 map[string]interface{}
	err = json.Unmarshal(resp.Value, &received3)
	require.Nil(err)
}
