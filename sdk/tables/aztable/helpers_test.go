// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEntity struct {
	PartitionKey string
	RowKey       string
	BigInt       uint64
	SmallInt     int32
	BigFloat     float64
	SmallFloat   float32
	BoolType     bool
	BinaryType   []byte
	StringType   string
	DateTimeType time.Time
}

func TestAddOdataAnnotations(t *testing.T) {
	testEnt := testEntity{
		PartitionKey: "pk",
		RowKey:       "rk",
		BigInt:       18446744073709551615,
		SmallInt:     10,
		BigFloat:     math.Pow(2, 33) + 1.1,
		SmallFloat:   3.14159,
		BoolType:     true,
		BinaryType:   []byte("BinaryType"),
		StringType:   "StringType",
		DateTimeType: time.Now(),
	}
	fmt.Println(testEnt.BigInt)
	copyEntity := testEntity{
		PartitionKey: "pk",
		RowKey:       "rk",
		BigInt:       18446744073709551615,
		SmallInt:     10,
		BigFloat:     math.Pow(2, 33) + 1.1,
		SmallFloat:   3.14159,
		BoolType:     true,
		BinaryType:   []byte("BinaryType"),
		StringType:   "StringType",
		DateTimeType: time.Now(),
	}
	marshalled, err := json.Marshal(testEnt)
	assert.Nil(t, err)
	var mapEntity map[string]interface{}
	err = json.Unmarshal(marshalled, &mapEntity)

	var copyEntityMap map[string]interface{}
	marshalledCopy, err := json.Marshal(copyEntity)
	assert.Nil(t, err)

	err = json.Unmarshal(marshalledCopy, &copyEntityMap)

	addOdataAnnotations(&mapEntity)

	assert.Equal(t, len(mapEntity), 2*len(copyEntityMap)-2)

	for k, v := range mapEntity {
		fmt.Println(k, v)
	}

	toOdataAnnotatedDictionary(&copyEntityMap)
	fmt.Println("\n\nSecond")
	for k, v := range copyEntityMap {
		fmt.Println(k, v)
	}
	// odata, ok := mapEntity["BigInt@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmInt64)

	// odata, ok = mapEntity["BigFloat@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmDouble)

	// odata, ok = mapEntity["SmallFloat@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmDouble)

	// odata, ok = mapEntity["BoolType@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmBoolean)

	// odata, ok = mapEntity["BinaryType@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmBinary)

	// odata, ok = mapEntity["StringType@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmString)

	// odata, ok = mapEntity["DateTimeType@odata.type"]
	// assert.True(t, ok)
	// assert.Equal(t, odata, edmDateTime)

}
