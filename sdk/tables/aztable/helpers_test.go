// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// type testEntity struct {
// 	Entity
// 	BigInt            uint64    `json:"BigInt"`
// 	BigIntType        string    `json:"BigInt@odata.type"`
// 	SmallInt          int32     `json:"SmallInt"`
// 	BigFloat          float64   `json:"BigFloat"`
// 	SmallFloat        float32   `json:"SmallFloat"`
// 	BoolType          bool      `json:"BoolType"`
// 	BinaryType        []byte    `json:"BinaryType"`
// 	StringType        string    `json:"StringType"`
// 	DateTimeType      time.Time `json:"DateTimeType"`
// 	DateTimeTypeOdata string    `json:"DateTimeType@odata.type"`
// }

// func TestAddOdataAnnotations(t *testing.T) {
// 	testEnt := testEntity{
// 		Entity: Entity{
// 			PartitionKey: "pk",
// 			RowKey:       "rk",
// 		},
// 		BigInt:            18446744073709551615,
// 		BigIntType:        "Edm.Int64",
// 		SmallInt:          10,
// 		BigFloat:          math.Pow(2, 33) + 1.1,
// 		SmallFloat:        3.14159,
// 		BoolType:          true,
// 		BinaryType:        []byte("BinaryType"),
// 		StringType:        "StringType",
// 		DateTimeType:      time.Now(),
// 		DateTimeTypeOdata: "Edm.DateTime",
// 	}
// 	fmt.Println(testEnt.BigInt)
// 	copyEntity := testEntity{
// 		Entity: Entity{
// 			PartitionKey: "pk",
// 			RowKey:       "rk",
// 		},
// 		BigInt:            18446744073709551615,
// 		BigIntType:        "Edm.Int64",
// 		SmallInt:          10,
// 		BigFloat:          math.Pow(2, 33) + 1.1,
// 		SmallFloat:        3.14159,
// 		BoolType:          true,
// 		BinaryType:        []byte("BinaryType"),
// 		StringType:        "StringType",
// 		DateTimeType:      time.Now(),
// 		DateTimeTypeOdata: "Edm.DateTime",
// 	}
// 	marshalled, err := json.Marshal(testEnt)
// 	assert.Nil(t, err)
// 	var mapEntity map[string]interface{}
// 	err = json.Unmarshal(marshalled, &mapEntity)

// 	var copyEntityMap map[string]interface{}
// 	marshalledCopy, err := json.Marshal(copyEntity)
// 	assert.Nil(t, err)

// 	err = json.Unmarshal(marshalledCopy, &copyEntityMap)

// 	addOdataAnnotations(&mapEntity)

// 	assert.Equal(t, len(mapEntity), 2*len(copyEntityMap)-2)

// 	for k, v := range mapEntity {
// 		fmt.Println(k, v)
// 	}

// 	toOdataAnnotatedDictionary(&copyEntityMap)
// 	fmt.Println("\n\nSecond")
// 	for k, v := range copyEntityMap {
// 		fmt.Println(k, v)
// 	}
// 	// odata, ok := mapEntity["BigInt@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmInt64)

// 	// odata, ok = mapEntity["BigFloat@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmDouble)

// 	// odata, ok = mapEntity["SmallFloat@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmDouble)

// 	// odata, ok = mapEntity["BoolType@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmBoolean)

// 	// odata, ok = mapEntity["BinaryType@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmBinary)

// 	// odata, ok = mapEntity["StringType@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmString)

// 	// odata, ok = mapEntity["DateTimeType@odata.type"]
// 	// assert.True(t, ok)
// 	// assert.Equal(t, odata, edmDateTime)

// }
