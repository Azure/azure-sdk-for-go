//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

// for testing, create struct with all the data types
// remove computer field, not supported anymore
type ComputerInfo struct {
	InputTime         time.Time
	Computer          string
	AdditionalContext int
}

// test for greater than 1 mb
// generate a file and read from it

func TestUpload(t *testing.T) {
	client := startTest(t)

	var data []ComputerInfo
	for i := 0; i < 10; i++ {
		data = append(data, ComputerInfo{
			InputTime:         time.Now().UTC(),
			Computer:          "Computer" + strconv.Itoa(i),
			AdditionalContext: i,
		})
	}
	data2, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	_, err = client.Upload(context.Background(), ruleID, stream, data2, nil)
	if err != nil {
		panic(err)
	}

}
