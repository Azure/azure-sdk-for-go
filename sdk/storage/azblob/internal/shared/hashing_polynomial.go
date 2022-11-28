//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"hash/crc64"
)

const CRC64Polynomial uint64 = 0x9A6C9329AC4BC9B5

var CRC64Table = crc64.MakeTable(CRC64Polynomial)
