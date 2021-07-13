package aztable

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type basicTestEntity struct {
	Entity
	Integer int32
	String  string
	Float   float64
}

func (s *tableClientLiveTests) TestAddBasicEntity() {
	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	basicEntity := basicTestEntity{
		Entity: Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Integer: 10,
		String:  "abcdef",
		Float:   3.14159,
	}

	marshalled, err := json.Marshal(basicEntity)
	assert.Nil(err)
	_, err = client.AddEntity(ctx, marshalled)
	assert.Nil(err)

	resp, err := client.GetEntity(ctx, "pk001", "rk001")
	assert.Nil(err)

	newEntity := resp.Value
	pk, ok := newEntity[partitionKey]
	assert.True(ok)
	assert.Equal(pk, "pk001")

	rk, ok := newEntity[rowKey]
	assert.True(ok)
	assert.Equal(rk, "rk001")

	queryString := "PartitionKey eq 'pk001'"
	queryOptions := QueryOptions{Filter: &queryString}
	pager := client.Query(queryOptions)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		// model := basicTestEntity{}
		for _, e := range resp.TableEntityQueryResponse.Value {
			pk, ok := e[partitionKey]
			assert.True(ok)
			assert.Equal(pk, "pk001")
		}

	}
}
