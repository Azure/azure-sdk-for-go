// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "time"

type Entity struct {
	PartitionKey string
	RowKey       string
	TimeStamp    time.Time
}
