// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
)

type edmEntity struct {
	Entity
	BigInt EdmInt64    `json:"BigInt"`
	Guid   EdmGuid     `json:"Guid"`
	Time   EdmDateTime `json:"Time"`
}

func createEdmEntity(count int, pk string) edmEntity {
	return edmEntity{
		Entity: Entity{
			PartitionKey: pk,
			RowKey:       fmt.Sprint(count),
		},
		BigInt: 1125899906842624,
		Guid:   "abcd-efgh-ijkl",
		Time:   EdmDateTime{time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC)},
	}
}

func (s *tableClientLiveTests) TestEdmMarshalling() {
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	EdmEntity := createEdmEntity(1, "partition")

	marshalled, err := json.Marshal(EdmEntity)
	require.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	require.Nil(err)

	fullMetadata := &QueryOptions{
		Format: OdataMetadataFormatApplicationJSONOdataFullmetadata.ToPtr(),
	}

	resp, err := client.GetEntity(ctx, "partition", fmt.Sprint(1), fullMetadata)
	require.Nil(err)
	receivedEntity := edmEntity{}
	err = json.Unmarshal(resp.Value, &receivedEntity)
	require.Nil(err)

	fmt.Println(receivedEntity)

	require.Equal(receivedEntity.BigInt, EdmEntity.BigInt)
	require.Equal(receivedEntity.Guid, EdmEntity.Guid)
	require.Equal(receivedEntity.Time, EdmEntity.Time)
	require.Equal(receivedEntity.PartitionKey, EdmEntity.PartitionKey)
	require.Equal(receivedEntity.RowKey, EdmEntity.RowKey)

}
