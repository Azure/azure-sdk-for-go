// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestAddBasicEntity(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
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
			pager := client.NewListEntitiesPager(&listOptions)
			count := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
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
			client, delete := initClientTest(t, service, true, tracing.Provider{})
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
			var received3 map[string]any
			err = json.Unmarshal(resp.Value, &received3)
			require.Nil(t, err)

		})
	}
}
func TestEntityQuotes(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			pk, err := createRandomName(t, "partition")
			require.NoError(t, err)

			edmEntity := EDMEntity{
				Entity: Entity{
					PartitionKey: pk,
					RowKey:       fmt.Sprint(1),
				},
				Properties: map[string]any{
					"SingleQuote":           "''",
					"DoubleQuote":           "\"\"",
					"JustSpaces":            "    ",
					"LeadingSpaces":         "   abc",
					"TrailingSpaces":        "abc     ",
					"LeadingTrailingSpaces": "    abc    ",
				},
			}

			marshalled, err := json.Marshal(edmEntity)
			require.Nil(t, err)
			_, err = client.AddEntity(ctx, marshalled, nil)
			require.Nil(t, err)

			resp, err := client.GetEntity(ctx, edmEntity.PartitionKey, edmEntity.RowKey, nil)
			require.Nil(t, err)
			var receivedEntity EDMEntity
			err = json.Unmarshal(resp.Value, &receivedEntity)
			require.Nil(t, err)

			require.Equal(t, edmEntity.PartitionKey, receivedEntity.PartitionKey)
			require.Equal(t, edmEntity.RowKey, receivedEntity.RowKey)
			require.Equal(t, edmEntity.Properties["SingleQuote"], receivedEntity.Properties["SingleQuote"])
			require.Equal(t, edmEntity.Properties["DoubleQuote"], receivedEntity.Properties["DoubleQuote"])
			require.Equal(t, edmEntity.Properties["JustSpaces"], receivedEntity.Properties["JustSpaces"])
			require.Equal(t, edmEntity.Properties["LeadingSpaces"], receivedEntity.Properties["LeadingSpaces"])
			require.Equal(t, edmEntity.Properties["TrailingSpaces"], receivedEntity.Properties["TrailingSpaces"])
			require.Equal(t, edmEntity.Properties["LeadingTrailingSpaces"], receivedEntity.Properties["LeadingTrailingSpaces"])

			// Unmarshal to raw json
			var received2 map[string]json.RawMessage
			err = json.Unmarshal(resp.Value, &received2)
			require.Nil(t, err)

			// Unmarshal to plain map
			var received3 map[string]any
			err = json.Unmarshal(resp.Value, &received3)
			require.Nil(t, err)
		})
	}
}

func TestEntityUnicode(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, delete := initClientTest(t, service, true, tracing.Provider{})
			defer delete()

			pk, err := createRandomName(t, "partition")
			require.NoError(t, err)

			edmEntity := EDMEntity{
				Entity: Entity{
					PartitionKey: pk,
					RowKey:       fmt.Sprint(1),
				},
				Properties: map[string]any{
					"Unicode": "ꀕ",
					"ꀕ":       "Unicode",
				},
			}

			marshalled, err := json.Marshal(edmEntity)
			require.Nil(t, err)
			_, err = client.AddEntity(ctx, marshalled, nil)
			require.Nil(t, err)

			resp, err := client.GetEntity(ctx, edmEntity.PartitionKey, edmEntity.RowKey, nil)
			require.Nil(t, err)
			var receivedEntity EDMEntity
			err = json.Unmarshal(resp.Value, &receivedEntity)
			require.Nil(t, err)

			require.Equal(t, edmEntity.PartitionKey, receivedEntity.PartitionKey)
			require.Equal(t, edmEntity.RowKey, receivedEntity.RowKey)
			require.Equal(t, edmEntity.Properties["Unicode"], receivedEntity.Properties["Unicode"])
			require.Equal(t, edmEntity.Properties["ꀕ"], receivedEntity.Properties["ꀕ"])

			// Unmarshal to raw json
			var received2 map[string]json.RawMessage
			err = json.Unmarshal(resp.Value, &received2)
			require.Nil(t, err)

			// Unmarshal to plain map
			var received3 map[string]any
			err = json.Unmarshal(resp.Value, &received3)
			require.Nil(t, err)
		})
	}
}

func TestPrepareKey(t *testing.T) {
	require.EqualValues(t, "unchanged", prepareKey("unchanged"))
	require.EqualValues(t, "sin''gle", prepareKey("sin'gle"))
	require.EqualValues(t, "''beginning", prepareKey("'beginning"))
	require.EqualValues(t, "end''", prepareKey("end'"))
	require.EqualValues(t, "''quoted''", prepareKey("'quoted'"))
	require.EqualValues(t, "d''''ouble", prepareKey("d''ouble"))
	require.EqualValues(t, "''", prepareKey("'"))
	require.EqualValues(t, "", prepareKey(""))
}
