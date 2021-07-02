// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"time"

	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	Entity
	BasicInt      int32
	LargeInt      int64
	StringValue   string
	DateTimeValue time.Time
}

func failOnError(err error, s *tableClientLiveTests) {
	if err != nil {
		assert.FailNow(s.T(), err.Error())
	}
}

// func (s *tableClientLiveTests) TestCustomEntity() {
// 	client, delete := s.init(true)
// 	defer delete()

// 	// Create a TestEntity
// 	testEntity := TestEntity{Entity: Entity{PartitionKey: "pk001", RowKey: "rk001"}, BasicInt: 10, LargeInt: int64(math.Pow(2, 34)), StringValue: "basicString", DateTimeValue: time.Now()}

// 	_, err := client.AddEntity(ctx, testEntity)
// 	failOnError(err, s)

// 	receivedEntity, err := client.GetEntity(ctx, "pk001", "rk001")
// 	assert.Nil(s.T(), err)
// 	assert.Equal(s.T(), testEntity.PartitionKey, receivedEntity.Value["PartitionKey"].(string))
// }
