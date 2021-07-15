// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const odataSuffix = "@odata.type"

func addOdataAnnotations(entity *map[string]interface{}) {
	for key, value := range *entity {
		if key == partitionKey || key == rowKey || strings.Contains(key, odataSuffix) {
			continue
		}

		switch value.(type) {
		case int64, int32, float64:
			e := getNumberType(value)
			addOdataType(entity, key, e)
		// case int32:
		// 	addOdataType(entity, key, edmInt32)
		// case float64:
		// 	addOdataType(entity, key, edmDouble)
		case []byte:
			addOdataType(entity, key, edmBinary)
		case bool:
			addOdataType(entity, key, edmBoolean)
		case time.Time:
			addOdataType(entity, key, edmDateTime)
		default:
			addOdataType(entity, key, edmString)
		}
	}
}

func addOdataType(entity *map[string]interface{}, key, edm string) {
	keyOdata := key + odataSuffix
	(*entity)[keyOdata] = edm
}

func getNumberType(value interface{}) string {
	// First convert interface to string
	strValue := fmt.Sprintf("%v", value)
	fmt.Println(strValue)

	v, err := strconv.ParseInt(strValue, 10, 64)
	if err == nil {
		fmt.Println("V: ", v)
		if v > int64(math.Pow(2, 32)) {
			return edmInt64
		}
		return edmInt32
	} else {
		fmt.Println("ERROR: ", err, value)
	}

	_, err = strconv.ParseFloat(strValue, 64)
	if err == nil {
		return edmDouble
	}

	// Don't know what to do with it so say it is a string
	return edmString
}
