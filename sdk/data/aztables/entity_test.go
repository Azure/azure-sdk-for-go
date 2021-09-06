// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBasicEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
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
			require.Nil(t, err)
			_, err = client.AddEntity(ctx, marshalled, nil)
			require.Nil(t, err)

			resp, err := client.GetEntity(ctx, "pk001", "rk001", nil)
			require.Nil(t, err)

			receivedEntity := basicTestEntity{}
			err = json.Unmarshal(resp.Value, &receivedEntity)
			require.Nil(t, err)
			require.Equal(t, receivedEntity.PartitionKey, "pk001")
			require.Equal(t, receivedEntity.RowKey, "rk001")

			queryString := "PartitionKey eq 'pk001'"
			listOptions := ListEntitiesOptions{Filter: &queryString}
			pager := client.List(&listOptions)
			count := 0
			for pager.NextPage(ctx) {
				resp := pager.PageResponse()
				for _, e := range resp.Entities {
					err = json.Unmarshal(e, &receivedEntity)
					require.NoError(t, err)
					require.Equal(t, receivedEntity.PartitionKey, "pk001")
					require.Equal(t, receivedEntity.RowKey, "rk001")
					count += 1
				}
			}
			require.Equal(t, count, 1)
		})
	}
}

func TestEdmMarshalling(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true)
			defer delete()

			edmEntity := createEdmEntity(1, "partition")

			marshalled, err := json.Marshal(edmEntity)
			require.Nil(t, err)
			_, err = client.AddEntity(ctx, marshalled, nil)
			require.Nil(t, err)

			resp, err := client.GetEntity(ctx, "partition", fmt.Sprint(1), nil)
			require.Nil(t, err)
			var receivedEntity EDMEntity
			err = json.Unmarshal(resp.Value, &receivedEntity)
			require.Nil(t, err)

			require.Equal(t, edmEntity.PartitionKey, receivedEntity.PartitionKey)
			require.Equal(t, edmEntity.RowKey, receivedEntity.RowKey)
			require.Equal(t, edmEntity.Properties["Bool"], receivedEntity.Properties["Bool"])
			require.Equal(t, edmEntity.Properties["Int32"], receivedEntity.Properties["Int32"])
			require.Equal(t, edmEntity.Properties["Int64"], receivedEntity.Properties["Int64"])
			require.Equal(t, edmEntity.Properties["Double"], receivedEntity.Properties["Double"])
			require.Equal(t, edmEntity.Properties["String"], receivedEntity.Properties["String"])
			require.Equal(t, edmEntity.Properties["Guid"], receivedEntity.Properties["Guid"])
			require.Equal(t, edmEntity.Properties["Binary"], receivedEntity.Properties["Binary"])
			requireSameDateTime(t, edmEntity.Properties["DateTime"], receivedEntity.Properties["DateTime"])

			// Unmarshal to raw json
			var received2 map[string]json.RawMessage
			err = json.Unmarshal(resp.Value, &received2)
			require.Nil(t, err)

			// Unmarshal to plain map
			var received3 map[string]interface{}
			err = json.Unmarshal(resp.Value, &received3)
			require.Nil(t, err)

		})
	}
}
