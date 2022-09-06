// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"CreateEntityTest":      {Register: insertTestRegister, New: NewInsertEntityTest},
		"ListEntitiesTest":      {Register: listTestRegister, New: NewListEntitiesTest},
		"CreateEntityBatchTest": {Register: batchTestRegister, New: NewBatchTest},
	})
}
